package ebml

import (
	"math/bits"
)

func vintOctetLength(b []byte) int {
	var zc int
	n := len(b)
	for i := 0; i < n; i++ {
		x := bits.LeadingZeros8(b[i])
		zc += x
		if x < 8 {
			break
		}
	}
	return zc + 1
}

func vintData(b []byte, l int) []byte {
	data := make([]byte, l)
	i := (l - 1) / 8
	copy(data, b[i:l+i])
	j := 8 - (l % 8)
	data[0] = clearBit(data[0], j)
	return data
}

func clearBit(n byte, pos int) byte {
	n &^= 1 << (pos % 8)
	return n
}

func vintDataAllOne(b []byte, w int) bool {
	var oc int
	for i := 0; i < w; i++ {
		bb := b[i]
		oc += bits.OnesCount8(bb)
	}
	return oc == (w*8 - w)
}
