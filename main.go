package main

import (
	"bufio"
	"fmt"
	"gbemu/cpu"
	"gbemu/mmu"
	"os"
)

func debugMode(cpu *cpu.CPU, breakPoint *uint16) bool {
	fmt.Printf("_")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	switch input {
	case "d":
		cpu.Dump()
		return debugMode(cpu, breakPoint)
	case "n":
		cpu.Execute()
		*breakPoint = cpu.GetPC()
		return true
	case "c":
		// loop until BIOS end
		*breakPoint = 0x100
		return true
	case "q":
		// quit
		return false
	default:
		return false
	}
}

func main() {
	mmu := mmu.New()
	cpu := cpu.New(mmu)
	var breakPoint uint16 = 0x93

	for {
		fmt.Printf("%#04x : ", cpu.GetPC())

		if cpu.GetPC() == breakPoint {
			isContinue := debugMode(cpu, &breakPoint)
			if !isContinue {
				break
			}
		} else {
			cpu.Execute()
		}

	}
}
