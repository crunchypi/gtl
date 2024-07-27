package core

import "io"

// -----------------------------------------------------------------------------
// Encoder.
// -----------------------------------------------------------------------------

// Encoder encodes values into binary form. Some commonly used encoders are:
//   - json.NewEncoder(bytes.NewBuffer(nil))
//   - gob.NewEncoder(bytes.NewBuffer(nil))
type Encoder interface {
	Encode(e any) error
}

// EncoderImpl lets you implement Encoder with a function. Place it into "Impl"
// and it will be called by the "Encode" method.
//
// Example:
//
//	func myEncoder() Encoder {
//	    return EncoderImpl{
//	        Impl: func(e any) error {
//	            // Your code.
//	        }
//	    }
//	}
type EncoderImpl struct {
	Impl func(e any) error
}

// Encode implements Encoder by deferring to the internal "Impl" func.
// If the internal "Impl" is not set, an io.ErrClosedPipe will be returned.
func (impl EncoderImpl) Encode(e any) error {
	if impl.Impl == nil {
		return io.ErrClosedPipe
	}

	return impl.Impl(e)
}

// -----------------------------------------------------------------------------
// Decoder.
// -----------------------------------------------------------------------------

// Decoder decodes values from binary form. Some commonly used encoders are:
//   - json.NewDecoder(bytes.NewBuffer(nil))
//   - gob.NewDecoder(bytes.NewBuffer(nil))
type Decoder interface {
	Decode(e any) error
}

// DecoderImpl lets you implement Decoder with a function. Place it into "Impl"
// and it will be called by the "Decode" method.
//
// Example:
//
//	func myDecoder() Decoder {
//	    return DecoderImpl{
//	        Impl: func(d any) error {
//	            // Your code.
//	        }
//	    }
//	}
type DecoderImpl struct {
	Impl func(d any) error
}

// Decode implements Decoder by deferring to the internal "Impl" func.
// If the internal "Impl" is not set, an io.EOF will be returned.
func (impl DecoderImpl) Decode(d any) error {
	if impl.Impl == nil {
		return io.EOF
	}

	return impl.Impl(d)
}
