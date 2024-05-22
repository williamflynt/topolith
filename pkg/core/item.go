package core

//go:generate stringer -type=ItemType -output=generated_itemtype_string.go

// ItemType is an iota that represents the type of Item.
type ItemType int

const (
	_ ItemType = iota
	Person
	Database
	Queue
	Blobstore
	Browser
	Mobile
	Server
	Device
	Code
)

// Item is a struct that represents a single entity in the world. If using the C4 diagrams methodology, an Item can be any of Person, Software System, Container, Component, or Code element.
type Item struct {
	Id        string   `json:"id" gen:"~"` // Id is the unique identifier of the Item. Often the same as Name.
	External  bool     `json:"external"`   // External is a boolean that represents whether the Item is external - a C4 diagrams concept.
	Type      ItemType `json:"type"`       // Type is the type of the Item.
	Name      string   `json:"name"`       // Name is the display name of the Item; meant for humans.
	Mechanism string   `json:"mechanism"`  // Mechanism is the method of implementation of the Item. This may not always be relevant for a strict C4 diagram.
	Expanded  string   `json:"expanded"`   // Expanded is the expanded description of the Item. This may not always be relevant for a strict C4 diagram.
}

// id returns the ID of the Item.
// We implement ActionTarget so that we can use Item as a target in Command.Action.
func (i Item) id() string {
	return i.Id
}