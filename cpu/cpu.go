package cpu

import (
	"fmt"
	"gbemu/mmu"
	"gbemu/utils"
)

var logger *utils.Logger = utils.NewLogger(false)

func printByte(opcode byte) {
	fmt.Printf("%#02x\n", opcode)
}

func printWord(opcode uint16) {
	fmt.Printf("%#04x\n", opcode)
}

type CPU struct {
	mmu        *mmu.MMU
	ticks      uint8
	TotalTicks uint32

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
	fmt.Printf("AF: %#02x%02x\n", cpu.a, cpu.f)
	fmt.Printf("BC: %#02x%02x\n", cpu.b, cpu.c)
	fmt.Printf("DE: %#02x%02x\n", cpu.d, cpu.e)
	fmt.Printf("HL: %#02x%02x\n", cpu.h, cpu.l)
	fmt.Printf("SP: %#04x\n", cpu.sp)
	fmt.Printf("PC: %#04x\n", cpu.pc)
	fmt.Printf("TotalTicks: %08d\n", cpu.TotalTicks)
	fmt.Printf("lcdc: %#02x\n", cpu.mmu.Read(0xff40))
	fmt.Printf("stat: %#02x\n", cpu.mmu.Read(0xff41))
	fmt.Printf("ly: %#02x\n", cpu.mmu.Read(0xff44))
	fmt.Printf("lyc: %#02x\n", cpu.mmu.Read(0xff45))
	fmt.Printf("instruction: %#02x\n", cpu.mmu.Read(cpu.pc))
	fmt.Println("--------------------")
	// fmt.Println("HRAM")
	// fmt.Println("--------------------")
	// fmt.Printf("[0xff05] = %#02x\n", cpu.mmu.Read(0xff05))
	// fmt.Printf("[0xff06] = %#02x\n", cpu.mmu.Read(0xff06))
	// fmt.Printf("[0xff07] = %#02x\n", cpu.mmu.Read(0xff07))
	// fmt.Printf("[0xff10] = %#02x\n", cpu.mmu.Read(0xff10))
	// fmt.Printf("[0xff11] = %#02x\n", cpu.mmu.Read(0xff11))
	// fmt.Printf("[0xff12] = %#02x\n", cpu.mmu.Read(0xff12))
	// fmt.Printf("[0xff14] = %#02x\n", cpu.mmu.Read(0xff14))
	// fmt.Printf("[0xff16] = %#02x\n", cpu.mmu.Read(0xff16))
	// fmt.Printf("[0xff17] = %#02x\n", cpu.mmu.Read(0xff17))
	// fmt.Printf("[0xff19] = %#02x\n", cpu.mmu.Read(0xff19))
	// fmt.Printf("[0xff1a] = %#02x\n", cpu.mmu.Read(0xff1a))
	// fmt.Printf("[0xff1b] = %#02x\n", cpu.mmu.Read(0xff1b))
	// fmt.Printf("[0xff1c] = %#02x\n", cpu.mmu.Read(0xff1c))
	// fmt.Printf("[0xff1e] = %#02x\n", cpu.mmu.Read(0xff1e))
	// fmt.Printf("[0xff20] = %#02x\n", cpu.mmu.Read(0xff20))
	// fmt.Printf("[0xff21] = %#02x\n", cpu.mmu.Read(0xff21))
	// fmt.Printf("[0xff22] = %#02x\n", cpu.mmu.Read(0xff22))
	// fmt.Printf("[0xff23] = %#02x\n", cpu.mmu.Read(0xff23))
	// fmt.Printf("[0xff24] = %#02x\n", cpu.mmu.Read(0xff24))
	// fmt.Printf("[0xff25] = %#02x\n", cpu.mmu.Read(0xff25))
	// fmt.Printf("[0xff26] = %#02x\n", cpu.mmu.Read(0xff26))
	// fmt.Printf("[0xff40] = %#02x\n", cpu.mmu.Read(0xff40))
	// fmt.Printf("[0xff42] = %#02x\n", cpu.mmu.Read(0xff42))
	// fmt.Printf("[0xff43] = %#02x\n", cpu.mmu.Read(0xff43))
	// fmt.Printf("[0xff45] = %#02x\n", cpu.mmu.Read(0xff45))
	// fmt.Printf("[0xff47] = %#02x\n", cpu.mmu.Read(0xff47))
	// fmt.Printf("[0xff48] = %#02x\n", cpu.mmu.Read(0xff48))
	// fmt.Printf("[0xff49] = %#02x\n", cpu.mmu.Read(0xff49))
	// fmt.Printf("[0xff4a] = %#02x\n", cpu.mmu.Read(0xff4a))
	// fmt.Printf("[0xff4b] = %#02x\n", cpu.mmu.Read(0xff4b))
	// fmt.Printf("[0xffff] = %#02x\n", cpu.mmu.Read(0xffff))
}

func (cpu *CPU) PrintNextIns() {
	fmt.Printf("Next instruction: %#02x\n", cpu.mmu.Read(cpu.pc))
}

func (cpu *CPU) Fetch() uint8 {
	res := cpu.mmu.Read(cpu.pc)
	cpu.pc++

	// after reaching 0x100, disable BIOS
	if cpu.pc >= 0x100 {
		cpu.mmu.IsBooting = false
	}

	return res
}

func (cpu *CPU) FetchWord() uint16 {
	low := cpu.Fetch()

	high := cpu.Fetch()

	return uint16(high)<<8 | uint16(low)
}

func (cpu *CPU) HandleInterrupts() {
	if !cpu.isIntEnabled {
		return
	}

	req := cpu.mmu.Read(0xff0f)
	enabled := cpu.mmu.Read(0xffff)
	if req == 0 {
		return
	}

	// bit 0: V-Blank
	// bit 1: LCD
	// bit 2: Timer
	// bit 3: Serial
	// bit 4: Joypad
	for i := 0; i < 5; i++ {
		if req&(1<<i) > 0 && enabled&(1<<i) > 0 {
			cpu.serviceInterrupt(i)
		}
	}
}

func (cpu *CPU) serviceInterrupt(interrupt int) {
	cpu.isIntEnabled = false

	// reset interrupt
	req := cpu.mmu.Read(0xff0f)
	req &= ^(uint8(1 << interrupt))
	cpu.mmu.Write(0xff0f, req)

	// save current pc
	cpu.pushd16(cpu.pc)

	switch interrupt {
	case 0:
		cpu.pc = 0x40
	case 1:
		cpu.pc = 0x48
	case 2:
		cpu.pc = 0x50
	case 4:
		cpu.pc = 0x60
	}
}
