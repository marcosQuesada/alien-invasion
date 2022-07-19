package game

import (
	"errors"
	"fmt"
	"io"
	"time"

	log "github.com/sirupsen/logrus"
)

type board interface {
	MoveToRandomNeighborhood(reporter func(string)) error
}

type runner struct {
	board  board
	writer io.Writer
	done   chan struct{}
}

func NewRunner(m board, w io.Writer, done chan struct{}) *runner {
	return &runner{
		board:  m,
		writer: w,
		done:   done,
	}
}

func (r *runner) Run() {
	_, _ = fmt.Fprint(r.writer, "Match Started")
	defer r.Terminate()

	for {
		select {
		case <-r.done:
			return
		default:
			err := r.board.MoveToRandomNeighborhood(func(s string) {
				_, _ = r.writer.Write([]byte(s))
			})
			if err == nil {
				time.Sleep(time.Millisecond * 50) // @TODO: Move to option
				continue
			}

			if err != nil && errors.Is(err, ErrMatchIsOver) {
				return
			}

			log.Errorf("unexpected error %v", err)
		}
	}
}

func (r *runner) Terminate() {
	_, _ = fmt.Fprint(r.writer, "Game Over")
	close(r.done)
}
