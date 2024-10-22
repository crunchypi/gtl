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

// -----------------------------------------------------------------------------
// Tests: NewOnceWriter.
// -----------------------------------------------------------------------------

func TestNewOnceWriterIdeal(t *testing.T) {
	rw := core.NewReadWriterFrom[Paged[string]]()
	vw := NewOnceWriter(
		NewOnceWriterArgs[string]{
			Writer: rw,
			Total:  4,
			Limit:  2,
		},
	)

	err := *new(error)
	val := Paged[string]{}

	vw.Write(context.Background(), "a")
	val, err = rw.Read(context.Background())
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("skip", 0, val.Skip, func(s string) { t.Fatal(s) })
	assertEq("limit", 2, val.Limit, func(s string) { t.Fatal(s) })
	assertEq("total", 4, val.Total, func(s string) { t.Fatal(s) })
	assertEq("val", "a", val.Val, func(s string) { t.Fatal(s) })

	vw.Write(context.Background(), "b")
	val, err = rw.Read(context.Background())
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("skip", 2, val.Skip, func(s string) { t.Fatal(s) })
	assertEq("limit", 2, val.Limit, func(s string) { t.Fatal(s) })
	assertEq("total", 4, val.Total, func(s string) { t.Fatal(s) })
	assertEq("val", "b", val.Val, func(s string) { t.Fatal(s) })

	vw.Write(context.Background(), "c")
	val, err = rw.Read(context.Background())
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("skip", 0, val.Skip, func(s string) { t.Fatal(s) })
	assertEq("limit", 0, val.Limit, func(s string) { t.Fatal(s) })
	assertEq("total", 0, val.Total, func(s string) { t.Fatal(s) })
	assertEq("val", "", val.Val, func(s string) { t.Fatal(s) })
}

func TestNewOnceWriterWithNilWriter(t *testing.T) {
	vw := NewOnceWriter(
		NewOnceWriterArgs[string]{
			Writer: nil,
			Total:  4,
			Limit:  2,
		},
	)

	err := vw.Write(context.Background(), "c")
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
}

func TestNewOnceWriterWithWriteErr(t *testing.T) {
	pw := core.WriterImpl[Paged[string]]{}
	vw := NewOnceWriter(
		NewOnceWriterArgs[string]{
			Writer: pw,
			Total:  4,
			Limit:  2,
		},
	)

	err := vw.Write(context.Background(), "c")
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
}

// -----------------------------------------------------------------------------
// Tests: NewContWriter.
// -----------------------------------------------------------------------------

func TestNewContWriterIdeal(t *testing.T) {
	rw := core.NewReadWriterFrom[Paged[string]]()
	vw := NewContWriter(
		NewContWriterArgs[string]{
			Reader: core.NewReaderFrom(1, 2),
			Writer: rw,
			Limit:  2,
		},
	)

	err := *new(error)
	val := Paged[string]{}

	// Write pages.
	err = vw.Write(context.Background(), "a")
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	err = vw.Write(context.Background(), "b")
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	err = vw.Write(context.Background(), "c")
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })

	// Read and validate.
	val, err = rw.Read(context.Background())
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("skip", 0, val.Skip, func(s string) { t.Fatal(s) })
	assertEq("limit", 1, val.Limit, func(s string) { t.Fatal(s) })
	assertEq("total", 1, val.Total, func(s string) { t.Fatal(s) })
	assertEq("val", "a", val.Val, func(s string) { t.Fatal(s) })

	val, err = rw.Read(context.Background())
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("skip", 0, val.Skip, func(s string) { t.Fatal(s) })
	assertEq("limit", 2, val.Limit, func(s string) { t.Fatal(s) })
	assertEq("total", 2, val.Total, func(s string) { t.Fatal(s) })
	assertEq("val", "b", val.Val, func(s string) { t.Fatal(s) })

	val, err = rw.Read(context.Background())
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })
	assertEq("skip", 0, val.Skip, func(s string) { t.Fatal(s) })
	assertEq("limit", 0, val.Limit, func(s string) { t.Fatal(s) })
	assertEq("total", 0, val.Total, func(s string) { t.Fatal(s) })
	assertEq("val", "", val.Val, func(s string) { t.Fatal(s) })
}

func TestNewContWriterWithNilReader(t *testing.T) {
	rw := core.NewReadWriterFrom[Paged[string]]()
	vw := NewContWriter(
		NewContWriterArgs[string]{
			Reader: nil,
			Writer: rw,
			Limit:  2,
		},
	)

	err := *new(error)
	val := Paged[string]{}

	// Write pages.
	err = vw.Write(context.Background(), "a")
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })

	// Read and validate.
	val, err = rw.Read(context.Background())
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })
	assertEq("skip", 0, val.Skip, func(s string) { t.Fatal(s) })
	assertEq("limit", 0, val.Limit, func(s string) { t.Fatal(s) })
	assertEq("total", 0, val.Total, func(s string) { t.Fatal(s) })
	assertEq("val", "", val.Val, func(s string) { t.Fatal(s) })
}

func TestNewContWriterWithNilWriter(t *testing.T) {
	vw := NewContWriter(
		NewContWriterArgs[string]{
			Reader: core.NewReaderFrom(1, 2),
			Writer: nil,
			Limit:  2,
		},
	)

	// Write pages.
	err := vw.Write(context.Background(), "a")
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })
}

func TestNewContWriterWithWriteErr(t *testing.T) {
	vw := NewContWriter(
		NewContWriterArgs[string]{
			Reader: core.NewReaderFrom(1, 2),
			Writer: core.WriterImpl[Paged[string]]{},
			Limit:  2,
		},
	)

	// Write pages.
	err := vw.Write(context.Background(), "a")
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })
}
