package maze

import (
	"fmt"
	"testing"

	"github.com/BaldiSlayer/rofl-lab2/pkg/bmstore"
	"github.com/stretchr/testify/require"
)

func TestThinWalled_mazeIterator(t *testing.T) {
	r := ThinWalled{height: 2, width: 2}
	visited := make(map[string]bool)

	f := func(y, x int) {
		visited[fmt.Sprintf("%d,%d", y, x)] = true
	}

	r.mazeIterator(bmstore.Store(outCells), f)
	require.Equal(t, 12, len(visited))

	visited = make(map[string]bool)

	r.mazeIterator(bmstore.Store(inCells), f)
	require.Equal(t, 4, len(visited))

	visited = make(map[string]bool)

	r.mazeIterator(bmstore.Store(inCells, outCells), f)
	require.Equal(t, 16, len(visited))
}
