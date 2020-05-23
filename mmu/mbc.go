package mmu

// currently only support MBC1
const (
	ROMONLY = iota
	MBC1
	MBC2
	MBC3
	MBC5

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
	} else if val == 0 {
		mmu.ramEnabled = false
		mmu.rtcEnabled = false
	}
}

func (mmu *MMU) changeLoROMBank(val uint8) {
	if val&0x1f == 0 {
		mmu.currentROMBank = mmu.currentROMBank&0xe0 | 1
		return
	}

	mmu.currentROMBank = mmu.currentROMBank&0xe0 | val&0x1f
}

func (mmu *MMU) changeHiROMBank(val uint8) {
	mmu.currentROMBank = mmu.currentROMBank&0x9f | (val & 0x3 << 5)
}

func (mmu *MMU) changeROMBankMBC3(val uint8) {
	mmu.currentROMBank = val & 0x7f
	if mmu.currentROMBank == 0 {
		mmu.currentROMBank = 1
	}
}

func (mmu *MMU) changeLoROMBankMBC5(val uint8) {
	mmu.currentROMBank = val
}

func (mmu *MMU) changeHiROMBankMBC5(val uint8) {
	if val&1 > 0 {
		mmu.hiCurrentROMBank = 1
	} else {
		mmu.hiCurrentROMBank = 0
	}
}

func (mmu *MMU) changeRAMBANK(val uint8) {
	mmu.currentRAMBank = val & 0x3
}

func (mmu *MMU) changeRAMBANKMBC3(val uint8) {
	mmu.currentRAMBank = val & 0x7
}

func (mmu *MMU) changeRAMBANKMBC5(val uint8) {
	mmu.currentRAMBank = val & 0xf
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
		if mmu.cartridgeType == MBC1 {
			mmu.changeLoROMBank(val)
		} else if mmu.cartridgeType == MBC3 {
			mmu.changeROMBankMBC3(val)
		} else if mmu.cartridgeType == MBC5 {
			if addr <= 0x2fff {
				mmu.changeLoROMBankMBC5(val)
			} else if addr <= 0x3fff {
				mmu.changeHiROMBankMBC5(val)
			}
		}

	case addr <= 0x5fff:
		if mmu.cartridgeType == MBC1 {
			if mmu.bankMode == romBankingMode {
				mmu.changeHiROMBank(val)
			} else if mmu.bankMode == ramBankingMode {
				mmu.changeRAMBANK(val)
			}
		} else if mmu.cartridgeType == MBC3 {
			if 0x8 <= val && val <= 0xc {
				mmu.rtcEnabled = true
				mmu.rtc = val
				return
			} else if val <= 0x7 {
				mmu.rtcEnabled = false
				mmu.changeRAMBANKMBC3(val)
			}
		} else if mmu.cartridgeType == MBC5 {
			mmu.changeRAMBANKMBC5(val)
		}

	case addr <= 0x7fff:
		if mmu.cartridgeType == MBC1 {
			mmu.changeBankingMode(val)
		}
	}
}
