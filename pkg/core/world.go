package core

const currentVersion = 1

type World struct {
	Version  int             `json:"version" gen:"~"` // Version is the version of the World.
	Items    map[string]Item `json:"items" gen:"~"`   // Items is a map of ID string to related Item.
	Rels     map[string]Rel  `json:"rels" gen:"~"`    // Rels is a map of `Rel.From.Id` to related Rel.
	History  []Command       `json:"history" gen:"~"` // History is a list of commands that have been executed.
	Tree     Tree            `json:"tree" gen:"~"`    // Tree is a tree representation of the World.
	Id       string          `json:"id"`              // Id is the unique identifier of the World.
	Name     string          `json:"name"`            // Name is the display name of the World; meant for humans.
	Expanded string          `json:"expanded"`        // Expanded is the expanded description of the World.
}
