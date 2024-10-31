package automata

import (
	"fmt"

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

func (dfa *DFA) GenerateDot() string {
	dot := "digraph DFA {\n"
	dot += "    rankdir=LR;\n"

	// Определение начального состояния
	dot += fmt.Sprintf("    start [label=\"Start\\n(начальное состояние)\", shape=doublecircle];\n")

	// Определение финальных состояний
	for state := range dfa.finalStates {
		dot += fmt.Sprintf("    final_%d_%d [label=\"Final\\n(финальное состояние)\", shape=doublecircle];\n", state.X, state.Y)
	}

	// Переходы между состояниями
	for _, symbols := range dfa.transitions.ts {
		for symbol, dst := range symbols {
			dot += fmt.Sprintf("    start -> final_%d_%d [label=\"%s\"];\n", dst.X, dst.Y, string(symbol))
		}
	}

	dot += "}\n"
	return dot
}

func (dfa *DFA) GetTransitions() Transitions {
	return dfa.transitions
}
