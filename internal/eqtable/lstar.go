package eqtable

import (
	"github.com/BaldiSlayer/rofl-lab2/internal/automata"
	"strings"
)

// OverMaze - таблица, которую мне передают при запросе на эквивалентность
type LStar struct {
	prefixes []string
	suffixes []string
	// [индекс_префикса][индекс_суффикса]
	answers []string
}

func NewLStar(parts TableParts) *LStar {
	answers := make([]string, 0, len(parts.Answers))

	for _, line := range parts.Answers {
		lineAnswers := strings.Builder{}

		for _, val := range line {
			if val {
				lineAnswers.WriteByte('+')
			} else {
				lineAnswers.WriteByte('-')
			}
		}

		answers = append(answers, lineAnswers.String())
	}

	return &LStar{
		prefixes: parts.Prefixes,
		suffixes: parts.Suffixes,
		answers:  answers,
	}
}

func getClassesOfEquivalence() {
	classesOfEquivalence := make([]int, 0)
	_ = classesOfEquivalence
}

func (lstar *LStar) ToDFA() *automata.DFA {
	aut := automata.NewEmptyDFA()

	_ = aut

	// aut.AddState()

	return &automata.DFA{}
}
