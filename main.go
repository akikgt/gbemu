package main

import (
	"fmt"
	"gbemu/cpu"
	"os"
)

func printOpcode(opcode byte) {
	fmt.Printf("%#02x\n", opcode)
}

func main() {
	f, err := os.Open("sgb_bios.bin")
	if err != nil {
		fmt.Println("Cannot open a file")
		os.Exit(1)
	}
	defer f.Close()

	// bios := make([]byte, 256)

	// f.Read(bios)
	// var pc uint16 = 0
	// for {
	// 	if pc > 30 {
	// 		break
	// 	}
	// 	b := bios[pc]
	// 	printOpcode(b)

	// 	pc++
	// }

	cpu := cpu.New()
	cpu.Demo()
}
