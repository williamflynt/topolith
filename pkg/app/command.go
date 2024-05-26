package app

import (
	"fmt"
	"github.com/williamflynt/topolith/pkg/errors"
	"github.com/williamflynt/topolith/pkg/grammar"
	"github.com/williamflynt/topolith/pkg/world"
	"regexp"
	"strings"
	"unicode"
)

var kvPattern = regexp.MustCompile(`\b(\w+)="?(\w+)"?\b`)

// CommandVerb represents the action to be performed.
type CommandVerb string

const (
	Create CommandVerb = "create"
	Set    CommandVerb = "set"
	Delete CommandVerb = "delete"
	Nest   CommandVerb = "nest"
	Free   CommandVerb = "free"
)

// Command is the interface that all app must implement.
type Command interface {
	Execute(w world.World) (fmt.Stringer, error)
	Undo(w world.World) error
	String() string
}

// CommandTarget represents the type of resource.
type CommandTarget string

const (
	ItemTarget CommandTarget = "item"
	RelTarget  CommandTarget = "rel"
)

// CommandBase is a base struct for common command fields.
type CommandBase struct {
	ResourceType CommandTarget
	Id           string
}

// --- COMMAND IMPLEMENTATIONS ---

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
	// TODO: Do this for all the other commands, too.
	return w.ItemCreate(c.Id, c.Params).Item()
}

func (c *ItemCreateCommand) Undo(w world.World) error {
	if c.noCreate {
		return nil
	}
	return w.ItemDelete(c.Id).Err()
}

func (c *ItemCreateCommand) String() string {
	return fmt.Sprintf(`%s %s "%s" %s`, c.ResourceType, Create, c.Id, itemParamsToString(c.Params))
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

func (c *ItemSetCommand) String() string {
	return fmt.Sprintf(`%s "%s" %s`, c.ResourceType, c.Id, itemParamsToString(c.Params))
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

func (c *ItemDeleteCommand) String() string {
	return fmt.Sprintf(`%s %s "%s"`, c.ResourceType, Delete, c.Id)
}

// ItemNestCommand represents a nest command for Item.
type ItemNestCommand struct {
	CommandBase
	ParentId    string
	oldParentId string
	noNest      bool
}

func (c *ItemNestCommand) Execute(w world.World) (fmt.Stringer, error) {
	oldParentId, found := w.Parent(c.Id)
	if !found {
		c.noNest = true
		return world.Item{}, errors.New("could not find Item").UseCode(errors.TopolithErrorNotFound).WithData(errors.KvPair{Key: "id", Value: c.Id})
	}
	c.oldParentId = oldParentId
	ww := w.Nest(c.Id, c.ParentId)
	return ww.Item()
}

func (c *ItemNestCommand) Undo(w world.World) error {
	if c.noNest {
		return nil
	}
	if c.oldParentId == "" {
		return w.Free(c.Id).Err()
	}
	return w.Nest(c.Id, c.oldParentId).Err()
}

func (c *ItemNestCommand) String() string {
	return fmt.Sprintf(`%s "%s" in "%s"`, Nest, c.Id, c.ParentId)
}

// ItemFreeCommand represents a free command for Item.
type ItemFreeCommand struct {
	CommandBase
	oldParentId string
}

func (c *ItemFreeCommand) Execute(w world.World) (fmt.Stringer, error) {
	oldParentId, found := w.Parent(c.Id)
	if !found {
		return world.Item{}, errors.New("could not find Item").UseCode(errors.TopolithErrorNotFound).WithData(errors.KvPair{Key: "id", Value: c.Id})
	}
	c.oldParentId = oldParentId
	return w.Free(c.Id).Item()
}

func (c *ItemFreeCommand) Undo(w world.World) error {
	if c.oldParentId == "" {
		return nil
	}
	return w.Nest(c.Id, c.oldParentId).Err()
}

func (c *ItemFreeCommand) String() string {
	return fmt.Sprintf(`%s "%s"`, Free, c.Id)
}

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

func (c *RelCreateCommand) String() string {
	return fmt.Sprintf(`%s %s "%s" "%s" %s`, c.ResourceType, Create, c.Id, c.ToId, relParamsToString(c.Params))
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

func (c *RelSetCommand) String() string {
	return fmt.Sprintf(`%s "%s" "%s" %s`, c.ResourceType, c.Id, c.ToId, relParamsToString(c.Params))
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

func (c *RelDeleteCommand) String() string {
	return fmt.Sprintf(`%s %s "%s" "%s"`, c.ResourceType, Delete, c.Id, c.ToId)
}

// --- EXPORTED FUNCTIONS ---

func MakeItemCreateCommand(id string, params world.ItemParams) (Command, error) {
	return &ItemCreateCommand{
		CommandBase: CommandBase{
			ResourceType: ItemTarget,
			Id:           id,
		},
		Params: params,
	}, nil
}

func MakeItemSetCommand(id string, params world.ItemParams) (Command, error) {
	return &ItemSetCommand{
		CommandBase: CommandBase{
			ResourceType: ItemTarget,
			Id:           id,
		},
		Params:    params,
		oldParams: world.ItemParams{},
	}, nil
}

func MakeItemDeleteCommand(id string) (Command, error) {
	return &ItemDeleteCommand{
		CommandBase: CommandBase{
			ResourceType: ItemTarget,
			Id:           id,
		},
		oldParams: world.ItemParams{},
	}, nil
}

func MakeItemFreeCommand(id string) (Command, error) {
	return &ItemFreeCommand{
		CommandBase: CommandBase{
			ResourceType: ItemTarget,
			Id:           id,
		},
	}, nil
}

func MakeNestCommand(id string, parentId string) (Command, error) {
	return &ItemNestCommand{
		CommandBase: CommandBase{
			ResourceType: ItemTarget,
			Id:           id,
		},
		ParentId: parentId,
	}, nil
}

func MakeRelCreateCommand(fromId string, toId string, params world.RelParams) (Command, error) {
	return &RelCreateCommand{
		CommandBase: CommandBase{
			ResourceType: RelTarget,
			Id:           fromId,
		},
		ToId:   toId,
		Params: params,
	}, nil
}

func MakeRelSetCommand(fromId string, toId string, params world.RelParams) (Command, error) {
	return &RelSetCommand{
		CommandBase: CommandBase{
			ResourceType: RelTarget,
			Id:           fromId,
		},
		ToId:      toId,
		Params:    params,
		oldParams: world.RelParams{},
	}, nil
}

func MakeRelDeleteCommand(fromId string, toId string) (Command, error) {
	return &RelDeleteCommand{
		CommandBase: CommandBase{
			ResourceType: RelTarget,
			Id:           fromId,
		},
		ToId:      toId,
		oldParams: world.RelParams{},
	}, nil
}

// --- INTERNAL FUNCTIONS ---

// inputToCommand converts a grammar.InputAttributes to a Command.
func inputToCommand(input grammar.InputAttributes) (Command, error) {
	switch input.ResourceType {
	case string(ItemTarget):
		switch input.Verb {
		case string(Create):
			params := world.ItemParams{}
			// TODO: Set the params
			return MakeItemCreateCommand(input.ResourceId, params)
		case string(Set):
			params := world.ItemParams{}
			// TODO: Set the params
			return MakeItemSetCommand(input.ResourceId, params)
		case string(Delete):
			return MakeItemDeleteCommand(input.ResourceId)
		case string(Nest):
			return MakeNestCommand(input.ResourceId, input.SecondaryId)
		case string(Free):
			return MakeItemFreeCommand(input.ResourceId)
		default:
			return nil, errors.New("unknown verb").UseCode(errors.TopolithErrorInvalid).WithData(errors.KvPair{Key: "verb", Value: input.Verb})
		}
	case string(RelTarget):
		switch input.Verb {
		case string(Create):
			params := world.RelParams{}
			// TODO: Set the params
			return MakeRelCreateCommand(input.ResourceId, input.SecondaryId, params)
		case string(Set):
			params := world.RelParams{}
			// TODO: Set the params
			return MakeRelSetCommand(input.ResourceId, input.SecondaryId, params)
		case string(Delete):
			return MakeRelDeleteCommand(input.ResourceId, input.SecondaryId)
		default:
			return nil, errors.New("unknown verb").UseCode(errors.TopolithErrorInvalid).WithData(errors.KvPair{Key: "verb", Value: input.Verb})
		}
	default:
		return nil, errors.New("unknown resource type").UseCode(errors.TopolithErrorInvalid).WithData(errors.KvPair{Key: "resourceType", Value: input.ResourceType})
	}
}

// parseCommand parses a Command from a string.
func parseCommand(s string) (Command, error) {
	firstWord, rest := scanWord(s)
	switch firstWord {
	case string(ItemTarget):
		secondWord, rest := scanWord(rest)
		switch secondWord {
		case string(Create):
			id, rest := scanWord(rest)
			params, err := itemParamsFromString(rest)
			if err != nil {
				return nil, err
			}
			return MakeItemCreateCommand(strings.TrimSpace(id), params)
		case string(Delete):
			id, rest := scanWord(rest)
			if rest != "" {
				return nil, errors.New("extra arguments after delete").UseCode(errors.TopolithErrorInvalid)
			}
			return MakeItemDeleteCommand(strings.TrimSpace(id))
		case string(Set):
		default:
			// Setting params.
			id, rest := scanWord(rest)
			remainingWords := strings.Split(strings.TrimSpace(rest), " ")
			possibleParams := kvPattern.FindAll([]byte(rest), -1)
			if len(possibleParams) != len(remainingWords) {
				return nil, errors.New("invalid params").UseCode(errors.TopolithErrorInvalid)
			}
			params, err := itemParamsFromString(rest)
			if err != nil {
				return nil, err
			}
			return MakeItemSetCommand(strings.TrimSpace(id), params)
		}
	case string(RelTarget):
		secondWord, rest := scanWord(rest)
		switch secondWord {
		case string(Create):
			id, rest := scanWord(rest)
			toId, rest := scanWord(rest)
			params, err := relParamsFromString(rest)
			if err != nil {
				return nil, err
			}
			return MakeRelCreateCommand(strings.TrimSpace(id), strings.TrimSpace(toId), params)
		case string(Delete):
			id, rest := scanWord(rest)
			toId, rest := scanWord(rest)
			if rest != "" {
				return nil, errors.New("extra arguments after delete").UseCode(errors.TopolithErrorInvalid)
			}
			return MakeRelDeleteCommand(strings.TrimSpace(id), strings.TrimSpace(toId))
		case string(Set):
		default:
			// Setting params.
			id, rest := scanWord(rest)
			toId, rest := scanWord(rest)
			remainingWords := strings.Split(strings.TrimSpace(rest), " ")
			possibleParams := kvPattern.FindAll([]byte(rest), -1)
			if len(possibleParams) != len(remainingWords) {
				return nil, errors.New("invalid params").UseCode(errors.TopolithErrorInvalid)
			}
			params, err := relParamsFromString(rest)
			if err != nil {
				return nil, err
			}
			return MakeRelSetCommand(strings.TrimSpace(id), strings.TrimSpace(toId), params)
		}
	case string(Nest):
		id, rest := scanWord(rest)
		shouldBeIn, rest := scanWord(rest)
		if shouldBeIn != "in" {
			return nil, errors.New("nest command must have 'in' between IDs").UseCode(errors.TopolithErrorInvalid)
		}
		pid, rest := scanWord(rest)
		if rest != "" {
			return nil, errors.New("extra arguments after nest").UseCode(errors.TopolithErrorInvalid)
		}
		return MakeNestCommand(strings.TrimSpace(id), strings.TrimSpace(pid))
	case string(Free):
		id, rest := scanWord(rest)
		if rest != "" {
			return nil, errors.New("extra arguments after free").UseCode(errors.TopolithErrorInvalid)
		}
		return MakeItemFreeCommand(strings.TrimSpace(id))
	default:
		return nil, errors.New("unknown command").UseCode(errors.TopolithErrorInvalid).WithData(errors.KvPair{Key: "command", Value: s})
	}
	return nil, errors.New("unknown command").UseCode(errors.TopolithErrorInvalid).WithData(errors.KvPair{Key: "command", Value: s})
}

func itemParamsToString(params world.ItemParams) string {
	components := make([]string, 0)
	if params.External != nil {
		components = append(components, fmt.Sprintf(`external=%t`, *params.External))
	}
	if params.Name != nil {
		components = append(components, fmt.Sprintf(`name="%s"`, *params.Name))
	}
	if params.Type != nil {
		components = append(components, fmt.Sprintf(`type="%s"`, *params.Type))
	}
	if params.Mechanism != nil {
		components = append(components, fmt.Sprintf(`mechanism="%s"`, *params.Mechanism))
	}
	if params.Expanded != nil {
		components = append(components, fmt.Sprintf(`expanded="%s"`, *params.Expanded))
	}
	return strings.Join(components, " ")
}

func itemParamsFromString(s string) (world.ItemParams, error) {
	elements := kvPattern.FindAllStringSubmatch(s, -1)
	params := world.ItemParams{}
	for _, element := range elements {
		key := element[1]
		value := strings.Trim(element[2], `'"`)
		switch key {
		case "external":
			external := value == "true"
			params.External = &external
		case "name":
			params.Name = &value
		case "type":
			if strings.EqualFold(value, string(ItemTarget)) {
				// Avoid weird lower-casing and set directly.
				v := string(ItemTarget)
				params.Type = &v
			}
			if strings.EqualFold(value, string(RelTarget)) {
				v := string(RelTarget)
				params.Type = &v
			}
		case "mechanism":
			params.Mechanism = &value
		case "expanded":
			params.Expanded = &value
		default:
			return params, fmt.Errorf("unknown key %s", key)
		}
	}
	return params, nil
}

func relParamsToString(params world.RelParams) string {
	components := make([]string, 0)
	if params.Verb != nil {
		components = append(components, fmt.Sprintf(`verb="%s"`, *params.Verb))
	}
	if params.Mechanism != nil {
		components = append(components, fmt.Sprintf(`mechanism="%s"`, *params.Mechanism))
	}
	if params.Async != nil {
		components = append(components, fmt.Sprintf(`async=%t`, *params.Async))
	}
	if params.Expanded != nil {
		components = append(components, fmt.Sprintf(`expanded="%s"`, *params.Expanded))
	}
	return strings.Join(components, " ")
}

func relParamsFromString(s string) (world.RelParams, error) {
	elements := kvPattern.FindAllStringSubmatch(s, -1)
	params := world.RelParams{}
	for _, element := range elements {
		key := element[1]
		value := strings.Trim(element[2], `'"`)
		switch key {
		case "verb":
			if strings.EqualFold(value, string(Create)) {
				v := string(Create)
				params.Verb = &v
			}
			if strings.EqualFold(value, string(Set)) {
				v := string(Set)
				params.Verb = &v
			}
			if strings.EqualFold(value, string(Delete)) {
				v := string(Delete)
				params.Verb = &v
			}
			if strings.EqualFold(value, string(Nest)) {
				v := string(Nest)
				params.Verb = &v
			}
			if strings.EqualFold(value, string(Free)) {
				v := string(Free)
				params.Verb = &v
			}
		case "mechanism":
			params.Mechanism = &value
		case "async":
			async := value == "true"
			params.Async = &async
		case "expanded":
			params.Expanded = &value
		default:
			return params, fmt.Errorf("unknown key %s", key)
		}
	}
	return params, nil
}

// scanWord scans the first word from a string, returning the first word and the rest.
// Scan over s until the first whitespace - take that as the first word.
// Ex: "item! 123123" -> "item!", " 123123"
func scanWord(s string) (string, string) {
	s = strings.TrimSpace(s)
	for i, r := range s {
		if unicode.IsSpace(r) {
			return s[:i], s[i:]
		}
	}
	return s, ""
}

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
