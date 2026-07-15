package input

import (
	"errors"
	"testing"
)

type mockController struct {
	makeRawCalled bool
	restoreCalled bool
	makeRawErr    error
	stateReturned *State
}

func (m *mockController) MakeRaw(fd int) (*State, error) {
	m.makeRawCalled = true
	return m.stateReturned, m.makeRawErr
}

func (m *mockController) Restore(fd int, state *State) error {
	m.restoreCalled = true
	return nil
}

func TestEnableRawMode_EnableAndReturnRestore(t *testing.T) {
	mock := &mockController{stateReturned: &State{}}

	restore, err := enableRawModeWith(mock, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// The restoration only occurs when a closure is invoqued
	if !mock.makeRawCalled {
		t.Error("MakeRaw should have been called")
	}

	if mock.restoreCalled {
		t.Error("Restore shouldn't have been called yet")
	}

	restore()
	if !mock.restoreCalled {
		t.Error("Restore should have been called after invoque a closure")
	}
}

func TestEnableRawMode_PropagateError(t *testing.T) {
	mock := &mockController{makeRawErr: errors.New("TTY failure")}

	_, err := enableRawModeWith(mock, 0)
	if err == nil {
		t.Fatal("expected an error when MakeRaw fails")
	}

	if mock.restoreCalled {
		t.Error("Restore shouldn't have been called when MakeRaw fails")
	}
}
