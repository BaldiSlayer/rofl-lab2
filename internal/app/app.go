package app

import (
	"bufio"
	"fmt"
	"github.com/BaldiSlayer/rofl-lab2/internal/mat"
	"github.com/BaldiSlayer/rofl-lab2/internal/mazegen"
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
	Action  func(s string)
}

func (lab2 *Lab2) addCliHandlers() []Handler {
	cmdGen := Handler{
		Checker: func(s string) bool {
			return s == "g"
		},
		Action: func(s string) {
			err := lab2.teacher.Generate()
			if err != nil {
				fmt.Println(err)
			}
		},
	}

	cmdPrint := Handler{
		Checker: func(s string) bool {
			return s == "p"
		},
		Action: func(s string) {
			maze, err := lab2.teacher.Print()
			if err != nil {
				fmt.Println(err)
			}

			for _, line := range maze {
				fmt.Println(line)
			}
		},
	}

	cmdInclude := Handler{
		Checker: func(s string) bool {
			return strings.HasPrefix(s, "i ")
		},
		Action: func(s string) {
			fmt.Println(lab2.teacher.Include(s[2:]))
		},
	}

	cmdEqual := Handler{
		Checker: func(s string) bool {
			return strings.HasPrefix(s, "e ")
		},
		Action: func(s string) {
			_ = s[2:]

			fmt.Println(lab2.teacher.Equal())
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
					cmd.Action(command)
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
	lab2.cli()
}
