package eqtable

import (
	"github.com/BaldiSlayer/rofl-lab2/internal/defaults"
	"github.com/BaldiSlayer/rofl-lab2/internal/wautomata"
	"sort"
	"strings"
)

var alphabet = defaults.GetAlphabet()

const (
	wordIn  = '+'
	wordOut = '-'
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
				lineAnswers.WriteByte(wordIn)
			} else {
				lineAnswers.WriteByte(wordOut)
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

func newClassesOfEquivalence(
	prefixes []string,
	answers []string,
) *ClassesOfEquivalence {
	classesOfEquivalence := make(map[string][]ClassMember, 0)
	wordToClass := make(map[string]string)

	for i, prefix := range prefixes {
		if _, ok := classesOfEquivalence[answers[i]]; !ok {
			classesOfEquivalence[answers[i]] = make([]ClassMember, 0)
		}

		classesOfEquivalence[answers[i]] = append(classesOfEquivalence[answers[i]], ClassMember{
			lineNumber: i,
			prefixVal:  prefix,
		})

		wordToClass[prefix] = answers[i]
	}

	// сортируем строки в каждом из классов эквивалентности
	for _, members := range classesOfEquivalence {
		sort.SliceStable(members, func(i, j int) bool {
			return len(members[i].prefixVal) < len(members[j].prefixVal)
		})
	}

	return &ClassesOfEquivalence{
		classesOfEquivalence: classesOfEquivalence,
		wordToClass:          wordToClass,
	}
}

func (lstar *LStar) ToDFA() *wautomata.DFA {
	aut := wautomata.NewEmptyDFA()

	epsIdx := -1
	for i, val := range lstar.suffixes {
		if val == "" {
			epsIdx = i

			break
		}
	}

	eqClasses := newClassesOfEquivalence(lstar.prefixes, lstar.answers)

	// добавляем состояния в автомат
	for ansString, members := range eqClasses.classesOfEquivalence {
		member := members[0]

		aut.AddState(member.prefixVal)

		// помечаем состояние финальным
		if ansString[epsIdx] == wordIn {
			aut.AddFinalState(member.prefixVal)
		}
	}

	if !aut.HasState("") {
		aut.AddState("")
	}

	for _, members := range eqClasses.classesOfEquivalence {
		for _, member := range members {
			if member.prefixVal == "ba" {
				u := 0

				_ = u
			}

			for _, letter := range alphabet {
				destState := member.prefixVal + string(letter)

				if eqClass, ok := eqClasses.wordToClass[destState]; ok {
					aut.AddTransition(
						member.prefixVal,
						eqClasses.classesOfEquivalence[eqClass][0].prefixVal,
						letter,
					)
				}
			}
		}
	}

	return aut
}
