package vm

import (
	"fmt"

	"github.com/Arch-4ng3l/Monkey/code"
	"github.com/Arch-4ng3l/Monkey/compiler"
	"github.com/Arch-4ng3l/Monkey/object"
)

const StackSize = 2048

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}

type Vm struct {
	constans     []object.Object
	instructions code.Instructions

	stack        []object.Object
	stackPointer int
}

func New(bytecode *compiler.Bytecode) *Vm {
	return &Vm{
		instructions: bytecode.Instructions,
		constans:     bytecode.Constants,
		stack:        make([]object.Object, StackSize),
		stackPointer: 0,
	}
}

func (vm *Vm) StackTop() object.Object {
	if vm.stackPointer == 0 {
		return nil
	}
	return vm.stack[vm.stackPointer-1]
}

func (vm *Vm) Run() error {
	for i := 0; i < len(vm.instructions); i++ {
		op := code.Opcode(vm.instructions[i])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[i+1:])
			i += 2

			err := vm.push(vm.constans[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd, code.OpDiv, code.OpMul, code.OpSub:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		case code.OpTrue:
			err := vm.push(True)
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := vm.push(False)
			if err != nil {
				return err
			}

		case code.OpPop:
			vm.pop()
		}

	}
	return nil
}

func (vm *Vm) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()
	leftType := left.Type()
	rightType := right.Type()

	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return vm.executeBinaryIntOperation(op, left, right)
	}
	return nil
}
func (vm *Vm) executeBinaryIntOperation(op code.Opcode, left, right object.Object) error {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	var res int
	switch op {
	case code.OpAdd:
		res = leftVal + rightVal
	case code.OpDiv:
		res = leftVal / rightVal
	case code.OpMul:
		res = leftVal * rightVal
	case code.OpSub:
		res = leftVal - rightVal
	default:
		return fmt.Errorf("unknown Operator for integers: %d", op)
	}

	return vm.push(&object.Integer{Value: res})
}

func (vm *Vm) push(obj object.Object) error {
	if vm.stackPointer >= StackSize {
		return fmt.Errorf("Stack Overflow")
	}
	vm.stack[vm.stackPointer] = obj
	vm.stackPointer++
	return nil
}

func (vm *Vm) pop() object.Object {
	obj := vm.stack[vm.stackPointer-1]
	vm.stackPointer--
	return obj
}

func (vm *Vm) LastPoppedStackElement() object.Object {
	return vm.stack[vm.stackPointer]
}
