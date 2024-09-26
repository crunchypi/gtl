package stats

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/crunchypi/gtl/core"
)

type StatsStreamed[T any] struct {
	Tag    string         `json:"tag"`
	Val    T              `json:"val"`
	Err    error          `json:"err"`
	CtxMap map[string]any `json:"ctx"`
	Stamp  time.Time      `json:"stamp"`
	Delta  time.Duration  `json:"delta"`
}

type StatsBatched struct {
	Tag    string         `json:"tag"`
	Len    int            `json:"len"`
	Err    error          `json:"err"`
	CtxMap map[string]any `json:"ctx"`
	Stamp  time.Time      `json:"stamp"`
	Delta  time.Duration  `json:"delta"`
}

type NewStreamedTeeReaderArgs[T, U any] struct {
	// Reader is what the func reads from. On nil, the func simply returns
	// a core.ReaderImpl[T], making it pointless. If it returns io.EOF,
	// then the returned reader will skip internal logic and also return io.EOF.
	Reader core.Reader[T]
	// Writer is where stats are written. On nil, stats are not written anywhere,
	// and the func simply passes along values from Reader. Errors coming from
	// here will be wrapped with any err returned from Reader using errors.Join.
	Writer core.Writer[StatsStreamed[U]]
	// Tag will be set to StatsStreamed.Tag which are sent to the Writer. If
	// empty, the value will be set to "<unset>".
	Tag string
	// Fmt defines what to set to StatsStreamed.Val which are sent to the Writer.
	// On nil, StatsStreamed.Val is set to the zero value of U.
	Fmt func(T) U
	// CtxKeys is used to extract values from the ctx given to the returned
	// Reader. These k:v pairs are set to StatsStreamed.CtxMap.
	CtxKeys []string
}

// NewStreamedTeeReader returns a Reader[T] which pulls from args.Reader, while
// writing stats to args.Writer. See args for details.
//
// Note that the errs returned from this reader should be checked with
// errors.Is(...), they may be from both args.Reader (e.g io.EOF), args.Writer
// (e.g io.ErrClosedPipe) or both (wrap). If the err is from args.Writer, then
// the value read from args.Reader may be still valid.
//
// Examples (interactive):
//   - https://go.dev/play/p/xQOOBB9vG0A
func NewStreamedTeeReader[T, U any](args NewStreamedTeeReaderArgs[T, U]) core.Reader[T] {
	if args.Reader == nil {
		return core.ReaderImpl[T]{}
	}
	if args.Writer == nil {
		args.Writer = core.WriterImpl[StatsStreamed[U]]{}
	}
	if args.Tag == "" {
		args.Tag = "<unset>"
	}
	if args.Fmt == nil {
		args.Fmt = func(v T) (r U) { return }
	}

	stamp := time.Now()
	return core.ReaderImpl[T]{
		Impl: func(ctx context.Context) (val T, err error) {
			val, err = args.Reader.Read(ctx)
			if err == io.EOF {
				return
			}

			stats := StatsStreamed[U]{}
			stats.Tag = args.Tag
			stats.Val = args.Fmt(val)
			stats.Err = err
			stats.CtxMap = make(map[string]any, len(args.CtxKeys))
			stats.Stamp = time.Now()
			stats.Delta = stats.Stamp.Sub(stamp)

			if ctx != nil {
				for _, key := range args.CtxKeys {
					stats.CtxMap[key] = ctx.Value(key)
				}
			}

			stamp = stats.Stamp
			err = errors.Join(err, args.Writer.Write(ctx, stats))
			return
		},
	}
}

type NewBatchedTeeReaderArgs[T any] struct {
	// Reader is what the func reads from. On nil, the func simply returns
	// a core.ReaderImpl[T], making it pointless. If it returns io.EOF,
	// then the returned reader will skip internal logic and also return io.EOF.
	Reader core.Reader[[]T]
	// Writer is where stats are written. On nil, stats are not written anywhere,
	// and the func simply passes along values from Reader. Errors coming from
	// here will be wrapped with any err returned from Reader using errors.Join.
	Writer core.Writer[StatsBatched]
	// Tag will be set to StatsStreamed.Tag which are sent to the Writer. If
	// empty, the value will be set to "<unset>".
	Tag string
	// CtxKeys is used to extract values from the ctx given to the returned
	// Reader. These k:v pairs are set to StatsStreamed.CtxMap.
	CtxKeys []string
}

// NewBatchedTeeReader returns a Reader[[]T] which pulls from args.Reader, while
// writing stats to args.Writer. It is similar to NewStreamedTeeReader but
// works with []T and writes stats containing "len" instead of "val".
// See args for details.
//
// Note that the errs returned from this reader should be checked with
// errors.Is(...), they may be from both args.Reader (e.g io.EOF), args.Writer
// (e.g io.ErrClosedPipe) or both (wrap). If the err is from args.Writer, then
// the value read from args.Reader may be still valid.
//
// Examples (interactive):
//   - https://go.dev/play/p/8T-eN52RPoE
func NewBatchedTeeReader[T any](args NewBatchedTeeReaderArgs[T]) core.Reader[[]T] {
	if args.Reader == nil {
		return core.ReaderImpl[[]T]{}
	}
	if args.Writer == nil {
		args.Writer = core.WriterImpl[StatsBatched]{}
	}
	if args.Tag == "" {
		args.Tag = "<unset>"
	}

	stamp := time.Now()
	return core.ReaderImpl[[]T]{
		Impl: func(ctx context.Context) (s []T, err error) {
			s, err = args.Reader.Read(ctx)
			if err == io.EOF {
				return
			}

			stats := StatsBatched{}
			stats.Tag = args.Tag
			stats.Len = len(s)
			stats.Err = err
			stats.CtxMap = make(map[string]any, len(args.CtxKeys))
			stats.Stamp = time.Now()
			stats.Delta = stats.Stamp.Sub(stamp)

			if ctx != nil {
				for _, key := range args.CtxKeys {
					stats.CtxMap[key] = ctx.Value(key)
				}
			}

			stamp = stats.Stamp
			err = errors.Join(err, args.Writer.Write(ctx, stats))
			return
		},
	}
}

type NewStreamedTeeWriterArgs[T, U any] struct {
	// WriterVals is what the returned Writer writes to. On nil, the func simply
	// returns a core.WriterImpl[T]{}, making it pointless. If it returns an
	// io.ErrClosedPipe, then the returned Writer returns early with that err.
	WriterVals core.Writer[T]
	// WriterStats is where stats are written. On nil, stats are not written
	// anywhere, and the returned Writer will simply pass along values to
	// WriterVals. Errors coming from here will be wrapped with any error
	// returned from WriterVals using errors.Join.
	WriterStats core.Writer[StatsStreamed[U]]
	// Tag will be set to StatsStreamed.Tag which are sent to the Writer. If
	// empty, the value will be set to "<unset>".
	Tag string
	// Fmt defines what to set to StatsStreamed.Val which are sent to the Writer.
	// On nil, StatsStreamed.Val is set to the zero value of U.
	Fmt func(T) U
	// CtxKeys is used to extract values from the ctx given to the returned
	// Reader. These k:v pairs are set to StatsStreamed.CtxMap.
	CtxKeys []string
}

// NewStreamedTeeWriter returns a Writer[T] which writes into args.WriterVals
// while writing stats to args.WriterStats. See args for details.
//
// Note that errs returned from this Writer should be checked with errors.Is(...),
// as they may be from args.WriterVals, args.WriterStats, or both. If the
// err is from args.WriterStats, then the value written to the returned Writer
// may have been written successfully.
//
// Examples (interactive):
//   - https://go.dev/play/p/8GYEViyq5hq
func NewStreamedTeeWriter[T, U any](args NewStreamedTeeWriterArgs[T, U]) core.Writer[T] {
	if args.WriterVals == nil {
		return core.WriterImpl[T]{}
	}
	if args.WriterStats == nil {
		args.WriterStats = core.WriterImpl[StatsStreamed[U]]{}
	}
	if args.Tag == "" {
		args.Tag = "<unset>"
	}
	if args.Fmt == nil {
		args.Fmt = func(v T) (r U) { return }
	}

	stamp := time.Now()
	return core.WriterImpl[T]{
		Impl: func(ctx context.Context, val T) (err error) {
			err = args.WriterVals.Write(ctx, val)
			if err == io.ErrClosedPipe {
				return
			}

			stats := StatsStreamed[U]{}
			stats.Tag = args.Tag
			stats.Val = args.Fmt(val)
			stats.Err = err
			stats.CtxMap = make(map[string]any, len(args.CtxKeys))
			stats.Stamp = time.Now()
			stats.Delta = stats.Stamp.Sub(stamp)

			if ctx != nil {
				for _, key := range args.CtxKeys {
					stats.CtxMap[key] = ctx.Value(key)
				}
			}

			stamp = stats.Stamp
			err = errors.Join(err, args.WriterStats.Write(ctx, stats))
			return
		},
	}
}
