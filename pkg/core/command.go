package core

// ActionTarget is an interface that represents a target of an action.
// It is Item or a Rel notionally, but anything with `id() string` will compile.
//
// Operations on the World.Tree are performed "in the background" by operating over Item.
// For example, creating a new Item automatically places it in the World.Tree.
type ActionTarget interface{ id() string }

// ActionVerb is a string that represents the action to perform on an ActionTarget.
type ActionVerb string

/*

Action Applicability Table

| ActionTarget | Create | Set | Clear | Nest | Delete |
| ------------ | ------ | --- | ----- | ---- | ------ |
| Item         | X      | X   | X     | X    | X      |
| Rel          | X      | X   | X     |      | X      |

Even where an ActionVerb is supported, not every action can be applied to any attribute of an ActionTarget.
For example, we do not support Set on the Item.Id attribute, as it is a unique identifier.
Code generation should enforce these constraints by only creating Action instances that are valid.
*/

const (
	Create ActionVerb = "create" // Create is the action verb for creating an ActionTarget in the World.
	Set    ActionVerb = "set"    // Set is the action verb for setting an attribute on an ActionTarget.
	Clear  ActionVerb = "clear"  // Clear is the action verb for clearing an attribute on an ActionTarget.
	Nest   ActionVerb = "nest"   // Nest is the action verb for nesting an ActionTarget within another ActionTarget in the World.Tree.
	Delete ActionVerb = "delete" // Delete is the action verb for deleting an ActionTarget from the World.
)

// Action is an interface that represents an action that can be performed on a resource.
type Action[T ActionTarget] struct {
	ResourceId string     `json:"resourceId"` // ResourceId is the ID of the resource that the command is targeting.
	Verb       ActionVerb `json:"verb"`       // Verb returns the name of the action to do.
	Params     T          `json:"params"`     // Resource is the resource that the command is targeting.
}

// Command is a struct that represents a command that can be executed on a resource.
type Command struct {
	Action Action[ActionTarget] `json:"action"` // Action is the action that the command is performing.
	Undo   Action[ActionTarget] `json:"undo"`   // Undo is the action that reverses Command.Action.
}
