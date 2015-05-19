// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package concrete_test

import (
	"math"
	"testing"

	"github.com/gonum/graph"
	"github.com/gonum/graph/concrete"
)

var (
	_ graph.Graph      = (*concrete.UndirectedGraph)(nil)
	_ graph.Undirected = (*concrete.UndirectedGraph)(nil)
)

func TestAssertMutableNotDirected(t *testing.T) {
	var g graph.MutableUndirected = concrete.NewUndirectedGraph(math.Inf(1))
	if _, ok := g.(graph.Directed); ok {
		t.Fatal("concrete.UndirectedGraph is directed, but a MutableGraph cannot safely be directed!")
	}
}

func TestMaxID(t *testing.T) {
	g := concrete.NewUndirectedGraph(math.Inf(1))
	nodes := make(map[graph.Node]struct{})
	for i := concrete.Node(0); i < 3; i++ {
		g.AddNode(i)
		nodes[i] = struct{}{}
	}
	g.RemoveNode(concrete.Node(0))
	delete(nodes, concrete.Node(0))
	g.RemoveNode(concrete.Node(2))
	delete(nodes, concrete.Node(2))
	n := g.NewNode()
	g.AddNode(n)
	if !g.Has(n) {
		t.Error("added node does not exist in graph")
	}
	if _, exists := nodes[n]; exists {
		t.Errorf("Created already existing node id: %v", n.ID())
	}
}
