package gpu

import "fmt"

type GPU struct {
	counter uint16

	vram [0x2000]uint8

	lcdc        uint8 // 0xff40
	stat        uint8 // 0xff41
	mode        uint8 // 0xff41 bit 0-1
	scy         uint8 // 0xff42
	scx         uint8 // 0xff43
	currentLine uint8 // 0xff44
	lyc         uint8 // 0xff45
	wy          uint8 // 0xff4a
	wx          uint8 // 0xff4b
}

func New() *GPU {
	gpu := &GPU{}

	return gpu
}

func (gpu *GPU) Read(addr uint16) uint8 {
	switch {
	case 0x8000 <= addr && addr <= 0x9fff:
		return gpu.vram[addr-0x8000]
	case addr == 0xff40:
		return gpu.lcdc
	case addr == 0xff41:
		return gpu.stat | gpu.mode
	case addr == 0xff42:
		return gpu.scy
	case addr == 0xff43:
		return gpu.scx
	case addr == 0xff44:
		return gpu.currentLine
	case addr == 0xff45:
		return gpu.lyc
	case addr == 0xff4a:
		return gpu.wy
	case addr == 0xff4b:
		return gpu.wx
	}

	return gpu.vram[addr]
}

func (gpu *GPU) Write(addr uint16, val uint8) {
	switch {
	case 0x8000 <= addr && addr <= 0x9fff:
		gpu.vram[addr-0x8000] = val
	case addr == 0xff40:
		gpu.lcdc = val
	case addr == 0xff41:
		gpu.stat = val & 0xfc // bit 1-0 are Read Only
	case addr == 0xff43:
		gpu.scx = val
	case addr == 0xff44:
		gpu.currentLine = val
	case addr == 0xff45:
		gpu.lyc = val
	case addr == 0xff4a:
		gpu.wy = val
	case addr == 0xff4b:
		gpu.wx = val
	}
}

func (gpu *GPU) Update(ticks uint8) {
	gpu.counter += uint16(ticks)

	// gpu.counter %= 456
	switch gpu.mode {

	// accessing OAM
	case 2:
		if gpu.counter >= 80 {
			gpu.counter -= 80
			gpu.mode = 3
		}

	// accessing VRAM
	case 3:
		if gpu.counter >= 172 {
			gpu.counter -= 172
			gpu.mode = 0
		}

	// horizontal blank
	case 0:
		if gpu.counter >= 204 {
			gpu.counter -= 204
			gpu.currentLine++

			if gpu.currentLine >= 143 {
				// enter vblank mode
				gpu.mode = 1
				// TODO: screen update
				// for i := 0; i < 8096; i++ {
				// 	if i%32 == 0 {
				// 		fmt.Printf("\n")
				// 		fmt.Printf("%#04x: ", i+0x8000)
				// 		fmt.Printf("%02x ", gpu.vram[i])
				// 		continue
				// 	}
				// 	fmt.Printf("%02x ", gpu.vram[i])
				// }

			} else {
				// back to accessing OAM mode
				gpu.mode = 2
			}
		}

	// vertical blank
	case 1:
		if gpu.counter >= 456 {
			gpu.counter -= 456
			gpu.currentLine++

			if gpu.currentLine > 153 {
				// back to accessing OAM mode
				gpu.mode = 2
				gpu.currentLine = 0
			}
		}

	}

	fmt.Printf("ly: %d\n", gpu.currentLine)
	fmt.Printf("GPU counter: %d\n", gpu.counter)
}
