package core

import (
	mapset "github.com/deckarep/golang-set/v2"
	"testing"
)

// TestNilPointerDereference tests that a totally empty Tree doesn't cause a nil pointer dereference.
func TestNilPointerDereference(t *testing.T) {
	nilTree := &tree{}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Item method panicked: %v", r)
		}
	}()
	_ = nilTree.Item()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Components method panicked: %v", r)
		}
	}()
	_ = nilTree.Components()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Parent method panicked: %v", r)
		}
	}()
	_ = nilTree.Parent()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Empty method panicked: %v", r)
		}
	}()
	_ = nilTree.Empty()
}

func TestEmptyTree(t *testing.T) {
	if !emptyTree.Empty() {
		t.Error("expected emptyTree to be empty")
	}
	if emptyTree.Item().Id != "" {
		t.Error("expected emptyTree to have no item")
	}
	if !emptyTree.Components().IsEmpty() {
		t.Error("expected emptyTree to have no components")
	}
	if emptyTree.Parent() != emptyTree {
		t.Error("expected emptyTree to have itself as parent")
	}
}

func TestTree_Item(t *testing.T) {
	root := createSampleTree()
	if root.Item().Id != "root" {
		t.Errorf("expected root item ID to be 'root', got '%s'", root.Item().Id)
	}
	if root.components.ToSlice()[0].Item().Id != "child1" {
		t.Errorf("expected first child item ID to be 'child1', got '%s'", root.components.ToSlice()[0].Item().Id)
	}
	if root.components.ToSlice()[1].Item().Id != "child2" {
		t.Errorf("expected second child item ID to be 'child2', got '%s'", root.components.ToSlice()[1].Item().Id)
	}
}

func TestTree_Components(t *testing.T) {
	root := createSampleTree()
	components := root.components.ToSlice()
	if len(components) != 2 {
		t.Errorf("expected 2 components, got %d", len(components))
	}
	if components[0].Item().Id != "child1" {
		t.Errorf("expected first component item ID to be 'child1', got '%s'", components[0].Item().Id)
	}
	if components[1].Item().Id != "child2" {
		t.Errorf("expected second component item ID to be 'child2', got '%s'", components[1].Item().Id)
	}
}

func TestTree_Parent(t *testing.T) {
	root := createSampleTree()
	child1 := root.components.ToSlice()[0]
	child2 := root.components.ToSlice()[1]
	if child1.Parent().Item().Id != "root" {
		t.Errorf("expected parent of child1 to be 'root', got '%s'", child1.Parent().Item().Id)
	}
	if child2.Parent().Item().Id != "root" {
		t.Errorf("expected parent of child2 to be 'root', got '%s'", child2.Parent().Item().Id)
	}
	if root.Parent() != emptyTree {
		t.Errorf("expected parent of root to be emptyTree, got '%s'", root.Parent().Item().Id)
	}
}

func TestTree_Empty(t *testing.T) {
	root := createSampleTree()
	if root.Empty() {
		t.Error("expected root not to be empty")
	}
	child1 := root.components.ToSlice()[0]
	if child1.Empty() {
		t.Error("expected child1 not to be empty")
	}
	child2 := root.components.ToSlice()[1]
	if child2.Empty() {
		t.Error("expected child2 not to be empty")
	}
	emptyChild := &tree{parent: root}
	if !emptyChild.Empty() {
		t.Error("expected emptyChild to be empty")
	}
}

// Helper function to create a sample Tree for testing.
func createSampleTree() *tree {
	root := &tree{item: &Item{Id: "root"}}
	child1 := &tree{item: &Item{Id: "child1"}, parent: root}
	child2 := &tree{item: &Item{Id: "child2"}, parent: root}
	root.components = mapset.NewSet[Tree](child1, child2)
	return root
}
