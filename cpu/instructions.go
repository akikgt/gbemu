package cpu

// d8  means immediate 8 bit data
// d16 means immediate 16 bit data
// s8  means signed immediate 8 bit data, which are added to pc
// r8  means 8 bit register
// r16 means 16 bit register
// m*  means data in memory address *
// b   means bit

func (cpu *CPU) getd8(src string) uint8 {
	var n uint8

	switch src {
	case "#":
		n = cpu.Fetch()
	case "(HL)":
		n = cpu.mmu.Read(cpu.getReg16("HL"))
	default:
		n = cpu.getReg8(src)
	}

	return n
}

func signExtend(a uint8) uint16 {
	return uint16(int8(a))
}

// param: Z, N, H, C
func (cpu *CPU) getFlag(f uint8) uint8 {
	return cpu.getReg8("F") >> f & 1
}

func checkHalfCarry(a, b, c uint8) uint8 {
	if ((a&0xf)+(b&0xf)+c)&0x10 == 0x10 {
		return SET
	}
	return RESET
}

func checkCarry(a, b, c uint8) uint8 {
	sum := uint16(a) + uint16(b) + uint16(c)

	if sum&0x100 == 0x100 {
		return SET
	}
	return RESET
}

func checkHalfBorrow(a, b, c uint8) uint8 {
	if (a & 0xf) < (b&0xf + c) {
		return RESET
	}
	return SET
}

func checkBorrow(a, b, c uint8) uint8 {
	if a < b+c {
		return RESET
	}
	return SET
}

func checkZero(a uint8) uint8 {
	if a == 0 {
		return SET
	}
	return RESET
}

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
	addr := cpu.getReg16(reg2)

	val := cpu.mmu.Read(addr)

	cpu.setReg8(reg1, val)

	logger.Log("LD %s, (%s)\n", reg1, reg2)
}

// LDmr16r8 put value into address r16
func (cpu *CPU) LDmr16r8(reg1, reg2 string) {
	logger.Log("LD (%s), %s\n", reg1, reg2)

	addr := cpu.getReg16(reg1)

	val := cpu.getReg8(reg2)

	cpu.mmu.Write(addr, val)
}

// LDmd16A put value A into address d16
func (cpu *CPU) LDmd16A() {
	addr := cpu.FetchWord()

	cpu.mmu.Write(addr, cpu.getReg8("A"))

	logger.Log("LD (%#02x), A\n", addr)
}

// LDAmd16 put value at address d16 into A
func (cpu *CPU) LDAmd16() {
	addr := cpu.FetchWord()

	val := cpu.mmu.Read(addr)

	cpu.setReg8("A", val)

	logger.Log("LD A, (%#02x)\n", addr)
}

// LDmCA put A into address 0xff00 + register C
func (cpu *CPU) LDmCA() {
	val := cpu.getReg8("A")

	addr := 0xff00 + uint16(cpu.getReg8("C"))

	cpu.mmu.Write(addr, val)

	logger.Log("LD (C), A\n")
}

// LDAmC put value at address 0xff00 + register C into A
func (cpu *CPU) LDAmC() {
	addr := 0xff00 + uint16(cpu.getReg8("C"))

	val := cpu.mmu.Read(addr)

	cpu.setReg8("A", val)

	logger.Log("LD A, (C)\n")
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

// LDImHLA put A into memory address HL. Increment HL
func (cpu *CPU) LDImHLA() {
	logger.Log("LDI (HL), A\n")

	cpu.LDmr16r8("HL", "A")

	cpu.INCr16("HL")
}

// LDIAmHL put value at address HL into A. Increment HL
func (cpu *CPU) LDIAmHL() {
	logger.Log("LDI A, (HL)\n")

	cpu.LDr8mr16("A", "HL")

	cpu.INCr16("HL")
}

// LDDmHLA put A into memory address HL. Decrement HL
func (cpu *CPU) LDDmHLA() {
	logger.Log("LDD (HL), A\n")

	cpu.LDmr16r8("HL", "A")

	cpu.DECr16("HL")
}

// LDDAmHL put value at address HL into A. Decrement HL
func (cpu *CPU) LDDAmHL() {
	logger.Log("LDD A, (HL)\n")

	cpu.LDr8mr16("A", "HL")

	cpu.DECr16("HL")
}

////////////////////////
// 16-bit

// LDr16d16 put value d16 into r16
func (cpu *CPU) LDr16d16(reg string) {
	nn := cpu.FetchWord()

	cpu.setReg16(reg, nn)

	logger.Log("LD %s, %#04x\n", reg, nn)
}

// LDmd16SP put SP at address d16
func (cpu *CPU) LDmd16SP() {
	addr := cpu.FetchWord()

	sp := cpu.getReg16("SP")

	cpu.mmu.WriteWord(addr, sp)

	logger.Log("LD (%#04x), SP\n", addr)
}

// LDr16r16 put reg2 into reg1
func (cpu *CPU) LDr16r16(reg1, reg2 string) {
	cpu.setReg16(reg1, cpu.getReg16(reg2))

	logger.Log("LD %s, %s\n", reg1, reg2)
}

// LDHLSPs8 put SP + s8 effective adress into HL
func (cpu *CPU) LDHLSPs8() {
	n := cpu.Fetch()

	sp := cpu.getReg16("SP")

	c := checkCarry(uint8(n), uint8(sp&0xff), 0)
	h := checkHalfCarry(uint8(n), uint8(sp&0xff), 0)
	cpu.setFlags(RESET, RESET, h, c)

	cpu.setReg16("HL", sp+signExtend(n))

	logger.Log("LDHL SP, %#02x\n", n)
}

// PUSHr16 decrement SP twice and push register r16 onto stack.
func (cpu *CPU) PUSHr16(reg string) {
	addr := cpu.getReg16("SP") - 2
	cpu.setReg16("SP", addr)

	cpu.mmu.WriteWord(addr, cpu.getReg16(reg))

	logger.Log("PUSH %s\n", reg)
}

// POPr16 pop two bytes off stack into register r16. Increment SP twice
func (cpu *CPU) POPr16(reg string) {
	addr := cpu.getReg16("SP")
	val := cpu.mmu.ReadWord(addr)

	cpu.setReg16("SP", addr+2)

	cpu.setReg16(reg, val)

	logger.Log("POP %s\n", reg)
}

//======================================================================
// Jumps
//======================================================================

// JPd16 jump to address d16
func (cpu *CPU) JPd16() {
	cpu.pc = cpu.FetchWord()
}

// JPHL jump to address contained in HL
func (cpu *CPU) JPHL() {
	cpu.pc = cpu.getReg16("HL")
}

// JRsd8 add sd8 to current address and jump to it
func (cpu *CPU) JRsd8() {
	cpu.pc += signExtend(cpu.Fetch())
}

// JRccs8 if current condition is true, add n to current address and jump to it
func (cpu *CPU) JRccs8(cc string) {
	n := cpu.Fetch()

	logger.Log("JR %s, %#04x\n", cc, n)

	switch cc {
	case "NZ":
		if testBit(Z, cpu.getReg8("F")) {
			return
		}
	case "Z":
		if !testBit(Z, cpu.getReg8("F")) {
			return
		}
	case "NC":
		if testBit(C, cpu.getReg8("F")) {
			return
		}
	case "C":
		if !testBit(C, cpu.getReg8("F")) {
			return
		}
	}

	cpu.pc = cpu.pc + signExtend(n)
}

//======================================================================
// Calls
//======================================================================

func (cpu *CPU) pushd16(d uint16) {
	addr := cpu.getReg16("SP") - 2
	cpu.setReg16("SP", addr)

	cpu.mmu.WriteWord(addr, d)
}

// CALLd16 push address of next instruction onto stack
// and then jump to address d16
func (cpu *CPU) CALLd16() {
	jumpTo := cpu.FetchWord()

	cpu.pushd16(cpu.pc)

	cpu.pc = jumpTo

	logger.Log("CALL %#04x\n", jumpTo)
}

//======================================================================
// Returns
//======================================================================

func (cpu *CPU) popd16() uint16 {
	addr := cpu.getReg16("SP")
	cpu.setReg16("SP", addr+2)

	return cpu.mmu.ReadWord(addr)
}

// RET pop two bytes from stack & jump to that address
func (cpu *CPU) RET() {
	cpu.pc = cpu.popd16()

	logger.Log("RET %#04x\n", cpu.pc)
}

//======================================================================
// Arithmetic
//======================================================================

////////////////////////
// 8-bit

// ADDr8 add r8 to A
func (cpu *CPU) ADDr8(reg string) {
	n := cpu.getd8(reg)
	a := cpu.getReg8("A")

	z := checkZero(a + n)
	h := checkHalfCarry(a, n, 0)
	c := checkCarry(a, n, 0)
	cpu.setFlags(z, RESET, h, c)

	cpu.setReg8("A", a+n)

	logger.Log("ADD %s(=%#02x)\n", reg, n)
}

// ADCr8 add r8 + carry flag to A
func (cpu *CPU) ADCr8(reg string) {
	n := cpu.getd8(reg) + cpu.getFlag(C)
	a := cpu.getReg8("A")

	z := checkZero(a + n)
	h := checkHalfCarry(a, n, 0)
	c := checkCarry(a, n, 0)
	cpu.setFlags(z, RESET, h, c)

	cpu.setReg8("A", a+n)

	logger.Log("ADC %s(n=%#02x)\n", reg, n)
}

// SUBr8 subtract r8 from A
func (cpu *CPU) SUBr8(reg string) {
	n := cpu.getd8(reg)
	a := cpu.getReg8("A")

	z := checkZero(a - n)
	h := checkHalfBorrow(a, n, 0)
	c := checkBorrow(a, n, 0)
	cpu.setFlags(z, SET, h, c)

	cpu.setReg8("A", a-n)

	logger.Log("SUB %s(=%#02x)\n", reg, n)
}

// SBCr8 subtract r8 + carry flag from A
func (cpu *CPU) SBCr8(reg string) {
	n := cpu.getd8(reg) + cpu.getFlag(C)
	a := cpu.getReg8("A")

	z := checkZero(a - n)
	h := checkHalfBorrow(a, n, 0)
	c := checkBorrow(a, n, 0)
	cpu.setFlags(z, SET, h, c)

	cpu.setReg8("A", a-n)

	logger.Log("SUB %s(=%#02x)\n", reg, n)
}

// XORr8 exclusive OR n with register A, result in A
func (cpu *CPU) XORr8(reg string) {
	n := cpu.getd8(reg)

	val := cpu.getReg8("A") ^ n

	z := checkZero(val)
	cpu.setFlags(z, RESET, RESET, RESET)

	cpu.setReg8("A", val)

	logger.Log("XOR %s\n", reg)
}

// INCr8 increment r8
func (cpu *CPU) INCr8(reg string) {
	n := cpu.getd8(reg)

	z := checkZero(n + 1)
	h := checkHalfCarry(n, 1, 0)
	cpu.setFlags(z, RESET, h, NA)

	cpu.setReg8(reg, n+1)

	logger.Log("INC %s\n", reg)
}

// DECr8 decrement r8
func (cpu *CPU) DECr8(reg string) {
	n := cpu.getd8(reg)

	z := checkZero(n - 1)
	h := checkHalfBorrow(n, 1, 0)
	cpu.setFlags(z, SET, h, NA)

	cpu.setReg8(reg, n-1)

	logger.Log("DEC %s\n", reg)
}

// CPr8 compare A with r8.
// This is basically an A - n subtraction instruction.
// but the result is thrown away
func (cpu *CPU) CPr8(reg string) {
	n := cpu.getd8(reg)
	a := cpu.getReg8("A")

	z := checkZero(a - n)
	h := checkHalfBorrow(a, n, 0)
	var c uint8 = RESET
	if a < n {
		c = SET
	}
	cpu.setFlags(z, SET, h, c)

	logger.Log("CP %s\n", reg)
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
	val := cpu.getd8(reg)

	if testBit(b, val) {
		cpu.setFlags(RESET, RESET, SET, NA)
	} else {
		cpu.setFlags(SET, RESET, SET, NA)
	}

	logger.Log("BIT %d, %s\n", b, reg)
}

// RLr8 rotate r8 left through carry flag
func (cpu *CPU) RLr8(reg string) {
	val := cpu.getd8(reg)

	res := val<<1 | cpu.getFlag(C)

	z := checkZero(res)

	var c uint8 = RESET
	if val>>7 == 1 {
		c = SET
	}

	cpu.setFlags(z, RESET, RESET, c)

	logger.Log("RL %s\n", reg)
}
