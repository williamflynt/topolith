package core

// Action is an interface that represents an action that can be performed on a resource.
type Action interface {
	Name() string // Name returns the name of the action.
}

type CommandTarget int

const (
	_ CommandTarget = iota
	WorldCommandTarget
	ItemCommandTarget
	RelCommandTarget
)

type Command struct {
	ResourceType CommandTarget `json:"resourceType"` // ResourceType is the type of resource that the command is targeting.
	ResourceId   string        `json:"resourceId"`   // ResourceId is the ID of the resource that the command is targeting.
	Action       Action        `json:"action"`       // Action is the action that the command is performing.
	Undo         Action        `json:"undo"`         // Undo is the action that undoes the command.
}
