package cpu

// d8  means immediate 8 bit data
// d16 means immediate 16 bit data
// r8  means 8 bit register
// r16 means 16 bit register

func testBit(b uint8, val uint8) bool {
	if val>>b&1 == 1 {
		return true
	}
	return false
}

func (cpu *CPU) setFlags(z, n, h, c uint8) {
	var newFlag uint8 = 0
	oldFlag := cpu.getReg8("F")

	for b := Z; b >= C; b-- {
		var status uint8

		switch b {
		case Z:
			status = z
		case N:
			status = n
		case H:
			status = h
		case C:
			status = c
		}

		switch status {
		case RESET:
			newFlag &= ^(1 << b)
		case SET:
			newFlag |= 1 << b
		case NA:
			newFlag |= oldFlag & (1 << b)
		}
	}

	cpu.setReg8("F", newFlag)
}

// NOP do nothing
func (cpu *CPU) NOP() {
	logger.Log("NOP\n")
}

// LDr8d8 put value d8 into r8
func (cpu *CPU) LDr8d8(reg string) {
	n := cpu.Fetch()
	logger.Log("LD %s, %#02x\n", reg, n)

	cpu.setReg8(reg, n)
}

// LDr8r8 put value reg2 into reg1
func (cpu *CPU) LDr8r8(reg1, reg2 string) {
	logger.Log("LD %s, %s\n", reg1, reg2)

	var val byte
	if reg2 == "(HL)" {
		addr := cpu.getReg16("HL")
		val = cpu.mmu.Read(addr)
	} else {
		val = cpu.getReg8(reg2)
	}

	cpu.setReg8(reg1, val)
}

// LDr16d16 put value d16 into r16
func (cpu *CPU) LDr16d16(reg string) {
	low := cpu.Fetch()
	high := cpu.Fetch()

	nn := uint16(high)<<8 | uint16(low)
	logger.Log("LD %s, %#02x\n", reg, nn)

	cpu.setReg16(reg, nn)
}

// XORr8 exclusive OR n with register A, result in A
func (cpu *CPU) XORr8(reg string) {
	var n uint8
	if reg == "#" {
		n = cpu.Fetch()
		logger.Log("XOR %#02x\n", n)
	} else {
		n = cpu.getReg8(reg)
		logger.Log("XOR %s\n", reg)
	}

	val := cpu.getReg8("A") ^ n

	if val == 0 {
		cpu.setFlags(1, 0, 0, 0)
	}

	cpu.setReg8("A", val)
}

// LDDmHLA put A into memory address HL. Decrement HL
func (cpu *CPU) LDDmHLA() {
	logger.Log("LDD (HL), A\n")

	addr := cpu.getReg16("HL")
	cpu.mmu.Write(addr, cpu.getReg8("A"))

	addr--
	cpu.setReg16("HL", addr)
}

// JRccd8 if current condition is true, add n to current address and jump to it
// n = one byte signed immediate value
func (cpu *CPU) JRccd8(cc string) {
	var n int8 = int8(cpu.Fetch())

	switch cc {
	case "NZ":
		if !testBit(Z, cpu.getReg8("F")) {
			cpu.pc = uint16(int32(cpu.pc) + int32(n))
		}
	}
}

////////////////////////
// CB prefixed

// BITbr8 test bit b in register r8
func (cpu *CPU) BITbr8(b uint8, reg string) {
	logger.Log("BIT %d, %s\n", b, reg)

	var val byte
	if reg == "(HL)" {
		addr := cpu.getReg16("HL")
		val = cpu.mmu.Read(addr)
	} else {
		val = cpu.getReg8(reg)
	}

	if testBit(b, val) {
		cpu.setFlags(RESET, RESET, SET, NA)
	} else {
		cpu.setFlags(SET, RESET, SET, NA)
	}
}
