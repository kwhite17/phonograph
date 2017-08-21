package model

import (
	"reflect"
	"sync"

	"fmt"
)

type Node interface {
	SetParent(parent Node)
	GetParent() Node
	SetChild(child Node)
	GetChild() Node
	GetValue() interface{}
	SetFoundAncestors(isVisited bool)
	SetFoundSuccessors(isVisited bool)
	HasFoundAncestors() bool
	HasFoundSuccessors() bool
	GetLock() *sync.Mutex
}

type ArtistNode struct {
	child           *ArtistNode
	parent          *ArtistNode
	Value           Artist
	foundAncestors  bool
	foundSuccessors bool
	lock            *sync.Mutex
}

type GenericNode struct {
	child           *GenericNode
	parent          *GenericNode
	Value           interface{}
	foundAncestors  bool
	foundSuccessors bool
	lock            *sync.Mutex
}

func (an *ArtistNode) GetChild() Node {
	return an.child
}

func (an *ArtistNode) GetParent() Node {
	return an.parent
}

func (an *ArtistNode) SetParent(parent Node) {
	if !IsNil(parent) {
		node := parent.(*ArtistNode)
		an.parent = node
	} else {
		an.parent = nil
	}
}

func (an *ArtistNode) SetChild(child Node) {
	if !IsNil(child) {
		node := child.(*ArtistNode)
		an.child = node
	} else {
		an.child = nil
	}
}

func (an *ArtistNode) SetFoundAncestors(isVisited bool) {
	an.foundAncestors = isVisited
}

func (an *ArtistNode) SetFoundSuccessors(isVisited bool) {
	an.foundSuccessors = isVisited
}

func (an *ArtistNode) HasFoundAncestors() bool {
	return an.foundAncestors
}

func (an *ArtistNode) HasFoundSuccessors() bool {
	return an.foundSuccessors
}

func (an *ArtistNode) GetLock() *sync.Mutex {
	if IsNil(an.lock) {
		an.lock = &sync.Mutex{}
	}
	return an.lock
}

func (an *ArtistNode) GetValue() interface{} {
	return an.Value
}

func (nn *GenericNode) GetChild() Node {
	return nn.child
}

func (nn *GenericNode) GetParent() Node {
	return nn.parent
}

func (nn *GenericNode) SetParent(parent Node) {
	if !IsNil(parent) {
		node := parent.(*GenericNode)
		nn.parent = node
	} else {
		nn.parent = nil
	}
}

func (nn *GenericNode) SetChild(child Node) {
	if !IsNil(child) {
		node := child.(*GenericNode)
		nn.child = node
	} else {
		nn.child = nil
	}
}

func (nn *GenericNode) GetLock() *sync.Mutex {
	if IsNil(nn.lock) {
		nn.lock = &sync.Mutex{}
	}
	return nn.lock
}

func (nn *GenericNode) SetFoundAncestors(isVisited bool) {
	nn.foundAncestors = isVisited
}

func (nn *GenericNode) SetFoundSuccessors(isVisited bool) {
	nn.foundSuccessors = isVisited
}

func (nn *GenericNode) HasFoundAncestors() bool {
	return nn.foundAncestors
}

func (nn *GenericNode) HasFoundSuccessors() bool {
	return nn.foundSuccessors
}

func (nn *GenericNode) string() string {
	return fmt.Sprintf("%v\n", nn.Value)
}

func (nn *GenericNode) GetValue() interface{} {
	return nn.Value
}

func IsNil(a interface{}) bool {
	defer func() { recover() }()
	return a == nil || reflect.ValueOf(a).IsNil()
}
