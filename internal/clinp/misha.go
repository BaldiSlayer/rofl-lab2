package clinp

import (
	"fmt"
	"github.com/BaldiSlayer/rofl-lab2/internal/mat"
	"log/slog"
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

			fmt.Println(table)
		}
	}
}
