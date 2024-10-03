package eventloop

import (
	"context"

	"github.com/crunchypi/gtl/core"
)

type NewArgs[T any] struct {
	Ctx    context.Context
	Reader core.Reader[T]
	Writer core.Writer[T]
}

// New spawns a new goroutine in which values are read from args.Reader and
// written to args.Writer, using args.Ctx as context.
//
// The reading and writing is done forever, or until either the reader or writer
// returns an err, or the returned cancel func is called.
// Note that errors coming from either args.Reader and args.Writer are not used
// for anything besides breaking the internal loop, you are intended to pick them
// up with decorators around the Reader and Writer. Also see pkg stats and log.
//
// Examples (interactive):
//   - https://go.dev/play/p/bPO8cOXpyqW
func New[T any](args NewArgs[T]) (ctx context.Context, ctxCancel context.CancelFunc) {
	if args.Ctx == nil {
		args.Ctx = context.Background()
	}

	ctx, ctxCancel = context.WithCancel(args.Ctx)

	ok := true
	ok = ok && args.Reader != nil
	ok = ok && args.Writer != nil
	if !ok {
		ctxCancel()
		return
	}

	go func() {
		defer ctxCancel()

		for {
			select {
			case <-ctx.Done():
				break
			default:
			}

			v, err := args.Reader.Read(ctx)
			if err != nil {
				break
			}

			err = args.Writer.Write(ctx, v)
			if err != nil {
				break
			}
		}
	}()

	return
}
