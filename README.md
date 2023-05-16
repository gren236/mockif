# Mockif

Short for "Mocking interfaces"!

Mockif generates a mock implementation for any interfaces found in Go package provided. It skips any test files (`*_test.go`).
For example, such interface:
```go
type Foo interface {
	Bar(a string, b []int) error
	Baz (a fmt.Scanner, b byte) []byte  
}
```
will yield this test implementation:
```go
type mockFoo struct {
	mBar func(a string, b []int) error
	mBaz func(a fmt.Scanner, b byte) []byte
}

func (fm mockFoo) Bar(a string, b []int) error {
	return fm.mBar(a, b)
}
func (fm mockFoo) Baz(a fmt.Scanner, b byte) []byte {
	return fm.mBaz(a, b)
}
```

Then it can be used to define the behavior of the mock implementation in vanilla Golang test suite.

For now, the `func` type is not supported for arguments and return parameters of methods.

## Usage

`mockif` tool is able to process 2 positional (not named) arguments:

1. Path to the directory that holds the package. Can be either relative or absolute:
```shell
mockif ./example
```
2. Name of the file with generated implementations (`mocks.go` by default):
```shell
mockif ./example mocks_test.go
```

## Installation

```shell
go install github.com/gren236/mockif@latest
```
