package errx

import (
	"fmt"
	"testing"
)

func TestPar(t *testing.T) {
	err := ParRun(
		func() *ErrX {
			return Eee().WithCode("c1").WithKind("k1").WithCaller()
		},
		func() *ErrX {
			return Eee().WithCode("c2").WithKind("k2")
		},
	)

	bts, _ := err.MarshalJSON()
	fmt.Printf("==> %s\n", bts)
}
