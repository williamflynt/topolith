package commands

import (
	"fmt"
	"github.com/williamflynt/topolith/pkg/errors"
	"github.com/williamflynt/topolith/pkg/topolith"
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

// Command is the interface that all commands must implement.
type Command interface {
	Execute(w topolith.World) error
	Undo(w topolith.World) error
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
	Params   topolith.ItemSetParams
	noCreate bool
}

func (c *ItemCreateCommand) Execute(w topolith.World) error {
	if _, ok := w.ItemFetch(c.Id); ok {
		c.noCreate = true
		return nil
	}
	return w.ItemCreate(c.Id, c.Params).Err()
}

func (c *ItemCreateCommand) Undo(w topolith.World) error {
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
	Params    topolith.ItemSetParams
	OldParams topolith.ItemSetParams
	noSet     bool
}

func (c *ItemSetCommand) Execute(w topolith.World) error {
	item, ok := w.ItemFetch(c.Id)
	if !ok {
		c.noSet = true
		return errors.New("could not find Item").UseCode(errors.TopolithErrorNotFound).WithData(errors.KvPair{"id", c.Id})
	}
	c.OldParams.External = boolPtr(item.External)
	c.OldParams.Name = strPtr(item.Name)
	c.OldParams.Type = strPtr(topolith.StringFromItemType(item.Type))
	c.OldParams.Mechanism = strPtr(item.Mechanism)
	c.OldParams.Expanded = strPtr(item.Expanded)
	return w.ItemSet(c.Id, c.Params).Err()
}

func (c *ItemSetCommand) Undo(w topolith.World) error {
	if c.noSet {
		return nil
	}
	return w.ItemSet(c.Id, c.OldParams).Err()
}

func (c *ItemSetCommand) String() string {
	return fmt.Sprintf(`%s "%s" %s`, c.ResourceType, c.Id, itemParamsToString(c.Params))
}

// ItemDeleteCommand represents a delete command for Item.
type ItemDeleteCommand struct {
	CommandBase
	OldParams topolith.ItemSetParams
	noDelete  bool
}

func (c *ItemDeleteCommand) Execute(w topolith.World) error {
	item, ok := w.ItemFetch(c.Id)
	if !ok {
		c.noDelete = true
		return nil
	}
	c.OldParams.External = boolPtr(item.External)
	c.OldParams.Name = strPtr(item.Name)
	c.OldParams.Type = strPtr(topolith.StringFromItemType(item.Type))
	c.OldParams.Mechanism = strPtr(item.Mechanism)
	c.OldParams.Expanded = strPtr(item.Expanded)
	return w.ItemDelete(c.Id).Err()
}

func (c *ItemDeleteCommand) Undo(w topolith.World) error {
	if c.noDelete {
		return nil
	}
	return w.ItemCreate(c.Id, c.OldParams).Err()
}

func (c *ItemDeleteCommand) String() string {
	return fmt.Sprintf(`%s %s "%s"`, c.ResourceType, Delete, c.Id)
}

// ItemNestCommand represents a nest command for Item.
type ItemNestCommand struct {
	CommandBase
	ParentId    string
	OldParentId string
	noNest      bool
}

func (c *ItemNestCommand) Execute(w topolith.World) error {
	oldParentId, found := w.Parent(c.Id)
	if !found {
		c.noNest = true
		return errors.New("could not find Item").UseCode(errors.TopolithErrorNotFound).WithData(errors.KvPair{Key: "id", Value: c.Id})
	}
	c.OldParentId = oldParentId
	return w.Nest(c.Id, c.ParentId).Err()
}

func (c *ItemNestCommand) Undo(w topolith.World) error {
	if c.noNest {
		return nil
	}
	if c.OldParentId == "" {
		return w.Free(c.Id).Err()
	}
	return w.Nest(c.Id, c.OldParentId).Err()
}

func (c *ItemNestCommand) String() string {
	return fmt.Sprintf(`%s "%s" in "%s"`, Nest, c.Id, c.ParentId)
}

// ItemFreeCommand represents a free command for Item.
type ItemFreeCommand struct {
	CommandBase
	OldParentId string
}

func (c *ItemFreeCommand) Execute(w topolith.World) error {
	oldParentId, found := w.Parent(c.Id)
	if !found {
		return errors.New("could not find Item").UseCode(errors.TopolithErrorNotFound).WithData(errors.KvPair{Key: "id", Value: c.Id})
	}
	c.OldParentId = oldParentId
	return w.Free(c.Id).Err()
}

func (c *ItemFreeCommand) Undo(w topolith.World) error {
	if c.OldParentId == "" {
		return nil
	}
	return w.Nest(c.Id, c.OldParentId).Err()
}

func (c *ItemFreeCommand) String() string {
	return fmt.Sprintf(`%s "%s"`, Free, c.Id)
}

// RelCreateCommand represents a create command for Rel.
type RelCreateCommand struct {
	CommandBase
	ToId     string
	Params   topolith.RelSetParams
	noCreate bool
}

func (c *RelCreateCommand) Execute(w topolith.World) error {
	if rels := w.RelFetch(c.Id, c.ToId, true); len(rels) > 0 {
		c.noCreate = true
		return nil
	}
	return w.RelCreate(c.Id, c.ToId, c.Params).Err()
}

func (c *RelCreateCommand) Undo(w topolith.World) error {
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
	Params    topolith.RelSetParams
	OldParams topolith.RelSetParams
	noSet     bool
}

func (c *RelSetCommand) Execute(w topolith.World) error {
	rels := w.RelFetch(c.Id, c.ToId, true)
	if len(rels) == 0 {
		c.noSet = true
		return errors.New("could not find Rel").UseCode(errors.TopolithErrorNotFound).WithData(errors.KvPair{Key: "id", Value: c.Id})
	}
	rel := rels[0]
	c.OldParams.Verb = strPtr(rel.Verb)
	c.OldParams.Mechanism = strPtr(rel.Mechanism)
	c.OldParams.Async = boolPtr(rel.Async)
	c.OldParams.Expanded = strPtr(rel.Expanded)
	return w.RelSet(c.Id, c.ToId, c.Params).Err()
}

func (c *RelSetCommand) Undo(w topolith.World) error {
	if c.noSet {
		return nil
	}
	return w.RelSet(c.Id, c.ToId, c.OldParams).Err()
}

func (c *RelSetCommand) String() string {
	return fmt.Sprintf(`%s "%s" "%s" %s`, c.ResourceType, c.Id, c.ToId, relParamsToString(c.Params))
}

// RelDeleteCommand represents a delete command for Rel.
type RelDeleteCommand struct {
	CommandBase
	OldParams topolith.RelSetParams
	noDelete  bool
}

func (c *RelDeleteCommand) Execute(w topolith.World) error {
	rels := w.RelFetch(c.Id, c.Id, true)
	if len(rels) == 0 {
		c.noDelete = true
		return nil
	}
	rel := rels[0]
	c.OldParams.Verb = strPtr(rel.Verb)
	c.OldParams.Mechanism = strPtr(rel.Mechanism)
	c.OldParams.Async = boolPtr(rel.Async)
	c.OldParams.Expanded = strPtr(rel.Expanded)
	return w.RelDelete(c.Id, c.Id).Err()
}

func (c *RelDeleteCommand) Undo(w topolith.World) error {
	if c.noDelete {
		return nil
	}
	return w.RelCreate(c.Id, c.Id, c.OldParams).Err()
}

func (c *RelDeleteCommand) String() string {
	return fmt.Sprintf(`%s %s "%s"`, c.ResourceType, Delete, c.Id)
}

// --- EXPORTED FUNCTIONS ---

// ParseCommand parses a Command from a string.
func ParseCommand(s string) (Command, error) {
	// TODO: Complete implementation. It turns out I'm making a grammar/protocol...
	firstWord := scanWord(s)
	switch firstWord {
	case string(ItemTarget):
		return nil, nil
	case string(RelTarget):
		return nil, nil
	case string(Nest):
		return nil, nil
	case string(Free):
		return nil, nil
	default:
		return nil, errors.New("unknown command").UseCode(errors.TopolithErrorInvalid).WithData(errors.KvPair{Key: "command", Value: s})
	}
}

// --- INTERNAL FUNCTIONS ---

func itemParamsToString(params topolith.ItemSetParams) string {
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

func itemParamsFromString(s string) (topolith.ItemSetParams, error) {
	elements := kvPattern.FindAllStringSubmatch(s, -1)
	params := topolith.ItemSetParams{}
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

func relParamsToString(params topolith.RelSetParams) string {
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

func relParamsFromString(s string) (topolith.RelSetParams, error) {
	elements := kvPattern.FindAllStringSubmatch(s, -1)
	params := topolith.RelSetParams{}
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

// scanWord scans the first word from a string.
// Scan over s until the first non-alpha character - take that as the first word.
// Ex: "item! 123123" -> "item"
func scanWord(s string) string {
	for i, r := range s {
		if !unicode.IsLetter(r) {
			return s[:i]
		}
	}
	return s
}

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
