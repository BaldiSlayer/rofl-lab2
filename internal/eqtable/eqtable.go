package eqtable

import (
	"github.com/BaldiSlayer/rofl-lab2/internal/cautomata"
)

type TableParts struct {
	Prefixes []string
	Suffixes []string
	// [индекс_префикса][индекс_суффикса]
	Answers [][]bool
}

type EqTable interface {
	ToDFA() *cautomata.DFA
}
