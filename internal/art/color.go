package art

import (
	"fmt"
	"html"
	"math"
	"strings"
)

// Colorize uses the same smooth gradient for ANSI terminals and HTML previews.
func Colorize(text, theme string, frame int, preview bool) (string, error) {
	return colorize(text, theme, frame, "", preview)
}

func ColorizePortrait(text, theme string, frame int, style string, preview bool) (string, error) {
	return colorize(text, theme, frame, style, preview)
}

func colorize(text, theme string, frame int, style string, preview bool) (string, error) {
	if theme != "tokyonight" && theme != "rainbow" && theme != "pastel" && theme != "none" {
		return "", fmt.Errorf("unknown color: use tokyonight, rainbow, pastel, none")
	}
	if theme == "none" {
		if preview {
			return html.EscapeString(text), nil
		}
		return text, nil
	}
	var out strings.Builder
	width := strings.IndexByte(text, '\n')
	if width < 0 {
		width = len(text)
	}
	height := max(1, strings.Count(text, "\n"))
	for y, line := range strings.SplitAfter(text, "\n") {
		for x, ch := range line {
			if ch == ' ' || ch == '\n' {
				out.WriteRune(ch)
				continue
			}
			r, g, b := gradient(x, y, frame, theme)
			if eyePixel(style, x, y, width, height) {
				r, g, b = 224, 175, 104
				if style == "reze" {
					r, g, b = 158, 206, 106
				}
			}
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
	if theme == "tokyonight" {
		colors := [][3]int{
			{122, 162, 247}, {187, 154, 247}, {247, 118, 142},
			{187, 154, 247}, {125, 207, 255}, {122, 162, 247},
		}
		position := float64(x)*0.025 + float64(y)*0.055 + float64(frame)/64*float64(len(colors))
		position = math.Mod(position, float64(len(colors)))
		index := int(position)
		fraction := position - float64(index)
		// Smoothstep gives adjacent palette segments a gentle transition.
		fraction = fraction * fraction * (3 - 2*fraction)
		left, right := colors[index], colors[(index+1)%len(colors)]
		mix := func(channel int) int {
			return int(math.Round(float64(left[channel])*(1-fraction) + float64(right[channel])*fraction))
		}
		return mix(0), mix(1), mix(2)
	}
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

// Eye rectangles are measured on the credited source crops, not drawn into the art.
func eyePixel(style string, x, y, width, height int) bool {
	var regions [][4]int
	var sourceWidth, sourceHeight int
	switch style {
	case "makima":
		sourceWidth, sourceHeight = 620, 620
		regions = [][4]int{{80, 175, 260, 280}, {350, 175, 540, 280}}
	case "reze":
		sourceWidth, sourceHeight = 515, 515
		regions = [][4]int{{35, 160, 160, 250}, {235, 80, 385, 180}}
	case "denji":
		sourceWidth, sourceHeight = 355, 365
		regions = [][4]int{{22, 117, 147, 200}, {163, 87, 284, 164}}
	default:
		return false
	}
	for _, region := range regions {
		left := region[0] * width / sourceWidth
		right := (region[2]*width + sourceWidth - 1) / sourceWidth
		top := region[1] * height / sourceHeight
		bottom := (region[3]*height + sourceHeight - 1) / sourceHeight
		if x >= left && x < right && y >= top && y < bottom {
			return true
		}
	}
	return false
}
