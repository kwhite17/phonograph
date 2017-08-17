package search

import "testing"

var graph = buildTestGraph()
var nodeList = make([]Node, 0)

func buildTestGraph() map[Node][]Node {
	for i := 0; i < 12; i++ {
		nodeList = append(nodeList, &NumberNode{Value: i})
	}
	graph := make(map[Node][]Node)
	graph[nodeList[0]] = []Node{nodeList[1], nodeList[2]}
	graph[nodeList[1]] = []Node{nodeList[0], nodeList[4]}
	graph[nodeList[2]] = []Node{nodeList[0], nodeList[3], nodeList[7]}
	graph[nodeList[3]] = []Node{nodeList[2], nodeList[7]}
	graph[nodeList[4]] = []Node{nodeList[1], nodeList[5], nodeList[6], nodeList[9]}
	graph[nodeList[5]] = []Node{nodeList[4], nodeList[11]}
	graph[nodeList[6]] = []Node{nodeList[4], nodeList[11]}
	graph[nodeList[7]] = []Node{nodeList[2], nodeList[3], nodeList[8]}
	graph[nodeList[8]] = []Node{nodeList[7]}
	graph[nodeList[9]] = []Node{nodeList[4]}
	graph[nodeList[10]] = []Node{}
	graph[nodeList[11]] = []Node{nodeList[5], nodeList[6]}

	return graph
}

func teardown() {
	for _, node := range nodeList {
		node.setFoundAncestors(false)
		node.setFoundSuccessors(false)
		node.setChild(nil)
		node.setParent(nil)
	}
}

func TestBasicSearch(t *testing.T) {
	defer teardown()
	source := nodeList[0]
	dest := nodeList[1]
	result := bfs(source, dest, graph)
	if result.Len() != 2 {
		t.Errorf("Expected result of length 2. Actual result length is %d", result.Len())
	}
}

func TestGrandparentSearch(t *testing.T) {
	defer teardown()
	source := nodeList[0]
	dest := nodeList[7]
	result := bfs(source, dest, graph)
	if result.Len() != 3 {
		t.Errorf("Expected result of length 3. Actual result length is %d", result.Len())
	}
}

func TestPathNotFound(t *testing.T) {
	defer teardown()
	source := nodeList[10]
	dest := nodeList[0]
	result := bfs(source, dest, graph)
	if result.Len() != 0 {
		t.Errorf("Expected result of length 0. Actual result length is %d", result.Len())
	}
}

func TestNilSource(t *testing.T) {
	defer teardown()
	dest := nodeList[0]
	result := bfs(nil, dest, graph)
	if result.Len() != 0 {
		t.Errorf("Expected result of length 0. Actual result length is %d", result.Len())
	}
}

func TestNilDestination(t *testing.T) {
	defer teardown()
	source := nodeList[0]
	result := bfs(source, nil, graph)
	if result.Len() != 0 {
		t.Errorf("Expected result of length 0. Actual result length is %d", result.Len())
	}
}

func TestMultipleShortestPaths(t *testing.T) {
	defer teardown()
	source := nodeList[0]
	dest := nodeList[11]
	result := bfs(source, dest, graph)
	if result.Len() != 5 {
		t.Errorf("Expected result of length 5. Actual result length is %d", result.Len())
	}
}

func TestBasicBidirectionalSearch(t *testing.T) {
	defer teardown()
	source := nodeList[0]
	dest := nodeList[1]
	result := BidirectionalBfs(source, dest, graph)
	if result.Len() != 2 {
		t.Errorf("Expected result of length 2. Actual result length is %d", result.Len())
	}
}

func TestGrandparentBidirectionalSearch(t *testing.T) {
	defer teardown()
	source := nodeList[0]
	dest := nodeList[7]
	result := BidirectionalBfs(source, dest, graph)
	if result.Len() != 3 {
		t.Errorf("Expected result of length 3. Actual result length is %d", result.Len())
	}
}

func TestBidirectionalNilSource(t *testing.T) {
	defer teardown()
	dest := nodeList[0]
	result := BidirectionalBfs(nil, dest, graph)
	if result.Len() != 0 {
		t.Errorf("Expected result of length 0. Actual result length is %d", result.Len())
	}
}

func TestBidirectionalNilDestination(t *testing.T) {
	defer teardown()
	source := nodeList[0]
	result := BidirectionalBfs(source, nil, graph)
	if result.Len() != 0 {
		t.Errorf("Expected result of length 0. Actual result length is %d", result.Len())
	}
}

func TestBidirectionalMultipleShortestPaths(t *testing.T) {
	defer teardown()
	source := nodeList[0]
	dest := nodeList[11]
	result := BidirectionalBfs(source, dest, graph)
	if result.Len() != 5 {
		t.Errorf("Expected result of length 5. Actual result length is %d", result.Len())
	}
}
