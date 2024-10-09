package sleep

import (
	"context"
	"testing"
	"time"

	"github.com/crunchypi/gtl/core"
)

var tvVerbose = false
var tvDuration = time.Millisecond * 100
var tvCtx = context.Background()

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
