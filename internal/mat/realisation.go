package mat

import (
	"errors"
	"github.com/BaldiSlayer/rofl-lab2/internal/maze"
	"github.com/BaldiSlayer/rofl-lab2/internal/mazegen"
	"github.com/BaldiSlayer/rofl-lab2/internal/models"
)

type Realization struct {
	width  int
	height int

	mazeGenerator mazegen.MazeGenerator
	maze          *maze.Maze
}

func NewRealization(gen mazegen.MazeGenerator, width, height int) *Realization {
	return &Realization{
		width:         width,
		height:        height,
		mazeGenerator: gen,
	}
}

func getPosAfterStep(i, j int, a byte) (int, int) {
	// из 0, 0 нельзя выйти из лабиринта
	if i == 0 && j == 0 {
		if a == 'W' || a == 'N' {
			return i, j
		}
	}

	switch a {
	case 'N':
		return i - 1, j
	case 'S':
		return i + 1, j
	case 'W':
		return i, j - 1
	case 'E':
		return i, j + 1
	}

	return i, j
}

// Include осуществляет проверку запроса на вхождение
func (r *Realization) Include(query string) (bool, error) {
	i, j := 0, 0

	for sPos := range query {
		newPosI, newPosJ := getPosAfterStep(i, j, query[sPos])

		if r.maze.IsOut(newPosI, newPosJ) {
			i, j = newPosI, newPosJ

			continue
		}

		if r.maze.IsEmpty(newPosI, newPosJ) {
			i, j = newPosI, newPosJ
		}
	}

	return r.maze.IsOut(i, j), nil
}

func (r *Realization) Equal() (models.EqualResponse, error) {
	return models.EqualResponse{}, nil
}

func (r *Realization) Generate() error {
	var err error

	r.maze, err = r.mazeGenerator.Generate(r.width, r.height)

	return err
}

func (r *Realization) Print() ([]string, error) {
	if r.maze == nil {
		return nil, errors.New("failed to print maze: no generated maze")
	}

	return r.maze.Print(), nil
}
