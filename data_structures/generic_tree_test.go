package data

import "testing"

func TestGenericTree_Basic(t *testing.T) {
	tree := NewGenericTree("root")
	if tree.Root == nil {
		t.Fatal("Root should not be nil")
	}
	if tree.Root.Value != "root" {
		t.Errorf("Expected root value 'root', got '%v'", tree.Root.Value)
	}

	child := tree.Root.AddChild("child1")
	if len(tree.Root.Children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(tree.Root.Children))
	}
	if child.Value != "child1" {
		t.Errorf("Expected child value 'child1', got '%v'", child.Value)
	}
	if child.Parent != tree.Root {
		t.Error("Child's parent should be root")
	}

	grandchild := child.AddChild("grandchild")
	if len(child.Children) != 1 {
		t.Errorf("Expected 1 grandchild, got %d", len(child.Children))
	}
	if grandchild.Value != "grandchild" {
		t.Errorf("Expected grandchild value 'grandchild', got '%v'", grandchild.Value)
	}
	if grandchild.Parent != child {
		t.Error("Grandchild's parent should be child")
	}

	path := grandchild.Path()
	if len(path) != 3 || path[0].Value != "root" || path[1].Value != "child1" || path[2].Value != "grandchild" {
		t.Errorf("Path incorrect, got values: %v, want: [root child1 grandchild]", []string{path[0].Value, path[1].Value, path[2].Value})
	}
}
