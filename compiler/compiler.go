package compiler

import (
	"fmt"

	"github.com/Arch-4ng3l/Monkey/ast"
	"github.com/Arch-4ng3l/Monkey/code"
	"github.com/Arch-4ng3l/Monkey/object"
)

type CompilationScope struct {
	instructions    code.Instructions
	lastInstruction EmittedInstruction
	prevInstruction EmittedInstruction
}

type Compiler struct {
	constants   []object.Object
	symbolTable *SymbolTable
	scopes      []CompilationScope
	scopeIdx    int
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

func New() *Compiler {
	mainScope := CompilationScope{
		instructions:    code.Instructions{},
		lastInstruction: EmittedInstruction{},
		prevInstruction: EmittedInstruction{},
	}
	return &Compiler{
		constants:   []object.Object{},
		symbolTable: NewSymbolTable(),
		scopes:      []CompilationScope{mainScope},
		scopeIdx:    0,
	}
}
func NewWithState(s *SymbolTable, constants []object.Object) *Compiler {
	comp := New()
	comp.symbolTable = s
	comp.constants = constants
	return comp
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.ExpresssionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(code.OpPop)

	case *ast.LetStatement:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}
		symbol := c.symbolTable.Define(node.Name.Value)
		c.emit(code.OpSetGlobal, symbol.Index)

	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}
		//Placeholder Value
		jmpNotTrue := c.emit(code.OpJmpNotTrue, 9999)

		err = c.Compile(node.If)
		if err != nil {
			return err
		}

		if c.lastInstructionIsPop() {
			c.removeLastPop()
		}

		jmp := c.emit(code.OpJmp, 9999)
		compiledPos := len(c.currentInstructions())
		c.changeOperand(jmpNotTrue, compiledPos)

		if node.Else == nil {
			c.emit(code.OpNull)

		} else {
			err := c.Compile(node.Else)
			if err != nil {
				return err
			}

			if c.lastInstructionIsPop() {
				c.removeLastPop()
			}
		}
		afterElsePos := len(c.currentInstructions())
		c.changeOperand(jmp, afterElsePos)

	case *ast.CallExpression:
		err := c.Compile(node.Function)
		if err != nil {
			return err
		}
		c.emit(code.OpCall)

	case *ast.ReturnStatement:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}

		c.emit(code.OpReturnValue)

	case *ast.FunctionLiteral:
		c.enterScope()
		err := c.Compile(node.Body)
		if err != nil {
			return err
		}
		if c.lastInstructionIsPop() {
			c.replaceLastPopWithReturn()
		}

		if !c.lastInstructionIs(code.OpReturnValue) {
			c.emit(code.OpReturn)
		}

		ins := c.leaveScope()
		compiledFn := &object.CompiledFunction{Instructions: ins}
		c.emit(code.OpConstant, c.addConstant(compiledFn))

	case *ast.BlockStatement:
		for i := range node.Statements {
			err := c.Compile(node.Statements[i])
			if err != nil {
				return err
			}
		}

	case *ast.InfixExpression:
		if node.Operator == "<" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}

			err = c.Compile(node.Left)
			if err != nil {
				return err
			}
			c.emit(code.OpGreaterThan)
			return nil
		}

		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case ">":
			c.emit(code.OpGreaterThan)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		}

	case *ast.IndexExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		err = c.Compile(node.Index)
		if err != nil {
			return err
		}
		c.emit(code.OpIndex)

	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("Unknown Operator %s", node.Operator)
		}

	case *ast.Ident:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}
		c.emit(code.OpGetGlobal, symbol.Index)

	case *ast.ArrayLiteral:
		for _, el := range node.Elements {
			err := c.Compile(el)
			if err != nil {
				return err
			}
		}
		c.emit(code.OpArray, len(node.Elements))

	case *ast.IntLiteral:
		integer := &object.Integer{Value: int(node.Value)}
		c.emit(code.OpConstant, c.addConstant(integer))
	case *ast.StrLiteral:
		str := &object.String{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(str))

	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	}

	return nil
}

func (c *Compiler) replaceLastPopWithReturn() {
	lastPos := c.scopes[c.scopeIdx].lastInstruction.Position
	c.replaceInstruction(lastPos, code.Make(code.OpReturnValue))

	c.scopes[c.scopeIdx].lastInstruction.Opcode = code.OpReturnValue
}

func (c *Compiler) enterScope() {
	scope := CompilationScope{
		instructions:    code.Instructions{},
		lastInstruction: EmittedInstruction{},
		prevInstruction: EmittedInstruction{},
	}
	c.scopes = append(c.scopes, scope)
	c.scopeIdx++
}
func (c *Compiler) leaveScope() code.Instructions {
	ins := c.currentInstructions()

	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIdx--
	return ins
}

func (c *Compiler) currentInstructions() code.Instructions {
	return c.scopes[c.scopeIdx].instructions
}

func (c *Compiler) changeOperand(opPos, operand int) {
	op := code.Opcode(c.currentInstructions()[opPos])
	newInstruction := code.Make(op, operand)
	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	ins := c.currentInstructions()

	for i := 0; i < len(newInstruction); i++ {
		ins[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) removeLastPop() {
	last := c.scopes[c.scopeIdx].lastInstruction
	prev := c.scopes[c.scopeIdx].prevInstruction
	oldIns := c.currentInstructions()
	newIns := oldIns[:last.Position]

	c.scopes[c.scopeIdx].instructions = newIns
	c.scopes[c.scopeIdx].lastInstruction = prev
}

func (c *Compiler) lastInstructionIsPop() bool {
	return c.scopes[c.scopeIdx].lastInstruction.Opcode == code.OpPop
}
func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
	if len(c.currentInstructions()) == 0 {
		return false
	}
	return c.scopes[c.scopeIdx].lastInstruction.Opcode == op
}

func (c *Compiler) emit(op code.Opcode, operangs ...int) int {
	ins := code.Make(op, operangs...)
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)
	return pos
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.scopes[c.scopeIdx].lastInstruction

	last := EmittedInstruction{Opcode: op, Position: pos}

	c.scopes[c.scopeIdx].prevInstruction = previous
	c.scopes[c.scopeIdx].lastInstruction = last
}

func (c *Compiler) addInstruction(ins []byte) int {
	pos := len(c.currentInstructions())

	updatedIns := append(c.currentInstructions(), ins...)

	c.scopes[c.scopeIdx].instructions = updatedIns

	return pos
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}

}
