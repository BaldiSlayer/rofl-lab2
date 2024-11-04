package models

type CounterExample struct {
	CounterExample string
}

// EqualResponse - ответ на запрос эквивалентности
type EqualResponse struct {
	Equal          bool
	CounterExample CounterExample
}

// Cell - клетка лабиринта
type Cell struct {
	X int
	Y int
}

// Vector вектор перемещения
type Vector struct {
	X int
	Y int
}
