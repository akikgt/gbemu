package main

import (
	"gbemu/cpu"
	"gbemu/mmu"
)

func main() {
	mmu := mmu.New()
	cpu := cpu.New(mmu)
	var breakPoint uint16 = 0x0b

	for {
		if cpu.GetPC() >= breakPoint {
			break
		}
		cpu.Execute()
	}
}
