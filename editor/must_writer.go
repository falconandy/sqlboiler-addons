package editor

import (
	"io"
)

type mustWriter struct {
	w io.Writer
}

func (w mustWriter) Write(s string) {
	_, err := w.w.Write([]byte(s))
	if err != nil {
		panic(err)
	}
}

func (w mustWriter) Writeln(s string) {
	w.Write(s)
	w.Write("\n")
}
