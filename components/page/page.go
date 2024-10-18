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

// Paged is a wrapper of Page, containing a generic value. Not used by this
// pkg, but defined due to its assumed usefulness since we sometimes want
// this data being fed through a ETL pipeline.
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
