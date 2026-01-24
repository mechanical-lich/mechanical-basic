package main

import (
	"github.com/mechanical-lich/mechanical-basic/pkg/basic"
)

func main() {
	mBasic := basic.NewMechanicalBasic()

	code := `
	function foobar():
		print "Inside foobar"
	endfunction

	function add(a, b):
		return a + b
	endfunction

	print "Outside foobar"
	foobar()
	print "After foobar"

	result = add(5.5, 7)
	print "Result of add(5.5, 7): " + result
	`

	err := mBasic.Run(code)
	if err != nil {
		panic(err)
	}
}
