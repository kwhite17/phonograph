package search

import (
	"sync"

	"github.com/kwhite17/phonograph/client"
	"github.com/kwhite17/phonograph/model"
)

type Node interface {
	expand(edges []Node) []Node
	setParent(parent Node)
	getParent() Node
	setChild(child Node)
	getChild() Node
	setFoundAncestors(isVisited bool)
	setFoundSuccessors(isVisited bool)
	hasFoundAncestors() bool
	hasFoundSuccessors() bool
	getLock() *sync.Mutex
}

type ArtistNode struct {
	child           *ArtistNode
	parent          *ArtistNode
	Value           model.Artist
	foundAncestors  bool
	foundSuccessors bool
	lock            *sync.Mutex
}

var sClient = client.InitSpotifyClient()

type NumberNode struct {
	child           *NumberNode
	parent          *NumberNode
	Value           int
	foundAncestors  bool
	foundSuccessors bool
	lock            *sync.Mutex
}

func (an *ArtistNode) getChild() Node {
	return an.child
}

func (an *ArtistNode) getParent() Node {
	return an.parent
}

func (an *ArtistNode) setParent(parent Node) {
	if !IsNil(parent) {
		node := parent.(*ArtistNode)
		an.parent = node
	} else {
		an.parent = nil
	}
}

func (an *ArtistNode) setChild(child Node) {
	if !IsNil(child) {
		node := child.(*ArtistNode)
		an.child = node
	} else {
		an.child = nil
	}
}

func (an *ArtistNode) setFoundAncestors(isVisited bool) {
	an.foundAncestors = isVisited
}

func (an *ArtistNode) setFoundSuccessors(isVisited bool) {
	an.foundSuccessors = isVisited
}

func (an *ArtistNode) hasFoundAncestors() bool {
	return an.foundAncestors
}

func (an *ArtistNode) hasFoundSuccessors() bool {
	return an.foundSuccessors
}

func (an *ArtistNode) expand(edges []Node) []Node {
	//convertEdges to songs
	collaborators := sClient.GetAssociatedArtists(an.Value)
	artistNodes := make([]Node, 0)
	for _, artist := range collaborators {
		artistNodes = append(artistNodes, &ArtistNode{Value: artist})
	}
	return artistNodes
}

func (an *ArtistNode) getLock() *sync.Mutex {
	if IsNil(an.lock) {
		an.lock = &sync.Mutex{}
	}
	return an.lock
}

func (nn *NumberNode) getChild() Node {
	return nn.child
}

func (nn *NumberNode) getParent() Node {
	return nn.parent
}

func (nn *NumberNode) setParent(parent Node) {
	if !IsNil(parent) {
		node := parent.(*NumberNode)
		nn.parent = node
	} else {
		nn.parent = nil
	}
}

func (nn *NumberNode) setChild(child Node) {
	if !IsNil(child) {
		node := child.(*NumberNode)
		nn.child = node
	} else {
		nn.child = nil
	}
}

func (nn *NumberNode) getLock() *sync.Mutex {
	if IsNil(nn.lock) {
		nn.lock = &sync.Mutex{}
	}
	return nn.lock
}

func (nn *NumberNode) setFoundAncestors(isVisited bool) {
	nn.foundAncestors = isVisited
}

func (nn *NumberNode) setFoundSuccessors(isVisited bool) {
	nn.foundSuccessors = isVisited
}

func (nn *NumberNode) hasFoundAncestors() bool {
	return nn.foundAncestors
}

func (nn *NumberNode) hasFoundSuccessors() bool {
	return nn.foundSuccessors
}

func (nn *NumberNode) expand(edges []Node) []Node {
	numberEdges := make([]Node, 0)
	for i := 0; i < len(edges); i++ {
		numberEdges = append(numberEdges, edges[i].(Node))
	}
	return numberEdges
}

func (nn *NumberNode) string() string {
	return string(nn.Value)
}
