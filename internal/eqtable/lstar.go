package eqtable

import (
	"github.com/BaldiSlayer/rofl-lab2/internal/automata"
	"sort"
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

type ClassMember struct {
	// lineNumber изначальный номер строки
	lineNumber int

	prefixVal string
}

type ClassesOfEquivalence struct {
	classesOfEquivalence map[string][]ClassMember
	wordToClass          map[string]string
}

func (eClasses *ClassesOfEquivalence) getClassesOfEquivalence(
	prefixes []string,
	answers []string,
) *ClassesOfEquivalence {
	classesOfEquivalence := make(map[string][]ClassMember, 0)
	wordToClass := make(map[string]string)

	for i, prefix := range prefixes {
		if _, ok := classesOfEquivalence[answers[i]]; !ok {
			classesOfEquivalence[answers[i]] = make([]ClassMember, 0)

			classesOfEquivalence[answers[i]] = append(classesOfEquivalence[answers[i]], ClassMember{
				lineNumber: i,
				prefixVal:  prefix,
			})

			wordToClass[prefix] = answers[i]
		}
	}

	// сортируем
	for _, members := range classesOfEquivalence {
		sort.Slice(members, func(i, j int) bool {
			return members[i].prefixVal < members[j].prefixVal
		})
	}

	return &ClassesOfEquivalence{
		classesOfEquivalence: classesOfEquivalence,
		wordToClass:          wordToClass,
	}
}

func (lstar *LStar) ToDFA() *automata.DFA {
	aut := automata.NewEmptyDFA()

	_ = aut

	// aut.AddState()

	return &automata.DFA{}
}
