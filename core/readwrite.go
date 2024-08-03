package core

import (
	"context"
	"io"
)

// -----------------------------------------------------------------------------
// New ReadWriter iface + impl.
// -----------------------------------------------------------------------------

// ReadWriter groups Reader[T] and Writer[U].
type ReadWriter[T, U any] interface {
	Reader[T]
	Writer[U]
}

// ReadWriterImpl lets you implement ReadWriter with functions. This is
// equivalent to using ReaderImpl and WriterImpl combined (see docs).
type ReadWriterImpl[T, U any] struct {
	ImplR func(context.Context) (T, error)
	ImplW func(context.Context, U) error
}

// Read implements the Reader[T] part of ReadWriter[T, U] by deferring logic
// to the internal ImplR func. If it's not set, an io.EOF is returned.
func (impl ReadWriterImpl[T, U]) Read(ctx context.Context) (r T, err error) {
	if impl.ImplR == nil {
		err = io.EOF
		return
	}

	return impl.ImplR(ctx)
}

// Write implements the Writer[U] part of ReadWriter[T, U] by deferring logic
// to the internal ImplW func. If it's not set, an io.ErrClosedPipe is returned.
func (impl ReadWriterImpl[T, U]) Write(ctx context.Context, v U) (err error) {
	if impl.ImplW == nil {
		err = io.ErrClosedPipe
		return
	}

	return impl.ImplW(ctx, v)
}
