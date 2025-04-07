package ebmltext

import (
	"bytes"
	"github.com/coding-socks/ebml/schema"
	"io"
	"testing"
)

func TestCode(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	wantID := schema.ElementID(0x1a45dfa3)
	if _, err := enc.WriteElementID(wantID); err != nil {
		t.Fatal(err)
	}
	wantDataSize := int64(0xff)
	if _, err := enc.WriteElementDataSize(wantDataSize, 1); err != nil {
		t.Fatal(err)
	}
	if _, err := enc.Write(bytes.Repeat([]byte{'0'}, int(wantDataSize))); err != nil {
		t.Fatal(err)
	}
	if want := []byte{0x1a, 0x45, 0xdf, 0xa3, 0x40, 0xff}; !bytes.Equal(buf.Bytes()[:buf.Len()-int(wantDataSize)], want) {
		t.Errorf("want %x, got %x", want, buf.Bytes())
	}

	dec := NewDecoder(bytes.NewReader(buf.Bytes()))
	gotID, err := dec.ReadElementID()
	if err != nil {
		t.Fatal(err)
	}
	dec.Release()
	if gotID != wantID {
		t.Errorf("gotID = %d, want %d", gotID, wantID)
	}
	gotDataSize, err := dec.ReadElementDataSize()
	if err != nil {
		t.Fatal(err)
	}
	dec.Release()
	if gotDataSize != wantDataSize {
		t.Errorf("gotDataSize = %d, want %d", gotDataSize, wantDataSize)
	}
	if _, err := io.CopyN(io.Discard, dec, gotDataSize); err != nil {
		t.Fatal(err)
	}
}
