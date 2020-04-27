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
	cpu.isIntEnabled = false

	return cpu
}

// Reset GB registers to initial state
func (cpu *CPU) Reset() {
	cpu.mmu.IsBooting = false

	cpu.pc = 0x100

	cpu.setReg16("AF", 0x01b0)
	cpu.setReg16("BC", 0x0013)
	cpu.setReg16("DE", 0x00d8)
	cpu.setReg16("HL", 0x014d)
	cpu.setReg16("SP", 0xfffe)
	cpu.mmu.Write(0xff05, 0x00)
	cpu.mmu.Write(0xff06, 0x00)
	cpu.mmu.Write(0xff07, 0x00)
	cpu.mmu.Write(0xff10, 0x80)
	cpu.mmu.Write(0xff11, 0xbf)
	cpu.mmu.Write(0xff12, 0xf3)
	cpu.mmu.Write(0xff14, 0xbf)
	cpu.mmu.Write(0xff16, 0x3f)
	cpu.mmu.Write(0xff17, 0x00)
	cpu.mmu.Write(0xff19, 0xbf)
	cpu.mmu.Write(0xff1a, 0x7f)
	cpu.mmu.Write(0xff1b, 0xff)
	cpu.mmu.Write(0xff1c, 0x9f)
	cpu.mmu.Write(0xff1e, 0xbf)
	cpu.mmu.Write(0xff20, 0xff)
	cpu.mmu.Write(0xff21, 0x00)
	cpu.mmu.Write(0xff22, 0x00)
	cpu.mmu.Write(0xff23, 0xbf)
	cpu.mmu.Write(0xff24, 0x77)
	cpu.mmu.Write(0xff25, 0xf3)
	cpu.mmu.Write(0xff26, 0xf1)
	cpu.mmu.Write(0xff40, 0x91)
	cpu.mmu.Write(0xff42, 0x00)
	cpu.mmu.Write(0xff43, 0x00)
	cpu.mmu.Write(0xff45, 0x00)
	cpu.mmu.Write(0xff47, 0xfc)
	cpu.mmu.Write(0xff48, 0xff)
	cpu.mmu.Write(0xff49, 0xff)
	cpu.mmu.Write(0xff4a, 0x00)
	cpu.mmu.Write(0xff4b, 0x00)
	cpu.mmu.Write(0xffff, 0x00)

	// for testing ROM
	// cpu.setReg16("AF", 0x1180)
	// cpu.setReg16("BC", 0x0000)
	// cpu.setReg16("DE", 0xff56)
	// cpu.setReg16("HL", 0x000d)
	// cpu.setReg16("SP", 0xfffe)
	// cpu.mmu.Test()

}

func (cpu *CPU) Dump() {
	fmt.Printf("\n")
	fmt.Println("--------------------")
	fmt.Printf("PC: %#04x\n", cpu.pc)
	fmt.Println("--------------------")
	fmt.Printf("AF: %#02x%02x\n", cpu.a, cpu.f)
	fmt.Printf("BC: %#02x%02x\n", cpu.b, cpu.c)
	fmt.Printf("DE: %#02x%02x\n", cpu.d, cpu.e)
	fmt.Printf("HL: %#02x%02x\n", cpu.h, cpu.l)
	fmt.Printf("SP: %#04x\n", cpu.sp)
	fmt.Printf("TotalTicks: %08d\n", cpu.TotalTicks)
	fmt.Printf("lcdc: %#02x\n", cpu.mmu.Read(0xff40))
	fmt.Printf("stat: %#02x\n", cpu.mmu.Read(0xff41))
	fmt.Printf("ly: %#02x\n", cpu.mmu.Read(0xff44))
	fmt.Printf("lyc: %#02x\n", cpu.mmu.Read(0xff45))
	fmt.Printf("instruction: %#02x\n", cpu.mmu.Read(cpu.pc))
	fmt.Printf("ie %#02x\n", cpu.mmu.Read(0xffff))
	fmt.Printf("if %#02x\n", cpu.mmu.Read(0xff0f))
	fmt.Printf("[0xff80] = %#02x\n", cpu.mmu.Read(0xff80))
	fmt.Printf("[0xff85] = %#02x\n", cpu.mmu.Read(0xff85))
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
	cpu.mmu.UpdateIntFlag()

	intFlag := cpu.mmu.Read(0xff0f)
	intEnabled := cpu.mmu.Read(0xffff)

	if !cpu.isIntEnabled {
		if cpu.halt && intFlag&intEnabled > 0 {
			cpu.halt = false
		}
		return
	}

	if intFlag == 0 {
		return
	}

	// bit 0: V-Blank
	// bit 1: LCD
	// bit 2: Timer
	// bit 3: Serial
	// bit 4: Joypad
	for i := 0; i < 5; i++ {
		if intFlag&(1<<i) > 0 && intEnabled&(1<<i) > 0 {
			cpu.serviceInterrupt(i)
		}
	}
}

func (cpu *CPU) serviceInterrupt(interrupt int) {
	cpu.isIntEnabled = false
	cpu.halt = false

	// reset interrupt flag
	intFlag := cpu.mmu.Read(0xff0f)
	intFlag &= ^(uint8(1 << interrupt))
	cpu.mmu.Write(0xff0f, intFlag)

	// save current pc
	cpu.pushd16(cpu.pc)

	// TODO: research. when interrupt occured, do I have to update ticks?
	// cpu.TotalTicks += 12

	switch interrupt {
	case 0:
		cpu.pc = 0x40
	case 1:
		cpu.pc = 0x48
	case 2:
		cpu.pc = 0x50
	case 3:
		cpu.pc = 0x58
	case 4:
		cpu.pc = 0x60
	}
}
