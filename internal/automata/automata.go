package automata

import (
	"github.com/BaldiSlayer/rofl-lab2/pkg/models"
)

func SpecialState() models.Cell {
	return models.Cell{X: -228, Y: -228}
}

type Transitions struct {
	ts map[models.Cell]map[byte]models.Cell
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

// DFA - детерминированный конечный автомат
type DFA struct {
	startState  models.Cell
	finalStates map[models.Cell]struct{}
	alphabet    []byte
	// [откуда][по_символу]куда_пришли
	transitions Transitions
	states      []models.Cell
}

func NewDFA(
	startState models.Cell,
	finalStates map[models.Cell]struct{},
	alphabet []byte,
	transitions Transitions,
	states []models.Cell,
) *DFA {
	return &DFA{
		startState:  startState,
		finalStates: finalStates,
		alphabet:    alphabet,
		transitions: transitions,
		states:      states,
	}
}
