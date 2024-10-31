package bmstore

func Store(a ...int) int {
	mask := 0

	for _, i := range a {
		mask |= 1 << i
	}

	return mask
}

func Check(mask, val int) bool {
	return (mask & (1 << val)) != 0
}
