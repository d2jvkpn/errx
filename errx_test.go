package errx

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

func TestErrx01(t *testing.T) {
	// 1.
	var e error

	e = fmt.Errorf("hello")
	e = fmt.Errorf("an error: %w", e)
	fmt.Printf("==> a1. %v\n", e)
	fmt.Printf("==> a2. %v\n", errors.Unwrap(e))

	e1, e2 := errors.New("hello"), errors.New("world")
	e = errors.Join(e1, e2)
	fmt.Printf("==> a3. %v\n", e)
	fmt.Printf("==> a4. %v\n", errors.Unwrap(e))

	// 2.
	var err *ErrX

	// errx = new(ErrX)
	err = NewErrX(errors.New("wrong"))
	err.WithCode("code42").WithKind("kind42").WithKind("kind_xx")

	fmt.Printf("==> b1. ErrX=%+#v\n", err)

	fmt.Printf("==> b2. errors=%v\n", err.Errors)

	// 3.
	err = fn01ErrX()
	fmt.Printf("==> c1. err is nil: %t\n", err == nil)

	e = fn02ErrX()
	err, _ = e.(*ErrX)
	fmt.Printf("==> c2. %t, %t\n", e == nil, err.IsNil())
	// false, true, true

	err = NewErrX(nil)
	e = err
	fmt.Printf("==> c3. is_nil=%t, e=%v\n", err.IsNil(), e)

	// 4.
	var bts []byte

	err = NewErrX(errors.New("e1"))
	err.WithErr(errors.New("e2")).WithKind("kind01").Trace()

	fmt.Printf("==> d3. ErrX=%v\n", err)

	bts, _ = json.Marshal(err)
	fmt.Printf("==> d3. json=%s\n", bts)

	e = testBizError(errors.New("account not found")).WithMsg("account not exists")
	err = ErrXFrom(e)
	err.WithErr(errors.New("sorry")).WithErr(nil)
	bts, _ = json.Marshal(err)
	fmt.Printf("==> d4. json=%s\n", bts)

	fmt.Printf("==> d5. debug=%s\n", err.Debug())
}

func fn01ErrX() (err *ErrX) {
	return nil
}

func fn02ErrX() (e error) {
	return fn01ErrX()
}

func testBizError(e error) (err *ErrX) {
	return NewErrX(e).Trace(2).WithCode("Biz").WithKind("NotFound")
}

func TestErrx02(t *testing.T) {
	var errx, e2 *ErrX

	errx = NewErrX(errors.New("..."), Code("biz_error")).WithKind("NotFound").Trace()

	fmt.Printf("==> 1. %+#v\n", errx)

	e2 = NewErrX(errors.New("an error"), Code("internal_error")).WithKind("DBError").Trace()

	errx.WithErr(e2)
	fmt.Printf("==> 2. %+#v\n", errx)

	bts, _ := json.Marshal(errx)
	fmt.Printf("==> 2. %s\n", bts)
}
