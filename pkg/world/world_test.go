package world

import (
	"fmt"
	"strings"
	"testing"
)

var simpleWorld = "$$world\nversion=1\nid=1\nname=worldname\nexpanded=\"this is expanded data\"\ntree{nil::[tree{item \"2\" external=false::[tree{item \"1\" external=false::[]}]} tree{item \"3\" external=false::[]}]}\nrel \"3\" \"2\"\nrel \"1\" \"2\"\nendworld$$"

func TestWorldSerde(t *testing.T) {
	w, err := FromString(simpleWorld)
	if err != nil {
		t.Fatalf("FromString failed: %v", err)
	}

	// Serialize and deserialize the world to get a new world.
	ser := w.String()
	w2, err := FromString(ser)
	if err != nil {
		t.Fatalf("FromString failed: %v", err)
	}

	if !WorldEqual(w, w2) {
		t.Fatalf("Worlds are not equal")
	}

	if t.Failed() {
		printDiff(w.String(), w2.String())
	}
}

// --- HELPERS ---

func printDiff(a, b string) {
	lines1 := strings.Split(a, "\n")
	lines2 := strings.Split(b, "\n")
	for i := 0; i < len(lines1) || i < len(lines2); i++ {
		if i >= len(lines1) {
			fmt.Printf("+ %s\n", lines2[i])
		} else if i >= len(lines2) {
			fmt.Printf("- %s\n", lines1[i])
		} else if lines1[i] != lines2[i] {
			fmt.Printf("- %s\n", lines1[i])
			fmt.Printf("+ %s\n", lines2[i])
		}
	}
}
