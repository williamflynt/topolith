package grammar

// InputAttributes is a struct that holds information from the parsed input to the REPL.
type InputAttributes struct {
	ResourceType string            `json:"resourceType"`
	ResourceId   string            `json:"resourceId"`
	SecondaryId  string            `json:"secondaryId"`
	Verb         string            `json:"verb"`
	Params       map[string]string `json:"params"`
	Strict       bool              `json:"strict"`
}

// Response is a struct that holds the response from our grammar.
type Response struct {
	Object string `json:"object"`
	Error  struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// Parse function to validate and pull information from the input to the REPL.
func Parse(s string) (*Parser, error) {
	p := &Parser{Buffer: s}
	if err := p.Init(); err != nil {
		return p, err
	}
	if err := p.Parse(); err != nil {
		return p, err
	}
	return p, nil
}
