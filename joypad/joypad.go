package joypad

const (
	START = iota
	SELECT
	A
	B
	DOWN
	UP
	LEFT
	RIGHT
)

type Joypad struct {
	state uint8

	buttonKeys    uint8
	directionKeys uint8

	ReqJoypadInt bool
}

func New() *Joypad {
	joypad := &Joypad{}

	joypad.state = 0xff

	joypad.buttonKeys = 0xf
	joypad.directionKeys = 0xf

	joypad.ReqJoypadInt = false

	return joypad
}

func (joypad *Joypad) Write(val uint8) {
	joypad.state = (joypad.state & 0xcf) | val&0x30 // bit 0 - 3 is Read Only, 6, 7 are not used
}

func (joypad *Joypad) Read() uint8 {
	if joypad.state&0x10 == 0 {
		return joypad.state&0xf0 | joypad.directionKeys
	}

	if joypad.state&0x20 == 0 {
		return joypad.state&0xf0 | joypad.buttonKeys
	}

	return joypad.state
}

func (joypad *Joypad) KeyPress(key uint8) {
	switch key {

	case DOWN:
		joypad.directionKeys &= 0x7
	case UP:
		joypad.directionKeys &= 0xb
	case LEFT:
		joypad.directionKeys &= 0xd
	case RIGHT:
		joypad.directionKeys &= 0xe

	case START:
		joypad.buttonKeys &= 0x7
	case SELECT:
		joypad.buttonKeys &= 0xb
	case B:
		joypad.buttonKeys &= 0xd
	case A:
		joypad.buttonKeys &= 0xe
	}

	joypad.ReqJoypadInt = true
}

func (joypad *Joypad) KeyRelease(key uint8) {
	switch key {

	case DOWN:
		joypad.directionKeys |= 0x8
	case UP:
		joypad.directionKeys |= 0x4
	case LEFT:
		joypad.directionKeys |= 0x2
	case RIGHT:
		joypad.directionKeys |= 0x1

	case START:
		joypad.buttonKeys |= 0x8
	case SELECT:
		joypad.buttonKeys |= 0x4
	case B:
		joypad.buttonKeys |= 0x2
	case A:
		joypad.buttonKeys |= 0x1
	}
}

func (joypad *Joypad) ReleaseAll() {
	joypad.KeyRelease(DOWN)
	joypad.KeyRelease(UP)
	joypad.KeyRelease(LEFT)
	joypad.KeyRelease(RIGHT)
	joypad.KeyRelease(START)
	joypad.KeyRelease(SELECT)
	joypad.KeyRelease(B)
	joypad.KeyRelease(A)
	joypad.ReqJoypadInt = false
}
