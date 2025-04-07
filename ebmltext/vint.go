package ebmltext

import (
	"io"
	"math/bits"
)

// AppendVint appends the binary representation of a VINT to b based on
// https://datatracker.ietf.org/doc/html/rfc8794.html#section-4
func AppendVint(num uint64, b []byte) (w int, err error) {
	w = ((bits.Len64(num) - 1) / 8) + 1
	if w > len(b) {
		return 0, io.ErrShortBuffer
	}
	offset := (w - 1) / 8
	for i := w - 1; i >= offset; i-- {
		// big endian
		b[i] = byte(num >> ((w - i - 1) * 8))
	}
	return w, nil
}

// AppendVintData appends the binary representation of a VINT to b based on
// https://datatracker.ietf.org/doc/html/rfc8794.html#section-4
func AppendVintData(num uint64, minW int, b []byte) (w int, err error) {
	w = ((bits.Len64(num) - 1) / 7) + 1
	if minW > w {
		w = minW
	}
	if w > len(b) {
		return 0, io.ErrShortBuffer
	}
	num |= 1 << ((8 * w) - w)
	return AppendVint(num, b)
}

// ReadVint reads the integer representation of a VINT from b based on
// https://datatracker.ietf.org/doc/html/rfc8794.html#section-4
func ReadVint(b []byte) (vint uint64, w int, err error) {
	w = 1
	n := len(b)
	for i := 0; i < n; i++ {
		x := bits.LeadingZeros8(b[i])
		w += x
		if x < 8 {
			break
		}
	}
	if w > n {
		return 0, 0, io.ErrShortBuffer
	}
	offset := w / 8
	for i := w - 1; i >= offset; i-- {
		// big endian
		vint |= uint64(b[i]) << ((w - i - 1) * 8)
	}
	return vint, w, nil
}

// ReadVintData reads the integer representation of a VINT from b based on
// https://datatracker.ietf.org/doc/html/rfc8794.html#section-4
func ReadVintData(b []byte) (vint uint64, w int, err error) {
	vint, w, err = ReadVint(b)
	if err != nil {
		return 0, 0, err
	}
	vint &^= 1 << ((w * 8) - w)
	return vint, w, nil
}

func vintAllOne(vint uint64) bool {
	w := ((bits.Len64(vint) - 1) / 7) + 1
	return vintDataAllOne(vint, w)
}

func vintDataAllOne(vint uint64, w int) bool {
	return bits.OnesCount64(vint<<((8*w)-w)) == ((8 * w) - w)
}

func makeVintDataAllOne(w int) uint64 {
	return (1 << ((w * 8) - w + 1)) - 1
}
