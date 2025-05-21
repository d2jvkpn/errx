package errx

import (
	"fmt"
	"testing"
)

func TestPar(t *testing.T) {
	err := ParRun(
		func() *ErrX {
			return Eee("k1").WithCode("c1").WithCaller()
		},
		func() *ErrX {
			return Eee("k2").WithCode("c2")
		},
	)

	bts, _ := err.MarshalJSON()
	fmt.Printf("==> %s\n", bts)
}
