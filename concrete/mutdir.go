// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package concrete

import (
	"github.com/gonum/graph"
)

// A Directed graph is a highly generalized MutableDirectedGraph.
//
// In most cases it's likely more desireable to use a graph specific to your
// problem domain.
type DirectedGraph struct {
	successors   map[int]map[int]WeightedEdge
	predecessors map[int]map[int]WeightedEdge
	nodeMap      map[int]graph.Node

	// Add/remove convenience variables
	maxID   int
	freeMap map[int]struct{}
}

func NewDirectedGraph() *DirectedGraph {
	return &DirectedGraph{
		successors:   make(map[int]map[int]WeightedEdge),
		predecessors: make(map[int]map[int]WeightedEdge),
		nodeMap:      make(map[int]graph.Node),
		maxID:        0,
		freeMap:      make(map[int]struct{}),
	}
}

/* Mutable Graph implementation */

func (g *DirectedGraph) NewNode() graph.Node {
	if g.maxID != maxInt {
		g.maxID++
		return Node(g.maxID)
	}

	// Implicitly checks if len(g.freeMap) == 0
	for id := range g.freeMap {
		return Node(id)
	}

	// I cannot foresee this ever happening, but just in case
	if len(g.nodeMap) == maxInt {
		panic("cannot allocate node: graph too large")
	}

	for i := 0; i < maxInt; i++ {
		if _, ok := g.nodeMap[i]; !ok {
			return Node(i)
		}
	}

	// Should not happen.
	panic("cannot allocate node id: no free id found")
}

// Adds a node to the graph. Implementation note: if you add a node close to or at
// the max int on your machine NewNode will become slower.
func (g *DirectedGraph) AddNode(n graph.Node) {
	g.nodeMap[n.ID()] = n
	g.successors[n.ID()] = make(map[int]WeightedEdge)
	g.predecessors[n.ID()] = make(map[int]WeightedEdge)

	delete(g.freeMap, n.ID())
	g.maxID = max(g.maxID, n.ID())
}

func (g *DirectedGraph) AddDirectedEdge(e graph.Edge, cost float64) {
	head, tail := e.Head(), e.Tail()
	if !g.Has(head) {
		g.AddNode(head)
	}

	if !g.Has(tail) {
		g.AddNode(tail)
	}

	g.successors[head.ID()][tail.ID()] = WeightedEdge{Edge: e, Cost: cost}
	g.predecessors[tail.ID()][head.ID()] = WeightedEdge{Edge: e, Cost: cost}
}

func (g *DirectedGraph) RemoveNode(n graph.Node) {
	if _, ok := g.nodeMap[n.ID()]; !ok {
		return
	}
	delete(g.nodeMap, n.ID())

	for succ := range g.successors[n.ID()] {
		delete(g.predecessors[succ], n.ID())
	}
	delete(g.successors, n.ID())

	for pred := range g.predecessors[n.ID()] {
		delete(g.successors[pred], n.ID())
	}
	delete(g.predecessors, n.ID())

	g.maxID-- // Fun facts: even if this ID doesn't exist this still works!
	g.freeMap[n.ID()] = struct{}{}
}

func (g *DirectedGraph) RemoveDirectedEdge(e graph.Edge) {
	head, tail := e.Head(), e.Tail()
	if _, ok := g.nodeMap[head.ID()]; !ok {
		return
	} else if _, ok := g.nodeMap[tail.ID()]; !ok {
		return
	}

	delete(g.successors[head.ID()], tail.ID())
	delete(g.predecessors[tail.ID()], head.ID())
}

func (g *DirectedGraph) EmptyGraph() {
	g.successors = make(map[int]map[int]WeightedEdge)
	g.predecessors = make(map[int]map[int]WeightedEdge)
	g.nodeMap = make(map[int]graph.Node)
}

/* Graph implementation */

func (g *DirectedGraph) From(n graph.Node) []graph.Node {
	if _, ok := g.successors[n.ID()]; !ok {
		return nil
	}

	successors := make([]graph.Node, len(g.successors[n.ID()]))
	i := 0
	for succ := range g.successors[n.ID()] {
		successors[i] = g.nodeMap[succ]
		i++
	}

	return successors
}

func (g *DirectedGraph) HasEdge(n, succ graph.Node) bool {
	if _, ok := g.nodeMap[n.ID()]; !ok {
		return false
	} else if _, ok := g.nodeMap[succ.ID()]; !ok {
		return false
	}
	_, ok := g.successors[n.ID()][succ.ID()]
	return ok
}

func (g *DirectedGraph) EdgeFromTo(n, succ graph.Node) graph.Edge {
	if _, ok := g.nodeMap[n.ID()]; !ok {
		return nil
	} else if _, ok := g.nodeMap[succ.ID()]; !ok {
		return nil
	}

	edge, ok := g.successors[n.ID()][succ.ID()]
	if !ok {
		return nil
	}
	return edge
}

func (g *DirectedGraph) To(n graph.Node) []graph.Node {
	if _, ok := g.successors[n.ID()]; !ok {
		return nil
	}

	predecessors := make([]graph.Node, len(g.predecessors[n.ID()]))
	i := 0
	for succ := range g.predecessors[n.ID()] {
		predecessors[i] = g.nodeMap[succ]
		i++
	}

	return predecessors
}

func (g *DirectedGraph) EdgeBetween(n, neigh graph.Node) graph.Edge {
	e := g.EdgeFromTo(n, neigh)
	if e != nil {
		return e
	}

	e = g.EdgeFromTo(neigh, n)
	if e != nil {
		return e
	}

	return nil
}

func (g *DirectedGraph) Has(n graph.Node) bool {
	_, ok := g.nodeMap[n.ID()]

	return ok
}

func (g *DirectedGraph) Degree(n graph.Node) int {
	if _, ok := g.nodeMap[n.ID()]; !ok {
		return 0
	}

	return len(g.successors[n.ID()]) + len(g.predecessors[n.ID()])
}

func (g *DirectedGraph) Order() int {
	return len(g.successors)
}

func (g *DirectedGraph) Nodes() []graph.Node {
	nodes := make([]graph.Node, len(g.successors))
	i := 0
	for _, n := range g.nodeMap {
		nodes[i] = n
		i++
	}

	return nodes
}

func (g *DirectedGraph) Cost(e graph.Edge) float64 {
	if s, ok := g.successors[e.Head().ID()]; ok {
		if we, ok := s[e.Tail().ID()]; ok {
			return we.Cost
		}
	}
	return inf
}

func (g *DirectedGraph) EdgeList() []graph.Edge {
	edgeList := make([]graph.Edge, 0, len(g.successors))
	edgeMap := make(map[int]map[int]struct{}, len(g.successors))
	for n, succMap := range g.successors {
		edgeMap[n] = make(map[int]struct{}, len(succMap))
		for succ, edge := range succMap {
			if doneMap, ok := edgeMap[succ]; ok {
				if _, ok := doneMap[n]; ok {
					continue
				}
			}
			edgeList = append(edgeList, edge)
			edgeMap[n][succ] = struct{}{}
		}
	}

	return edgeList
}
