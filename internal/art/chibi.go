package art

import "fmt"

func Chibi(style string) (string, error) {
	switch style {
	case "makima":
		return `       .----.
      / /| | \
     | /_|_|_|
     | | @ @ |
     | |  -  |
     |  '---' |
     () /\_/\ ()
     ()|  V  |()
       |__|__|
`, nil
	case "reze":
		return `       .----.
      / .--. \
     / /    \ \
    | | o  o | |
    | |  _   | |
     \ '----' /
      '--[]--'
       /|  |\
      (_|__|_)
`, nil
	case "denji":
		return `       /\/\/\
      /      \
     |  o  o  |
     |   v    |
      '------'
       /\__/\
      /| >< |\
     (_|____|_)
`, nil
	default:
		return "", fmt.Errorf("unknown style: use makima, reze, denji")
	}
}

func ColorizeChibi(text, theme string, frame int, style string, preview bool) (string, error) {
	return colorize(text, theme, frame, "chibi:"+style, preview)
}

func chibiColor(style string, row int, ch rune) (int, int, int) {
	if ch == '@' || ch == 'o' {
		if style == "reze" {
			return 158, 206, 106
		}
		return 224, 175, 104
	}
	if row < 3 {
		switch style {
		case "makima":
			return 247, 118, 142
		case "denji":
			return 224, 175, 104
		default:
			return 187, 154, 247
		}
	}
	if row < 6 {
		return 192, 202, 245
	}
	return 122, 162, 247
}
