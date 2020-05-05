package mmu

import (
	"fmt"
	"gbemu/gpu"
	"gbemu/joypad"
	"gbemu/timer"
)

type MMU struct {
	bios      [0x100]uint8
	cartridge []byte

	memory [0x10000]uint8

	IsBooting bool

	gpu    *gpu.GPU
	timer  *timer.Timer
	joypad *joypad.Joypad

	cartridgeType  uint8
	currentROMBank uint8
	currentRAMBank uint8
}

func New(gpu *gpu.GPU, timer *timer.Timer, joypad *joypad.Joypad) *MMU {
	mmu := &MMU{
		gpu:    gpu,
		timer:  timer,
		joypad: joypad,
	}

	mmu.IsBooting = true

	mmu.bios = [0x100]uint8{
		// BIOS
		0x31, 0xFE, 0xFF, 0xAF, 0x21, 0xFF, 0x9F, 0x32, 0xCB, 0x7C, 0x20, 0xFB, 0x21, 0x26, 0xFF, 0x0E,
		0x11, 0x3E, 0x80, 0x32, 0xE2, 0x0C, 0x3E, 0xF3, 0xE2, 0x32, 0x3E, 0x77, 0x77, 0x3E, 0xFC, 0xE0,
		0x47, 0x11, 0x04, 0x01, 0x21, 0x10, 0x80, 0x1A, 0xCD, 0x95, 0x00, 0xCD, 0x96, 0x00, 0x13, 0x7B,
		0xFE, 0x34, 0x20, 0xF3, 0x11, 0xD8, 0x00, 0x06, 0x08, 0x1A, 0x13, 0x22, 0x23, 0x05, 0x20, 0xF9,
		0x3E, 0x19, 0xEA, 0x10, 0x99, 0x21, 0x2F, 0x99, 0x0E, 0x0C, 0x3D, 0x28, 0x08, 0x32, 0x0D, 0x20,
		0xF9, 0x2E, 0x0F, 0x18, 0xF3, 0x67, 0x3E, 0x64, 0x57, 0xE0, 0x42, 0x3E, 0x91, 0xE0, 0x40, 0x04,
		0x1E, 0x02, 0x0E, 0x0C, 0xF0, 0x44, 0xFE, 0x90, 0x20, 0xFA, 0x0D, 0x20, 0xF7, 0x1D, 0x20, 0xF2,
		0x0E, 0x13, 0x24, 0x7C, 0x1E, 0x83, 0xFE, 0x62, 0x28, 0x06, 0x1E, 0xC1, 0xFE, 0x64, 0x20, 0x06,
		0x7B, 0xE2, 0x0C, 0x3E, 0x87, 0xF2, 0xF0, 0x42, 0x90, 0xE0, 0x42, 0x15, 0x20, 0xD2, 0x05, 0x20,
		0x4F, 0x16, 0x20, 0x18, 0xCB, 0x4F, 0x06, 0x04, 0xC5, 0xCB, 0x11, 0x17, 0xC1, 0xCB, 0x11, 0x17,
		0x05, 0x20, 0xF5, 0x22, 0x23, 0x22, 0x23, 0xC9, 0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B,
		0x03, 0x73, 0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D, 0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E,
		0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99, 0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC,
		0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E, 0x3c, 0x42, 0xB9, 0xA5, 0xB9, 0xA5, 0x42, 0x3C,
		0x21, 0x04, 0x01, 0x11, 0xA8, 0x00, 0x1A, 0x13, 0xBE, 0x20, 0xFE, 0x23, 0x7D, 0xFE, 0x34, 0x20,
		0xF5, 0x06, 0x19, 0x78, 0x86, 0x23, 0x05, 0x20, 0xFB, 0x86, 0x20, 0xFE, 0x3E, 0x01, 0xE0, 0x50,
	}

	mmu.memory = [0x10000]uint8{
		// BIOS
		0x31, 0xFE, 0xFF, 0xAF, 0x21, 0xFF, 0x9F, 0x32, 0xCB, 0x7C, 0x20, 0xFB, 0x21, 0x26, 0xFF, 0x0E,
		0x11, 0x3E, 0x80, 0x32, 0xE2, 0x0C, 0x3E, 0xF3, 0xE2, 0x32, 0x3E, 0x77, 0x77, 0x3E, 0xFC, 0xE0,
		0x47, 0x11, 0x04, 0x01, 0x21, 0x10, 0x80, 0x1A, 0xCD, 0x95, 0x00, 0xCD, 0x96, 0x00, 0x13, 0x7B,
		0xFE, 0x34, 0x20, 0xF3, 0x11, 0xD8, 0x00, 0x06, 0x08, 0x1A, 0x13, 0x22, 0x23, 0x05, 0x20, 0xF9,
		0x3E, 0x19, 0xEA, 0x10, 0x99, 0x21, 0x2F, 0x99, 0x0E, 0x0C, 0x3D, 0x28, 0x08, 0x32, 0x0D, 0x20,
		0xF9, 0x2E, 0x0F, 0x18, 0xF3, 0x67, 0x3E, 0x64, 0x57, 0xE0, 0x42, 0x3E, 0x91, 0xE0, 0x40, 0x04,
		0x1E, 0x02, 0x0E, 0x0C, 0xF0, 0x44, 0xFE, 0x90, 0x20, 0xFA, 0x0D, 0x20, 0xF7, 0x1D, 0x20, 0xF2,
		0x0E, 0x13, 0x24, 0x7C, 0x1E, 0x83, 0xFE, 0x62, 0x28, 0x06, 0x1E, 0xC1, 0xFE, 0x64, 0x20, 0x06,
		0x7B, 0xE2, 0x0C, 0x3E, 0x87, 0xF2, 0xF0, 0x42, 0x90, 0xE0, 0x42, 0x15, 0x20, 0xD2, 0x05, 0x20,
		0x4F, 0x16, 0x20, 0x18, 0xCB, 0x4F, 0x06, 0x04, 0xC5, 0xCB, 0x11, 0x17, 0xC1, 0xCB, 0x11, 0x17,
		0x05, 0x20, 0xF5, 0x22, 0x23, 0x22, 0x23, 0xC9, 0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B,
		0x03, 0x73, 0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D, 0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E,
		0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99, 0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC,
		0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E, 0x3c, 0x42, 0xB9, 0xA5, 0xB9, 0xA5, 0x42, 0x3C,
		0x21, 0x04, 0x01, 0x11, 0xA8, 0x00, 0x1A, 0x13, 0xBE, 0x20, 0xFE, 0x23, 0x7D, 0xFE, 0x34, 0x20,
		0xF5, 0x06, 0x19, 0x78, 0x86, 0x23, 0x05, 0x20, 0xFB, 0x86, 0x20, 0xFE, 0x3E, 0x01, 0xE0, 0x50,

		// Cartridge
		0x00, 0x00, 0x00, 0x00, // padding
		// logo start from 0x0104
		0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B, 0x03, 0x73, 0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D,
		0x00, 0x08, 0x11, 0x1F, 0x88, 0x89, 0x00, 0x0E, 0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD, 0xD9, 0x99,
		0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC, 0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E,
	}

	// TODO: ff44 means current scan line. update it dynamically
	mmu.memory[0xff44] = 0x90

	mmu.memory[0xff0f] = 0xe1

	return mmu
}

func (mmu *MMU) Load(buf []byte) {
	mmu.cartridge = buf

	mmu.cartridgeType = mmu.getCartridgeType()

	mmu.currentROMBank = 1
}

func (mmu *MMU) getCartridgeType() uint8 {
	// TODO: support more cartridge type
	switch mmu.cartridge[0x147] {

	case 0:
		return ROMONLY
	case 0x01, 0x02, 0x03:
		return MBC1
	}

	return 0
}

func (mmu *MMU) Read(addr uint16) uint8 {
	switch {
	case addr < 0x100:
		if mmu.IsBooting {
			return mmu.bios[addr]
		}
		return mmu.cartridge[addr]

	// Cartridge ROM, bank 0
	case 0x100 <= addr && addr <= 0x3fff:
		return mmu.cartridge[addr]

	// Cartridge ROM, other banks
	case 0x4000 <= addr && addr <= 0x7fff:
		return mmu.cartridge[addr+uint16(mmu.currentROMBank-1)*0x4000]

	// VRAM
	case 0x8000 <= addr && addr <= 0x9fff:
		return mmu.gpu.Read(addr)

	// joypad
	case addr == 0xff00:
		return mmu.joypad.Read()

	// Timer
	case 0xff04 <= addr && addr <= 0xff07:
		return mmu.timer.Read(addr)

	// LCD
	case 0xff40 <= addr && addr <= 0xff4b:
		return mmu.gpu.Read(addr)

	// OAM
	case 0xfe00 <= addr && addr <= 0xfe9f:
		return mmu.gpu.Read(addr)

	}

	return mmu.memory[addr]
}

func (mmu *MMU) Write(addr uint16, val uint8) {
	switch {
	case addr < 0x100:
		if mmu.IsBooting {
			mmu.bios[addr] = val
			return
		}
		mmu.cartridge[addr] = val
		return

	// VRAM
	case 0x8000 <= addr && addr <= 0x9fff:
		mmu.gpu.Write(addr, val)
		return

	case 0xa000 <= addr && addr <= 0xbfff:
		mmu.memory[addr] = val

	// joypad
	case addr == 0xff00:
		mmu.joypad.Write(val)
		return

	// Timer
	case 0xff04 <= addr && addr <= 0xff07:
		mmu.timer.Write(addr, val)
		return

	// LCD
	case 0xff40 <= addr && addr <= 0xff4b:
		if addr == 0xff46 {
			mmu.dmaTransfer(val)
			return
		}
		mmu.gpu.Write(addr, val)
		return

	// OAM
	case 0xfe00 <= addr && addr <= 0xfe9f:
		mmu.gpu.Write(addr, val)
		return

	case addr == 0xff0f:
		mmu.memory[addr] = val&0x1f | 0xe0
		return

	case addr == 0xff02 && val == 0x81:
		fmt.Printf("%c", mmu.memory[0xff01])
		return
	}

	mmu.memory[addr] = val
}

func (mmu *MMU) ReadWord(addr uint16) uint16 {
	return uint16(mmu.Read(addr)) | uint16(mmu.Read(addr+1))<<8
}

func (mmu *MMU) WriteWord(addr uint16, val uint16) {
	mmu.Write(addr, uint8(val&0xff))
	mmu.Write(addr+1, uint8((val>>8)&0xff))
}

// DMA transfer
// The written value specifies the transfer source address divided by 100h
// src: XX00-XX9f
// dst: fe00-fe9f
func (mmu *MMU) dmaTransfer(val uint8) {
	var src uint16 = uint16(val) << 8
	for i := 0; i < 0xa0; i++ {
		mmu.Write(0xfe00+uint16(i), mmu.Read(src+uint16(i)))
	}
}

func (mmu *MMU) UpdateIntFlag() {
	intFlag := mmu.Read(0xff0f)

	if mmu.gpu.ReqVBlankInt {
		intFlag |= 1
	}

	if mmu.gpu.ReqLCDInt {
		intFlag |= 1 << 1
	}

	if mmu.timer.ReqTimerInt {
		intFlag |= 1 << 2
	}

	if mmu.joypad.ReqJoypadInt {
		intFlag |= 1 << 4
		mmu.joypad.ReqJoypadInt = false
	}

	mmu.Write(0xff0f, intFlag)
}

func (mmu *MMU) Test() {
	mmu.gpu.Test()
}
