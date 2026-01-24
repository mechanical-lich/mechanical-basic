package basic

import (
	"github.com/mechanical-lich/mechanical-basic/internal/basic"
)

type MechBasic struct {
	interpreter *basic.Interpreter
}

func NewMechanicalBasic() *MechBasic {

	return &MechBasic{
		interpreter: basic.NewInterpreter(),
	}
}

func (mb *MechBasic) RegisterFunc(name string, function basic.ExternalFunc) {
	mb.interpreter.RegisterFunction(name, function)
}

func (mb *MechBasic) Run(code string) error {
	return mb.interpreter.Interpret(code)
}
