// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package concrete

import (
	"github.com/gonum/graph"
)

// A dense graph is a graph such that all IDs are in a contiguous block from 0 to
// TheNumberOfNodes-1. It uses an adjacency matrix and should be relatively fast for both access
// and writing.
type DirectedDenseGraph struct {
	adjacencyMatrix []float64
	numNodes        int
}

// Creates a dense graph with the proper number of nodes. If passable is true all nodes will have
// an edge with unit cost, otherwise every node will start unconnected (cost of +Inf).
func NewDirectedDenseGraph(numNodes int, passable bool) *DirectedDenseGraph {
	g := &DirectedDenseGraph{adjacencyMatrix: make([]float64, numNodes*numNodes), numNodes: numNodes}
	if passable {
		for i := range g.adjacencyMatrix {
			g.adjacencyMatrix[i] = 1
		}
	} else {
		for i := range g.adjacencyMatrix {
			g.adjacencyMatrix[i] = inf
		}
	}

	return g
}

func (g *DirectedDenseGraph) Has(n graph.Node) bool {
	return n.ID() < g.numNodes
}

func (g *DirectedDenseGraph) Degree(n graph.Node) int {
	deg := 0
	for i := 0; i < g.numNodes; i++ {
		if g.adjacencyMatrix[i*g.numNodes+n.ID()] != inf {
			deg++
		}

		if g.adjacencyMatrix[n.ID()*g.numNodes+i] != inf {
			deg++
		}
	}

	return deg
}

func (g *DirectedDenseGraph) Nodes() []graph.Node {
	nodes := make([]graph.Node, g.numNodes)
	for i := 0; i < g.numNodes; i++ {
		nodes[i] = Node(i)
	}

	return nodes
}

func (g *DirectedDenseGraph) DirectedEdgeList() []graph.Edge {
	edges := make([]graph.Edge, 0, len(g.adjacencyMatrix))
	for i := 0; i < g.numNodes; i++ {
		for j := 0; j < g.numNodes; j++ {
			if g.adjacencyMatrix[i*g.numNodes+j] != inf {
				edges = append(edges, Edge{Node(i), Node(j)})
			}
		}
	}

	return edges
}

func (g *DirectedDenseGraph) From(n graph.Node) []graph.Node {
	neighbors := make([]graph.Node, 0)
	for i := 0; i < g.numNodes; i++ {
		if g.adjacencyMatrix[n.ID()*g.numNodes+i] != inf {
			neighbors = append(neighbors, Node(i))
		}
	}

	return neighbors
}

func (g *DirectedDenseGraph) To(n graph.Node) []graph.Node {
	neighbors := make([]graph.Node, 0)
	for i := 0; i < g.numNodes; i++ {
		if g.adjacencyMatrix[i*g.numNodes+n.ID()] != inf {
			neighbors = append(neighbors, Node(i))
		}
	}

	return neighbors
}

func (g *DirectedDenseGraph) HasEdge(n, succ graph.Node) bool {
	return g.adjacencyMatrix[n.ID()*g.numNodes+succ.ID()] != inf
}

func (g *DirectedDenseGraph) EdgeFromTo(n, succ graph.Node) graph.Edge {
	if g.HasEdge(n, succ) {
		return Edge{n, succ}
	}
	return nil
}

func (g *DirectedDenseGraph) Cost(e graph.Edge) float64 {
	return g.adjacencyMatrix[e.Head().ID()*g.numNodes+e.Tail().ID()]
}

// Sets the cost of an edge. If the cost is +Inf, it will remove the edge,
// if directed is true, it will only remove the edge one way. If it's false it will change the cost
// of the edge from succ to node as well.
func (g *DirectedDenseGraph) SetEdgeCost(e graph.Edge, cost float64, directed bool) {
	g.adjacencyMatrix[e.Head().ID()*g.numNodes+e.Tail().ID()] = cost
}

// Equivalent to SetEdgeCost(edge, math.Inf(1), directed)
func (g *DirectedDenseGraph) RemoveEdge(e graph.Edge, directed bool) {
	g.SetEdgeCost(e, inf, directed)
}

func (g *DirectedDenseGraph) Crunch() {}
