package history

import (
	"fmt"
	"github.com/williamflynt/topolith/pkg/errors"
	"github.com/williamflynt/topolith/pkg/topolith"
)

// CommandVerb represents the action to be performed.
type CommandVerb string

const (
	Create CommandVerb = "create"
	Set    CommandVerb = "set"
	Delete CommandVerb = "delete"
	Nest   CommandVerb = "nest"
)

// Command is the interface that all commands must implement.
type Command interface {
	Execute(w topolith.World) error
	Undo(w topolith.World) error
	String() string
	FromString(s string) (Command, error)
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
	Params topolith.ItemSetParams
}

func (c *ItemCreateCommand) Execute(w topolith.World) error {
	// Perform create operation
	fmt.Printf("Creating %v with Id %s\n", c.ResourceType, c.Id)
	return nil
}

func (c *ItemCreateCommand) Undo(w topolith.World) error {
	// Undo create operation
	fmt.Printf("Undo creating %v with Id %s\n", c.ResourceType, c.Id)
	// TODO Implement
	return nil
}

func (c *ItemCreateCommand) String() string {
	return fmt.Sprintf("Create %v with Id %s", c.ResourceType, c.Id)
}

func (c *ItemCreateCommand) FromString(s string) (Command, error) {
	// TODO Implement
	panic("not implemented")
}

// ItemSetCommand represents a set command for Item.
type ItemSetCommand struct {
	CommandBase
	Params topolith.ItemSetParams
}

func (c *ItemSetCommand) Execute(w topolith.World) error {
	// Perform set operation
	fmt.Printf("Setting Item with Id %s and params %+v\n", c.Id, c.Params)
	return nil
}

func (c *ItemSetCommand) Undo(w topolith.World) error {
	// Undo set operation
	fmt.Printf("Undo setting Item with Id %s\n", c.Id)
	// TODO Implement
	return nil
}

func (c *ItemSetCommand) String() string {
	return fmt.Sprintf("Set Item with Id %s and params %+v", c.Id, c.Params)
}

func (c *ItemSetCommand) FromString(s string) (Command, error) {
	// TODO Implement
	panic("not implemented")
}

// ItemDeleteCommand represents a delete command for Item.
type ItemDeleteCommand struct {
	CommandBase
}

func (c *ItemDeleteCommand) Execute(w topolith.World) error {
	// Perform delete operation
	fmt.Printf("Deleting Item with Id %s\n", c.Id)
	return nil
}

func (c *ItemDeleteCommand) Undo(w topolith.World) error {
	// Undo delete operation
	fmt.Printf("Undo deleting Item with Id %s\n", c.Id)
	// TODO Implement
	return nil
}

func (c *ItemDeleteCommand) String() string {
	return fmt.Sprintf("Delete Item with Id %s", c.Id)
}

func (c *ItemDeleteCommand) FromString(s string) (Command, error) {
	// TODO Implement
	panic("not implemented")
}

// ItemNestCommand represents a nest command for Item.
type ItemNestCommand struct {
	CommandBase
	ParentId string
}

func (c *ItemNestCommand) Execute(w topolith.World) error {
	// Perform nest operation
	fmt.Printf("Nesting Item with Id %s under Item with Id %s\n", c.Id, c.ParentId)
	return nil
}

func (c *ItemNestCommand) Undo(w topolith.World) error {
	// Undo nest operation
	fmt.Printf("Undo nesting Item with Id %s under Item with Id %s\n", c.Id, c.ParentId)
	// TODO Implement
	return nil
}

func (c *ItemNestCommand) String() string {
	return fmt.Sprintf("Nest Item with Id %s under Item with Id %s", c.Id, c.ParentId)
}

func (c *ItemNestCommand) FromString(s string) (Command, error) {
	// TODO Implement
	panic("not implemented")
}

// RelCreateCommand represents a create command for Rel.
type RelCreateCommand struct {
	CommandBase
	Params topolith.RelSetParams
}

func (c *RelCreateCommand) Execute(w topolith.World) error {
	// Perform create operation
	fmt.Printf("Creating Rel from %s to %s with params %+v\n", c.Id, c.Id, c.Params)
	return nil
}

func (c *RelCreateCommand) Undo(w topolith.World) error {
	// Undo create operation
	fmt.Printf("Undo creating Rel from %s to %s\n", c.Id, c.Id)
	// TODO Implement
	return nil
}

func (c *RelCreateCommand) String() string {
	return fmt.Sprintf("Create Rel from %s to %s with params %+v", c.Id, c.Id, c.Params)
}

func (c *RelCreateCommand) FromString(s string) (Command, error) {
	// TODO Implement
	panic("not implemented")
}

// RelSetCommand represents a set command for Rel.
type RelSetCommand struct {
	CommandBase
	Params topolith.RelSetParams
}

func (c *RelSetCommand) Execute(w topolith.World) error {
	// Perform set operation
	fmt.Printf("Setting Rel from %s to %s with params %+v\n", c.Id, c.Id, c.Params)
	return nil
}

func (c *RelSetCommand) Undo(w topolith.World) error {
	// Undo set operation
	fmt.Printf("Undo setting Rel from %s to %s\n", c.Id, c.Id)
	// TODO Implement
	return nil
}

func (c *RelSetCommand) String() string {
	return fmt.Sprintf("Set Rel from %s to %s with params %+v", c.Id, c.Id, c.Params)
}

func (c *RelSetCommand) FromString(s string) (Command, error) {
	// TODO Implement
	panic("not implemented")
}

// RelDeleteCommand represents a delete command for Rel.
type RelDeleteCommand struct {
	CommandBase
}

func (c *RelDeleteCommand) Execute(w topolith.World) error {
	// Perform delete operation
	fmt.Printf("Deleting Rel from %s to %s\n", c.Id, c.Id)
	return nil
}

func (c *RelDeleteCommand) Undo(w topolith.World) error {
	// Undo delete operation
	fmt.Printf("Undo deleting Rel from %s to %s\n", c.Id, c.Id)
	// TODO Implement
	return nil
}

func (c *RelDeleteCommand) String() string {
	return fmt.Sprintf("Delete Rel from %s to %s", c.Id, c.Id)
}

func (c *RelDeleteCommand) FromString(s string) (Command, error) {
	// TODO Implement
	panic("not implemented")
}

// CommandFactory creates commands based on the verb and resource type.
func CommandFactory(resourceType CommandTarget, id string, verb CommandVerb, params interface{}) (Command, error) {
	switch resourceType {
	case ItemTarget:
		switch verb {
		case Create:
			p, ok := params.(topolith.ItemSetParams)
			if !ok {
				return nil, errors.New("params must be of type ItemSetParams").UseCode(errors.TopolithErrorInvalid)
			}
			return &ItemCreateCommand{
				CommandBase: CommandBase{
					ResourceType: resourceType,
					Id:           id,
				},
				Params: p,
			}, nil
		case Set:
			p, ok := params.(topolith.ItemSetParams)
			if !ok {
				return nil, errors.New("params must be of type ItemSetParams").UseCode(errors.TopolithErrorInvalid)
			}
			return &ItemSetCommand{
				CommandBase: CommandBase{
					ResourceType: resourceType,
					Id:           id,
				},
				Params: p,
			}, nil
		case Delete:
			return &ItemDeleteCommand{
				CommandBase: CommandBase{
					ResourceType: resourceType,
					Id:           id,
				},
			}, nil
		case Nest:
			pid, ok := params.(string)
			if !ok {
				return nil, errors.New("params must be of type string").UseCode(errors.TopolithErrorInvalid)
			}
			return &ItemNestCommand{
				CommandBase: CommandBase{
					ResourceType: resourceType,
					Id:           id,
				},
				ParentId: pid,
			}, nil
		default:
			return nil, errors.New("unknown verb for Item").UseCode(errors.TopolithErrorInvalid)
		}
	case RelTarget:
		switch verb {
		case Create:
			p, ok := params.(topolith.RelSetParams)
			if !ok {
				return nil, errors.New("params must be of type RelSetParams").UseCode(errors.TopolithErrorInvalid)
			}
			return &RelCreateCommand{
				CommandBase: CommandBase{
					ResourceType: resourceType,
					Id:           id,
				},
				Params: p,
			}, nil
		case Set:
			p, ok := params.(topolith.RelSetParams)
			if !ok {
				return nil, errors.New("params must be of type RelSetParams").UseCode(errors.TopolithErrorInvalid)
			}
			return &RelSetCommand{
				CommandBase: CommandBase{
					ResourceType: resourceType,
					Id:           id,
				},
				Params: p,
			}, nil
		case Delete:
			return &RelDeleteCommand{
				CommandBase: CommandBase{
					ResourceType: resourceType,
					Id:           id,
				},
			}, nil
		default:
			return nil, errors.New("unknown verb for Rel").UseCode(errors.TopolithErrorInvalid)
		}
	default:
		return nil, errors.New("unknown resource type").UseCode(errors.TopolithErrorInvalid)
	}
}
