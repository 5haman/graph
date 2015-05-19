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
//
// This graph implements the CrunchGraph, but since it's naturally dense this is superfluous.
type UndirectedDenseGraph struct {
	adjacencyMatrix []float64
	numNodes        int
}

// Creates a dense graph with the proper number of nodes. If passable is true all nodes will have
// an edge with unit cost, otherwise every node will start unconnected (cost of +Inf).
func NewUndirectedDenseGraph(numNodes int, passable bool) *UndirectedDenseGraph {
	g := &UndirectedDenseGraph{adjacencyMatrix: make([]float64, numNodes*numNodes), numNodes: numNodes}
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

func (g *UndirectedDenseGraph) Has(n graph.Node) bool {
	return n.ID() < g.numNodes
}

func (g *UndirectedDenseGraph) Degree(n graph.Node) int {
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

func (g *UndirectedDenseGraph) Nodes() []graph.Node {
	nodes := make([]graph.Node, g.numNodes)
	for i := 0; i < g.numNodes; i++ {
		nodes[i] = Node(i)
	}

	return nodes
}

func (g *UndirectedDenseGraph) From(n graph.Node) []graph.Node {
	neighbors := make([]graph.Node, 0)
	for i := 0; i < g.numNodes; i++ {
		if g.adjacencyMatrix[i*g.numNodes+n.ID()] != inf ||
			g.adjacencyMatrix[n.ID()*g.numNodes+i] != inf {
			neighbors = append(neighbors, Node(i))
		}
	}

	return neighbors
}

func (g *UndirectedDenseGraph) HasEdge(n, neighbor graph.Node) bool {
	return g.adjacencyMatrix[neighbor.ID()*g.numNodes+n.ID()] != inf || g.adjacencyMatrix[n.ID()*g.numNodes+neighbor.ID()] != inf
}

func (g *UndirectedDenseGraph) EdgeBetween(n, neighbor graph.Node) graph.Edge {
	if g.HasEdge(n, neighbor) {
		return Edge{n, neighbor}
	}
	return nil
}

func (g *UndirectedDenseGraph) Cost(e graph.Edge) float64 {
	return g.adjacencyMatrix[e.Head().ID()*g.numNodes+e.Tail().ID()]
}

// Sets the cost of an edge. If the cost is +Inf, it will remove the edge,
// if directed is true, it will only remove the edge one way. If it's false it will change the cost
// of the edge from succ to node as well.
func (g *UndirectedDenseGraph) SetEdgeCost(e graph.Edge, cost float64, directed bool) {
	g.adjacencyMatrix[e.Head().ID()*g.numNodes+e.Tail().ID()] = cost
	g.adjacencyMatrix[e.Tail().ID()*g.numNodes+e.Head().ID()] = cost
}

// Equivalent to SetEdgeCost(edge, math.Inf(1), directed)
func (g *UndirectedDenseGraph) RemoveEdge(e graph.Edge, directed bool) {
	g.SetEdgeCost(e, inf, directed)
}

func (g *UndirectedDenseGraph) Crunch() {}
