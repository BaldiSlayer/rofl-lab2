package mat

import (
	"errors"
	"fmt"
	"github.com/BaldiSlayer/rofl-lab2/internal/maze"

	"github.com/BaldiSlayer/rofl-lab2/internal/mazegen"
	"github.com/BaldiSlayer/rofl-lab2/internal/models"
)

type Realization struct {
	width  int
	height int

	mazeGenerator *mazegen.LightWallsGenerator
	maze          *maze.ThinWalled
}

func NewRealization(gen *mazegen.LightWallsGenerator, width, height int) *Realization {
	return &Realization{
		width:         width,
		height:        height,
		mazeGenerator: gen,
	}
}

// walk возвращает позицию после прохождения по лабиринту по пути query
func (r *Realization) walk(query string) (int, int) {
	i, j := 0, 0

	for sPos := range query {
		i, j = r.maze.GetPosAfterStep(i, j, query[sPos])
	}

	return i, j
}

// Include осуществляет проверку запроса на вхождение
func (r *Realization) Include(query string) (bool, error) {
	return r.maze.IsOut(r.walk(query)), nil
}

type reachableResponse struct {
	allReachable  bool
	notReachableX int
	notReachableY int
}

// allCellsAreReachable проверяет, чтобы все клетки лабиринта были достижимы
func (r *Realization) allCellsAreReachable(prefixes []string) (reachableResponse, error) {
	reachableCells := make(map[int]struct{})

	for _, prefix := range prefixes {
		i, j := r.walk(prefix)
		reachableCells[i*r.width+j] = struct{}{}
	}

	for i := 0; i < r.width*r.height; i++ {
		if _, ok := reachableCells[i]; !ok {
			return reachableResponse{
				allReachable:  false,
				notReachableX: i % r.width,
				notReachableY: i / r.width,
			}, nil
		}
	}

	return reachableResponse{
		allReachable: true,
	}, nil
}

func (r *Realization) Equal(prefixes []string, suffixes []string, matrix [][]bool) (models.EqualResponse, error) {
	// сначала проверяем, что по префиксам мы доходим до всех клеток
	allReachableResult, err := r.allCellsAreReachable(prefixes)
	if err != nil {
		return models.EqualResponse{}, fmt.Errorf("failed to check cells achievability: %w", err)
	}

	if !allReachableResult.allReachable {
		// теперь нужно найти путь от непосещенной клетки до старта и от нее же до выхода
		// сконкатенировать два этих пути, это и будет контрпримером

		return models.EqualResponse{
			Equal: false,
			CounterExample: models.CounterExample{
				CounterExample: "",
			},
		}, nil
	}

	// потом строим ДКА и проводим проверки уже с ним

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
