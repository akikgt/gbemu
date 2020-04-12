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
	fmt.Printf("%#04x\n", opcode)
}

type CPU struct {
	mmu        *mmu.MMU
	ticks      uint8
	totalTicks uint32

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

	halt         bool
	stop         bool
	isIntEnabled bool
}

// New return CPU
func New(mmu *mmu.MMU) *CPU {
	cpu := &CPU{mmu: mmu}

	cpu.halt = false
	cpu.stop = false
	cpu.isIntEnabled = true

	return cpu
}

func (cpu *CPU) Dump() {
	fmt.Printf("\n")
	fmt.Println("--------------------")
	fmt.Printf("A: %#02x F: %#02x\n", cpu.a, cpu.f)
	fmt.Printf("B: %#02x C: %#02x\n", cpu.b, cpu.c)
	fmt.Printf("D: %#02x E: %#02x\n", cpu.d, cpu.e)
	fmt.Printf("H: %#02x L: %#02x\n", cpu.h, cpu.l)
	fmt.Printf("SP: %#04x\n", cpu.sp)
	fmt.Printf("PC: %#04x\n", cpu.pc)
	fmt.Printf("TotalTicks: %08d\n", cpu.totalTicks)
	fmt.Println("--------------------")
}

func (cpu *CPU) PrintNextIns() {
	fmt.Printf("Next instruction: %#02x\n", cpu.mmu.Read(cpu.pc))
}

func (cpu *CPU) Fetch() uint8 {
	res := cpu.mmu.Read(cpu.pc)
	cpu.pc++

	return res
}

func (cpu *CPU) FetchWord() uint16 {
	low := cpu.Fetch()

	high := cpu.Fetch()

	return uint16(high)<<8 | uint16(low)
}
