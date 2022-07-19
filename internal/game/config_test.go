package game

import (
	"errors"
	"log"
	"reflect"
	"testing"
)

func TestParseEdgesFromLineGetsMappedAsUnidirectionalEdges(t *testing.T) {
	samples := []struct {
		line                     []string
		totalUnidirectionalRoads int
		edges                    map[DirectionType]CityName
		error                    error
	}{
		{
			line:  []string{"fewfwefewf"},
			error: ErrInvalidLineFormat,
		},
		{
			line:  []string{"non-existent-0-0", "east=X-1-1"},
			error: ErrInvalidLineFormat,
		},
		{
			line:  []string{"northX-0-0", "east=X-1-1"},
			error: ErrInvalidLineFormat,
		},
		{
			line:                     []string{"north=X-0-0", "east=X-1-1"},
			totalUnidirectionalRoads: 2,
			edges: map[DirectionType]CityName{
				NorthDirection: CityName("X-0-0"),
				EastDirection:  CityName("X-1-1"),
			},
		},
	}

	for _, sample := range samples {
		ed, err := parseEdges(sample.line)
		if sample.error != nil && !errors.Is(err, sample.error) {
			t.Fatalf("expected error %v", sample.error)
		}

		if expected, got := sample.totalUnidirectionalRoads, len(ed); expected != got {
			log.Fatalf("total edges dos not match, expected %d got %d", expected, got)
		}

		for directionType, cityName := range sample.edges {
			v, ok := ed[directionType]
			if !ok {
				t.Fatalf("Direction type %s not found on edges", directionType)
			}

			if v != cityName {
				log.Fatalf("City Names does not match, expected %s got %s", cityName, v)
			}
		}
	}
}

func TestParseLineDetectBadFormattedLines(t *testing.T) {
	samples := []struct {
		line  string
		error error
	}{
		{
			line:  "",
			error: ErrInvalidLineFormat,
		},
		{
			line:  "xxxxxxx",
			error: ErrInvalidLineFormat,
		},
		{
			line:  "xxxx  xxx",
			error: ErrInvalidLineFormat,
		},
		{
			line:  "xxxx  north=bar",
			error: nil,
		},
		{
			line:  "xxxx  north=",
			error: ErrInvalidLineFormat,
		},
	}

	for _, sample := range samples {
		err := parseLine(sample.line, nil)
		if sample.error != nil && !errors.Is(err, sample.error) {
			t.Fatalf("expected error %v", sample.error)
		}
	}
}

func TestParseLinePopulatesAllDescribedCitiesAndRoadsAsBidirectionalWithSuccess(t *testing.T) {
	samples := []struct {
		line        string
		totalCities int
		planetMap   *PlanetMap
	}{
		{
			line:        "X-1-0 north=X-0-0 east=X-1-1",
			totalCities: 3,
			planetMap: &PlanetMap{
				cities: map[CityName]*City{
					CityName("X-1-0"): NewCity("X-1-0"),
					CityName("X-0-0"): NewCity("X-0-0"),
					CityName("X-1-1"): NewCity("X-1-1"),
				},
				roads: map[CityName][]*Road{
					CityName("X-1-0"): []*Road{
						NorthDirection: {Remote: CityName("X-0-0")},
						EastDirection:  {Remote: CityName("X-1-1")},
						SouthDirection: nil,
						WestDirection:  nil,
					},
					CityName("X-1-1"): []*Road{
						NorthDirection: nil,
						EastDirection:  nil,
						SouthDirection: nil,
						WestDirection:  {Remote: CityName("X-1-0")},
					},
					CityName("X-0-0"): []*Road{
						NorthDirection: nil,
						EastDirection:  nil,
						SouthDirection: {Remote: CityName("X-1-0")},
						WestDirection:  nil,
					},
				},
			},
		},
		{
			line:        "Foo north=Bar west=Baz south=Qu-ux",
			totalCities: 4,
			planetMap: &PlanetMap{
				cities: map[CityName]*City{
					CityName("Foo"):   NewCity("Foo"),
					CityName("Bar"):   NewCity("Bar"),
					CityName("Baz"):   NewCity("Baz"),
					CityName("Qu-ux"): NewCity("Qu-ux"),
				},
				roads: map[CityName][]*Road{
					CityName("Foo"): []*Road{
						NorthDirection: {Remote: CityName("Bar")},
						EastDirection:  nil,
						SouthDirection: {Remote: CityName("Qu-ux")},
						WestDirection:  {Remote: CityName("Baz")},
					},
					CityName("Bar"): []*Road{
						NorthDirection: nil,
						EastDirection:  nil,
						SouthDirection: {Remote: CityName("Foo")},
						WestDirection:  nil,
					},
					CityName("Baz"): []*Road{
						NorthDirection: nil,
						EastDirection:  {Remote: CityName("Foo")},
						SouthDirection: nil,
						WestDirection:  nil,
					},
					CityName("Qu-ux"): []*Road{
						NorthDirection: {Remote: CityName("Foo")},
						EastDirection:  nil,
						SouthDirection: nil,
						WestDirection:  nil,
					},
				},
			},
		},
		{
			line:        "Bar south=Foo west=Bee",
			totalCities: 3,
			planetMap: &PlanetMap{
				cities: map[CityName]*City{
					CityName("Bar"): NewCity("Bar"),
					CityName("Foo"): NewCity("Foo"),
					CityName("Bee"): NewCity("Bee"),
				},
				roads: map[CityName][]*Road{
					CityName("Bar"): []*Road{
						NorthDirection: nil,
						EastDirection:  nil,
						SouthDirection: {Remote: CityName("Foo")},
						WestDirection:  {Remote: CityName("Bee")},
					},
					CityName("Foo"): []*Road{
						NorthDirection: {Remote: CityName("Bar")},
						EastDirection:  nil,
						SouthDirection: nil,
						WestDirection:  nil,
					},
					CityName("Bee"): []*Road{
						NorthDirection: nil,
						EastDirection:  {Remote: CityName("Bar")},
						SouthDirection: nil,
						WestDirection:  nil,
					},
				},
			},
		},
	}

	for _, sample := range samples {
		m := newEmptyMap()
		if err := parseLine(sample.line, m); err != nil {
			t.Fatalf("Unexxpected error %v on line %s", err, sample.line)
		}

		if expected, got := sample.totalCities, len(m.cities); expected != got {
			log.Fatalf("total cities dos not match, expected %d got %d", expected, got)
		}

		if expected, got := len(m.roads), len(sample.planetMap.roads); expected != got {
			log.Fatalf("total edges dos not match, expected %d got %d", expected, got)
		}

		for cityName, roads := range sample.planetMap.roads {
			v, ok := m.roads[cityName]
			if !ok {
				t.Fatalf("unable to find road %s on PlanetMap", cityName)
			}

			if !reflect.DeepEqual(v, roads) {
				t.Fatal("Final PlanetMap and Expected map does not match")
			}
		}
	}
}

func TestLoadMapFromConfigFileWithDuplicatedEdgeDefinitions(t *testing.T) {
	filename := "test.csv"
	m, err := LoaDefinitionsFromFile(filename)
	if err != nil {
		t.Fatalf("unable to load file %s error %v", filename, err)
	}

	if expected, got := 4, len(m.cities); expected != got {
		t.Fatalf("total cities does not match, expected %d got %d", expected, got)
	}

	if expected, got := 4, len(m.roads); expected != got {
		t.Fatalf("total road map does not match, expected %d got %d", expected, got)
	}

	city10 := CityName("CITY-1-0")
	roadsCity10, ok := m.roads[city10]
	if !ok {
		t.Fatalf("Road %s not found", city10)
	}
	if roadsCity10[NorthDirection] == nil || roadsCity10[EastDirection] == nil {
		t.Fatal("expected value")
	}
	if roadsCity10[SouthDirection] != nil || roadsCity10[WestDirection] != nil {
		t.Fatal("not expected value")
	}

	city11 := CityName("CITY-1-1")
	roadsCity11, ok := m.roads[city11]
	if !ok {
		t.Fatalf("Road %s not found", city11)
	}
	if roadsCity11[NorthDirection] == nil || roadsCity11[WestDirection] == nil {
		t.Fatal("expected value")
	}
	if roadsCity11[SouthDirection] != nil || roadsCity11[EastDirection] != nil {
		t.Fatal("not expected value")
	}

	city00 := CityName("CITY-0-0")
	roadsCity00, ok := m.roads[city00]
	if !ok {
		t.Fatalf("Road %s not found", city00)
	}
	if roadsCity00[SouthDirection] == nil || roadsCity00[EastDirection] == nil {
		t.Fatal("expected value")
	}
	if roadsCity00[NorthDirection] != nil || roadsCity00[WestDirection] != nil {
		t.Fatal("not expected value")
	}

	city01 := CityName("CITY-0-1")
	roadsCity01, ok := m.roads[city01]
	if !ok {
		t.Fatalf("Road %s not found", city01)
	}
	if roadsCity01[SouthDirection] == nil || roadsCity01[WestDirection] == nil {
		t.Fatal("expected value")
	}
	if roadsCity01[NorthDirection] != nil || roadsCity01[EastDirection] != nil {
		t.Fatal("not expected value")
	}
}
