package art

import (
	"image"
	"image/color"
	"strings"
	"testing"
)

func TestPortraitsAndBlink(t *testing.T) {
	for _, style := range []string{"makima", "reze"} {
		awake, err := Portrait(style, "happy", 0)
		if err != nil {
			t.Fatal(err)
		}
		blink, err := Portrait(style, "happy", 7)
		if err != nil || awake == blink || len(awake) != len(blink) {
			t.Fatalf("%s: blink should change eyes without changing frame size", style)
		}
		if strings.Contains(awake, "{{") {
			t.Fatal("unresolved template placeholder")
		}
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
}
