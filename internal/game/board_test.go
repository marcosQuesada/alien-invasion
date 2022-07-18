package game

import "testing"

func TestOnAddVisitorWithEmptyCityReportsSuccess(t *testing.T) {
	c := NewCity("foo")
	a := NewAlien("alien-1", 100)
	if err := c.AddVisitor(a); err != nil {
		t.Fatalf("unexpected error adding visitor %v", err)
	}
}

func TestOnAddVisitorWithVisitedCityReportsWar(t *testing.T) {
	c := NewCity("foo")
	a1 := NewAlien("alien-1", 100)
	_ = c.AddVisitor(a1)
	a2 := NewAlien("alien-2", 100)
	err := c.AddVisitor(a2)
	if err == nil {
		t.Fatal("expected war error")
	}

	e, ok := err.(*ErrCityWarStarted)
	if !ok {
		t.Fatalf("unexpected error type got %T", err)
	}

	if expected, got := c.name, e.city; expected != got {
		t.Fatalf("values do not match, expected %s got %s", expected, got)
	}

	if expected, got := a1.name, e.alienA; expected != got {
		t.Fatalf("values do not match, expected %s got %s", expected, got)
	}

	if expected, got := a2.name, e.alienB; expected != got {
		t.Fatalf("values do not match, expected %s got %s", expected, got)
	}
}
