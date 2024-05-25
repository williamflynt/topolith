package main

import (
	"github.com/c-bata/go-prompt"
	"github.com/williamflynt/topolith/pkg/app"
)

// completer handles the autocompletion for the REPL.
func completer(app app.App) prompt.Completer {
	return func(d prompt.Document) []prompt.Suggest {
		text := d.TextBeforeCursor()

		suggestions := []prompt.Suggest{
			{Text: ".save", Description: "Store the world"},
			{Text: ".load", Description: "Load a new world"},
			{Text: "item", Description: "Manage items"},
			{Text: "rel", Description: "Manage relationships"},
			{Text: "world", Description: "Manage the world"},
			{Text: "in?", Description: "Check item containment"},
			{Text: "nest", Description: "Nest items"},
			{Text: "free", Description: "Free items"},
			{Text: "undo", Description: "Undo last action"},
			{Text: "redo", Description: "Redo reversed action"},
		}

		return prompt.FilterHasPrefix(suggestions, text, true)
	}
}
