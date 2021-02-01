package matroska

import (
	"github.com/coding-socks/ebml"
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
