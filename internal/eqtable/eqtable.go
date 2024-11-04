package eqtable

import (
	"github.com/BaldiSlayer/rofl-lab2/internal/automata"
)

type TableParts struct {
	Prefixes []string
	Suffixes []string
	// [индекс_префикса][индекс_суффикса]
	Answers [][]bool
}

type EqTable interface {
	ToDFA() *automata.DFA
}
