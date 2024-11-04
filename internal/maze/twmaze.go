package maze

import (
	"fmt"
	"github.com/BaldiSlayer/rofl-lab2/pkg/bmstore"

	"github.com/BaldiSlayer/rofl-lab2/internal/cautomata"
	"github.com/BaldiSlayer/rofl-lab2/internal/defaults"
	"github.com/BaldiSlayer/rofl-lab2/pkg/models"
)

const (
	outCells = iota + 1
	inCells
)

type LightWallCell struct {
	leftState  bool
	rightState bool
	upState    bool
	downState  bool
}

// Left возвращает можно ли пойти из клетки налево
func (c *LightWallCell) Left() bool {
	return c.leftState
}

// Right возвращает можно ли пойти из клетки направо
func (c *LightWallCell) Right() bool {
	return c.rightState
}

// Up возвращает можно ли пойти из клетки вверх
func (c *LightWallCell) Up() bool {
	return c.upState
}

// Down возвращает можно ли пойти из клетки вниз
func (c *LightWallCell) Down() bool {
	return c.downState
}

type ThinWalled struct {
	maze   [][]LightWallCell
	height int
	width  int
}

func NewThinWalled(width, height int, maze [][]LightWallCell) *ThinWalled {
	return &ThinWalled{
		maze:   maze,
		width:  width,
		height: height,
	}
}

func vecToLetter(vec models.Vector) byte {
	return map[models.Vector]byte{
		models.Vector{X: 0, Y: -1}: 'N',
		models.Vector{X: 0, Y: 1}:  'S',
		models.Vector{X: -1, Y: 0}: 'W',
		models.Vector{X: 1, Y: 0}:  'E',
	}[vec]
}

// IsSpecial проверяет является ли клетка особым состоянием
// особое состояние - любая клетка вне лабиринта и его каймы в одну клетку
func (w *ThinWalled) IsSpecial(cell models.Cell) bool {
	return cell.X < -1 || cell.X > w.width || cell.Y < -1 || cell.Y > w.height
}

func (w *ThinWalled) GetDest(src models.Cell, a byte) models.Cell {
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

// symbolCanGo проверяет можно ли пройти из src по символу symbol
func (w *ThinWalled) canGoFromSrc(src models.Cell, symbol byte) bool {
	switch symbol {
	case 'N':
		return w.maze[src.Y][src.X].Up()
	case 'S':
		return w.maze[src.Y][src.X].Down()
	case 'W':
		return w.maze[src.Y][src.X].Left()
	case 'E':
		return w.maze[src.Y][src.X].Right()
	}

	return false
}

// canGoToDst проверяет можно ли попасть в dst по символу symbol
func (w *ThinWalled) canGoToDst(dst models.Cell, symbol byte) bool {
	switch symbol {
	case 'N':
		return w.maze[dst.Y][dst.X].Down()
	case 'S':
		return w.maze[dst.Y][dst.X].Up()
	case 'W':
		return w.maze[dst.Y][dst.X].Right()
	case 'E':
		return w.maze[dst.Y][dst.X].Left()
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
	dst := w.GetDest(src, a)

	if w.CanGo(src, dst) {
		return dst
	}

	return src
}

func (w *ThinWalled) IsOut(cell models.Cell) bool {
	return cell.Y < 0 || cell.Y >= len(w.maze) || cell.X < 0 || cell.X >= len(w.maze[0])
}

// IsBorder является ли клетка "каймой лабиринта"
func (w *ThinWalled) IsBorder(cell models.Cell) bool {
	return cell.X == -1 || cell.Y == -1 || cell.X == w.width || cell.Y == w.height
}

func (w *ThinWalled) Print() {
	for _, layer := range w.maze {
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
	w.maze[y1][x1].downState = true
	w.maze[y2][x2].upState = true
}

func (w *ThinWalled) DeleteHorizontalWall(x1, y1, x2, y2 int) {
	w.maze[y1][x1].rightState = true
	w.maze[y2][x2].leftState = true
}

// MakeExit создает выход в (row, col) путем удаления стены
// применяется для удаления стен в граничных точках, поэтому удаляет
// лишь одну стенку
func (w *ThinWalled) MakeExit(row, col int) {
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

// addTransitions добавляет для
func (w *ThinWalled) addTransitions(
	transitions cautomata.Transitions,
	i, j int,
) cautomata.Transitions {
	directions := defaults.GetDirections()
	alphabet := defaults.GetAlphabet()

	for idx, dir := range directions {
		src := models.Cell{X: j, Y: i}
		dst := models.Cell{X: j + dir.X, Y: i + dir.Y}

		// в спец клетку(и) можем попасть только из "каймы", там точно нет стенок
		if w.IsSpecial(dst) {
			transitions.Add(src, defaults.SpecialState(), alphabet[idx])

			continue
		}

		// если можем пройти - добавить переход в следующую клетку, иначе петля
		if w.CanGo(
			src,
			dst,
		) {
			transitions.Add(src, dst, alphabet[idx])
		} else {
			transitions.Add(src, src, alphabet[idx])
		}
	}

	return transitions
}

// mazeIterator итерируется по лабиринту и для каждой клетки применяет функцию f
func (w *ThinWalled) mazeIterator(mode int, f func(y, x int)) {
	if bmstore.Check(mode, inCells) {
		for i := 0; i < w.height; i++ {
			for j := 0; j < w.width; j++ {
				f(i, j)
			}
		}
	}

	if bmstore.Check(mode, outCells) {
		// добавим сверху
		j := -1
		for i := 0; i < w.height; i++ {
			f(i, j)
		}

		// добавим справа
		j = w.width
		for i := 0; i < w.height; i++ {
			f(i, j)
		}

		// добавим сверху
		i := -1
		for j := 0; j < w.width; j++ {
			f(i, j)
		}

		// добавим снизу
		i = w.height
		for j := 0; j < w.width; j++ {
			f(i, j)
		}

		// уголки
		f(-1, -1)
		f(-1, w.width)
		f(w.height, -1)
		f(w.height, w.width)
	}
}

// getAllStates получает все состояния для ДКА
func (w *ThinWalled) getAllStates() map[models.Cell]struct{} {
	states := make(map[models.Cell]struct{}, 1+w.width+w.height)
	states[defaults.SpecialState()] = struct{}{}

	w.mazeIterator(bmstore.Store(inCells, outCells), func(y, x int) {
		states[models.Cell{X: x, Y: y}] = struct{}{}
	})

	return states
}

// getAllStates получает все финальные состояния для ДКА
func (w *ThinWalled) getFinalStates() map[models.Cell]struct{} {
	finalStates := make(map[models.Cell]struct{})
	finalStates[defaults.SpecialState()] = struct{}{}

	w.mazeIterator(bmstore.Store(outCells), func(y, x int) {
		finalStates[models.Cell{X: x, Y: y}] = struct{}{}
	})

	return finalStates
}

// getTransitions получает переходы автомата
func (w *ThinWalled) getTransitions() cautomata.Transitions {
	transitions := cautomata.NewTransitions()

	w.mazeIterator(bmstore.Store(inCells, outCells), func(y, x int) {
		transitions = w.addTransitions(transitions, y, x)
	})

	return transitions
}

// ToDFA переводит лабиринт в детерминированный конечный автомат
func (w *ThinWalled) ToDFA() *cautomata.DFA {
	transitions := w.getTransitions()

	finalStates := w.getFinalStates()

	allStates := w.getAllStates()

	return cautomata.NewDFA(
		models.Cell{X: 0, Y: 0},
		finalStates,
		defaults.GetAlphabet(),
		transitions,
		allStates,
	)
}
