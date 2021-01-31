// Package vint implements Variable-Size Integer manipulation
// functions for some of the predeclared unsigned integer types.
//
// https://tools.ietf.org/html/rfc8794#section-4
package vint

import "math/big"

// A form value describes the internal representation.
type form byte

const (
	finite form = iota
	inf
)

var Inf = &Vint{
	f: inf,
	i: big.NewInt(0),
	w: 0,
	d: big.NewInt(0),
}

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

func (v Vint) IsInf() bool {
	return v.f == inf
}

// Data32 returns the VINT_DATA portion of a Variable-Size uint32.
//
// https://tools.ietf.org/html/rfc8794#section.4.3
func Data32(x uint32, w int) uint32 {
	if w > 4 || w < 1 {
		panic("vint: invalid VINT_WIDTH for 32 bit vint")
	}
	pos := (w * 8) - w
	return x - (1 << pos)
}

// Data64 returns the VINT_DATA portion of a Variable-Size uint64.
//
// https://tools.ietf.org/html/rfc8794#section.4.3
func Data64(x uint64, w int) uint64 {
	if w > 8 || w < 1 {
		panic("vint: invalid VINT_WIDTH for 32 bit vint")
	}
	pos := (w * 8) - w
	return x - (1 << pos)
}
