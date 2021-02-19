// Package vint implements Variable-Size Integer manipulation
// functions for some of the predeclared unsigned integer types.
//
// https://tools.ietf.org/html/rfc8794#section-4
package vint

import (
	"math/big"
	"math/bits"
)

// A form value describes the internal representation.
type form byte

const (
	finite form = iota
	inf
)

type Vint struct {
	f form
	i *big.Int
	w int
	d *big.Int
}

func NewVint(b []byte) *Vint {
	v := &Vint{
		f: finite,
		i: new(big.Int).SetBytes(b),
		w: len(b),
	}
	marker := (v.w * 8) - v.w
	v.d = new(big.Int).SetBit(v.i, marker, 0)
	return v
}

func (v Vint) Val() *big.Int {
	return v.i
}

func (v Vint) Width() int {
	return v.w
}

func (v Vint) Data() *big.Int {
	return v.d
}

func AllOne(b []byte, w int) bool {
	var oc int
	for _, bb := range b {
		oc += bits.OnesCount8(bb)
	}
	return oc == (w*8 - w)
}

func ShorterAvailable(b []byte, w int) bool {
	tz := (w * 8) - (len(b) * 8)
	for _, bb := range b {
		x := bits.LeadingZeros8(bb)
		tz += x
		if x < 8 {
			break
		}
	}
	return (tz - w) >= 8
}
