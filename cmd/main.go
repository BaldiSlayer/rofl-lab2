package main

import (
	"github.com/BaldiSlayer/rofl-lab2/internal/app"
)

const (
	width  = 1
	height = 1
)

func main() {
	mat := app.NewLab2(width, height)

	mat.Run()
}
