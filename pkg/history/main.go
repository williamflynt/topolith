package history

import (
	"github.com/williamflynt/topolith/pkg/errors"
	"github.com/williamflynt/topolith/pkg/topolith"
)

type History interface {
	World() topolith.World // World returns the topolith.World associated with this History.
	Commands() []Command   // Commands returns the list of Command that have been executed for the present state of the topolith.World.
	Undo() (error, int)    // Undo reverses the last operation on the World. If there are no operations to undo, noop. Return any error that occurred and the number of operations left to undo.
	Redo() (error, int)    // Redo executes the most recently reversed operation on the World. If there are no operations to redo, noop. Return any error that occurred and the number of operations left to redo.
	CanUndo() bool         // CanUndo indicates whether more Command objects exist to Undo.
	CanRedo() bool         // CanRedo indicates whether more Command objects exist to Redo.
}

// history implements History.
type history struct {
	world       topolith.World // world is the topolith.World associated with this History.
	commands    []Command      // commands is a list of Command that have been executed.
	commandsIdx int            // commandsIdx is the index of the last executed Command in the commands list. It must initialize to -1.
}

func NewHistory(world topolith.World) (History, error) {
	if world == nil {
		return nil, errors.New("cannot create History with nil World").UseCode(errors.TopolithErrorInvalid)
	}
	return history{
		world:       world,
		commands:    make([]Command, 0),
		commandsIdx: -1,
	}, nil
}

func (h history) World() topolith.World {
	return h.world
}

func (h history) Commands() []Command {
	return h.commands[:h.commandsIdx+1]
}

func (h history) Undo() (error, int) {
	if h.commandsIdx < 0 {
		return nil, 0
	}
	if err := h.commands[h.commandsIdx].Undo(h.world); err != nil {
		// We aren't going to validate state of the World. But a problem happened.
		// Clear commands, reset commandsIdx, and return the error.
		h.commands = make([]Command, 0)
		h.commandsIdx = -1
		return err, 0
	}
	h.commandsIdx--
	return nil, len(h.commands) - h.commandsIdx - 1
}

func (h history) Redo() (error, int) {
	if h.commandsIdx >= len(h.commands)-1 {
		return nil, 0
	}
	if err := h.commands[h.commandsIdx].Execute(h.world); err != nil {
		// We aren't going to validate state of the World. But a problem happened.
		// Clear commands, reset commandsIdx, and return the error.
		h.commands = make([]Command, 0)
		h.commandsIdx = -1
		return err, 0
	}
	h.commandsIdx++
	return nil, len(h.commands) - h.commandsIdx - 1
}

func (h history) CanUndo() bool {
	return h.commandsIdx >= 0
}

func (h history) CanRedo() bool {
	return h.commandsIdx < len(h.commands)-1
}
