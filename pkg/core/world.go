package core

import (
	"github.com/williamflynt/topolith/internal/errors"
	"slices"
)

const currentVersion = 1

// WorldInfo is an interface that represents the basic, descriptive information of a World.
type WorldInfo interface {
	Version() int     // Version is the version of the World.
	Id() string       // Id is the unique identifier of the World.
	Name() string     // Name is the display name of the World; meant for humans.
	Expanded() string // Expanded is the expanded description of the World.
}

// WorldOperations is an interface that represents operations over the state of the World.
type WorldOperations interface {
	ItemCreate(id string) WorldWithItem                     // ItemCreate creates a new Item in the World, or retrieves it if already exists.
	ItemDelete(id string) World                             // ItemDelete deletes an Item from the World. If the item doesn't exist, noop.
	ItemFetch(id string) (Item, bool)                       // ItemFetch fetches an Item from the World. Returns an "okay" boolean, which is true only if the Item exists.
	ItemSetExternal(id string, external bool) WorldWithItem // ItemSetExternal sets Item.External. Returns the empty Item if the Item doesn't exist.
	ItemSetType(id string, itemType string) WorldWithItem   // ItemSetType sets Item.Type. Returns the empty Item if the Item doesn't exist.
	ItemSetName(id, name string) WorldWithItem              // ItemSetName sets Item.Name. Returns the empty Item if the Item doesn't exist.
	ItemSetMechanism(id, mechanism string) WorldWithItem    // ItemSetMechanism sets Item.Mechanism. Returns the empty Item if the Item doesn't exist.
	ItemSetExpanded(id, expanded string) WorldWithItem      // ItemSetExpanded sets Item.Expanded. Returns the empty Item if the Item doesn't exist.

	RelCreate(fromId, toId string) WorldWithRel                  // RelCreate creates a new Rel in the World, or retrieves it if already exists. Returns the empty Rel if either Item doesn't exist.
	RelDelete(fromId, toId string) World                         // RelDelete deletes a Rel from the World. If the Rel doesn't exist, noop.
	RelFetch(fromId, toId string, strict bool) []Rel             // RelFetch fetches a Rel from the World. It will traverse the internal World Tree to find the first Rel that matches the fromId OR any descendent of the associated Item, and the toId or any descendent of the associated Item. If strict is true, it will only return the Rel if the fromId and toId match exactly.
	RelSetVerb(fromId, toId, verb string) WorldWithRel           // RelSetVerb sets Rel.Verb. Returns the empty Rel if the Rel doesn't exist.
	RelSetMechanism(fromId, toId, mechanism string) WorldWithRel // RelSetMechanism sets Rel.Mechanism. Returns the empty Rel if the Rel doesn't exist.
	RelSetAsync(fromId, toId string, async bool) WorldWithRel    // RelSetAsync sets Rel.Async. Returns the empty Rel if the Rel doesn't exist.
	RelSetExpanded(fromId, toId, expanded string) WorldWithRel   // RelSetExpanded sets Rel.Expanded. Returns the empty Rel if the Rel doesn't exist.

	In(childId, parentId string, strict bool) bool // In checks if a child Item is nested anywhere under a parent Item. If strict is true, it will only return true if the childId and parentId match exactly.
	Nest(childId, parentId string) WorldWithItem   // Nest nests a child Item under a parent Item. If the parent doesn't exist, noop.
	Free(childId, parentId string) WorldWithItem

	Undo() World // Undo reverses the last operation on the World. If there are no operations to undo, noop and return the same World.
	Redo() World // Redo executes the most recently reversed operation on the World. If there are no operations to redo, noop and return the same World.
	Err() error  // Err returns an error if the last operation failed, or nil if it succeeded.
}

// World is an interface that represents the state of the world.
type World interface {
	WorldInfo
	WorldOperations
}

// WorldWithItem is an interface that represents a World with a possible Item from the last operation.
type WorldWithItem interface {
	World
	Item() (Item, error)
}

// WorldWithRel is an interface that represents a World with a possible Rel from the last operation.
type WorldWithRel interface {
	World
	Rel() (Rel, error)
}

// WorldWithBoth is an interface that represents a World with a possible Item and Rel from the last operation.
type WorldWithBoth interface {
	WorldWithItem
	WorldWithRel
}

// world is the internal implementation of the WorldWithBoth interface.
// It necessarily implements World.
type world struct {
	Version_  int    `json:"version"`  // Version_ is the version of the World.
	Id_       string `json:"id"`       // Id_ is the unique identifier of the World.
	Name_     string `json:"name"`     // Name_ is the display name of the World; meant for humans.
	Expanded_ string `json:"expanded"` // Expanded_ is the expanded description of the World.

	Items   map[string]Item `json:"items"`   // Items is a map of ID string to related Item.
	Rels    map[string]Rel  `json:"rels"`    // Rels is a map of `Rel.From.Id` to related Rel.
	History []Command       `json:"history"` // History is a list of commands that have been executed.
	Tree    Tree            `json:"tree"`    // Tree is a tree representation of the World.

	historyIdx int   // historyIdx is the index of the last executed command in the History list.
	latestItem *Item // latestItem is the last Item that was created or modified. This will be returned by the Item() method.
	latestRel  *Rel  // latestRel is the last Rel that was created or modified. This will be returned by the Rel() method.
	latestErr  error // latestErr is any error that occurred during the most recent operation.
}

func CreateWorld(name string) World {
	return world{
		Version_:  currentVersion,
		Id_:       name,
		Name_:     name,
		Expanded_: "",
		Items:     make(map[string]Item),
		Rels:      make(map[string]Rel),
		History:   make([]Command, 0),
		Tree:      newTree(nil, nil),
	}
}

func (w world) Version() int {
	return w.Version_
}

func (w world) Id() string {
	return w.Id_
}

func (w world) Name() string {
	return w.Name_
}

func (w world) Expanded() string {
	return w.Expanded_
}

func (w world) ItemCreate(id string) WorldWithItem {
	w.resetLatestTrackers()
	if id == "" {
		w.latestErr = errors.New("id cannot be empty")
		return w
	}
	if existing, ok := w.Items[id]; ok {
		w.latestItem = &existing
		return w
	}
	item := Item{
		Id: id,
	}
	w.Items[id] = item
	w.latestItem = &item
	if err := w.Tree.AddOrMove(&item); err != nil {
		// This shouldn't happen if we're properly syncing the Items map with Tree...
		w.latestErr = err
	}
	return w
}

func (w world) ItemDelete(id string) World {
	w.resetLatestTrackers()
	if id == "" {
		w.latestErr = errors.New("id cannot be empty")
		return w
	}
	delete(w.Items, id)
	w.Tree.Delete(id)
	for k, rel := range w.Rels {
		if rel.From.Id == id || rel.To.Id == id {
			delete(w.Rels, k)
		}
	}
	return w
}

func (w world) ItemFetch(id string) (Item, bool) {
	item, ok := w.Items[id]
	return item, ok
}

func (w world) ItemSetExternal(id string, external bool) WorldWithItem {
	w.resetLatestTrackers()
	item, ok := w.ItemFetch(id)
	if !ok {
		w.latestErr = errors.
			New("item not found").
			UseCode(errors.TopolithErrorNotFound).
			WithData(errors.KvPair{Key: "id", Value: id})
		return w
	}
	item.External = external
	w.Items[id] = item
	w.latestItem = &item
	return w
}

func (w world) ItemSetType(id string, itemType string) WorldWithItem {
	w.resetLatestTrackers()
	item, ok := w.ItemFetch(id)
	if !ok {
		w.latestErr = errors.
			New("item not found").
			UseCode(errors.TopolithErrorNotFound).
			WithData(errors.KvPair{Key: "id", Value: id})
		return w
	}
	w.latestItem = &item
	item.Type = itemTypeFromString(itemType)
	w.Items[id] = item
	w.latestItem = &item
	return w
}

func (w world) ItemSetName(id, name string) WorldWithItem {
	w.resetLatestTrackers()
	item, ok := w.ItemFetch(id)
	if !ok {
		w.latestErr = errors.
			New("item not found").
			UseCode(errors.TopolithErrorNotFound).
			WithData(errors.KvPair{Key: "id", Value: id})
		return w
	}
	item.Name = name
	w.Items[id] = item
	w.latestItem = &item
	return w
}

func (w world) ItemSetMechanism(id, mechanism string) WorldWithItem {
	w.resetLatestTrackers()
	item, ok := w.ItemFetch(id)
	if !ok {
		w.latestErr = errors.
			New("item not found").
			UseCode(errors.TopolithErrorNotFound).
			WithData(errors.KvPair{Key: "id", Value: id})
		return w
	}
	item.Mechanism = mechanism
	w.Items[id] = item
	w.latestItem = &item
	return w
}

func (w world) ItemSetExpanded(id, expanded string) WorldWithItem {
	w.resetLatestTrackers()
	item, ok := w.ItemFetch(id)
	if !ok {
		w.latestErr = errors.
			New("item not found").
			UseCode(errors.TopolithErrorNotFound).
			WithData(errors.KvPair{Key: "id", Value: id})
		return w
	}
	item.Expanded = expanded
	w.Items[id] = item
	w.latestItem = &item
	return w
}

func (w world) RelCreate(fromId, toId string) WorldWithRel {
	w.resetLatestTrackers()
	fromItem, ok := w.ItemFetch(fromId)
	if !ok {
		w.latestErr = errors.
			New("fromId for Rel not found").
			UseCode(errors.TopolithErrorNotFound).
			WithData(errors.KvPair{Key: "fromId", Value: fromId})
		return w
	}
	toItem, ok := w.ItemFetch(toId)
	if !ok {
		w.latestErr = errors.
			New("fromId for Rel not found").
			UseCode(errors.TopolithErrorNotFound).
			WithData(errors.KvPair{Key: "fromId", Value: fromId})
		return w
	}
	rel := Rel{
		From: fromItem,
		To:   toItem,
	}
	w.Rels[rel.id()] = rel
	w.latestRel = &rel
	return w
}

func (w world) RelDelete(fromId, toId string) World {
	w.resetLatestTrackers()
	delete(w.Rels, relIdFromIds(fromId, toId))
	return w
}

func (w world) RelFetch(fromId, toId string, strict bool) []Rel {
	w.resetLatestTrackers()
	if strict {
		return []Rel{w.Rels[relIdFromIds(fromId, toId)]}
	}
	rels := make([]Rel, 0)
	leftIds := append(w.Tree.GetDescendantIds(fromId), fromId)
	rightIds := append(w.Tree.GetDescendantIds(toId), toId)
	for _, rel := range w.Rels {
		if slices.Contains(leftIds, rel.From.Id) && slices.Contains(rightIds, rel.To.Id) {
			rels = append(rels, rel)
		}
	}
	return rels
}

func (w world) RelSetVerb(fromId, toId, verb string) WorldWithRel {
	w.resetLatestTrackers()
	rel, ok := w.Rels[relIdFromIds(fromId, toId)]
	if !ok {
		w.latestErr = errors.
			New("rel not found").
			UseCode(errors.TopolithErrorNotFound).
			WithData(errors.KvPair{Key: "fromId", Value: fromId}, errors.KvPair{Key: "toId", Value: toId})
		return w
	}
	rel.Verb = verb
	w.Rels[rel.id()] = rel
	w.latestRel = &rel
	return w
}

func (w world) RelSetMechanism(fromId, toId, mechanism string) WorldWithRel {
	w.resetLatestTrackers()
	rel, ok := w.Rels[relIdFromIds(fromId, toId)]
	if !ok {
		w.latestErr = errors.
			New("rel not found").
			UseCode(errors.TopolithErrorNotFound).
			WithData(errors.KvPair{Key: "fromId", Value: fromId}, errors.KvPair{Key: "toId", Value: toId})
		return w
	}
	rel.Mechanism = mechanism
	w.Rels[rel.id()] = rel
	w.latestRel = &rel
	return w
}

func (w world) RelSetAsync(fromId, toId string, async bool) WorldWithRel {
	w.resetLatestTrackers()
	rel, ok := w.Rels[relIdFromIds(fromId, toId)]
	if !ok {
		w.latestErr = errors.
			New("rel not found").
			UseCode(errors.TopolithErrorNotFound).
			WithData(errors.KvPair{Key: "fromId", Value: fromId}, errors.KvPair{Key: "toId", Value: toId})
		return w
	}
	rel.Async = async
	w.Rels[rel.id()] = rel
	w.latestRel = &rel
	return w
}

func (w world) RelSetExpanded(fromId, toId, expanded string) WorldWithRel {
	w.resetLatestTrackers()
	rel, ok := w.Rels[relIdFromIds(fromId, toId)]
	if !ok {
		w.latestErr = errors.
			New("rel not found").
			UseCode(errors.TopolithErrorNotFound).
			WithData(errors.KvPair{Key: "fromId", Value: fromId}, errors.KvPair{Key: "toId", Value: toId})
		return w
	}
	rel.Expanded = expanded
	w.Rels[rel.id()] = rel
	w.latestRel = &rel
	return w
}

func (w world) In(childId, parentId string, strict bool) bool {
	w.resetLatestTrackers()
	tree, ok := w.Tree.Find(parentId)
	if !ok {
		return false
	}
	return tree.Has(childId, strict)
}

func (w world) Nest(childId, parentId string) WorldWithItem {
	w.resetLatestTrackers()
	item, ok := w.ItemFetch(childId)
	if !ok {
		w.latestErr = errors.
			New("childId for Nest not found").
			UseCode(errors.TopolithErrorNotFound).
			WithData(errors.KvPair{Key: "childId", Value: childId})
		return w
	}
	w.latestItem = &item
	if _, ok := w.ItemFetch(parentId); !ok {
		w.latestErr = errors.
			New("parentId for Nest not found").
			UseCode(errors.TopolithErrorNotFound).
			WithData(errors.KvPair{Key: "parentId", Value: parentId})
		return w
	}
	tree, ok := w.Tree.Find(parentId)
	if !ok {
		// The parent Item exists, but its entry in our World Tree doesn't.
		// This shouldn't happen, so we need to error out fast to fix the underlying issue.
		// That's why we don't just add the Item (which does exist) to the Tree.
		w.latestErr = errors.
			New("parentId for Nest not found in Tree").
			UseCode(errors.TopolithErrorBadSyncState).
			WithData(errors.KvPair{Key: "parentId", Value: parentId})
		return w
	}
	w.latestErr = tree.AddOrMove(&item)
	return w
}

func (w world) Free(childId, parentId string) WorldWithItem {
	w.resetLatestTrackers()
	item, ok := w.ItemFetch(childId)
	if !ok {
		w.latestErr = errors.
			New("childId for Free not found").
			UseCode(errors.TopolithErrorNotFound).
			WithData(errors.KvPair{Key: "childId", Value: childId})
		return w
	}
	w.latestItem = &item
	if _, ok := w.ItemFetch(parentId); !ok {
		w.latestErr = errors.
			New("parentId for Free not found").
			UseCode(errors.TopolithErrorNotFound).
			WithData(errors.KvPair{Key: "parentId", Value: parentId})
		return w
	}
	if _, ok = w.Tree.Find(parentId); !ok {
		// The parent Item exists, but its entry in our World Tree doesn't.
		// This shouldn't happen - see note on Nest.
		w.latestErr = errors.
			New("parentId for Free not found in Tree").
			UseCode(errors.TopolithErrorBadSyncState).
			WithData(errors.KvPair{Key: "parentId", Value: parentId})
		return w
	}
	w.latestErr = w.Tree.AddOrMove(&item)
	return w
}

func (w world) Undo() World {
	//TODO implement me
	panic("implement me")
}

func (w world) Redo() World {
	//TODO implement me
	panic("implement me")
}

func (w world) Err() error {
	return w.latestErr
}

func (w world) Item() (Item, error) {
	return *w.latestItem, w.latestErr
}

func (w world) Rel() (Rel, error) {
	return *w.latestRel, w.latestErr
}

// --- INTERNAL HELPERS ---

// resetLatestTrackers resets the latestItem, latestRel, and latestErr fields.
// We do this before every operation to ensure that we don't accidentally return stale values
// from our Item() and Rel() methods.
func (w world) resetLatestTrackers() {
	w.latestItem = &Item{}
	w.latestRel = &Rel{}
	w.latestErr = nil
}
