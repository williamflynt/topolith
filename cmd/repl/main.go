package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/williamflynt/topolith/pkg/app"
	"github.com/williamflynt/topolith/pkg/grammar"
	"github.com/williamflynt/topolith/pkg/world"
	"strings"
)

// executor handles the unparsed input to the REPL.
func executor(app app.App) prompt.Executor {
	return func(input string) {
		input = strings.TrimSpace(input)
		if input == "" {
			return
		}

		p, err := grammar.Parse(input)
		if err != nil {
			fmt.Println("invalid input:", err)
			p.PrintSyntaxTree()
			return
		}

		resp := app.Exec(input)
		fmt.Println(resp)
	}
}

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

func main() {
	core, err := app.NewApp(world.CreateWorld("default-world"))
	if err != nil {
		fmt.Println("error creating app:", err)
		return
	}

	fmt.Println("Interactive console. Type 'Ctrl-D' to quit.")
	p := prompt.New(
		executor(core),
		completer(core),
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("topolith"),
	)
	p.Run()
}
