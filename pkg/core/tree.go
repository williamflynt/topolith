package core

var EmptyTree = tree{}

// Tree is an interface that represents a tree structure for Item.
// It represents the hierarchy a World, where each Item can only be part of a single parent Item.
type Tree interface {
	Item() Item         // Item returns the Item of this Tree.
	Components() []Tree // Components returns the pieces of the Item of this Tree as Tree items.
	Parent() Tree       // Parent returns the parent of this Tree. An empty Tree is returned if this Tree has no parent.
	Empty() bool        // Empty returns whether this Tree has no Item and no Components.
}

// tree implements Tree.
type tree struct {
	item       *Item
	components []*tree
	parent     *tree
}

func (t *tree) Item() Item {
	if t.item == nil {
		return Item{}
	}
	return *t.item
}

func (t *tree) Components() []Tree {
	var components []Tree
	for _, c := range t.components {
		components = append(components, c)
	}
	return components
}

func (t *tree) Parent() Tree {
	if t.parent == nil {
		return &EmptyTree
	}
	return t.parent
}

func (t *tree) Empty() bool {
	return t.item == nil && (t.components == nil || len(t.components) == 0)
}
