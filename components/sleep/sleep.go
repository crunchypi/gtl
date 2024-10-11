package sleep

import (
	"context"
	"time"

	"github.com/crunchypi/gtl/core"
)

type NewStaticReaderArgs[T any] struct {
	Reader core.Reader[T]
	Delay  time.Duration
}

// NewStaticReader returns a reader which wraps args.Reader with sleep/delay,
// the duration being set with args.Delay. When reading from the returned reader,
// the operation happens either after args.Delay has elapsed, or if ctx is done.
//
// Examples (interactive):
//   - https://go.dev/play/p/KsWoywcqcB6
func NewStaticReader[T any](args NewStaticReaderArgs[T]) core.Reader[T] {
	if args.Reader == nil {
		return core.ReaderImpl[T]{}
	}

	return core.ReaderImpl[T]{
		Impl: func(ctx context.Context) (v T, err error) {
			if ctx == nil {
				ctx = context.Background()
			}

			select {
			case <-ctx.Done():
			case <-time.After(args.Delay):
			}

			return args.Reader.Read(ctx)
		},
	}
}

type NewDynamicReaderArgs[T any] struct {
	Reader core.Reader[T]
	Delay  time.Duration
}

// NewDynamicReader returns a reader which wraps args.Reader with sleep/delay,
// the minimum duration being set with args.Delay.
//
// Unlike NewStaticReader, this one has a couple extra properties. Firstly it
// tries to adjust the sleep duration to a constant args.Delay, it does so by
// subtracting args.Delay by the time it took to read from args.Reader.
// Secondly, you may set ctx value "bounds" if you know how many times args.Reader
// may be called before io.EOF, in that case the whole stream will be read in
// approximately args.Delay. This is useful if you want a complete ETL pipeline
// to take a specific amount of time.
//
// Examples (interactive):
//   - https://go.dev/play/p/bCj5Z1wdKMu
func NewDynamicReader[T any](args NewDynamicReaderArgs[T]) core.Reader[T] {
	if args.Reader == nil {
		return core.ReaderImpl[T]{}
	}

	return core.ReaderImpl[T]{
		Impl: func(ctx context.Context) (v T, err error) {
			if ctx == nil {
				ctx = context.Background()
			}

			ts := time.Now()
			v, err = args.Reader.Read(ctx)
			if err != nil {
				return
			}

			d := args.Delay
			if bounds, ok := ctx.Value("bounds").(int); ok {
				d /= time.Duration(bounds)
			}

			select {
			case <-ctx.Done():
			case <-time.After(d - time.Now().Sub(ts)):
			}

			return
		},
	}
}

type NewStaticWriterArgs[T any] struct {
	Writer core.Writer[T]
	Delay  time.Duration
}

// NewStaticWriter returns a Writer which writes to args.Writer and then sleeps
// for the duration defined with args.Delay, or until ctx is done.
//
// Examples (interactive):
//   - https://go.dev/play/p/cRdqr85gAh2
func NewStaticWriter[T any](args NewStaticWriterArgs[T]) core.Writer[T] {
	if args.Writer == nil {
		return core.WriterImpl[T]{}
	}

	return core.WriterImpl[T]{
		Impl: func(ctx context.Context, val T) (err error) {
			if ctx == nil {
				ctx = context.Background()
			}

			err = args.Writer.Write(ctx, val)
			select {
			case <-ctx.Done():
			case <-time.After(args.Delay):
			}

			return
		},
	}
}
