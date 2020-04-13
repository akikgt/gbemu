package main

import (
	"bufio"
	"fmt"
	c "gbemu/cpu"
	g "gbemu/gpu"
	m "gbemu/mmu"
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
		cpu.Execute()
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
		return false
	}
}

const (
	screenWidth  = 160
	screenHeight = 144

	// GB CPU is 4194304Hz. To get 60FPS, 4194304/60
	maxTicks = 69905
)

var (
	gpu *g.GPU = g.New()
	mmu *m.MMU = m.New(gpu)
	cpu *c.CPU = c.New(mmu)
	// breakPoint uint16 = 0x100
	// breakPoint uint16 = 0x282a // tetris end of tileset loading
)

func update(screen *ebiten.Image) error {

	// reset TotalTicks every update
	cpu.TotalTicks = 0

	for cpu.TotalTicks < maxTicks {
		// fmt.Printf("%#04x : ", cpu.GetPC())
		// cpu.TestFlags()

		ticks := cpu.Execute()
		gpu.Update(ticks)
		// if cpu.GetPC() == breakPoint && !mmu.IsBooting {
		// 	cpu.Dump()
		// 	isContinue := debugMode(cpu, &breakPoint)
		// 	if !isContinue {
		// 		break
		// 	}
		// } else {
		// 	ticks := cpu.Execute()
		// 	gpu.Update(ticks)
		// }

	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// gpu.DisplayTileSets()
	screen.ReplacePixels(gpu.Pixels)

	msg := fmt.Sprintf("TPS = %0.2f\nFPS = %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS())
	ebitenutil.DebugPrint(screen, msg)

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

	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Game Boy Emulator"); err != nil {
		log.Fatal(err)
	}
}
