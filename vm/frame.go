package vm

import (
	"monkey/code"
	"monkey/object"
)

const MaxFrame = 1024

// Frame is a data structure that holds execution-relevant information, like the instructions and the instruction pointer.
// In compiler or interpreter literature, this data structure is also called activation record.
type Frame struct {
	cl          *object.Closure // points to the closure referenced by the frame.
	ip          int             // is the instruction pointer in this frame for this function.
	basePointer int             // is the pointer value before execution of a function .
}

func NewFrame(cl *object.Closure, basePointer int) *Frame {
	f := &Frame{cl: cl, ip: -1, basePointer: basePointer}
	return f
}

func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}
