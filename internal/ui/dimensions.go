package ui

import (
	"os"

	"golang.org/x/term"
)

type Size struct {
	Height int
	Width  int
}

type sizeRender interface {
	GetSize(fd int) (width, height int, err error)
}

type realSizeRender struct{}

func (realSizeRender) GetSize(fd int) (int, int, error) {
	return term.GetSize(fd)
}

func getTerminalSizeWith(sr sizeRender, fd int) (Size, error) {
	width, height, err := sr.GetSize(fd)
	if err != nil {
		return Size{}, err
	}

	return Size{Height: height, Width: width}, nil
}

func GetTerminalSize() (Size, error) {
	return getTerminalSizeWith(realSizeRender{}, int(os.Stdout.Fd()))
}
