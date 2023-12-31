package vm

import (
	"github.com/Arch-4ng3l/Monkey/code"
	"github.com/Arch-4ng3l/Monkey/object"
)

type Frame struct {
	fn          *object.CompiledFunction
	ip          int
	BasePointer int
}

func NewFrame(fn *object.CompiledFunction, basePointer int) *Frame {
	return &Frame{
		fn:          fn,
		ip:          -1,
		BasePointer: basePointer,
	}
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
