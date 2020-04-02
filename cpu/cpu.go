package cpu

import "fmt"

func printByte(opcode byte) {
	fmt.Printf("%#02x\n", opcode)
}

func printWord(opcode uint16) {
	fmt.Printf("%#02x\n", opcode)
}

type CPU struct {
	clock int

	// registers
	a byte // accumulator
	f byte // flags
	b byte
	c byte
	d byte
	e byte
	h byte
	l byte

	pc uint16 // program counter
	sp uint16 // stack pointer
}

func (cpu *CPU) Demo() {
	cpu.setA(0x55)
	printWord(cpu.getHL())
	cpu.setBC(0x1234)
	printByte(cpu.getB())
	printByte(Z)
	cpu.LDr8d8('B', 10)
	printByte(cpu.getB())
}

func New() *CPU {
	cpu := &CPU{}

	return cpu
}

// LD nn, n
func (cpu *CPU) LDr8d8(reg byte, n byte) {
	fmt.Printf("LD %c, %d\n", reg, n)

	switch reg {
	case 'B':
		cpu.setB(n)
	case 'C':
		cpu.setC(n)
	case 'D':
		cpu.setD(n)
	case 'E':
		cpu.setE(n)
	case 'H':
		cpu.setH(n)
	case 'L':
		cpu.setL(n)
	}
}

// func (cpu *CPU) LDr8r8(reg1 byte, reg2 byte) {

// }
