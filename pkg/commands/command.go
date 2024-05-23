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
	oldParams topolith.ItemSetParams
	noSet     bool
}

func (c *ItemSetCommand) Execute(w topolith.World) error {
	item, ok := w.ItemFetch(c.Id)
	if !ok {
		c.noSet = true
		return errors.New("could not find Item").UseCode(errors.TopolithErrorNotFound).WithData(errors.KvPair{Key: "id", Value: c.Id})
	}
	c.oldParams.External = boolPtr(item.External)
	c.oldParams.Name = strPtr(item.Name)
	c.oldParams.Type = strPtr(topolith.StringFromItemType(item.Type))
	c.oldParams.Mechanism = strPtr(item.Mechanism)
	c.oldParams.Expanded = strPtr(item.Expanded)
	return w.ItemSet(c.Id, c.Params).Err()
}

func (c *ItemSetCommand) Undo(w topolith.World) error {
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
	oldParams topolith.ItemSetParams
	noDelete  bool
}

func (c *ItemDeleteCommand) Execute(w topolith.World) error {
	item, ok := w.ItemFetch(c.Id)
	if !ok {
		c.noDelete = true
		return nil
	}
	c.oldParams.External = boolPtr(item.External)
	c.oldParams.Name = strPtr(item.Name)
	c.oldParams.Type = strPtr(topolith.StringFromItemType(item.Type))
	c.oldParams.Mechanism = strPtr(item.Mechanism)
	c.oldParams.Expanded = strPtr(item.Expanded)
	return w.ItemDelete(c.Id).Err()
}

func (c *ItemDeleteCommand) Undo(w topolith.World) error {
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

func (c *ItemNestCommand) Execute(w topolith.World) error {
	oldParentId, found := w.Parent(c.Id)
	if !found {
		c.noNest = true
		return errors.New("could not find Item").UseCode(errors.TopolithErrorNotFound).WithData(errors.KvPair{Key: "id", Value: c.Id})
	}
	c.oldParentId = oldParentId
	return w.Nest(c.Id, c.ParentId).Err()
}

func (c *ItemNestCommand) Undo(w topolith.World) error {
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

func (c *ItemFreeCommand) Execute(w topolith.World) error {
	oldParentId, found := w.Parent(c.Id)
	if !found {
		return errors.New("could not find Item").UseCode(errors.TopolithErrorNotFound).WithData(errors.KvPair{Key: "id", Value: c.Id})
	}
	c.oldParentId = oldParentId
	return w.Free(c.Id).Err()
}

func (c *ItemFreeCommand) Undo(w topolith.World) error {
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
	oldParams topolith.RelSetParams
	noSet     bool
}

func (c *RelSetCommand) Execute(w topolith.World) error {
	rels := w.RelFetch(c.Id, c.ToId, true)
	if len(rels) == 0 {
		c.noSet = true
		return errors.New("could not find Rel").UseCode(errors.TopolithErrorNotFound).WithData(errors.KvPair{Key: "id", Value: c.Id})
	}
	rel := rels[0]
	c.oldParams.Verb = strPtr(rel.Verb)
	c.oldParams.Mechanism = strPtr(rel.Mechanism)
	c.oldParams.Async = boolPtr(rel.Async)
	c.oldParams.Expanded = strPtr(rel.Expanded)
	return w.RelSet(c.Id, c.ToId, c.Params).Err()
}

func (c *RelSetCommand) Undo(w topolith.World) error {
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
	oldParams topolith.RelSetParams
	noDelete  bool
}

func (c *RelDeleteCommand) Execute(w topolith.World) error {
	rels := w.RelFetch(c.Id, c.Id, true)
	if len(rels) == 0 {
		c.noDelete = true
		return nil
	}
	rel := rels[0]
	c.oldParams.Verb = strPtr(rel.Verb)
	c.oldParams.Mechanism = strPtr(rel.Mechanism)
	c.oldParams.Async = boolPtr(rel.Async)
	c.oldParams.Expanded = strPtr(rel.Expanded)
	return w.RelDelete(c.Id, c.ToId).Err()
}

func (c *RelDeleteCommand) Undo(w topolith.World) error {
	if c.noDelete {
		return nil
	}
	return w.RelCreate(c.Id, c.Id, c.oldParams).Err()
}

func (c *RelDeleteCommand) String() string {
	return fmt.Sprintf(`%s %s "%s" "%s"`, c.ResourceType, Delete, c.Id, c.ToId)
}

// --- INTERNAL FUNCTIONS ---

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
