package core

import (
	"context"
	"io"
)

// -----------------------------------------------------------------------------
// New Reader iface + impl.
// -----------------------------------------------------------------------------

// Reader reads T, it is intended as a generic variant of io.Reader.
type Reader[T any] interface {
	Read(context.Context) (T, error)
}

// ReaderImpl lets you implement Reader with a function. Place it into "impl"
// and it will be called by the "Read" method.
//
// Example:
//
//	func myReader() Reader[int] {
//	    return ReaderImpl[int]{
//	        Impl: func(ctx context.Context) (int, error) {
//	            // Your implementation.
//	        },
//	    }
//	}
type ReaderImpl[T any] struct {
	Impl func(context.Context) (T, error)
}

// Read implements Reader by deferring to the internal "Impl" func.
// If the internal "Impl" is not set, an io.EOF will be returned.
func (impl ReaderImpl[T]) Read(ctx context.Context) (r T, err error) {
	if impl.Impl == nil {
		err = io.EOF
		return
	}

	return impl.Impl(ctx)
}
