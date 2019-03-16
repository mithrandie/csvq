package query

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/mithrandie/csvq/lib/parser"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/file"
)

type Transaction struct {
	Session *Session

	Environment *cmd.Environment
	Flags       *cmd.Flags

	WaitTimeout   time.Duration
	RetryDelay    time.Duration
	FileContainer *file.Container

	CachedViews      ViewMap
	UncommittedViews *UncommittedViewMap

	SelectedViews []*View
	AffectedRows  int

	AutoCommit bool
}

func NewTransaction(ctx context.Context, defaultWaitTimeout time.Duration, retryDelay time.Duration, session *Session) (*Transaction, error) {
	environment, err := cmd.NewEnvironment(ctx, defaultWaitTimeout, retryDelay)
	if err != nil {
		return nil, NewTransactionError(err.Error())
	}
	flags := cmd.NewFlags(environment)

	if err := cmd.LoadPalette(environment); err != nil {
		return nil, NewTransactionError(err.Error())
	}

	return &Transaction{
		Session:          session,
		Environment:      environment,
		Flags:            flags,
		WaitTimeout:      file.DefaultWaitTimeout,
		RetryDelay:       file.DefaultRetryDelay,
		FileContainer:    file.NewContainer(),
		CachedViews:      make(ViewMap, 10),
		UncommittedViews: NewUncommittedViewMap(),
		SelectedViews:    nil,
		AffectedRows:     0,
		AutoCommit:       false,
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

func (tx *Transaction) Commit(filter *Filter, expr parser.Expression) error {
	createdFiles, updatedFiles := tx.UncommittedViews.UncommittedFiles()

	createFileInfo := make([]*FileInfo, 0, len(createdFiles))
	updateFileInfo := make([]*FileInfo, 0, len(updatedFiles))

	if 0 < len(createdFiles) {
		for _, fileinfo := range createdFiles {
			view, _ := tx.CachedViews.Get(parser.Identifier{Literal: fileinfo.Path})

			fp := view.FileInfo.Handler.FileForUpdate()
			if err := fp.Truncate(0); err != nil {
				return NewSystemError(err.Error())
			}
			if _, err := fp.Seek(0, io.SeekStart); err != nil {
				return NewSystemError(err.Error())
			}

			_, err := EncodeView(fp, view, fileinfo, tx.Flags)
			if err != nil {
				return NewCommitError(expr, err.Error())
			}
			createFileInfo = append(createFileInfo, view.FileInfo)
		}
	}

	if 0 < len(updatedFiles) {
		for _, fileinfo := range updatedFiles {
			view, _ := tx.CachedViews.Get(parser.Identifier{Literal: fileinfo.Path})

			fp := view.FileInfo.Handler.FileForUpdate()
			if err := fp.Truncate(0); err != nil {
				return NewSystemError(err.Error())
			}
			if _, err := fp.Seek(0, io.SeekStart); err != nil {
				return NewSystemError(err.Error())
			}

			if _, err := EncodeView(fp, view, fileinfo, tx.Flags); err != nil {
				return NewCommitError(expr, err.Error())
			}

			updateFileInfo = append(updateFileInfo, view.FileInfo)
		}
	}

	for _, f := range createFileInfo {
		if err := tx.FileContainer.Commit(f.Handler); err != nil {
			return NewCommitError(expr, err.Error())
		}
		tx.UncommittedViews.Unset(f)
		tx.Session.LogNotice(fmt.Sprintf("Commit: file %q is created.", f.Path), tx.Flags.Quiet)
	}
	for _, f := range updateFileInfo {
		if err := tx.FileContainer.Commit(f.Handler); err != nil {
			return NewCommitError(expr, err.Error())
		}
		tx.UncommittedViews.Unset(f)
		tx.Session.LogNotice(fmt.Sprintf("Commit: file %q is updated.", f.Path), tx.Flags.Quiet)
	}

	msglist := filter.TempViews.Store(tx.UncommittedViews.UncommittedTempViews())
	if 0 < len(msglist) {
		tx.Session.LogNotice(strings.Join(msglist, "\n"), tx.Flags.Quiet)
	}
	tx.UncommittedViews.Clean()
	if err := tx.ReleaseResources(); err != nil {
		return NewCommitError(expr, err.Error())
	}
	return nil
}

func (tx *Transaction) Rollback(filter *Filter, expr parser.Expression) error {
	createdFiles, updatedFiles := tx.UncommittedViews.UncommittedFiles()

	if 0 < len(createdFiles) {
		for _, fileinfo := range createdFiles {
			tx.Session.LogNotice(fmt.Sprintf("Rollback: file %q is deleted.", fileinfo.Path), tx.Flags.Quiet)
		}
	}

	if 0 < len(updatedFiles) {
		for _, fileinfo := range updatedFiles {
			tx.Session.LogNotice(fmt.Sprintf("Rollback: file %q is restored.", fileinfo.Path), tx.Flags.Quiet)
		}
	}

	if filter != nil {
		msglist := filter.TempViews.Restore(tx.UncommittedViews.UncommittedTempViews())
		if 0 < len(msglist) {
			tx.Session.LogNotice(strings.Join(msglist, "\n"), tx.Flags.Quiet)
		}
	}
	tx.UncommittedViews.Clean()
	if err := tx.ReleaseResources(); err != nil {
		return NewRollbackError(expr, err.Error())
	}
	return nil
}

func (tx *Transaction) ReleaseResources() error {
	if err := tx.CachedViews.Clean(tx.FileContainer); err != nil {
		return err
	}
	if err := tx.FileContainer.UnlockAll(); err != nil {
		return err
	}
	return nil
}

func (tx *Transaction) ReleaseResourcesWithErrors() error {
	var errs []error
	if err := tx.CachedViews.CleanWithErrors(tx.FileContainer); err != nil {
		errs = append(errs, err.(*file.ForcedUnlockError).Errors...)
	}
	if err := tx.FileContainer.UnlockAllWithErrors(); err != nil {
		errs = append(errs, err.(*file.ForcedUnlockError).Errors...)
	}

	if errs != nil {
		return file.NewForcedUnlockError(errs)
	}
	return nil
}
