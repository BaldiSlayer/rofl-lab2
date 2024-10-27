package models

type CounterExample struct {
	CounterExample string
}

type EqualResponse struct {
	Equal          bool
	CounterExample CounterExample
}
