package joypad

type Joypad struct {
	state uint8

	ReqJoypadInt bool
}

func New() *Joypad {
	joypad := &Joypad{}

	joypad.state = 0xff

	joypad.ReqJoypadInt = false

	return joypad
}

func (joypad *Joypad) Write(val uint8) {
	joypad.state = val
}

func (joypad *Joypad) Read() uint8 {
	return joypad.state
}
