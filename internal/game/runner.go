package game

import (
	"context"
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
	writer io.WriteCloser
}

func NewRunner(m board, w io.WriteCloser) *runner {
	return &runner{
		board:  m,
		writer: w,
	}
}

func (r *runner) Run(ctx context.Context) {
	_, _ = fmt.Fprint(r.writer, "Match Started")
	defer r.Terminate()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := r.board.MoveToRandomNeighborhood(func(s string) {
				_, _ = r.writer.Write([]byte(s))
			})
			if err == nil {
				time.Sleep(time.Millisecond * 500) // @TODO: REMOVE IT
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
	_, _ = fmt.Fprint(r.writer, "Match Over")
	_, _ = r.writer.Write([]byte("DUMP MAP"))
	_ = r.writer.Close()
}
