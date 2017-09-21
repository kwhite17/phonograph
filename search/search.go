package search

import (
	"container/list"
	"fmt"

	"github.com/kwhite17/phonograph/client"
	"github.com/kwhite17/phonograph/model"
)

func BidirectionalBfs(source model.Node, dest model.Node, cli client.Client) []model.Node {
	sourceChan := make(chan model.Node)
	destChan := make(chan model.Node)
	go bidirectionalBfsHelper(source, dest, sourceChan, cli, true)
	go bidirectionalBfsHelper(dest, source, destChan, cli, false)
	result := make([]model.Node, 0)
	select {
	case dResult := <-destChan:
		if dResult != nil {
			result = organizeResult(dResult, source, dest)
			return result
		}
	case sResult := <-sourceChan:
		if sResult != nil {
			result = organizeResult(sResult, source, dest)
			return result
		}
	}
	return result
}

func bidirectionalBfsHelper(start model.Node, end model.Node, resultChan chan model.Node, cli client.Client, isSource bool) {
	if !model.IsNil(start) && !model.IsNil(end) {
		queue := list.New()
		queue.PushBack(start)
		for queue.Len() != 0 {
			element := queue.Remove(queue.Front()).(model.Node)
			if element != end {
				element.GetLock().Lock()
				if isSource {
					element.SetFoundSuccessors(true)
				} else {
					element.SetFoundAncestors(true)
				}
				nextElements := cli.Expand(element)
				element.GetLock().Unlock()
				for i := 0; i < len(nextElements); i++ {
					cur := nextElements[i]
					cur.GetLock().Lock()
					if isSource {
						if model.IsNil(cur.GetParent()) && !cur.HasFoundSuccessors() && cur != start {
							cur.SetParent(element)
						}
						if !cur.HasFoundSuccessors() {
							queue.PushBack(cur)
						}
					} else {
						if model.IsNil(cur.GetChild()) && !cur.HasFoundSuccessors() && cur != start {
							cur.SetChild(element)
						}
						if !cur.HasFoundSuccessors() {
							queue.PushBack(cur)
						}
					}
					if !model.IsNil(cur.GetParent()) && !model.IsNil(cur.GetChild()) {
						resultChan <- cur
						cur.GetLock().Unlock()
						return
					}
					cur.GetLock().Unlock()
				}
			} else {
				element.GetLock().Lock()
				resultChan <- element
				element.GetLock().Unlock()
				return
			}
		}
	}
	resultChan <- nil
}

func organizeResult(commonNode model.Node, source model.Node, dest model.Node) []model.Node {
	finalList := make([]model.Node, 0)
	if model.IsNil(commonNode) {
		return finalList
	}
	curNode := commonNode
	for !model.IsNil(curNode) {
		fmt.Println(curNode)
		finalList = append([]model.Node{curNode}, finalList...)
		curNode = curNode.GetParent()
	}
	curNode = commonNode.GetChild()
	for !model.IsNil(curNode) {
		fmt.Println(curNode)
		finalList = append(finalList, curNode)
		curNode = curNode.GetChild()
	}
	return finalList
}

func bfs(source model.Node, dest model.Node, cli client.Client) []model.Node {
	result := make([]model.Node, 0)
	if bfsHelper(source, dest, cli) {
		cur := dest
		for !model.IsNil(cur) {
			result = append([]model.Node{cur}, result...)
			cur = cur.GetParent()
		}
	}
	return result
}

func bfsHelper(source model.Node, dest model.Node, cli client.Client) bool {
	if model.IsNil(source) {
		return false
	}
	if model.IsNil(dest) {
		return false
	}
	queue := list.New()
	queue.PushBack(source)
	for queue.Len() != 0 {
		element := queue.Remove(queue.Front()).(model.Node)
		element.SetFoundSuccessors(true)
		nextElements := cli.Expand(element)
		for i := 0; i < len(nextElements); i++ {
			cur := nextElements[i]
			if model.IsNil(cur.GetParent()) && !cur.HasFoundSuccessors() {
				cur.SetParent(element)
				if cur == dest {
					return true
				}
				queue.PushBack(cur)
			}
		}
	}
	return false
}
