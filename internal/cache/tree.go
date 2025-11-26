package cache

import "strings"

type Node struct {
	children map[string]*Node
	metadata map[string]any
	data     []byte
	Path     string
}

var ROOT_NODE = NewNode("root")

func NewNode(name string) *Node {
	return &Node{
		children: make(map[string]*Node),
		Path:     name,
	}
}

func (n *Node) calculateStats() (count int, totalSize int64) {
	if len(n.data) > 0 {
		count = 1
		totalSize = int64(len(n.data))
	}

	for _, child := range n.children {
		c, s := child.calculateStats()
		count += c
		totalSize += s
	}

	return count, totalSize
}

func (n *Node) GetData() []byte {
	return n.data
}

func (n *Node) IsLeaf() bool {
	if (n) == nil {
		return false
	}
	return len(n.children) == 0
}

func (n *Node) GetChild(name string) *Node {
	return n.children[name]
}

func (n *Node) GetChildren() map[string]*Node {
	return n.children
}

func (n *Node) GetChildrenAt(slug string) map[string]*Node {
	child := n.Lookup(slug)
	if child == nil {
		return make(map[string]*Node)
	}
	return child.children
}

func (n *Node) AddChild(child *Node) {
	if n.children == nil {
		n.children = make(map[string]*Node)
	}
	n.children[child.Path] = child
}

func (n *Node) Lookup(path string) *Node {
	if path == "" || path == "/" {
		return n
	}

	elements := strings.Split(path, "/")
	currentNode := n

	for _, elem := range elements {
		if elem == "" {
			continue
		}

		child := currentNode.GetChild(elem)
		if child == nil {
			return nil
		}
		currentNode = child
	}
	return currentNode
}

func (n *Node) Insert(path string, data []byte) {
	elements := strings.Split(path, "/")
	current := n

	for _, elem := range elements {
		if elem == "" {
			continue
		}

		child := current.GetChild(elem)
		if child == nil {
			child = NewNode(elem)
			current.AddChild(child)
		}
		current = child
	}
	current.data = data
}
