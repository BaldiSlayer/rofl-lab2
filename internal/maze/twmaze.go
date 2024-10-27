package maze

import (
	"container/list"
	"fmt"
	"github.com/BaldiSlayer/rofl-lab2/internal/models"
)

type Wall struct {
	X1 int
	X2 int
	Y1 int
	Y2 int
}

// TODO заменить на хранение информации в битах числа, сейчас мне лень
type LightWallCell struct {
	leftState  bool
	rightState bool
	upState    bool
	downState  bool
}

// left проверяет есть ли стена слева
func (c *LightWallCell) left() bool {
	return c.leftState
}

// right проверяет есть ли стена справа
func (c *LightWallCell) right() bool {
	return c.rightState
}

// up проверяет есть ли стена сверху
func (c *LightWallCell) up() bool {
	return c.upState
}

// down проверяет есть ли стена снизу
func (c *LightWallCell) down() bool {
	return c.downState
}

type ThinWalled struct {
	Maze [][]LightWallCell
}

func (w *ThinWalled) getDest(i, j int, a byte) (int, int) {
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
func (w *ThinWalled) GetPosAfterStep(i, j int, a byte) (int, int) {
	destI, destJ := w.getDest(i, j, a)

	// не хочу получать fatalpanic
	if w.IsOut(i, j) {
		if w.IsOut(destI, destJ) {
			return destI, destJ
		}

		// проверка "запустит ли" меня клетка лабиринта
		switch a {
		case 'N':
			if w.Maze[destI][destJ].down() {
				return destI, destJ
			}
			return i, j
		case 'S':
			if w.Maze[destI][destJ].up() {
				return destI, destJ
			}
			return i, j
		case 'W':
			if w.Maze[destI][destJ].right() {
				return destI, destJ
			}
			return i, j
		case 'E':
			if w.Maze[destI][destJ].left() {
				return destI, destJ
			}
			return i, j
		}
	}

	switch a {
	case 'N':
		if w.Maze[i][j].up() {
			return destI, destJ
		}
		return i, j
	case 'S':
		if w.Maze[i][j].down() {
			return destI, destJ
		}
		return i, j
	case 'W':
		if w.Maze[i][j].left() {
			return destI, destJ
		}
		return i, j
	case 'E':
		if w.Maze[i][j].right() {
			return destI, destJ
		}
		return i, j
	}

	return i, j
}

func (w *ThinWalled) IsOut(i, j int) bool {
	return i < 0 || i >= len(w.Maze) || j < 0 || j >= len(w.Maze[0])
}

func (w *ThinWalled) Print() {
	for _, layer := range w.Maze {
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

func (w *ThinWalled) MakeVerticalWall(x1, y1, x2, y2 int) {
	w.Maze[y1][x1].downState = true
	w.Maze[y2][x2].upState = true
}

func (w *ThinWalled) MakeHorizontalWall(x1, y1, x2, y2 int) {
	w.Maze[y1][x1].rightState = true
	w.Maze[y2][x2].leftState = true
}

func (w *ThinWalled) cellCordsToInt(x, y int) int {
	return y*len(w.Maze[0]) + x
}

// isWallInWay проверяет, будет ли стена на пути (если да, то мы не можем пойти)
func (w *ThinWalled) isWallInWay(dest models.Cell, actionNumber int) bool {
	// если клетка вне - париться не надо, нам никто не мешает
	if w.IsOut(dest.Y, dest.X) {
		return false
	}

	// логика такова:
	// если мы хотим пойти наверх, то нам нужно посмотреть, есть ли у верхней клетки стена снизу
	// если вниз, посмотреть, есть ли у нижней клетки стена сверху и т.д.
	// я делаю проверку только по клетке, в которую перехожу, так как стены в лабиринте двусторонние
	// и можно обойтись одной проверкой, так как мы можем быть изначально вне лабиринта и пытаться в него зайти
	isWallDirs := [...]func() bool{
		w.Maze[dest.Y][dest.X].down,
		w.Maze[dest.Y][dest.X].up,
		w.Maze[dest.Y][dest.X].right,
		w.Maze[dest.Y][dest.X].left,
	}

	return isWallDirs[actionNumber]()
}

// GetPath находит путь от start до end
func (w *ThinWalled) GetPath(start, end models.Cell) string {
	// height, width := len(w.Maze), len(w.Maze[0])

	// up, down, left, right
	directions := []models.Cell{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	queue := list.New()
	queue.PushBack(start)

	visited := make(map[models.Cell]struct{})
	visited[start] = struct{}{}

	prevs := make(map[models.Cell]models.Cell)

	// делаем вид, что тут пришли сами в себя
	prevs[start] = start

	for queue.Len() > 0 {
		current := queue.Front().Value.(models.Cell)
		queue.Remove(queue.Front())

		for number, dir := range directions {
			next := models.Cell{
				X: current.X + dir.X,
				Y: current.Y + dir.Y,
			}

			// мы должны уметь ходить в эту клетку
			if w.isWallInWay(next, number) {

			}
		}
	}

	// восстановление пути
	path := ""

	return path
}

func (w *ThinWalled) MakeExit(row, col int) {
	if row == 0 {
		w.Maze[row][col].upState = true
	}

	if row == len(w.Maze)-1 {
		w.Maze[row][col].downState = true
	}

	if col == 0 {
		w.Maze[row][col].leftState = true
	}

	if col == len(w.Maze[0])-1 {
		w.Maze[row][col].rightState = true
	}
}
