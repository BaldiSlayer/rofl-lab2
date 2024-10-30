package mat

import (
	"errors"
	"fmt"
	"github.com/BaldiSlayer/rofl-lab2/internal/automata"
	"github.com/BaldiSlayer/rofl-lab2/internal/maze"
	"github.com/BaldiSlayer/rofl-lab2/internal/mazegen"
	"github.com/BaldiSlayer/rofl-lab2/internal/models"
)

const (
	allCells = iota
	onlyOut
	onlyIn
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

// TODO перепроверь границы
func isSpecial(cell models.Cell, width, height int) bool {
	return cell.X < -1 || cell.X > width || cell.Y < -1 || cell.Y > height
}

// addTransitions добавляет для
func (r *Realization) addTransitions(
	transitions map[models.Cell]map[string]models.Cell,
	i, j int,
) map[models.Cell]map[string]models.Cell {
	directions := []models.Vector{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	alphabet := []string{"N", "S", "W", "E"}

	for idx, dir := range directions {
		src := models.Cell{X: j, Y: i}
		dst := models.Cell{X: j + dir.X, Y: i + dir.Y}

		// в спец клетку(и) можем попасть только из "каймы", там точно нет стенок
		if isSpecial(dst, r.width, r.height) {
			transitions[src][alphabet[idx]] = automata.SpecialState()

			continue
		}

		// если можем пройти - добавить переход в некст клетку, иначе петля
		if r.maze.CanGo(
			src,
			dst,
		) {
			transitions[src][alphabet[idx]] = dst
		} else {
			transitions[src][alphabet[idx]] = src
		}
	}

	return transitions
}

// mazeIterator итерируется по лабиринту, i это y, j это x
// TODO сделать mode битовой маской
func (r *Realization) mazeIterator(mode int, f func(y, x int)) {
	if mode != onlyOut {
		for i := 0; i < r.height; i++ {
			for j := 0; j < r.width; j++ {
				f(i, j)
			}
		}
	}

	if mode != onlyIn {
		// добавим сверху
		j := -1
		for i := 0; i < r.height; i++ {
			f(i, j)
		}

		// добавим справа
		j = r.width
		for i := 0; i < r.height; i++ {
			f(i, j)
		}

		// добавим сверху
		i := -1
		for j := 0; j < r.width; j++ {
			f(i, j)
		}

		// добавим снизу
		i = r.height
		for j := 0; j < r.width; j++ {
			f(i, j)
		}

		// уголки
		f(-1, -1)
		f(-1, r.width)
		f(r.height, -1)
		f(r.height, r.width)
	}
}

func (r *Realization) getAllStates() []models.Cell {
	states := make([]models.Cell, 0, 1+r.width+r.height)
	states = append(states, automata.SpecialState())

	r.mazeIterator(allCells, func(y, x int) {
		states = append(states, models.Cell{X: x, Y: y})
	})

	return states
}

func (r *Realization) getFinalStates() map[models.Cell]struct{} {
	finalStates := make(map[models.Cell]struct{})
	finalStates[automata.SpecialState()] = struct{}{}

	r.mazeIterator(onlyOut, func(y, x int) {
		finalStates[models.Cell{X: x, Y: y}] = struct{}{}
	})

	return finalStates
}

func (r *Realization) getTransitions() map[models.Cell]map[string]models.Cell {
	transitions := make(map[models.Cell]map[string]models.Cell)

	r.mazeIterator(allCells, func(y, x int) {
		transitions = r.addTransitions(transitions, y, x)
	})

	return transitions
}

// toDFA переводит r.maze в детерминированный конечный автомат
// TODO а это точно надо навешивать как метод на Realization или можно сделать
// toDFA(maze)
func (r *Realization) toDFA() *automata.DFA {
	return automata.NewDFA(
		models.Cell{X: 0, Y: 0},
		r.getFinalStates(),
		[]string{"N", "S", "W", "E"},
		r.getTransitions(),
		r.getAllStates(),
	)
}

func (r *Realization) Generate() error {
	var err error

	r.maze, err = r.mazeGenerator.Generate(r.width, r.height)

	r.mazeDFA = r.toDFA()

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
