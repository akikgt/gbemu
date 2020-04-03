package main

import (
	"gbemu/cpu"
	"gbemu/mmu"
)

func main() {
	mmu := mmu.New()
	cpu := cpu.New(mmu)
	var breakPoint uint16 = 0x08

	for {
		if cpu.GetPC() >= breakPoint {
			break
		}
		cpu.Execute()
	}
}
