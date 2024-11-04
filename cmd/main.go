package main

import (
	"github.com/BaldiSlayer/rofl-lab2/internal/app"
	"github.com/BaldiSlayer/rofl-lab2/internal/clinp"
)

const (
	width  = 2
	height = 2
)

func main() {
	mat := app.NewLab2(width, height)

	mat.Run(
		&clinp.Misha{},
	)
}
