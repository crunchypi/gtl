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