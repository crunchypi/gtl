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
