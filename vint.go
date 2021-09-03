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
	copy(data, b[:l])
	i := (l - 1) / 8
	j := 8 - (l % 8)
	data[i] = clearBit(data[i], j)
	return data
}

func clearBit(n byte, pos int) byte {
	n &^= 1 << (pos % 8)
	return n
}

func vintDataAllOne(b []byte, l int) bool {
	var oc int
	for i := 0; i < l; i++ {
		oc += bits.OnesCount8(b[i])
	}
	return oc == (l*8 - l)
}
