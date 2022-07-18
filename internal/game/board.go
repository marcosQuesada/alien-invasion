package game

import log "github.com/sirupsen/logrus"

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

func (p *PlanetMap) Dump() string {
	return "DUMPING PLANET MAP" // @TODO: Fullfill
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

func (c *City) Name() CityName { // @TODO CLEAN
	return c.name
}

func (c *City) AddVisitor(a *Alien) error {
	log.Infof("Add visitor on city %s alienName %s", c.name, a.name)

	c.visitors[a.name] = a
	if len(c.visitors) < maxAliensByCity {
		return nil
	}

	log.Printf("War started on City %s 2 aliens", c.name)

	v := []AlienName{}
	for _, visitor := range c.visitors { // @TODO: Max visitors 2...
		v = append(v, visitor.name)
	}

	return &ErrCityWarStarted{
		city:   c.name,
		alienA: v[0],
		alienB: v[1],
	}
}
