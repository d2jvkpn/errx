package errx

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrx(t *testing.T) {
	// 1.
	var err error

	err = fmt.Errorf("hello")
	err = fmt.Errorf("an error: %w", err)
	fmt.Printf("==> a1. %v\n", err)
	fmt.Printf("==> a2. %v\n", errors.Unwrap(err))

	e1, e2 := errors.New("hello"), errors.New("world")
	err = errors.Join(e1, e2)
	fmt.Printf("==> a3. %v\n", err)
	fmt.Printf("==> a4. %v\n", errors.Unwrap(err))

	// 2.
	var (
		ek *ErrKind
		ec *ErrCode
	)

	err = &ErrKind{Kind: "Fn1"}
	fmt.Printf("==> b1. %v\n", err)
	fmt.Printf("==> b2. %t\n", errors.As(err, &ek))
	fmt.Printf("==> b3. %v\n", ek)

	ec = &ErrCode{Code: "invlaid"}
	err = errors.Join(err, ec)

	fmt.Printf("==> b4. %v\n", err)

	// 3.
	var (
		kind *ErrKind
		code *ErrCode
		errx *ErrX
	)

	_ = errors.As(err, &kind)
	_ = errors.As(err, &code)

	fmt.Printf("==> c1. code=%v, kind=%v\n", code, kind)

	// errx = new(ErrX)
	errx = NewErrX(errors.New("wrong"))
	errx.WithCode("code42").WithKind("kind42").WithKind("kind_xx")

	fmt.Printf("==> c2. ErrX=%+#v\n", errx)

	fmt.Printf("==> c3. code=%s, kind=%s, msg=%s\n", errx.GetCode(), errx.GetKind(), errx.GetMsg())

	fmt.Printf("==> c4. raw_errors=%v\n", errx.GetRawErrors())

	// 4.
	errx = Fn01ErrX()
	fmt.Printf("==> d1. errx is nil: %t\n", errx == nil)

	var e error = Fn02ErrX()
	errx, _ = e.(*ErrX)
	fmt.Printf("==> d2. %t, %t, %t\n", e == nil, errx.IsNil(), len(errx.GetRawErrors()) == 0)
	// false, true, true

	// 4.
	var bts []byte

	errx = NewErrX(errors.New("e1"))
	errx.WithRaw(errors.New("e2")).WithKind("kind01").Trace()

	fmt.Printf("==> d3. errx=%v\n", errx)

	bts, _ = errx.MarshalJSON()
	fmt.Printf("==> d3. json=%s\n", bts)

	err = testBizError(errors.New("account not found")).WithMsg("account not exists")
	errx = ErrXFrom(err)
	errx.WithRaw(errors.New("sorry")).WithRaw(nil)
	bts, _ = errx.MarshalJSON()
	fmt.Printf("==> d4. json=%s\n", bts)

	fmt.Printf("==> d5. respone=%s\n", errx.Response())
	fmt.Printf("==> d5. debug=%s\n", errx.Debug())
}

func Fn01ErrX() (errx *ErrX) {
	return nil
}

func Fn02ErrX() (err error) {
	return Fn01ErrX()
}

func testBizError(e error) (errx *ErrX) {
	return NewErrX(e).Trace(2).WithCode("Biz").WithKind("NotFound")
}
