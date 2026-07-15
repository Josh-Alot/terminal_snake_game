package input

import (
	"os"

	"golang.org/x/term"
)

type State = term.State
type realController struct{}

type rawModeController interface {
	MakeRaw(fd int) (*State, error)
	Restore(fd int, state *State) error
}

func (realController) MakeRaw(fd int) (*State, error) {
	return term.MakeRaw(fd)
}

func (realController) Restore(fd int, state *State) error {
	return term.Restore(fd, state)
}

func EnableRawMode() (func(), error) {
	return enableRawModeWith(realController{}, int(os.Stdin.Fd()))
}

func enableRawModeWith(controller rawModeController, fd int) (func(), error) {
	oldState, err := controller.MakeRaw(fd)
	if err != nil {
		return nil, err
	}

	restore := func() {
		_ = controller.Restore(fd, oldState)
	}

	return restore, nil
}
