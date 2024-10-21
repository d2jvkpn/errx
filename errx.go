package errx

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
)

type ErrX struct {
	Kind string `json:"kind"`
	Code string `json:"code"`
	Msg  string `json:"msg"`

	errors []error
	Caller string // fn::file::line
}

type Option func(*ErrX)

func NewErrX(e error, options ...Option) (err *ErrX) {
	if e == nil {
		return nil
	}

	err = &ErrX{errors: []error{e}}

	for _, opt := range options {
		opt(err)
	}

	return err
}

func NewErrXxx(options ...Option) (err *ErrX) {
	err = &ErrX{errors: []error{errors.New("...")}}

	for _, opt := range options {
		opt(err)
	}

	return err
}

/*
func Eee() error {
	return errors.New("...")
}
*/

func Kind(str string) Option {
	return func(self *ErrX) {
		self.Kind = str
	}
}

func Code(str string) Option {
	return func(self *ErrX) {
		self.Code = str
	}
}

func Msg(str string) Option {
	return func(self *ErrX) {
		self.Msg = str
	}
}

// checks if the input is an ErrX
func ErrXFrom(e error, xs ...bool) (err *ErrX, ok bool) {
	if e == nil {
		return nil, false
	}

	if len(xs) > 0 && xs[0] {
		if err, ok = e.(*ErrX); !ok {
			err = NewErrX(e)
		}
	} else {
		err = NewErrX(e)
	}

	return err, ok
}

func (self *ErrX) Apply(options ...Option) *ErrX {
	for i := range options {
		options[i](self)
	}

	return self
}

func (self *ErrX) WithCaller(skips ...int) *ErrX {
	var (
		skip int = 1
		pc   uintptr
	)

	if len(skips) > 0 {
		skip = skips[0]
	}

	pc, file, line, ok := runtime.Caller(skip)

	if ok {
		self.Caller = fmt.Sprintf(
			"%s::%s::%d",
			filepath.Base(runtime.FuncForPC(pc).Name()),
			file,
			line,
		)
	}

	return self
}

func (self *ErrX) WithErr(errs ...error) *ErrX {
	for i := range errs {
		if errs[i] != nil {
			self.errors = append(self.errors, errs[i])
		}
	}

	return self
}

func (self *ErrX) WithKind(str string) *ErrX {
	self.Kind = str

	return self
}

func (self *ErrX) WithCode(str string) *ErrX {
	self.Code = str

	return self
}

func (self *ErrX) WithMsg(str string) *ErrX {
	self.Msg = str

	return self
}

func (self *ErrX) MarshalErrors() (msgs []json.RawMessage) {
	var (
		ok  bool
		e   error
		msg json.RawMessage
	)

	for _, e = range self.errors {
		// data.Errors = append(data.Errors, fmt.Sprintf("%v", e))
		if _, ok = e.(*ErrX); ok {
			msg, _ = json.Marshal(&e)
		} else {
			msg, _ = json.Marshal(e.Error())
		}

		msgs = append(msgs, msg)
	}

	return msgs
}

func (self *ErrX) NumberOfErrors() int {
	if self == nil {
		return 0
	}

	return len(self.errors)
}

func (self ErrX) MarshalJSON() ([]byte, error) {
	data := struct {
		Kind string `json:"kind"`
		Code string `json:"code"`
		Msg  string `json:"msg"`

		Errors []json.RawMessage `json:"errors"`
		Caller string            `json:"caller,omitempty"`
	}{
		Kind: self.Kind,
		Code: self.Code,
		Msg:  self.Msg,

		Errors: self.MarshalErrors(),
		Caller: self.Caller,
	}

	return json.Marshal(data)
}

// iterate errors.As(self.errors[i], target): target = New(CustomError)
func (self *ErrX) As(target any) bool {
	for i := range self.errors {
		if errors.As(self.errors[i], target) {
			return true
		}
	}

	return false
}

// iterate errors.Is(self.errors[i], target): target = pkg.ErrorNotFound
func (self *ErrX) Is(e error) bool {
	for i := range self.errors {
		if errors.Is(self.errors[i], e) {
			return true
		}
	}

	return false
}

// copy errors of the *ErrX
func (self *ErrX) CopyErrors() (errs []error) {
	errs = make([]error, len(self.errors))

	for i := range self.errors {
		errs[i] = self.errors[i]
	}

	return errs
}

// compare the kind and code of two *ErrX
func (self *ErrX) Equals(other *ErrX) bool {
	switch {
	case self == nil && other == nil:
		return true
	case self != nil && other != nil:
	default:
		return false
	}

	return self.Kind == other.Kind && self.Code == other.Code
}
