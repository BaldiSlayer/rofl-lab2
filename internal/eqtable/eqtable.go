package eqtable

import (
	"github.com/BaldiSlayer/rofl-lab2/internal/automata"
	"github.com/BaldiSlayer/rofl-lab2/internal/maze"
)

type TableParts struct {
	Prefixes []string
	Suffixes []string
	// [индекс_префикса][индекс_суффикса]
	Answers [][]bool
}

type EqTable interface {
	ToDFA(maze *maze.ThinWalled) *automata.DFA
}
