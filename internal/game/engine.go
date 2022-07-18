package game

import (
	"errors"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
)

var ErrNoAvailableRoads = errors.New("Error No available roads")
var ErrMatchIsOver = errors.New("Match is Over")

type Randomizer interface {
	RandomPosition(max int) int
}

type engine struct {
	planetMap *PlanetMap
	players   map[AlienName]*Alien
	random    Randomizer
	mutex     *sync.Mutex // @TODO: Wait until final concurrency analysis scenario
}

func NewEngine(st *PlanetMap, r Randomizer) *engine {
	return &engine{
		planetMap: st,
		random:    r,
		players:   make(map[AlienName]*Alien),
		mutex:     &sync.Mutex{},
	}
}

// @TODO
func (m *engine) Populate(totalPlayers int) {

}

func (m *engine) AssignRandomPosition(a *Alien) (exit bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.players[a.name] = a
	var err error
	a.position, err = m.assignMapRandomPosition(a)

	return errors.Is(err, ErrMatchIsOver)
}

func (m *engine) assignMapRandomPosition(a *Alien) (CityName, error) {
	var cities []*City // @TODO: Iterative on First Round, SHOULD WE CATCH IT?
	for _, city := range m.planetMap.cities {
		cities = append(cities, city)
	}
	log.Printf("Assign Random Position to alien %s", a.name)

	city := cities[m.random.RandomPosition(len(cities))]
	err := city.AddVisitor(a)
	if we, ok := err.(*ErrCityWarStarted); ok {

		if err := m.setCityOnWar(we.city, we.alienA, we.alienB); err != nil {
			return city.name, ErrMatchIsOver
		}

		return city.name, nil
	}

	return city.name, nil
}

// MoveToRandomNeighborhood Iterates on all players
func (m *engine) MoveToRandomNeighborhood(reporter func(string)) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.players) == 0 {
		return ErrMatchIsOver
	}

	for _, a := range m.players {
		cm, err := m.moveToRandomNeighborhood(a)
		// Alien on deadlock city
		if err != nil && errors.Is(err, ErrNoAvailableRoads) {
			reporter(err.Error())
			delete(m.players, a.name)
			continue
		}
		a.position = cm

		if we, ok := err.(*ErrCityWarStarted); ok {
			reporter(err.Error())
			if err := m.setCityOnWar(we.city, we.alienA, we.alienB); err != nil {
				return ErrMatchIsOver
			}
		}
	}

	return nil
}

func (m *engine) moveToRandomNeighborhood(a *Alien) (CityName, error) {
	if err := a.Iterate(); err != nil {
		log.Errorf("Alien %s Exits Game, max iterations achieved", a.name)
		delete(m.players, a.name)
		return "", nil
	}

	roads, ok := m.planetMap.roads[a.position]
	if !ok || len(roads) == 0 {
		// @TODO: Isolated alien
		return "", fmt.Errorf("alien %s on position %s no Roads available ", a.name, a.position)
	}

	var availableRoads []*Road // @TODO: THIS ¿?¿?¿?
	for _, r := range roads {
		if r == nil {
			continue
		}

		availableRoads = append(availableRoads, r)
	}

	if len(availableRoads) == 0 {
		return "", fmt.Errorf("alien %s on city %s error %w", a.name, a.position, ErrNoAvailableRoads)
	}

	road := availableRoads[m.random.RandomPosition(len(availableRoads))]
	destinationCity, ok := m.planetMap.cities[road.Remote]
	if !ok {
		return "", fmt.Errorf("no destination city %s found in the map", road.Remote)
	}

	err := destinationCity.AddVisitor(a)
	if err == nil {
		city, ok := m.planetMap.cities[a.position]
		if ok {
			city.RemoveVisitor(a.name)
		}
		return destinationCity.name, nil
	}

	return destinationCity.name, err

}

func (m *engine) SetCityOnWar(c CityName, a1, a2 AlienName) error {
	m.mutex.Lock()
	m.mutex.Unlock()

	return m.setCityOnWar(c, a1, a2)
}

func (m *engine) setCityOnWar(c CityName, a1, a2 AlienName) error {
	m.destroyCity(c)
	delete(m.players, a1)
	delete(m.players, a2)
	if len(m.players) == 0 {
		return ErrMatchIsOver
	}
	return nil
}

func (m *engine) destroyCity(c CityName) {
	if _, ok := m.planetMap.cities[c]; !ok {
		return
	}

	for direction, road := range m.planetMap.roads[c] {
		if road == nil {
			continue
		}
		// Remove opposite roads !
		dir := DirectionType(direction)
		m.planetMap.roads[road.Remote][dir.Opposite()] = nil
	}
	delete(m.planetMap.roads, c)
	delete(m.planetMap.cities, c)
}
