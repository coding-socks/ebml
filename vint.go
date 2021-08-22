// Package vint implements Variable-Size Integer manipulation
// functions for some of the predeclared unsigned integer types.
//
// https://tools.ietf.org/html/rfc8794#section-4
package ebml

import (
	"math/bits"
)

func vintOctetLength(vint []byte) int {
	var zc int
	n := len(vint)
	for i := 0; i < n; i++ {
		x := bits.LeadingZeros8(vint[0])
		zc += x
		if x < 8 {
			break
		}
	}
	return bits.LeadingZeros8(vint[0]) + 1
}

func vintData(vint []byte, l int) []byte {
	data := make([]byte, l)
	copy(data, vint[:l])
	data[0] = clearBit(data[0], 8-l)
	return data
}

func clearBit(n byte, pos int) byte {
	n &^= 1 << (pos % 8)
	return n
}

func vintDataAllZero(b []byte, w int) bool {
	var zc int
	for i := 0; i < w; i++ {
		bb := b[i]
		zc += 8 - bits.OnesCount8(bb)
	}
	return zc == (w * 8)
}

func vintDataAllOne(b []byte, w int) bool {
	var oc int
	for i := 0; i < w; i++ {
		bb := b[i]
		oc += bits.OnesCount8(bb)
	}
	return oc == (w*8 - w)
}
