package tiny

import (
	"strings"
)

type node struct {
	pattern  string
	segment  string
	children []*node
	isWild   bool
}

func (n *node) matchChild(segment string) *node {
	for _, child := range n.children {
		if child.segment == segment || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(segment string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.segment == segment || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, segments []string, height int) {
	if len(segments) == height {
		n.pattern = pattern
		return
	}
	segment := segments[height]
	child := n.matchChild(segment)
	if child == nil {
		child = &node{segment: segment, isWild: segment[0] == ':' || segment[0] == '*' || (segment[0] == '{' && segment[len(segment)-1] == '}')}
		n.children = append(n.children, child)
	}
	child.insert(pattern, segments, height+1)
}

func (n *node) search(segments []string, height int) *node {
	if len(segments) == height || strings.HasPrefix(n.segment, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	segment := segments[height]
	children := n.matchChildren(segment)
	for _, child := range children {
		result := child.search(segments, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
