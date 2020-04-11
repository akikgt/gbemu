package gpu

import (
	"fmt"
	"gbemu/mmu"
)

type GPU struct {
	mmu *mmu.MMU

	counter uint16
}

func New(mmu *mmu.MMU) *GPU {
	gpu := &GPU{mmu: mmu}

	return gpu
}

func (gpu *GPU) Update(ticks uint8) {
	gpu.counter += uint16(ticks)

	gpu.counter %= 456

	fmt.Printf("GPU counter: %d\n", gpu.counter)
}
