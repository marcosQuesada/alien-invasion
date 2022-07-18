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
