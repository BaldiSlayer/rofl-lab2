package eqtable

import (
	"github.com/BaldiSlayer/rofl-lab2/internal/automata"
	"github.com/BaldiSlayer/rofl-lab2/internal/maze"
)

// OverMaze - таблица, которую мне передают при запросе на эквивалентность
type LStar struct {
	prefixes []string
	suffixes []string
	// [индекс_префикса][индекс_суффикса]
	answers [][]bool
}

func NewLStar(parts TableParts) *LStar {
	return &LStar{
		prefixes: parts.Prefixes,
		suffixes: parts.Suffixes,
		answers:  parts.Answers,
	}
}

func (lstar *LStar) ToDFA(maze *maze.ThinWalled) *automata.DFA {
	return &automata.DFA{}
}
