package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/falconandy/sqlboiler-addons/editor"
)

func main() {
	var buf bytes.Buffer
	for _, arg := range os.Args[1:] {
		buf.Reset()
		b := editor.NewBuilder(arg)
		err := b.Build(&buf)
		if err != nil {
			log.Fatalln(err)
		}

		err = ioutil.WriteFile(filepath.Join(arg, "boil_editors.go"), buf.Bytes(), 0644)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
