package mazegen

import (
	"math/rand"
)

type SimpleGenerator struct {
}

type Cell struct {
	X, Y int
}

func isInBounds(x, y, width, height int) bool {
	return x >= 0 && y >= 0 && x < width && y < height
}

func shuffleDirections() {
	var directions = []Cell{
		{X: -1, Y: 0},
		{X: 1, Y: 0},
		{X: 0, Y: -1},
		{X: 0, Y: 1},
	}

	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})
}

func generateMaze(grid [][]bool, x, y, height, width int) {
	var directions = []Cell{
		{X: -1, Y: 0},
		{X: 1, Y: 0},
		{X: 0, Y: -1},
		{X: 0, Y: 1},
	}

	grid[x][y] = true
	shuffleDirections()

	for _, dir := range directions {
		nx, ny := x+dir.X, y+dir.Y
		if isInBounds(nx, ny, height, width) &&
			!grid[nx][ny] && isInBounds(nx+dir.X, ny+dir.Y, height, width) &&
			!grid[nx+dir.X][ny+dir.Y] {
			generateMaze(grid, nx+dir.X, ny+dir.Y, height, width)
			grid[nx][ny] = true
		}
	}
}
