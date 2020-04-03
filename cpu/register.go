package cpu

// Flag register bits
const (
	Z = 7 // Zero Flag
	N = 6 // Subtract Flag
	H = 5 // Half Carry Flag
	C = 4 // Carry Flag
)

func (cpu *CPU) getReg8(reg string) byte {
	switch reg {
	case "A":
		return cpu.a
	case "F":
		return cpu.f
	case "B":
		return cpu.b
	case "C":
		return cpu.c
	case "D":
		return cpu.d
	case "E":
		return cpu.e
	case "H":
		return cpu.h
	case "L":
		return cpu.l
	}
	return 0
}

func (cpu *CPU) setReg8(reg string, val byte) {
	switch reg {
	case "A":
		cpu.a = val
	case "F":
		cpu.f = val
	case "B":
		cpu.b = val
	case "C":
		cpu.c = val
	case "D":
		cpu.d = val
	case "E":
		cpu.e = val
	case "H":
		cpu.h = val
	case "L":
		cpu.l = val
	}
}

func (cpu *CPU) getReg16(reg string) uint16 {
	switch reg {
	case "AF":
		return uint16(cpu.a)<<8 | uint16(cpu.f)
	case "BC":
		return uint16(cpu.b)<<8 | uint16(cpu.c)
	case "DE":
		return uint16(cpu.d)<<8 | uint16(cpu.e)
	case "HL":
		return uint16(cpu.h)<<8 | uint16(cpu.l)
	}
	return 0
}

func (cpu *CPU) setReg16(reg string, val uint16) {
	switch reg {
	case "AF":
		cpu.a = byte(val >> 8 & 0xff)
		cpu.f = byte(val & 0xff)
	case "BC":
		cpu.b = byte(val >> 8 & 0xff)
		cpu.c = byte(val & 0xff)
	case "DE":
		cpu.d = byte(val >> 8 & 0xff)
		cpu.e = byte(val & 0xff)
	case "HL":
		cpu.h = byte(val >> 8 & 0xff)
		cpu.l = byte(val & 0xff)
	}
}
