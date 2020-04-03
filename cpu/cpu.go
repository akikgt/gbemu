package cpu

import (
	"fmt"
	"gbemu/mmu"
)

func printByte(opcode byte) {
	fmt.Printf("%#02x\n", opcode)
}

func printWord(opcode uint16) {
	fmt.Printf("%#02x\n", opcode)
}

type CPU struct {
	mmu   *mmu.MMU
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
	// cpu.LDr8d8("B", 10)
	// printByte(cpu.getReg8("B"))
	// cpu.setReg16("AF", 0x1234)
	// printByte(cpu.getReg8("A"))
	// printByte(cpu.getReg8("F"))
	// cpu.LDr8r8("A", "(HL)")
	// printWord(cpu.getReg16("AF"))
	// printByte(cpu.mmu.Read(0))
}

// New return CPU
func New(mmu *mmu.MMU) *CPU {
	cpu := &CPU{mmu: mmu}
	return cpu
}

func (cpu *CPU) Dump() {
	fmt.Println("--------------------")
	fmt.Printf("A: %#02x F: %#02x\n", cpu.a, cpu.f)
	fmt.Printf("B: %#02x C: %#02x\n", cpu.b, cpu.c)
	fmt.Printf("D: %#02x E: %#02x\n", cpu.d, cpu.e)
	fmt.Printf("H: %#02x L: %#02x\n", cpu.h, cpu.l)
	fmt.Printf("SP: %#02x\n", cpu.sp)
	fmt.Printf("PC: %#02x\n", cpu.pc)
	fmt.Println("--------------------")
}

func (cpu *CPU) Fetch() uint8 {
	res := cpu.mmu.Read(cpu.pc)
	cpu.pc++

	return res
}

func (cpu *CPU) Execute() {
	opcode := cpu.Fetch()

	switch opcode {
	case 0x00:
		cpu.NOP()

	// 8-bit loads
	case 0x06:
		cpu.LDr8d8("B")

	// 16-bit loads
	case 0x01:
		cpu.LDr16d16("BC")
	case 0x11:
		cpu.LDr16d16("DE")
	case 0x21:
		cpu.LDr16d16("HL")
	case 0x31:
		cpu.LDr16d16("SP")

	// 8-bit ALU
	case 0xaf:
		cpu.XORn("A")
	default:
		fmt.Printf("unknown opcode: %#02x\n", opcode)
	}

	cpu.Dump() // for debug
}
