package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Godsdar/anime-ascii/internal/art"
)

func terminalSize() (columns, rows int) {
	// stty is available on macOS and Linux; redirected input uses a small fallback.
	columns, rows = 80, 30
	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()
	cmd := exec.CommandContext(ctx, "stty", "size")
	cmd.Stdin = os.Stdin
	output, err := cmd.Output()
	if err != nil {
		return
	}
	var actualRows, actualColumns int
	if _, err := fmt.Sscanf(string(output), "%d %d", &actualRows, &actualColumns); err == nil && actualRows > 0 && actualColumns > 0 {
		columns, rows = actualColumns, actualRows
	}
	return
}

func compactWidth(style string, columns, rows int) (int, error) {
	for _, width := range []int{48, 40, 32} {
		text, err := art.Portrait(style, width)
		if err != nil {
			return 0, err
		}
		if width <= columns-2 && strings.Count(text, "\n") <= rows-2 {
			return width, nil
		}
	}
	return 0, fmt.Errorf("terminal is too small for readable eyes; widen it to at least 34 columns and 19 rows")
}
