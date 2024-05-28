package world

import "testing"

func TestRelFromString(t *testing.T) {
	item1 := Item{Id: "abc123"}
	item2 := Item{Id: "this is my rel"}
	sampleRel := "rel abc123 \"this is my rel\" verb=reads mechanism=HTTPS async=true expanded=\"this is expanded\""
	rel, err := RelFromString(&world{
		Version_:  1,
		Id_:       "sampleWorld",
		Name_:     "Sample World",
		Expanded_: "This is a world for testing!",
		Items:     map[string]Item{"abc123": item1, "this is my rel": item2},
		Rels:      make(map[string]Rel),
	}, sampleRel)
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if rel.From.Id != "abc123" {
		t.Errorf("expected from to be 'abc123', got '%s'", rel.From.Id)
	}
	if rel.To.Id != "this is my rel" {
		t.Errorf("expected to to be 'this is my rel', got '%s'", rel.To.Id)
	}
	if rel.Verb != "reads" {
		t.Errorf("expected verb to be 'reads', got '%s'", rel.Verb)
	}
	if rel.Mechanism != "HTTPS" {
		t.Errorf("expected mechanism to be 'HTTPS', got '%s'", rel.Mechanism)
	}
	if !rel.Async {
		t.Error("expected async to be true")
	}
	if rel.Expanded != "this is expanded" {
		t.Errorf("expected expanded to be 'this is expanded', got '%s'", rel.Expanded)
	}
}

func TestRelSerde(t *testing.T) {
	item1 := Item{Id: "abc123"}
	item2 := Item{Id: "this is my rel"}
	rel := Rel{
		From:      item1,
		To:        item2,
		Verb:      "reads",
		Mechanism: "HTTPS",
		Async:     true,
		Expanded:  "this is expanded",
	}
	serialized := rel.String()
	deserialized, err := RelFromString(&world{
		Version_:  1,
		Id_:       "sampleWorld",
		Name_:     "Sample World",
		Expanded_: "This is a world for testing!",
		Items:     map[string]Item{"abc123": item1, "this is my rel": item2},
		Rels:      make(map[string]Rel),
	}, serialized)
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if rel.From.Id != deserialized.From.Id {
		t.Errorf("expected from to be 'abc123', got '%s'", deserialized.From.Id)
	}
	if rel.To.Id != deserialized.To.Id {
		t.Errorf("expected to to be 'this is my rel', got '%s'", deserialized.To.Id)
	}
	if rel.Verb != deserialized.Verb {
		t.Errorf("expected verb to be 'reads', got '%s'", deserialized.Verb)
	}
	if rel.Mechanism != deserialized.Mechanism {
		t.Errorf("expected mechanism to be 'HTTPS', got '%s'", deserialized.Mechanism)
	}
	if rel.Async != deserialized.Async {
		t.Errorf("expected async to be true, got '%t'", deserialized.Async)
	}
	if rel.Expanded != deserialized.Expanded {
		t.Errorf("expected expanded to be 'this is expanded', got '%s'", deserialized.Expanded)
	}
	serialized2 := deserialized.String()
	if serialized != serialized2 {
		t.Errorf("expected serialized to be:\n%s\ngot: \n%s\n", serialized, serialized2)
	}
}
