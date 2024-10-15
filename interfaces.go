package errx

import (
	"fmt"
	"strings"
)

type Error interface {
	Error() string
	IsNil() bool
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

func (self *ErrX) Error() string {
	var (
		strs    []string
		builder strings.Builder
	)

	if self == nil || len(self.errors) == 0 {
		return "<nil>"
	}

	strs = make([]string, 0, 7)

	// kind, code and msg
	if self.Kind != "" {
		strs = append(strs, fmt.Sprintf("kind=%q", self.Kind))
	}

	if self.Code != "" {
		strs = append(strs, fmt.Sprintf("code=%q", self.Code))
	}

	if self.Msg != "" {
		strs = append(strs, fmt.Sprintf("msg=%q", self.Msg))
	}

	// fn, file and line
	if self.fn != "" {
		strs = append(strs, fmt.Sprintf("fn=%q", self.fn))
	}

	if self.file != "" {
		strs = append(strs, fmt.Sprintf("file=%q", self.file))
	}

	if self.line > 0 {
		strs = append(strs, fmt.Sprintf("lint=%d", self.line))
	}

	// errors=
	// - "error1"
	// - "error2"
	strs = append(strs, "errors=")

	builder.Grow(64)
	builder.WriteString(strings.Join(strs, "; "))

	for _, bts := range self.MarshalErrors() {
		builder.WriteString("\n- ")
		builder.Write(bts)
	}

	return builder.String()
}
