package app

import (
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/williamflynt/topolith/pkg/errors"
	"github.com/williamflynt/topolith/pkg/grammar"
	"github.com/williamflynt/topolith/pkg/world"
	"strconv"
	"strings"
)

// TODO: Update all commands to use Flags where appropriate.

// Command is the interface that all app must implement.
//
// The CommandVerb is implicit in the **type** of the Command.
// Each type of Command will have its own unique structure and behavior.
// This is so that we make impossible states unrepresentable.
type Command interface {
	Execute(w world.World) (fmt.Stringer, error) // Execute runs the command on the given world.World. Return the resource object(s) or response, and an error if any.
	Undo(w world.World) error                    // Undo reverts the changes made by the command on the given world.World. Return an error if any. For non-mutating commands, this is a noop.
	fmt.Stringer
}

// CommandVerb represents the action to be performed.
// They are taken from the grammar.
type CommandVerb string

const (
	Create        CommandVerb = "create"          // Create command is used to create a new resource.
	Fetch         CommandVerb = "fetch"           // Fetch command is used to retrieve a resource.
	Set           CommandVerb = "set"             // Set command is used to update a resource.
	Clear         CommandVerb = "clear"           // Clear command is used to remove specific attributes from a resource.
	Delete        CommandVerb = "delete"          // Delete command is used to remove a resource.
	List          CommandVerb = "list"            // List command is used to retrieve a list of resources.
	Nest          CommandVerb = "nest"            // Nest command is used to nest a world.Item under another.
	Free          CommandVerb = "free"            // Free command is used to remove a world.Item from its parent and return it to the world.Tree root.
	Exists        CommandVerb = "exists"          // Exists command is used to check if a resource exists.
	ToQuery       CommandVerb = "to?"             // ToQuery command is used to retrieve all the world.Rel that have a relationship to the given world.Item.
	FromQuery     CommandVerb = "from?"           // FromQuery command is used to retrieve all the world.Rel that have a relationship from the given world.Item.
	InQuery       CommandVerb = "in?"             // InQuery command is used to retrieve all the world.Item that are nested under the given world.Item.
	CreateOrFetch CommandVerb = "create-or-fetch" // CreateOrFetch command is used to create a new resource if it doesn't exist, or fetch it if it does.
	CreateOrSet   CommandVerb = "create-or-set"   // CreateOrSet command is used to create a new resource if it doesn't exist, or set the given attributes if it does.
)

// CommandFlag represents a flag for a command.
type CommandFlag string

const (
	Strict  CommandFlag = "strict"  // Strict flag is used to indicate that the command should only be executed if the resource already exists, or to strictly interpret IDs (ie: not consider children/parents).
	Verbose CommandFlag = "verbose" // Verbose flag is used to indicate that the command should return more information.
	Ids     CommandFlag = "ids"     // Ids flag is used to indicate that the command should return the IDs of the resources, rather than the resources themselves.
)

// CommandTarget represents the type of resource we're executing the command on.
// Queries will return this type of resource, and Mutations will be executed over this type of resource.
// Commands that update the world.Tree should list the ItemTarget type.
//
// If we are executing an Exists command, we will return true/false rather than the resource itself.
//
// A request for Ids will return a list of IDs rather than the resources themselves.
type CommandTarget string

const (
	WorldTarget CommandTarget = "world"
	ItemTarget  CommandTarget = "item"
	RelTarget   CommandTarget = "rel"
)

// StringerList is a helper type to allow for a list of fmt.Stringer to be joined into a single string.
// The result of calling String() on a StringerList is a newline-separated string of the String() representation of each element.
type StringerList[T fmt.Stringer] []T

func (l StringerList[T]) String() string {
	strs := make([]string, len(l))
	for i, s := range l {
		strs[i] = s.String()
	}
	return strings.Join(strs, "\n")
}

type BoolStringer bool

func (b BoolStringer) String() string {
	return fmt.Sprintf("%t", b)
}

// CommandBase is a base struct for common command fields.
type CommandBase struct {
	InputAttributes grammar.InputAttributes
	ResourceType    CommandTarget
	Id              string
	Flags           mapset.Set[CommandFlag]
}

func (c *CommandBase) String() string {
	return c.InputAttributes.Raw
}

// --- COMMAND IMPLEMENTATIONS ---

// WorldFetchCommand represents a fetch command for the whole World.
type WorldFetchCommand struct {
	InputAttributes grammar.InputAttributes
}

func (c *WorldFetchCommand) Execute(w world.World) (fmt.Stringer, error) {
	return w, nil
}

func (c *WorldFetchCommand) Undo(w world.World) error {
	return nil
}

func (c *WorldFetchCommand) String() string {
	return c.InputAttributes.Raw
}

/* Item Commands */

// ItemCreateCommand represents a create command.
type ItemCreateCommand struct {
	CommandBase
	Params   world.ItemParams
	noCreate bool
}

func (c *ItemCreateCommand) Execute(w world.World) (fmt.Stringer, error) {
	if item, ok := w.ItemFetch(c.Id); ok {
		c.noCreate = true
		return item, nil
	}
	return w.ItemCreate(c.Id, c.Params).Item()
}

func (c *ItemCreateCommand) Undo(w world.World) error {
	if c.noCreate {
		return nil
	}
	return w.ItemDelete(c.Id).Err()
}

// ItemFetchCommand represents a fetch command for Item.
type ItemFetchCommand struct {
	CommandBase
}

func (c *ItemFetchCommand) Execute(w world.World) (fmt.Stringer, error) {
	item, ok := w.ItemFetch(c.Id)
	if !ok {
		return world.Item{}, errors.New("could not find Item").UseCode(errors.TopolithErrorNotFound).WithData(errors.KvPair{Key: "id", Value: c.Id})
	}
	return item, nil
}

func (c *ItemFetchCommand) Undo(w world.World) error {
	return nil
}

// ItemListCommand represents a list command for Item.
type ItemListCommand struct {
	CommandBase
	Limit int
}

func (c *ItemListCommand) Execute(w world.World) (fmt.Stringer, error) {
	items := w.ItemList(c.Limit)
	return StringerList[world.Item](items), nil
}

func (c *ItemListCommand) Undo(w world.World) error {
	return nil
}

// ItemSetCommand represents a set command for Item.
type ItemSetCommand struct {
	CommandBase
	Params    world.ItemParams
	oldParams world.ItemParams
	noSet     bool
}

func (c *ItemSetCommand) Execute(w world.World) (fmt.Stringer, error) {
	item, ok := w.ItemFetch(c.Id)
	if !ok {
		c.noSet = true
		return world.Item{}, errors.New("could not find Item").UseCode(errors.TopolithErrorNotFound).WithData(errors.KvPair{Key: "id", Value: c.Id})
	}
	c.oldParams.External = boolPtr(item.External)
	c.oldParams.Name = strPtr(item.Name)
	c.oldParams.Type = strPtr(world.StringFromItemType(item.Type))
	c.oldParams.Mechanism = strPtr(item.Mechanism)
	c.oldParams.Expanded = strPtr(item.Expanded)
	return w.ItemSet(c.Id, c.Params).Item()
}

func (c *ItemSetCommand) Undo(w world.World) error {
	if c.noSet {
		return nil
	}
	return w.ItemSet(c.Id, c.oldParams).Err()
}

// ItemClearCommand represents a clear command for Item - a modified set command.
type ItemClearCommand struct {
	CommandBase
	Params    world.ItemParams
	oldParams world.ItemParams
	noSet     bool
}

func (c *ItemClearCommand) Execute(w world.World) (fmt.Stringer, error) {
	item, ok := w.ItemFetch(c.Id)
	if !ok {
		c.noSet = true
		return world.Item{}, errors.New("could not find Item").UseCode(errors.TopolithErrorNotFound).WithData(errors.KvPair{Key: "id", Value: c.Id})
	}
	c.oldParams.External = boolPtr(item.External)
	c.oldParams.Name = strPtr(item.Name)
	c.oldParams.Type = strPtr(world.StringFromItemType(item.Type))
	c.oldParams.Mechanism = strPtr(item.Mechanism)
	c.oldParams.Expanded = strPtr(item.Expanded)
	return w.ItemSet(c.Id, c.Params).Item()
}

func (c *ItemClearCommand) Undo(w world.World) error {
	if c.noSet {
		return nil
	}
	return w.ItemSet(c.Id, c.oldParams).Err()
}

// ItemDeleteCommand represents a delete command for Item.
type ItemDeleteCommand struct {
	CommandBase
	oldParams world.ItemParams
	noDelete  bool
}

func (c *ItemDeleteCommand) Execute(w world.World) (fmt.Stringer, error) {
	item, ok := w.ItemFetch(c.Id)
	if !ok {
		c.noDelete = true
		return world.Item{}, nil
	}
	c.oldParams.External = boolPtr(item.External)
	c.oldParams.Name = strPtr(item.Name)
	c.oldParams.Type = strPtr(world.StringFromItemType(item.Type))
	c.oldParams.Mechanism = strPtr(item.Mechanism)
	c.oldParams.Expanded = strPtr(item.Expanded)
	return world.Item{}, w.ItemDelete(c.Id).Err()
}

func (c *ItemDeleteCommand) Undo(w world.World) error {
	if c.noDelete {
		return nil
	}
	return w.ItemCreate(c.Id, c.oldParams).Err()
}

// ItemNestCommand represents a nest command for Item.
type ItemNestCommand struct {
	CommandBase
	Ids          []string
	ParentId     string
	oldParentIds map[string]string
	noNest       map[string]bool
}

func (c *ItemNestCommand) Execute(w world.World) (fmt.Stringer, error) {
	oldParentIds := make(map[string]string)
	noNest := make(map[string]bool)
	errs := make([]error, 0)
	for _, id := range c.Ids {
		oldParentId, found := w.Parent(id)
		if !found {
			noNest[id] = true
			errs = append(errs, errors.New("could not find Item").UseCode(errors.TopolithErrorNotFound).WithData(errors.KvPair{Key: "id", Value: id}))
			continue
		}
		oldParentIds[id] = oldParentId // Empty string if root.
		if oldParentId == c.ParentId {
			noNest[id] = true
			continue
		}
		w.Nest(id, c.ParentId)
	}
	if len(errs) > 0 {
		return BoolStringer(false), errors.Join(errs...)
	}
	return BoolStringer(true), nil
}

func (c *ItemNestCommand) Undo(w world.World) error {
	for id, oldParentId := range c.oldParentIds {
		if oldParentId == "" {
			w.Free(id)
			continue
		}
		w.Nest(id, oldParentId)
	}
	return nil
}

// ItemFreeCommand represents a free command for Item.
type ItemFreeCommand struct {
	CommandBase
	Ids          []string
	oldParentIds map[string]string
}

func (c *ItemFreeCommand) Execute(w world.World) (fmt.Stringer, error) {
	oldParentIds := make(map[string]string)
	errs := make([]error, 0)

	for _, id := range c.Ids {
		oldParentId, found := w.Parent(id)
		if !found {
			errs = append(errs, errors.New("could not find Item").UseCode(errors.TopolithErrorNotFound).WithData(errors.KvPair{Key: "id", Value: id}))
			continue
		}
		oldParentIds[id] = oldParentId // Empty string if root.
		w.Free(id)
	}

	if len(errs) > 0 {
		return BoolStringer(false), errors.Join(errs...)
	}
	return BoolStringer(true), nil
}

func (c *ItemFreeCommand) Undo(w world.World) error {
	for id, oldParentId := range c.oldParentIds {
		if oldParentId == "" {
			continue
		}
		w.Nest(id, oldParentId)
	}
	return nil
}

// ItemExistsCommand represents an exists command for Item.
type ItemExistsCommand struct {
	CommandBase
}

func (c *ItemExistsCommand) Execute(w world.World) (fmt.Stringer, error) {
	_, ok := w.ItemFetch(c.Id)
	return BoolStringer(ok), nil
}

func (c *ItemExistsCommand) Undo(w world.World) error {
	return nil
}

// ItemCreateOrFetchCommand represents a create-or-fetch command for Item.
type ItemCreateOrFetchCommand struct {
	CommandBase
	noCreate bool
}

func (c *ItemCreateOrFetchCommand) Execute(w world.World) (fmt.Stringer, error) {
	item, ok := w.ItemFetch(c.Id)
	if ok {
		c.noCreate = true
		return item, nil
	}
	return w.ItemCreate(c.Id, world.ItemParams{}).Item()
}

func (c *ItemCreateOrFetchCommand) Undo(w world.World) error {
	if c.noCreate {
		return nil
	}
	return w.ItemDelete(c.Id).Err()
}

// ItemCreateOrSetCommand represents a create-or-set command for Item.
type ItemCreateOrSetCommand struct {
	CommandBase
	Params    world.ItemParams
	oldParams world.ItemParams
	noCreate  bool
}

func (c *ItemCreateOrSetCommand) Execute(w world.World) (fmt.Stringer, error) {
	if item, ok := w.ItemFetch(c.Id); ok {
		c.noCreate = true
		c.oldParams.External = boolPtr(item.External)
		c.oldParams.Name = strPtr(item.Name)
		c.oldParams.Type = strPtr(world.StringFromItemType(item.Type))
		c.oldParams.Mechanism = strPtr(item.Mechanism)
		c.oldParams.Expanded = strPtr(item.Expanded)
		return w.ItemSet(c.Id, c.Params).Item()
	}
	return w.ItemCreate(c.Id, c.Params).Item()
}

func (c *ItemCreateOrSetCommand) Undo(w world.World) error {
	if c.noCreate {
		return w.ItemSet(c.Id, c.oldParams).Err()
	}
	return w.ItemDelete(c.Id).Err()
}

// ItemComponentsListCommand represents a list command for the components of an Item.
type ItemComponentsListCommand struct {
	CommandBase
}

func (c *ItemComponentsListCommand) Execute(w world.World) (fmt.Stringer, error) {
	items, ok := w.ItemComponents(c.Id)
	if !ok {
		return StringerList[world.Item](items), errors.New("could not find Item").UseCode(errors.TopolithErrorNotFound).WithData(errors.KvPair{Key: "id", Value: c.Id})
	}
	return StringerList[world.Item](items), nil
}

func (c *ItemComponentsListCommand) Undo(w world.World) error {
	return nil
}

// ItemInQueryCommand represents an in-query command for Item.
type ItemInQueryCommand struct {
	CommandBase
	ParentId string
}

func (c *ItemInQueryCommand) Execute(w world.World) (fmt.Stringer, error) {
	strict := c.Flags.Contains(Strict)
	isInThere := w.In(c.Id, c.ParentId, strict)
	return BoolStringer(isInThere), nil
}

func (c *ItemInQueryCommand) Undo(w world.World) error {
	return nil
}

/* Rel Commands */

// RelCreateCommand represents a create command for Rel.
type RelCreateCommand struct {
	CommandBase
	ToId     string
	Params   world.RelParams
	noCreate bool
}

func (c *RelCreateCommand) Execute(w world.World) (fmt.Stringer, error) {
	if rels := w.RelFetch(c.Id, c.ToId, true); len(rels) > 0 {
		c.noCreate = true
		return world.Rel{}, nil
	}
	return w.RelCreate(c.Id, c.ToId, c.Params).Rel()
}

func (c *RelCreateCommand) Undo(w world.World) error {
	if c.noCreate {
		return nil
	}
	return w.RelDelete(c.Id, c.ToId).Err()
}

// RelFetchCommand represents a fetch command for Rel.
type RelFetchCommand struct {
	CommandBase
	ToId string
}

func (c *RelFetchCommand) Execute(w world.World) (fmt.Stringer, error) {
	strict := c.Flags.Contains(Strict)
	rels := w.RelFetch(c.Id, c.ToId, strict)
	if len(rels) == 0 {
		return world.Rel{}, errors.New("could not find Rel").UseCode(errors.TopolithErrorNotFound).WithData(errors.KvPair{Key: "id", Value: c.Id})
	}
	return rels[0], nil
}

func (c *RelFetchCommand) Undo(w world.World) error {
	return nil
}

// RelListCommand represents a list command for Rel.
type RelListCommand struct {
	CommandBase
	Limit int
}

func (c *RelListCommand) Execute(w world.World) (fmt.Stringer, error) {
	rels := w.RelList(c.Limit)
	return StringerList[world.Rel](rels), nil
}

func (c *RelListCommand) Undo(w world.World) error {
	return nil
}

// RelClearCommand represents a clear command for Rel - a modified set command.
type RelClearCommand struct {
	CommandBase
	ToId      string
	Params    world.RelParams
	oldParams world.RelParams
	noSet     bool
}

func (c *RelClearCommand) Execute(w world.World) (fmt.Stringer, error) {
	rels := w.RelFetch(c.Id, c.ToId, true)
	if len(rels) == 0 {
		c.noSet = true
		return world.Rel{}, errors.New("could not find Rel").UseCode(errors.TopolithErrorNotFound).WithData(errors.KvPair{Key: "id", Value: c.Id})
	}
	rel := rels[0]
	c.oldParams.Verb = strPtr(rel.Verb)
	c.oldParams.Mechanism = strPtr(rel.Mechanism)
	c.oldParams.Async = boolPtr(rel.Async)
	c.oldParams.Expanded = strPtr(rel.Expanded)
	return w.RelSet(c.Id, c.ToId, c.Params).Rel()
}

func (c *RelClearCommand) Undo(w world.World) error {
	if c.noSet {
		return nil
	}
	return w.RelSet(c.Id, c.ToId, c.oldParams).Err()
}

// RelExistsCommand represents an exists command for Rel.
type RelExistsCommand struct {
	CommandBase
	ToId string
}

func (c *RelExistsCommand) Execute(w world.World) (fmt.Stringer, error) {
	rels := w.RelFetch(c.Id, c.ToId, true)
	return BoolStringer(len(rels) > 0), nil
}

func (c *RelExistsCommand) Undo(w world.World) error {
	return nil
}

// RelToQueryCommand represents a to-query command for Rel.
type RelToQueryCommand struct {
	CommandBase
}

func (c *RelToQueryCommand) Execute(w world.World) (fmt.Stringer, error) {
	strict := c.Flags.Contains(Strict)
	rels := w.RelTo(c.Id, strict)
	return StringerList[world.Rel](rels), nil
}

func (c *RelToQueryCommand) Undo(w world.World) error {
	return nil
}

// RelFromQueryCommand represents a from-query command for Rel.
type RelFromQueryCommand struct {
	CommandBase
}

func (c *RelFromQueryCommand) Execute(w world.World) (fmt.Stringer, error) {
	strict := c.Flags.Contains(Strict)
	rels := w.RelFrom(c.Id, strict)
	return StringerList[world.Rel](rels), nil
}

func (c *RelFromQueryCommand) Undo(w world.World) error {
	return nil
}

// RelCreateOrFetchCommand represents a create-or-fetch command for Rel.
type RelCreateOrFetchCommand struct {
	CommandBase
	ToId     string
	noCreate bool
}

func (c *RelCreateOrFetchCommand) Execute(w world.World) (fmt.Stringer, error) {
	rels := w.RelFetch(c.Id, c.ToId, true)
	if len(rels) > 0 {
		c.noCreate = true
		return rels[0], nil
	}
	return w.RelCreate(c.Id, c.ToId, world.RelParams{}).Rel()
}

func (c *RelCreateOrFetchCommand) Undo(w world.World) error {
	if c.noCreate {
		return nil
	}
	return w.RelDelete(c.Id, c.ToId).Err()
}

// RelCreateOrSetCommand represents a create-or-set command for Rel.
type RelCreateOrSetCommand struct {
	CommandBase
	ToId      string
	Params    world.RelParams
	oldParams world.RelParams
	noCreate  bool
}

func (c *RelCreateOrSetCommand) Execute(w world.World) (fmt.Stringer, error) {
	rels := w.RelFetch(c.Id, c.ToId, true)
	if len(rels) > 0 {
		c.noCreate = true
		rel := rels[0]
		c.oldParams.Verb = strPtr(rel.Verb)
		c.oldParams.Mechanism = strPtr(rel.Mechanism)
		c.oldParams.Async = boolPtr(rel.Async)
		c.oldParams.Expanded = strPtr(rel.Expanded)
		return w.RelSet(c.Id, c.ToId, c.Params).Rel()
	}
	return w.RelCreate(c.Id, c.ToId, c.Params).Rel()
}

func (c *RelCreateOrSetCommand) Undo(w world.World) error {
	if c.noCreate {
		return w.RelSet(c.Id, c.ToId, c.oldParams).Err()
	}
	return w.RelDelete(c.Id, c.ToId).Err()
}

// RelSetCommand represents a set command for Rel.
type RelSetCommand struct {
	CommandBase
	ToId      string
	Params    world.RelParams
	oldParams world.RelParams
	noSet     bool
}

func (c *RelSetCommand) Execute(w world.World) (fmt.Stringer, error) {
	rels := w.RelFetch(c.Id, c.ToId, true)
	if len(rels) == 0 {
		c.noSet = true
		return world.Rel{}, errors.New("could not find Rel").UseCode(errors.TopolithErrorNotFound).WithData(errors.KvPair{Key: "id", Value: c.Id})
	}
	rel := rels[0]
	c.oldParams.Verb = strPtr(rel.Verb)
	c.oldParams.Mechanism = strPtr(rel.Mechanism)
	c.oldParams.Async = boolPtr(rel.Async)
	c.oldParams.Expanded = strPtr(rel.Expanded)
	return w.RelSet(c.Id, c.ToId, c.Params).Rel()
}

func (c *RelSetCommand) Undo(w world.World) error {
	if c.noSet {
		return nil
	}
	return w.RelSet(c.Id, c.ToId, c.oldParams).Err()
}

// RelDeleteCommand represents a delete command for Rel.
type RelDeleteCommand struct {
	CommandBase
	ToId      string
	oldParams world.RelParams
	noDelete  bool
}

func (c *RelDeleteCommand) Execute(w world.World) (fmt.Stringer, error) {
	rels := w.RelFetch(c.Id, c.Id, true)
	if len(rels) == 0 {
		c.noDelete = true
		return world.Rel{}, nil
	}
	rel := rels[0]
	c.oldParams.Verb = strPtr(rel.Verb)
	c.oldParams.Mechanism = strPtr(rel.Mechanism)
	c.oldParams.Async = boolPtr(rel.Async)
	c.oldParams.Expanded = strPtr(rel.Expanded)
	return world.Rel{}, w.RelDelete(c.Id, c.ToId).Err()
}

func (c *RelDeleteCommand) Undo(w world.World) error {
	if c.noDelete {
		return nil
	}
	return w.RelCreate(c.Id, c.Id, c.oldParams).Err()
}

// --- EXPORTED FUNCTIONS ---

// InputToCommand converts a grammar.InputAttributes to a Command.
func InputToCommand(input grammar.InputAttributes) (Command, error) {
	base := CommandBase{
		InputAttributes: input,
		ResourceType:    CommandTarget(input.ResourceType),
		Id:              input.ResourceId,
		Flags:           mapset.NewSet[CommandFlag](),
	}
	for _, flag := range input.Flags {
		base.Flags.Add(CommandFlag(flag))
	}

	switch base.ResourceType {
	case WorldTarget:
		return &WorldFetchCommand{InputAttributes: input}, nil
	case ItemTarget:
		return itemCommand(base, input)
	case RelTarget:
		return relCommand(base, input)
	default:
		return nil, errors.New("invalid resource type").UseCode(errors.TopolithErrorInvalid).WithData(errors.KvPair{Key: "resourceType", Value: input.ResourceType})
	}
}

// --- INTERNAL FUNCTIONS ---

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func itemCommand(base CommandBase, input grammar.InputAttributes) (Command, error) {
	switch CommandVerb(input.Verb) {
	case Create:
		return &ItemCreateCommand{CommandBase: base, Params: world.ItemParamsFromInput(input)}, nil
	case Fetch:
		return &ItemFetchCommand{CommandBase: base}, nil
	case List:
		return &ItemListCommand{CommandBase: base, Limit: limitFromInput(input)}, nil
	case Set:
		return &ItemSetCommand{CommandBase: base, Params: world.ItemParamsFromInput(input)}, nil
	case Clear:
		return &ItemClearCommand{CommandBase: base, Params: world.ItemParamsFromInput(input)}, nil
	case Delete:
		return &ItemDeleteCommand{CommandBase: base}, nil
	case Nest:
		return &ItemNestCommand{CommandBase: base, Ids: input.ResourceIds, ParentId: input.SecondaryIds[0], oldParentIds: make(map[string]string), noNest: make(map[string]bool)}, nil
	case Free:
		return &ItemFreeCommand{CommandBase: base, Ids: input.ResourceIds, oldParentIds: make(map[string]string)}, nil
	case Exists:
		return &ItemExistsCommand{CommandBase: base}, nil
	case InQuery:
		return &ItemInQueryCommand{CommandBase: base, ParentId: input.SecondaryIds[0]}, nil
	case CreateOrFetch:
		return &ItemCreateOrFetchCommand{CommandBase: base}, nil
	case CreateOrSet:
		return &ItemCreateOrSetCommand{CommandBase: base, Params: world.ItemParamsFromInput(input)}, nil
	default:
		return nil, errors.New("invalid verb").UseCode(errors.TopolithErrorInvalid).WithData(errors.KvPair{Key: "verb", Value: input.Verb}, errors.KvPair{Key: "resourceType", Value: input.ResourceType})
	}
}

func relCommand(base CommandBase, input grammar.InputAttributes) (Command, error) {
	switch CommandVerb(input.Verb) {
	case Create:
		return &RelCreateCommand{CommandBase: base, ToId: input.SecondaryIds[0], Params: world.RelParamsFromInput(input)}, nil
	case Fetch:
		return &RelFetchCommand{CommandBase: base, ToId: input.SecondaryIds[0]}, nil
	case List:
		return &RelListCommand{CommandBase: base, Limit: limitFromInput(input)}, nil
	case Set:
		return &RelSetCommand{CommandBase: base, ToId: input.SecondaryIds[0], Params: world.RelParamsFromInput(input)}, nil
	case Clear:
		return &RelClearCommand{CommandBase: base, ToId: input.SecondaryIds[0], Params: world.RelParamsFromInput(input)}, nil
	case Delete:
		return &RelDeleteCommand{CommandBase: base, ToId: input.SecondaryIds[0]}, nil
	case Exists:
		return &RelExistsCommand{CommandBase: base, ToId: input.SecondaryIds[0]}, nil
	case ToQuery:
		return &RelToQueryCommand{CommandBase: base}, nil
	case FromQuery:
		return &RelFromQueryCommand{CommandBase: base}, nil
	case CreateOrFetch:
		return &RelCreateOrFetchCommand{CommandBase: base, ToId: input.SecondaryIds[0]}, nil
	case CreateOrSet:
		return &RelCreateOrSetCommand{CommandBase: base, ToId: input.SecondaryIds[0], Params: world.RelParamsFromInput(input)}, nil
	default:
		return nil, errors.New("invalid verb").UseCode(errors.TopolithErrorInvalid).WithData(errors.KvPair{Key: "verb", Value: input.Verb}, errors.KvPair{Key: "resourceType", Value: input.ResourceType})
	}
}

func limitFromInput(input grammar.InputAttributes) int {
	v, ok := input.Params["limit"]
	if !ok {
		return 0
	}
	x, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return x
}
