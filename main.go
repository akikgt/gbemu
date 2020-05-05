package main

import (
	"bufio"
	"fmt"
	c "gbemu/cpu"
	g "gbemu/gpu"
	j "gbemu/joypad"
	m "gbemu/mmu"
	t "gbemu/timer"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func debugMode(cpu *c.CPU, breakPoint *uint16) bool {
	fmt.Printf("_")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	switch input {
	case "d":
		cpu.Dump()
		return debugMode(cpu, breakPoint)
	case "i":
		cpu.PrintNextIns()
		return debugMode(cpu, breakPoint)
	case "n":
		ticks := cpu.Execute()
		gpu.Update(ticks)
		timer.Update(ticks)
		cpu.HandleInterrupts()
		*breakPoint = cpu.GetPC()
		return true
	case "c":
		// loop until end
		// end of BIOS = 0x100 or end of GB RAM = 0x10000 65536
		*breakPoint = 0x100
		return true
	case "q":
		// quit
		return false
	default:
		ticks := cpu.Execute()
		gpu.Update(ticks)
		timer.Update(ticks)
		cpu.HandleInterrupts()
		*breakPoint = cpu.GetPC()
		return true
	}
}

const (
	screenWidth  = 160
	screenHeight = 144

	// GB CPU is 4194304Hz. To get 60FPS, 4194304/60
	maxTicks = 69905
)

var (
	gpu    *g.GPU    = g.New()
	timer  *t.Timer  = t.New()
	joypad *j.Joypad = j.New()
	mmu    *m.MMU    = m.New(gpu, timer, joypad)
	cpu    *c.CPU    = c.New(mmu)

	breakPoint uint16 = 0xc370
	// breakPoint uint16 = 0x29fa
	// 0x2a24 でff80の値が実機と違う
	// after 0x034c tetris load all tiles
	// breakPoint uint16 = 0x282a // tetris end of tileset loading
)

func update(screen *ebiten.Image) error {

	// reset TotalTicks every update
	cpu.TotalTicks = 0

	for cpu.TotalTicks < maxTicks {

		ticks := cpu.Execute()
		gpu.Update(ticks)
		timer.Update(ticks)
		cpu.HandleInterrupts()
		continue

		if cpu.GetPC() == breakPoint && !mmu.IsBooting {
			cpu.Dump()
			isContinue := debugMode(cpu, &breakPoint)
			if !isContinue {
				break
			}
		} else {
			ticks := cpu.Execute()
			gpu.Update(ticks)
			timer.Update(ticks)
			cpu.HandleInterrupts()
		}

	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// gpu.DisplayTileSets()
	screen.ReplacePixels(gpu.Pixels)

	msg := fmt.Sprintf("TPS = %0.2f\nFPS = %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS())
	ebitenutil.DebugPrint(screen, msg)

	// joypad
	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		joypad.KeyPress(j.DOWN)
	} else if ebiten.IsKeyPressed(ebiten.KeyK) {
		joypad.KeyPress(j.UP)
	} else if ebiten.IsKeyPressed(ebiten.KeyH) {
		joypad.KeyPress(j.LEFT)
	} else if ebiten.IsKeyPressed(ebiten.KeyL) {
		joypad.KeyPress(j.RIGHT)
	} else if ebiten.IsKeyPressed(ebiten.KeyF) {
		joypad.KeyPress(j.START)
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		joypad.KeyPress(j.SELECT)
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		joypad.KeyPress(j.B)
	} else if ebiten.IsKeyPressed(ebiten.KeyA) {
		joypad.KeyPress(j.A)
	} else {
		joypad.KeyRelease(j.DOWN)
		joypad.KeyRelease(j.UP)
		joypad.KeyRelease(j.LEFT)
		joypad.KeyRelease(j.RIGHT)
		joypad.KeyRelease(j.START)
		joypad.KeyRelease(j.SELECT)
		joypad.KeyRelease(j.B)
		joypad.KeyRelease(j.A)
		joypad.ReqJoypadInt = false
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Too few arguments. Please provide GameBoy ROM")
		os.Exit(1)
	}

	fp, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	buf := make([]byte, 0x10000)
	nb, err := fp.Read(buf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully read %d byte\n", nb)

	mmu.Load(buf)

	cpu.Reset()

	if err := ebiten.Run(update, screenWidth, screenHeight, 3, "Game Boy Emulator"); err != nil {
		log.Fatal(err)
	}
}
