package main

import (
	"fmt"
	"gbemu/cpu"
	"gbemu/mmu"
)

func main() {
	mmu := mmu.New()
	cpu := cpu.New(mmu)
	var breakPoint uint16 = 0xa8

	for {
		fmt.Printf("%#04x : ", cpu.GetPC())

		if cpu.GetPC() >= breakPoint {
			cpu.Dump()
			break
		}
		cpu.Execute()
	}
}
