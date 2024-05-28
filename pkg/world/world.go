package world

import (
	"encoding/json"
	"fmt"
	"github.com/williamflynt/topolith/pkg/errors"
	"slices"
	"strconv"
	"strings"
)

const currentVersion = 1

// Info is an interface that represents the basic, descriptive information of a World.
type Info interface {
	Version() int     // Version is the version of the World.
	Id() string       // Id is the unique identifier of the World.
	Name() string     // Name is the display name of the World; meant for humans.
	Expanded() string // Expanded is the expanded description of the World.
}

// Operations is an interface that represents operations over the state of the World.
type Operations interface {
	ItemCreate(id string, params ItemParams) WorldWithItem // ItemCreate creates a new Item in the World, or retrieves it if already exists.
	ItemDelete(id string) World                            // ItemDelete deletes an Item from the World. If the item doesn't exist, noop.
	ItemFetch(id string) (Item, bool)                      // ItemFetch fetches an Item from the World. Returns an "okay" boolean, which is true only if the Item exists.
	ItemList(limit int) []Item                             // ItemList returns a list of Items in the World, up to the given limit. A 0 indicates no limit.
	ItemSet(id string, params ItemParams) WorldWithItem    // ItemSet sets the not-nil attributes from ItemParams on Item that has the given ID.

	RelCreate(fromId, toId string, params RelParams) WorldWithRel // RelCreate creates a new Rel in the World, or retrieves it if already exists. Returns the empty Rel if either Item doesn't exist.
	RelDelete(fromId, toId string) World                          // RelDelete deletes a Rel from the World. If the Rel doesn't exist, noop.
	RelFetch(fromId, toId string, strict bool) []Rel              // RelFetch fetches a Rel from the World. It will traverse the internal World Tree to find the first Rel that matches the fromId OR any descendent of the associated Item, and the toId or any descendent of the associated Item. If strict is true, it will only return the Rel if the fromId and toId match exactly.
	RelTo(toId string, strict bool) []Rel                         // RelTo fetches a Rel from the World. It will traverse the internal World Tree to find the first Rel that matches the fromId OR any descendent of the associated Item, and the toId or any descendent of the associated Item. If strict is true, it will only return the Rel if the fromId and toId match exactly.
	RelFrom(fromId string, strict bool) []Rel                     // RelFrom fetches a Rel from the World. It will traverse the internal World Tree to find the first Rel that matches the fromId OR any descendent of the associated Item, and the toId or any descendent of the associated Item. If strict is true, it will only return the Rel if the fromId and toId match exactly.
	RelList(limit int) []Rel                                      // RelList returns a list of Rels in the World, up to the given limit. A 0 indicates no limit.
	RelSet(fromId, toId string, params RelParams) WorldWithRel    // RelSet sets the not-nil attributes from RelParams on Rel that has the given fromId and toId.

	In(childId, parentId string, strict bool) bool // In checks if a child Item is nested anywhere under a parent Item. If strict is true, it will only return true if the childId and parentId match exactly.
	Parent(childId string) (string, bool)          // Parent returns the ID of the parent Item of the given child Item. An empty string is returned if the child Item has no parent. The okay boolean is false if the childId isn't found.
	Components(childId string) ([]string, bool)    // Components returns the IDs of the child Items of the given parent Item. An empty slice is returned if the parent Item has no children. The okay boolean is false if the parent Item isn't found.
	ItemParent(id string) (Item, bool)             // ItemParent returns the ID of the parent Item of the given child Item. An empty Item is returned if the child Item has no parent. The okay boolean is false if the childId isn't found.
	ItemComponents(id string) ([]Item, bool)       // ItemComponents returns the IDs of the child Items of the given parent Item. An empty slice is returned if the parent Item has no children. The okay boolean is false if the parent Item isn't found.
	Nest(childId, parentId string) WorldWithItem   // Nest nests a child Item under a parent Item. If the parent doesn't exist, noop.
	Free(childId string) WorldWithItem             // Free removes an Item from its parent to the root. If the Item doesn't exist, noop.

	Err() error // Err returns an error if the last operation failed, or nil if it succeeded.
}

// World is an interface that represents the state of the world.
type World interface {
	Info
	Operations
	fmt.Stringer
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

	Items map[string]Item `json:"items"` // Items is a map of ID string to related Item.
	Rels  map[string]Rel  `json:"rels"`  // Rels is a map of `Rel.From.Id` to related Rel.
	Tree  Tree            `json:"tree"`  // Tree is a tree representation of the World.

	latestItem *Item // latestItem is the last Item that was created or modified. This will be returned by the Item() method.
	latestRel  *Rel  // latestRel is the last Rel that was created or modified. This will be returned by the Rel() method.
	latestErr  error // latestErr is any error that occurred during the most recent operation.
}

func CreateWorld(name string) World {
	return &world{
		Version_:  currentVersion,
		Id_:       name,
		Name_:     name,
		Expanded_: "",
		Items:     make(map[string]Item),
		Rels:      make(map[string]Rel),
		Tree:      newTree(nil, nil),

		latestItem: &Item{},
		latestRel:  &Rel{},
	}
}

func (w *world) String() string {
	treeString := w.Tree.String()
	allRels := make([]string, 0)
	for _, rel := range w.Rels {
		allRels = append(allRels, rel.String())
	}
	return fmt.Sprintf("$$world\n%s\n%s\nendworld$$",
		treeString,
		strings.Join(allRels, "\n"),
	)
}

func (w *world) Version() int {
	return w.Version_
}

func (w *world) Id() string {
	return w.Id_
}

func (w *world) Name() string {
	return w.Name_
}

func (w *world) Expanded() string {
	return w.Expanded_
}

func (w *world) ItemCreate(id string, params ItemParams) WorldWithItem {
	w.resetLatestTrackers()
	if id == "" {
		w.latestErr = errors.New("id cannot be empty")
		return w
	}
	if existing, ok := w.Items[id]; ok {
		w.latestItem = &existing
		// Check params against existing, and create an error if they don't match.
		if err := equalItemParams(existing, params); err != nil {
			w.latestErr = err
		}
		return w
	}
	item := Item{
		Id: id,
	}
	w.Items[id] = item
	w.ItemSet(id, params) // After we set in the tracking map on World.
	w.latestItem = &item
	if err := w.Tree.AddOrMove(&item); err != nil {
		// This shouldn't happen if we're properly syncing the Items map with Tree...
		w.latestErr = err
	}
	return w
}

func (w *world) ItemDelete(id string) World {
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

func (w *world) ItemFetch(id string) (Item, bool) {
	item, ok := w.Items[id]
	return item, ok
}

func (w *world) ItemList(limit int) []Item {
	items := make([]Item, 0)
	for _, item := range w.Items {
		if limit > 0 && len(items) >= limit {
			break
		}
		items = append(items, item)
	}
	return items
}

func (w *world) ItemSet(id string, params ItemParams) WorldWithItem {
	w.resetLatestTrackers()
	item, ok := w.Items[id]
	if !ok {
		w.latestErr = errors.
			New("item not found").
			UseCode(errors.TopolithErrorNotFound).
			WithData(errors.KvPair{Key: "id", Value: id})
		return w
	}
	if params.Name != nil {
		item.Name = *params.Name
	}
	if params.Expanded != nil {
		item.Expanded = *params.Expanded
	}
	if params.External != nil {
		item.External = *params.External
	}
	if params.Type != nil {
		if iotaType, err := strconv.Atoi(*params.Type); err == nil {
			item.Type = ItemType(iotaType)
		} else {
			w.latestErr = err
		}
	}
	if params.Mechanism != nil {
		item.Mechanism = *params.Mechanism
	}
	w.Items[id] = item
	w.latestItem = &item
	return w
}

func (w *world) ItemParent(childId string) (Item, bool) {
	w.resetLatestTrackers()
	tree, ok := w.Tree.Find(childId)
	if !ok {
		return Item{}, false
	}
	parent := tree.Parent().Item()
	return parent, true
}

func (w *world) ItemComponents(parentId string) ([]Item, bool) {
	w.resetLatestTrackers()
	tree, ok := w.Tree.Find(parentId)
	if !ok {
		return []Item{}, false
	}
	components := tree.Components().ToSlice()
	items := make([]Item, len(components))
	for i, c := range components {
		items[i] = c.Item()
	}
	return items, true
}

func (w *world) RelCreate(fromId, toId string, params RelParams) WorldWithRel {
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
	existing, ok := w.Rels[relIdFromIds(fromId, toId)]
	if ok {
		w.latestRel = &existing
		// Check params against existing, and create an error if they don't match.
		if err := equalRelParams(existing, params); err != nil {
			w.latestErr = err
		}
		return w
	}
	rel := Rel{
		From: fromItem,
		To:   toItem,
	}
	w.Rels[rel.id()] = rel
	w.RelSet(fromId, toId, params) // After we set in the tracking map on World.
	w.latestRel = &rel
	return w
}

func (w *world) RelDelete(fromId, toId string) World {
	w.resetLatestTrackers()
	delete(w.Rels, relIdFromIds(fromId, toId))
	return w
}

func (w *world) RelFetch(fromId, toId string, strict bool) []Rel {
	w.resetLatestTrackers()
	if strict {
		rel, ok := w.Rels[relIdFromIds(fromId, toId)]
		if ok {
			return []Rel{rel}
		}
		return []Rel{}
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

func (w *world) RelTo(toId string, strict bool) []Rel {
	rels := make([]Rel, 0)
	rightIds := []string{toId}
	if strict {
		rightIds = append(w.Tree.GetDescendantIds(toId), toId)
	}
	for _, rel := range w.Rels {
		if slices.Contains(rightIds, rel.To.Id) {
			rels = append(rels, rel)
		}
	}
	return rels
}

func (w *world) RelFrom(fromId string, strict bool) []Rel {
	rels := make([]Rel, 0)
	leftIds := []string{fromId}
	if strict {
		leftIds = append(w.Tree.GetDescendantIds(fromId), fromId)
	}
	for k, rel := range w.Rels {
		if slices.Contains(leftIds, k) {
			rels = append(rels, rel)
		}
	}
	return rels
}

func (w *world) RelList(limit int) []Rel {
	rels := make([]Rel, 0)
	for _, rel := range w.Rels {
		if limit > 0 && len(rels) >= limit {
			break
		}
		rels = append(rels, rel)
	}
	return rels
}

func (w *world) RelSet(fromId, toId string, params RelParams) WorldWithRel {
	w.resetLatestTrackers()
	rel, ok := w.Rels[relIdFromIds(fromId, toId)]
	if !ok {
		w.latestErr = errors.
			New("rel not found").
			UseCode(errors.TopolithErrorNotFound).
			WithData(errors.KvPair{Key: "fromId", Value: fromId}, errors.KvPair{Key: "toId", Value: toId})
		return w
	}
	if params.Verb != nil {
		rel.Verb = *params.Verb
	}
	if params.Mechanism != nil {
		rel.Mechanism = *params.Mechanism
	}
	if params.Async != nil {
		rel.Async = *params.Async
	}
	if params.Expanded != nil {
		rel.Expanded = *params.Expanded
	}
	w.Rels[rel.id()] = rel
	w.latestRel = &rel
	return w
}

func (w *world) In(childId, parentId string, strict bool) bool {
	w.resetLatestTrackers()
	tree, ok := w.Tree.Find(parentId)
	if !ok {
		return false
	}
	return tree.Has(childId, strict)
}

func (w *world) Parent(childId string) (string, bool) {
	w.resetLatestTrackers()
	tree, ok := w.Tree.Find(childId)
	if !ok {
		return "", false
	}
	return tree.Parent().Item().Id, true
}

func (w *world) Components(parentId string) ([]string, bool) {
	w.resetLatestTrackers()
	tree, ok := w.Tree.Find(parentId)
	if !ok {
		return []string{}, false
	}
	components := tree.Components().ToSlice()
	ids := make([]string, len(components))
	for i, c := range components {
		ids[i] = c.Item().Id
	}
	return ids, true
}

func (w *world) Nest(childId, parentId string) WorldWithItem {
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

func (w *world) Free(childId string) WorldWithItem {
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
	w.latestErr = w.Tree.AddOrMove(&item)
	return w
}

func (w *world) Err() error {
	return w.latestErr
}

func (w *world) Item() (Item, error) {
	if w.latestItem == nil {
		return Item{}, w.latestErr
	}
	return *w.latestItem, w.latestErr
}

func (w *world) Rel() (Rel, error) {
	if w.latestRel == nil {
		return Rel{}, w.latestErr
	}
	return *w.latestRel, w.latestErr
}

// --- INTERNAL HELPERS ---

func equalItemParams(existing Item, params ItemParams) error {
	if (params.Name != nil && *params.Name != existing.Name) ||
		(params.Expanded != nil && *params.Expanded != existing.Expanded) ||
		(params.External != nil && *params.External != existing.External) ||
		(params.Type != nil && strconv.Itoa(int(existing.Type)) != *params.Type) ||
		(params.Mechanism != nil && *params.Mechanism != existing.Mechanism) {
		existingJson, _ := json.Marshal(existing)
		paramsJson, _ := json.Marshal(params)
		return errors.
			New("parameter mismatch").
			UseCode(errors.TopolithErrorConflict).
			WithData(
				errors.KvPair{Key: "object", Value: "Item"},
				errors.KvPair{Key: "existing", Value: string(existingJson)},
				errors.KvPair{Key: "params", Value: string(paramsJson)},
			)
	}
	return nil
}

func equalRelParams(existing Rel, params RelParams) error {
	if (params.Verb != nil && *params.Verb != existing.Verb) ||
		(params.Mechanism != nil && *params.Mechanism != existing.Mechanism) ||
		(params.Async != nil && *params.Async != existing.Async) ||
		(params.Expanded != nil && *params.Expanded != existing.Expanded) {
		existingJson, _ := json.Marshal(existing)
		paramsJson, _ := json.Marshal(params)
		return errors.
			New("parameter mismatch").
			UseCode(errors.TopolithErrorConflict).
			WithData(
				errors.KvPair{Key: "object", Value: "Rel"},
				errors.KvPair{Key: "existing", Value: string(existingJson)},
				errors.KvPair{Key: "params", Value: string(paramsJson)},
			)
	}
	return nil
}

// resetLatestTrackers resets the latestItem, latestRel, and latestErr fields.
// We do this before every operation to ensure that we don't accidentally return stale values
// from our Item() and Rel() methods.
func (w *world) resetLatestTrackers() {
	w.latestItem = &Item{}
	w.latestRel = &Rel{}
	w.latestErr = nil
}
