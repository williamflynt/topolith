package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/williamflynt/topolith/pkg/app"
	"github.com/williamflynt/topolith/pkg/grammar"
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

		resp, err := app.Exec(input)
		fmt.Println(resp)
		if err != nil {
			fmt.Println(err)
		}
	}
}
