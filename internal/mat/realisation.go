package mat

import (
	"errors"

	"github.com/BaldiSlayer/rofl-lab2/internal/mazegen"
	"github.com/BaldiSlayer/rofl-lab2/internal/models"
)

type Realization struct {
	width  int
	height int

	mazeGenerator *mazegen.LightWallsGenerator
	maze          *mazegen.ThinWalledMaze
}

func NewRealization(gen *mazegen.LightWallsGenerator, width, height int) *Realization {
	return &Realization{
		width:         width,
		height:        height,
		mazeGenerator: gen,
	}
}

// Include осуществляет проверку запроса на вхождение
func (r *Realization) Include(query string) (bool, error) {
	i, j := 0, 0

	for sPos := range query {
		i, j = r.maze.GetPosAfterStep(i, j, query[sPos])
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

	r.maze.Print()

	return nil, nil
}
