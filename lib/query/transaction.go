package query

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/file"
	"github.com/mithrandie/csvq/lib/parser"

	"github.com/mithrandie/go-text/color"
)

type Transaction struct {
	Session *Session

	Environment *cmd.Environment
	Palette     *color.Palette
	Flags       *cmd.Flags

	WaitTimeout   time.Duration
	RetryDelay    time.Duration
	FileContainer *file.Container

	cachedViews      ViewMap
	uncommittedViews UncommittedViews

	operationMutex   *sync.Mutex
	viewLoadingMutex *sync.Mutex
	stdinIsLocked    bool

	PreparedStatements PreparedStatementMap

	SelectedViews []*View
	AffectedRows  int

	AutoCommit bool
}

func NewTransaction(ctx context.Context, defaultWaitTimeout time.Duration, retryDelay time.Duration, session *Session) (*Transaction, error) {
	environment, err := cmd.NewEnvironment(ctx, defaultWaitTimeout, retryDelay)
	if err != nil {
		return nil, ConvertLoadConfigurationError(err)
	}
	flags := cmd.NewFlags(environment)

	palette, err := cmd.NewPalette(environment)
	if err != nil {
		return nil, ConvertLoadConfigurationError(err)
	}
	palette.Disable()

	return &Transaction{
		Session:            session,
		Environment:        environment,
		Palette:            palette,
		Flags:              flags,
		WaitTimeout:        file.DefaultWaitTimeout,
		RetryDelay:         file.DefaultRetryDelay,
		FileContainer:      file.NewContainer(),
		cachedViews:        NewViewMap(),
		uncommittedViews:   NewUncommittedViews(),
		operationMutex:     &sync.Mutex{},
		viewLoadingMutex:   &sync.Mutex{},
		stdinIsLocked:      false,
		PreparedStatements: NewPreparedStatementMap(),
		SelectedViews:      nil,
		AffectedRows:       0,
		AutoCommit:         false,
	}, nil
}

func (tx *Transaction) UpdateWaitTimeout(waitTimeout float64, retryDelay time.Duration) {
	d, err := time.ParseDuration(strconv.FormatFloat(waitTimeout, 'f', -1, 64) + "s")
	if err != nil {
		d = file.DefaultWaitTimeout
	}

	tx.WaitTimeout = d
	tx.RetryDelay = retryDelay
	tx.Flags.SetWaitTimeout(waitTimeout)
}

func (tx *Transaction) UseColor(useColor bool) {
	if useColor {
		tx.Palette.Enable()
	} else {
		tx.Palette.Disable()
	}
	tx.Flags.SetColor(useColor)
}

func (tx *Transaction) Commit(ctx context.Context, filter *Filter, expr parser.Expression) error {
	tx.operationMutex.Lock()
	defer tx.operationMutex.Unlock()

	createdFiles, updatedFiles := tx.uncommittedViews.UncommittedFiles()

	createFileInfo := make([]*FileInfo, 0, len(createdFiles))
	updateFileInfo := make([]*FileInfo, 0, len(updatedFiles))

	if 0 < len(createdFiles) {
		for _, fileinfo := range createdFiles {
			view, _ := tx.cachedViews.Get(parser.Identifier{Literal: fileinfo.Path})

			fp, _ := view.FileInfo.Handler.FileForUpdate()
			if err := fp.Truncate(0); err != nil {
				return NewSystemError(err.Error())
			}
			if _, err := fp.Seek(0, io.SeekStart); err != nil {
				return NewSystemError(err.Error())
			}

			_, err := EncodeView(ctx, fp, view, fileinfo, tx)
			if err != nil {
				return NewCommitError(expr, err.Error())
			}
			createFileInfo = append(createFileInfo, view.FileInfo)
		}
	}

	if 0 < len(updatedFiles) {
		for _, fileinfo := range updatedFiles {
			view, _ := tx.cachedViews.Get(parser.Identifier{Literal: fileinfo.Path})

			fp, _ := view.FileInfo.Handler.FileForUpdate()
			if err := fp.Truncate(0); err != nil {
				return NewSystemError(err.Error())
			}
			if _, err := fp.Seek(0, io.SeekStart); err != nil {
				return NewSystemError(err.Error())
			}

			if _, err := EncodeView(ctx, fp, view, fileinfo, tx); err != nil {
				return NewCommitError(expr, err.Error())
			}

			updateFileInfo = append(updateFileInfo, view.FileInfo)
		}
	}

	for _, f := range createFileInfo {
		if err := tx.FileContainer.Commit(f.Handler); err != nil {
			return NewCommitError(expr, err.Error())
		}
		tx.uncommittedViews.Unset(f)
		tx.LogNotice(fmt.Sprintf("Commit: file %q is created.", f.Path), tx.Flags.Quiet)
	}
	for _, f := range updateFileInfo {
		if err := tx.FileContainer.Commit(f.Handler); err != nil {
			return NewCommitError(expr, err.Error())
		}
		tx.uncommittedViews.Unset(f)
		tx.LogNotice(fmt.Sprintf("Commit: file %q is updated.", f.Path), tx.Flags.Quiet)
	}

	msglist := filter.tempViews.Store(tx.Session, tx.uncommittedViews.UncommittedTempViews())
	if 0 < len(msglist) {
		tx.LogNotice(strings.Join(msglist, "\n"), tx.quietForTemporaryViews(expr))
	}
	tx.uncommittedViews.Clean()
	tx.UnlockStdin()
	if err := tx.ReleaseResources(); err != nil {
		return NewCommitError(expr, err.Error())
	}
	return nil
}

func (tx *Transaction) Rollback(filter *Filter, expr parser.Expression) error {
	tx.operationMutex.Lock()
	defer tx.operationMutex.Unlock()

	createdFiles, updatedFiles := tx.uncommittedViews.UncommittedFiles()

	if 0 < len(createdFiles) {
		for _, fileinfo := range createdFiles {
			tx.LogNotice(fmt.Sprintf("Rollback: file %q is deleted.", fileinfo.Path), tx.Flags.Quiet)
		}
	}

	if 0 < len(updatedFiles) {
		for _, fileinfo := range updatedFiles {
			tx.LogNotice(fmt.Sprintf("Rollback: file %q is restored.", fileinfo.Path), tx.Flags.Quiet)
		}
	}

	if filter != nil {
		msglist := filter.tempViews.Restore(tx.uncommittedViews.UncommittedTempViews())
		if 0 < len(msglist) {
			tx.LogNotice(strings.Join(msglist, "\n"), tx.quietForTemporaryViews(expr))
		}
	}
	tx.uncommittedViews.Clean()
	tx.UnlockStdin()
	if err := tx.ReleaseResources(); err != nil {
		return NewRollbackError(expr, err.Error())
	}
	return nil
}

func (tx *Transaction) quietForTemporaryViews(expr parser.Expression) bool {
	return tx.Flags.Quiet || expr == nil
}

func (tx *Transaction) ReleaseResources() error {
	if err := tx.cachedViews.Clean(tx.FileContainer); err != nil {
		return err
	}
	if err := tx.FileContainer.CloseAll(); err != nil {
		return err
	}
	tx.UnlockStdin()
	return nil
}

func (tx *Transaction) ReleaseResourcesWithErrors() error {
	var errs []error
	if err := tx.cachedViews.CleanWithErrors(tx.FileContainer); err != nil {
		errs = append(errs, err.(*file.ForcedUnlockError).Errors...)
	}
	if err := tx.FileContainer.CloseAllWithErrors(); err != nil {
		errs = append(errs, err.(*file.ForcedUnlockError).Errors...)
	}
	tx.UnlockStdin()
	return file.NewForcedUnlockError(errs)
}

func (tx *Transaction) LockStdinContext(ctx context.Context) error {
	tctx, cancel := file.GetTimeoutContext(ctx, tx.WaitTimeout)
	defer cancel()

	err := tx.Session.stdinLocker.LockContext(tctx)
	if err == nil {
		tx.stdinIsLocked = true
	}
	return err
}

func (tx *Transaction) UnlockStdin() {
	if tx.stdinIsLocked {
		tx.stdinIsLocked = false
		_ = tx.Session.stdinLocker.Unlock()
	}
}

func (tx *Transaction) RLockStdinContext(ctx context.Context) error {
	tctx, cancel := file.GetTimeoutContext(ctx, tx.WaitTimeout)
	defer cancel()

	return tx.Session.stdinLocker.RLockContext(tctx)
}

func (tx *Transaction) RUnlockStdin() {
	_ = tx.Session.stdinLocker.RUnlock()
}

func (tx *Transaction) Error(s string) string {
	if tx.Palette != nil {
		return tx.Palette.Render(cmd.ErrorEffect, s)
	}
	return s
}

func (tx *Transaction) Warn(s string) string {
	if tx.Palette != nil {
		return tx.Palette.Render(cmd.WarnEffect, s)
	}
	return s
}

func (tx *Transaction) Notice(s string) string {
	if tx.Palette != nil {
		return tx.Palette.Render(cmd.NoticeEffect, s)
	}
	return s
}

func (tx *Transaction) Log(log string, quiet bool) {
	if !quiet {
		if err := tx.Session.WriteToStdoutWithLineBreak(log); err != nil {
			println(err.Error())
		}
	}
}

func (tx *Transaction) LogNotice(log string, quiet bool) {
	if !quiet {
		if err := tx.Session.WriteToStdoutWithLineBreak(tx.Notice(log)); err != nil {
			println(err.Error())
		}
	}
}

func (tx *Transaction) LogWarn(log string, quiet bool) {
	if !quiet {
		if err := tx.Session.WriteToStdoutWithLineBreak(tx.Warn(log)); err != nil {
			println(err.Error())
		}
	}
}

func (tx *Transaction) LogError(log string) {
	if err := tx.Session.WriteToStderrWithLineBreak(tx.Error(log)); err != nil {
		println(err.Error())
	}
}
