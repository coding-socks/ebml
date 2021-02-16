// +build codegen

// This program generates schema.go.

package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/coding-socks/ebml/internal/schema"
	"golang.org/x/tools/imports"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var header = []byte(`// Code generated by go run make_definition.go. DO NOT EDIT.

package ebml

type Definition struct {
	ID       []byte
	Type     string
	Name     string
	Default  interface{}
	Children []Definition
}

func NewDefinition(id string, t, name string, def interface{}, children []Definition) Definition {
	bid, err := hex.DecodeString(id)
	if err != nil {
		panic(err)
	}
	return Definition{ID: bid, Type: t, Name: name, Default: def, Children: children}
}

`)

func main() {
	filename := "definition.go"
	buf := bytes.NewBuffer(header)

	gen(buf)

	out, err := imports.Process(filename, buf.Bytes(), nil)
	if err != nil {
		err = ioutil.WriteFile(filename, buf.Bytes(), 0666)
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal(err)
	}

	err = ioutil.WriteFile(filename, out, 0666)
	if err != nil {
		log.Fatal(err)
	}
}

func gen(w io.Writer) {
	root := schema.NewTreeNode(schema.Element{
		Type: schema.TypeMaster,
		Name: "Document",
	})
	for _, fp := range []string{filepath.Join(".", "ebml.xml")} {
		var s schema.Schema
		func() {
			f, err := os.Open(fp)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			if err := xml.NewDecoder(f).Decode(&s); err != nil {
				log.Fatal(err)
			}
		}()
		for _, el := range s.Element {
			if el.Path == `\(-\)Void` || el.Path == `\(1-\)CRC-32` {
				// TODO: Implement support for Void and CRC-32 tags
				continue
			}
			p := strings.Split(el.Path, `\`)[1:]
			branch := root
			lastIndex := len(p) - 1
			for _, s := range p[:lastIndex] {
				node := branch.Get(s)
				if node == nil {
					node = schema.NewTreeNode(el)
					branch.Put(s, node)
				}
				branch = node
			}
			branch.Put(p[lastIndex], schema.NewTreeNode(el))
		}
	}
	fmt.Fprint(w, "var headerDefinition = ")
	root.VisitAll(func(node *schema.TreeNode) {
		writeDocType(w, node, true)
	})
}

func writeDocType(w io.Writer, node *schema.TreeNode, root bool) {
	id := strings.TrimPrefix(node.El.ID, "0x")
	def := node.El.Default
	if def == "" {
		def = "nil"
	} else if node.El.Type == schema.TypeString {
		def = fmt.Sprintf("%q", def)
	}
	fmt.Fprintf(
		w, "\n\tNewDefinition(%[2]q, %[3]s, \"%[1]s\", %[4]s, ",
		node.El.Name, id, schema.ResolveType(node.El.Type), def,
	)
	if node.El.Type == schema.TypeMaster {
		fmt.Fprint(w, "[]Definition{")
		node.VisitAll(func(n *schema.TreeNode) {
			writeDocType(w, n, false)
		})
		fmt.Fprint(w, "\n}")
	} else {
		fmt.Fprint(w, "nil")
	}
	if root {
		fmt.Fprint(w, ")")
	} else {
		fmt.Fprint(w, "),")
	}
}

func writeHeaderDefinition(w io.Writer, node *schema.TreeNode) {
	fmt.Fprintf(w, "\n\t%s,", node.El.Name)
	node.VisitAll(func(n *schema.TreeNode) {
		writeHeaderDefinition(w, n)
	})
}
