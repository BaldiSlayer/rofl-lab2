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

func (w *ThinWalled) getDest(src models.Cell, a byte) models.Cell {
	switch a {
	case 'N':
		return models.Cell{X: src.X, Y: src.Y - 1}
	case 'S':
		return models.Cell{X: src.X, Y: src.Y + 1}
	case 'W':
		return models.Cell{X: src.X - 1, Y: src.Y}
	case 'E':
		return models.Cell{X: src.X + 1, Y: src.Y}
	}

	return src
}

func vecToLetter(vec models.Vector) string {
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
	symbol := vecToLetter(models.Vector{
		X: dst.X - src.X,
		Y: dst.Y - src.Y,
	})

	// если обе вне - никто не мешает
	if w.IsOut(src) && w.IsOut(dst) {
		return true
	}

	// если клетка вне лабиринта, то нужно посмотреть обратную стену
	if w.IsOut(src) {
		return w.canGoToDst(dst, symbol)
	}

	return w.canGoFromSrc(src, symbol)
}

// GetPosAfterStep получает позицию, которая будет после прохода по символу a из src
func (w *ThinWalled) GetPosAfterStep(src models.Cell, a byte) models.Cell {
	dst := w.getDest(src, a)

	if w.CanGo(src, dst) {
		return dst
	}

	return src
}

func (w *ThinWalled) IsOut(cell models.Cell) bool {
	return cell.Y < 0 || cell.Y >= len(w.Maze) || cell.X < 0 || cell.X >= len(w.Maze[0])
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

// restorePath занимается восстановлением пути от start до end по даннным о предках prev
func restorePath(start, end models.Cell, prev map[models.Cell]models.Cell) string {
	path := ""

	cur := end
	for cur != start {
		path = vecToLetter(models.Vector{
			X: cur.X - prev[cur].X,
			Y: cur.Y - prev[cur].Y,
		}) + path

		cur = prev[cur]
	}

	return path
}

// bfsOnMaze реализует алгоритм поиска в ширину, возвращает обратно данные о предках для каждой клетки
func (w *ThinWalled) bfsOnMaze(start, end models.Cell) map[models.Cell]models.Cell {
	// Up, Down, Left, Right, это важно, иначе сломается isWallInWay
	directions := []models.Vector{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}

	queue := list.New()
	queue.PushBack(start)

	visited := make(map[models.Cell]struct{})
	visited[start] = struct{}{}

	prev := make(map[models.Cell]models.Cell)
	prev[start] = start

	for queue.Len() > 0 {
		current := queue.Front().Value.(models.Cell)
		queue.Remove(queue.Front())

		if current == end {
			break
		}

		for _, dir := range directions {
			next := models.Cell{
				X: current.X + dir.X,
				Y: current.Y + dir.Y,
			}

			// если можем перейти в клетку и там еще не были
			if _, vis := visited[next]; w.CanGo(current, next) && !vis {
				visited[next] = struct{}{}

				queue.PushBack(next)
				prev[next] = current
			}
		}
	}

	return prev
}

// GetPath находит путь от start до end
func (w *ThinWalled) GetPath(start, end models.Cell) string {
	return restorePath(
		start,
		end,
		w.bfsOnMaze(start, end),
	)
}

// MakeExit создает выход в (row, col) путем удаления стены
// применяется для удаления стен в граничных точках, поэтому удаляет
// лишь одну стенку
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
