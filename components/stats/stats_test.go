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

func TestNewStreamedTeeWriterIdeal(t *testing.T) {
	rwv := core.NewReadWriterFrom[string]()
	rws := core.NewReadWriterFrom[StatsStreamed[string]]()

	w := NewStreamedTeeWriter(
		NewStreamedTeeWriterArgs[string, string]{
			WriterVals:  rwv,
			WriterStats: rws,
			Tag:         "test",
			Fmt:         func(v string) string { return v },
			CtxKeys:     []string{tvCtxKey},
		},
	)

	// Vars.
	val := ""
	err := *new(error)
	stats := StatsStreamed[string]{}

	// Call: 1st.
	err = w.Write(tvCtx, "test1")
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	// Eval: Val
	val, err = rwv.Read(tvCtx)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test1", val, func(s string) { t.Fatal(s) })

	// Eval: Stats.
	stats, err = rws.Read(tvCtx)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("tag", stats.Tag, "test", func(s string) { t.Fatal(s) })
	assertEq("val", stats.Val, "test1", func(s string) { t.Fatal(s) })
	assertEq("ctx", stats.CtxMap, tvCtxMap, func(s string) { t.Fatal(s) })

	// Eval empty.
	_, err = rwv.Read(tvCtx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })

	_, err = rws.Read(tvCtx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
}

func TestNewStreamedTeeWriterWithNilWriterOfVals(t *testing.T) {
	rws := core.NewReadWriterFrom[StatsStreamed[string]]()

	w := NewStreamedTeeWriter(
		NewStreamedTeeWriterArgs[string, string]{
			WriterVals:  nil,
			WriterStats: rws,
			Tag:         "test",
			Fmt:         func(v string) string { return v },
			CtxKeys:     []string{tvCtxKey},
		},
	)

	// Vars.
	err := *new(error)

	// Call: 1st.
	err = w.Write(tvCtx, "test1")
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })

	// Eval: Stats.
	_, err = rws.Read(tvCtx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
}

func TestNewStreamedTeeWriterWithNilWriterOfStats(t *testing.T) {
	rwv := core.NewReadWriterFrom[string]()

	w := NewStreamedTeeWriter(
		NewStreamedTeeWriterArgs[string, string]{
			WriterVals:  rwv,
			WriterStats: nil,
			Tag:         "test",
			Fmt:         func(v string) string { return v },
			CtxKeys:     []string{tvCtxKey},
		},
	)

	// Vars.
	val := ""
	err := *new(error)

	// Call: 1st.
	err = w.Write(tvCtx, "test1")
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })

	// Eval: Val
	val, err = rwv.Read(tvCtx)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test1", val, func(s string) { t.Fatal(s) })

	// Eval empty.
	_, err = rwv.Read(tvCtx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
}

func TestNewStreamedTeeWriterWithUnsetTag(t *testing.T) {
	rwv := core.NewReadWriterFrom[string]()
	rws := core.NewReadWriterFrom[StatsStreamed[string]]()

	w := NewStreamedTeeWriter(
		NewStreamedTeeWriterArgs[string, string]{
			WriterVals:  rwv,
			WriterStats: rws,
			Tag:         "",
			Fmt:         func(v string) string { return v },
			CtxKeys:     []string{tvCtxKey},
		},
	)

	// Vars.
	val := ""
	err := *new(error)
	stats := StatsStreamed[string]{}

	// Call: 1st.
	err = w.Write(tvCtx, "test1")
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	// Eval: Val
	val, err = rwv.Read(tvCtx)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test1", val, func(s string) { t.Fatal(s) })

	// Eval: Stats.
	stats, err = rws.Read(tvCtx)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("tag", stats.Tag, "<unset>", func(s string) { t.Fatal(s) })
	assertEq("val", stats.Val, "test1", func(s string) { t.Fatal(s) })
	assertEq("ctx", stats.CtxMap, tvCtxMap, func(s string) { t.Fatal(s) })

	// Eval empty.
	_, err = rwv.Read(tvCtx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })

	_, err = rws.Read(tvCtx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
}

func TestNewStreamedTeeWriterWithNilFmt(t *testing.T) {
	rwv := core.NewReadWriterFrom[string]()
	rws := core.NewReadWriterFrom[StatsStreamed[string]]()

	w := NewStreamedTeeWriter(
		NewStreamedTeeWriterArgs[string, string]{
			WriterVals:  rwv,
			WriterStats: rws,
			Tag:         "test",
			Fmt:         nil,
			CtxKeys:     []string{tvCtxKey},
		},
	)

	// Vars.
	val := ""
	err := *new(error)
	stats := StatsStreamed[string]{}

	// Call: 1st.
	err = w.Write(tvCtx, "test1")
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	// Eval: Val
	val, err = rwv.Read(tvCtx)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test1", val, func(s string) { t.Fatal(s) })

	// Eval: Stats.
	stats, err = rws.Read(tvCtx)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("tag", stats.Tag, "test", func(s string) { t.Fatal(s) })
	assertEq("val", stats.Val, "", func(s string) { t.Fatal(s) })
	assertEq("ctx", stats.CtxMap, tvCtxMap, func(s string) { t.Fatal(s) })

	// Eval empty.
	_, err = rwv.Read(tvCtx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })

	_, err = rws.Read(tvCtx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
}

func TestNewStreamedTeeWriterWithErrClosedPipe(t *testing.T) {
	rws := core.NewReadWriterFrom[StatsStreamed[string]]()

	w := NewStreamedTeeWriter(
		NewStreamedTeeWriterArgs[string, string]{
			WriterVals:  core.WriterImpl[string]{},
			WriterStats: rws,
			Tag:         "test",
			Fmt:         func(v string) string { return v },
			CtxKeys:     []string{tvCtxKey},
		},
	)

	// Vars.
	err := *new(error)

	// Call: 1st.
	err = w.Write(tvCtx, "test1")
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })

	_, err = rws.Read(tvCtx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
}

func TestNewBatchedTeeWriterIdeal(t *testing.T) {
	rwv := core.NewReadWriterFrom[[]string]()
	rws := core.NewReadWriterFrom[StatsBatched]()

	w := NewBatchedTeeWriter(
		NewBatchedTeeWriterArgs[string]{
			WriterVals:  rwv,
			WriterStats: rws,
			Tag:         "test",
			CtxKeys:     []string{tvCtxKey},
		},
	)

	// Vars.
	val := []string{}
	err := *new(error)
	stats := StatsBatched{}

	// Call: 1st.
	err = w.Write(tvCtx, []string{"test1"})
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	// Eval: Val
	val, err = rwv.Read(tvCtx)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", []string{"test1"}, val, func(s string) { t.Fatal(s) })

	// Eval: Stats.
	stats, err = rws.Read(tvCtx)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("tag", stats.Tag, "test", func(s string) { t.Fatal(s) })
	assertEq("len", stats.Len, 1, func(s string) { t.Fatal(s) })
	assertEq("ctx", stats.CtxMap, tvCtxMap, func(s string) { t.Fatal(s) })

	// Eval empty.
	_, err = rwv.Read(tvCtx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })

	_, err = rws.Read(tvCtx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
}

func TestNewBatchedTeeWriterWithNilWriterOfVals(t *testing.T) {
	rws := core.NewReadWriterFrom[StatsBatched]()

	w := NewBatchedTeeWriter(
		NewBatchedTeeWriterArgs[string]{
			WriterVals:  nil,
			WriterStats: rws,
			Tag:         "test",
			CtxKeys:     []string{tvCtxKey},
		},
	)

	// Vars.
	err := *new(error)

	// Call: 1st.
	err = w.Write(tvCtx, []string{"test1"})
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })

	// Eval: Stats.
	_, err = rws.Read(tvCtx)
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })
}

func TestNewBatchedTeeWriterWithNilWriterOfStats(t *testing.T) {
	rwv := core.NewReadWriterFrom[[]string]()

	w := NewBatchedTeeWriter(
		NewBatchedTeeWriterArgs[string]{
			WriterVals:  rwv,
			WriterStats: nil,
			Tag:         "test",
			CtxKeys:     []string{tvCtxKey},
		},
	)

	// Vars.
	val := []string{}
	err := *new(error)

	// Call: 1st.
	err = w.Write(tvCtx, []string{"test1"})
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })

	// Eval: Val
	val, err = rwv.Read(tvCtx)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", []string{"test1"}, val, func(s string) { t.Fatal(s) })

	// Eval empty.
	_, err = rwv.Read(tvCtx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
}

func TestNewBatchedTeeWriterWithUnsetTag(t *testing.T) {
	rwv := core.NewReadWriterFrom[[]string]()
	rws := core.NewReadWriterFrom[StatsBatched]()

	w := NewBatchedTeeWriter(
		NewBatchedTeeWriterArgs[string]{
			WriterVals:  rwv,
			WriterStats: rws,
			Tag:         "",
			CtxKeys:     []string{tvCtxKey},
		},
	)

	// Vars.
	val := []string{}
	err := *new(error)
	stats := StatsBatched{}

	// Call: 1st.
	err = w.Write(tvCtx, []string{"test1"})
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })

	// Eval: Val
	val, err = rwv.Read(tvCtx)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", []string{"test1"}, val, func(s string) { t.Fatal(s) })

	// Eval: Stats.
	stats, err = rws.Read(tvCtx)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("tag", stats.Tag, "<unset>", func(s string) { t.Fatal(s) })
	assertEq("len", stats.Len, 1, func(s string) { t.Fatal(s) })
	assertEq("ctx", stats.CtxMap, tvCtxMap, func(s string) { t.Fatal(s) })

	// Eval empty.
	_, err = rwv.Read(tvCtx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })

	_, err = rws.Read(tvCtx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
}

func TestNewBatchedTeeWriterWithErrClosedPipe(t *testing.T) {
	rws := core.NewReadWriterFrom[StatsBatched]()

	w := NewBatchedTeeWriter(
		NewBatchedTeeWriterArgs[string]{
			WriterVals:  core.WriterImpl[[]string]{},
			WriterStats: rws,
			Tag:         "test",
			CtxKeys:     []string{tvCtxKey},
		},
	)

	// Vars.
	err := *new(error)

	// Call: 1st.
	err = w.Write(tvCtx, []string{"test1"})
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })

	_, err = rws.Read(tvCtx)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
}
