package search

import (
	"container/list"
	"reflect"
)

func bidirectionalBfs(source Node, dest Node, graph map[Node][]Node) *list.List {
	resultChan := make(chan Node)
	go bidirectionalBfsHelper(source, dest, resultChan, graph, true)
	go bidirectionalBfsHelper(dest, source, resultChan, graph, false)
	result := list.New()
	for i := 0; i < 2; i++ {
		commonNode := <-resultChan
		switch commonNode {
		case nil:
			continue
		default:
			result = organizeResult(commonNode, source, dest)
		}
	}
	return result
}

func bidirectionalBfsHelper(start Node, end Node, resultChan chan Node, g map[Node][]Node, isSource bool) {
	if !isNil(start) {
		queue := list.New()
		queue.PushBack(start)
		for queue.Len() != 0 {
			element := queue.Remove(queue.Front()).(Node)
			if element != end {
				nextElements := make([]Node, 0)
				element.getLock().Lock()
				if isSource {
					element.setFoundSuccessors(true)
				} else {
					element.setFoundAncestors(true)
				}
				if g != nil {
					nextElements = element.expand(g[element])
				} else {
					nextElements = element.expand(nil)
				}
				element.getLock().Unlock()
				for i := 0; i < len(nextElements); i++ {
					cur := nextElements[i]
					cur.getLock().Lock()
					if isSource {
						if isNil(cur.getParent()) && !cur.hasFoundSuccessors() && cur != start {
							cur.setParent(element)
						}
						if !cur.hasFoundSuccessors() {
							queue.PushBack(cur)
						}
					} else {
						if isNil(cur.getChild()) && !cur.hasFoundAncestors() && cur != start {
							cur.setChild(element)
						}
						if !cur.hasFoundAncestors() {
							queue.PushBack(cur)
						}
					}
					if !isNil(cur.getParent()) && !isNil(cur.getChild()) {
						resultChan <- cur
						cur.getLock().Unlock()
						return
					}
					cur.getLock().Unlock()
				}
			} else {
				element.getLock().Lock()
				resultChan <- element
				element.getLock().Unlock()
				return
			}
		}
	}
	resultChan <- nil
}

func organizeResult(commonNode Node, source Node, dest Node) *list.List {
	finalList := list.New()
	if isNil(commonNode) {
		return finalList
	}
	curNode := commonNode
	for !isNil(curNode) {
		finalList.PushFront(curNode)
		curNode = curNode.getParent()
	}
	curNode = commonNode.getChild()
	for !isNil(curNode) {
		finalList.PushBack(curNode)
		curNode = curNode.getChild()
	}
	return finalList
}

func bfs(source Node, dest Node, graph map[Node][]Node) *list.List {
	result := list.New()
	if bfsHelper(source, dest, graph) {
		cur := dest
		for !isNil(cur) {
			result.PushFront(cur)
			cur = cur.getParent()
		}
	}
	return result
}

func bfsHelper(source Node, dest Node, graph map[Node][]Node) bool {
	if isNil(source) {
		return false
	}
	if isNil(dest) {
		return false
	}
	queue := list.New()
	queue.PushBack(source)
	for queue.Len() != 0 {
		element := queue.Remove(queue.Front()).(Node)
		element.setFoundSuccessors(true)
		nextElements := make([]Node, 0)
		if !isNil(graph) {
			nextElements = element.expand(graph[element])
		} else {
			nextElements = element.expand(nil)
		}
		for i := 0; i < len(nextElements); i++ {
			cur := nextElements[i]
			if isNil(cur.getParent()) && !cur.hasFoundSuccessors() {
				cur.setParent(element)
				if cur == dest {
					return true
				}
				queue.PushBack(cur)
			}
		}
	}
	return false
}

func isNil(a interface{}) bool {
	defer func() { recover() }()
	return a == nil || reflect.ValueOf(a).IsNil()
}
