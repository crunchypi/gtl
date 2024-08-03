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

// -----------------------------------------------------------------------------
// New ReadWriteCloser iface + impl.
// -----------------------------------------------------------------------------

// ReadWriteCloser groups Reader[T] and Writer[U] with io.Closer.
type ReadWriteCloser[T, U any] interface {
	io.Closer
	Reader[T]
	Writer[U]
}

// ReadWriteCloserImpl lets you implement ReadWriteCloser with functions.
// This is similar to ReadWriterImpl but lets you implement io.Closer as well.
type ReadWriteCloserImpl[T, U any] struct {
	ImplC func() error
	ImplR func(context.Context) (T, error)
	ImplW func(context.Context, U) error
}

// Close implements io.Close by deferring to the internal ImplC func.
// If the internal ImplC func is nil, nothing will happen.
func (impl ReadWriteCloserImpl[T, U]) Close() (err error) {
	if impl.ImplC == nil {
		return
	}

	return impl.ImplC()
}

// Read implements Reader[T] by deferring logic to the internal ImplR func.
// If it's not set, an io.EOF is returned.
func (impl ReadWriteCloserImpl[T, U]) Read(ctx context.Context) (r T, err error) {
	if impl.ImplR == nil {
		err = io.EOF
		return
	}

	return impl.ImplR(ctx)
}

// Write implements Writer[U] by deferring logic to the internal ImplW func.
// If it's not set, an io.ErrClosedPipe is returned.
func (impl ReadWriteCloserImpl[T, U]) Write(ctx context.Context, v U) (err error) {
	if impl.ImplW == nil {
		err = io.ErrClosedPipe
		return
	}

	return impl.ImplW(ctx, v)
}
