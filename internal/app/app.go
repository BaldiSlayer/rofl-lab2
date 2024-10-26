package app

import (
	"bufio"
	"fmt"
	"github.com/BaldiSlayer/rofl-lab2/internal/mat"
	"github.com/BaldiSlayer/rofl-lab2/internal/mazegen"
	"log/slog"
	"os"
	"strings"
)

type Lab2 struct {
	teacher mat.MAT
	width   int
	height  int
}

func NewLab2(width, height int) *Lab2 {
	return &Lab2{
		teacher: mat.NewRealization(
			mazegen.NewSimpleGenerator(),
			width,
			height,
		),
	}
}

type Handler struct {
	Checker func(s string) bool
	Action  func(s string) error
}

func (lab2 *Lab2) addCliHandlers() []Handler {
	cmdGen := Handler{
		Checker: func(s string) bool {
			return s == "g"
		},
		Action: func(s string) error {
			err := lab2.teacher.Generate()
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmdPrint := Handler{
		Checker: func(s string) bool {
			return s == "p"
		},
		Action: func(s string) error {
			maze, err := lab2.teacher.Print()
			if err != nil {
				return err
			}

			for _, line := range maze {
				fmt.Println(line)
			}

			return nil
		},
	}

	cmdInclude := Handler{
		Checker: func(s string) bool {
			return strings.HasPrefix(s, "i ")
		},
		Action: func(s string) error {
			res, err := lab2.teacher.Include(s[2:])
			if err != nil {
				return err
			}

			fmt.Println(res)

			return nil
		},
	}

	cmdEqual := Handler{
		Checker: func(s string) bool {
			return strings.HasPrefix(s, "e ")
		},
		Action: func(s string) error {
			_ = s[2:]

			res, err := lab2.teacher.Equal()
			if err != nil {
				return err
			}

			fmt.Println(res)

			return nil
		},
	}

	return []Handler{
		cmdGen,
		cmdPrint,
		cmdInclude,
		cmdEqual,
	}
}

func (lab2 *Lab2) cli() {
	commands := lab2.addCliHandlers()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		if scanner.Scan() {
			command := scanner.Text()

			for _, cmd := range commands {
				if cmd.Checker(command) {
					err := cmd.Action(command)
					if err != nil {
						slog.Error("error while do action", "error", err)
					}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
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
