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
