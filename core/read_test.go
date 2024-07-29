package core

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"testing"
)

// -----------------------------------------------------------------------------
// Reader impl.
// -----------------------------------------------------------------------------

func TestReaderImplReadIdeal(t *testing.T) {
	r := ReaderImpl[int]{}
	r.Impl = func(ctx context.Context) (int, error) { return 1, nil }

	val, err := r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 1, val, func(s string) { t.Fatal(s) })
}

func TestReaderImplReadWithNilImpl(t *testing.T) {
	r := ReaderImpl[int]{}

	val, err := r.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

// -----------------------------------------------------------------------------
// ReadCloser impl.
// -----------------------------------------------------------------------------

func TestReadCloserImplReadIdeal(t *testing.T) {
	rc := ReadCloserImpl[int]{}
	rc.ImplR = func(ctx context.Context) (int, error) { return 1, nil }

	val, err := rc.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 1, val, func(s string) { t.Fatal(s) })
}

func TestReadCloserImplReadWithNilImpl(t *testing.T) {
	rc := ReadCloserImpl[int]{}

	val, err := rc.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

func TestReadCloserImplCloseIdeal(t *testing.T) {
	rc := ReadCloserImpl[int]{}
	rc.ImplC = func() error { return nil }

	err := rc.Close()
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
}

func TestReadCloserImplCloseWithNilImpl(t *testing.T) {
	rc := ReadCloserImpl[int]{}

	err := rc.Close()
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
}

// -----------------------------------------------------------------------------
// Constructors.
// -----------------------------------------------------------------------------

func TestNewReaderFromIdeal(t *testing.T) {
	r := NewReaderFrom(1, 2)

	err := *new(error)
	val := 0

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 1, val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", 2, val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", 0, val, func(s string) { t.Fatal(s) })
}

func TestNewReaderFromBytesIdeal(t *testing.T) {
	b := bytes.NewBuffer(nil)
	json.NewEncoder(b).Encode("test1")
	json.NewEncoder(b).Encode("test2")

	f := func(r io.Reader) Decoder { return json.NewDecoder(r) }
	r := NewReaderFromBytes[string](b)(f)

	err := *new(error)
	val := ""

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test1", val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test2", val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", "", val, func(s string) { t.Fatal(s) })
}

func TestNewReaderFromBytesWithNilReader(t *testing.T) {
	r := NewReaderFromBytes[string](nil)(nil)

	err := *new(error)
	val := ""

	val, err = r.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", "", val, func(s string) { t.Fatal(s) })
}

func TestNewReaderFromBytesWithNilDecoder(t *testing.T) {
	b := bytes.NewBuffer(nil)
	json.NewEncoder(b).Encode("test1")
	json.NewEncoder(b).Encode("test2")

	r := NewReaderFromBytes[string](b)(nil)

	err := *new(error)
	val := ""

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test1", val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", *new(error), err, func(s string) { t.Fatal(s) })
	assertEq("val", "test2", val, func(s string) { t.Fatal(s) })

	val, err = r.Read(nil)
	assertEq("err", io.EOF, err, func(s string) { t.Fatal(s) })
	assertEq("val", "", val, func(s string) { t.Fatal(s) })
}
