package game

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var ErrInvalidLineFormat = errors.New("invalid line format")

func LoaDefinitions(filename string) (*PlanetMap, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open file %s error %w", filename, err)

	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	n := NewEmptyMap()
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

	m.Cities[origin] = NewCity(origin)
	for directionType, dest := range edges {
		if _, ok := m.Cities[dest]; !ok {
			m.Cities[dest] = NewCity(dest)
		}

		if _, ok := m.Roads[origin]; !ok {
			m.Roads[origin] = make([]*Road, 4)
		}
		if _, ok := m.Roads[dest]; !ok {
			m.Roads[dest] = make([]*Road, 4)
		}

		m.Roads[origin][directionType] = &Road{Remote: dest}
		m.Roads[dest][directionType.Opposite()] = &Road{Remote: origin}
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
