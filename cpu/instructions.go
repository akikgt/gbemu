package cpu

func (cpu *CPU) setFlags(Z, N, H, C uint8) {
	var newFlag uint8 = 0

	newFlag |= Z << 7
	newFlag |= N << 6
	newFlag |= H << 5
	newFlag |= C << 4

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

// XORn exclusive OR n with register A, result in A
func (cpu *CPU) XORn(reg string) {
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
