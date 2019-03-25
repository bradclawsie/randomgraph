// Package randomgraph creates different varieties of randomly generated graphs.
// TODO: make sure reulting graphs are not partitioned.
package randomgraph

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Vertex is an alias for a string label.
type Vertex string

// Vertices is a list of Vertex.
type Vertices []Vertex

// Graph maps a Vertex (source) to destination Vertices.
type Graph map[Vertex]Vertices

// RandSeed seeds the rng.
func RandSeed() {
	rand.Seed(time.Now().UnixNano())
}

// ValidVertices detects nil, singleton, or vertex lists with repeated elements.
func ValidVertices(vs Vertices) error {
	if vs == nil || len(vs) < 2 {
		return errors.New("nil or inadequate vertices")
	}
	seen := make(map[Vertex]bool)
	for _, v := range vs {
		if _, duplicate := seen[v]; duplicate {
			return fmt.Errorf("repeated vertex %v", v)
		}
		seen[v] = true
	}
	return nil
}

// UndirectedCyclic creates a random graph from a set of Vertices.
// Each edge a -> b has a corresponding edge b -> a (undirected, cyclic).
func UndirectedCyclic(vs Vertices) (Graph, error) {

	vsErr := ValidVertices(vs)
	if vsErr != nil {
		return nil, vsErr
	}

	vsShuf := make(Vertices, len(vs))
	copy(vsShuf, vs)
	l := len(vsShuf)
	vsSrc := make(Vertices, len(vs))
	copy(vsSrc, vs)
	ps := rand.Perm(len(vsSrc))
	for i, j := range ps {
		vsSrc[i] = vs[j]
	}
	m := make(map[Vertex]map[Vertex]bool)

	for _, src := range vsSrc {

		// Get a random shuffle of indices for vsShuf.
		ps := rand.Perm(l)
		// Swap elements to get a random permutation of vsShuf.
		for i, j := range ps {
			vsShuf[i] = vsSrc[j]
		}
		// Pick a random number of vertices in vsShuf to be
		// joined with src by an edge.
		pick := 0
		for {
			pick = rand.Intn(l)
			if pick != 0 {
				break
			}
		}
		if _, exists := m[src]; !exists {
			m[src] = make(map[Vertex]bool)
		}
		picked := 0
		for _, dst := range vsShuf {
			if dst != src {
				m[src][dst] = true
				if _, exists := m[dst]; !exists {
					m[dst] = make(map[Vertex]bool)
				}
				m[dst][src] = true
				picked++
			}
			if picked == pick {
				break
			}
		}
	}

	// Turn m into a Graph.
	g := make(Graph)
	for src, dsts := range m {
		g[src] = make(Vertices, len(dsts))
		i := 0
		for dst := range m[src] {
			g[src][i] = dst
			i++
		}
	}
	return g, nil
}

// DirectedAcyclic creates a random graph from a set of Vertices.
// No cycles should be created in resulting graph.
func DirectedAcyclic(vs Vertices) (Graph, error) {

	vsErr := ValidVertices(vs)
	if vsErr != nil {
		return nil, vsErr
	}

	vsShuf := make(Vertices, len(vs))
	copy(vsShuf, vs)
	l := len(vsShuf)
	vsSrc := make(Vertices, len(vs))
	copy(vsSrc, vs)
	ps := rand.Perm(len(vsSrc))
	for i, j := range ps {
		vsSrc[i] = vs[j]
	}
	m := make(map[Vertex]map[Vertex]bool)

	// canAdd is a BFS to determine if connecting srcTry to dstTry
	// makes a cycle.
	canAdd := func(srcTry, dstTry Vertex) bool {
		if srcTry == dstTry {
			return false
		}
		q := make(Vertices, 0)
		q = append(q, dstTry)
		for len(q) != 0 {
			head := q[0]
			if head == srcTry {
				return false
			}
			q = q[1:]
			if _, headExists := m[head]; headExists {
				qCopy := make(Vertices, len(q))
				copy(qCopy, q)
				for k := range m[head] {
					duplicate := false
					for _, qElt := range q {
						if qElt == k {
							duplicate = true
							break
						}
					}
					if !duplicate {
						qCopy = append(qCopy, k)
					}
				}
				q = qCopy
			}
		}
		return true
	}

	for _, src := range vsSrc {

		// Get a random shuffle of indices for vsShuf.
		ps := rand.Perm(l)
		// Swap elements to get a random permutation of vsShuf.
		for i, j := range ps {
			vsShuf[i] = vsSrc[j]
		}
		if _, exists := m[src]; !exists {
			m[src] = make(map[Vertex]bool)
		}
		for _, dst := range vsShuf {
			if dst != src {
				if canAdd(src, dst) {
					m[src][dst] = true
				}
			}
		}
	}

	// Turn m into a Graph.
	g := make(Graph)
	for src, dsts := range m {
		g[src] = make(Vertices, len(dsts))
		i := 0
		for dst := range m[src] {
			g[src][i] = dst
			i++
		}
	}
	return g, nil
}
