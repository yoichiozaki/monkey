package vm

import (
	"monkey/code"
	"monkey/object"
)

const MaxFrame = 1024

// Frame is a data structure that holds execution-relevant information, like the instructions and the instruction pointer.
// In compiler or interpreter literature, this data structure is also called activation record.
type Frame struct {
	fn          *object.CompiledFunction // points to the compiled function referenced by the frame.
	ip          int                      // is the instruction pointer in this frame for this function.
	basePointer int                      // is the pointer value before execution of a function .
}

func NewFrame(fn *object.CompiledFunction, basePointer int) *Frame {
	f := &Frame{fn: fn, ip: -1, basePointer: basePointer}
	return f
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
