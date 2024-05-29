package persistance

import (
	"github.com/williamflynt/topolith/pkg/world"
	"testing"
)

func TestEndToEndFileSaveLoad(t *testing.T) {
	// Setup
	dir := t.TempDir()
	fp := &filePersistence{directory: dir}
	fp.SetSourcePath(dir)

	w := world.CreateWorld("test-world")

	// Test
	err := fp.Save(w)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load
	w2, err := fp.Load("test-world")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Assert
	if w.Name() != w2.Name() {
		t.Fatalf("World names don't match: %s != %s", w.Name(), w2.Name())
	}
}
