package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"strings"
)

// Executor function to handle commands.
func executor(input string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}

	args := strings.Fields(input)
	command := args[0]
	switch command {
	case "save":
		handleWorldCommand(args[1:])
	case "load":
		handleWorldCommand(args[1:])
	case "world":
		handleWorldCommand(args[1:])
	case "item":
		handleItemCommand(args[1:])
	case "rel":
		handleRelCommand(args[1:])
	case "in?":
		handleInCommand(args[1:])
	case "nest":
		handleNestCommand(args[1:])
	case "free":
		handleFreeCommand(args[1:])
	case "undo":
		handleUndoCommand(args[1:])
	case "redo":
		handleUndoCommand(args[1:])
	default:
		fmt.Println("Unknown command:", command)
	}
}

func handleSaveCommand(args []string) {
	fmt.Println("Save command:", args)
}

func handleLoadCommand(args []string) {
	fmt.Println("Load command:", args)
}

func handleWorldCommand(args []string) {
	fmt.Println("World command:", args)
}

func handleItemCommand(args []string) {
	fmt.Println("Item command:", args)
}

func handleRelCommand(args []string) {
	fmt.Println("Rel command:", args)
}

func handleInCommand(args []string) {
	fmt.Println("In? command:", args)
}

func handleNestCommand(args []string) {
	fmt.Println("Nest command:", args)
}

func handleFreeCommand(args []string) {
	fmt.Println("Free command:", args)
}

func handleUndoCommand(args []string) {
	fmt.Println("Undo command:", args)
}

func handleRedoCommand(args []string) {
	fmt.Println("Redo command:", args)
}

// Completer function for autocompletion
func completer(d prompt.Document) []prompt.Suggest {
	text := d.TextBeforeCursor()
	words := strings.Fields(text)

	if len(words) == 0 {
		return []prompt.Suggest{
			{Text: "save", Description: "Store the world"},
			{Text: "load", Description: "Load a new world"},
			{Text: "world", Description: "Manage the world"},
			{Text: "item", Description: "Manage items"},
			{Text: "rel", Description: "Manage relationships"},
			{Text: "in?", Description: "Check item containment"},
			{Text: "nest", Description: "Nest items"},
			{Text: "free", Description: "Free items"},
			{Text: "undo", Description: "Undo last action"},
			{Text: "redo", Description: "Undo last action"},
		}
	}

	command := words[0]

	switch command {
	case "world":
		return prompt.FilterHasPrefix([]prompt.Suggest{
			{Text: "--pretty", Description: "Pretty print the world"},
		}, d.GetWordBeforeCursor(), true)

	case "item":
		return prompt.FilterHasPrefix([]prompt.Suggest{
			// TODO: get IDs from the world
		}, d.GetWordBeforeCursor(), true)

	case "rel":
		return prompt.FilterHasPrefix([]prompt.Suggest{
			// TODO: get IDs from the world
		}, d.GetWordBeforeCursor(), true)

	case "in?":
		return prompt.FilterHasPrefix([]prompt.Suggest{
			// TODO: get IDs from the world
		}, d.GetWordBeforeCursor(), true)

	case "nest":
		return prompt.FilterHasPrefix([]prompt.Suggest{
			// TODO: get IDs from the world
		}, d.GetWordBeforeCursor(), true)

	case "free":
		return prompt.FilterHasPrefix([]prompt.Suggest{
			// TODO: get IDs from the world
		}, d.GetWordBeforeCursor(), true)

	case "undo":
		return []prompt.Suggest{}

	case "redo":
		return []prompt.Suggest{}

	default:
		return []prompt.Suggest{}
	}
}

func main() {
	fmt.Println("Interactive console. Type 'exit' or 'Ctrl-D' to quit.")
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle("interactive-console"),
	)
	p.Run()
}
