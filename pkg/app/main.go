package app

import (
	"fmt"
	"github.com/williamflynt/topolith/pkg/errors"
	"github.com/williamflynt/topolith/pkg/grammar"
	"github.com/williamflynt/topolith/pkg/persistance"
	"github.com/williamflynt/topolith/pkg/world"
)

type App interface {
	World() world.World                   // World returns the world.World associated with this App.
	Exec(s string) string                 // Exec parses the given string to a valid Command and executes it. Return a string response in accordance with our grammar.
	History() []Command                   // History returns the list of Command that have been executed for the present state of the world.World.
	CanUndo() bool                        // CanUndo indicates whether more Command objects exist to Undo.
	CanRedo() bool                        // CanRedo indicates whether more Command objects exist to Redo.
	Persistence() persistance.Persistence // Persistence returns the persistance.Persistence object associated with this App.
}

func NewApp(world world.World) (App, error) {
	if world == nil {
		return nil, errors.New("cannot create App with nil World").UseCode(errors.TopolithErrorInvalid)
	}
	return &app{
		world:       world,
		commands:    make([]Command, 0),
		commandsIdx: -1,
		persistance: persistance.NewFilePersistence(),
	}, nil
}

// app implements App.
type app struct {
	world       world.World // world is the world.World associated with this App.
	commands    []Command   // commands is a list of Command that have been executed.
	commandsIdx int         // commandsIdx is the index of the last executed Command in the commands list. It must initialize to -1.
	persistance persistance.Persistence
}

func (h *app) World() world.World {
	return h.world
}

func (h *app) Exec(s string) string {
	p, err := grammar.Parse(s)
	if err != nil || p.StmtType != "Command" {
		if p != nil {
			p.PrintSyntaxTree()
		}
		return errors.New("invalid input").UseCode(errors.TopolithErrorInvalid).WithError(err).WithDescription("invalid input").WithData(errors.KvPair{Key: "input", Value: s}).String()

	}
	c, err := InputToCommand(p.InputAttributes)
	if err != nil {
		return errors.New("invalid input").UseCode(errors.TopolithErrorInvalid).WithError(err).WithDescription("invalid input").WithData(errors.KvPair{Key: "input", Value: s}).String()
	}
	stringerObj, err := h.exec(c)
	if err != nil {
		return errors.New("error executing command").UseCode(errors.TopolithErrorCommandErr).WithError(err).WithDescription("unexpected error executing command").WithData(errors.KvPair{Key: "input", Value: s}).String()
	}
	response := okString(stringerObj, err)
	if p, err := grammar.Parse(response); err != nil || p.StmtType != "Response" {
		if p != nil {
			p.PrintSyntaxTree()
		}
		return errors.New("error generating response").UseCode(errors.TopolithErrorInternal).WithError(err).WithDescription("unexpected error generating response").WithData(errors.KvPair{Key: "input", Value: s}).String()
	}
	return response
}

func (h *app) History() []Command {
	return h.commands[:h.commandsIdx+1]
}

func (h *app) CanUndo() bool {
	return h.commandsIdx >= 0
}

func (h *app) CanRedo() bool {
	return h.commandsIdx < len(h.commands)-1
}

func (h *app) Persistence() persistance.Persistence {
	return h.persistance
}

// --- INTERNAL ---

func (h *app) exec(c Command) (fmt.Stringer, error) {
	h.commands = append(h.commands, c)
	h.commandsIdx++
	return c.Execute(h.world)
}

func (h *app) undo() (error, int) {
	if h.commandsIdx < 0 {
		return nil, 0
	}
	if err := h.commands[h.commandsIdx].Undo(h.world); err != nil {
		// We aren't going to validate state of the World. But a problem happened.
		// Clear app, reset commandsIdx, and return the error.
		h.commands = make([]Command, 0)
		h.commandsIdx = -1
		return err, 0
	}
	h.commandsIdx--
	return nil, len(h.commands) - h.commandsIdx - 1
}

func (h *app) redo() (error, int) {
	if h.commandsIdx >= len(h.commands)-1 {
		return nil, 0
	}
	_, err := h.commands[h.commandsIdx].Execute(h.world)
	if err != nil {
		// We aren't going to validate state of the World. But a problem happened.
		// Clear app, reset commandsIdx, and return the error.
		h.commands = make([]Command, 0)
		h.commandsIdx = -1
		return err, 0
	}
	h.commandsIdx++
	return nil, len(h.commands) - h.commandsIdx - 1
}

func errOrEmpty(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func okString(o fmt.Stringer, err error) string {
	return fmt.Sprintf("%s\n200 ok %s", o, errOrEmpty(err))
}
