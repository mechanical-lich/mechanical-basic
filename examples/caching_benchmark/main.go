package main

import (
	"time"

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

	// First run
	startTime := time.Now()
	err := mBasic.Run(code)
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(startTime)
	println("First run took:", elapsed.String())

	// Second run
	startTime = time.Now()
	err = mBasic.Run(code)
	if err != nil {
		panic(err)
	}
	elapsed = time.Since(startTime)
	println("Second run took:", elapsed.String())

	// Third run
	startTime = time.Now()
	err = mBasic.Run(code)
	if err != nil {
		panic(err)
	}
	elapsed = time.Since(startTime)
	println("Second run took:", elapsed.String())
}
