package world

import (
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/williamflynt/topolith/pkg/errors"
	"github.com/williamflynt/topolith/pkg/grammar"
	"sort"
	"strings"
)

var emptyTree = newTree(nil, nil)

// Tree is an interface that represents a tree structure for Item.
// It represents the hierarchy a World, where each Item can only be part of a single parent Item.
type Tree interface {
	Item() Item                          // Item returns the Item of this Tree.
	AddOrMove(item *Item) error          // AddOrMove adds an Item to Tree.Components. If the Item already exists in a different Tree, it will be moved with all its Components.
	Has(id string, strict bool) bool     // Has returns whether the Item with the given ID is in this Tree. If not strict, search Component trees as well.
	Find(id string) (Tree, bool)         // Find returns the Tree with the given ID, if it exists in this Tree or its Components.
	Delete(id string)                    // Delete removes the Item with the given ID from this Tree. It will transfer Components to its parent.
	GetDescendantIds(id string) []string // GetDescendantIds returns all the descendant IDs for the Tree with Item.Id matching the given id. If the Item doesn't exist, returns an empty slice.
	Components() mapset.Set[Tree]        // Components returns the pieces of the Item of this Tree as Tree items.
	Parent() Tree                        // Parent returns the parent of this Tree. An empty Tree is returned if this Tree has no parent.
	Root() Tree                          // Root returns the root of this Tree.
	Empty() bool                         // Empty returns whether this Tree has no Item and no Components.
	fmt.Stringer
}

// TreeFromString returns a Tree from a string representation.
// It is the inverse of Tree.String().
//
// We can also get all the items in a World from parsing the Tree string!
//
// Example of a Tree string:
// `tree{nil::[tree{item "2" external=false::[tree{item "1" external=false::[]}]} tree{item "3" external=false::[]}]}`
func TreeFromString(s string) (Tree, map[string]Item, error) {
	if s == "" {
		return emptyTree, map[string]Item{}, nil
	}

	p, err := grammar.Parse(s)
	if err != nil {
		return nil, map[string]Item{}, errors.New("error parsing Tree").UseCode(errors.TopolithErrorInvalid).WithError(err).WithDescription("error parsing Tree").WithData(errors.KvPair{Key: "input", Value: s})
	}

	itemMap := make(map[string]Item)
	for _, itemStr := range p.ItemStrings {
		item, err := ItemFromString(itemStr)
		if err != nil {
			return nil, itemMap, err
		}
		itemMap[item.Id] = item
	}

	t, err := convertNodeToTree(p.Tree, itemMap, nil)
	return t, itemMap, err
}

func TreeEqual(t1, t2 Tree) bool {
	// Basic item comparison.
	if !ItemEqual(t1.Item(), t2.Item()) {
		return false
	}
	// Compare number of components.
	if t1.Components().Cardinality() != t2.Components().Cardinality() {
		return false
	}
	// Sort because set order is not guaranteed.
	components1Slice := t1.Components().ToSlice()
	components2Slice := t2.Components().ToSlice()
	sort.Slice(components1Slice, func(i, j int) bool {
		return components1Slice[i].Item().Id < components1Slice[j].Item().Id
	})
	sort.Slice(components2Slice, func(i, j int) bool {
		return components1Slice[i].Item().Id < components1Slice[j].Item().Id
	})
	// Compare recursively.
	for i, c1 := range components1Slice {
		c2 := components2Slice[i]
		if !ItemEqual(c1.Item(), c2.Item()) {
			return false
		}
		if !TreeEqual(c1, c2) {
			return false
		}
	}
	return true
}

// convertNodeToTree recursively converts a grammar.Node to a Tree.
func convertNodeToTree(node grammar.Node, itemMap map[string]Item, parent *tree) (*tree, error) {
	var itemPtr *Item
	if parent != nil {
		item, ok := itemMap[node.Id]
		if !ok {
			return parent, errors.New("Item not found in map").UseCode(errors.TopolithErrorNotFound).WithDescription("Item not found in map").WithData(errors.KvPair{Key: "id", Value: node.Id})
		}
		itemPtr = &item
	}
	delete(itemMap, node.Id) // Prevent cycles.
	t := &tree{
		item:       itemPtr,
		components: mapset.NewSet[Tree](),
		parent:     parent,
	}
	for _, childNode := range node.Children {
		childTree, err := convertNodeToTree(childNode, itemMap, t)
		if err != nil {
			return t, err
		}
		t.components.Add(childTree)
	}
	return t, nil
}

// tree implements Tree.
type tree struct {
	item       *Item
	components mapset.Set[Tree]
	parent     *tree
}

func (t *tree) String() string {
	treeString := "tree{"

	if t.item != nil {
		treeString += t.item.String()
	} else {
		treeString += "nil"
	}

	treeString += "::["
	componentTrees := make([]string, 0)
	for _, c := range t.components.ToSlice() {
		componentTrees = append(componentTrees, c.String())
	}
	if len(componentTrees) > 0 {
		treeString += strings.Join(componentTrees, " ")
	}
	treeString += "]}"
	return treeString
}

func (t *tree) AddOrMove(item *Item) error {
	if item == nil {
		return errors.New("cannot add nil Item to Tree").UseCode(errors.TopolithErrorInvalid)
	}
	if t.Has(item.Id, true) {
		// Already in this specific Tree - noop.
		return nil
	}
	if node, ok := t.Root().Find(item.Id); ok {
		// Already in a different Tree.
		// Just move the node to this Tree.
		node.Parent().Components().Remove(node)
		t.components.Add(node)
		return nil
	}
	// Not in any Tree. Make it and add.
	itemTree := newTree(item, t)
	t.components.Add(itemTree)
	return nil
}

func (t *tree) Has(id string, strict bool) bool {
	if id == "" {
		return false
	}
	inThisNode := t.item != nil && t.item.Id == id
	if inThisNode || strict {
		return inThisNode
	}
	for _, c := range t.components.ToSlice() {
		if c.Has(id, false) {
			return true
		}
	}
	return false
}

func (t *tree) Find(id string) (Tree, bool) {
	if id == "" {
		return nil, false
	}
	if t.item != nil && t.item.Id == id {
		return t, true
	}
	for _, c := range t.components.ToSlice() {
		if found, ok := c.Find(id); ok {
			return found, true
		}
	}
	return nil, false
}

func (t *tree) Delete(id string) {
	if id == "" {
		return
	}
	found, ok := t.Find(id)
	if !ok {
		// The item doesn't exist here.
		return
	}
	foundComponents := found.Components().ToSlice()
	found.Parent().Components().Remove(found)
	for _, c := range foundComponents {
		found.Parent().Components().Add(c)
	}
}

func (t *tree) GetDescendantIds(id string) []string {
	if id == "" {
		return []string{}
	}
	found, ok := t.Find(id)
	if !ok {
		return []string{}
	}
	descendantIds := make([]string, 0)
	for _, c := range found.Components().ToSlice() {
		descendantIds = append(descendantIds, c.Item().Id)
		descendantIds = append(descendantIds, c.GetDescendantIds(c.Item().Id)...)
	}
	return descendantIds
}

func (t *tree) Item() Item {
	if t.item == nil {
		return Item{}
	}
	return *t.item
}

func (t *tree) Components() mapset.Set[Tree] {
	return t.components
}

func (t *tree) Root() Tree {
	if t.parent == nil {
		return t
	}
	return t.parent.Root()
}

func (t *tree) Parent() Tree {
	if t.parent == nil {
		return emptyTree
	}
	return t.parent
}

func (t *tree) Empty() bool {
	return t.item == nil && (t.components == nil || t.components.IsEmpty())
}

// --- INTERNAL HELPERS ---

func newTree(item *Item, parent *tree) *tree {
	return &tree{
		item:       item,
		components: mapset.NewSet[Tree](),
		parent:     parent,
	}
}
