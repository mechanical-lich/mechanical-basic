package main

import (
	"fmt"

	"github.com/mechanical-lich/mechanical-basic/pkg/basic"
	"github.com/mechanical-lich/mechanical-basic/pkg/functions"
)

func main() {
	mBasic := basic.NewMechanicalBasic()

	x := 0
	y := 0
	setX := func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, nil
		}

		newX, err := functions.EnsureInt(args[0])
		if err != nil {
			return nil, err
		}

		x = newX
		return nil, nil
	}

	setY := func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, nil
		}

		newY, err := functions.EnsureInt(args[0])
		if err != nil {
			return nil, err
		}

		y = newY
		return nil, nil
	}

	mBasic.RegisterFunc("setX", setX)
	mBasic.RegisterFunc("setY", setY)

	code := `
	function init():
		print "Init called"
	endfunction

	function hit(x, y, forceX, forceY):
		setX(x + forceX)
		setY(y + forceY)
	endfunction

	function update():
		print "Update called"
	endfunction
	`

	// Load the code
	err := mBasic.Load(code)
	if err != nil {
		panic(err)
	}

	fmt.Println("X and Y before calls:", x, y)

	// Call functions
	fmt.Println("Calling init....")
	_, err = mBasic.Call("init")
	if err != nil {
		panic(err)
	}

	fmt.Println("Calling hit....")
	_, err = mBasic.Call("hit", 10, 20, 5, -3)
	if err != nil {
		panic(err)
	}

	fmt.Println("Calling update....")
	_, err = mBasic.Call("update")
	if err != nil {
		panic(err)
	}

	fmt.Println("\nX and Y after calls:", x, y)
}
