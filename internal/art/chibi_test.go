package art

import (
	"strings"
	"testing"
)

func TestChibiSizeAndEyes(t *testing.T) {
	for _, style := range []string{"makima", "reze", "denji"} {
		text, err := Chibi(style)
		if err != nil {
			t.Fatal(err)
		}
		lines := strings.Split(strings.TrimSuffix(text, "\n"), "\n")
		if len(lines) > 10 {
			t.Fatalf("%s must fit in 10 rows", style)
		}
		for _, line := range lines {
			if len(line) > 20 {
				t.Fatalf("%s must fit in 20 columns", style)
			}
		}
		eye := "o"
		if style == "makima" {
			eye = "@"
		}
		if strings.Count(text, eye) != 2 {
			t.Fatalf("%s needs two distinct eyes", style)
		}
		colored, err := ColorizeChibi(text, "tokyonight", 0, style, false)
		if err != nil || !strings.Contains(colored, "\x1b[38;2;") {
			t.Fatalf("%s needs terminal color", style)
		}
	}
}
