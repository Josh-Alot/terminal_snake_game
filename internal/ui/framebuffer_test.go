package ui

import (
	"bytes"
	"strings"
	"testing"
)

func TestFrameBuffer_FlushEmitsSetCells(t *testing.T) {
	fb := NewFrameBuffer(10, 5)
	fb.Set(2, 1, 'X')
	fb.Set(5, 3, 'Y')

	var buf bytes.Buffer
	fb.Flush(&buf)

	got := buf.String()
	for _, want := range []string{"\033[1;2HX", "\033[3;5HY"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected output to contain %q, got %q", want, got)
		}
	}
}

func TestFrameBuffer_NoDiffNoOutput(t *testing.T) {
	fb := NewFrameBuffer(10, 5)
	fb.Set(1, 1, 'A')

	var buf bytes.Buffer
	fb.Flush(&buf)
	buf.Reset()

	fb.Flush(&buf)

	if got := buf.String(); got != "" {
		t.Errorf("expected no output on a diff-free flush, got %q", got)
	}
}

func TestFrameBuffer_FirstFlushDrawsAll(t *testing.T) {
	fb := NewFrameBuffer(10, 5)
	fb.Set(3, 2, 'S')
	fb.Set(7, 4, 'N')

	var buf bytes.Buffer
	fb.Flush(&buf)

	got := buf.String()
	for _, want := range []string{"\033[2;3HS", "\033[4;7HN"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected first flush to contain %q, got %q", want, got)
		}
	}
}

func TestFrameBuffer_ClearProducesSpaces(t *testing.T) {
	fb := NewFrameBuffer(10, 5)
	fb.Set(2, 2, 'X')

	var buf bytes.Buffer
	fb.Flush(&buf)
	buf.Reset()

	fb.Clear()
	fb.Flush(&buf)

	if got := buf.String(); !strings.Contains(got, "\033[2;2H ") {
		t.Errorf("expected flush after Clear to emit a space at the cleared cell, got %q", got)
	}
}

func TestFrameBuffer_OutOfBoundsIgnored(t *testing.T) {
	fb := NewFrameBuffer(10, 5)

	// None of these may panic.
	fb.Set(0, 0, 'X')
	fb.Set(11, 1, 'X')
	fb.Set(1, 6, 'X')
	fb.Set(-3, -2, 'X')

	var buf bytes.Buffer
	fb.Flush(&buf)

	if got := buf.String(); got != "" {
		t.Errorf("expected out-of-bounds writes to be ignored, got %q", got)
	}
}

func TestFrameBuffer_WriteStringMultibyte(t *testing.T) {
	fb := NewFrameBuffer(10, 5)
	fb.WriteString(2, 1, "██")

	var buf bytes.Buffer
	fb.Flush(&buf)

	// The UTF-8 bytes of █ must be emitted contiguously: interleaving
	// cursor moves inside a rune would corrupt the output.
	if got := buf.String(); !strings.Contains(got, "\033[1;2H██") {
		t.Errorf("expected contiguous %q after cursor move, got %q", "██", got)
	}
}

func TestFrameBuffer_DrawBox(t *testing.T) {
	fb := NewFrameBuffer(6, 4)
	fb.DrawBox(Size{Width: 6, Height: 4})

	var buf bytes.Buffer
	fb.Flush(&buf)

	got := buf.String()
	for _, want := range []string{
		"\033[1;1H┌────┐", // top border with both corners
		"\033[4;1H└────┘", // bottom border with both corners
		"\033[2;1H│",      // left side
		"\033[2;6H│",      // right side
	} {
		if !strings.Contains(got, want) {
			t.Errorf("expected box to contain %q, got %q", want, got)
		}
	}
}
