package core

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"testing"
)

func assertEq[T any](subject string, a T, b T, f func(string)) {
	if f == nil {
		return
	}

	ab, _ := json.Marshal(a)
	bb, _ := json.Marshal(b)

	as := string(ab)
	bs := string(bb)

	if as == bs {
		return
	}

	s := "unexpected '%v':\n\twant: '%v'\n\thave: '%v'\n"
	f(fmt.Sprintf(s, subject, as, bs))
}

func TestEncoderImplEncodeIdeal(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	enc := EncoderImpl{Impl: gob.NewEncoder(buf).Encode}

	err := enc.Encode("test")
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test", string(buf.Bytes()[4:]), func(s string) { t.Fatal(s) })
}

func TestEncoderImplEncodeWithNilImpl(t *testing.T) {
	enc := EncoderImpl{}

	err := enc.Encode("test")
	assertEq("err", io.ErrClosedPipe, err, func(s string) { t.Fatal(s) })
}
