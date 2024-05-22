package core

import (
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/williamflynt/topolith/internal/errors"
)

var emptyTree = newTree(nil, nil)

// Tree is an interface that represents a tree structure for Item.
// It represents the hierarchy a World, where each Item can only be part of a single parent Item.
type Tree interface {
	Item() Item                     // Item returns the Item of this Tree.
	Add(item *Item) error           // Add adds an Item to Tree.Components, after validating that the Item is not already in the Tree.
	In(id string, strict bool) bool // In returns whether the Item with the given ID is in this Tree. If not strict, search Component trees as well.
	Find(id string) (Tree, bool)    // Find returns the Tree with the given ID in this Tree or Components.
	Delete(id string)               // Delete removes the Item with the given ID from this Tree. It will transfer Components to its parent.
	Components() mapset.Set[Tree]   // Components returns the pieces of the Item of this Tree as Tree items.
	Parent() Tree                   // Parent returns the parent of this Tree. An empty Tree is returned if this Tree has no parent.
	Empty() bool                    // Empty returns whether this Tree has no Item and no Components.
}

// tree implements Tree.
type tree struct {
	item       *Item
	components mapset.Set[Tree]
	parent     *tree
}

func (t *tree) Add(item *Item) error {
	if item == nil {
		return errors.New("cannot add nil item to tree").UseCode(errors.TopolithErrorInvalid)
	}
	if t.In(item.Id, true) {
		// Already in this specific Tree - noop.
		return nil
	}
	if t.In(item.Id, false) {
		// Already in a different Tree.
		return errors.New("item with id already in tree").UseCode(errors.TopolithErrorConflict)
	}
	itemTree := newTree(item, t)
	t.components.Add(itemTree)
	return nil
}

func (t *tree) In(id string, strict bool) bool {
	if id == "" {
		return false
	}
	inThisNode := t.item != nil && t.item.Id == id
	if inThisNode || strict {
		return inThisNode
	}
	for _, c := range t.components.ToSlice() {
		if c.In(id, false) {
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

func (t *tree) Item() Item {
	if t.item == nil {
		return Item{}
	}
	return *t.item
}

func (t *tree) Components() mapset.Set[Tree] {
	return t.components
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
