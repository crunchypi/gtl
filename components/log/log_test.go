package log

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/crunchypi/gtl/core"
)

var tvVerbose = false
var tvCtxKey = "testKey"
var tvCtxVal = "testVal"
var tvCtx = context.WithValue(context.Background(), tvCtxKey, tvCtxVal)
var tvErr = errors.New("test error")

// Tests printing to stdout are annoying.
func init() {
	if !tvVerbose {
		defaultLogger = slog.New(
			slog.NewJSONHandler(
				bytes.NewBuffer(nil),
				nil,
			),
		)
	}
}

// -----------------------------------------------------------------------------
// Utils: Reader.
// -----------------------------------------------------------------------------

// wraps r, will return an err when r returns an io.EOF
func tfNewErredReader[T any](r core.Reader[T]) core.Reader[T] {
	return core.ReaderImpl[T]{
		Impl: func(ctx context.Context) (val T, err error) {
			val, err = r.Read(ctx)
			if err == io.EOF {
				err = tvErr
			}

			return val, err
		},
	}
}

func tfReadAll[T any](ctx context.Context, r core.Reader[T]) ([]T, error) {
	var v T
	var s = make([]T, 0, 8)
	var err error

	for v, err = r.Read(ctx); err == nil; v, err = r.Read(ctx) {
		s = append(s, v)
	}

	return s, err
}

// -----------------------------------------------------------------------------
// Utils: Writer.
// -----------------------------------------------------------------------------

func tfNewNopWriter[T any]() core.Writer[T] {
	return core.WriterImpl[T]{
		Impl: func(ctx context.Context, v T) error {
			return nil
		},
	}
}

func tfWriteSlice[T any](ctx context.Context, s []T, w core.Writer[T]) error {
	err := *new(error)
	for _, v := range s {
		err = w.Write(ctx, v)
		if err != nil {
			break
		}
	}

	return err
}

// -----------------------------------------------------------------------------
// Tests: NewStreamedReader
// -----------------------------------------------------------------------------

func TestNewStreamedReaderIdeal(t *testing.T) {
	tfReadAll(
		tvCtx,
		NewStreamedReader(
			NewStreamedReaderArgs[string]{
				Reader:  core.NewReaderFrom("a", "b"),
				Logger:  defaultLogger,
				Msg:     "test",
				Fmt:     func(s string) any { return s },
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

func TestNewStreamedReaderWithNilReader(t *testing.T) {
	tfReadAll(
		tvCtx,
		NewStreamedReader(
			NewStreamedReaderArgs[string]{
				Reader:  nil,
				Logger:  defaultLogger,
				Msg:     "test",
				Fmt:     func(s string) any { return s },
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

func TestNewStreamedReaderWithNilLogger(t *testing.T) {
	tfReadAll(
		tvCtx,
		NewStreamedReader(
			NewStreamedReaderArgs[string]{
				Reader:  core.NewReaderFrom("a", "b"),
				Logger:  nil,
				Msg:     "test",
				Fmt:     func(s string) any { return s },
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

func TestNewStreamedReaderWithEmptyMsg(t *testing.T) {
	tfReadAll(
		tvCtx,
		NewStreamedReader(
			NewStreamedReaderArgs[string]{
				Reader:  core.NewReaderFrom("a", "b"),
				Logger:  defaultLogger,
				Msg:     "",
				Fmt:     func(s string) any { return s },
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

func TestNewStreamedReaderWithNilFmt(t *testing.T) {
	tfReadAll(
		tvCtx,
		NewStreamedReader(
			NewStreamedReaderArgs[string]{
				Reader:  core.NewReaderFrom("a", "b"),
				Logger:  defaultLogger,
				Msg:     "test",
				Fmt:     nil,
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

func TestNewStreamedReaderWithNilCtxKeys(t *testing.T) {
	tfReadAll(
		tvCtx,
		NewStreamedReader(
			NewStreamedReaderArgs[string]{
				Reader:  core.NewReaderFrom("a", "b"),
				Logger:  defaultLogger,
				Msg:     "test",
				Fmt:     func(s string) any { return s },
				CtxKeys: nil,
			},
		),
	)
}

func TestNewStreamedReaderWithErr(t *testing.T) {
	tfReadAll(
		tvCtx,
		NewStreamedReader(
			NewStreamedReaderArgs[string]{
				Reader:  tfNewErredReader(core.NewReaderFrom("a", "b")),
				Logger:  defaultLogger,
				Msg:     "test",
				Fmt:     func(s string) any { return s },
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

func TestNewStreamedReaderWithNilCtx(t *testing.T) {
	tfReadAll(
		nil,
		NewStreamedReader(
			NewStreamedReaderArgs[string]{
				Reader:  core.NewReaderFrom("a", "b"),
				Logger:  defaultLogger,
				Msg:     "test",
				Fmt:     func(s string) any { return s },
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

// -----------------------------------------------------------------------------
// Tests: NewBatchedReader
// -----------------------------------------------------------------------------

func TestNewBatchedReaderIdeal(t *testing.T) {
	r := core.NewReaderFrom("a", "b")

	tfReadAll(
		tvCtx,
		NewBatchedReader(
			NewBatchedReaderArgs[string]{
				Reader:  core.NewReaderWithBatching(r, 1),
				Logger:  defaultLogger,
				Msg:     "test",
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

func TestNewBatchedReaderWithNilReader(t *testing.T) {
	tfReadAll(
		tvCtx,
		NewBatchedReader(
			NewBatchedReaderArgs[string]{
				Reader:  nil,
				Logger:  defaultLogger,
				Msg:     "test",
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

func TestNewBatchedReaderWithNilLogger(t *testing.T) {
	r := core.NewReaderFrom("a", "b")

	tfReadAll(
		tvCtx,
		NewBatchedReader(
			NewBatchedReaderArgs[string]{
				Reader:  core.NewReaderWithBatching(r, 1),
				Logger:  nil,
				Msg:     "test",
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

func TestNewBatchedReaderWithEmptyMsg(t *testing.T) {
	r := core.NewReaderFrom("a", "b")

	tfReadAll(
		tvCtx,
		NewBatchedReader(
			NewBatchedReaderArgs[string]{
				Reader:  core.NewReaderWithBatching(r, 1),
				Logger:  defaultLogger,
				Msg:     "",
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

func TestNewBatchedReaderWithNilCtxKeys(t *testing.T) {
	r := core.NewReaderFrom("a", "b")

	tfReadAll(
		tvCtx,
		NewBatchedReader(
			NewBatchedReaderArgs[string]{
				Reader:  core.NewReaderWithBatching(r, 1),
				Logger:  defaultLogger,
				Msg:     "test",
				CtxKeys: nil,
			},
		),
	)
}

func TestNewBatchedReaderWithErr(t *testing.T) {
	r := core.NewReaderFrom("a", "b")

	tfReadAll(
		tvCtx,
		NewBatchedReader(
			NewBatchedReaderArgs[string]{
				Reader:  tfNewErredReader(core.NewReaderWithBatching(r, 1)),
				Logger:  defaultLogger,
				Msg:     "test",
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

func TestNewBatchedReaderWithNilCtx(t *testing.T) {
	r := core.NewReaderFrom("a", "b")

	tfReadAll(
		nil,
		NewBatchedReader(
			NewBatchedReaderArgs[string]{
				Reader:  core.NewReaderWithBatching(r, 1),
				Logger:  defaultLogger,
				Msg:     "test",
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

// -----------------------------------------------------------------------------
// Tests: NewStreamedWriter
// -----------------------------------------------------------------------------

func TestNewStreamedWriterIdeal(t *testing.T) {
	w := tfNewNopWriter[string]()

	tfWriteSlice(
		tvCtx,
		[]string{"a", "b", "c"},
		NewStreamedWriter(
			NewStreamedWriterArgs[string]{
				Writer:  w,
				Logger:  defaultLogger,
				Msg:     "test",
				Fmt:     func(s string) any { return s },
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

func TestNewStreamedWriterWithNilWriter(t *testing.T) {
	tfWriteSlice(
		tvCtx,
		[]string{"a", "b", "c"},
		NewStreamedWriter(
			NewStreamedWriterArgs[string]{
				Writer:  nil,
				Logger:  defaultLogger,
				Msg:     "test",
				Fmt:     func(s string) any { return s },
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

func TestNewStreamedWriterWithNilLogger(t *testing.T) {
	w := tfNewNopWriter[string]()

	tfWriteSlice(
		tvCtx,
		[]string{"a", "b", "c"},
		NewStreamedWriter(
			NewStreamedWriterArgs[string]{
				Writer:  w,
				Logger:  nil,
				Msg:     "test",
				Fmt:     func(s string) any { return s },
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

func TestNewStreamedWriterWithEmptyMsg(t *testing.T) {
	w := tfNewNopWriter[string]()

	tfWriteSlice(
		tvCtx,
		[]string{"a", "b", "c"},
		NewStreamedWriter(
			NewStreamedWriterArgs[string]{
				Writer:  w,
				Logger:  defaultLogger,
				Msg:     "",
				Fmt:     func(s string) any { return s },
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

func TestNewStreamedWriterWithNilFmt(t *testing.T) {
	w := tfNewNopWriter[string]()

	tfWriteSlice(
		tvCtx,
		[]string{"a", "b", "c"},
		NewStreamedWriter(
			NewStreamedWriterArgs[string]{
				Writer:  w,
				Logger:  defaultLogger,
				Msg:     "test",
				Fmt:     nil,
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

func TestNewStreamedWriterWithNilCtxKeys(t *testing.T) {
	w := tfNewNopWriter[string]()

	tfWriteSlice(
		tvCtx,
		[]string{"a", "b", "c"},
		NewStreamedWriter(
			NewStreamedWriterArgs[string]{
				Writer:  w,
				Logger:  defaultLogger,
				Msg:     "test",
				Fmt:     func(s string) any { return s },
				CtxKeys: nil,
			},
		),
	)
}

func TestNewStreamedWriterWithErrClosedPipe(t *testing.T) {
	tfWriteSlice(
		tvCtx,
		[]string{"a", "b", "c"},
		NewStreamedWriter(
			NewStreamedWriterArgs[string]{
				Writer:  core.WriterImpl[string]{},
				Logger:  defaultLogger,
				Msg:     "test",
				Fmt:     func(s string) any { return s },
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

func TestNewStreamedWriterWithErrUnknown(t *testing.T) {
	w := core.WriterImpl[string]{}
	w.Impl = func(ctx context.Context, s string) error { return tvErr }

	tfWriteSlice(
		tvCtx,
		[]string{"a", "b", "c"},
		NewStreamedWriter(
			NewStreamedWriterArgs[string]{
				Writer:  w,
				Logger:  defaultLogger,
				Msg:     "test",
				Fmt:     func(s string) any { return s },
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

func TestNewStreamedWriterWithNilCtx(t *testing.T) {
	w := tfNewNopWriter[string]()

	tfWriteSlice(
		nil,
		[]string{"a", "b", "c"},
		NewStreamedWriter(
			NewStreamedWriterArgs[string]{
				Writer:  w,
				Logger:  defaultLogger,
				Msg:     "test",
				Fmt:     func(s string) any { return s },
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

// -----------------------------------------------------------------------------
// Tests: NewBatchedWriter
// -----------------------------------------------------------------------------

func TestNewBatchedWriterIdeal(t *testing.T) {
	w := tfNewNopWriter[[]string]()

	tfWriteSlice(
		tvCtx,
		[][]string{{"a", "b"}, {"c"}},
		NewBatchedWriter(
			NewBatchedWriterArgs[string]{
				Writer:  w,
				Logger:  defaultLogger,
				Msg:     "test",
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

func TestNewBatchedWriterWithNilWriter(t *testing.T) {
	tfWriteSlice(
		tvCtx,
		[][]string{{"a", "b"}, {"c"}},
		NewBatchedWriter(
			NewBatchedWriterArgs[string]{
				Writer:  nil,
				Logger:  defaultLogger,
				Msg:     "test",
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

func TestNewBatchedWriterWithNilLogger(t *testing.T) {
	w := tfNewNopWriter[[]string]()

	tfWriteSlice(
		tvCtx,
		[][]string{{"a", "b"}, {"c"}},
		NewBatchedWriter(
			NewBatchedWriterArgs[string]{
				Writer:  w,
				Logger:  nil,
				Msg:     "test",
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

func TestNewBatchedWriterWithEmptyMsg(t *testing.T) {
	w := tfNewNopWriter[[]string]()

	tfWriteSlice(
		tvCtx,
		[][]string{{"a", "b"}, {"c"}},
		NewBatchedWriter(
			NewBatchedWriterArgs[string]{
				Writer:  w,
				Logger:  defaultLogger,
				Msg:     "",
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

func TestNewBatchedWriterWithNilCtxKeys(t *testing.T) {
	w := tfNewNopWriter[[]string]()

	tfWriteSlice(
		tvCtx,
		[][]string{{"a", "b"}, {"c"}},
		NewBatchedWriter(
			NewBatchedWriterArgs[string]{
				Writer:  w,
				Logger:  defaultLogger,
				Msg:     "test",
				CtxKeys: nil,
			},
		),
	)
}

func TestNewBatchedWriterWithErrClosedPipe(t *testing.T) {
	w := core.WriterImpl[[]string]{}

	tfWriteSlice(
		tvCtx,
		[][]string{{"a", "b"}, {"c"}},
		NewBatchedWriter(
			NewBatchedWriterArgs[string]{
				Writer:  w,
				Logger:  defaultLogger,
				Msg:     "test",
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

func TestNewBatchedWriterWithErrUnknown(t *testing.T) {
	w := core.WriterImpl[[]string]{}
	w.Impl = func(ctx context.Context, s []string) error { return tvErr }

	tfWriteSlice(
		tvCtx,
		[][]string{{"a", "b"}, {"c"}},
		NewBatchedWriter(
			NewBatchedWriterArgs[string]{
				Writer:  w,
				Logger:  defaultLogger,
				Msg:     "test",
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}

func TestNewBatchedWriterWithNilCtx(t *testing.T) {
	w := tfNewNopWriter[[]string]()

	tfWriteSlice(
		nil,
		[][]string{{"a", "b"}, {"c"}},
		NewBatchedWriter(
			NewBatchedWriterArgs[string]{
				Writer:  w,
				Logger:  defaultLogger,
				Msg:     "test",
				CtxKeys: []string{tvCtxKey},
			},
		),
	)
}
