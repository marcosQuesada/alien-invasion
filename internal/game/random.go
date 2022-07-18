package game

import (
	"math/rand"
	"time"
)

type randomProvider struct{}

func NewRandomProvider() Randomizer {
	return &randomProvider{}
}

func (r *randomProvider) RandomPosition(max int) int {
	s := rand.NewSource(time.Now().UnixNano())
	return rand.New(s).Intn(max)
}
