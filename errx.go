package errx

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"path/filepath"
	"runtime"
)

type ErrX struct {
	errors []error

	Line int    `json:"line"`
	Fn   string `json:"fn"`
	File string `json:"file"`

	Kind string `json:"kind"`
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

type Option func(*ErrX)

func NewErrX(e error, options ...Option) (err *ErrX) {
	err = &ErrX{errors: make([]error, 0, 1)}

	if e != nil {
		err.errors = append(err.errors, e)
	}

	for _, opt := range options {
		opt(err)
	}

	return err
}

func Trace(skips ...int) Option {
	if len(skips) == 0 {
		skips = []int{1}
	}

	return func(self *ErrX) {
		self.Trace(skips...)
	}
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
func ErrXFrom(e error) (err *ErrX) {
	var ok bool

	if err, ok = e.(*ErrX); ok {
		return err
	}

	err = NewErrX(e)

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

	pc, self.File, self.Line, _ = runtime.Caller(skip)
	self.Fn = filepath.Base(runtime.FuncForPC(pc).Name())

	return self
}

func (self *ErrX) IsNil() bool {
	if self == nil {
		return true
	}

	return len(self.errors) == 0
}

func (self *ErrX) Error() string {
	if self.IsNil() {
		return "<nil>"
	}

	return errors.Join(self.errors...).Error()
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
		if e == nil {
			continue
		}

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

func (self *ErrX) MarshalJSON() ([]byte, error) {
	data := struct {
		Errors []json.RawMessage `json:"errors"`
		Line   int               `json:"line,omitempty"`
		Fn     string            `json:"fn,omitempty"`
		File   string            `json:"file,omitempty"`

		Kind string `json:"kind"`
		Code string `json:"code"`
		Msg  string `json:"msg"`
	}{
		Errors: self.MarshalErrors(),
		Line:   self.Line,
		Fn:     self.Fn,
		File:   self.File,

		Code: self.Code,
		Kind: self.Kind,
		Msg:  self.Msg,
	}

	return json.Marshal(data)
}

func (self *ErrX) ErrKC() (string, string) {
	return self.Kind, self.Code
}

func (self *ErrX) Debug() string {
	var (
		strs    []string
		builder strings.Builder
	)

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

	if self.Line > 0 {
		strs = append(strs, fmt.Sprintf("lint=%d", self.Line))
	}

	if self.Fn != "" {
		strs = append(strs, fmt.Sprintf("fn=%q", self.Fn))
	}

	if self.File != "" {
		strs = append(strs, fmt.Sprintf("file=%q", self.File))
	}

	builder.Grow(64)
	builder.WriteString(strings.Join(strs, "; "))

	builder.WriteString("\nerrors:")

	for _, bts := range self.MarshalErrors() {
		builder.WriteString("\n- ")
		builder.Write(bts)
	}

	return builder.String()
}

func ParallelRun(funcs ...func() *ErrX) (err *ErrX) {
	var (
		errs []*ErrX
		wg   sync.WaitGroup
	)

	errs = make([]*ErrX, len(funcs))
	wg.Add(len(funcs))

	for i := range funcs {
		go func(i int) {
			errs[i] = funcs[i]()
			wg.Done()
		}(i)
	}
	wg.Wait()

	for i := range errs {
		if errs[i] == nil {
			continue
		}

		if err == nil {
			err = errs[i]
		} else {
			err.WithErr(errs[i])
		}
	}

	return err
}

func ParallelRun2(funcs ...func() error) (err *ErrX) {
	var (
		hasError bool
		ok       bool
		errs     []error
		wg       sync.WaitGroup
	)

	errs = make([]error, len(funcs))
	wg.Add(len(funcs))

	for i := range funcs {
		go func(i int) {
			errs[i] = funcs[i]()
			wg.Done()
		}(i)
	}
	wg.Wait()

	hasError = false
	for i := range errs {
		if errs[i] == nil {
			continue
		}
		hasError = true

		if err != nil {
			if err, ok = errs[i].(*ErrX); ok {
				errs[i] = nil
			}
		}
	}

	if !hasError {
		return nil
	}

	if err == nil {
		err = NewErrX(errors.Join(errs...))
	} else {
		err.WithErr(errs...)
	}

	return err
}
