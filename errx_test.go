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
	err.WithKind("kind42").WithKind("kind_xx").WithCode("code42")

	fmt.Printf("==> b1. ErrX: %+#v\n", err)

	fmt.Printf("==> b2. errors: %v\n", err.errors)

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

	fmt.Printf("==> d3. ErrX: %v\n", err)

	bts, _ = json.Marshal(err)
	fmt.Printf("==> d3. json: %s\n", bts)

	e = testBizError(errors.New("account not found")).WithMsg("account not exists")
	err = ErrXFrom(e)
	err.WithErr(errors.New("sorry")).WithErr(nil)
	bts, _ = json.Marshal(err)
	fmt.Printf("==> d4. json: %s\n", bts)

	fmt.Printf("==> d5. debug: %s\n", err)
}

func fn01ErrX() (err *ErrX) {
	return nil
}

func fn02ErrX() (e error) {
	return fn01ErrX()
}

func testBizError(e error) (err *ErrX) {
	return NewErrX(e).Trace(2).WithKind("biz_error").WithCode("NotFound")
}

func TestErrx02(t *testing.T) {
	var (
		err1, err2 *ErrX
		bts        []byte
	)

	err1 = NewErrXxx(Code("NotFound")).WithKind("biz_error").Trace()

	fmt.Printf("==> 1. err1: %+#v\n", err1)

	err2 = NewErrX(errors.New("an error"), Code("DBError")).WithKind("internal_error").Trace()

	err1.WithErr(err2)
	fmt.Printf("==> 2. err1: %+#v\n", err1)

	bts, _ = json.Marshal(err1)
	fmt.Printf("==> 3. err1 json: %s\n", bts)

	fmt.Printf("==> 4a. err1 debug: %s\n", err1)
	fmt.Printf("==> 4b. err1 debug: %+v\n", err1)

	err1 = nil
	bts, _ = json.Marshal(err1)
	fmt.Printf("==> 5. err1 json: %s\n", bts)

	fmt.Printf("==> 6. err1 debug: %s\n", err1)
}

func TestErr03(t *testing.T) {
	var eE Error

	eE = NewErrXxx()
	fmt.Printf("==> 1. Error: %t, %v\n", eE.IsNil(), eE)

	eE = NewErrX(nil)
	fmt.Printf("==> 2. Error: %t, %v\n", eE.IsNil(), eE)
}

type _Err1 struct {
	str string
}

func (self _Err1) Error() string {
	return self.str
}

func TestErr04(t *testing.T) {
	err1 := _Err1{str: "_Err1"}
	// err2 := _Err1{str: "_Err2"}

	var e1 error
	e1 = err1

	var e2 _Err1

	fmt.Printf("==> 1. %t\n", errors.Is(e1, e2)) // false

	e1 = fmt.Errorf("...: %w", err1)
	fmt.Printf("==> 2. %t\n", errors.Is(e1, e2)) // false

	e2 = _Err1{str: "_Err1"}
	fmt.Printf("==> 3. %t\n", errors.Is(e1, e2)) // true

	var e3 _Err1
	fmt.Printf("==> 4. %t\n", errors.As(e1, &e3)) // true
	fmt.Printf("==> 5. %v\n", e3)

	e1 = errors.Join(err1, errors.New("an_error"))
	fmt.Printf("==> 6. %t\n", errors.As(e1, &e3)) // true
	fmt.Printf("==> 7. %v\n", e3)

	var e4 _Err1
	e1 = errors.Join(errors.New("an_error"), err1)
	fmt.Printf("==> 8. %t\n", errors.As(e1, &e4)) // true
	fmt.Printf("==> 9. %v\n", e4)
}
