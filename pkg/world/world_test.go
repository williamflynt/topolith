package world

import (
	"fmt"
	"strings"
	"testing"
)

var simpleWorld = "$$world\nversion=1\nid=1\nname=worldname\nexpanded=\"this is expanded data\"\ntree{nil::[tree{item \"2\"::[tree{item \"1\"::[]}]} tree{item \"3\"::[]}]}\nrel \"3\" \"2\"\nrel \"1\" \"2\"\nendworld$$"
var simpleWorld2 = "$$world\nversion=1\nid=test-world\nname=\nexpanded=\ntree{nil::[tree{item \"item-6\"::[tree{item \"item-1\"::[]}]} tree{item \"item-8\"::[tree{item \"item-3\"::[]}]} tree{item \"item-9\"::[tree{item \"item-4\"::[]}]} tree{item \"item-10\"::[tree{item \"item-5\"::[]}]} tree{item \"item-7\"::[tree{item \"item-2\"::[]}]}]}\nrel \"item-2\" \"item-1\"\nrel \"item-4\" \"item-3\"\nrel \"item-5\" \"item-4\"\nrel \"item-6\" \"item-5\"\nrel \"item-10\" \"item-9\"\nrel \"item-3\" \"item-2\"\nrel \"item-7\" \"item-6\"\nrel \"item-8\" \"item-7\"\nrel \"item-9\" \"item-8\"\nendworld$$"

var worlds = []string{simpleWorld, simpleWorld2}

func TestWorldSerde(t *testing.T) {
	for i, s := range worlds {
		t.Run(fmt.Sprintf(`TestWorldSerde-%d`, i), func(t *testing.T) {
			w, err := FromString(s)
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
		})
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
