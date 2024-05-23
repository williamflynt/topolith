package topolith

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

func itemTypeFromString(s string) ItemType {
	switch s {
	case "person":
		return Person
	case "database":
		return Database
	case "queue":
		return Queue
	case "blobstore":
		return Blobstore
	case "browser":
		return Browser
	case "mobile":
		return Mobile
	case "server":
		return Server
	case "device":
		return Device
	case "code":
		return Code
	default:
		return 0
	}
}

// Item is a struct that represents a single entity in the world. If using the C4 diagrams methodology, an Item can be any of Person, Software System, Container, Component, or Code element.
type Item struct {
	Id        string   `json:"id"`        // Id is the unique identifier of the Item. Often the same as Name.
	External  bool     `json:"external"`  // External is a boolean that represents whether the Item is external - a C4 diagrams concept.
	Type      ItemType `json:"type"`      // Type is the type of the Item.
	Name      string   `json:"name"`      // Name is the display name of the Item; meant for humans.
	Mechanism string   `json:"mechanism"` // Mechanism is the method of implementation of the Item. This may not always be relevant for a strict C4 diagram.
	Expanded  string   `json:"expanded"`  // Expanded is the expanded description of the Item. This may not always be relevant for a strict C4 diagram.
}

// id returns the ID of the Item.
func (i Item) id() string {
	return i.Id
}

// ItemSetParams is a struct that represents the parameters that can be set on an Item.
type ItemSetParams struct {
	External  *bool   `json:"external"`
	Type      *string `json:"type"`
	Name      *string `json:"name"`
	Mechanism *string `json:"mechanism"`
	Expanded  *string `json:"expanded"`
}
