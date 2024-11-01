package app

import (
	"bufio"
	"fmt"
	"github.com/BaldiSlayer/rofl-lab2/internal/mat"
	"github.com/BaldiSlayer/rofl-lab2/internal/mazegen"
	"log/slog"
	"os"
)

type Lab2 struct {
	teacher mat.MAT
	width   int
	height  int
}

func NewLab2(width, height int) *Lab2 {
	return &Lab2{
		teacher: mat.NewRealization(
			mazegen.NewLightWallsGenerator(),
			width,
			height,
		),
	}
}

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

func (lab2 *Lab2) processCommands(ch chan string) {
	for message := range ch {
		if message == "g" {
			err := lab2.teacher.Generate()
			if err != nil {
				slog.Error("failed to generate", "error", err)
			}
		}

		if message == "p" {
			visualize, err := lab2.teacher.Visualize()
			if err != nil {
				slog.Error("failed to visualize", "error", err)
			}

			for _, line := range visualize {
				fmt.Println(line)
			}
		}

		if message == "isin" {
			query := readInclude(ch)

			answer, _ := lab2.teacher.Include(query)
			if answer {
				fmt.Println("True")
			} else {
				fmt.Println("False")
			}
		}

		if message == "table" {
			readedTable := readTable(ch)

			fmt.Println(readedTable)
		}
	}
}

func (lab2 *Lab2) cli() {
	scanner := bufio.NewScanner(os.Stdin)

	ch := make(chan string)

	go func(ch chan string) {
		lab2.processCommands(ch)
	}(ch)

	for {
		if scanner.Scan() {
			ch <- scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "error reading input:", err)

			break
		}
	}
}

func (lab2 *Lab2) httpServer() {
	// TODO impelement at some time
}

func (lab2 *Lab2) Run() {
	// TODO это нужно выбирать в зависимости от значения флага
	lab2.cli()
}
