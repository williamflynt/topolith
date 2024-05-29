package world

import (
	"fmt"
	"github.com/williamflynt/topolith/pkg/errors"
	"github.com/williamflynt/topolith/pkg/grammar"
	"strings"
)

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

func ItemTypeFromString(s string) ItemType {
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

func StringFromItemType(t ItemType) string {
	switch t {
	case Person:
		return "person"
	case Database:
		return "database"
	case Queue:
		return "queue"
	case Blobstore:
		return "blobstore"
	case Browser:
		return "browser"
	case Mobile:
		return "mobile"
	case Server:
		return "server"
	case Device:
		return "device"
	case Code:
		return "code"
	default:
		return ""
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

func (i Item) String() string {
	item := fmt.Sprintf(`item "%s" external=%t`, i.Id, i.External)
	paramRepr := make([]string, 0)
	if i.Type > 0 {
		paramRepr = append(paramRepr, fmt.Sprintf(`type=%s`, StringFromItemType(i.Type)))
	}
	if i.Name != "" {
		paramRepr = append(paramRepr, fmt.Sprintf(`name="%s"`, i.Name))
	}
	if i.Mechanism != "" {
		paramRepr = append(paramRepr, fmt.Sprintf(`mechanism="%s"`, i.Mechanism))
	}
	if i.Expanded != "" {
		paramRepr = append(paramRepr, fmt.Sprintf(`expanded="%s"`, i.Expanded))
	}
	if len(paramRepr) > 0 {
		item += " " + strings.Join(paramRepr, " ")
	}
	return item
}

// ItemFromString returns an Item from the string representation in accordance with grammar.Parser.
func ItemFromString(s string) (Item, error) {
	p, err := grammar.Parse(s)
	if err != nil {
		return Item{}, errors.New("error parsing Item").UseCode(errors.TopolithErrorInvalid).WithError(err).WithDescription("error parsing Item").WithData(errors.KvPair{Key: "input", Value: s})
	}
	return itemSet(Item{Id: p.InputAttributes.ResourceId}, ItemParamsFromInput(p.InputAttributes))
}

func ItemEqual(i1, i2 Item) bool {
	if i1.Id != i2.Id {
		return false
	}
	if i1.External != i2.External {
		return false
	}
	if i1.Type != i2.Type {
		return false
	}
	if i1.Name != i2.Name {
		return false
	}
	if i1.Mechanism != i2.Mechanism {
		return false
	}
	if i1.Expanded != i2.Expanded {
		return false
	}
	return true
}

// id returns the ID of the Item.
func (i Item) id() string {
	return i.Id
}

// ItemParams is a struct that represents the parameters that can be set on an Item.
type ItemParams struct {
	External  *bool   `json:"external"`
	Type      *string `json:"type"`
	Name      *string `json:"name"`
	Mechanism *string `json:"mechanism"`
	Expanded  *string `json:"expanded"`
}

func ItemParamsFromInput(input grammar.InputAttributes) ItemParams {
	params := ItemParams{}
	if v, ok := input.Params["external"]; ok {
		params.External = boolPtr(v == "true")
	}
	if v, ok := input.Params["type"]; ok {
		params.Type = strPtr(v)
	}
	if v, ok := input.Params["name"]; ok {
		params.Name = strPtr(v)
	}
	if v, ok := input.Params["mechanism"]; ok {
		params.Mechanism = strPtr(v)
	}
	if v, ok := input.Params["expanded"]; ok {
		params.Expanded = strPtr(v)
	}
	return params
}
