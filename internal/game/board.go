package game

import (
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
)

const maxAliensByCity = 2

type Road struct {
	Remote CityName
}

type CityName string

type City struct {
	name     CityName
	visitors map[AlienName]*Alien
}

type PlanetMap struct {
	cities map[CityName]*City
	roads  map[CityName][]*Road // Indexed by Direction as INT
}

func (p *PlanetMap) Dump(r io.Writer) {
	for name, _ := range p.cities {
		res := string(name)
		if p.roads[name][NorthDirection] != nil {
			res = fmt.Sprintf("%s %s=%s", res, North, p.roads[name][NorthDirection].Remote)
		}
		if p.roads[name][EastDirection] != nil {
			res = fmt.Sprintf("%s %s=%s", res, East, p.roads[name][EastDirection].Remote)
		}
		if p.roads[name][SouthDirection] != nil {
			res = fmt.Sprintf("%s %s=%s", res, South, p.roads[name][SouthDirection].Remote)
		}
		if p.roads[name][WestDirection] != nil {
			res = fmt.Sprintf("%s %s=%s", res, West, p.roads[name][WestDirection].Remote)
		}
		_, _ = r.Write([]byte(res))
	}
}

func newEmptyMap() *PlanetMap {
	return &PlanetMap{
		cities: make(map[CityName]*City),
		roads:  make(map[CityName][]*Road),
	}
}

func NewCity(name CityName) *City {
	return &City{
		name:     name,
		visitors: make(map[AlienName]*Alien),
	}
}

func (c *City) AddVisitor(a *Alien) error {
	log.Infof("%s from %s visits city %s", a.name, a.position, c.name)

	c.visitors[a.name] = a
	if len(c.visitors) < maxAliensByCity {
		return nil
	}

	log.Printf("War started on City %s 2 aliens", c.name)

	v := []AlienName{}
	for _, visitor := range c.visitors {
		v = append(v, visitor.name)
	}

	return &ErrCityWarStarted{
		city:   c.name,
		alienA: v[0],
		alienB: v[1],
	}
}

func (c *City) RemoveVisitor(a AlienName) {
	delete(c.visitors, a)
}
