package dsu

// DSU реализация структуры данных система непересекающихся множеств
type DSU struct {
	parent []int
	rank   []int
}

func New(n int) *DSU {
	parent := make([]int, n)
	for i := range parent {
		parent[i] = i
	}

	rank := make([]int, n)

	return &DSU{
		parent: parent,
		rank:   rank,
	}
}

func (d *DSU) Find(a int) int {
	if d.parent[a] != a {
		d.parent[a] = d.Find(d.parent[a])
	}

	return d.parent[a]
}

func (d *DSU) Union(a int, b int) {
	rootA := d.Find(a)
	rootB := d.Find(b)

	if rootA != rootB {
		if d.rank[rootA] > d.rank[rootB] {
			d.parent[rootB] = rootA
		} else if d.rank[rootA] < d.rank[rootB] {
			d.parent[rootA] = rootB
		} else {
			d.parent[rootB] = rootA
			d.rank[rootA] += 1
		}
	}
}
