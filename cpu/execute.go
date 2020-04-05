package cpu

func parseBit(opcode, base uint8) uint8 {
	// (higher 4-bit - base) * 2 + bit3
	return (opcode>>4-base)*2 + (opcode >> 3 & 1)
}

func parseReg(opcode uint8) string {
	val := opcode & 0xf

	switch val {
	case 0x7, 0xf:
		return "A"
	case 0x0, 0x8:
		return "B"
	case 0x1, 0x9:
		return "C"
	case 0x2, 0xa:
		return "D"
	case 0x3, 0xb:
		return "E"
	case 0x4, 0xc:
		return "H"
	case 0x5, 0xd:
		return "L"
	case 0x6, 0xe:
		return "(HL)"
	}

	return "unknown"
}

func (cpu *CPU) Execute() {
	opcode := cpu.Fetch()

	switch opcode {
	case 0x00:
		cpu.NOP()

	// 8-bit loads
	// LD nn, n
	case 0x06:
		cpu.LDr8d8("B")
	case 0x0e:
		cpu.LDr8d8("C")
	case 0x16:
		cpu.LDr8d8("D")
	case 0x1e:
		cpu.LDr8d8("E")
	case 0x26:
		cpu.LDr8d8("H")
	case 0x2e:
		cpu.LDr8d8("L")
	case 0x3e:
		cpu.LDr8d8("A")

	// LD r1, r2
	case 0x78:
		cpu.LDr8r8("A", "B")
	case 0x79:
		cpu.LDr8r8("A", "C")
	case 0x7a:
		cpu.LDr8r8("A", "D")
	case 0x7b:
		cpu.LDr8r8("A", "E")
	case 0x7c:
		cpu.LDr8r8("A", "H")
	case 0x7d:
		cpu.LDr8r8("A", "L")
	case 0x7f:
		cpu.LDr8r8("A", "A")
	case 0x40:
		cpu.LDr8r8("B", "B")
	case 0x41:
		cpu.LDr8r8("B", "C")
	case 0x42:
		cpu.LDr8r8("B", "D")
	case 0x43:
		cpu.LDr8r8("B", "E")
	case 0x44:
		cpu.LDr8r8("B", "H")
	case 0x45:
		cpu.LDr8r8("B", "L")
	case 0x47:
		cpu.LDr8r8("B", "A")
	case 0x48:
		cpu.LDr8r8("C", "B")
	case 0x49:
		cpu.LDr8r8("C", "C")
	case 0x4a:
		cpu.LDr8r8("C", "D")
	case 0x4b:
		cpu.LDr8r8("C", "E")
	case 0x4c:
		cpu.LDr8r8("C", "H")
	case 0x4d:
		cpu.LDr8r8("C", "L")
	case 0x4f:
		cpu.LDr8r8("C", "A")
	case 0x50:
		cpu.LDr8r8("D", "B")
	case 0x51:
		cpu.LDr8r8("D", "C")
	case 0x52:
		cpu.LDr8r8("D", "D")
	case 0x53:
		cpu.LDr8r8("D", "E")
	case 0x54:
		cpu.LDr8r8("D", "H")
	case 0x55:
		cpu.LDr8r8("D", "L")
	case 0x57:
		cpu.LDr8r8("D", "A")
	case 0x58:
		cpu.LDr8r8("E", "B")
	case 0x59:
		cpu.LDr8r8("E", "C")
	case 0x5a:
		cpu.LDr8r8("E", "D")
	case 0x5b:
		cpu.LDr8r8("E", "E")
	case 0x5c:
		cpu.LDr8r8("E", "H")
	case 0x5d:
		cpu.LDr8r8("E", "L")
	case 0x5f:
		cpu.LDr8r8("E", "A")
	case 0x60:
		cpu.LDr8r8("H", "B")
	case 0x61:
		cpu.LDr8r8("H", "C")
	case 0x62:
		cpu.LDr8r8("H", "D")
	case 0x63:
		cpu.LDr8r8("H", "E")
	case 0x64:
		cpu.LDr8r8("H", "H")
	case 0x65:
		cpu.LDr8r8("H", "L")
	case 0x67:
		cpu.LDr8r8("H", "A")
	case 0x68:
		cpu.LDr8r8("L", "B")
	case 0x69:
		cpu.LDr8r8("L", "C")
	case 0x6a:
		cpu.LDr8r8("L", "D")
	case 0x6b:
		cpu.LDr8r8("L", "E")
	case 0x6c:
		cpu.LDr8r8("L", "H")
	case 0x6d:
		cpu.LDr8r8("L", "L")
	case 0x6f:
		cpu.LDr8r8("L", "A")

	// Put value at memory into r8
	case 0x0a:
		cpu.LDr8mr16("A", "BC")
	case 0x1a:
		cpu.LDr8mr16("A", "DE")
	case 0x7e:
		cpu.LDr8mr16("A", "HL")
	case 0x46:
		cpu.LDr8mr16("B", "HL")
	case 0x4e:
		cpu.LDr8mr16("C", "HL")
	case 0x56:
		cpu.LDr8mr16("D", "HL")
	case 0x5e:
		cpu.LDr8mr16("E", "HL")
	case 0x66:
		cpu.LDr8mr16("H", "HL")
	case 0x6e:
		cpu.LDr8mr16("L", "HL")

	case 0x2a:
		cpu.LDIAmHL()
	case 0x3a:
		cpu.LDDAmHL()

	case 0xf0:
		cpu.LDHAmd8()
	case 0xf2:
		cpu.LDAmC()
	case 0xfa:
		cpu.LDAmd16()

	// Put value into memory
	case 0x02:
		cpu.LDmr16r8("BC", "A")
	case 0x12:
		cpu.LDmr16r8("DE", "A")
	case 0x70:
		cpu.LDmr16r8("HL", "B")
	case 0x71:
		cpu.LDmr16r8("HL", "C")
	case 0x72:
		cpu.LDmr16r8("HL", "D")
	case 0x73:
		cpu.LDmr16r8("HL", "E")
	case 0x74:
		cpu.LDmr16r8("HL", "H")
	case 0x75:
		cpu.LDmr16r8("HL", "L")
	case 0x77:
		cpu.LDmr16r8("HL", "A")

	case 0x22:
		cpu.LDImHLA()
	case 0x32:
		cpu.LDDmHLA()

	case 0x36:
		cpu.LDmHLd8()
	case 0xe0:
		cpu.LDHmd8A()
	case 0xe2:
		cpu.LDmCA()
	case 0xea:
		cpu.LDmd16A()

	// 16-bit loads
	case 0x01:
		cpu.LDr16d16("BC")
	case 0x11:
		cpu.LDr16d16("DE")
	case 0x21:
		cpu.LDr16d16("HL")
	case 0x31:
		cpu.LDr16d16("SP")
	case 0x08:
		cpu.LDmd16SP()
	case 0xf8:
		cpu.LDHLSPs8()
	case 0xf9:
		cpu.LDr16r16("SP", "HL")
	case 0xf5:
		cpu.PUSHr16("AF")
	case 0xc5:
		cpu.PUSHr16("BC")
	case 0xd5:
		cpu.PUSHr16("DE")
	case 0xe5:
		cpu.PUSHr16("HL")
	case 0xf1:
		cpu.POPr16("AF")
	case 0xc1:
		cpu.POPr16("BC")
	case 0xd1:
		cpu.POPr16("DE")
	case 0xe1:
		cpu.POPr16("HL")

	// 8-bit ALU
	// Rotate & Shifts
	case 0x17:
		cpu.RLr8("A")

	// ADD n
	case 0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87:
		reg := parseReg(opcode)
		cpu.ADDr8(reg)
	case 0xc6:
		cpu.ADDr8("#")

	// ADC n
	case 0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f:
		reg := parseReg(opcode)
		cpu.ADCr8(reg)
	case 0xce:
		cpu.ADCr8("#")

	// SUB n
	case 0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97:
		reg := parseReg(opcode)
		cpu.SUBr8(reg)
	case 0xd6:
		cpu.SUBr8("#")

	// SBC n
	case 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f:
		reg := parseReg(opcode)
		cpu.SBCr8(reg)
	case 0xde:
		cpu.ADCr8("#")

	// AND n
	case 0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7:
		reg := parseReg(opcode)
		cpu.ANDr8(reg)
	case 0xe6:
		cpu.ANDr8("#")

	// OR n
	case 0xb0, 0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7:
		reg := parseReg(opcode)
		cpu.ORr8(reg)
	case 0xf6:
		cpu.ORr8("#")

	// XOR n
	case 0xa8, 0xa9, 0xaa, 0xab, 0xac, 0xad, 0xae, 0xaf:
		reg := parseReg(opcode)
		cpu.XORr8(reg)
	case 0xee:
		cpu.XORr8("#")

	// INC n
	case 0x3c:
		cpu.INCr8("A")
	case 0x04:
		cpu.INCr8("B")
	case 0x0c:
		cpu.INCr8("C")
	case 0x14:
		cpu.INCr8("D")
	case 0x1c:
		cpu.INCr8("E")
	case 0x24:
		cpu.INCr8("H")
	case 0x2c:
		cpu.INCr8("L")
	case 0x34:
		cpu.INCr8("(HL)")

	// DEC n
	case 0x3d:
		cpu.DECr8("A")
	case 0x05:
		cpu.DECr8("B")
	case 0x0d:
		cpu.DECr8("C")
	case 0x15:
		cpu.DECr8("D")
	case 0x1d:
		cpu.DECr8("E")
	case 0x25:
		cpu.DECr8("H")
	case 0x2d:
		cpu.DECr8("L")
	case 0x35:
		cpu.DECr8("(HL)")

	// CP n
	case 0xb8, 0xb9, 0xba, 0xbb, 0xbc, 0xbd, 0xbe, 0xbf:
		reg := parseReg(opcode)
		cpu.CPr8(reg)
	case 0xfe:
		cpu.CPr8("#")

	// 16-bit ALU
	// ADD nn
	case 0x09:
		cpu.ADDHLr16("BC")
	case 0x19:
		cpu.ADDHLr16("DE")
	case 0x29:
		cpu.ADDHLr16("HL")
	case 0x39:
		cpu.ADDHLr16("SP")

	// INC nn
	case 0x03:
		cpu.INCr16("BC")
	case 0x13:
		cpu.INCr16("DE")
	case 0x23:
		cpu.INCr16("HL")
	case 0x33:
		cpu.INCr16("SP")
	// DEC nn
	case 0x0B:
		cpu.DECr16("BC")
	case 0x1B:
		cpu.DECr16("DE")
	case 0x2B:
		cpu.DECr16("HL")
	case 0x3B:
		cpu.DECr16("SP")

	///////////
	// Jumps
	// JP nn
	case 0xc3:
		cpu.JPd16()

	// JP (HL)
	case 0xe9:
		cpu.JPHL()

	// JR n
	case 0x18:
		cpu.JRsd8()

	// JR cc, n
	case 0x20:
		cpu.JRccs8("NZ")
	case 0x28:
		cpu.JRccs8("Z")
	case 0x30:
		cpu.JRccs8("NC")
	case 0x38:
		cpu.JRccs8("C")

	///////////
	// Calls
	case 0xcd:
		cpu.CALLd16()

	case 0xc4:
		cpu.CALLccd16("NZ")
	case 0xcc:
		cpu.CALLccd16("Z")
	case 0xd4:
		cpu.CALLccd16("NC")
	case 0xdc:
		cpu.CALLccd16("C")

	///////////
	// Restarts
	case 0xc7:
		cpu.RSTd16(0x00)
	case 0xcf:
		cpu.RSTd16(0x08)
	case 0xd7:
		cpu.RSTd16(0x10)
	case 0xdf:
		cpu.RSTd16(0x18)
	case 0xe7:
		cpu.RSTd16(0x20)
	case 0xef:
		cpu.RSTd16(0x28)
	case 0xf7:
		cpu.RSTd16(0x30)
	case 0xff:
		cpu.RSTd16(0x38)

	///////////
	// Returns
	case 0xc9:
		cpu.RET()

	// CB-prefixed
	case 0xcb:
		logger.Log("CB-prefixed\n")
		cpu.CBPrefixed()

	default:
		logger.Log("unknown opcode: %#02x\n", opcode)
	}
}

func (cpu *CPU) CBPrefixed() {
	opcode := cpu.Fetch()

	reg := parseReg(opcode)

	switch {
	// RL
	case 0x10 <= opcode && opcode <= 0x17:
		cpu.RLr8(reg)
	// BIT
	case 0x40 <= opcode && opcode <= 0x7f:
		b := parseBit(opcode, 4)
		cpu.BITbr8(b, reg)
	}
}
