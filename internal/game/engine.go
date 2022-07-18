package game

import (
	"fmt"
	"io"
)

type engine struct {
	writer io.WriteCloser
}

func NewGame(w io.WriteCloser) *engine {
	return &engine{
		writer: w,
	}
}

func (g *engine) Run() {
	defer g.writer.Close()
	_, _ = g.writer.Write([]byte("XXXXXX"))
	_, _ = fmt.Fprintf(g.writer, "Heeello")
	_, _ = g.writer.Write([]byte("Planet Destroyed"))
}
