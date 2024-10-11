package errx

import (
	"encoding/json"
	"errors"
	"fmt"

	"path/filepath"
	"runtime"
)

// #### 1. ErrX
type ErrX struct {
	Item error `json:"-"`

	Line int    `json:"line,omitempty"`
	Fn   string `json:"fn,omitempty"`
	File string `json:"file,omitempty"`
}

func NewErrX(e error) (errx *ErrX) {
	return &ErrX{Item: &ErrRaw{Errors: []error{e}}}
}

// checks if e is an ErrX
func ErrXFrom(e error) (errx *ErrX) {
	var ok bool

	if errx, ok = e.(*ErrX); ok {
		return errx
	}

	errx = &ErrX{Item: NewErrRaw(e)}

	return errx
}

func (self *ErrX) Trace(skips ...int) *ErrX {
	var (
		skip int
		pc   uintptr
	)

	if skip = 1; len(skips) > 0 {
		skip = skips[0]
	}

	pc, self.File, self.Line, _ = runtime.Caller(skip)
	self.Fn = filepath.Base(runtime.FuncForPC(pc).Name())

	return self
}

func (self *ErrX) IsNil() bool {
	if self == nil {
		return true
	}

	if errs := self.GetRawErrors(); len(errs) == 0 {
		return true
	}

	return self.Item == nil
}

func (self *ErrX) Error() string {
	if self.IsNil() {
		return "<nil>"
	}

	return self.Item.Error()
}

// #### 2. ErrRaw
type ErrRaw struct {
	Errors []error `json:"errors"`
}

func NewErrRaw(e error) (er *ErrRaw) {
	er = &ErrRaw{Errors: make([]error, 0, 1)}

	if e != nil {
		er.Errors = append(er.Errors, e)
	}

	return er
}

func (self *ErrRaw) Error() string {
	if len(self.Errors) == 0 {
		return "<nil>"
	}

	return errors.Join(self.Errors...).Error()
}

func (self *ErrRaw) AddErr(e error) *ErrRaw {
	if e != nil {
		self.Errors = append(self.Errors, e)
	}

	return self
}

func (self *ErrX) WithRaw(e error) *ErrX {
	var er *ErrRaw

	if errors.As(self.Item, &er) {
		er.AddErr(e)
	} else {
		self.Item = errors.Join(self.Item, NewErrRaw(e))
	}

	return self
}

func (self *ErrX) GetRawErrors() []error {
	var er *ErrRaw

	if self == nil || self.Item == nil {
		return nil
	}

	if errors.As(self.Item, &er) {
		return er.Errors
	} else {
		return nil
	}
}

// #### 3. ErrCode
type ErrCode struct {
	Code string `json:"code"`
}

func (self *ErrCode) Error() string {
	return self.Code
}

func (self *ErrX) WithCode(str string) *ErrX {
	var ec *ErrCode

	if errors.As(self.Item, &ec) {
		ec.Code = str
	} else {
		self.Item = errors.Join(self.Item, &ErrCode{Code: str})
	}

	return self
}

func (self *ErrX) GetCode() string {
	var ec *ErrCode

	if errors.As(self.Item, &ec) {
		return ec.Code
	} else {
		return ""
	}
}

// #### 4. ErrKind
type ErrKind struct {
	Kind string `json:"kind"`
}

func (self *ErrKind) Error() string {
	return self.Kind
}

func (self *ErrX) WithKind(str string) *ErrX {
	var ek *ErrKind

	if errors.As(self.Item, &ek) {
		ek.Kind = str
	} else {
		self.Item = errors.Join(self.Item, &ErrKind{Kind: str})
	}

	return self
}

func (self *ErrX) GetKind() string {
	var ek *ErrKind

	if errors.As(self.Item, &ek) {
		return ek.Kind
	} else {
		return ""
	}
}

// #### 5. ErrMsg
type ErrMsg struct {
	Msg string `json:"msg"`
}

func (self *ErrMsg) Error() string {
	return self.Msg
}

func (self *ErrX) WithMsg(str string) *ErrX {
	var em *ErrMsg

	if errors.As(self.Item, &em) {
		em.Msg = str
	} else {
		self.Item = errors.Join(self.Item, &ErrMsg{Msg: str})
	}

	return self
}

func (self *ErrX) GetMsg() string {
	var em *ErrMsg

	if errors.As(self.Item, &em) {
		return em.Msg
	} else {
		return ""
	}
}

// #### 5. JSON
func (self *ErrX) MarshalJSON() (bts []byte, e error) {
	data := struct {
		Errors []string `json:"errors"`

		Code string `json:"code,omitempty"`
		Kind string `json:"kind,omitempty"`
		Msg  string `json:"msg,omitempty"`

		Line int    `json:"line,omitempty"`
		Fn   string `json:"fn,omitempty"`
		File string `json:"file,omitempty"`
	}{
		Code: self.GetCode(),
		Kind: self.GetKind(),
		Msg:  self.GetMsg(),

		Line: self.Line,
		Fn:   self.Fn,
		File: self.File,
	}

	for _, e := range self.GetRawErrors() {
		data.Errors = append(data.Errors, fmt.Sprintf("%v", e))
	}

	return json.Marshal(data)
}

func (self *ErrX) Response() (bts json.RawMessage) {
	data := struct {
		Code string `json:"code,omitempty"`
		Kind string `json:"kind,omitempty"`
		Msg  string `json:"msg,omitempty"`
	}{
		Code: self.GetCode(),
		Kind: self.GetKind(),
		Msg:  self.GetMsg(),
	}

	bts, _ = json.Marshal(data)
	return bts
}

func (self *ErrX) Debug() (bts json.RawMessage) {
	data := struct {
		Errors []string `json:"errors"`

		Line int    `json:"line,omitempty"`
		Fn   string `json:"fn,omitempty"`
		File string `json:"file,omitempty"`
	}{
		Line: self.Line,
		Fn:   self.Fn,
		File: self.File,
	}

	for _, e := range self.GetRawErrors() {
		data.Errors = append(data.Errors, fmt.Sprintf("%v", e))
	}

	bts, _ = json.Marshal(data)
	return bts
}
