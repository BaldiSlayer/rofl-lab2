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

func (c *LightWallCell) Left() bool {
	return c.leftState
}

func (c *LightWallCell) Right() bool {
	return c.rightState
}

func (c *LightWallCell) Up() bool {
	return c.upState
}

func (c *LightWallCell) Down() bool {
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

func (w *ThinWalled) vecToLetter(vec models.Vector) string {
	return map[models.Vector]string{
		models.Vector{X: 0, Y: -1}: "N",
		models.Vector{X: 0, Y: 1}:  "S",
		models.Vector{X: -1, Y: 0}: "W",
		models.Vector{X: 1, Y: 0}:  "E",
	}[vec]
}

// symbolCanGo проверяет можно ли пройти из src по символу symbol
func (w *ThinWalled) canGoFromSrc(src models.Cell, symbol string) bool {
	switch symbol {
	case "N":
		return w.Maze[src.Y][src.X].Up()
	case "S":
		return w.Maze[src.Y][src.X].Down()
	case "W":
		return w.Maze[src.Y][src.X].Left()
	case "E":
		return w.Maze[src.Y][src.X].Right()
	}

	return false
}

// canGoToDst проверяет можно ли попасть в dst по символу symbol
func (w *ThinWalled) canGoToDst(dst models.Cell, symbol string) bool {
	switch symbol {
	case "N":
		return w.Maze[dst.Y][dst.X].Down()
	case "S":
		return w.Maze[dst.Y][dst.X].Up()
	case "W":
		return w.Maze[dst.Y][dst.X].Right()
	case "E":
		return w.Maze[dst.Y][dst.X].Left()
	}

	return false
}

// CanGo - проверяет можем ли мы пройти от src до dst
func (w *ThinWalled) CanGo(src, dst models.Cell) bool {
	symbol := w.vecToLetter(models.Vector{
		X: dst.X - src.X,
		Y: dst.Y - dst.Y,
	})

	// если обе вне - никто не мешает
	if w.IsOut(src.Y, src.X) && w.IsOut(dst.Y, dst.X) {
		return true
	}

	// если клетка вне лабиринта, то нужно посмотреть обратную стену
	if w.IsOut(src.Y, src.X) {
		return w.canGoToDst(dst, symbol)
	}

	return w.canGoFromSrc(src, symbol)
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
			if w.Maze[destI][destJ].Down() {
				return destI, destJ
			}
			return i, j
		case 'S':
			if w.Maze[destI][destJ].Up() {
				return destI, destJ
			}
			return i, j
		case 'W':
			if w.Maze[destI][destJ].Right() {
				return destI, destJ
			}
			return i, j
		case 'E':
			if w.Maze[destI][destJ].Left() {
				return destI, destJ
			}
			return i, j
		}
	}

	switch a {
	case 'N':
		if w.Maze[i][j].Up() {
			return destI, destJ
		}
		return i, j
	case 'S':
		if w.Maze[i][j].Down() {
			return destI, destJ
		}
		return i, j
	case 'W':
		if w.Maze[i][j].Left() {
			return destI, destJ
		}
		return i, j
	case 'E':
		if w.Maze[i][j].Right() {
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
			if cell.Right() {
				rowStr += " "
			} else {
				rowStr += "|"
			}
		}

		fmt.Println(rowStr)

		bottomStr := ""

		for _, cell := range layer {
			if cell.Down() {
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
			w.Maze[dest.Y][dest.X].Down,
			w.Maze[dest.Y][dest.X].Up,
			w.Maze[dest.Y][dest.X].Right,
			w.Maze[dest.Y][dest.X].Left,
		}

		return !isWallDirs[actionNumber]()
	}

	isWallDirs := [...]func() bool{
		w.Maze[src.Y][src.X].Up,
		w.Maze[src.Y][src.X].Down,
		w.Maze[src.Y][src.X].Left,
		w.Maze[src.Y][src.X].Right,
	}

	return !isWallDirs[actionNumber]()
}

// GetPath находит путь от start до end
func (w *ThinWalled) GetPath(start, end models.Cell) string {
	// Up, Down, Left, Right, это важно, иначе сломается isWallInWay
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
