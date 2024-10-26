package main

import (
	"github.com/BaldiSlayer/rofl-lab2/internal/app"
	"github.com/BaldiSlayer/rofl-lab2/internal/mazegen"
	"os"
)

const (
	width  = 3
	height = 3
)

func main() {

	mzG := mazegen.NewLightWallsGenerator()

	mz, _ := mzG.Generate(5, 5)

	mz.Print()

	os.Exit(1)

	mat := app.NewLab2(3, 3)

	mat.Run()
}
