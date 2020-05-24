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

	breakPoint uint16 = 0xffff
)

func update(screen *ebiten.Image) error {

	// reset TotalTicks every update
	cpu.TotalTicks = 0

	for cpu.TotalTicks < maxTicks {
		ticks := cpu.Execute()
		gpu.Update(ticks)
		timer.Update(ticks)
		cpu.HandleInterrupts()
		// if cpu.GetPC() == 0x0233 {
		// 	cpu.SetPC(0x0236)
		// }
		// fmt.Printf("%04x\n", cpu.GetPC())
	}

	// fmt.Printf("%04x\n", cpu.GetPC())
	// gpu.DisplayTileSets()
	gpu.DumpColorPalette()
	// fmt.Printf("%04x\n", cpu.GetPC())
	// mmu.PrintCurrentRomBank()
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	screen.ReplacePixels(gpu.Pixels)

	// for debug, TPS, FPS
	// msg := fmt.Sprintf("TPS = %0.2f\nFPS = %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS())
	// ebitenutil.DebugPrint(screen, msg)

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
		joypad.ReleaseAll()
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

	buf := make([]byte, 0x1000000)
	nb, err := fp.Read(buf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Successfully read %d byte\n", nb)
	// for i := 0; i < 4096; i++ {
	// 	fmt.Printf("%#04x\n", buf[uint32(i+0x4000)+uint32(58)<<14])
	// }
	// os.Exit(1)

	mmu.Load(buf)

	cpu.Reset()
	if len(os.Args) == 3 && os.Args[2] == "--color" {
		cpu.SetCGBMode()
		gpu.SetCGBMode()
	}

	if err := ebiten.Run(update, screenWidth, screenHeight, 3, "Game Boy Emulator"); err != nil {
		log.Fatal(err)
	}
}
