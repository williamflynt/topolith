package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/williamflynt/topolith/pkg/app"
	"github.com/williamflynt/topolith/pkg/world"
)

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
