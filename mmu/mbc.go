package mmu

// currently only support MBC1
const (
	ROMONLY = iota
	MBC1
	MBC2

	// current bank mode
	romBankingMode
	ramBankingMode
)

func (mmu *MMU) enableRAMBank(val uint8) {
	if mmu.cartridgeType == ROMONLY {
		return
	}

	if val&0xf == 0xa {
		mmu.ramEnabled = true
	} else {
		mmu.ramEnabled = false
	}

}

func (mmu *MMU) changeLoROMBANK(val uint8) {
	if val&0x1f == 0 {
		mmu.currentROMBank = mmu.currentROMBank&0xe0 | 1
		return
	}

	mmu.currentROMBank = mmu.currentROMBank&0xe0 | val&0x1f
	// fmt.Println(mmu.currentROMBank)
}

func (mmu *MMU) changeHiROMBANK(val uint8) {
	mmu.currentROMBank = mmu.currentROMBank&0x9f | val&0x60
}

func (mmu *MMU) changeRAMBANK(val uint8) {
	mmu.currentRAMBank = val & 0x3
}

func (mmu *MMU) changeBankingMode(val uint8) {
	switch val & 1 {
	case 0:
		mmu.bankMode = romBankingMode
		mmu.currentRAMBank = 0
	case 1:
		mmu.bankMode = ramBankingMode
	}
}

func (mmu *MMU) handleMBC(addr uint16, val uint8) {
	switch {
	case addr <= 0x1fff:
		mmu.enableRAMBank(val)

	case addr <= 0x3fff:
		mmu.changeLoROMBANK(val)

	case addr <= 0x5fff:
		if mmu.bankMode == romBankingMode {
			mmu.changeHiROMBANK(val)
		} else if mmu.bankMode == ramBankingMode {
			mmu.changeRAMBANK(val)
		}

	case addr <= 0x7fff:
		mmu.changeBankingMode(val)
	}
}
