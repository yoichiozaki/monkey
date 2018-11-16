package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte

type Opcode byte

const (
	OpConstant      Opcode = iota // sets constant value in constant pool.
	OpAdd                         // pops 2 topmost elements from off the stack and adds them, pushes back on the top of the stack.
	OpSub                         // pops 2 topmost elements from off the stack and subtracts them, pushes back on the top of the stack.
	OpMul                         // pops 2 topmost elements from off the stack and multiplies them, pushes back on the top of the stack.
	OpDiv                         // pops 2 topmost elements from off the stack and divides them, pushes back on the top of the stack.
	OpPop                         // makes the stack clean after every expression statement.
	OpTrue                        // pushes an *object.Boolean(true) on to the stack.
	OpFalse                       // pushed an *object.Boolean(false) on to the stack.
	OpEqual                       // pops 2 topmost elements from off the stack and compares them, pushes back the result on the top of the stack.
	OpNotEqual                    // pops 2 topmost elements from off the stack and compares them, pushes back the result on the top of the stack.
	OpGreaterThan                 // pops 2 topmost elements from off the stack and compares them, pushes back the result on the top of the stack.
	OpMinus                       // pops 1 topmost element from off the stack and negates it, pushes back the result on the top of the stack.
	OpBang                        // pops 1 topmost element from off the stack and negates it, pushes back the result on the top of the stack.
	OpJumpNotTruthy               // jumps to a certain address if the topmost element on the stack is not truthy
	OpJump                        // jumps whatever the topmost element of the stack is.
	OpNull                        // pushes an *object.Null on to the stack.
	OpGetGlobal                   // gets global variable bound to its operand.
	OpSetGlobal                   // sets global variable bound to its operand.
	OpGetLocal                    // gets global variable bound to its operand.
	OpSetLocal                    // sets local variable bound to its operand.
	OpArray                       // tells how many elements the array has.
	OpHash                        // tells how many keys and values the hash has.
	OpIndex                       // pops 2 topmost elements off from the stack and performs the index operation, puts the result back on.
	OpCall                        // calls function.
	OpReturnValue                 // returns from function with return value. The returned value sits on top of the stack.
	OpReturn                      // return from function with no explicit return value, but implicit vm.Null.
)

type Definition struct {
	Name         string
	OperandWidth []int
}

var definitions = map[Opcode]*Definition{
	OpConstant:      {"OpConstant", []int{2}},
	OpAdd:           {"OpAdd", []int{}},
	OpSub:           {"OpSub", []int{}},
	OpMul:           {"OpMul", []int{}},
	OpDiv:           {"OpDiv", []int{}},
	OpPop:           {"OpPop", []int{}},
	OpTrue:          {"OpTrue", []int{}},
	OpFalse:         {"OpFalse", []int{}},
	OpEqual:         {"OpEqual", []int{}},
	OpNotEqual:      {"OpNotEqual", []int{}},
	OpGreaterThan:   {"OpGreaterThan", []int{}},
	OpMinus:         {"OpMinus", []int{}},
	OpBang:          {"OpBang", []int{}},
	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}},
	OpJump:          {"OpJump", []int{2}},
	OpNull:          {"OpNull", []int{}},
	OpGetGlobal:     {"OpGetGlobal", []int{2}},
	OpSetGlobal:     {"OpSetGlobal", []int{2}},
	OpGetLocal:      {"OpGetLocal", []int{1}},
	OpSetLocal:      {"OpSetLocal", []int{1}},
	OpArray:         {"OpArray", []int{2}},
	OpHash:          {"OpHash", []int{2}},
	OpIndex:         {"OpIndex", []int{}},
	OpCall:          {"OpCall", []int{1}},
	OpReturnValue:   {"OpReturnValue", []int{}},
	OpReturn:        {"OpReturn", []int{}},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d is undefined.", op)
	}
	return def, nil
}

func Make(op Opcode, operands ...int) []byte { // note: constant value is indexing with its order in the constant pool.
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}
	instructionLen := 1
	for _, w := range def.OperandWidth {
		instructionLen += w
	}
	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)
	offset := 1
	for i, o := range operands {
		width := def.OperandWidth[i]
		switch width {
		case 1:
			instruction[offset] = byte(o)
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}
	return instruction
}

func (ins Instructions) String() string {
	var out bytes.Buffer
	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s.\n", err)
			continue
		}
		operands, read := ReadOperands(def, ins[i+1:])
		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))
		i += 1 + read
	}
	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidth)
	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d.\n", len(operands), operandCount)
	}
	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}
	return fmt.Sprintf("ERROR: unhandled operandCount for %s is there.\n", def.Name)
}
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidth))
	offset := 0
	for i, width := range def.OperandWidth {
		switch width {
		case 1:
			operands[i] = int(ReadUint8(ins[offset:]))
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}
		offset += width
	}
	return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

func ReadUint8(ins Instructions) uint8 {
	return uint8(ins[0])
}
