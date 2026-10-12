package art

import (
	"fmt"
	"html"
	"math"
	"strings"
)

// Colorize uses the same smooth gradient for ANSI terminals and HTML previews.
func Colorize(text, theme string, frame int, preview bool) (string, error) {
	if theme != "rainbow" && theme != "pastel" && theme != "none" {
		return "", fmt.Errorf("unknown color: use rainbow, pastel, none")
	}
	if theme == "none" {
		if preview {
			return html.EscapeString(text), nil
		}
		return text, nil
	}
	var out strings.Builder
	for y, line := range strings.SplitAfter(text, "\n") {
		for x, ch := range line {
			if ch == ' ' || ch == '\n' {
				out.WriteRune(ch)
				continue
			}
			r, g, b := gradient(x, y, frame, theme)
			if preview {
				fmt.Fprintf(&out, `<span style="color:rgb(%d,%d,%d)">%s</span>`, r, g, b, html.EscapeString(string(ch)))
			} else {
				fmt.Fprintf(&out, "\x1b[38;2;%d;%d;%dm%c", r, g, b, ch)
			}
		}
		if !preview && strings.HasSuffix(line, "\n") {
			// Reset after each line so piped text never colors subsequent output.
			out.WriteString("\x1b[0m")
		}
	}
	if !preview {
		out.WriteString("\x1b[0m")
	}
	return out.String(), nil
}

func gradient(x, y, frame int, theme string) (int, int, int) {
	phase := float64(x)*0.17 + float64(y)*0.24 + float64(frame)*2*math.Pi/64
	mid, amplitude := 170.0, 80.0
	if theme == "pastel" {
		mid, amplitude = 205, 45
	}
	channel := func(offset float64) int {
		return int(math.Round(mid + amplitude*math.Sin(phase+offset)))
	}
	return channel(0), channel(2 * math.Pi / 3), channel(4 * math.Pi / 3)
}
