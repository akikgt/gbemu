package cpu

// Flag register bits
const (
	Z = 7 // Zero Flag
	N = 6 // Subtract Flag
	H = 5 // Half Carry Flag
	C = 4 // Carry Flag
)

func (cpu *CPU) getA() byte {
	return cpu.a
}

func (cpu *CPU) setA(val byte) {
	cpu.a = val
}

func (cpu *CPU) getF() byte {
	return cpu.f
}

func (cpu *CPU) setF(val byte) {
	cpu.f = val
}

func (cpu *CPU) getAF() uint16 {
	return uint16(cpu.a)<<8 | uint16(cpu.f)
}

func (cpu *CPU) setAF(val uint16) {
	cpu.a = byte(val >> 8)
	cpu.f = byte(val & 0xff)
}

func (cpu *CPU) getB() byte {
	return cpu.b
}

func (cpu *CPU) setB(val byte) {
	cpu.b = val
}

func (cpu *CPU) getC() byte {
	return cpu.c
}

func (cpu *CPU) setC(val byte) {
	cpu.c = val
}

func (cpu *CPU) getBC() uint16 {
	return uint16(cpu.b)<<8 | uint16(cpu.c)
}

func (cpu *CPU) setBC(val uint16) {
	cpu.b = byte(val >> 8)
	cpu.c = byte(val & 0xff)
}

func (cpu *CPU) getD() byte {
	return cpu.d
}

func (cpu *CPU) setD(val byte) {
	cpu.d = val
}

func (cpu *CPU) getE() byte {
	return cpu.e
}

func (cpu *CPU) setE(val byte) {
	cpu.e = val
}

func (cpu *CPU) getDE() uint16 {
	return uint16(cpu.d)<<8 | uint16(cpu.e)
}

func (cpu *CPU) setDE(val uint16) {
	cpu.d = byte(val >> 8)
	cpu.e = byte(val & 0xff)
}

func (cpu *CPU) getH() byte {
	return cpu.h
}

func (cpu *CPU) setH(val byte) {
	cpu.h = val
}

func (cpu *CPU) getL() byte {
	return cpu.l
}

func (cpu *CPU) setL(val byte) {
	cpu.l = val
}

func (cpu *CPU) getHL() uint16 {
	return uint16(cpu.h)<<8 | uint16(cpu.l)
}

func (cpu *CPU) setHL(val uint16) {
	cpu.h = byte(val >> 8)
	cpu.l = byte(val & 0xff)
}
