package cpu

// d8  means immediate 8 bit data
// d16 means immediate 16 bit data
// s8  means signed immediate 8 bit data, which are added to pc
// r8  means 8 bit register
// r16 means 16 bit register
// m*  means data in memory address *

func testBit(b uint8, val uint8) bool {
	if val>>b&1 == 1 {
		return true
	}
	return false
}

// NOP do nothing
func (cpu *CPU) NOP() {
	logger.Log("NOP\n")
}

//======================================================================
// Load
//======================================================================

// LDr8d8 put value d8 into r8
func (cpu *CPU) LDr8d8(reg string) {
	n := cpu.Fetch()
	logger.Log("LD %s, %#02x\n", reg, n)

	cpu.setReg8(reg, n)
}

// LDmHLd8 put value d8 into address HL
func (cpu *CPU) LDmHLd8() {
	n := cpu.Fetch()

	addr := cpu.getReg16("HL")

	cpu.mmu.Write(addr, n)

	logger.Log("LD (HL), %#02x\n", n)
}

// LDr8r8 put value reg2 into reg1
func (cpu *CPU) LDr8r8(reg1, reg2 string) {
	logger.Log("LD %s, %s\n", reg1, reg2)

	val := cpu.getReg8(reg2)

	cpu.setReg8(reg1, val)
}

// LDr8mr16 put value at address r16 into r8
func (cpu *CPU) LDr8mr16(reg1, reg2 string) {
	logger.Log("LD %s, (%s)", reg1, reg2)

	addr := cpu.getReg16(reg2)

	val := cpu.mmu.Read(addr)

	cpu.setReg8(reg1, val)
}

// LDmr16r8 put value into address r16
func (cpu *CPU) LDmr16r8(reg1, reg2 string) {
	logger.Log("LD (%s), %s", reg1, reg2)

	addr := cpu.getReg16(reg1)

	val := cpu.getReg8(reg2)

	cpu.mmu.Write(addr, val)
}

// LDmd16A put value A into address d16
func (cpu *CPU) LDmd16A() {
	low := cpu.Fetch()
	high := cpu.Fetch()

	addr := uint16(high)<<8 | uint16(low)

	cpu.mmu.Write(addr, cpu.getReg8("A"))

	logger.Log("LD (%#02x), A\n", addr)
}

// LDAmd16 put value at address d16 into A
func (cpu *CPU) LDAmd16() {
	low := cpu.Fetch()
	high := cpu.Fetch()

	addr := uint16(high)<<8 | uint16(low)

	val := cpu.mmu.Read(addr)

	cpu.setReg8("A", val)

	logger.Log("LD A, (%#02x)\n", addr)
}

// LDmCA put A into address 0xff00 + register C
func (cpu *CPU) LDmCA() {
	logger.Log("LD (C), A")

	val := cpu.getReg8("A")

	addr := 0xff00 + uint16(cpu.getReg8("C"))

	cpu.mmu.Write(addr, val)
}

// LDAmC put value at address 0xff00 + register C into A
func (cpu *CPU) LDAmC() {
	logger.Log("LD A, (C)")

	addr := 0xff00 + uint16(cpu.getReg8("C"))

	val := cpu.mmu.Read(addr)

	cpu.setReg8("A", val)
}

// LDHAmd8 put value at address 0xff00 + d8 into A
func (cpu *CPU) LDHAmd8() {
	addr := 0xff00 + uint16(cpu.Fetch())

	val := cpu.mmu.Read(addr)

	cpu.setReg8("A", val)

	logger.Log("LD A, (%#02x)\n", addr)
}

// LDHmd8A put value A into address 0xff00 + d8
func (cpu *CPU) LDHmd8A() {
	addr := 0xff00 + uint16(cpu.Fetch())

	cpu.mmu.Write(addr, cpu.getReg8("A"))

	logger.Log("LD (%#02x), A\n", addr)
}

// LDDmHLA put A into memory address HL. Decrement HL
func (cpu *CPU) LDDmHLA() {
	logger.Log("LDD (HL), A\n")

	addr := cpu.getReg16("HL")
	cpu.mmu.Write(addr, cpu.getReg8("A"))

	addr--
	cpu.setReg16("HL", addr)
}

// LDr16d16 put value d16 into r16
func (cpu *CPU) LDr16d16(reg string) {
	low := cpu.Fetch()
	high := cpu.Fetch()

	nn := uint16(high)<<8 | uint16(low)
	logger.Log("LD %s, %#02x\n", reg, nn)

	cpu.setReg16(reg, nn)
}

//======================================================================
// Jumps
//======================================================================

// JRccs8 if current condition is true, add n to current address and jump to it
func (cpu *CPU) JRccs8(cc string) {
	var n int8 = int8(cpu.Fetch())

	switch cc {
	case "NZ":
		if !testBit(Z, cpu.getReg8("F")) {
			cpu.pc = uint16(int32(cpu.pc) + int32(n))
		}
	case "Z":
		if testBit(Z, cpu.getReg8("F")) {
			cpu.pc = uint16(int32(cpu.pc) + int32(n))
		}
	case "NC":
		if !testBit(C, cpu.getReg8("F")) {
			cpu.pc = uint16(int32(cpu.pc) + int32(n))
		}
	case "C":
		if testBit(C, cpu.getReg8("F")) {
			cpu.pc = uint16(int32(cpu.pc) + int32(n))
		}
	}
}

//======================================================================
// Arithmetic
//======================================================================

////////////////////////
// 8-bit

// XORr8 exclusive OR n with register A, result in A
func (cpu *CPU) XORr8(reg string) {
	var n uint8
	switch reg {
	case "#":
		n = cpu.Fetch()
		logger.Log("XOR %#02x\n", n)
	case "(HL)":
		addr := cpu.getReg16(reg)
		n = cpu.mmu.Read(addr)
		logger.Log("XOR %s\n", reg)
	default:
		n = cpu.getReg8(reg)
		logger.Log("XOR %s\n", reg)
	}

	val := cpu.getReg8("A") ^ n

	if val == 0 {
		cpu.setFlags(SET, RESET, RESET, RESET)
	} else {
		cpu.setFlags(RESET, RESET, RESET, RESET)
	}

	cpu.setReg8("A", val)
}

////////////////////////
// 16-bit

// INCr16 increment r16
func (cpu *CPU) INCr16(reg string) {
	cpu.setReg16(reg, cpu.getReg16(reg)+1)

	logger.Log("INC %s\n", reg)
}

// DECr16 decrement r16
func (cpu *CPU) DECr16(reg string) {
	cpu.setReg16(reg, cpu.getReg16(reg)-1)

	logger.Log("DEC %s\n", reg)
}

//======================================================================
// CB prefixed
//======================================================================

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
