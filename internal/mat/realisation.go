package mat

import (
	"errors"

	"github.com/BaldiSlayer/rofl-lab2/internal/automata"
	"github.com/BaldiSlayer/rofl-lab2/internal/defaults"
	"github.com/BaldiSlayer/rofl-lab2/internal/eqtable"
	"github.com/BaldiSlayer/rofl-lab2/internal/maze"
	"github.com/BaldiSlayer/rofl-lab2/internal/mazegen"
	"github.com/BaldiSlayer/rofl-lab2/pkg/models"
)

type Implementation struct {
	width  int
	height int

	mazeGenerator *mazegen.LightWallsGenerator
	maze          *maze.ThinWalled

	mazeDFA *automata.DFA
}

func NewImplementation(gen *mazegen.LightWallsGenerator, width, height int) *Implementation {
	return &Implementation{
		width:         width,
		height:        height,
		mazeGenerator: gen,
	}
}

// walk возвращает позицию после прохождения по лабиринту по пути query
func (r *Implementation) walk(query string) models.Cell {
	start := models.Cell{X: 0, Y: 0}

	for sPos := range query {
		start = r.maze.GetPosAfterStep(start, query[sPos])
	}

	return start
}

// Include осуществляет проверку запроса на вхождение
func (r *Implementation) Include(query string) (bool, error) {
	return r.mazeDFA.Include(query, r.maze.IsSpecial), nil
}

func (r *Implementation) genCounterForStates(state models.Cell) string {
	// необходимо вернуть путь от стейта до (-1, -1) и пойти в SpecialState путем выхода из каймы
	if state == defaults.SpecialState() {
		return r.mazeDFA.GetPath(state, false, true) + "N"
	}

	// необходимо пройти от старта до выхода через state, для этого
	// мы найдем путь state->start, перевернем его
	// и сконкатенируем с state->final
	return r.mazeDFA.GetPath(state, true, true) +
		r.mazeDFA.GetPath(state, false, false)
}

func (r *Implementation) genCounterForFinalState(state models.Cell) string {
	// необходимо вернуть путь от стейта до (-1, -1) и пойти в SpecialState путем выхода из каймы
	if state == defaults.SpecialState() {
		return r.mazeDFA.GetPath(state, false, true) + "N"
	}

	// пройдем из этого финального состояния до стартового
	return r.mazeDFA.GetPath(state, true, true)
}

func (r *Implementation) genCounterForTransition(src automata.Transition, dst models.Cell) string {
	if dst == defaults.SpecialState() {
		// иду обратно, чтобы не путать лернер и не возвращать ему слишком далекие клетки
		return r.mazeDFA.GetPath(src.Src, true, true) + string(src.Symbol) +
			map[byte]string{
				'N': "S",
				'S': "N",
				'W': "E",
				'E': "W",
			}[src.Symbol]
	}

	// идем от старта до src потом делаем переход по нужному символу
	// затем идем от места, куда перешли до финальной клетки
	return r.mazeDFA.GetPath(src.Src, true, true) +
		string(src.Symbol) + r.mazeDFA.GetPath(dst, false, false)
}

// getNonEqualFinalStatePath ищет контрпример по наличию состояния
func (r *Implementation) getNonEqualStatePath(dfaFromTable *automata.DFA) string {
	for mazeState := range r.mazeDFA.States() {
		if !dfaFromTable.HasState(mazeState) {
			return r.genCounterForStates(mazeState)
		}
	}

	return ""
}

// getNonEqualFinalStatePath ищет контрпример по наличию финального состояния
func (r *Implementation) getNonEqualFinalStatePath(dfaFromTable *automata.DFA) string {
	for mazeFinalState := range r.mazeDFA.GetFinalStates() {
		if !dfaFromTable.HasFinalState(mazeFinalState) {
			return r.genCounterForFinalState(mazeFinalState)
		}
	}

	return ""
}

// getNonEqualTransitions ищет контрпример по переходам
func (r *Implementation) getNonEqualTransitions(dfaFromTable *automata.DFA) string {
	mazeTransitions := r.mazeDFA.Transitions()

	for src := range mazeTransitions {
		if !dfaFromTable.HasTransition(src) {
			return r.genCounterForTransition(src, mazeTransitions[src])
		}
	}

	return ""
}

// getCounterExample ищет контрпример
func (r *Implementation) getCounterExample(dfaFromTable *automata.DFA) string {
	counterGens := []func(dfaFromTable *automata.DFA) string{
		r.getNonEqualStatePath,
		r.getNonEqualFinalStatePath,
		r.getNonEqualTransitions,
	}

	for _, generator := range counterGens {
		if val := generator(dfaFromTable); val != "" {
			return val
		}
	}

	return ""
}

func (r *Implementation) Equal(tableParts eqtable.TableParts) (models.EqualResponse, error) {
	// TODO надо придумать, как тут сделать покрасивее
	eqTable := eqtable.NewOverMaze(tableParts, r.maze)

	// context.TODO хуй

	dfaFromTable := eqTable.ToDFA()

	counter := r.getCounterExample(dfaFromTable)
	if counter != "" {
		return models.EqualResponse{
			Equal: false,
			CounterExample: models.CounterExample{
				CounterExample: counter,
			},
		}, nil
	}

	return models.EqualResponse{
		Equal: true,
	}, nil
}

func (r *Implementation) Generate() error {
	var err error

	r.maze, err = r.mazeGenerator.Generate(r.width, r.height)

	r.mazeDFA = r.maze.ToDFA()

	return err
}

func (r *Implementation) Visualize() ([]string, error) {
	if r.maze == nil {
		return nil, errors.New("failed to print maze: no generated maze")
	}

	r.maze.Print()

	return nil, nil
}
