package eventloop

import (
	"context"
	"testing"
	"time"

	"github.com/crunchypi/gtl/core"
)

func newWriterWithNop[T any]() core.Writer[T] {
	return core.WriterImpl[T]{
		Impl: func(ctx context.Context, v T) error {
			return nil
		},
	}
}

func newWriterWithSleep1s[T any](w core.Writer[T]) core.Writer[T] {
	return core.WriterImpl[T]{
		Impl: func(ctx context.Context, v T) error {
			time.Sleep(time.Second)
			return w.Write(ctx, v)
		},
	}
}

func TestNewIdeal(t *testing.T) {
	args := NewArgs[int]{}
	args.Ctx = context.Background()
	args.Reader = core.NewReaderFrom(1, 2, 3)
	args.Writer = newWriterWithNop[int]()

	ctx, _ := New(args)

	select {
	case <-ctx.Done():
	case <-time.After(time.Second * 3):
		t.Fatal("test hung")
	}
}

func TestNewWithNilCtx(t *testing.T) {
	args := NewArgs[int]{}
	args.Ctx = nil
	args.Reader = core.NewReaderFrom(1, 2, 3)
	args.Writer = newWriterWithNop[int]()

	ctx, _ := New(args)

	select {
	case <-ctx.Done():
	case <-time.After(time.Second * 3):
		t.Fatal("test hung")
	}
}

func TestNewWithNilReader(t *testing.T) {
	args := NewArgs[int]{}
	args.Ctx = context.Background()
	args.Reader = nil
	args.Writer = newWriterWithNop[int]()

	ctx, _ := New(args)

	select {
	case <-ctx.Done():
	case <-time.After(time.Second * 3):
		t.Fatal("test hung")
	}
}

func TestNewWithNilWriter(t *testing.T) {
	args := NewArgs[int]{}
	args.Ctx = context.Background()
	args.Reader = core.NewReaderFrom(1, 2, 3)
	args.Writer = nil

	ctx, _ := New(args)

	select {
	case <-ctx.Done():
	case <-time.After(time.Second * 3):
		t.Fatal("test hung")
	}
}

func TestNewWithReaderErr(t *testing.T) {
	args := NewArgs[int]{}
	args.Ctx = context.Background()
	args.Reader = core.ReaderImpl[int]{}
	args.Writer = newWriterWithNop[int]()

	ctx, _ := New(args)

	select {
	case <-ctx.Done():
	case <-time.After(time.Second * 3):
		t.Fatal("test hung")
	}
}

func TestNewWithWriterErr(t *testing.T) {
	args := NewArgs[int]{}
	args.Ctx = context.Background()
	args.Reader = core.NewReaderFrom(1, 2, 3)
	args.Writer = core.WriterImpl[int]{}

	ctx, _ := New(args)

	select {
	case <-ctx.Done():
	case <-time.After(time.Second * 3):
		t.Fatal("test hung")
	}
}

func TestNewWithCancel(t *testing.T) {
	args := NewArgs[int]{}
	args.Ctx = context.Background()
	args.Reader = core.NewReaderFrom(1, 2, 3)
	args.Writer = newWriterWithSleep1s(newWriterWithNop[int]())

	ctx, ctxCancel := New(args)
	ctxCancel()

	select {
	case <-ctx.Done():
	case <-time.After(time.Second * 3):
		t.Fatal("test hung")
	}
}
