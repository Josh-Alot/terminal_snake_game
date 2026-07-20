package ui

import (
	"fmt"
	"io"
	"strings"
)

// FrameBuffer is a diff-based renderer: the game builds a full frame in
// the buffer and Flush emits only the cells that changed since the
// previous flush, minimising cursor movement and bytes written.
//
// Each terminal cell holds one rune (not one byte) so that multi-byte
// glyphs used by the game (█, *, ─, │, ┌…) map to exactly one cell.
type FrameBuffer struct {
	width, height int
	cur           [][]rune // current frame, one rune per terminal cell
	prev          [][]rune // last flushed frame
}

func NewFrameBuffer(width, height int) *FrameBuffer {
	newPlane := func(fill rune) [][]rune {
		plane := make([][]rune, height)
		for y := range plane {
			plane[y] = make([]rune, width)
			for x := range plane[y] {
				plane[y][x] = fill
			}
		}
		return plane
	}

	return &FrameBuffer{
		width:  width,
		height: height,
		cur:    newPlane(' '),
		prev:   newPlane(' '), // first flush diffs every non-space cell
	}
}

// Set writes ch at (x, y); coordinates are 1-based and out-of-bounds
// writes are ignored.
func (fb *FrameBuffer) Set(x, y int, ch rune) {
	if x < 1 || x > fb.width || y < 1 || y > fb.height {
		return
	}
	fb.cur[y-1][x-1] = ch
}

// WriteString writes s into consecutive cells starting at (x, y), one
// rune per cell.
func (fb *FrameBuffer) WriteString(x, y int, s string) {
	for _, r := range s {
		fb.Set(x, y, r)
		x++
	}
}

// Clear resets the current frame to spaces; prev is left untouched so
// the next flush emits the erased cells.
func (fb *FrameBuffer) Clear() {
	for y := range fb.cur {
		for x := range fb.cur[y] {
			fb.cur[y][x] = ' '
		}
	}
}

// DrawBox draws a border of the given size into the current frame,
// mirroring ui.DrawBox. Boxes smaller than 3x3 are skipped.
func (fb *FrameBuffer) DrawBox(size Size) {
	if size.Width < 3 || size.Height < 3 {
		return
	}

	fb.WriteString(1, 1, "┌"+strings.Repeat("─", size.Width-2)+"┐")
	for y := 2; y < size.Height; y++ {
		fb.WriteString(1, y, "│")
		fb.WriteString(size.Width, y, "│")
	}
	fb.WriteString(1, size.Height, "└"+strings.Repeat("─", size.Width-2)+"┘")
}

// Flush emits the diff between cur and prev in a single write, then
// copies cur into prev. Consecutive changed cells on the same row are
// grouped under one cursor move so multi-byte UTF-8 runes (█, ─, │…)
// are reassembled contiguously.
func (fb *FrameBuffer) Flush(w io.Writer) {
	var out strings.Builder

	for y := 1; y <= fb.height; y++ {
		runStart := 0
		for x := 1; x <= fb.width+1; x++ { // x == width+1 flushes an open run
			diff := x <= fb.width && fb.cur[y-1][x-1] != fb.prev[y-1][x-1]
			if diff && runStart == 0 {
				runStart = x
			}
			if !diff && runStart != 0 {
				fmt.Fprintf(&out, "\033[%d;%dH%s", y, runStart, string(fb.cur[y-1][runStart-1:x-1]))
				runStart = 0
			}
		}
	}

	io.WriteString(w, out.String())

	for y := range fb.cur {
		copy(fb.prev[y], fb.cur[y])
	}
}
