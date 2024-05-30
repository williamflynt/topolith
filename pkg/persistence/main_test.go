package persistence

import (
	"fmt"
	"github.com/williamflynt/topolith/pkg/world"
	"testing"
)

func TestEndToEndFileSaveLoad(t *testing.T) {
	dir := t.TempDir()
	fp := &filePersistence{directory: dir}
	fp.SetSourcePath(dir)

	worldName := "test-world"
	w := world.CreateWorld(worldName)

	// Add items and relationships
	for i := 1; i <= 10; i++ {
		_ = w.ItemCreate(fmt.Sprintf("item-%d", i), world.ItemParams{})
		if i > 1 {
			_ = w.RelCreate(fmt.Sprintf("item-%d", i), fmt.Sprintf("item-%d", i-1), world.RelParams{})
		}
	}

	// Nest items
	for i := 1; i <= 5; i++ {
		w.Nest(fmt.Sprintf("item-%d", i), fmt.Sprintf("item-%d", i+5))
	}

	// Test
	err := fp.Save(w)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load
	w2, err := fp.Load(worldName)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Assert
	if !world.WorldEqual(w, w2) {
		t.Fatalf("Worlds are not equal")
	}

	if t.Failed() {
		fmt.Printf("\tOriginal:\n\n%s", w.String())
		fmt.Printf("\tLoaded:\n\n%s", w2.String())
	}
}
