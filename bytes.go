package ebml

import "math/bits"

func allOneVint(b []byte, w int) bool {
	var oc int
	for _, bb := range b {
		oc += bits.OnesCount8(bb)
	}
	return oc == (w*8 - w)
}

func shorterAvailableVint(b []byte, w int) bool {
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
