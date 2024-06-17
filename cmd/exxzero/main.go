package main

import (
	"fmt"
	"sort"
)

type edge struct {
	weight, from, to int
}

type DSU struct {
	parent, rank []int
}

func newDSU(n int) *DSU {
	parent := make([]int, n)
	rank := make([]int, n)
	for i := range parent {
		parent[i] = -1
		rank[i] = 1
	}
	return &DSU{parent, rank}
}

func (dsu *DSU) find(i int) int {
	if dsu.parent[i] == -1 {
		return i
	}
	dsu.parent[i] = dsu.find(dsu.parent[i])
	return dsu.parent[i]
}

func (dsu *DSU) unite(x, y int) {
	s1 := dsu.find(x)
	s2 := dsu.find(y)
	if s1 != s2 {
		if dsu.rank[s1] < dsu.rank[s2] {
			dsu.parent[s1] = s2
		} else if dsu.rank[s1] > dsu.rank[s2] {
			dsu.parent[s2] = s1
		} else {
			dsu.parent[s2] = s1
			dsu.rank[s1]++
		}
	}
}

type Graph struct {
	edges []edge
	V     int
}

func newGraph(V int) *Graph {
	return &Graph{V: V}
}

func (g *Graph) addEdge(from, to, weight int) {
	g.edges = append(g.edges, edge{weight, from, to})
}

func (g *Graph) kruskalsMST() {
	sort.Slice(g.edges, func(i, j int) bool {
		return g.edges[i].weight < g.edges[j].weight
	})

	dsu := newDSU(g.V)
	mstCost := 0
	fmt.Println("Following are the edges in the constructed MST")

	for _, e := range g.edges {
		if dsu.find(e.from) != dsu.find(e.to) {
			dsu.unite(e.from, e.to)
			mstCost += e.weight
			fmt.Printf("%d -- %d == %d\n", e.from, e.to, e.weight)
		}
	}
	fmt.Printf("Minimum Cost Spanning Tree: %d\n", mstCost)
}

func main() {
	g := newGraph(18)

	edgesToAdd := []struct {
		from, to, weight int
	}{
		{1, 2, 1}, {1, 3, 2}, {1, 4, 4}, {2, 5, 4}, {2, 7, 2}, {3, 6, 6}, {3, 5, 7},
		{4, 6, 2}, {4, 7, 3}, {5, 8, 7}, {5, 9, 5}, {6, 8, 7}, {6, 10, 3}, {7, 9, 5},
		{7, 10, 4}, {8, 11, 4}, {9, 11, 1}, {10, 11, 4},
	}

	for _, e := range edgesToAdd {
		g.addEdge(e.from, e.to, e.weight)
	}

	g.kruskalsMST()
}
