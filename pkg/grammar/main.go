package grammar

import "strings"

// TODO(wf 27 May 2024): We shouldn't be setting raw strings
//  (ex: `Flags` -> []Flag)
//  Make ResourceType, Flag, etc... constants
//  Then work those into the `.peg` file (the Go code parts).

// Node is a struct that represents a node in a tree.
// We use this to extract the structure of the world.Tree when we parse a World or standalone tree.
type Node struct {
	Id       string
	Children []Node
}

// InputAttributes is a struct that holds information from the parsed input to the REPL.
type InputAttributes struct {
	ResourceType string            `json:"resourceType"`
	ResourceId   string            `json:"resourceId"`
	ResourceIds  []string          `json:"resourceIds"`
	SecondaryIds []string          `json:"secondaryIds"`
	Verb         string            `json:"verb"`
	Params       map[string]string `json:"params"`
	Flags        []string          `json:"flags"`

	Raw string `json:"raw"`
}

// Response is a struct that holds the response from our grammar.
type Response struct {
	Object ResponseObject `json:"object"`
	Status ResponseStatus `json:"status"`
}

type ResponseObject struct {
	Type string `json:"type"`
	Repr string `json:"repr"`
}

type ResponseStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Parse function to validate and pull information from the input to the REPL.
func Parse(s string) (*Parser, error) {
	p := &Parser{
		Buffer: s,
		InputAttributes: InputAttributes{
			ResourceIds:  make([]string, 0),
			SecondaryIds: make([]string, 0),
			Verb:         "",
			Params:       make(map[string]string),
			Flags:        make([]string, 0),
		},
	}
	if err := p.Init(); err != nil {
		return p, err
	}
	if err := p.Parse(); err != nil {
		return p, err
	}
	p.Execute()
	return p, nil
}

// --- INTERNAL ---

func cleanString(s string) string {
	return strings.Trim(strings.TrimSpace(s), "\"")
}
