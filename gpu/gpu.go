package gpu

const (
	screenWidth  = 160
	screenHeight = 144

	// palette colors for Non CGB mode
	white     uint8 = 0
	ligthGray uint8 = 1
	darkGray  uint8 = 2
	black     uint8 = 3
)

type GPU struct {
	counter uint16

	vram [0x2000]uint8

	lcdc uint8 // 0xff40 LCD control
	stat uint8 // 0xff41 LCDC status
	mode uint8 // 0xff41 bit 0-1
	scy  uint8 // 0xff42
	scx  uint8 // 0xff43
	ly   uint8 // 0xff44 current Y-coordinate
	lyc  uint8 // 0xff45
	dma  uint8 // 0xff46
	bgp  uint8 // 0xff47 bg palette data
	obp0 uint8 // 0xff48
	obp1 uint8 // 0xff49
	wy   uint8 // 0xff4a
	wx   uint8 // 0xff4b

	Pixels      []byte
	frameBuffer []byte
	tileSets    [384][8][8]uint8
}

func New() *GPU {
	gpu := &GPU{}

	gpu.Pixels = make([]byte, screenHeight*screenWidth*4) // 4 = RGBA
	gpu.ResetFrame()

	gpu.obp0 = 0xff
	gpu.obp1 = 0xff

	gpu.stat = 0x80

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

// DisplayTileSets displays all tile sets. It's only for debug mode
func (gpu *GPU) DisplayTileSets() {
	for y := 0; y < screenHeight; y++ {
		for x := 0; x < screenWidth; x++ {
			tileNum := (y/8)*(screenWidth/8) + x/8
			colorNum := gpu.tileSets[tileNum][y%8][x%8]

			gpu.paintPixel(y*screenWidth+x, colorNum)
		}
	}
}

func (gpu *GPU) renderTiles() {
	var base uint16 = 0x1800

	y := (gpu.scy + gpu.ly) & 255
	var tileRow uint16 = uint16(y/8) * 32

	for lx := 0; lx < 160; lx++ {
		x := uint16(lx) + uint16(gpu.scx)
		tileCol := x / 8
		tileAddr := base + tileRow + tileCol
		tileNum := gpu.vram[tileAddr]

		// (y, x) is coordinate in 256 * 256 full background
		colorNum := gpu.tileSets[tileNum][y%8][x%8]

		// (ly, lx) is coordinate in 160 * 144 screen
		coord := int(gpu.ly)*screenWidth + lx

		gpu.paintPixel(coord, colorNum)
	}
}

func (gpu *GPU) paintPixel(coord int, colorNum uint8) {
	color := gpu.getColor(colorNum)

	red, green, blue := getRGB(color)

	gpu.Pixels[coord*4+0] = red   // R
	gpu.Pixels[coord*4+1] = green // G
	gpu.Pixels[coord*4+2] = blue  // B
	gpu.Pixels[coord*4+3] = 0xff  // A
}

func (gpu *GPU) getColor(colorNum uint8) uint8 {
	var color uint8

	switch colorNum {
	case 0:
		color = gpu.bgp & 0x3
	case 1:
		color = (gpu.bgp & 0xc) >> 2
	case 2:
		color = (gpu.bgp & 0x30) >> 4
	case 3:
		color = (gpu.bgp & 0xc0) >> 6
	}

	switch color {
	case 0:
		return white
	case 1:
		return ligthGray
	case 2:
		return darkGray
	case 3:
		return black
	}

	return white
}

func getRGB(color uint8) (uint8, uint8, uint8) {
	switch color {
	case white:
		return 0xff, 0xff, 0xff
	case ligthGray:
		return 0xcc, 0xcc, 0xcc
	case darkGray:
		return 0x77, 0x77, 0x77
	case black:
		return 0x00, 0x00, 0x00
	}

	return 0, 0, 0
}

func (gpu *GPU) Read(addr uint16) uint8 {
	if 0x8000 <= addr && addr <= 0x9fff {
		return gpu.vram[addr-0x8000]
	}

	switch addr {
	case 0xff40:
		return gpu.lcdc
	case 0xff41:
		return gpu.stat | gpu.mode
	case 0xff42:
		return gpu.scy
	case 0xff43:
		return gpu.scx
	case 0xff44:
		return gpu.ly
	case 0xff45:
		return gpu.lyc
	case 0xff46:
		return gpu.dma
	case 0xff47:
		return gpu.bgp
	case 0xff48:
		return gpu.obp0
	case 0xff49:
		return gpu.obp1
	case 0xff4a:
		return gpu.wy
	case 0xff4b:
		return gpu.wx
	}

	return gpu.vram[addr]
}

func (gpu *GPU) Write(addr uint16, val uint8) {

	if 0x8000 <= addr && addr <= 0x9fff {
		gpu.vram[addr-0x8000] = val
		gpu.updateTileSets()

		return
	}

	switch addr {
	case 0xff40:
		gpu.lcdc = val
	case 0xff41:
		// bit 2-0 are Read Only
		// bit 7 is always set
		gpu.stat = val&0xf8 | gpu.stat&0x07 | 1<<7
	case 0xff42:
		gpu.scy = val
	case 0xff43:
		gpu.scx = val
	case 0xff44:
		gpu.ly = 0 // ReadOnly. Writing will reset the counter
	case 0xff45:
		gpu.lyc = val
	case 0xff46:
		gpu.dma = val
	case 0xff47:
		gpu.bgp = val
	case 0xff48:
		gpu.obp0 = val
	case 0xff49:
		gpu.obp1 = val
	case 0xff4a:
		gpu.wy = val
	case 0xff4b:
		gpu.wx = val
	}
}

func (gpu *GPU) isLCDEnabled() bool {
	return gpu.lcdc&0x80 != 0
}

func (gpu *GPU) compareLYC() {
	// update stat bit-2 coincidence flag
	if gpu.ly == gpu.lyc {
		gpu.stat |= 1 << 2
	} else {
		gpu.stat &= ^(uint8(1 << 2))
	}
}

func (gpu *GPU) Update(ticks uint8) {
	if !gpu.isLCDEnabled() {
		return
	}

	gpu.counter += uint16(ticks)

	// the mode goes through 2 -> 3 -> 0 -> ...
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
			gpu.ly++

			if gpu.ly >= 143 {
				gpu.mode = 1
			} else {
				gpu.mode = 2
				gpu.renderTiles()
			}
		}

	// vertical blank
	case 1:
		if gpu.counter >= 456 {
			gpu.counter -= 456
			gpu.ly++

			if gpu.ly > 153 {
				gpu.mode = 2
				gpu.ly = 0
			}
		}

	}

	gpu.compareLYC()
	// fmt.Printf("ly: %d\n", gpu.ly)
	// fmt.Printf("scy: %d\n", gpu.scy)
	// fmt.Printf("GPU counter: %d\n", gpu.counter)
}
