package app

import (
	"bufio"
	"fmt"
	"github.com/BaldiSlayer/rofl-lab2/internal/clinp"
	"github.com/BaldiSlayer/rofl-lab2/internal/mat"
	"github.com/BaldiSlayer/rofl-lab2/internal/mazegen"
	"os"
)

type Lab2 struct {
	teacher        mat.MAT
	inputProcessor clinp.InputProcessor
	width          int
	height         int
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

func (lab2 *Lab2) cli() {
	scanner := bufio.NewScanner(os.Stdin)

	ch := make(chan string)

	go func() {
		lab2.inputProcessor.ProcessCommands(
			lab2.teacher,
			ch,
		)
	}()

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
	lab2.inputProcessor = &clinp.Misha{}

	// TODO это нужно выбирать в зависимости от значения флага
	lab2.cli()
}
