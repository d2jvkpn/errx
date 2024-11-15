# errx
An extension of Golang's error type that implements error tracking and encoding (kind &amp; code).


#### C01. structure
```go
type ErrX struct {
	// internal data: error array
	errors []error

	// fields for identification and api response
	Kind string `json:"kind"`
	Code string `json:"code"`
	Msg  string `json:"msg"`

	// optinal tracer
	Caller string // fn::file::line
}
```
