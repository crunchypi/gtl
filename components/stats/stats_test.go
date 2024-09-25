package stats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/crunchypi/gtl/core"
)

var tvCtxKey = "testKey"
var tvCtxVal = "testVal"
var tvCtxMap = map[string]any{tvCtxKey: tvCtxVal}
var tvCtx = context.WithValue(context.Background(), tvCtxKey, tvCtxVal)
var tvErr = errors.New("test error")

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

func TestNewStreamedTeeReaderIdeal(t *testing.T) {
	rw := core.NewReadWriterFrom[StatsStreamed[string]]()

	r := NewStreamedTeeReader(
		NewStreamedTeeReaderArgs[string, string]{
			Reader:  core.NewReadWriterFrom("test1"),
			Writer:  rw,
			Tag:     "test",
			Fmt:     func(s string) string { return s },
			CtxKeys: []string{tvCtxKey},
		},
	)

	// Vars
	val := ""
	err := *new(error)
	stat := StatsStreamed[string]{}

	// Call: 1st
	val, err = r.Read(tvCtx)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test1", val, func(s string) { t.Fatal(s) })

	stat, err = rw.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("tag", stat.Tag, "test", func(s string) { t.Fatal(s) })
	assertEq("val", stat.Val, "test1", func(s string) { t.Fatal(s) })
	assertEq("ctx", stat.CtxMap, tvCtxMap, func(s string) { t.Fatal(s) })

	// Call: 2nd.
	val, err = r.Read(tvCtx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })

	stat, err = rw.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
}

func TestNewStreamedTeeReaderWithNilReader(t *testing.T) {
	rw := core.NewReadWriterFrom[StatsStreamed[string]]()

	r := NewStreamedTeeReader(
		NewStreamedTeeReaderArgs[string, string]{
			Reader:  nil,
			Writer:  rw,
			Tag:     "test",
			Fmt:     func(s string) string { return s },
			CtxKeys: []string{tvCtxKey},
		},
	)

	// Vars
	err := *new(error)

	// Call: 1st
	_, err = r.Read(tvCtx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })

	_, err = rw.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
}

func TestNewStreamedTeeReaderWithNilWriter(t *testing.T) {
	r := NewStreamedTeeReader(
		NewStreamedTeeReaderArgs[string, string]{
			Reader:  core.NewReadWriterFrom("test1"),
			Writer:  nil,
			Tag:     "test",
			Fmt:     func(s string) string { return s },
			CtxKeys: []string{tvCtxKey},
		},
	)

	// Vars
	val := ""
	err := *new(error)

	// Call: 1st
	val, err = r.Read(tvCtx)
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })
	assertEq("val", "test1", val, func(s string) { t.Fatal(s) })

	// Call: 2nd.
	val, err = r.Read(tvCtx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
}

func TestNewStreamedTeeReaderWithUnsetTag(t *testing.T) {
	rw := core.NewReadWriterFrom[StatsStreamed[string]]()

	r := NewStreamedTeeReader(
		NewStreamedTeeReaderArgs[string, string]{
			Reader:  core.NewReadWriterFrom("test1"),
			Writer:  rw,
			Tag:     "",
			Fmt:     func(s string) string { return s },
			CtxKeys: []string{tvCtxKey},
		},
	)

	// Vars
	val := ""
	err := *new(error)
	stat := StatsStreamed[string]{}

	// Call: 1st
	val, err = r.Read(tvCtx)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test1", val, func(s string) { t.Fatal(s) })

	stat, err = rw.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("tag", stat.Tag, "<unset>", func(s string) { t.Fatal(s) })
	assertEq("val", stat.Val, "test1", func(s string) { t.Fatal(s) })
	assertEq("ctx", stat.CtxMap, tvCtxMap, func(s string) { t.Fatal(s) })

	// Call: 2nd.
	val, err = r.Read(tvCtx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })

	stat, err = rw.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
}

func TestNewStreamedTeeReaderWithNilFmt(t *testing.T) {
	rw := core.NewReadWriterFrom[StatsStreamed[string]]()

	r := NewStreamedTeeReader(
		NewStreamedTeeReaderArgs[string, string]{
			Reader:  core.NewReadWriterFrom("test1"),
			Writer:  rw,
			Tag:     "test",
			Fmt:     nil,
			CtxKeys: []string{tvCtxKey},
		},
	)

	// Vars
	val := ""
	err := *new(error)
	stat := StatsStreamed[string]{}

	// Call: 1st
	val, err = r.Read(tvCtx)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test1", val, func(s string) { t.Fatal(s) })

	stat, err = rw.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("tag", stat.Tag, "test", func(s string) { t.Fatal(s) })
	assertEq("val", stat.Val, "", func(s string) { t.Fatal(s) })
	assertEq("ctx", stat.CtxMap, tvCtxMap, func(s string) { t.Fatal(s) })

	// Call: 2nd.
	val, err = r.Read(tvCtx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })

	stat, err = rw.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
}

func TestNewBatchedTeeReaderIdeal(t *testing.T) {
	rw := core.NewReadWriterFrom[StatsBatched]()

	r := NewBatchedTeeReader(
		NewBatchedTeeReaderArgs[string]{
			Reader:  core.NewReadWriterFrom([]string{"test1"}),
			Writer:  rw,
			Tag:     "test",
			CtxKeys: []string{tvCtxKey},
		},
	)

	// Vars.
	val := []string{}
	err := *new(error)
	stat := StatsBatched{}

	// Call: 1st
	val, err = r.Read(tvCtx)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", []string{"test1"}, val, func(s string) { t.Fatal(s) })

	stat, err = rw.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("tag", stat.Tag, "test", func(s string) { t.Fatal(s) })
	assertEq("len", stat.Len, 1, func(s string) { t.Fatal(s) })
	assertEq("ctx", stat.CtxMap, tvCtxMap, func(s string) { t.Fatal(s) })

	// Call: 2nd.
	val, err = r.Read(tvCtx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })

	stat, err = rw.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
}

func TestNewBatchedTeeReaderWithNilReader(t *testing.T) {
	rw := core.NewReadWriterFrom[StatsBatched]()

	r := NewBatchedTeeReader(
		NewBatchedTeeReaderArgs[string]{
			Reader:  nil,
			Writer:  rw,
			Tag:     "test",
			CtxKeys: []string{tvCtxKey},
		},
	)

	// Vars
	err := *new(error)

	// Call: 1st
	_, err = r.Read(tvCtx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })

	_, err = rw.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
}

func TestNewBatchedTeeReaderWithNilWriter(t *testing.T) {
	r := NewBatchedTeeReader(
		NewBatchedTeeReaderArgs[string]{
			Reader:  core.NewReadWriterFrom([]string{"test1"}),
			Writer:  nil,
			Tag:     "test",
			CtxKeys: []string{tvCtxKey},
		},
	)

	// Vars.
	val := []string{}
	err := *new(error)

	// Call: 1st
	val, err = r.Read(tvCtx)
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })
	assertEq("val", []string{"test1"}, val, func(s string) { t.Fatal(s) })

	// Call: 2nd.
	val, err = r.Read(tvCtx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
}

func TestNewBatchedTeeReaderWithUnsetTag(t *testing.T) {
	rw := core.NewReadWriterFrom[StatsBatched]()

	r := NewBatchedTeeReader(
		NewBatchedTeeReaderArgs[string]{
			Reader:  core.NewReadWriterFrom([]string{"test1"}),
			Writer:  rw,
			Tag:     "",
			CtxKeys: []string{tvCtxKey},
		},
	)

	// Vars.
	val := []string{}
	err := *new(error)
	stat := StatsBatched{}

	// Call: 1st
	val, err = r.Read(tvCtx)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", []string{"test1"}, val, func(s string) { t.Fatal(s) })

	stat, err = rw.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("tag", stat.Tag, "<unset>", func(s string) { t.Fatal(s) })
	assertEq("len", stat.Len, 1, func(s string) { t.Fatal(s) })
	assertEq("ctx", stat.CtxMap, tvCtxMap, func(s string) { t.Fatal(s) })

	// Call: 2nd.
	val, err = r.Read(tvCtx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })

	stat, err = rw.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
}
