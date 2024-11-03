package clinp

import (
	"fmt"
	"github.com/BaldiSlayer/rofl-lab2/internal/defaults"
	"github.com/BaldiSlayer/rofl-lab2/internal/eqtable"
	"github.com/BaldiSlayer/rofl-lab2/internal/mat"
	"log/slog"
	"strings"
)

type Misha struct{}

func readInclude(ch chan string) string {
	return <-ch
}

func readTable(ch chan string) []string {
	tableLines := make([]string, 0)

	for message := range ch {
		if message == "end" {
			break
		}

		tableLines = append(tableLines, message)
	}

	return tableLines
}

func splitTableToParts(table []string) eqtable.TableParts {
	prefixes := make([]string, 0, len(table))

	for _, str := range table {
		// берем из каждой строки первый
		first := strings.Split(str, " ")[0]

		if first == defaults.EpsilonSymbol {
			prefixes = append(prefixes, "")

			continue
		}

		prefixes = append(prefixes, first)
	}

	suffixes := make([]string, 0)

	for _, suf := range strings.Split(table[0], " ") {
		if suf == defaults.EpsilonSymbol {
			suffixes = append(suffixes, "")

			continue
		}

		suffixes = append(suffixes, suf)
	}

	answers := make([][]bool, 0, len(table)-1)

	for _, str := range table[1:] {
		newAnsString := make([]bool, 0)

		for _, val := range strings.Split(str, " ")[1:] {
			if val == "0" {
				newAnsString = append(newAnsString, false)
			} else {
				newAnsString = append(newAnsString, true)
			}
		}

		answers = append(answers, newAnsString)
	}

	return eqtable.TableParts{
		Prefixes: prefixes,
		Suffixes: suffixes,
		Answers:  answers,
	}
}

func (m *Misha) ProcessCommands(teacher mat.MAT, commandsChan chan string) {
	for message := range commandsChan {
		if message == "g" {
			err := teacher.Generate()
			if err != nil {
				slog.Error("failed to generate", "error", err)
			}
		}

		if message == "p" {
			visualize, err := teacher.Visualize()
			if err != nil {
				slog.Error("failed to visualize", "error", err)
			}

			for _, line := range visualize {
				fmt.Println(line)
			}
		}

		if message == "isin" {
			query := readInclude(commandsChan)

			answer, _ := teacher.Include(query)
			if answer {
				fmt.Println("True")
			} else {
				fmt.Println("False")
			}
		}

		if message == "table" {
			table := readTable(commandsChan)

			tableParts := splitTableToParts(table)

			eqTable := eqtable.NewOverMaze(tableParts)

			equal, err := teacher.Equal(eqTable)
			if err != nil {
				slog.Error("failed to check equality", "error", err)
			}

			if equal.Equal {
				fmt.Println("TRUE")
			} else {
				fmt.Println(equal.CounterExample.CounterExample)
			}
		}
	}
}
