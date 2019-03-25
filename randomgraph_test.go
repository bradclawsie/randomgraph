package randomgraph

import (
	"testing"
)

func TestValidVertices(t *testing.T) {
	t.Run("Vertices=nil", func(t *testing.T) {
		err := ValidVertices(Vertices(nil))
		if err == nil {
			t.Error("nil vertices")
		}
	})
	t.Run("Vertices=singleton", func(t *testing.T) {
		err := ValidVertices(Vertices([]Vertex{"a"}))
		if err == nil {
			t.Error("singleton vertices")
		}
	})
	t.Run("Vertices=duplicates", func(t *testing.T) {
		err := ValidVertices(Vertices([]Vertex{"a", "b", "c", "b"}))
		if err == nil {
			t.Error("duplicate vertices")
		}
	})
	t.Run("Vertices=ok", func(t *testing.T) {
		err := ValidVertices(Vertices([]Vertex{"a", "b", "c"}))
		if err != nil {
			t.Error(err)
		}
	})
}

func TestUndirectedCyclic(t *testing.T) {

	RandSeed()

	f := func(vs Vertices) {
		g, err := UndirectedCyclic(vs)
		if err != nil {
			t.Error(err)
		}
		count := 0
		for src, dsts := range g {
			count++
			for _, dst := range dsts {
				if _, dstExists := g[dst]; !dstExists {
					t.Errorf("dst %v not in graph", dst)
				}
				srcFound := false
				for _, dstDst := range g[dst] {
					if dstDst == src {
						srcFound = true
					}
				}
				if !srcFound {
					t.Errorf("dst %v edge back to %v not found", dst, src)
				}
			}
		}
		if count != len(vs) {
			t.Errorf("count of gsrc edges %d not equal to %d", count, len(vs))
		}
	}

	t.Run("Vertices=Simple", func(t *testing.T) {
		f(Vertices([]Vertex{"a", "b", "c", "d", "e", "f"}))
	})

	t.Run("Vertices=Edge", func(t *testing.T) {
		f(Vertices([]Vertex{"a", "b"}))
	})
}

func TestDirectedAcyclic(t *testing.T) {

	RandSeed()

	f := func(vs Vertices) {
		g, err := DirectedAcyclic(vs)
		if err != nil {
			t.Error(err)
		}

		count := 0
		for _, dsts := range g {
			count++
			for _, dst := range dsts {
				if _, dstExists := g[dst]; !dstExists {
					t.Errorf("dst %v not in graph", dst)
				}
			}
		}
		if count != len(vs) {
			t.Errorf("count of gsrc edges %d not equal to %d", count, len(vs))
		}
	}

	t.Run("Vertices=Simple", func(t *testing.T) {
		f(Vertices([]Vertex{"a", "b", "c", "d", "e", "f"}))
	})

	t.Run("Vertices=Edge", func(t *testing.T) {
		f(Vertices([]Vertex{"a", "b"}))
	})
}
