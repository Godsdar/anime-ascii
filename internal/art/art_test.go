package art

import (
	"image"
	"image/color"
	"strings"
	"testing"
)

func TestConvertedPortraits(t *testing.T) {
	var previous string
	for _, style := range []string{"makima", "reze"} {
		text, err := Portrait(style)
		if err != nil {
			t.Fatal(err)
		}
		lines := strings.Split(strings.TrimSuffix(text, "\n"), "\n")
		if len(lines) < 50 || len(lines) > 80 {
			t.Fatalf("%s: unexpected portrait height", style)
		}
		for _, line := range lines {
			if len(line) != 100 {
				t.Fatalf("%s: expected 100 columns", style)
			}
			for _, ch := range line {
				if !strings.ContainsRune(" .:-=+*#%@", ch) {
					t.Fatalf("%s: invalid ASCII density character", style)
				}
			}
		}
		if text == previous {
			t.Fatal("portraits must differ")
		}
		previous = text
	}
	if _, err := Portrait("unknown"); err == nil {
		t.Fatal("unknown character must fail")
	}
}

func TestColorize(t *testing.T) {
	text := " /\\\n<^>\n"
	plain, err := Colorize(text, "none", 0, false)
	if err != nil || plain != text {
		t.Fatal("plain output must remain unchanged")
	}
	colored, err := Colorize(text, "rainbow", 0, false)
	if err != nil || !strings.Contains(colored, "\x1b[38;2;") || !strings.HasSuffix(colored, "\x1b[0m") {
		t.Fatal("truecolor output must set and reset ANSI colors")
	}
	next, _ := Colorize(text, "rainbow", 1, false)
	if next == colored {
		t.Fatal("gradient must change between animation frames")
	}
	preview, _ := Colorize("<script>", "pastel", 0, true)
	if strings.Contains(preview, "<script>") || !strings.Contains(preview, "&lt;") {
		t.Fatal("HTML preview must escape text")
	}
	if _, err := Colorize(text, "unknown", 0, false); err == nil {
		t.Fatal("invalid color should fail")
	}
}

func TestTokyoNightPalette(t *testing.T) {
	r, g, b := gradient(0, 0, 0, "tokyonight")
	if r != 122 || g != 162 || b != 247 {
		t.Fatalf("expected Tokyo Night blue, got %d %d %d", r, g, b)
	}
	for frame := range 64 {
		for _, point := range [][2]int{{0, 0}, {30, 10}, {50, 26}} {
			r, g, b := gradient(point[0], point[1], frame, "tokyonight")
			for _, value := range []int{r, g, b} {
				if value < 0 || value > 255 {
					t.Fatalf("invalid RGB component: %d", value)
				}
			}
		}
	}
	colored, err := Colorize("Makima", "tokyonight", 0, false)
	if err != nil || !strings.Contains(colored, "\x1b[38;2;122;162;247m") {
		t.Fatal("Tokyo Night must be available in ANSI output")
	}
}

func TestRenderBlackWhiteAndTransparency(t *testing.T) {
	img := image.NewNRGBA(image.Rect(0, 0, 2, 2))
	for y := 0; y < 2; y++ {
		img.SetNRGBA(0, y, color.NRGBA{R: 255, G: 255, B: 255, A: 0})
		img.SetNRGBA(1, y, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	}
	got := Render(img, 10)
	want := strings.Repeat("     @@@@@\n", 5)
	if got != want {
		t.Fatalf("unexpected rendering: %q", got)
	}
	inverted := RenderAdjusted(img, 10, true, 1.5)
	if inverted != strings.Repeat("@@@@@     \n", 5) {
		t.Fatalf("unexpected inverted rendering: %q", inverted)
	}
}
