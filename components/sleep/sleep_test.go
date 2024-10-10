package sleep

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/crunchypi/gtl/core"
)

var tvVerbose = false
var tvDuration = time.Millisecond * 100
var tvCtx = context.Background()

func tfNewRandomReader[T any](r core.Reader[T], d time.Duration) core.Reader[T] {
	return core.ReaderImpl[T]{
		Impl: func(ctx context.Context) (T, error) {
			time.Sleep(time.Duration(rand.Intn(int(d))))
			return r.Read(ctx)
		},
	}
}

// -----------------------------------------------------------------------------
// Tests for: NewStaticReader
// -----------------------------------------------------------------------------

func TestNewStaticReaderIdeal(t *testing.T) {
	vr := core.NewReaderFrom(1, 2, 3)
	sr := NewStaticReader(NewStaticReaderArgs[int]{vr, tvDuration})

	ts := time.Now()
	for _, err := sr.Read(tvCtx); err == nil; _, err = sr.Read(tvCtx) {
	}

	if tvVerbose {
		t.Log(time.Since(ts))
	}
}

func TestNewStaticReaderWithNilReader(t *testing.T) {
	sr := NewStaticReader(NewStaticReaderArgs[int]{Delay: tvDuration})

	ts := time.Now()
	for _, err := sr.Read(tvCtx); err == nil; _, err = sr.Read(tvCtx) {
	}

	if tvVerbose {
		t.Log(time.Since(ts))
	}
}

func TestNewStaticReaderWithNegativeDuration(t *testing.T) {
	vr := core.NewReaderFrom(1, 2, 3)
	sr := NewStaticReader(NewStaticReaderArgs[int]{vr, -tvDuration})

	ts := time.Now()
	for _, err := sr.Read(tvCtx); err == nil; _, err = sr.Read(tvCtx) {
	}

	if tvVerbose {
		t.Log(time.Since(ts))
	}
}

func TestNewStaticReaderWithNilCtx(t *testing.T) {
	vr := core.NewReaderFrom(1, 2, 3)
	sr := NewStaticReader(NewStaticReaderArgs[int]{vr, tvDuration})

	ts := time.Now()
	for _, err := sr.Read(nil); err == nil; _, err = sr.Read(nil) {
	}

	if tvVerbose {
		t.Log(time.Since(ts))
	}
}

// -----------------------------------------------------------------------------
// Tests for: NewDynamicReader
// -----------------------------------------------------------------------------

func TestNewDynamicReaderIdeal(t *testing.T) {
	vr := core.NewReaderFrom(1, 2, 3)
	fr := tfNewRandomReader(vr, tvDuration/3)
	sr := NewDynamicReader(NewDynamicReaderArgs[int]{fr, tvDuration})

	ts := time.Now()
	for _, err := sr.Read(tvCtx); err == nil; _, err = sr.Read(tvCtx) {
	}

	if tvVerbose {
		t.Log(time.Since(ts))
	}
}

func TestNewDynamicReaderWithNilReader(t *testing.T) {
	sr := NewDynamicReader(NewDynamicReaderArgs[int]{Delay: tvDuration})

	ts := time.Now()
	for _, err := sr.Read(tvCtx); err == nil; _, err = sr.Read(tvCtx) {
	}

	if tvVerbose {
		t.Log(time.Since(ts))
	}
}

func TestNewDynamicReaderWithNegativeDuration(t *testing.T) {
	vr := core.NewReaderFrom(1, 2, 3)
	fr := tfNewRandomReader(vr, tvDuration/3)
	sr := NewDynamicReader(NewDynamicReaderArgs[int]{fr, -tvDuration})

	ts := time.Now()
	for _, err := sr.Read(tvCtx); err == nil; _, err = sr.Read(tvCtx) {
	}

	if tvVerbose {
		t.Log(time.Since(ts))
	}
}

func TestNewDynamicReaderWithNilCtx(t *testing.T) {
	vr := core.NewReaderFrom(1, 2, 3)
	fr := tfNewRandomReader(vr, tvDuration/3)
	sr := NewDynamicReader(NewDynamicReaderArgs[int]{fr, tvDuration})

	ts := time.Now()
	for _, err := sr.Read(nil); err == nil; _, err = sr.Read(nil) {
	}

	if tvVerbose {
		t.Log(time.Since(ts))
	}
}

func TestNewDynamicReaderWithBounds(t *testing.T) {
	vr := core.NewReaderFrom(1, 2, 3)
	fr := tfNewRandomReader(vr, tvDuration/3)
	sr := NewDynamicReader(NewDynamicReaderArgs[int]{fr, tvDuration})

	ts := time.Now()
	ctx := context.WithValue(tvCtx, "bounds", 3)
	for _, err := sr.Read(ctx); err == nil; _, err = sr.Read(ctx) {
	}

	if tvVerbose {
		t.Log(time.Since(ts))
	}
}
