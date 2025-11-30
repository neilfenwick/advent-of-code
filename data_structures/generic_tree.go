package data

// GenericTree is a generic, idiomatic tree structure.
type GenericTree[T any] struct {
	Root *GenericTreeNode[T]
}

// GenericTreeNode represents a node in a generic tree.
type GenericTreeNode[T any] struct {
	Parent   *GenericTreeNode[T]
	Children []*GenericTreeNode[T]
	Value    T
}

// NewGenericTree creates a new tree with a root node.
func NewGenericTree[T any](value T) *GenericTree[T] {
	return &GenericTree[T]{Root: &GenericTreeNode[T]{Value: value}}
}

// AddChild adds a child to the given parent node.
func (n *GenericTreeNode[T]) AddChild(value T) *GenericTreeNode[T] {
	child := &GenericTreeNode[T]{Value: value, Parent: n}
	n.Children = append(n.Children, child)
	return child
}

// Path returns the path from the root to this node as a slice.
func (n *GenericTreeNode[T]) Path() []*GenericTreeNode[T] {
	var path []*GenericTreeNode[T]
	for curr := n; curr != nil; curr = curr.Parent {
		path = append([]*GenericTreeNode[T]{curr}, path...)
	}
	return path
}
