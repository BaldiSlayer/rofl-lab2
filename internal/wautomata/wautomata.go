package wautomata

import (
	"github.com/BaldiSlayer/rofl-lab2/internal/defaults"
)

type Transition struct {
	Src    string
	Symbol byte
}

type Transitions struct {
	ts map[Transition]string
}

func NewEmptyTransitions() Transitions {
	return Transitions{
		ts: make(map[Transition]string),
	}
}

func NewTransitions(ts map[Transition]string) Transitions {
	return Transitions{
		ts: ts,
	}
}

func (t *Transitions) Add(src, dst string, symbol byte) {
	t.ts[Transition{Src: src, Symbol: symbol}] = dst
}

func (t *Transitions) Has(transition Transition) bool {
	_, ok := t.ts[transition]

	return ok
}

// DFA - детерминированный конечный автомат
type DFA struct {
	startState  string
	finalStates map[string]struct{}
	alphabet    []byte
	// [откуда][по_символу]куда_пришли
	transitions Transitions
	states      map[string]struct{}
}

func NewEmptyDFA() *DFA {
	return &DFA{
		startState:  "",
		finalStates: make(map[string]struct{}),
		alphabet:    defaults.GetAlphabet(),
		transitions: NewEmptyTransitions(),
		states:      make(map[string]struct{}),
	}
}

func New(
	startState string,
	finalStates map[string]struct{},
	alphabet []byte,
	transitions Transitions,
	states map[string]struct{},
) *DFA {
	return &DFA{
		startState:  startState,
		finalStates: finalStates,
		alphabet:    alphabet,
		transitions: transitions,
		states:      states,
	}
}

func NewDFA(
	startState string,
	finalStates map[string]struct{},
	alphabet []byte,
	transitions Transitions,
	states map[string]struct{},
) *DFA {
	return &DFA{
		startState:  startState,
		finalStates: finalStates,
		alphabet:    alphabet,
		transitions: transitions,
		states:      states,
	}
}

func (dfa *DFA) States() map[string]struct{} {
	return dfa.states
}

func (dfa *DFA) GetFinalStates() map[string]struct{} {
	return dfa.finalStates
}

func (dfa *DFA) Transitions() map[Transition]string {
	return dfa.transitions.ts
}

func (dfa *DFA) HasState(state string) bool {
	_, ok := dfa.states[state]

	return ok
}

func (dfa *DFA) HasFinalState(state string) bool {
	_, ok := dfa.finalStates[state]

	return ok
}

func (dfa *DFA) HasTransition(t Transition) bool {
	return dfa.transitions.Has(t)
}

func (dfa *DFA) AddState(state string) {
	dfa.states[state] = struct{}{}
}

func (dfa *DFA) AddTransition(src, dst string, symbol byte) {
	dfa.transitions.Add(src, dst, symbol)
}

func (dfa *DFA) AddFinalState(state string) {
	dfa.finalStates[state] = struct{}{}
}
