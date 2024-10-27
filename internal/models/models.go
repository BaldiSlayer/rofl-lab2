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
