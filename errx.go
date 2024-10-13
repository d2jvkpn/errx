package errx

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"path/filepath"
	"runtime"
)

type ErrX struct {
	Kind string `json:"kind"`
	Code string `json:"code"`
	Msg  string `json:"msg"`

	errors []error
	fn     string
	file   string
	line   int
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

func Kind(str string) Option {
	return func(self *ErrX) {
		self.WithKind(str)
	}
}

func Code(str string) Option {
	return func(self *ErrX) {
		self.WithCode(str)
	}
}

func Msg(str string) Option {
	return func(self *ErrX) {
		self.WithMsg(str)
	}
}

// checks if the input is an ErrX
func ErrXFrom(e error, options ...Option) (err *ErrX) {
	var ok bool

	if e == nil {
		return nil
	}

	if err, ok = e.(*ErrX); !ok {
		err = NewErrX(e)
	}

	for _, opt := range options {
		opt(err)
	}

	return err
}

func (self *ErrX) Trace(skips ...int) *ErrX {
	var (
		skip int
		pc   uintptr
	)

	if skip = 1; len(skips) > 0 {
		skip = skips[0]
	}

	pc, self.file, self.line, _ = runtime.Caller(skip)
	self.fn = filepath.Base(runtime.FuncForPC(pc).Name())

	return self
}

func (self *ErrX) IsNil() bool {
	if self == nil {
		return true
	}

	return len(self.errors) == 0
}

/*
func (self *ErrX) Error() string {
	if self.IsNil() {
		return "<nil>"
	}

	return errors.Join(self.errors...).Error()
}
*/

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
		Fn     string            `json:"fn,omitempty"`
		File   string            `json:"file,omitempty"`
		Line   int               `json:"line,omitempty"`
	}{
		Kind: self.Kind,
		Code: self.Code,
		Msg:  self.Msg,

		Errors: self.MarshalErrors(),
		Fn:     self.fn,
		File:   self.file,
		Line:   self.line,
	}

	return json.Marshal(data)
}

func (self *ErrX) ErrKC() (string, string) {
	return self.Kind, self.Code
}

func (self *ErrX) Error() string {
	var (
		strs    []string
		builder strings.Builder
	)

	if self == nil {
		return "<nil>"
	}

	strs = make([]string, 0, 6)

	if self.Kind != "" {
		strs = append(strs, fmt.Sprintf("kind=%q", self.Kind))
	}

	if self.Code != "" {
		strs = append(strs, fmt.Sprintf("code=%q", self.Code))
	}

	if self.Msg != "" {
		strs = append(strs, fmt.Sprintf("msg=%q", self.Msg))
	}

	if self.fn != "" {
		strs = append(strs, fmt.Sprintf("fn=%q", self.fn))
	}

	if self.file != "" {
		strs = append(strs, fmt.Sprintf("file=%q", self.file))
	}

	if self.line > 0 {
		strs = append(strs, fmt.Sprintf("lint=%d", self.line))
	}

	strs = append(strs, "errors=")

	builder.Grow(64)
	builder.WriteString(strings.Join(strs, "; "))

	for _, bts := range self.MarshalErrors() {
		builder.WriteString("\n- ")
		builder.Write(bts)
	}

	return builder.String()
}

func (self *ErrX) As(target any) bool {
	for i := range self.errors {
		if errors.As(self.errors[i], target) {
			return true
		}
	}

	return false
}
