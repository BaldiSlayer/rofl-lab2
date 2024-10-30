package mat

import (
	"fmt"
	"testing"
)

func TestRealization_mazeIterator(t *testing.T) {
	r := Realization{height: 2, width: 2}
	visited := make(map[string]bool)

	f := func(y, x int) {
		visited[fmt.Sprintf("%d,%d", y, x)] = true
	}

	r.mazeIterator(onlyOut, f)
	if len(visited) != 12 {
		t.Errorf("Expected 8 visited positions, got %d", len(visited))
	}

	visited = make(map[string]bool)

	r.mazeIterator(onlyIn, f)
	if len(visited) != 4 {
		t.Errorf("Expected 4 visited positions, got %d", len(visited))
	}

	visited = make(map[string]bool)

	r.mazeIterator(allCells, f)
	if len(visited) != 16 {
		t.Errorf("Expected 12 visited positions, got %d", len(visited))
	}
}
