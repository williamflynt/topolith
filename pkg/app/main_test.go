package app

import (
	"github.com/williamflynt/topolith/pkg/grammar"
	"github.com/williamflynt/topolith/pkg/world"
	"testing"
)

func TestBareWorldResponse(t *testing.T) {
	testApp, err := NewApp(world.CreateWorld("test-world"))
	if err != nil {
		t.Fatalf("error creating app: %v", err)
	}
	response := testApp.Exec("world")
	// TODO - a perfectly fine WORLDOBJECT response is parsing as a Command (world fetch) here - why?
	p, err := grammar.Parse(response)
	if err != nil {
		t.Fatalf("error parsing response: %v", err)
	}
	if p.Response.Status.Code != 200 {
		t.Fatalf("expected 200 status code, got %d", p.Response.Status.Code)
	}
	if p.Response.Object.Type != "world" {
		t.Fatalf("expected world type, got %s", p.Response.Object.Type)
	}
}
