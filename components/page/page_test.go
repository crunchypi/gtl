package page

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/crunchypi/gtl/core"
)

func assertEq[T any](subject string, want T, have T, f func(string)) {
	if f == nil {
		return
	}

	ab, _ := json.Marshal(want)
	bb, _ := json.Marshal(have)

	as := string(ab)
	bs := string(bb)

	if as == bs {
		return
	}

	s := "unexpected '%v':\n\twant: '%v'\n\thave: '%v'\n"
	f(fmt.Sprintf(s, subject, as, bs))
}

// -----------------------------------------------------------------------------
// Tests: NewOnceReader.
// -----------------------------------------------------------------------------

func TestNewOnceReaderIdealEven(t *testing.T) {
	pr := NewOnceReader(NewOnceReaderArgs{Total: 4, Limit: 2})

	err := *new(error)
	val := Page{}

	val, err = pr.Read(context.Background())
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("skip", 0, val.Skip, func(s string) { t.Fatal(s) })
	assertEq("limit", 2, val.Limit, func(s string) { t.Fatal(s) })
	assertEq("total", 4, val.Total, func(s string) { t.Fatal(s) })

	val, err = pr.Read(context.Background())
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("skip", 2, val.Skip, func(s string) { t.Fatal(s) })
	assertEq("limit", 2, val.Limit, func(s string) { t.Fatal(s) })
	assertEq("total", 4, val.Total, func(s string) { t.Fatal(s) })

	val, err = pr.Read(context.Background())
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("skip", 0, val.Skip, func(s string) { t.Fatal(s) })
	assertEq("limit", 0, val.Limit, func(s string) { t.Fatal(s) })
	assertEq("total", 0, val.Total, func(s string) { t.Fatal(s) })
}

func TestNewOnceReaderIdealOdd(t *testing.T) {
	pr := NewOnceReader(NewOnceReaderArgs{Total: 5, Limit: 3})

	err := *new(error)
	val := Page{}

	val, err = pr.Read(context.Background())
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("skip", 0, val.Skip, func(s string) { t.Fatal(s) })
	assertEq("limit", 3, val.Limit, func(s string) { t.Fatal(s) })
	assertEq("total", 5, val.Total, func(s string) { t.Fatal(s) })

	val, err = pr.Read(context.Background())
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("skip", 3, val.Skip, func(s string) { t.Fatal(s) })
	assertEq("limit", 2, val.Limit, func(s string) { t.Fatal(s) })
	assertEq("total", 5, val.Total, func(s string) { t.Fatal(s) })

	val, err = pr.Read(context.Background())
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("skip", 0, val.Skip, func(s string) { t.Fatal(s) })
	assertEq("limit", 0, val.Limit, func(s string) { t.Fatal(s) })
	assertEq("total", 0, val.Total, func(s string) { t.Fatal(s) })
}

// -----------------------------------------------------------------------------
// Tests: NewContReader.
// -----------------------------------------------------------------------------

func TestNewContReader(t *testing.T) {
	lr := core.NewReaderFrom(1, 2, 3)
	pr := NewContReader(NewContReaderArgs{Reader: lr, Limit: 2})

	err := *new(error)
	val := Page{}

	val, err = pr.Read(context.Background())
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("skip", 0, val.Skip, func(s string) { t.Fatal(s) })
	assertEq("limit", 1, val.Limit, func(s string) { t.Fatal(s) })
	assertEq("total", 1, val.Total, func(s string) { t.Fatal(s) })

	val, err = pr.Read(context.Background())
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("skip", 0, val.Skip, func(s string) { t.Fatal(s) })
	assertEq("limit", 2, val.Limit, func(s string) { t.Fatal(s) })
	assertEq("total", 2, val.Total, func(s string) { t.Fatal(s) })

	val, err = pr.Read(context.Background())
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("skip", 0, val.Skip, func(s string) { t.Fatal(s) })
	assertEq("limit", 2, val.Limit, func(s string) { t.Fatal(s) })
	assertEq("total", 3, val.Total, func(s string) { t.Fatal(s) })

	val, err = pr.Read(context.Background())
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("skip", 2, val.Skip, func(s string) { t.Fatal(s) })
	assertEq("limit", 1, val.Limit, func(s string) { t.Fatal(s) })
	assertEq("total", 3, val.Total, func(s string) { t.Fatal(s) })

	val, err = pr.Read(context.Background())
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("skip", 0, val.Skip, func(s string) { t.Fatal(s) })
	assertEq("limit", 0, val.Limit, func(s string) { t.Fatal(s) })
	assertEq("total", 0, val.Total, func(s string) { t.Fatal(s) })
}

func TestNewContReaderWithNilReader(t *testing.T) {
	pr := NewContReader(NewContReaderArgs{Reader: nil, Limit: 2})

	err := *new(error)
	val := Page{}

	val, err = pr.Read(context.Background())
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("skip", 0, val.Skip, func(s string) { t.Fatal(s) })
	assertEq("limit", 0, val.Limit, func(s string) { t.Fatal(s) })
	assertEq("total", 0, val.Total, func(s string) { t.Fatal(s) })
}
