package ui

import (
	"errors"
	"testing"
)

type mockSizeRender struct {
	width      int
	height     int
	err        error
	calledWith int
}

func (m *mockSizeRender) GetSize(fd int) (int, int, error) {
	m.calledWith = fd
	return m.width, m.height, m.err
}

func TestGetTerminalSize_ReturnDimensions(t *testing.T) {
	mock := &mockSizeRender{width: 800, height: 600}

	size, err := getTerminalSizeWith(mock, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if size.Width != 800 || size.Height != 600 {
		t.Errorf("getSizeWith: expected %dx%d, got %dx%d", 800, 600, size.Width, size.Height)
	}

	if mock.calledWith != 0 {
		t.Errorf("GetSize called with fd %d, want 0", mock.calledWith)
	}
}

func TestGetTerminalSize_PropagateError(t *testing.T) {
	mock := &mockSizeRender{err: errors.New("not a terminal")}

	_, err := getTerminalSizeWith(mock, 0)
	if err == nil {
		t.Fatal("expected an error when GetSize fails")
	}
}
