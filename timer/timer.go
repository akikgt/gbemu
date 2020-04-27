package timer

type Timer struct {
	div  uint8 // 0xff04
	tima uint8 // 0xff05
	tma  uint8 // 0xff06
	tac  uint8 // 0xff07

	dividerCounter uint16
	timerCounter   uint16

	ReqTimerInt bool
}

func New() *Timer {
	timer := &Timer{}

	timer.ReqTimerInt = false

	return timer
}

func (timer *Timer) Write(addr uint16, val uint8) {
	switch addr {
	case 0xff04:
		// writing any value to this register resets it to 0
		timer.div = 0
	case 0xff05:
		timer.tima = val
	case 0xff06:
		timer.tma = val
	case 0xff07:
		timer.tac = val
	}
}

func (timer *Timer) Read(addr uint16) uint8 {
	switch addr {
	case 0xff04:
		return timer.div
	case 0xff05:
		return timer.tima
	case 0xff06:
		return timer.tma
	case 0xff07:
		return timer.tac
	}

	return 0
}

func (timer *Timer) isTimerEnabled() bool {
	return timer.tac&0x4 > 0
}

func (timer *Timer) getThreshold() uint16 {
	switch timer.tac & 0x3 {
	case 0:
		// freq 4096. GB CPU speed is 4194304Hz
		// 4194304 / 4096 = 1024
		return 1024
	case 1:
		return 16 // freq 262144
	case 2:
		return 64 // freq 65536
	case 3:
		return 256 // freq 16384
	}

	return 1024
}

func (timer *Timer) updateDiv(ticks uint8) {
	timer.dividerCounter += uint16(ticks)

	if timer.dividerCounter >= 0xff {
		timer.dividerCounter = 0
		timer.div++
	}
}

func (timer *Timer) Update(ticks uint8) {
	timer.ReqTimerInt = false

	timer.updateDiv(ticks)

	if !timer.isTimerEnabled() {
		return
	}

	timer.timerCounter += uint16(ticks)

	if timer.timerCounter < timer.getThreshold() {
		return
	}

	// reset counter
	timer.timerCounter = 0

	// update tima
	if timer.tima == 0xff {
		timer.tima = timer.tma
		timer.ReqTimerInt = true
	} else {
		timer.tima++
	}
}
