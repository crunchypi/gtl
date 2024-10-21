package page

import (
	"context"
	"io"

	"github.com/crunchypi/gtl/core"
)

// Page represents pagination directives.
type Page struct {
	Skip  int // Pagination skip.
	Limit int // Pagination limit.
	Total int // Total number of items.
}

// Paged is a wrapper of Page, containing a generic value.
type Paged[T any] struct {
	Page
	Val T
}

type NewOnceReaderArgs struct {
	Total int
	Limit int
}

// NewOnceReader returns a Reader of pagination directives which supports
// paging from 0 to 'total' with the given 'limit', then returns an io.EOF.
// It is useful for e.g paging through a database if you know the total size.
//
// Examples (interactive):
//   - https://go.dev/play/p/NOuwlVmJwbg
func NewOnceReader(args NewOnceReaderArgs) core.Reader[Page] {
	var skip int

	return core.ReaderImpl[Page]{
		Impl: func(ctx context.Context) (p Page, err error) {
			if skip >= args.Total {
				return p, io.EOF
			}

			p.Skip = skip
			p.Limit = args.Limit
			p.Total = args.Total

			if skip+args.Limit > args.Total {
				p.Limit -= (skip + args.Limit) - args.Total
			}

			skip += args.Limit
			return
		},
	}
}

type NewContReaderArgs struct {
	// Ints from here are passed as args.Total for NewOnceReader.
	Reader core.Reader[int]
	// Limit is passed to args.Limit for NewOnceReader.
	Limit int
}

// NewContReader passes ints from args.Reader to NewOnceReader, from which pages
// are returned here. When all pages are read, a new int from args.Reader is
// passed to a new NewOnceReader, and so on.
//
// Examples (interactive):
//   - https://go.dev/play/p/Dk2hZM7Wxi7
func NewContReader(args NewContReaderArgs) core.Reader[Page] {
	if args.Reader == nil {
		return core.ReaderImpl[Page]{}
	}

	ln := 0                                  // Last n
	pr := NewOnceReader(NewOnceReaderArgs{}) // Page reader

	return core.ReaderImpl[Page]{
		Impl: func(ctx context.Context) (p Page, err error) {
			// First/next.
			p, err = pr.Read(ctx)
			if err != nil {

				// Read next bound.
				ln, err = args.Reader.Read(ctx)
				if err != nil {
					return
				}

				// Retry.
				pr = NewOnceReader(
					NewOnceReaderArgs{
						Total: ln,
						Limit: args.Limit,
					},
				)

				p, err = pr.Read(ctx)
			}

			return p, err
		},
	}
}

type NewOnceWriterArgs[T any] struct {
	Writer core.Writer[Paged[T]]
	// Total is passed to NewOnceReader
	Total int
	// Limit is passed to NewOnceReader
	Limit int
}

// NewOnceWriter returns a writer which writes values to args.Writer along with
// pagination directives, based on args.Total and args.Limit, which are passed
// to NewOnceReader under the hood. When all pages are written, the writer
// returned here will give an io.ErrClosedPipe.
//
// Examples (interactive):
//   - https://go.dev/play/p/RfhamjAXEFE
func NewOnceWriter[T any](args NewOnceWriterArgs[T]) core.Writer[T] {
	if args.Writer == nil {
		return core.WriterImpl[T]{}
	}

	pr := NewOnceReader(NewOnceReaderArgs{Total: args.Total, Limit: args.Limit})
	return core.WriterImpl[T]{
		Impl: func(ctx context.Context, val T) (err error) {
			p, err := pr.Read(ctx)
			if err != nil {
				return io.ErrClosedPipe
			}

			return args.Writer.Write(ctx, Paged[T]{Page: p, Val: val})
		},
	}
}
