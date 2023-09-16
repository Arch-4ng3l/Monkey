package vm

import (
	"fmt"

	"github.com/Arch-4ng3l/Monkey/code"
	"github.com/Arch-4ng3l/Monkey/compiler"
	"github.com/Arch-4ng3l/Monkey/object"
)

const StackSize = 2048
const GlobalSize = 65536

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}

type Vm struct {
	constans     []object.Object
	instructions code.Instructions

	stack        []object.Object
	stackPointer int
	globals      []object.Object
}

func New(bytecode *compiler.Bytecode) *Vm {
	return &Vm{
		instructions: bytecode.Instructions,
		constans:     bytecode.Constants,
		stack:        make([]object.Object, StackSize),
		stackPointer: 0,
		globals:      make([]object.Object, GlobalSize),
	}
}
func NewWithGLobalStore(bytecode *compiler.Bytecode, s []object.Object) *Vm {
	vm := New(bytecode)
	vm.globals = s
	return vm
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
		case code.OpJmp:
			pos := int(code.ReadUint16(vm.instructions[i+1:]))
			i = pos - 1
		case code.OpJmpNotTrue:
			pos := int(code.ReadUint16(vm.instructions[i+1:]))
			i += 2
			condition := vm.pop()
			if !isTrue(condition) {
				i = pos - 1
			}
		case code.OpIndex:
			index := vm.pop()
			left := vm.pop()
			err := vm.executeIndex(left, index)
			if err != nil {
				return err
			}
		case code.OpArray:
			numElements := int(code.ReadUint16(vm.instructions[i+1:]))
			i += 2
			arr := vm.buildArr(vm.stackPointer-numElements, vm.stackPointer)
			vm.stackPointer = vm.stackPointer - numElements

			err := vm.push(arr)
			if err != nil {
				return err
			}

		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[i+1:])
			i += 2

			err := vm.push(vm.constans[constIndex])
			if err != nil {
				return err
			}
		case code.OpSetGlobal:
			globalIdx := code.ReadUint16(vm.instructions[i+1:])
			i += 2

			vm.globals[globalIdx] = vm.pop()
		case code.OpGetGlobal:
			globalIdx := code.ReadUint16(vm.instructions[i+1:])
			i += 2
			err := vm.push(vm.globals[globalIdx])

			if err != nil {
				return err
			}
		case code.OpAdd, code.OpDiv, code.OpMul, code.OpSub:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			err := vm.executeComparision(op)
			if err != nil {
				return err
			}
		case code.OpBang:
			err := vm.executeBangOperator()
			if err != nil {
				return err
			}
		case code.OpMinus:
			err := vm.executeMinusOperator()
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
		case code.OpNull:
			err := vm.push(Null)
			if err != nil {
				return err
			}

		case code.OpPop:
			vm.pop()
		}

	}
	return nil
}

func (vm *Vm) executeIndex(left, index object.Object) error {
	if index.Type() != object.INTEGER_OBJ {
		return fmt.Errorf("Cant use %s As indext", index.Type())
	}
	if left.Type() == object.ARR_OBJ {
		return vm.executeArrIdx(left, index)
	}
	if left.Type() == object.STR_OBJ {
		return vm.executeStrIdx(left, index)
	}

	return fmt.Errorf("Cant index %s with %s", left.Type(), index.Type())
}

func (vm *Vm) executeStrIdx(left, index object.Object) error {
	str := left.(*object.String)
	i := index.(*object.Integer).Value
	max := len(str.Value) - 1
	if i < 0 || i > max {
		return vm.push(Null)
	}
	return vm.push(&object.String{Value: string(str.Value[i])})
}

func (vm *Vm) executeArrIdx(left, index object.Object) error {
	arr := left.(*object.Array)
	i := index.(*object.Integer).Value
	max := len(arr.Elements) - 1
	if i < 0 || i > max {
		return vm.push(Null)
	}
	return vm.push(arr.Elements[i])
}

func (vm *Vm) buildArr(startIdx, endIdx int) object.Object {
	elements := make([]object.Object, endIdx-startIdx)
	for i := startIdx; i < endIdx; i++ {
		elements[i-startIdx] = vm.stack[i]
	}
	return &object.Array{Elements: elements}
}

func (vm *Vm) executeMinusOperator() error {
	operand := vm.pop()
	if operand.Type() != object.INTEGER_OBJ {
		return fmt.Errorf("Operand Is not An Integer Object: %s", operand.Type())
	}

	val := operand.(*object.Integer).Value

	return vm.push(&object.Integer{Value: -val})
}

func (vm *Vm) executeBangOperator() error {
	operand := vm.pop()
	switch operand {
	case True:
		return vm.push(False)
	case False:
		return vm.push(True)
	case Null:
		return vm.push(True)
	default:
		return vm.push(False)
	}

}
func (vm *Vm) executeComparision(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return vm.executeIntegerComparision(op, left, right)
	}
	switch op {
	case code.OpEqual:
		vm.push(vm.boolToBoolObject(left == right))
	case code.OpNotEqual:
		vm.push(vm.boolToBoolObject(left != right))
	}

	return nil
}
func (vm *Vm) executeIntegerComparision(op code.Opcode, left, right object.Object) error {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch op {
	case code.OpGreaterThan:
		vm.push(vm.boolToBoolObject(leftVal > rightVal))
	case code.OpEqual:
		vm.push(vm.boolToBoolObject(leftVal == rightVal))
	case code.OpNotEqual:
		vm.push(vm.boolToBoolObject(leftVal != rightVal))
	}

	return nil
}
func (vm *Vm) boolToBoolObject(input bool) *object.Boolean {
	if input {
		return True
	}
	return False
}

func (vm *Vm) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()
	leftType := left.Type()
	rightType := right.Type()

	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return vm.executeBinaryIntOperation(op, left, right)
	}
	if leftType == object.STR_OBJ && rightType == object.STR_OBJ {
		return vm.executeBinaryStrOperation(op, left, right)
	}
	return nil
}
func (vm *Vm) executeBinaryStrOperation(op code.Opcode, left, right object.Object) error {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	var res string
	switch op {
	case code.OpAdd:
		res = leftVal + rightVal
		return vm.push(&object.String{Value: res})
	default:
		return fmt.Errorf("unknown Operator for Strings: %d", op)
	}

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

func isTrue(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
}
