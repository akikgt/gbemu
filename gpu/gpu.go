package gpu

const (
	screenWidth  = 160
	screenHeight = 144
)

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
	bgp         uint8 // 0xff47
	obp0        uint8 // 0xff48
	obp1        uint8 // 0xff49
	wy          uint8 // 0xff4a
	wx          uint8 // 0xff4b

	Pixels      []byte
	frameBuffer []byte
	tileSets    [384][8][8]uint8
}

func New() *GPU {
	gpu := &GPU{}

	gpu.Pixels = make([]byte, screenHeight*screenWidth*4)
	for y := 0; y < screenHeight; y++ {
		for x := 0; x < screenWidth; x++ {
			gpu.Pixels[(y*screenWidth+x)*4+0] = 0xff // R
			gpu.Pixels[(y*screenWidth+x)*4+1] = 0xff // G
			gpu.Pixels[(y*screenWidth+x)*4+2] = 0xff // B
			gpu.Pixels[(y*screenWidth+x)*4+3] = 0xff // A
		}
	}
	gpu.frameBuffer = make([]byte, screenHeight*screenWidth*4)
	for y := 0; y < screenHeight; y++ {
		for x := 0; x < screenWidth; x++ {
			gpu.frameBuffer[(y*screenWidth+x)*4+0] = 0xff // R
			gpu.frameBuffer[(y*screenWidth+x)*4+1] = 0xff // G
			gpu.frameBuffer[(y*screenWidth+x)*4+2] = 0xff // B
			gpu.frameBuffer[(y*screenWidth+x)*4+3] = 0xff // A
		}
	}

	return gpu
}

func (gpu *GPU) ResetFrame() {
	for y := 0; y < screenHeight; y++ {
		for x := 0; x < screenWidth; x++ {
			gpu.Pixels[(y*screenWidth+x)*4+0] = 0xff // R
			gpu.Pixels[(y*screenWidth+x)*4+1] = 0xff // G
			gpu.Pixels[(y*screenWidth+x)*4+2] = 0xff // B
			gpu.Pixels[(y*screenWidth+x)*4+3] = 0xff // A
		}
	}
}

func (gpu *GPU) DisplayTileSets() {
	// get tile sets
	for i := 0; i < 384; i++ {

		for y := 0; y < 8; y++ {
			// each tile data is 16 byte
			data1 := gpu.vram[i*16+y*2]
			data2 := gpu.vram[i*16+y*2+1]

			for x := 0; x < 8; x++ {
				b := 7 - x
				color := (data2>>b&1)<<1 | (data1 >> b & 1)

				gpu.tileSets[i][y][x] = color
			}

		}
	}

	// display
	for y := 0; y < screenHeight; y++ {
		for x := 0; x < screenWidth; x++ {
			tileNum := (y/8)*(screenWidth/8) + x/8
			color := gpu.tileSets[tileNum][y%8][x%8]

			switch color {
			case 0:
				gpu.Pixels[(y*screenWidth+x)*4+0] = 0xff // R
				gpu.Pixels[(y*screenWidth+x)*4+1] = 0xff // G
				gpu.Pixels[(y*screenWidth+x)*4+2] = 0xff // B
				gpu.Pixels[(y*screenWidth+x)*4+3] = 0xff // A
			case 1:
				gpu.Pixels[(y*screenWidth+x)*4+0] = 0x88 // R
				gpu.Pixels[(y*screenWidth+x)*4+1] = 0xcc // G
				gpu.Pixels[(y*screenWidth+x)*4+2] = 0x44 // B
				gpu.Pixels[(y*screenWidth+x)*4+3] = 0xff // A
			case 2:
				gpu.Pixels[(y*screenWidth+x)*4+0] = 0x77 // R
				gpu.Pixels[(y*screenWidth+x)*4+1] = 0x77 // G
				gpu.Pixels[(y*screenWidth+x)*4+2] = 0x77 // B
				gpu.Pixels[(y*screenWidth+x)*4+3] = 0xff // A
			case 3:
				gpu.Pixels[(y*screenWidth+x)*4+0] = 0x00 // R
				gpu.Pixels[(y*screenWidth+x)*4+1] = 0x00 // G
				gpu.Pixels[(y*screenWidth+x)*4+2] = 0x00 // B
				gpu.Pixels[(y*screenWidth+x)*4+3] = 0xff // A
			}

		}
	}
}

func (gpu *GPU) updateTileSets() {
	for i := 0; i < 384; i++ {
		for y := 0; y < 8; y++ {
			// each tile data is 16 byte
			data1 := gpu.vram[i*16+y*2]
			data2 := gpu.vram[i*16+y*2+1]

			for x := 0; x < 8; x++ {
				b := 7 - x
				color := (data2>>b&1)<<1 | (data1 >> b & 1)
				gpu.tileSets[i][y][x] = color
			}

		}
	}
}

func (gpu *GPU) RenderFrame() {
	// for y := 0; y < screenHeight; y++ {
	// 	for x := 0; x < screenWidth; x++ {
	// 		fmt.Printf("%02x ", gpu.Pixels[(y*screenWidth+x)*4+0])
	// 	}
	// 	fmt.Printf("\n")
	// }
	gpu.Pixels = gpu.frameBuffer
}

func (gpu *GPU) RenderTiles() {
	var base uint16 = 0x1800
	y := (gpu.scy + gpu.currentLine) & 255
	var tileRow uint16 = uint16(y/8) * 32
	for px := 0; px < 160; px++ {
		x := uint8(px) + gpu.scx
		tileCol := x / 8
		tileAddr := base + tileRow + uint16(tileCol)
		tileNum := gpu.vram[tileAddr]

		// how to designate y affects scrolling logo...TODO: research
		color := gpu.tileSets[tileNum][y%8][px%8]

		pixel := int(gpu.currentLine)*screenWidth + int(px)
		switch color {
		case 0:
			gpu.frameBuffer[pixel*4+0] = 0xff // R
			gpu.frameBuffer[pixel*4+1] = 0xff // G
			gpu.frameBuffer[pixel*4+2] = 0xff // B
			gpu.frameBuffer[pixel*4+3] = 0xff // A
		case 1:
			gpu.frameBuffer[pixel*4+0] = 0x88 // R
			gpu.frameBuffer[pixel*4+1] = 0xcc // G
			gpu.frameBuffer[pixel*4+2] = 0x44 // B
			gpu.frameBuffer[pixel*4+3] = 0xff // A
		case 2:
			gpu.frameBuffer[pixel*4+0] = 0x77 // R
			gpu.frameBuffer[pixel*4+1] = 0x77 // G
			gpu.frameBuffer[pixel*4+2] = 0x77 // B
			gpu.frameBuffer[pixel*4+3] = 0xff // A
		case 3:
			gpu.frameBuffer[pixel*4+0] = 0x00 // R
			gpu.frameBuffer[pixel*4+1] = 0x00 // G
			gpu.frameBuffer[pixel*4+2] = 0x00 // B
			gpu.frameBuffer[pixel*4+3] = 0xff // A
		default:
			gpu.frameBuffer[pixel*4+0] = 0xff // R
			gpu.frameBuffer[pixel*4+1] = 0xff // G
			gpu.frameBuffer[pixel*4+2] = 0xff // B
			gpu.frameBuffer[pixel*4+3] = 0xff // A
		}
	}
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
	case addr == 0xff47:
		return gpu.bgp
	case addr == 0xff48:
		return gpu.obp0
	case addr == 0xff49:
		return gpu.obp1
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
		gpu.updateTileSets()
	case addr == 0xff40:
		gpu.lcdc = val
	case addr == 0xff41:
		gpu.stat = val & 0xfc // bit 1-0 are Read Only
	case addr == 0xff42:
		gpu.scy = val
	case addr == 0xff43:
		gpu.scx = val
	case addr == 0xff44:
		gpu.currentLine = val
	case addr == 0xff45:
		gpu.lyc = val
	case addr == 0xff47:
		gpu.bgp = val
	case addr == 0xff48:
		gpu.obp0 = val
	case addr == 0xff49:
		gpu.obp1 = val
	case addr == 0xff4a:
		gpu.wy = val
	case addr == 0xff4b:
		gpu.wx = val
	}
}

func (gpu *GPU) Update(ticks uint8) {
	if gpu.lcdc&0x80 == 0 {
		return
	}

	gpu.counter += uint16(ticks)

	switch gpu.mode {

	// accessing OAM
	case 2:
		if gpu.counter >= 80 {
			gpu.counter -= 80
			gpu.mode = 3
			gpu.RenderTiles()
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
				// TODO: complete screen update

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

	// fmt.Printf("ly: %d\n", gpu.currentLine)
	// fmt.Printf("scy: %d\n", gpu.scy)
	// fmt.Printf("GPU counter: %d\n", gpu.counter)
}
