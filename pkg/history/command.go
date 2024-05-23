package history

import (
	"fmt"
	"github.com/williamflynt/topolith/pkg/errors"
	"github.com/williamflynt/topolith/pkg/topolith"
)

// ActionVerb represents the action to be performed.
type ActionVerb string

const (
	Create ActionVerb = "create"
	Set    ActionVerb = "set"
	Clear  ActionVerb = "clear"
	Nest   ActionVerb = "nest"
	Delete ActionVerb = "delete"
)

// Command is the interface that all commands must implement.
type Command interface {
	Execute(w topolith.World) error
	Undo(w topolith.World) error
	String() string
	FromString(s string) (Command, error)
}

// CommandTarget represents the type of resource.
type CommandTarget int

const (
	ItemTarget CommandTarget = iota
	RelTarget
)

// CommandBase is a base struct for common command fields.
type CommandBase struct {
	ResourceType CommandTarget
	ID           string
}

// CreateCommand represents a create command.
type CreateCommand struct {
	CommandBase
}

func (c *CreateCommand) Execute(w topolith.World) error {
	// Perform create operation
	fmt.Printf("Creating %v with ID %s\n", c.ResourceType, c.ID)
	return nil
}

func (c *CreateCommand) Undo(w topolith.World) error {
	// Undo create operation
	fmt.Printf("Undo creating %v with ID %s\n", c.ResourceType, c.ID)
	// TODO Implement
	return nil
}

func (c *CreateCommand) String() string {
	return fmt.Sprintf("Create %v with ID %s", c.ResourceType, c.ID)
}

func (c *CreateCommand) FromString(s string) (Command, error) {
	// TODO Implement
	panic("not implemented")
}

// SetItemCommand represents a set command for Item.
type SetItemCommand struct {
	CommandBase
	Params topolith.ItemSetParams
}

func (c *SetItemCommand) Execute(w topolith.World) error {
	// Perform set operation
	fmt.Printf("Setting Item with ID %s and params %+v\n", c.ID, c.Params)
	return nil
}

func (c *SetItemCommand) Undo(w topolith.World) error {
	// Undo set operation
	fmt.Printf("Undo setting Item with ID %s\n", c.ID)
	// TODO Implement
	return nil
}

func (c *SetItemCommand) String() string {
	return fmt.Sprintf("Set Item with ID %s and params %+v", c.ID, c.Params)
}

func (c *SetItemCommand) FromString(s string) (Command, error) {
	// TODO Implement
	panic("not implemented")
}

// SetRelCommand represents a set command for Rel.
type SetRelCommand struct {
	CommandBase
	Params topolith.RelSetParams
}

func (c *SetRelCommand) Execute(w topolith.World) error {
	// Perform set operation
	fmt.Printf("Setting Rel from %s to %s with params %+v\n", c.ID, c.ID, c.Params)
	return nil
}

func (c *SetRelCommand) Undo(w topolith.World) error {
	// Undo set operation
	fmt.Printf("Undo setting Rel from %s to %s\n", c.ID, c.ID)
	// TODO Implement
	return nil
}

func (c *SetRelCommand) String() string {
	return fmt.Sprintf("Set Rel from %s to %s with params %+v", c.ID, c.ID, c.Params)
}

func (c *SetRelCommand) FromString(s string) (Command, error) {
	// TODO Implement
	panic("not implemented")
}

// TODO Additional commands like ClearCommand, NestCommand, DeleteCommand can be similarly defined.

// CommandFactory creates commands based on the verb and resource type.
func CommandFactory(resourceType CommandTarget, id string, verb ActionVerb, params interface{}) (Command, error) {
	switch verb {
	case Create:
		return &CreateCommand{CommandBase: CommandBase{ResourceType: resourceType, ID: id}}, nil
	case Set:
		switch resourceType {
		case ItemTarget:
			p, ok := params.(topolith.ItemSetParams)
			if !ok {
				return nil, errors.New("invalid parameters for set item command")
			}
			return &SetItemCommand{CommandBase: CommandBase{ResourceType: ItemTarget, ID: id}, Params: p}, nil
		case RelTarget:
			p, ok := params.(topolith.RelSetParams)
			if !ok {
				return nil, errors.New("invalid parameters for set rel command")
			}
			return &SetRelCommand{CommandBase: CommandBase{ResourceType: RelTarget, ID: id}, Params: p}, nil
		default:
			return nil, errors.New("unknown resource type")
		}
	// Implement other cases for Clear, Nest, and Delete similarly.
	case Clear, Nest, Delete:
		// Handle other verbs similarly
	default:
		return nil, errors.New("unknown verb")
	}
	return nil, errors.New("command not created")
}

func main() {
	// Example usage:
	createCmd, err := CommandFactory(ItemTarget, "item1", Create, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(createCmd.String())

	w := topolith.CreateWorld("myWorld")
	createCmd.Execute(w)
	createCmd.Undo(w)

	itemParams := topolith.ItemSetParams{
		Name:      stringPtr("myname"),
		Expanded:  stringPtr("this is an item"),
		External:  boolPtr(true),
		Type:      stringPtr("TypeA"),
		Mechanism: stringPtr("MechanismA"),
	}
	setItemCmd, err := CommandFactory(ItemTarget, "item1", Set, itemParams)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(setItemCmd.String())
	setItemCmd.Execute(w)
	setItemCmd.Undo(w)
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

// Dummy functions to simulate state retrieval and restoration.
func getItemState(id string) interface{} {
	// Return dummy previous state
	return topolith.ItemSetParams{Name: stringPtr("oldName")}
}

func restoreItemState(id string, state interface{}) {
	// Restore the state
	fmt.Printf("Restored state for item %s: %+v\n", id, state)
}

func getRelState(id string) interface{} {
	// Return dummy previous state
	return topolith.RelSetParams{Verb: stringPtr("oldVerb")}
}

func restoreRelState(id string, state interface{}) {
	// Restore the state
	fmt.Printf("Restored state for rel %s: %+v\n", id, state)
}
