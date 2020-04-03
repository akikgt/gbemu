package main

import (
	"fmt"
	"gbemu/cpu"
	"gbemu/mmu"
)

func printOpcode(opcode byte) {
	fmt.Printf("%#02x\n", opcode)
}

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
