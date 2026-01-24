package main

import (
	"fmt"

	"github.com/mechanical-lich/mechanical-basic/pkg/basic"
)

func main() {
	mBasic := basic.NewMechanicalBasic()

	code := `
	print "Hello from Mechanical BASIC!"
	`
	// Without custom print function, output goes to standard output
	fmt.Println("BEFORE:")
	mBasic.Run(code)

	// With it goes to the provided function
	logMessage := func(msg any) {
		fmt.Println("[LOG:]", msg)
	}

	// RegisterFunc(name, arguementCount, function pointer)
	mBasic.SetPrintFunc(logMessage)

	fmt.Println("AFTER:")
	mBasic.Run(code)

}
