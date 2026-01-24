package main

import (
	"fmt"

	"github.com/mechanical-lich/mechanical-basic/pkg/basic"
)

func main() {
	mBasic := basic.NewMechanicalBasic()

	code := `
	let counter = 0
	function count():
		counter = counter + 1
		print "Count is now: " + counter
	endfunction

	`

	// Load the code
	err := mBasic.Load(code)
	if err != nil {
		panic(err)
	}

	// Call functions
	for i := 0; i < 3; i++ {
		fmt.Println("Calling count function, iteration:", i)
		_, err = mBasic.Call("count")
		if err != nil {
			panic(err)
		}
	}

}
