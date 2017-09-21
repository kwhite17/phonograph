package search

import (
	"testing"

	"github.com/kwhite17/phonograph/client"
	"github.com/kwhite17/phonograph/model"
)

var graph = buildTestGraph()
var nodeList = make([]model.Node, 0)

func buildTestGraph() client.GraphClient {
	for i := 0; i < 12; i++ {
		nodeList = append(nodeList, &model.GenericNode{Value: i})
	}
	graph := make(map[interface{}][]interface{})
	graph[nodeList[0]] = []interface{}{nodeList[1], nodeList[2]}
	graph[nodeList[1]] = []interface{}{nodeList[0], nodeList[4]}
	graph[nodeList[2]] = []interface{}{nodeList[0], nodeList[3], nodeList[7]}
	graph[nodeList[3]] = []interface{}{nodeList[2], nodeList[7]}
	graph[nodeList[4]] = []interface{}{nodeList[1], nodeList[5], nodeList[6], nodeList[9]}
	graph[nodeList[5]] = []interface{}{nodeList[4], nodeList[11]}
	graph[nodeList[6]] = []interface{}{nodeList[4], nodeList[11]}
	graph[nodeList[7]] = []interface{}{nodeList[2], nodeList[3], nodeList[8]}
	graph[nodeList[8]] = []interface{}{nodeList[7]}
	graph[nodeList[9]] = []interface{}{nodeList[4]}
	graph[nodeList[10]] = []interface{}{}
	graph[nodeList[11]] = []interface{}{nodeList[5], nodeList[6]}
	gCli := client.GraphClient{Graph: graph}
	return gCli
}

func teardown() {
	for _, node := range nodeList {
		node.SetFoundAncestors(false)
		node.SetFoundSuccessors(false)
		node.SetChild(nil)
		node.SetParent(nil)
	}
}

func TestBasicSearch(t *testing.T) {
	defer teardown()
	source := nodeList[0]
	dest := nodeList[1]
	result := bfs(source, dest, graph)
	if len(result) != 2 {
		t.Errorf("Expected result of length 2. Actual result length is %d", len(result))
	}
}

func TestGrandparentSearch(t *testing.T) {
	defer teardown()
	source := nodeList[0]
	dest := nodeList[7]
	result := bfs(source, dest, graph)
	if len(result) != 3 {
		t.Errorf("Expected result of length 3. Actual result length is %d", len(result))
	}
}

func TestPathNotFound(t *testing.T) {
	defer teardown()
	source := nodeList[10]
	dest := nodeList[0]
	result := bfs(source, dest, graph)
	if len(result) != 0 {
		t.Errorf("Expected result of length 0. Actual result length is %d", len(result))
	}
}

func TestNilSource(t *testing.T) {
	defer teardown()
	dest := nodeList[0]
	result := bfs(nil, dest, graph)
	if len(result) != 0 {
		t.Errorf("Expected result of length 0. Actual result length is %d", len(result))
	}
}

func TestNilDestination(t *testing.T) {
	defer teardown()
	source := nodeList[0]
	result := bfs(source, nil, graph)
	if len(result) != 0 {
		t.Errorf("Expected result of length 0. Actual result length is %d", len(result))
	}
}

func TestMultipleShortestPaths(t *testing.T) {
	defer teardown()
	source := nodeList[0]
	dest := nodeList[11]
	result := bfs(source, dest, graph)
	if len(result) != 5 {
		t.Errorf("Expected result of length 5. Actual result length is %d", len(result))
	}
}

func TestBasicBidirectionalSearch(t *testing.T) {
	defer teardown()
	source := nodeList[0]
	dest := nodeList[1]
	result := BidirectionalBfs(source, dest, graph)
	if len(result) != 2 {
		t.Errorf("Expected result of length 2. Actual result length is %d", len(result))
	}
}

func TestGrandparentBidirectionalSearch(t *testing.T) {
	defer teardown()
	source := nodeList[0]
	dest := nodeList[7]
	result := BidirectionalBfs(source, dest, graph)
	if len(result) != 3 {
		t.Errorf("Expected result of length 3. Actual result length is %d", len(result))
	}
}

func TestBidirectionalNilSource(t *testing.T) {
	defer teardown()
	dest := nodeList[0]
	result := BidirectionalBfs(nil, dest, graph)
	if len(result) != 0 {
		t.Errorf("Expected result of length 0. Actual result length is %d", len(result))
	}
}

func TestBidirectionalNilDestination(t *testing.T) {
	defer teardown()
	source := nodeList[0]
	result := BidirectionalBfs(source, nil, graph)
	if len(result) != 0 {
		t.Errorf("Expected result of length 0. Actual result length is %d", len(result))
	}
}

func TestBidirectionalMultipleShortestPaths(t *testing.T) {
	defer teardown()
	source := nodeList[0]
	dest := nodeList[11]
	result := BidirectionalBfs(source, dest, graph)
	if len(result) != 5 {
		t.Errorf("Expected result of length 5. Actual result length is %d", len(result))
	}
}
