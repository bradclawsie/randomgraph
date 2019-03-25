# randomgraph
Generate random graphs.

```
randomgraph.RandSeed()
vs := randomgraph.Vertices([]randomgraph.Vertex{"a","b","c","d","e"})
graph, err := randomgraph.UndirectedCyclic(vs)
graph, err := randomgraph.DirectedAcyclic(vs)
```
