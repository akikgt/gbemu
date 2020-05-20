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

	vram0 [0x2000]uint8
	vram1 [0x2000]uint8 // for CGB mode
	oam   [0xa0]uint8

	lcdc uint8 // 0xff40 LCD control
	stat uint8 // 0xff41 LCDC status
	scy  uint8 // 0xff42
	scx  uint8 // 0xff43
	ly   uint8 // 0xff44 current Y-coordinate
	lyc  uint8 // 0xff45
	bgp  uint8 // 0xff47 bg palette data
	obp0 uint8 // 0xff48
	obp1 uint8 // 0xff49
	wy   uint8 // 0xff4a
	wx   uint8 // 0xff4b
	vbk  uint8 // 0xff4f VRAM bank

	Pixels    []byte
	tileSets  [384][8][8]uint8
	tileSets2 [384][8][8]uint8

	ReqVBlankInt bool
	ReqLCDInt    bool

	cgbMode bool
	cbgp    [0x40]uint8
	cbpidx  uint8
	// cobp [0x80]uint8
}

func New() *GPU {
	gpu := &GPU{}

	gpu.Pixels = make([]byte, screenHeight*screenWidth*4) // 4 = RGBA
	gpu.ResetFrame()

	gpu.obp0 = 0xff
	gpu.obp1 = 0xff

	gpu.stat = 0x85

	gpu.ReqVBlankInt = false
	gpu.ReqLCDInt = false

	gpu.cbpidx = 0
	return gpu
}

func (gpu *GPU) SetCGBMode() {
	gpu.cgbMode = true
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
			data1 := gpu.vram0[i*16+y*2]
			data2 := gpu.vram0[i*16+y*2+1]

			for x := 0; x < 8; x++ {
				b := 7 - x
				color := (data2>>b&1)<<1 | (data1 >> b & 1)
				gpu.tileSets[i][y][x] = color
			}
		}

		if !gpu.cgbMode {
			continue
		}

		// update tileSets for CGB mode
		for y := 0; y < 8; y++ {
			// each tile data is 16 byte
			data1 := gpu.vram1[i*16+y*2]
			data2 := gpu.vram1[i*16+y*2+1]

			for x := 0; x < 8; x++ {
				b := 7 - x
				color := (data2>>b&1)<<1 | (data1 >> b & 1)
				gpu.tileSets2[i][y][x] = color
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

			gpu.paintPixel(y*screenWidth+x, colorNum, gpu.bgp)
		}
	}
}

func (gpu *GPU) isWindowEnabled() bool {
	return gpu.lcdc&0x20 != 0 && gpu.wy <= gpu.ly
}

func (gpu *GPU) renderScanline() {
	if gpu.lcdc&0x1 > 0 {
		gpu.renderBG()
	}

	if gpu.lcdc&0x2 > 0 {
		gpu.renderSprites()
	}
}

func (gpu *GPU) renderSprites() {
	for i := 0; i < 40; i++ {
		y := gpu.oam[i*4] - 16
		x := gpu.oam[i*4+1] - 8
		tileNum := gpu.oam[i*4+2]
		attributes := gpu.oam[i*4+3]

		var height uint8 = 8
		if gpu.lcdc&0x4 > 0 {
			// use 8x16 mode
			height = 16
			// modify tileNum
			//
			// In 8x16 mode, the lower bit of the tile number is ignored.
			// Ie. the upper 8x8 tile is "NN AND FEh", and the lower 8x8 tile is "NN OR 01h".
			// reference: pandoc. (https://bgb.bircd.org/pandocs.htm#powerupsequence)
			tileNum &= 0xfe
		}

		// check current line includes sprite
		if !(y <= gpu.ly && gpu.ly < y+height) {
			continue
		}

		tileY := gpu.ly - y
		// Y flip
		if attributes>>6&1 == 1 {
			tileY = height - tileY - 1
		}

		// 8x16mode: check the tile is upper 8x8 tile or lower?
		// if lower(tileY == 8 - 15), add 1 to tileNum
		if tileY > 7 {
			tileNum |= 1
		}

		for lx := 0; lx < 8; lx++ {
			// X flip
			tileX := lx
			if attributes>>5&1 == 1 {
				tileX = 7 - tileX
			}

			colorNum := gpu.tileSets[tileNum][tileY%8][tileX%8]

			// for sprites, colorNum 0 means transparent
			if colorNum == 0 {
				continue
			}

			coord := int(gpu.ly)*screenWidth + (int(x) + lx)

			if (int(x) + lx) > 160 {
				continue
			}

			// priority
			if attributes>>7&1 == 1 {
				if !gpu.isSpritePrior(coord) {
					continue
				}
			}

			// change palette based on the attribute bit4
			palette := gpu.obp0
			if attributes>>4&1 == 1 {
				palette = gpu.obp1
			}

			gpu.paintPixel(coord, colorNum, palette)
		}
	}
}

func (gpu *GPU) isSpritePrior(coord int) bool {
	// check background color and if the back ground color is color num 0
	// return true

	// get background colors
	red := gpu.Pixels[coord*4+0]
	green := gpu.Pixels[coord*4+1]
	blue := gpu.Pixels[coord*4+2]

	// get base color of background palette
	baseColor := gpu.getNGBColor(0, gpu.bgp)
	baseRed, baseGreen, baseBlue := getMonochrome(baseColor)

	// check current background color is color num 0(base color) or not
	if red == baseRed && green == baseGreen && blue == baseBlue {
		return true
	}

	return false
}

func (gpu *GPU) renderBG() {
	var base uint16
	if gpu.isWindowEnabled() {
		if gpu.lcdc&0x40 != 0 {
			base = 0x1c00
		} else {
			base = 0x1800
		}
	} else {
		if gpu.lcdc&0x08 != 0 {
			base = 0x1c00
		} else {
			base = 0x1800
		}
	}

	var y uint16
	if gpu.isWindowEnabled() {
		y = uint16(gpu.ly - gpu.wy)
	} else {
		y = uint16((gpu.scy + gpu.ly) & 255)
	}

	tileRow := y / 8

	for lx := 0; lx < 160; lx++ {

		var x uint16
		if gpu.isWindowEnabled() && lx >= int(gpu.wx-7) {
			x = uint16(lx - int(gpu.wx-7))
		} else {
			x = (uint16(lx) + uint16(gpu.scx)) & 255
		}

		tileCol := x / 8
		tileAddr := base + tileRow*32 + tileCol
		var tileNum uint16 = uint16(gpu.vram0[tileAddr])

		// read BG map attributes
		attributes := gpu.vram1[tileAddr]
		paletteNum := attributes & 0x7
		tileBankNum := attributes & 0x8

		// select tile data 0=8800-97FF or 1=8000-8FFF
		if gpu.lcdc&0x10 == 0 && tileNum < 128 {
			tileNum += 256
		}

		// (y, x) is coordinate in 256 * 256 full background
		colorNum := gpu.tileSets[tileNum][y%8][x%8]
		if tileBankNum == 1 {
			colorNum = gpu.tileSets2[tileNum][y%8][x%8]
		}

		// (ly, lx) is coordinate in 160 * 144 screen
		coord := int(gpu.ly)*screenWidth + lx

		if gpu.cgbMode {
			gpu.paintColorPixel(coord, colorNum, paletteNum)
		} else {
			gpu.paintPixel(coord, colorNum, gpu.bgp)
		}
	}
}

func (gpu *GPU) paintColorPixel(coord int, colorNum uint8, palette uint8) {
	red, green, blue := gpu.getRGB(colorNum, palette)

	gpu.Pixels[coord*4+0] = red   // R
	gpu.Pixels[coord*4+1] = green // G
	gpu.Pixels[coord*4+2] = blue  // B
	gpu.Pixels[coord*4+3] = 0xff  // A
}

func (gpu *GPU) getRGB(colorNum, palette uint8) (uint8, uint8, uint8) {
	var color uint16 = uint16(gpu.cbgp[palette*8+2*colorNum]) | uint16(gpu.cbgp[palette*8+2*colorNum+1])<<8
	red := uint8(color&0x1f) << 3
	green := uint8(color>>5&0x1f) << 3
	blue := uint8(color>>10&0x1f) << 3
	return red, green, blue
}

func (gpu *GPU) paintPixel(coord int, colorNum uint8, palette uint8) {
	color := gpu.getNGBColor(colorNum, palette)

	red, green, blue := getMonochrome(color)

	gpu.Pixels[coord*4+0] = red   // R
	gpu.Pixels[coord*4+1] = green // G
	gpu.Pixels[coord*4+2] = blue  // B
	gpu.Pixels[coord*4+3] = 0xff  // A
}

func (gpu *GPU) getNGBColor(colorNum, palette uint8) uint8 {
	var color uint8

	switch colorNum {
	case 0:
		color = palette & 0x3
	case 1:
		color = (palette & 0xc) >> 2
	case 2:
		color = (palette & 0x30) >> 4
	case 3:
		color = (palette & 0xc0) >> 6
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

func getMonochrome(color uint8) (uint8, uint8, uint8) {
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
		if gpu.vbk == 1 {
			return gpu.vram1[addr-0x8000]
		}
		return gpu.vram0[addr-0x8000]
	}

	if 0xfe00 <= addr && addr <= 0xfe9f {
		return gpu.oam[addr-0xfe00]
	}

	switch addr {
	case 0xff40:
		return gpu.lcdc
	case 0xff41:
		return gpu.stat
	case 0xff42:
		return gpu.scy
	case 0xff43:
		return gpu.scx
	case 0xff44:
		return gpu.ly
	case 0xff45:
		return gpu.lyc
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

	return gpu.vram0[addr]
}

func (gpu *GPU) Write(addr uint16, val uint8) {

	if 0x8000 <= addr && addr <= 0x9fff {
		if gpu.vbk == 1 {
			gpu.vram1[addr-0x8000] = val
		} else {
			gpu.vram0[addr-0x8000] = val
		}
		gpu.updateTileSets()

		return
	}

	if 0xfe00 <= addr && addr <= 0xfe9f {
		gpu.oam[addr-0xfe00] = val

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
	case 0xff4f:
		gpu.vbk = val & 1

	// Background palette data
	case 0xff68:
		gpu.cbpidx = val
	case 0xff69:
		idx := gpu.cbpidx & 0x3f
		gpu.cbgp[idx] = val
		if gpu.cbpidx&0x80 > 0 {
			// Auto Increment
			gpu.cbpidx++
		}
	}
}

func (gpu *GPU) isLCDEnabled() bool {
	return gpu.lcdc&0x80 > 0
}

func (gpu *GPU) compareLYC() {
	// update stat bit-2 coincidence flag
	if gpu.ly == gpu.lyc {
		gpu.stat |= 1 << 2

		if gpu.stat>>6&1 == 1 {
			gpu.ReqLCDInt = true
		}

	} else {
		gpu.stat &= ^(uint8(1 << 2))
	}
}

func (gpu *GPU) updateLCDInterrupt() {
	mode := gpu.stat & 0x3

	switch mode {
	case 0: // H-Blank interrupt
		if gpu.stat>>3&1 == 1 {
			gpu.ReqLCDInt = true
		}
	case 1: // V-Blank interrupt
		if gpu.stat>>4&1 == 1 {
			gpu.ReqLCDInt = true
		}
	case 2: // OAM interrupt
		if gpu.stat>>5&1 == 1 {
			gpu.ReqLCDInt = true
		}
	}
}

func (gpu *GPU) Update(ticks uint8) {
	gpu.ReqLCDInt = false
	gpu.ReqVBlankInt = false

	if !gpu.isLCDEnabled() {
		// reference: https://www.reddit.com/r/Gameboy/comments/a1c8h0/what_happens_when_a_gameboy_screen_is_disabled/
		gpu.counter = 0
		gpu.ly = 0
		gpu.stat = gpu.stat & 0xf8 // enter mode 0.
		// http://www.codeslinger.co.uk/pages/projects/gameboy/lcd.html
		// says the mode should be 1. but I found Dr.mario won't past the menu if I set it to 1
		return
	}

	gpu.counter += uint16(ticks)

	// the mode goes through 2 -> 3 -> 0 -> ...
	switch gpu.stat & 0x03 {

	// accessing OAM
	case 2:
		if gpu.counter >= 80 {
			gpu.counter -= 80
			gpu.stat = gpu.stat&0xf8 | 3
		}

	// accessing VRAM
	case 3:
		if gpu.counter >= 172 {
			gpu.counter -= 172

			gpu.stat = gpu.stat & 0xf8
			gpu.updateLCDInterrupt()

			gpu.renderScanline()
		}

	// horizontal blank
	case 0:
		if gpu.counter >= 204 {
			gpu.counter -= 204
			gpu.ly++

			if gpu.ly >= 144 {
				// enter v-blank mode
				gpu.stat = gpu.stat&0xf8 | 1
				gpu.ReqVBlankInt = true
			} else {
				gpu.stat = gpu.stat&0xf8 | 2
			}

			gpu.updateLCDInterrupt()
		}

	// vertical blank
	case 1:
		if gpu.counter >= 456 {
			gpu.counter -= 456
			gpu.ly++

			if gpu.ly >= 154 {
				gpu.stat = gpu.stat&0xf8 | 2
				gpu.ly = 0

				gpu.updateLCDInterrupt()
			}
		}

	}

	gpu.compareLYC()
}
