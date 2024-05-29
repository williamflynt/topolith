package world

import (
	"fmt"
	"github.com/williamflynt/topolith/pkg/errors"
	"github.com/williamflynt/topolith/pkg/grammar"
	"strings"
)

const RelIdSeparator = "::"

// Rel is a struct that represents a relationship between two Item. It might be represented in diagrams as an arrow between two Item renderings.
type Rel struct {
	From      Item   `json:"from"`      // From is the source Item.
	To        Item   `json:"to"`        // To is the destination Item.
	Verb      string `json:"verb"`      // Verb is the action that the relationship represents (ex: reads, uses).
	Mechanism string `json:"mechanism"` // Mechanism is the method of implementation of the relationship (ex: HTTPS, JSON).
	Async     bool   `json:"async"`     // Async is a boolean that represents whether the relationship is asynchronous.
	Expanded  string `json:"expanded"`  // Expanded is expanded information on the relationship.
}

func (r Rel) String() string {
	rel := fmt.Sprintf(`rel "%s" "%s"`, r.From.Id, r.To.Id)
	paramRepr := make([]string, 0)
	if r.Verb != "" {
		paramRepr = append(paramRepr, fmt.Sprintf(`verb="%s"`, r.Verb))
	}
	if r.Mechanism != "" {
		paramRepr = append(paramRepr, fmt.Sprintf(`mechanism="%s"`, r.Mechanism))
	}
	if r.Async {
		paramRepr = append(paramRepr, fmt.Sprintf(`async=%t`, r.Async))
	}
	if r.Expanded != "" {
		paramRepr = append(paramRepr, fmt.Sprintf(`expanded="%s"`, r.Expanded))
	}
	if len(paramRepr) > 0 {
		rel += " " + strings.Join(paramRepr, " ")
	}
	if _, err := grammar.Parse(rel); err != nil {
		panic(err)
	}
	return rel
}

func RelFromString(w World, s string) (Rel, error) {
	p, err := grammar.Parse(s)
	if err != nil {
		return Rel{}, fmt.Errorf("error parsing Rel: %w", err)
	}
	fromItem, ok := w.ItemFetch(p.InputAttributes.ResourceId)
	if !ok {
		return Rel{}, errors.New("from Item not found").UseCode(errors.TopolithErrorNotFound).WithDescription("FromItem not found").WithData(errors.KvPair{Key: "input", Value: s}, errors.KvPair{Key: "fromId", Value: p.InputAttributes.ResourceId})
	}
	toItem, ok := w.ItemFetch(p.InputAttributes.SecondaryIds[0])
	if !ok {
		return Rel{}, errors.New("to Item not found").UseCode(errors.TopolithErrorNotFound).WithDescription("ToItem not found").WithData(errors.KvPair{Key: "input", Value: s}, errors.KvPair{Key: "toId", Value: p.InputAttributes.SecondaryIds[0]})
	}
	return relSet(Rel{From: fromItem, To: toItem}, RelParamsFromInput(p.InputAttributes))
}

func RelEqual(r1, r2 Rel) bool {
	if !ItemEqual(r1.From, r2.From) {
		return false
	}
	if !ItemEqual(r1.To, r2.To) {
		return false
	}
	if r1.Verb != r2.Verb {
		return false
	}
	if r1.Mechanism != r2.Mechanism {
		return false
	}
	if r1.Async != r2.Async {
		return false
	}
	if r1.Expanded != r2.Expanded {
		return false
	}
	return true
}

// id returns the ID of the Rel.
func (r Rel) id() string {
	return relIdFromIds(r.From.Id, r.To.Id)
}

func relIdFromIds(fromId, toId string) string {
	return fromId + RelIdSeparator + toId
}

// RelParams is a struct that represents the parameters that can be set on a Rel.
type RelParams struct {
	Verb      *string `json:"verb"`
	Mechanism *string `json:"mechanism"`
	Async     *bool   `json:"async"`
	Expanded  *string `json:"expanded"`
}

func RelParamsFromInput(input grammar.InputAttributes) RelParams {
	params := RelParams{}
	if v, ok := input.Params["verb"]; ok {
		params.Verb = strPtr(v)
	}
	if v, ok := input.Params["mechanism"]; ok {
		params.Mechanism = strPtr(v)
	}
	if v, ok := input.Params["async"]; ok {
		params.Async = boolPtr(v == "true")
	}
	if v, ok := input.Params["expanded"]; ok {
		params.Expanded = strPtr(v)
	}
	return params
}
