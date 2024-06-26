package grammar

import (
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/williamflynt/topolith/pkg/errors"
	"reflect"
	"testing"
)

// TODO(wf 27 May 2024): More robust testing for commands.
// TODO(wf 27 May 2024): Test responses, errors, World representation and parsing.

var simpleTree = "tree{nil::[tree{item \"2\"::[tree{item \"1\"::[]}]} tree{item \"3\"::[]}]}"
var simpleWorld = "version=1\nid=1\nname=worldname\nexpanded=\"this is expanded data\"\ntree{nil::[tree{item \"2\"::[tree{item \"1\"::[]}]} tree{item \"3\"::[]}]}\nrel \"3\" \"2\"\nrel \"1\" \"2\""

var testCommands = []struct {
	In  string
	Err bool
	Out InputAttributes
}{
	{In: "item create abc123", Err: false, Out: InputAttributes{ResourceType: "item", ResourceId: "abc123", ResourceIds: []string{}, SecondaryIds: []string{}, Verb: "create", Params: map[string]string{}, Flags: []string{}}},
	{In: `item create "my abc123"`, Err: false, Out: InputAttributes{ResourceType: "item", ResourceId: "my abc123", ResourceIds: []string{}, SecondaryIds: []string{}, Verb: "create", Params: map[string]string{}, Flags: []string{}}},
	{In: `item create "my abc123" name=John`, Err: false, Out: InputAttributes{ResourceType: "item", ResourceId: "my abc123", ResourceIds: []string{}, SecondaryIds: []string{}, Verb: "create", Params: map[string]string{"name": "John"}, Flags: []string{}}},
	{In: "item set abc123 name=John", Err: false, Out: InputAttributes{ResourceType: "item", ResourceId: "abc123", ResourceIds: []string{}, SecondaryIds: []string{}, Verb: "set", Params: map[string]string{"name": "John"}, Flags: []string{}}},
	{In: "item clear abc123 name", Err: false, Out: InputAttributes{ResourceType: "item", ResourceId: "abc123", ResourceIds: []string{}, SecondaryIds: []string{}, Verb: "clear", Params: map[string]string{"name": ""}, Flags: []string{}}},
	{In: "item delete abc123", Err: false, Out: InputAttributes{ResourceType: "item", ResourceId: "abc123", ResourceIds: []string{}, SecondaryIds: []string{}, Verb: "delete", Params: map[string]string{}, Flags: []string{}}},
	{In: "rel create abc123 def456", Err: false, Out: InputAttributes{ResourceType: "rel", ResourceId: "abc123", ResourceIds: []string{}, SecondaryIds: []string{"def456"}, Verb: "create", Params: map[string]string{}, Flags: []string{}}},
	{In: "rel set abc123 def456 verb=likes", Err: false, Out: InputAttributes{ResourceType: "rel", ResourceId: "abc123", ResourceIds: []string{}, SecondaryIds: []string{"def456"}, Verb: "set", Params: map[string]string{"verb": "likes"}, Flags: []string{}}},
	{In: "rel clear abc123 def456 verb", Err: false, Out: InputAttributes{ResourceType: "rel", ResourceId: "abc123", ResourceIds: []string{}, SecondaryIds: []string{"def456"}, Verb: "clear", Params: map[string]string{"verb": ""}, Flags: []string{}}},
	{In: "rel delete abc123 def456", Err: false, Out: InputAttributes{ResourceType: "rel", ResourceId: "abc123", ResourceIds: []string{}, SecondaryIds: []string{"def456"}, Verb: "delete", Params: map[string]string{}, Flags: []string{}}},
	{In: "free abc123 def456", Err: false, Out: InputAttributes{ResourceType: "item", ResourceId: "", ResourceIds: []string{"abc123", "def456"}, SecondaryIds: []string{}, Verb: "free", Params: map[string]string{}, Flags: []string{}}},
	{In: "nest abc123 def456 in ghi789", Err: false, Out: InputAttributes{ResourceType: "item", ResourceId: "", ResourceIds: []string{"abc123", "def456"}, SecondaryIds: []string{"ghi789"}, Verb: "nest", Params: map[string]string{}, Flags: []string{}}},
	{In: "item fetch abc123", Err: false, Out: InputAttributes{ResourceType: "item", ResourceId: "abc123", ResourceIds: []string{}, SecondaryIds: []string{}, Verb: "fetch", Params: map[string]string{}, Flags: []string{}}},
	{In: "rel fetch abc123 def456", Err: false, Out: InputAttributes{ResourceType: "rel", ResourceId: "abc123", ResourceIds: []string{}, SecondaryIds: []string{"def456"}, Verb: "fetch", Params: map[string]string{}, Flags: []string{}}},
	{In: "item in abc123", Err: false, Out: InputAttributes{ResourceType: "item", ResourceId: "abc123", ResourceIds: []string{}, SecondaryIds: []string{}, Verb: "list", Params: map[string]string{}, Flags: []string{}}},
	{In: "items in abc123", Err: false, Out: InputAttributes{ResourceType: "item", ResourceId: "abc123", ResourceIds: []string{}, SecondaryIds: []string{}, Verb: "list", Params: map[string]string{}, Flags: []string{}}},
	{In: "world", Err: false, Out: InputAttributes{ResourceType: "world", ResourceId: "", ResourceIds: []string{}, SecondaryIds: []string{}, Verb: "fetch", Params: map[string]string{}, Flags: []string{}}},
	{In: "item list 10", Err: false, Out: InputAttributes{ResourceType: "item", ResourceId: "", ResourceIds: []string{}, SecondaryIds: []string{}, Verb: "list", Params: map[string]string{"limit": "10"}, Flags: []string{}}},
	{In: "items list 10", Err: false, Out: InputAttributes{ResourceType: "item", ResourceId: "", ResourceIds: []string{}, SecondaryIds: []string{}, Verb: "list", Params: map[string]string{"limit": "10"}, Flags: []string{}}},
	{In: "rel list", Err: false, Out: InputAttributes{ResourceType: "rel", ResourceId: "", ResourceIds: []string{}, SecondaryIds: []string{}, Verb: "list", Params: map[string]string{}, Flags: []string{}}},
	{In: "to? abc123", Err: false, Out: InputAttributes{ResourceType: "rel", ResourceId: "abc123", ResourceIds: []string{}, SecondaryIds: []string{}, Verb: "to?", Params: map[string]string{}, Flags: []string{}}},
	{In: "from? abc123", Err: false, Out: InputAttributes{ResourceType: "rel", ResourceId: "abc123", ResourceIds: []string{}, SecondaryIds: []string{}, Verb: "from?", Params: map[string]string{}, Flags: []string{}}},
	{In: "in? abc123 def456", Err: false, Out: InputAttributes{ResourceType: "item", ResourceId: "abc123", ResourceIds: []string{}, SecondaryIds: []string{"def456"}, Verb: "in?", Params: map[string]string{}, Flags: []string{}}},
	{In: "item? abc123", Err: false, Out: InputAttributes{ResourceType: "item", ResourceId: "abc123", ResourceIds: []string{}, SecondaryIds: []string{}, Verb: "exists", Params: map[string]string{}, Flags: []string{}}},
	{In: "rel? abc123 def456", Err: false, Out: InputAttributes{ResourceType: "rel", ResourceId: "abc123", ResourceIds: []string{}, SecondaryIds: []string{"def456"}, Verb: "exists", Params: map[string]string{}, Flags: []string{}}},
	{In: "item abc123", Err: false, Out: InputAttributes{ResourceType: "item", ResourceId: "abc123", ResourceIds: []string{}, SecondaryIds: []string{}, Verb: "create-or-fetch", Params: map[string]string{}, Flags: []string{}}},
	{In: "rel abc123 def456", Err: false, Out: InputAttributes{ResourceType: "rel", ResourceId: "abc123", ResourceIds: []string{}, SecondaryIds: []string{"def456"}, Verb: "create-or-fetch", Params: map[string]string{}, Flags: []string{}}},
	{In: "item abc123 name=John", Err: false, Out: InputAttributes{ResourceType: "item", ResourceId: "abc123", ResourceIds: []string{}, SecondaryIds: []string{}, Verb: "create-or-set", Params: map[string]string{"name": "John"}, Flags: []string{}}},
	{In: "rel abc123 def456 verb=likes", Err: false, Out: InputAttributes{ResourceType: "rel", ResourceId: "abc123", ResourceIds: []string{}, SecondaryIds: []string{"def456"}, Verb: "create-or-set", Params: map[string]string{"verb": "likes"}, Flags: []string{}}},
	{In: "item create abc123 -strict --verbose", Err: false, Out: InputAttributes{ResourceType: "item", ResourceId: "abc123", ResourceIds: []string{}, SecondaryIds: []string{}, Verb: "create", Params: map[string]string{}, Flags: []string{"strict", "verbose"}}},
}

var testResponses = []struct {
	In  string
	Err bool
	Out Response
}{
	{In: "item abc123\n$$$$\n200 ok \"all ok\"", Err: false, Out: Response{Object: ResponseObject{Type: "item", Repr: `item abc123`}, Status: ResponseStatus{Code: 200, Message: "all ok"}}},
	{In: "rel abc123 \"def456\" verb=\"writes to\" async=false\n$$$$\n200 ok \"all ok\"", Err: false, Out: Response{Object: ResponseObject{Type: "rel", Repr: `rel abc123 "def456" verb="writes to" async=false`}, Status: ResponseStatus{Code: 200, Message: "all ok"}}},
	{In: "$$world\n" + simpleWorld + "\nendworld$$" + "\n$$$$\n200 ok ", Err: false, Out: Response{Object: ResponseObject{Type: "world", Repr: simpleWorld}, Status: ResponseStatus{Code: 200, Message: ""}}},
	{In: "a b c d" + "\n$$$$\n200 ok ", Err: false, Out: Response{Object: ResponseObject{Type: "ids", Repr: `["a","b","c","d"]`}, Status: ResponseStatus{Code: 200, Message: ""}}},
	{In: "1 2 3 4" + "\n$$$$\n200 ok ", Err: false, Out: Response{Object: ResponseObject{Type: "ids", Repr: `["1","2","3","4"]`}, Status: ResponseStatus{Code: 200, Message: ""}}},
}

var testTrees = []struct {
	In   string
	Err  bool
	Tree Node
}{
	{In: simpleTree, Err: false, Tree: Node{
		"nil", []Node{
			{Id: "2", Children: []Node{
				{Id: "1", Children: []Node{}},
			}},
			{Id: "3", Children: []Node{}},
		}}},
}

func TestCommands(t *testing.T) {
	for _, c := range testCommands {
		t.Run(c.In, func(t *testing.T) {
			p, err := Parse(c.In)

			if p.StmtType != "Command" {
				t.Errorf("expected StmtType to be 'Command', but got: '%s'", p.StmtType)
			}

			// All commands should have a non-empty values for resource and verb (aka action type).
			if p.InputAttributes.ResourceType == "" {
				t.Errorf("expected InputAttributes.ResourceType, but got ''")
			}
			if p.InputAttributes.Verb == "" {
				t.Errorf("expected InputAttributes.Verb, but got ''")
			}
			if p.InputAttributes.Raw != c.In {
				t.Errorf("expected InputAttributes.Raw to equal the input text")
			}

			// Checking specific output.
			if p.InputAttributes.ResourceType != c.Out.ResourceType {
				t.Errorf("expected resource type: '%s', but got: '%s'", c.Out.ResourceType, p.InputAttributes.ResourceType)
			}
			if p.InputAttributes.ResourceId != c.Out.ResourceId {
				t.Errorf("expected resource id: '%s', but got: '%s'", c.Out.ResourceId, p.InputAttributes.ResourceId)
			}
			if !mapset.NewSet(p.InputAttributes.ResourceIds...).Equal(mapset.NewSet(c.Out.ResourceIds...)) {
				fmt.Println(fmt.Sprintf("PARSED: %v\nEXPECT: %v", p.InputAttributes.ResourceIds, c.Out.ResourceIds))
				t.Errorf("expected resource ids: '%v', but got: '%v'", c.Out.ResourceIds, p.InputAttributes.ResourceIds)
			}
			if !mapset.NewSet(p.InputAttributes.SecondaryIds...).Equal(mapset.NewSet(c.Out.SecondaryIds...)) {
				fmt.Println(fmt.Sprintf("PARSED: %v\nEXPECT: %v", p.InputAttributes.SecondaryIds, c.Out.SecondaryIds))
				t.Errorf("expected secondary ids: '%v', but got: '%v'", c.Out.SecondaryIds, p.InputAttributes.SecondaryIds)
			}
			if p.InputAttributes.Verb != c.Out.Verb {
				t.Errorf("expected verb: '%s', but got: '%s'", c.Out.Verb, p.InputAttributes.Verb)
			}
			if !reflect.DeepEqual(p.InputAttributes.Params, c.Out.Params) {
				fmt.Println(fmt.Sprintf("PARSED: %v\nEXPECT: %v", p.InputAttributes.Params, c.Out.Params))
				t.Errorf("expected params: '%v', but got: '%v'", c.Out.Params, p.InputAttributes.Params)
			}
			if !mapset.NewSet(p.InputAttributes.Flags...).Equal(mapset.NewSet(c.Out.Flags...)) {
				fmt.Println(fmt.Sprintf("PARSED: %v\nEXPECT: %v", p.InputAttributes.Flags, c.Out.Flags))
				t.Errorf("expected flags: '%v', but got: '%v'", c.Out.Flags, p.InputAttributes.Flags)
			}

			// Checking whether we expected to catch an error.
			if c.Err && err == nil {
				if err == nil {
					t.Errorf("expected error for command: '%s', but got none", c.In)
				}
			}
			if !c.Err && err != nil {
				if err != nil {
					t.Errorf("did not expect error for command: '%s', but got: %s", c.In, err)
				}
			}

			if t.Failed() {
				p.PrintSyntaxTree()
			}
		})
	}
}

func TestResponses(t *testing.T) {
	for _, c := range testResponses {
		t.Run(c.In, func(t *testing.T) {
			p, err := Parse(c.In)

			if p.StmtType != "Response" {
				t.Errorf("expected StmtType to be 'Response', but got: '%s'", p.StmtType)
			}

			// Checking specific output.
			if p.Response.Object.Type != c.Out.Object.Type {
				t.Errorf("expected object type: '%s', but got: '%s'", c.Out.Object.Type, p.Response.Object.Type)
			}
			if p.Response.Object.Repr != c.Out.Object.Repr {
				t.Errorf("expected object repr:\n%s\n\ngot:\n%s", c.Out.Object.Repr, p.Response.Object.Repr)
			}
			if p.Response.Status.Code != c.Out.Status.Code {
				t.Errorf("expected status code: '%d', but got: '%d'", c.Out.Status.Code, p.Response.Status.Code)
			}
			if p.Response.Status.Message != c.Out.Status.Message {
				t.Errorf("expected status message: '%s', but got: '%s'", c.Out.Status.Message, p.Response.Status.Message)
			}

			// Checking whether we expected to catch an error.
			if c.Err && err == nil {
				if err == nil {
					t.Errorf("expected error for response: '%s', but got none", c.In)
				}
			}
			if !c.Err && err != nil {
				if err != nil {
					t.Errorf("did not expect error for response: '%s', but got: %s", c.In, err)
				}
			}

			if t.Failed() {
				p.PrintSyntaxTree()
			}
		})
	}
}

func TestTrees(t *testing.T) {
	for _, testTree := range testTrees {
		t.Run(testTree.In, func(t *testing.T) {
			p, err := Parse(testTree.In)

			if p.StmtType != "Tree" {
				t.Errorf("expected StmtType to be 'Tree', but got: '%s'", p.StmtType)
			}
			if p.Response.Object.Type != "tree" {
				t.Errorf("expected StmtType to be 'Tree', but got: '%s'", p.StmtType)
			}

			// Checking specific output.
			if !reflect.DeepEqual(p.Tree, testTree.Tree) {
				t.Errorf("expected tree map: '%v', but got: '%v'", testTree.Tree, p.Tree)
			}

			// Checking whether we expected to catch an error.
			if testTree.Err && err == nil {
				if err == nil {
					t.Errorf("expected error for tree: '%s', but got none", testTree.In)
				}
			}
			if !testTree.Err && err != nil {
				if err != nil {
					t.Errorf("did not expect error for tree: '%s', but got: %s", testTree.In, err)
				}
			}

			if t.Failed() {
				p.PrintSyntaxTree()
			}
		})
	}
}

func TestStandaloneErrors(t *testing.T) {
	testErr := errors.New("test error").UseCode(errors.TopolithErrorInternal).WithError(errors.New("inner error")).WithDescription("test description").WithData(errors.KvPair{Key: "key", Value: "value"})
	p, err := Parse(testErr.String())
	if err != nil {
		t.Fatalf("error parsing error string: %v", err)
	}
	if p.StmtType != "Status" {
		t.Errorf("expected StmtType to be 'Status', but got: '%s'", p.StmtType)
	}
}
