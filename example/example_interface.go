package example

import "fmt"

type Foo interface {
	Bar(a string, b []int) error
	Baz(a fmt.Scanner, b byte) []byte
}
