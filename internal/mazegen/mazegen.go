package mazegen

import "github.com/BaldiSlayer/rofl-lab2/internal/maze"

type MazeGenerator interface {
	Generate(width, height int) (*maze.ThinWalled, error)
}
