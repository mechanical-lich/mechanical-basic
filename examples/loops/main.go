package main

import (
	"github.com/mechanical-lich/mechanical-basic/pkg/basic"
)

func main() {
	mBasic := basic.NewMechanicalBasic()

	code := `
	let x = 10
	print "Starting: " + x
	for i = 1 to 5
		x = x + 1
		print "Loop " + i + ": " + x
	next i
	print "Final: " + x
	`

	err := mBasic.Run(code)
	if err != nil {
		panic(err)
	}
}
