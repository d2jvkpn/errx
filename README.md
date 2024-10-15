# errx
An extension of Golang's error type that implements error tracking and encoding (kind &amp; code).


#### C01. structure
```go
type ErrX struct {
	Kind string `json:"kind"`
	Code string `json:"code"`
	Msg  string `json:"msg"`

	errors []error
	fn     string
	file   string
	line   int
}
```
