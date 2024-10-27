package mazegen

import (
	"github.com/BaldiSlayer/rofl-lab2/internal/dsu"
	"github.com/BaldiSlayer/rofl-lab2/internal/maze"
	"math/rand"
	"time"
)

type LightWallsGenerator struct{}

func NewLightWallsGenerator() *LightWallsGenerator {
	return &LightWallsGenerator{}
}

func (l *LightWallsGenerator) Generate(width, height int) (*maze.ThinWalled, error) {
	mazeField := make([][]maze.LightWallCell, height)
	for i := 0; i < width; i++ {
		mazeField[i] = make([]maze.LightWallCell, width)
	}

	generatedMaze := maze.ThinWalled{
		Maze: mazeField,
	}

	walls := make([]maze.Wall, 0)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if x < width-1 {
				walls = append(walls, maze.Wall{
					X1: x,
					Y1: y,
					X2: x + 1,
					Y2: y,
				})
			}

			if y < height-1 {
				walls = append(walls, maze.Wall{
					X1: x,
					Y1: y,
					X2: x,
					Y2: y + 1,
				})
			}
		}
	}

	rand.Seed(time.Now().UnixNano())

	// Используем Shuffle для перемешивания среза
	rand.Shuffle(len(walls), func(i, j int) {
		walls[i], walls[j] = walls[j], walls[i]
	})

	// TODO return error if > int max
	disDU := dsu.New(width * height)

	for _, wallInst := range walls {
		cell1 := wallInst.Y1*width + wallInst.X1
		cell2 := wallInst.Y2*width + wallInst.X2

		if disDU.Find(cell1) != disDU.Find(cell2) {
			disDU.Union(cell1, cell2)

			if wallInst.X1 == wallInst.X2 {
				if wallInst.Y1 < wallInst.Y2 {
					generatedMaze.MakeVerticalWall(
						wallInst.X1,
						wallInst.Y1,
						wallInst.X2,
						wallInst.Y2,
					)
				}
			} else {
				if wallInst.X1 < wallInst.X2 {
					generatedMaze.MakeVerticalWall(
						wallInst.X1,
						wallInst.Y1,
						wallInst.X2,
						wallInst.Y2,
					)
				}
			}
		}
	}

	generatedMaze.MakeExit(0, rand.Intn(width))
	generatedMaze.MakeExit(height-1, rand.Intn(width))

	generatedMaze.MakeExit(rand.Intn(height), 0)
	generatedMaze.MakeExit(rand.Intn(height), width-1)

	return &generatedMaze, nil
}
