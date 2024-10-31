package eqtable

import (
	"github.com/BaldiSlayer/rofl-lab2/internal/automata"
	"github.com/BaldiSlayer/rofl-lab2/internal/defaults"
	"github.com/BaldiSlayer/rofl-lab2/internal/maze"
	"github.com/BaldiSlayer/rofl-lab2/pkg/models"
)

type EqTable interface {
	ToDFA(maze *maze.ThinWalled) *automata.DFA
}

// OverMaze - таблица, которую мне передают при запросе на эквивалентность
type OverMaze struct {
	prefixes []string
	suffixes []string
	// [индекс_префикса][индекс_суффикса]
	answers [][]bool
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

func (table *OverMaze) wordIterate(startState models.Cell, word string, maze *maze.ThinWalled) models.Cell {
	transitions := automata.NewTransitions()

	states := make(map[models.Cell]struct{})
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
			transitions.Add(curState, nextState, byte(letter))

			continue
		}

		// если то, куда мы переходим - специальное состояние, то нужно задать соответствующее значение
		if maze.IsSpecial(nextState) {
			nextState = automata.SpecialState()
		}

		transitions.Add(curState, nextState, byte(letter))

		states[nextState] = struct{}{}

		// передвинулись
		curState = nextState
	}

	return curState
}

func (table *OverMaze) ToDFA(maze *maze.ThinWalled) *automata.DFA {
	startState := models.Cell{X: 0, Y: 0}

	// инициализация финальных состояний
	finalStates := make(map[models.Cell]struct{})
	// инициализация переходов
	transitions := automata.NewTransitions()
	// инициализация
	states := make([]models.Cell, 0)

	// обходим полученную таблицу
	for word, included := range table.getWords() {
		// стоим в начале
		stateAfterWalk := table.wordIterate(startState, word, maze)

		// если на пересечении стояла единичка, добавляем в финальные состояния
		if included {
			finalStates[stateAfterWalk] = struct{}{}
		}
	}

	return automata.NewDFA(
		startState,
		finalStates,
		defaults.GetAlphabet(),
		transitions,
		states,
	)
}
