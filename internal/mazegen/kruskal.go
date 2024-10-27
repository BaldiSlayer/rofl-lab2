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

// TODO заменить на хранение информации в битах числа, сейчас мне лень
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

type ThinWalledMaze struct {
	maze [][]lightWallCell
}

func (w *ThinWalledMaze) getDest(i, j int, a byte) (int, int) {
	switch a {
	case 'N':
		return i - 1, j
	case 'S':
		return i + 1, j
	case 'W':
		return i, j - 1
	case 'E':
		return i, j + 1
	}

	return i, j
}

// TODO refactor
func (w *ThinWalledMaze) GetPosAfterStep(i, j int, a byte) (int, int) {
	destI, destJ := w.getDest(i, j, a)

	// не хочу получать fatalpanic
	if w.IsOut(i, j) {
		if w.IsOut(destI, destJ) {
			return destI, destJ
		}

		// проверка "запустит ли" меня клетка лабиринта
		switch a {
		case 'N':
			if w.maze[destI][destJ].down() {
				return destI, destJ
			}
			return i, j
		case 'S':
			if w.maze[destI][destJ].up() {
				return destI, destJ
			}
			return i, j
		case 'W':
			if w.maze[destI][destJ].right() {
				return destI, destJ
			}
			return i, j
		case 'E':
			if w.maze[destI][destJ].left() {
				return destI, destJ
			}
			return i, j
		}
	}

	switch a {
	case 'N':
		if w.maze[i][j].up() {
			return destI, destJ
		}
		return i, j
	case 'S':
		if w.maze[i][j].down() {
			return destI, destJ
		}
		return i, j
	case 'W':
		if w.maze[i][j].left() {
			return destI, destJ
		}
		return i, j
	case 'E':
		if w.maze[i][j].right() {
			return destI, destJ
		}
		return i, j
	}

	return i, j
}

func (w *ThinWalledMaze) IsOut(i, j int) bool {
	return i < 0 || i >= len(w.maze) || j < 0 || j >= len(w.maze[0])
}

func (w *ThinWalledMaze) Print() {
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

func (w *ThinWalledMaze) makeVerticalWall(x1, y1, x2, y2 int) {
	w.maze[y1][x1].downState = true
	w.maze[y2][x2].upState = true
}

func (w *ThinWalledMaze) makeHorizontalWall(x1, y1, x2, y2 int) {
	w.maze[y1][x1].rightState = true
	w.maze[y2][x2].leftState = true
}

type LightWallsGenerator struct{}

func NewLightWallsGenerator() *LightWallsGenerator {
	return &LightWallsGenerator{}
}

func (l *LightWallsGenerator) Generate(width, height int) (*ThinWalledMaze, error) {
	mazeField := make([][]lightWallCell, height)
	for i := 0; i < width; i++ {
		mazeField[i] = make([]lightWallCell, width)
	}

	generatedMaze := ThinWalledMaze{
		maze: mazeField,
	}

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

	for _, wallInst := range walls {
		cell1 := wallInst.y1*width + wallInst.x1
		cell2 := wallInst.y2*width + wallInst.x2

		if disDU.Find(cell1) != disDU.Find(cell2) {
			disDU.Union(cell1, cell2)

			if wallInst.x1 == wallInst.x2 {
				if wallInst.y1 < wallInst.y2 {
					generatedMaze.makeVerticalWall(
						wallInst.x1,
						wallInst.y1,
						wallInst.x2,
						wallInst.y2,
					)
				}
			} else {
				if wallInst.x1 < wallInst.x2 {
					generatedMaze.makeHorizontalWall(
						wallInst.x1,
						wallInst.y1,
						wallInst.x2,
						wallInst.y2,
					)
				}
			}
		}
	}

	generatedMaze.makeExit(0, rand.Intn(width))
	generatedMaze.makeExit(height-1, rand.Intn(width))

	generatedMaze.makeExit(rand.Intn(height), 0)
	generatedMaze.makeExit(rand.Intn(height), width-1)

	return &generatedMaze, nil
}

func (w *ThinWalledMaze) makeExit(row, col int) {
	if row == 0 {
		w.maze[row][col].upState = true
	}

	if row == len(w.maze)-1 {
		w.maze[row][col].downState = true
	}

	if col == 0 {
		w.maze[row][col].leftState = true
	}

	if col == len(w.maze[0])-1 {
		w.maze[row][col].rightState = true
	}
}
