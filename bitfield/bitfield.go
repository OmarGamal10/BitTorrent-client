package bitfield

import "errors"

// bitfield of indices 0 from left
type Bitfield []byte

func (bf *Bitfield) IsAvailable(idx int) bool {
	byteIdx := idx / 8
	bitIdx := idx % 8
	return (*bf)[byteIdx]>>(7-bitIdx)&1 != 0
}

func (bf *Bitfield) Set(idx int) error {
	byteIdx := idx / 8
	if byteIdx >= len(*bf) {
		return errors.New("Invald piece index")
	}
	bitIdx := idx % 8
	(*bf)[byteIdx] |= 1 << (7 - bitIdx)
	return nil
}
