package models

type CounterExample struct {
	CounterExample string
}

// EqualResponse - ответ на запрос эквивалентности
type EqualResponse struct {
	Equal          bool
	CounterExample CounterExample
}

type Cell struct {
	X int
	Y int
}

type Transition struct {
	// Src - откуда идет переход
	Src int
	// Dst - куда идет переход
	Dst int
	// Symbol - по какому символу идет переход (хотя мб было бы лучше хранить тут рну)
	Symbol string
}

// Vector вектор перемещения
type Vector struct {
	X int
	Y int
}
