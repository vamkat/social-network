package commonerrors

import (
	"errors"
	"testing"
)

func Benchmark_NewError(b *testing.B) {
	err := errors.New("error not found")
	for i := 0; b.Loop(); i++ {
		New(ErrNotFound, err, i)
	}
}
