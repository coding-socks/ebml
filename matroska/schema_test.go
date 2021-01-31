package matroska

import (
	"github.com/coding-socks/ebml"
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestDecode(t *testing.T) {
	f, err := os.Open(filepath.Join("..", "test5.mkv"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var d Document
	if err := ebml.NewDecoder(f).Decode(&d); err != nil {
		t.Fatal(err)
	}
	log.Printf("%+v", d.EBML)
}

func TestElement(t *testing.T) {
	f, err := os.Open(filepath.Join("..", "test5.mkv"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	d := ebml.NewDecoder(f)
	var el ebml.Element
	for err == nil {
		el, err = d.Element()
		if err == nil {
			log.Printf("id: %8x; size: %d", el.ID.Val(), el.DataSize.Data())
		}
	}
	if err != io.EOF {
		t.Fatal(err)
	}
}
