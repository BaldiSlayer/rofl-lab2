package automata

import (
	"github.com/BaldiSlayer/rofl-lab2/internal/defaults"
	"github.com/BaldiSlayer/rofl-lab2/pkg/models"
)

func SpecialState() models.Cell {
	return models.Cell{X: -228, Y: -228}
}

type Transitions struct {
	ts map[models.Cell]map[byte]models.Cell
}

func NewTransitionsFromMap(ts map[models.Cell]map[byte]models.Cell) Transitions {
	return Transitions{
		ts: ts,
	}
}

func NewTransitions() Transitions {
	return Transitions{
		ts: make(map[models.Cell]map[byte]models.Cell),
	}
}

func (t *Transitions) Add(src, dst models.Cell, symbol byte) {
	if _, ok := t.ts[src]; !ok {
		t.ts[src] = make(map[byte]models.Cell)
	}

	t.ts[src][symbol] = dst
}

func (t *Transitions) Has(src models.Cell, symbol byte) bool {
	srcTs := t.ts[src]

	if srcTs == nil {
		return false
	}

	_, ok := srcTs[symbol]

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

func (dfa *DFA) Transitions() map[models.Cell]map[byte]models.Cell {
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

func (dfa *DFA) HasTransition(src models.Cell, symbol byte) bool {
	return dfa.transitions.Has(src, symbol)
}

func (dfa *DFA) AddState(state models.Cell) {
	dfa.states[state] = struct{}{}
}

func (dfa *DFA) AddTransition(src, dst models.Cell, symbol byte) {
	if _, ok := dfa.transitions.ts[src]; !ok {
		dfa.transitions.ts[src] = make(map[byte]models.Cell)
	}

	dfa.transitions.ts[src][symbol] = dst
}

func (dfa *DFA) AddFinalState(state models.Cell) {
	dfa.finalStates[state] = struct{}{}
}
