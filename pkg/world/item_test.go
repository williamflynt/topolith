package world

import (
	"fmt"
	"testing"
)

var testItems = []Item{
	{Id: "abc123", Name: "Test Item", Expanded: "This is a test item.", External: true},
	{Id: "abc123", Name: "Test Item", Expanded: "This is a test item."},
}

func TestItemSerde(t *testing.T) {
	for i, item := range testItems {
		t.Run(fmt.Sprintf(`TestItemSerde-%d`, i), func(t *testing.T) {
			serialized := item.String()
			deserialized, err := ItemFromString(serialized)
			if err != nil {
				t.Error("unexpected error:", err)
			}
			if item.Id != deserialized.Id {
				t.Errorf("expected ID to be 'abc123', got '%s'", deserialized.Id)
			}
			if item.Name != deserialized.Name {
				t.Errorf("expected name to be 'Test Item', got '%s'", deserialized.Name)
			}
			if item.Expanded != deserialized.Expanded {
				t.Errorf("expected expanded to be 'This is a test item.', got '%s'", deserialized.Expanded)
			}
			if item.External != deserialized.External {
				t.Errorf("expected external to be true, got '%t'", deserialized.External)
			}
			serialized2 := deserialized.String()
			if serialized != serialized2 {
				t.Errorf("expected serialized to be:\n%s\ngot: \n%s\n", serialized, serialized2)
			}
		})
	}
}
