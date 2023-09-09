package vm

import (
	"fmt"

	"github.com/Arch-4ng3l/Monkey/code"
	"github.com/Arch-4ng3l/Monkey/compiler"
	"github.com/Arch-4ng3l/Monkey/object"
)

const StackSize = 2048

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
		case code.OpAdd:
			right := vm.pop()
			left := vm.pop()
			leftVal := left.(*object.Integer).Value
			rightVal := right.(*object.Integer).Value
			result := leftVal + rightVal
			vm.push(&object.Integer{Value: result})
		}

	}
	return nil
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
