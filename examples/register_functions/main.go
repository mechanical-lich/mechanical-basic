package main

import (
	"errors"
	"fmt"
	"math"

	"github.com/mechanical-lich/mechanical-basic/pkg/basic"
	"github.com/mechanical-lich/mechanical-basic/pkg/functions"
)

func main() {
	mBasic := basic.NewMechanicalBasic()

	x := 124
	getX := func(args ...interface{}) (interface{}, error) {
		if len(args) > 0 {
			return nil, nil
		}
		return x, nil
	}

	setX := func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, errors.New("incorrect arguements")
		}

		newX, err := functions.EnsureInt(args[0])
		if err != nil {
			return nil, err
		}

		x = newX
		return nil, nil
	}

	pow := func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 {
			return nil, nil
		}

		// Check if input is float or int

		a, err := functions.EnsureFloat(args[0])
		if err != nil {
			fmt.Println("Error converting first argument:", err)
			return nil, err
		}

		b, err := functions.EnsureFloat(args[1])
		if err != nil {
			fmt.Println("Error converting second argument:", err)
			return nil, err
		}

		return math.Pow(a, b), nil
	}

	// RegisterFunc(name, arguementCount, function pointer)
	mBasic.RegisterFunc("getX", getX)
	mBasic.RegisterFunc("setX", setX)
	mBasic.RegisterFunc("pow", pow)

	var code string
	// Impacting state outside of the script example
	code = `
	print getX()
	setX(42)
	print getX()
	`
	mBasic.Run(code)

	// Multiple arguement example
	code = `let result = pow(2, 3)
	print "2 to the 3rd power is " + result
	`
	mBasic.Run(code)
}
