package cpu

import (
	"fmt"
	"gbemu/mmu"
	"gbemu/utils"
)

var logger *utils.Logger = utils.NewLogger(true)

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

// New return CPU
func New(mmu *mmu.MMU) *CPU {
	cpu := &CPU{mmu: mmu}
	return cpu
}

func (cpu *CPU) Dump() {
	fmt.Printf("\n")
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
