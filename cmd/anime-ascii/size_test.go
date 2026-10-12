package main

import (
	"strings"
	"testing"

	"github.com/Godsdar/anime-ascii/internal/art"
)

func TestCompactPortraitFitsTerminal(t *testing.T) {
	for _, tc := range []struct {
		style            string
		cols, rows, want int
	}{
		{"makima", 80, 24, 40},
		{"reze", 80, 24, 40},
		{"makima", 38, 20, 32},
		{"reze", 80, 30, 48},
		{"makima", 200, 70, 48},
		{"denji", 80, 24, 40},
	} {
		width, err := compactWidth(tc.style, tc.cols, tc.rows)
		if err != nil || width != tc.want {
			t.Fatalf("%s in %dx%d: got %d, want %d (%v)", tc.style, tc.cols, tc.rows, width, tc.want, err)
		}
		text, _ := art.Portrait(tc.style, width)
		if width > tc.cols-2 || strings.Count(text, "\n") > tc.rows-2 {
			t.Fatal("automatic portrait must leave room for the prompt")
		}
	}
}
