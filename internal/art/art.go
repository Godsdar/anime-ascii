package art

import (
	"fmt"
	"image"
	"math"
	"strings"
)

// Portrait keeps every generated frame the same size for terminal animation.
func Portrait(style, mood string, frame int) (string, error) {
	eyes, mouth := "o", "w"
	switch mood {
	case "happy":
		eyes = "^"
	case "sleepy":
		eyes, mouth = "-", "_"
	case "surprised":
		eyes, mouth = "O", "o"
	default:
		return "", fmt.Errorf("unknown mood: use happy, sleepy, surprised")
	}
	if frame%8 == 7 {
		eyes = "-"
	}
	var template string
	switch style {
	case "neko", "cat":
		template = `       .                    +
            /\_/\
           ( {{E}}.{{E}} )
           / >{{M}}< \
          (  / \  )
           \(___)/   ~
        .   "   "       *
`
	case "chibi", "bunny":
		template = `       +                    .
             /\ /\
            (  V  )
           ( {{E}}.{{E}} )
           /  {{M}}  \
          (  / \  )
           (_/ \_)     *
`
	case "mecha":
		template = `        .                 +
            \  |  /
             [===]
            [{{E}}-{{E}}]
          /| [{{M}}] |\
         (_|=====|_)
            |_ _|
            |_|_|     *
`
	case "fox":
		template = `       .                    *
            /\   /\
           /  \_/  \
           \ {{E}}.{{E}} /
            \ {{M}} /
          / (   ) \___
         (_/(___)\___~)
        +   "   "       .
`
	default:
		return "", fmt.Errorf("unknown style: use neko, bunny, fox, mecha (chibi is an alias for bunny)")
	}
	return strings.NewReplacer("{{E}}", eyes, "{{M}}", mouth).Replace(template), nil
}

// Render averages each source region; transparent pixels blend onto black.
func Render(img image.Image, width int) string {
	b := img.Bounds()
	height := max(1, min(80, int(math.Round(float64(b.Dy())/float64(b.Dx())*float64(width)*0.5))))
	const palette = " .:-=+*#%@"
	var out strings.Builder
	for y := 0; y < height; y++ {
		y0 := b.Min.Y + y*b.Dy()/height
		y1 := min(b.Max.Y, max(y0+1, b.Min.Y+(y+1)*b.Dy()/height))
		for x := 0; x < width; x++ {
			x0 := b.Min.X + x*b.Dx()/width
			x1 := min(b.Max.X, max(x0+1, b.Min.X+(x+1)*b.Dx()/width))
			var sum float64
			for sy := y0; sy < y1; sy++ {
				for sx := x0; sx < x1; sx++ {
					r, g, blue, _ := img.At(sx, sy).RGBA()
					sum += 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(blue)
				}
			}
			brightness := sum / float64((x1-x0)*(y1-y0)) / 65535
			out.WriteByte(palette[min(len(palette)-1, int(brightness*float64(len(palette)-1)))])
		}
		out.WriteByte('\n')
	}
	return out.String()
}
