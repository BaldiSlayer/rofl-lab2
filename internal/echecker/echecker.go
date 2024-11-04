package echecker

// EqualityChecker делает проверку на эквивалентность автомата полученного из лабиринта и из таблицы
type EqualityChecker interface {
	GetCounterExample() string
}
