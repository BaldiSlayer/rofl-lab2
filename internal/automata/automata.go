package automata

import (
	"container/list"

	"github.com/BaldiSlayer/rofl-lab2/internal/defaults"
	"github.com/BaldiSlayer/rofl-lab2/pkg/models"
)

type Transition struct {
	Src    models.Cell
	Symbol byte
}

type Transitions struct {
	ts map[Transition]models.Cell
}

func NewTransitions() Transitions {
	return Transitions{
		ts: make(map[Transition]models.Cell),
	}
}

func (t *Transitions) Add(src, dst models.Cell, symbol byte) {
	t.ts[Transition{Src: src, Symbol: symbol}] = dst
}

func (t *Transitions) Has(transition Transition) bool {
	_, ok := t.ts[transition]

	return ok
}

// DFA - детерминированный конечный автомат
type DFA struct {
	startState  models.Cell
	finalStates map[models.Cell]struct{}
	alphabet    []byte
	// [откуда][по_символу]куда_пришли
	transitions Transitions
	states      map[models.Cell]struct{}
}

func NewEmptyDFA() *DFA {
	return &DFA{
		startState:  models.Cell{X: 0, Y: 0},
		finalStates: make(map[models.Cell]struct{}),
		alphabet:    defaults.GetAlphabet(),
		transitions: NewTransitions(),
		states:      make(map[models.Cell]struct{}),
	}
}

func NewDFA(
	startState models.Cell,
	finalStates map[models.Cell]struct{},
	alphabet []byte,
	transitions Transitions,
	states map[models.Cell]struct{},
) *DFA {
	return &DFA{
		startState:  startState,
		finalStates: finalStates,
		alphabet:    alphabet,
		transitions: transitions,
		states:      states,
	}
}

func (dfa *DFA) States() map[models.Cell]struct{} {
	return dfa.states
}

func (dfa *DFA) GetFinalStates() map[models.Cell]struct{} {
	return dfa.finalStates
}

func (dfa *DFA) Transitions() map[Transition]models.Cell {
	return dfa.transitions.ts
}

func (dfa *DFA) HasState(state models.Cell) bool {
	_, ok := dfa.states[state]

	return ok
}

func (dfa *DFA) HasFinalState(state models.Cell) bool {
	_, ok := dfa.finalStates[state]

	return ok
}

func (dfa *DFA) HasTransition(t Transition) bool {
	return dfa.transitions.Has(t)
}

func (dfa *DFA) AddState(state models.Cell) {
	dfa.states[state] = struct{}{}
}

func (dfa *DFA) AddTransition(src, dst models.Cell, symbol byte) {
	dfa.transitions.Add(src, dst, symbol)
}

func (dfa *DFA) AddFinalState(state models.Cell) {
	dfa.finalStates[state] = struct{}{}
}

func vecToLetter(vec models.Vector) byte {
	return map[models.Vector]byte{
		models.Vector{X: 0, Y: -1}: 'N',
		models.Vector{X: 0, Y: 1}:  'S',
		models.Vector{X: -1, Y: 0}: 'W',
		models.Vector{X: 1, Y: 0}:  'E',
	}[vec]
}

func symbolToVec(symbol byte) models.Vector {
	return map[byte]models.Vector{
		'N': {X: 0, Y: -1},
		'S': {X: 0, Y: 1},
		'W': {X: -1, Y: 0},
		'E': {X: 1, Y: 0},
	}[symbol]
}

func restorePath(start, end models.Cell, prev map[models.Cell]models.Cell) string {
	path := ""

	cur := end
	for cur != start {
		path = string(vecToLetter(models.Vector{
			X: cur.X - prev[cur].X,
			Y: cur.Y - prev[cur].Y,
		})) + path

		cur = prev[cur]
	}

	return path
}

func (dfa *DFA) bfs(
	src models.Cell,
	// получше бы назвать
	bfsEnder func(current models.Cell) bool,
) (models.Cell, map[models.Cell]models.Cell) {
	queue := list.New()

	queue.PushBack(src)

	visited := make(map[models.Cell]struct{})
	visited[src] = struct{}{}

	prev := make(map[models.Cell]models.Cell)
	prev[src] = src

	for queue.Len() > 0 {
		current := queue.Front().Value.(models.Cell)
		queue.Remove(queue.Front())

		if bfsEnder(current) {
			return current, prev
		}

		// обходим всех соседей
		for _, letter := range defaults.GetAlphabet() {
			nextState := dfa.transitions.ts[Transition{Src: current, Symbol: letter}]

			// если не посещали
			if _, ok := visited[nextState]; !ok {
				visited[nextState] = struct{}{}

				queue.PushBack(nextState)
			}
		}
	}

	return models.Cell{X: 0, Y: 0}, nil
}

// GetPath возвращает путь по ДКА от состояния Src до стартового или финального
// состояния в зависимости от выставленного флага (хм, может лучше передавать функцию чекер)
func (dfa *DFA) GetPath(src models.Cell, start, reversed bool) string {
	// если хотим придти в стартовое состояние
	bfsEnder := func(current models.Cell) bool {
		return current == dfa.startState
	}

	if !start {
		bfsEnder = func(current models.Cell) bool {
			_, ok := dfa.finalStates[current]

			return ok
		}
	}

	// чето функции получились не такие красивые, как хотелось =(

	end, prev := dfa.bfs(src, bfsEnder)

	ans := restorePath(src, end, prev)

	if reversed {
		return reverse(ans)
	}

	return ans
}

func reverse(s string) string {
	runes := []rune(s)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

// Include проверяет, распознает ли автомат слово
func (dfa *DFA) Include(query string, isSpecial func(cell models.Cell) bool) bool {
	curState := defaults.GetStartState()
	nonDKAState := models.Cell{}

	for _, letter := range query {
		vec := symbolToVec(byte(letter))

		if curState == defaults.SpecialState() {
			nonDKAState = models.Cell{X: nonDKAState.X + vec.X, Y: nonDKAState.Y + vec.Y}

			// если смогли вернуться в лабиринт/кайму
			if !isSpecial(nonDKAState) {
				curState = nonDKAState
			}

			continue
		}

		nonDKAState = curState
		curState = dfa.Transitions()[Transition{Src: curState, Symbol: byte(letter)}]
		if curState == defaults.SpecialState() {
			nonDKAState = models.Cell{X: nonDKAState.X + vec.X, Y: nonDKAState.Y + vec.Y}
		}
	}

	_, ok := dfa.finalStates[curState]

	return ok
}
