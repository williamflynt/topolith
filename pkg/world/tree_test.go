package world

import (
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"sort"
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
	if root.Item().Id != "" {
		t.Errorf("expected root item ID to be '', got '%s'", root.Item().Id)
	}
	rootComponentIds := mapset.NewSet[string]()
	for _, c := range root.components.ToSlice() {
		rootComponentIds.Add(c.Item().Id)
	}
	if !rootComponentIds.Contains("child1") {
		t.Error("expected 'child1' in components")
	}
	if !rootComponentIds.Contains("child2") {
		t.Error("expected 'child2' in components")
	}
}

func TestTree_Components(t *testing.T) {
	root := createSampleTree()
	components := root.components.ToSlice()
	if len(components) != 2 {
		t.Errorf("expected 2 components, got %d", len(components))
	}
	idSet := mapset.NewSet[string]("child1", "child2")
	if !idSet.Contains(components[0].Item().Id) {
		t.Errorf("expected '%s' in component IDs, not found", components[0].Item().Id)
	}
	if !idSet.Contains(components[1].Item().Id) {
		t.Errorf("expected '%s' in component IDs, not found", components[1].Item().Id)
	}
}

func TestTree_Parent(t *testing.T) {
	root := createSampleTree()
	child1 := root.components.ToSlice()[0]
	child2 := root.components.ToSlice()[1]
	if child1.Parent().Item().Id != "" {
		t.Errorf("expected parent of child1 to be 'root', got '%s'", child1.Parent().Item().Id)
	}
	if child2.Parent().Item().Id != "" {
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

func TestTreeFromString(t *testing.T) {
	simpleTree := "tree{nil::[tree{item \"2\"::[tree{item \"1\"::[]}]} tree{item \"3\"::[]}]}"
	parsed, _, err := TreeFromString(simpleTree)
	if parsed == nil {
		t.Fatal("expected parsed tree not to be nil")
	}
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// TEST the root Node.
	if parsed.Item().Id != "" {
		t.Errorf("expected root item ID to be 'nil', got '%s'", parsed.Item().Id)
	}
	components := parsed.Components().ToSlice()
	if len(components) != 2 {
		t.Errorf("expected 2 components, got %d", len(components))
	}

	// FIND Node for Item 2.
	tree2 := components[0]
	if tree2.Item().Id != "2" {
		tree2 = components[1]
	}
	if tree2.Item().Id != "2" {
		t.Fatalf("expected tree2 item ID to be '2', not found (got '%s')", tree2.Item().Id)
	}

	// FIND Node for Item 3.
	tree3 := components[0]
	if tree3.Item().Id != "3" {
		tree3 = components[1]
	}
	if tree3.Item().Id != "3" {
		t.Fatalf("expected tree3 item ID to be '3', not found (got '%s')", tree3.Item().Id)
	}

	// TEST NODE 2.
	if tree2.Components().IsEmpty() {
		t.Error("expected tree2 to have 1 component")
	}
	tree2_1 := tree2.Components().ToSlice()[0]
	if tree2_1.Item().Id != "1" {
		t.Errorf("expected tree2_1 item ID to be '1', got '%s'", tree2_1.Item().Id)
	}
	if !tree2_1.Components().IsEmpty() {
		t.Error("expected tree2_1 to have no components")
	}

	// TEST NODE 3.
	if !tree3.Components().IsEmpty() {
		t.Error("expected tree3 to have no components")
	}
}

func TestTreeSerde(t *testing.T) {
	root := createSampleTree()
	serialized := root.String()
	parsed, _, err := TreeFromString(serialized)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if root.Item().Id != parsed.Item().Id {
		t.Errorf("expected root item ID to be '%s', got '%s'", root.Item().Id, parsed.Item().Id)
	}
	if root.Components().Cardinality() != parsed.Components().Cardinality() {
		t.Errorf("expected %d components, got %d", root.Components().Cardinality(), parsed.Components().Cardinality())
	}
	serialized2 := parsed.String()
	if len(serialized) != len(serialized2) {
		t.Errorf("expected serialized strings to match lengths, they didn't")
	}
	// Hack around comparing structs with different memory addresses, and just sort and compare strings.
	rootComponents := root.Components().ToSlice()
	sort.Slice(rootComponents, func(i, j int) bool {
		return rootComponents[i].Item().Id < rootComponents[j].Item().Id
	})
	parsedComponents := parsed.Components().ToSlice()
	sort.Slice(parsedComponents, func(i, j int) bool {
		return parsedComponents[i].Item().Id < parsedComponents[j].Item().Id
	})
	if len(rootComponents) != len(parsedComponents) {
		t.Errorf("expected %d components, got %d", len(rootComponents), len(parsedComponents))
	}
	for i := range rootComponents {
		if rootComponents[i].Item().String() != parsedComponents[i].Item().String() {
			t.Errorf("component mismatch for ID '%s'", rootComponents[i].Item().Id)
			fmt.Println(rootComponents[i].Item().String())
			fmt.Println(parsedComponents[i].Item().String())
		}
	}
	if t.Failed() {
		fmt.Println(serialized)
		fmt.Println(serialized2)
	}
}

// Helper function to create a sample Tree for testing.
func createSampleTree() *tree {
	root := &tree{item: nil, components: mapset.NewSet[Tree]()}
	child1 := &tree{item: &Item{Id: "child1"}, parent: root, components: mapset.NewSet[Tree]()}
	child2 := &tree{item: &Item{Id: "child2"}, parent: root, components: mapset.NewSet[Tree]()}
	root.components = mapset.NewSet[Tree](child1, child2)
	return root
}
