package log

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/crunchypi/gtl/core"
)

var defaultLogger *slog.Logger

func init() {
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})
	defaultLogger = slog.New(h)
}

func ctxValues(ctx context.Context, keys []string) map[string]any {
	m := make(map[string]any, len(keys))

	if ctx == nil {
		return m
	}

	for _, k := range keys {
		m[k] = ctx.Value(k)
	}

	return m
}

type NewStreamedReaderArgs[T any] struct {
	Reader  core.Reader[T] // On nil, the func returns core.ReaderImpl[T]
	Logger  *slog.Logger   // On nil, will use a default logger.
	Msg     string         // On "" , will set the log "msg" to "<unset>"
	Fmt     func(T) any    // On nil, will set the log "val" to the value of T.
	CtxKeys []string       // On nil, will set the log "ctx" to nil.
}

// NewStreamedReader returns a reader which wraps args.Reader with logging.
//
// Logging format:
//
//	{"time":"...","level":"...","msg":"...","err":"...","val":"...","ctx":{...}}
//
// Logging format details:
//   - "time": Format depends on args.Logger. Default is RFC3999.
//   - "level": Normally INFO, may be ERROR on errs other than io.EOF.
//   - "msg": Set to args.Msg. Will be "<unset>" if not set.
//   - "err": Set to read errs.
//   - "val": Values from args.Reader, formatted by args.Fmt.
//   - "ctx": Key-val pairs from args.CtxKeys and ctx given to the reader.
//
// Be sure to check out docs for NewStreamedReaderArgs, as behaviour depends
// a bit on what the arg struct contains.
func NewStreamedReader[T any](args NewStreamedReaderArgs[T]) core.Reader[T] {
	if args.Reader == nil {
		return core.ReaderImpl[T]{}
	}
	if args.Logger == nil {
		args.Logger = defaultLogger
	}
	if args.Msg == "" {
		args.Msg = "<unset>"
	}
	if args.Fmt == nil {
		args.Fmt = func(v T) any { return v }
	}

	return core.ReaderImpl[T]{
		Impl: func(ctx context.Context) (val T, err error) {
			val, err = args.Reader.Read(ctx)
			if err == io.EOF {
				return
			}

			l := args.Logger
			l = l.With("err", err)
			l = l.With("val", args.Fmt(val))
			l = l.With("ctx", ctxValues(ctx, args.CtxKeys))

			f := l.Info
			if err != nil {
				f = l.Error
			}

			f(args.Msg)
			return val, err
		},
	}
}

type NewBatchedReaderArgs[T any] struct {
	Reader  core.Reader[[]T] // On nil, the func returns core.ReaderImpl[[]T]
	Logger  *slog.Logger     // On nil, will use a default logger.
	Msg     string           // On "" , will set the log "msg" to "<unset>"
	CtxKeys []string         // On nil, will set the log "ctx" to nil.
}

// NewBatchedReader returns a reader which wraps args.Reader with logging.
//
// Logging format:
//
//	{"time":"...","level":"...","msg":"...","err":"...","len":"...","ctx":{...}}
//
// Logging format details:
//   - "time": Format depends on args.Logger. Default is RFC3999.
//   - "level": Normally INFO, may be ERROR on errs other than io.EOF.
//   - "msg": Set to args.Msg. Will be "<unset>" if not set.
//   - "err": Set to read errs.
//   - "len": The len of values read from args.Reader.
//   - "ctx": Key-val pairs from args.CtxKeys and ctx given to the reader.
//
// Be sure to check out docs for NewBatchedReaderArgs, as behaviour depends
// a bit on what the arg struct contains.
func NewBatchedReader[T any](args NewBatchedReaderArgs[T]) core.Reader[[]T] {
	if args.Reader == nil {
		return core.ReaderImpl[[]T]{}
	}
	if args.Logger == nil {
		args.Logger = defaultLogger
	}
	if args.Msg == "" {
		args.Msg = "<unset>"
	}

	return core.ReaderImpl[[]T]{
		Impl: func(ctx context.Context) (s []T, err error) {
			s, err = args.Reader.Read(ctx)
			if err == io.EOF {
				return
			}

			l := args.Logger
			l = l.With("err", err)
			l = l.With("len", len(s))
			l = l.With("ctx", ctxValues(ctx, args.CtxKeys))

			f := l.Info
			if err != nil {
				f = l.Error
			}

			f(args.Msg)
			return s, err
		},
	}
}

type NewStreamedWriterArgs[T any] struct {
	Writer  core.Writer[T] // On nil, the func returns core.WriterImpl[T]
	Logger  *slog.Logger   // On nil, will use a default logger.
	Msg     string         // On "" , will set the log "msg" to "<unset>"
	Fmt     func(T) any    // On nil, will set the log "val" to the value of T.
	CtxKeys []string       // On nil, will set the log "ctx" to nil.
}

// NewStreamedWriter returns a writer which accepts values and passes them to
// args.Writer while logging with args.Logger.
//
// Logging format:
//
//	{"time":"...","level":"...","msg":"...","err":"...","val":"...","ctx":{...}}
//
// Logging format details:
//   - "time": Format depends on args.Logger. Default is RFC3999.
//   - "level": Normally INFO, may be ERROR on errs other than io.ErrClosedPipe.
//   - "msg": Set to args.Msg. Will be "<unset>" if not set.
//   - "err": Set to write errs.
//   - "val": Set to values put into this writer, formatted by args.Fmt.
//   - "ctx": Key-val pairs from args.CtxKeys and ctx given to the reader.
//
// Be sure to check out docs for NewStreamedWriterArgs, as behaviour depends
// a bit on what the arg struct contains.
func NewStreamedWriter[T any](args NewStreamedWriterArgs[T]) core.Writer[T] {
	if args.Writer == nil {
		return core.WriterImpl[T]{}
	}
	if args.Logger == nil {
		args.Logger = defaultLogger
	}
	if args.Msg == "" {
		args.Msg = "<unset>"
	}
	if args.Fmt == nil {
		args.Fmt = func(v T) any { return v }
	}

	return core.WriterImpl[T]{
		Impl: func(ctx context.Context, val T) (err error) {
			err = args.Writer.Write(ctx, val)
			if err == io.ErrClosedPipe {
				return
			}

			l := args.Logger
			l = l.With("err", err)
			l = l.With("val", args.Fmt(val))
			l = l.With("ctx", ctxValues(ctx, args.CtxKeys))

			f := l.Info
			if err != nil {
				f = l.Error
			}

			f(args.Msg)
			return err
		},
	}
}
