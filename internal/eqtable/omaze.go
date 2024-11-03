package eqtable

import (
	"github.com/BaldiSlayer/rofl-lab2/internal/automata"
	"github.com/BaldiSlayer/rofl-lab2/internal/maze"
	"github.com/BaldiSlayer/rofl-lab2/pkg/models"
)

// OverMaze - таблица, которую мне передают при запросе на эквивалентность
type OverMaze struct {
	prefixes []string
	suffixes []string
	// [индекс_префикса][индекс_суффикса]
	answers [][]bool
}

func NewOverMaze(parts TableParts) *OverMaze {
	return &OverMaze{
		prefixes: parts.Prefixes,
		suffixes: parts.Suffixes,
		answers:  parts.Answers,
	}
}

// GetWords -
func (table *OverMaze) getWords() map[string]bool {
	mp := make(map[string]bool)

	for i, pref := range table.prefixes {
		for j, suf := range table.suffixes {
			mp[pref+suf] = table.answers[i][j]
		}
	}

	return mp
}

func wordIterate(
	startState models.Cell,
	word string,
	maze *maze.ThinWalled,
	aut *automata.DFA,
) models.Cell {
	curState := startState

	for _, letter := range word {
		nextState := maze.GetPosAfterStep(curState, byte(letter))

		// если мы в специальном состоянии, то просто перемещаемся дальше
		if maze.IsSpecial(curState) {
			// если мы вышли из специального состояния, то передвинемся
			if !maze.IsSpecial(nextState) {
				curState = nextState
			}

			continue
		}

		canGo := maze.CanGo(curState, nextState)

		// если не можем пойти, то оставляем тот же стейт и добавляем петлю
		if !canGo {
			aut.AddTransition(curState, curState, byte(letter))

			continue
		}

		// если то, куда мы переходим - специальное состояние, то нужно задать соответствующее значение
		if maze.IsSpecial(nextState) {
			nextState = automata.SpecialState()
		}

		aut.AddTransition(curState, nextState, byte(letter))
		aut.AddState(nextState)

		// передвигаемся
		curState = nextState
	}

	return curState
}

func (table *OverMaze) ToDFA(maze *maze.ThinWalled) *automata.DFA {
	startState := models.Cell{X: 0, Y: 0}

	aut := automata.NewEmptyDFA()

	// обходим полученную таблицу
	for word, included := range table.getWords() {
		// стоим в начале
		stateAfterWalk := wordIterate(startState, word, maze, aut)

		// если на пересечении стояла единичка, добавляем в финальные состояния
		if included {
			aut.AddFinalState(stateAfterWalk)
		}
	}

	return aut
}
