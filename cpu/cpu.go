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
	cpu.LDr8d8("B", 10)
	printByte(cpu.getReg8("B"))
	cpu.setReg16("AF", 0x1234)
	printByte(cpu.getReg8("A"))
	printByte(cpu.getReg8("F"))
	cpu.LDr8r8("A", "(HL)")
	printWord(cpu.getReg16("AF"))
}

// New return CPU
func New() *CPU {
	cpu := &CPU{}

	return cpu
}

// LDr8d8 put value d8 into r8
func (cpu *CPU) LDr8d8(reg string, n byte) {
	fmt.Printf("LD %s, %d\n", reg, n)

	cpu.setReg8(reg, n)
}

// LDr8r8 put value reg2 into reg1
func (cpu *CPU) LDr8r8(reg1 string, reg2 string) {
	fmt.Printf("LD %s, %s\n", reg1, reg2)

	var val byte
	if reg2 == "(HL)" {
		// TODO read from memory
		val = 0x00
		addr := cpu.getReg16("HL")
		fmt.Printf("read byte from address: %x", addr)
		// val = cpu.mmu.Read(addr)
	} else {
		val = cpu.getReg8(reg2)
	}

	cpu.setReg8(reg1, val)
}
