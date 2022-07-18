package game

import (
	"errors"
	"testing"
)

func TestOnDestroyCityRemovesCityAndAllRelatedRoadsFromMap(t *testing.T) {
	filename := "test.csv"
	m, err := LoaDefinitionsFromFile(filename)
	if err != nil {
		t.Fatalf("unable to load file %s error %v", filename, err)
	}

	st := newStaticProvider()
	r := NewEngine(m, st)
	city := CityName("CITY-1-0")
	a1 := NewAlien("Alien-1", 10)
	a1.position = city
	r.players[a1.name] = a1

	a2 := NewAlien("Alien-2", 10)
	a2.position = city
	r.players[a2.name] = a2

	err = r.setCityOnWar(city, a1.name, a2.name)
	if err == nil {
		t.Fatal("expected error")
	}

	if !errors.Is(err, ErrMatchIsOver) {
		t.Fatalf("unexpected error type, gt %T", err)
	}
}

type fakeRandomizer struct {
	alternative bool
	carry       int
}

func newStaticProvider() *fakeRandomizer {
	return &fakeRandomizer{}
}

func newAltProvider() *fakeRandomizer {
	return &fakeRandomizer{alternative: true}
}

func (f *fakeRandomizer) RandomPosition(max int) int {
	f.carry++
	if f.carry%2 == 0 {
		return 0
	}
	return 1
}
