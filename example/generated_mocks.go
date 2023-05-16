package example

import (
	"fmt"
)

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
