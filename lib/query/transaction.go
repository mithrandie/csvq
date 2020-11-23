package query

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mithrandie/csvq/lib/cmd"
	"github.com/mithrandie/csvq/lib/file"
	"github.com/mithrandie/csvq/lib/parser"
	"github.com/mithrandie/csvq/lib/value"

	"github.com/mithrandie/go-text/color"
	"github.com/mithrandie/go-text/fixedlen"
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

	flagMutex *sync.RWMutex

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
		flagMutex:          &sync.RWMutex{},
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

func (tx *Transaction) Commit(ctx context.Context, scope *ReferenceScope, expr parser.Expression) error {
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

			if _, err := EncodeView(ctx, fp, view, fileinfo.ExportOptions(tx), tx.Palette); err != nil {
				return NewCommitError(expr, err.Error())
			}

			if !tx.Flags.ExportOptions.StripEndingLineBreak && !(fileinfo.Format == cmd.FIXED && fileinfo.SingleLine) {
				if _, err := fp.Write([]byte(tx.Flags.ExportOptions.LineBreak.Value())); err != nil {
					return NewCommitError(expr, err.Error())
				}
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

			if _, err := EncodeView(ctx, fp, view, fileinfo.ExportOptions(tx), tx.Palette); err != nil {
				return NewCommitError(expr, err.Error())
			}

			if !tx.Flags.ExportOptions.StripEndingLineBreak && !(fileinfo.Format == cmd.FIXED && fileinfo.SingleLine) {
				if _, err := fp.Write([]byte(tx.Flags.ExportOptions.LineBreak.Value())); err != nil {
					return NewCommitError(expr, err.Error())
				}
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

	msglist := scope.StoreTemporaryTable(tx.Session, tx.uncommittedViews.UncommittedTempViews())
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

func (tx *Transaction) Rollback(scope *ReferenceScope, expr parser.Expression) error {
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

	if scope != nil {
		msglist := scope.RestoreTemporaryTable(tx.uncommittedViews.UncommittedTempViews())
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

var errNotAllowdFlagFormat = errors.New("not allowed flag format")
var errInvalidFlagName = errors.New("invalid flag name")

func (tx *Transaction) SetFormatFlag(value interface{}, outFile string) error {
	return tx.setFlag(cmd.FormatFlag, value, outFile)
}

func (tx *Transaction) SetFlag(key string, value interface{}) error {
	return tx.setFlag(key, value, "")
}

func (tx *Transaction) setFlag(key string, value interface{}, outFile string) error {
	tx.flagMutex.Lock()
	defer tx.flagMutex.Unlock()

	var err error

	switch strings.ToUpper(key) {
	case cmd.RepositoryFlag:
		if s, ok := value.(string); ok {
			err = tx.Flags.SetRepository(s)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.TimezoneFlag:
		if s, ok := value.(string); ok {
			err = tx.Flags.SetLocation(s)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.DatetimeFormatFlag:
		if s, ok := value.(string); ok {
			tx.Flags.SetDatetimeFormat(s)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.AnsiQuotesFlag:
		if b, ok := value.(bool); ok {
			tx.Flags.SetAnsiQuotes(b)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.StrictEqualFlag:
		if b, ok := value.(bool); ok {
			tx.Flags.SetStrictEqual(b)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.WaitTimeoutFlag:
		if f, ok := value.(float64); ok {
			tx.UpdateWaitTimeout(f, file.DefaultRetryDelay)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.ImportFormatFlag:
		if s, ok := value.(string); ok {
			err = tx.Flags.SetImportFormat(s)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.DelimiterFlag:
		if s, ok := value.(string); ok {
			err = tx.Flags.SetDelimiter(s)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.DelimiterPositionsFlag:
		if s, ok := value.(string); ok {
			err = tx.Flags.SetDelimiterPositions(s)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.JsonQueryFlag:
		if s, ok := value.(string); ok {
			tx.Flags.SetJsonQuery(s)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.EncodingFlag:
		if s, ok := value.(string); ok {
			err = tx.Flags.SetEncoding(s)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.NoHeaderFlag:
		if b, ok := value.(bool); ok {
			tx.Flags.SetNoHeader(b)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.WithoutNullFlag:
		if b, ok := value.(bool); ok {
			tx.Flags.SetWithoutNull(b)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.FormatFlag:
		if s, ok := value.(string); ok {
			err = tx.Flags.SetFormat(s, outFile)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.ExportEncodingFlag:
		if s, ok := value.(string); ok {
			err = tx.Flags.SetWriteEncoding(s)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.ExportDelimiterFlag:
		if s, ok := value.(string); ok {
			err = tx.Flags.SetWriteDelimiter(s)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.ExportDelimiterPositionsFlag:
		if s, ok := value.(string); ok {
			err = tx.Flags.SetWriteDelimiterPositions(s)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.WithoutHeaderFlag:
		if b, ok := value.(bool); ok {
			tx.Flags.SetWithoutHeader(b)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.LineBreakFlag:
		if s, ok := value.(string); ok {
			err = tx.Flags.SetLineBreak(s)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.EncloseAllFlag:
		if b, ok := value.(bool); ok {
			tx.Flags.SetEncloseAll(b)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.JsonEscapeFlag:
		if s, ok := value.(string); ok {
			err = tx.Flags.SetJsonEscape(s)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.PrettyPrintFlag:
		if b, ok := value.(bool); ok {
			tx.Flags.SetPrettyPrint(b)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.StripEndingLineBreakFlag:
		if b, ok := value.(bool); ok {
			tx.Flags.SetStripEndingLineBreak(b)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.EastAsianEncodingFlag:
		if b, ok := value.(bool); ok {
			tx.Flags.SetEastAsianEncoding(b)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.CountDiacriticalSignFlag:
		if b, ok := value.(bool); ok {
			tx.Flags.SetCountDiacriticalSign(b)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.CountFormatCodeFlag:
		if b, ok := value.(bool); ok {
			tx.Flags.SetCountFormatCode(b)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.ColorFlag:
		if b, ok := value.(bool); ok {
			tx.UseColor(b)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.QuietFlag:
		if b, ok := value.(bool); ok {
			tx.Flags.SetQuiet(b)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.LimitRecursion:
		if i, ok := value.(int64); ok {
			tx.Flags.SetLimitRecursion(i)
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.CPUFlag:
		if i, ok := value.(int64); ok {
			tx.Flags.SetCPU(int(i))
		} else {
			err = errNotAllowdFlagFormat
		}
	case cmd.StatsFlag:
		if b, ok := value.(bool); ok {
			tx.Flags.SetStats(b)
		} else {
			err = errNotAllowdFlagFormat
		}
	default:
		err = errInvalidFlagName
	}

	return err
}

func (tx *Transaction) GetFlag(key string) (value.Primary, bool) {
	tx.flagMutex.RLock()
	defer tx.flagMutex.RUnlock()

	var val value.Primary
	var ok = true

	switch strings.ToUpper(key) {
	case cmd.RepositoryFlag:
		val = value.NewString(tx.Flags.Repository)
	case cmd.TimezoneFlag:
		val = value.NewString(tx.Flags.Location)
	case cmd.DatetimeFormatFlag:
		s := ""
		if 0 < len(tx.Flags.DatetimeFormat) {
			list := make([]string, 0, len(tx.Flags.DatetimeFormat))
			for _, f := range tx.Flags.DatetimeFormat {
				list = append(list, "\""+f+"\"")
			}
			s = "[" + strings.Join(list, ", ") + "]"
		}
		val = value.NewString(s)
	case cmd.AnsiQuotesFlag:
		val = value.NewBoolean(tx.Flags.AnsiQuotes)
	case cmd.StrictEqualFlag:
		val = value.NewBoolean(tx.Flags.StrictEqual)
	case cmd.WaitTimeoutFlag:
		val = value.NewFloat(tx.Flags.WaitTimeout)
	case cmd.ImportFormatFlag:
		val = value.NewString(tx.Flags.ImportOptions.Format.String())
	case cmd.DelimiterFlag:
		val = value.NewString(string(tx.Flags.ImportOptions.Delimiter))
	case cmd.DelimiterPositionsFlag:
		s := fixedlen.DelimiterPositions(tx.Flags.ImportOptions.DelimiterPositions).String()
		if tx.Flags.ImportOptions.SingleLine {
			s = "S" + s
		}
		val = value.NewString(s)
	case cmd.JsonQueryFlag:
		val = value.NewString(tx.Flags.ImportOptions.JsonQuery)
	case cmd.EncodingFlag:
		val = value.NewString(tx.Flags.ImportOptions.Encoding.String())
	case cmd.NoHeaderFlag:
		val = value.NewBoolean(tx.Flags.ImportOptions.NoHeader)
	case cmd.WithoutNullFlag:
		val = value.NewBoolean(tx.Flags.ImportOptions.WithoutNull)
	case cmd.FormatFlag:
		val = value.NewString(tx.Flags.ExportOptions.Format.String())
	case cmd.ExportEncodingFlag:
		val = value.NewString(tx.Flags.ExportOptions.Encoding.String())
	case cmd.ExportDelimiterFlag:
		val = value.NewString(string(tx.Flags.ExportOptions.Delimiter))
	case cmd.ExportDelimiterPositionsFlag:
		s := fixedlen.DelimiterPositions(tx.Flags.ExportOptions.DelimiterPositions).String()
		if tx.Flags.ExportOptions.SingleLine {
			s = "S" + s
		}
		val = value.NewString(s)
	case cmd.WithoutHeaderFlag:
		val = value.NewBoolean(tx.Flags.ExportOptions.WithoutHeader)
	case cmd.LineBreakFlag:
		val = value.NewString(tx.Flags.ExportOptions.LineBreak.String())
	case cmd.EncloseAllFlag:
		val = value.NewBoolean(tx.Flags.ExportOptions.EncloseAll)
	case cmd.JsonEscapeFlag:
		val = value.NewString(cmd.JsonEscapeTypeToString(tx.Flags.ExportOptions.JsonEscape))
	case cmd.PrettyPrintFlag:
		val = value.NewBoolean(tx.Flags.ExportOptions.PrettyPrint)
	case cmd.StripEndingLineBreakFlag:
		val = value.NewBoolean(tx.Flags.ExportOptions.StripEndingLineBreak)
	case cmd.EastAsianEncodingFlag:
		val = value.NewBoolean(tx.Flags.ExportOptions.EastAsianEncoding)
	case cmd.CountDiacriticalSignFlag:
		val = value.NewBoolean(tx.Flags.ExportOptions.CountDiacriticalSign)
	case cmd.CountFormatCodeFlag:
		val = value.NewBoolean(tx.Flags.ExportOptions.CountFormatCode)
	case cmd.ColorFlag:
		val = value.NewBoolean(tx.Flags.ExportOptions.Color)
	case cmd.QuietFlag:
		val = value.NewBoolean(tx.Flags.Quiet)
	case cmd.LimitRecursion:
		val = value.NewInteger(tx.Flags.LimitRecursion)
	case cmd.CPUFlag:
		val = value.NewInteger(int64(tx.Flags.CPU))
	case cmd.StatsFlag:
		val = value.NewBoolean(tx.Flags.Stats)
	default:
		ok = false
	}
	return val, ok
}
