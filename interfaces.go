package errx

import (
// "fmt"
)

type Error interface {
	Error() string
	IsNil() bool
}
