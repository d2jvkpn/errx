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

	strs = make([]string, 0, 5)

	if self.Kind != "" {
		strs = append(strs, fmt.Sprintf("kind=%q", self.Kind))
	}

	if self.Code != "" {
		strs = append(strs, fmt.Sprintf("code=%q", self.Code))
	}

	if self.Msg != "" {
		strs = append(strs, fmt.Sprintf("msg=%q", self.Msg))
	}

	if self.Caller != "" {
		strs = append(strs, fmt.Sprintf("caller=%q", self.Caller))
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
