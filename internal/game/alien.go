package game

import (
	"errors"
	"fmt"
	"sync/atomic"
)

var ErrMaxIterationsExceed = errors.New("max iterations exceed")

type ErrCityWarStarted struct {
	city   CityName
	alienA AlienName
	alienB AlienName
}

func (w *ErrCityWarStarted) Error() string {
	return fmt.Sprintf("%s has been destroyed by alien %s and alien %s!", w.city, w.alienA, w.alienB)
}

type AlienName string

type Alien struct {
	name          AlienName
	maxIterations int
	position      CityName
	iterations    int32
}

func NewAlien(name string, maxIterations int) *Alien {
	return &Alien{
		name:          AlienName(name),
		maxIterations: maxIterations,
	}
}

func (a *Alien) Iterate() error {
	if atomic.AddInt32(&a.iterations, 1) > a.iterations {
		return fmt.Errorf("alien %s %w", a.name, ErrMaxIterationsExceed)
	}

	return nil
}
