package game

import (
	"fmt"
	"strings"
)

type DirectionType int

const (
	NorthDirection DirectionType = iota
	EastDirection
	SouthDirection
	WestDirection
)

const (
	North = "north"
	East  = "east"
	South = "south"
	West  = "west"
)

func DirectionTypeFromString(s string) (DirectionType, error) {
	s = strings.ToLower(s)
	switch s {
	case North:
		return NorthDirection, nil
	case East:
		return EastDirection, nil
	case South:
		return SouthDirection, nil
	case West:
		return WestDirection, nil
	}

	return 9, fmt.Errorf("unable to find direction from type %s", s)
}

func (d DirectionType) String() string {
	switch d {
	case NorthDirection:
		return North
	case EastDirection:
		return East
	case SouthDirection:
		return South
	case WestDirection:
		return West
	}

	return "" // Cannot happen directions are assumed already validated
}

func (d DirectionType) Opposite() DirectionType {
	switch d {
	case NorthDirection:
		return SouthDirection
	case EastDirection:
		return WestDirection
	case SouthDirection:
		return NorthDirection
	case WestDirection:
		return EastDirection
	}

	return 9 // Cannot happen directions are assumed already validated
}

type CityName string

type City struct {
	name     CityName
	visitors map[string]struct{}
}

func NewCity(name CityName) *City {
	return &City{
		name:     name,
		visitors: make(map[string]struct{}),
	}
}

func (c *City) Name() CityName {
	return c.name
}

type Road struct {
	Remote CityName
}

type PlanetMap struct {
	Cities map[CityName]*City
	Roads  map[CityName][]*Road // Indexed by Direction as INT

	//mutex *sync.Mutex
}

func NewEmptyMap() *PlanetMap {
	return &PlanetMap{
		Cities: make(map[CityName]*City),
		Roads:  make(map[CityName][]*Road),
		//mutex:  &sync.Mutex{},
	}
}
