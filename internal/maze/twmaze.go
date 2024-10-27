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
// true - можно пройти
type LightWallCell struct {
	leftState  bool
	rightState bool
	upState    bool
	downState  bool
}

// left
func (c *LightWallCell) left() bool {
	return c.leftState
}

// right
func (c *LightWallCell) right() bool {
	return c.rightState
}

// up
func (c *LightWallCell) up() bool {
	return c.upState
}

// down
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

func (w *ThinWalled) DeleteVerticalWall(x1, y1, x2, y2 int) {
	w.Maze[y1][x1].downState = true
	w.Maze[y2][x2].upState = true
}

func (w *ThinWalled) DeleteHorizontalWall(x1, y1, x2, y2 int) {
	w.Maze[y1][x1].rightState = true
	w.Maze[y2][x2].leftState = true
}

func (w *ThinWalled) cellCordsToInt(x, y int) int {
	return y*len(w.Maze[0]) + x
}

// isWallInWay проверяет, будет ли стена на пути (если да, то мы не можем пойти)
func (w *ThinWalled) isWallInWay(src, dest models.Cell, actionNumber int) bool {
	// если обе клетки вне - нам ничто не мешает
	if w.IsOut(src.Y, src.X) && w.IsOut(dest.Y, dest.X) {
		return false
	}

	if w.IsOut(src.Y, src.X) {
		isWallDirs := [...]func() bool{
			w.Maze[dest.Y][dest.X].down,
			w.Maze[dest.Y][dest.X].up,
			w.Maze[dest.Y][dest.X].right,
			w.Maze[dest.Y][dest.X].left,
		}

		return !isWallDirs[actionNumber]()
	}

	isWallDirs := [...]func() bool{
		w.Maze[src.Y][src.X].up,
		w.Maze[src.Y][src.X].down,
		w.Maze[src.Y][src.X].left,
		w.Maze[src.Y][src.X].right,
	}

	return !isWallDirs[actionNumber]()
}

// GetPath находит путь от start до end
func (w *ThinWalled) GetPath(start, end models.Cell) string {
	// up, down, left, right, это важно, иначе сломается isWallInWay
	directions := []models.Cell{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}

	queue := list.New()
	queue.PushBack(start)

	visited := make(map[models.Cell]struct{})
	visited[start] = struct{}{}

	prev := make(map[models.Cell]models.Cell)

	// делаем вид, что тут пришли сами в себя
	prev[start] = start

	for queue.Len() > 0 {
		current := queue.Front().Value.(models.Cell)
		queue.Remove(queue.Front())

		if current == end {
			break
		}

		for number, dir := range directions {
			next := models.Cell{
				X: current.X + dir.X,
				Y: current.Y + dir.Y,
			}

			_, vis := visited[next]

			if !w.isWallInWay(current, next, number) && !vis {
				visited[next] = struct{}{}

				queue.PushBack(next)
				prev[next] = current
			}
		}
	}

	// восстановление пути
	path := ""

	cur := end
	for cur != start {
		parent := prev[cur]

		step := models.Cell{
			X: cur.X - parent.X,
			Y: cur.Y - parent.Y,
		}

		down := models.Cell{
			X: 0,
			Y: 1,
		}

		up := models.Cell{
			X: 0,
			Y: -1,
		}

		left := models.Cell{
			X: -1,
			Y: 0,
		}

		right := models.Cell{
			X: 1,
			Y: 0,
		}

		if step == down {
			path = "S" + path
		}

		if step == up {
			path = "N" + path
		}

		if step == left {
			path = "W" + path
		}

		if step == right {
			path = "E" + path
		}

		cur = parent
	}

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
