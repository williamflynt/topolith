package grammar

import "testing"

var testCommands = []struct {
	Input       string
	ExpectError bool
}{
	// Valid Commands
	{`item create abc123`, false},
	{`item create abc123 external=true`, false},
	{`item clear abc123 external name`, false},
	{`item clear abc123 external`, false},
	{`item set abc123 name="this is the name"`, false},
	{`item delete abc123`, false},
	{`item abc123`, false},
	{`item abc123 external=true`, false},
	{`item abc123 name=myname`, false},
	{`item abc123 name="my name"`, false},
	{`nest abc123 in "this is my ID"`, false},
	{`rel create abc123 def456`, false},
	{`rel set abc123 def456 verb="connect"`, false},
	{`rel delete abc123 def456`, false},
	{`rel fetch abc123 def456`, false},
	{`rel abc123 def456`, false},
	{`item fetch abc123`, false},
	{`item in abc123`, false},
	{`item? abc123`, false},
	{`in? abc123`, false},
	{`rel abc123`, false},
	{`rel? abc123 def456`, false},
	{`item create abc123 type=person`, false},
	{`rel create abc123 def456 verb="connect"`, false},
	{`rel set abc123 def456 async=true`, false},
	{`rel clear abc123 honkhonn async`, false},

	// Invalid Commands
	{`item create`, true},
	{`item set`, true},
	{`item delete`, true},
	{`item clear abc123 external name="ok"`, true},
	{`item clear abc123`, true},
	{`item create abc123 unknown=true`, true},
	{`item set abc123 unknown="value"`, true},
	{`rel create abc123`, true},
	{`rel set abc123`, true},
	{`rel delete abc123`, true},
	{`item create abc123 external=invalid`, true},
	{`rel create abc123 def456 verb=`, true},
	{`item create abc123 type=unknown`, true},
	{`item set abc123 name=`, true},
	{`nest abc123 in`, true},
	{`rel abc123 def456 verb`, true},
	{`rel create abc123 def456 mechanism=`, true},
	{`item create abc123 type=`, true},
	{`rel set abc123 def456 async=notabool`, true},
	{`rel clear abc123 honkhonn`, true},
	{`rel clear abc123 honkhonn async=true`, true},
}

func parse(s string) (*parser, error) {
	p := &parser{
		resourceType: "",
		resourceId:   "",
		secondaryId:  "",
		verb:         "",
		params:       make(map[string]string),
		strict:       false,
		Buffer:       s,
		buffer:       nil,
		rules:        [67]func() bool{},
		parse:        nil,
		reset:        nil,
		Pretty:       false,
	}
	if err := p.Init(); err != nil {
		return p, err
	}
	if err := p.Parse(); err != nil {
		return p, err
	}
	return p, nil
}

func TestCommands(t *testing.T) {
	for _, c := range testCommands {
		t.Run(c.Input, func(t *testing.T) {
			p, err := parse(c.Input)
			if c.ExpectError {
				if err == nil {
					t.Errorf("expected error for command: '%s', but got none", c.Input)
					p.PrintSyntaxTree()
					return
				}
			} else {
				if err != nil {
					t.Errorf("did not expect error for command: '%s', but got: %s", c.Input, err)
				}
			}
			p.PrintSyntaxTree()
		})
	}
}
