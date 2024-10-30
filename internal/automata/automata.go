package automata

import "github.com/BaldiSlayer/rofl-lab2/internal/models"

func SpecialState() models.Cell {
	return models.Cell{X: -228, Y: -228}
}

// DFA - детерминированный конечный автомат
type DFA struct {
	startState  models.Cell
	finalStates map[models.Cell]struct{}
	alphabet    []string
	// [откуда][по_символу]куда_пришли
	transitions map[models.Cell]map[string]models.Cell
	states      []models.Cell
}

func NewDFA(
	startState models.Cell,
	finalStates map[models.Cell]struct{},
	alphabet []string,
	transitions map[models.Cell]map[string]models.Cell,
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
