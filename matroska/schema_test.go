package matroska

import (
	"github.com/coding-socks/ebml"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func downloadTestFile(filename, source string) error {
	resp, err := http.Get(source)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	f, err := os.Create(filepath.Join(".", filename))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

func TestDecode(t *testing.T) {
	var err error
	var f *os.File
	testFiles := []struct {
		name     string
		filename string
		source   string
	}{
		{
			name:     "Live stream recording",
			filename: "test4.mkv",
			source:   "https://github.com/Matroska-Org/matroska-test-files/blob/master/test_files/test4.mkv?raw=true",
		},
		{
			name:     "Multiple audio/subtitles",
			filename: "test5.mkv",
			source:   "https://github.com/Matroska-Org/matroska-test-files/blob/master/test_files/test5.mkv?raw=true",
		},
		{
			name:     "Different EBML head sizes & cue-less seeking",
			filename: "test6.mkv",
			source:   "https://github.com/Matroska-Org/matroska-test-files/blob/master/test_files/test6.mkv?raw=true",
		},
	}
	for _, tt := range testFiles {
		p := filepath.Join(".", tt.filename)
		_, err := os.Stat(p)
		if err == nil {
			continue
		}
		if !os.IsNotExist(err) {
			t.Fatal(err)
		}
		if err := downloadTestFile(tt.filename, tt.source); err != nil {
			t.Fatal(err)
		}
	}
	for _, tt := range testFiles {
		t.Run(tt.name, func(t *testing.T) {
			f, err = os.Open(filepath.Join(".", tt.filename))
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			var d Document
			if err = ebml.NewDecoder(f).Decode(&d); err != nil {
				t.Error(err)
			}
			log.Printf("%+v %+v", d.EBML, d.Segment.Info)
		})
	}
}
