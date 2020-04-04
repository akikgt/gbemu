package main

import (
	"gbemu/cpu"
	"gbemu/mmu"
)

func main() {
	mmu := mmu.New()
	cpu := cpu.New(mmu)
	var breakPoint uint16 = 0x11

	for {
		if cpu.GetPC() >= breakPoint {
			cpu.Dump()
			break
		}
		cpu.Execute()
	}
}
