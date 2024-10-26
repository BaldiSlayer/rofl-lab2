package mazegen

import (
	"fmt"
	"github.com/BaldiSlayer/rofl-lab2/internal/dsu"
	"math/rand"
	"time"
)

type wall struct {
	x1 int
	x2 int
	y1 int
	y2 int
}

type lightWallCell struct {
	leftState  bool
	rightState bool
	upState    bool
	downState  bool
}

func (c *lightWallCell) left() bool {
	return c.leftState
}

func (c *lightWallCell) right() bool {
	return c.rightState
}

func (c *lightWallCell) up() bool {
	return c.upState
}

func (c *lightWallCell) down() bool {
	return c.downState
}

type LightWalls struct {
	maze [][]lightWallCell
}

func (w *LightWalls) Print() {
	for _, layer := range w.maze {
		rowStr := ""

		for _, cell := range layer {
			if cell.right() {
				rowStr += " "
			} else {
				rowStr += "|"
			}
		}

		fmt.Println(rowStr)

		bottomStr := ""

		for _, cell := range layer {
			if cell.down() {
				bottomStr += " "
			} else {
				bottomStr += "-"
			}
		}

		fmt.Println(bottomStr)
	}
}

type LightWallsGenerator struct{}

func NewLightWallsGenerator() *LightWallsGenerator {
	return &LightWallsGenerator{}
}

func (l *LightWallsGenerator) Generate(width, height int) (*LightWalls, error) {
	walls := make([]wall, 0)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if x < width-1 {
				walls = append(walls, wall{
					x1: x,
					y1: y,
					x2: x + 1,
					y2: y,
				})
			}

			if y < height-1 {
				walls = append(walls, wall{
					x1: x,
					y1: y,
					x2: x,
					y2: y + 1,
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

	maze := make([][]lightWallCell, height)
	for i := 0; i < width; i++ {
		maze[i] = make([]lightWallCell, width)
	}

	for _, wallInst := range walls {
		cell1 := wallInst.y1*width + wallInst.x1
		cell2 := wallInst.y2*width + wallInst.x2

		if disDU.Find(cell1) != disDU.Find(cell2) {
			disDU.Union(cell1, cell2)

			if wallInst.x1 == wallInst.x2 {
				if wallInst.y1 < wallInst.y2 {
					maze[wallInst.y1][wallInst.x1].downState = true
					maze[wallInst.y2][wallInst.x2].upState = true
				}

				if wallInst.y1 < wallInst.y2 {
					maze[wallInst.y1][wallInst.x1].rightState = true
					maze[wallInst.y2][wallInst.x2].leftState = true
				}
			}
		}

	}

	return &LightWalls{
		maze: maze,
	}, nil
}
