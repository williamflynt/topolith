package app

import (
	"github.com/williamflynt/topolith/pkg/errors"
	"github.com/williamflynt/topolith/pkg/world"
)

type App interface {
	World() world.World // World returns the world.World associated with this App.
	History() []Command // Commands returns the list of Command that have been executed for the present state of the world.World.
	Undo() (error, int) // Undo reverses the last operation on the World. If there are no operations to undo, noop. Return any error that occurred and the number of operations left to undo.
	Redo() (error, int) // Redo executes the most recently reversed operation on the World. If there are no operations to redo, noop. Return any error that occurred and the number of operations left to redo.
	CanUndo() bool      // CanUndo indicates whether more Command objects exist to Undo.
	CanRedo() bool      // CanRedo indicates whether more Command objects exist to Redo.
}

// app implements App.
type app struct {
	world       world.World // world is the world.World associated with this App.
	commands    []Command   // commands is a list of Command that have been executed.
	commandsIdx int         // commandsIdx is the index of the last executed Command in the commands list. It must initialize to -1.
}

func NewHistory(world world.World) (App, error) {
	if world == nil {
		return nil, errors.New("cannot create App with nil World").UseCode(errors.TopolithErrorInvalid)
	}
	return &app{
		world:       world,
		commands:    make([]Command, 0),
		commandsIdx: -1,
	}, nil
}

func (h *app) World() world.World {
	return h.world
}

func (h *app) History() []Command {
	return h.commands[:h.commandsIdx+1]
}

func (h *app) Undo() (error, int) {
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

func (h *app) Redo() (error, int) {
	if h.commandsIdx >= len(h.commands)-1 {
		return nil, 0
	}
	if err := h.commands[h.commandsIdx].Execute(h.world); err != nil {
		// We aren't going to validate state of the World. But a problem happened.
		// Clear app, reset commandsIdx, and return the error.
		h.commands = make([]Command, 0)
		h.commandsIdx = -1
		return err, 0
	}
	h.commandsIdx++
	return nil, len(h.commands) - h.commandsIdx - 1
}

func (h *app) CanUndo() bool {
	return h.commandsIdx >= 0
}

func (h *app) CanRedo() bool {
	return h.commandsIdx < len(h.commands)-1
}
