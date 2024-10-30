package mat

import (
	"errors"
	"fmt"
	"github.com/BaldiSlayer/rofl-lab2/internal/automata"
	"github.com/BaldiSlayer/rofl-lab2/internal/maze"
	"github.com/BaldiSlayer/rofl-lab2/internal/mazegen"
	"github.com/BaldiSlayer/rofl-lab2/internal/models"
)

type Realization struct {
	width  int
	height int

	mazeGenerator *mazegen.LightWallsGenerator
	maze          *maze.ThinWalled

	mazeDFA *automata.DFA
}

func NewRealization(gen *mazegen.LightWallsGenerator, width, height int) *Realization {
	return &Realization{
		width:         width,
		height:        height,
		mazeGenerator: gen,
	}
}

// walk возвращает позицию после прохождения по лабиринту по пути query
func (r *Realization) walk(query string) models.Cell {
	start := models.Cell{X: 0, Y: 0}

	for sPos := range query {
		start = r.maze.GetPosAfterStep(start, query[sPos])
	}

	return start
}

// Include осуществляет проверку запроса на вхождение
func (r *Realization) Include(query string) (bool, error) {
	return r.maze.IsOut(r.walk(query)), nil
}

type reachableResponse struct {
	allReachable     bool
	notReachableCell models.Cell
}

// allCellsAreReachable проверяет, чтобы все клетки лабиринта были достижимы
func (r *Realization) allCellsAreReachable(prefixes []string) (reachableResponse, error) {
	// храним информацию о клетках, в которые мы смогли прийти
	reachableCells := make(map[models.Cell]struct{})

	for _, prefix := range prefixes {
		cell := r.walk(prefix)

		reachableCells[cell] = struct{}{}
	}

	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			cell := models.Cell{X: x, Y: y}

			if _, ok := reachableCells[cell]; !ok {
				return reachableResponse{
					allReachable:     false,
					notReachableCell: cell,
				}, nil
			}
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

	r.mazeDFA = r.maze.ToDFA()

	return err
}

func (r *Realization) Print() ([]string, error) {
	if r.maze == nil {
		return nil, errors.New("failed to print maze: no generated maze")
	}

	r.maze.Print()

	return nil, nil
}

func (r *Realization) tableToDFA() *automata.DFA {
	alphabet := []string{"N", "S", "W", "E"}
	_ = alphabet

	directions := []models.Vector{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	_ = directions

	startState := models.Cell{X: 0, Y: 0}
	_ = startState

	return &automata.DFA{}
}
