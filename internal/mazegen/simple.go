package mazegen

import (
	"github.com/BaldiSlayer/rofl-lab2/internal/maze"
	"math/rand"
	"time"
)

type SimpleGenerator struct {
}

type Cell struct {
	X, Y int
}

var directions = []Cell{
	{X: -1, Y: 0},
	{X: 1, Y: 0},
	{X: 0, Y: -1},
	{X: 0, Y: 1},
}

func isInBounds(x, y, width, height int) bool {
	return x >= 0 && y >= 0 && x < width && y < height
}

func shuffleDirections() {
	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})
}

func generateMaze(grid [][]bool, x, y, height, width int) {
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

func NewSimpleGenerator() *SimpleGenerator {
	return &SimpleGenerator{}
}

func (sg *SimpleGenerator) Generate(width, height int) (*maze.Maze, error) {
	rand.Seed(time.Now().UnixNano())

	grid := make([][]bool, height)
	for i := range grid {
		grid[i] = make([]bool, width)
	}

	generateMaze(grid, 0, 0, height, width)

	y := height - 1
	for x := 0; x < width; x++ {
		if !grid[y][x] && grid[y-1][x] {
			grid[y][x] = true

			break
		}
	}

	return maze.New(grid), nil
}
