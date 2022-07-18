package game

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var ErrInvalidLineFormat = errors.New("invalid line format")

func LoaDefinitionsFromFile(filename string) (*PlanetMap, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open file %s error %w", filename, err)

	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	n := newEmptyMap()
	for scanner.Scan() {
		line := scanner.Text()
		err := parseLine(line, n)
		if err != nil {
			return nil, fmt.Errorf("raw line %s error %w", line, err)
		}
	}

	return n, nil
}

func parseLine(raw string, m *PlanetMap) error {
	line := strings.ReplaceAll(raw, "\n", "")
	parts := strings.Split(line, " ")
	if len(parts) < 2 {
		return ErrInvalidLineFormat
	}

	origin := CityName(parts[0])
	edges, err := parseEdges(parts[1:])
	if err != nil {
		return fmt.Errorf("unable to parse row, error %w", err)
	}

	m.cities[origin] = NewCity(origin)
	for directionType, dest := range edges {
		if _, ok := m.cities[dest]; !ok {
			m.cities[dest] = NewCity(dest)
		}

		if _, ok := m.roads[origin]; !ok {
			m.roads[origin] = make([]*Road, 4)
		}
		if _, ok := m.roads[dest]; !ok {
			m.roads[dest] = make([]*Road, 4)
		}

		m.roads[origin][directionType] = &Road{Remote: dest}
		m.roads[dest][directionType.Opposite()] = &Road{Remote: origin}
	}

	return nil
}

func parseEdges(raw []string) (map[DirectionType]CityName, error) {
	res := map[DirectionType]CityName{}
	for _, s := range raw {
		ep := strings.Split(s, "=")
		if len(ep) != 2 {
			return res, ErrInvalidLineFormat
		}

		d, err := DirectionTypeFromString(ep[0])
		if err != nil {
			return res, ErrInvalidLineFormat
		}

		res[d] = CityName(ep[1])
	}

	return res, nil
}
