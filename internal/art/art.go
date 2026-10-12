package art

import (
	"embed"
	"fmt"
	"image"
	"math"
	"strings"
)

//go:embed portraits/*.txt
var portraits embed.FS

func Portrait(style string, width int) (string, error) {
	if style != "makima" && style != "reze" && style != "denji" {
		return "", fmt.Errorf("unknown style: use makima, reze, denji")
	}
	switch width {
	case 80, 100, 160, 200, 240:
	default:
		return "", fmt.Errorf("portrait width must be 80, 100, 160, 200 or 240")
	}
	data, err := portraits.ReadFile(fmt.Sprintf("portraits/%s-%d.txt", style, width))
	return string(data), err
}

// Render averages each source region; transparent pixels blend onto black.
func Render(img image.Image, width int) string {
	return RenderAdjusted(img, width, false, 1)
}

func RenderAdjusted(img image.Image, width int, invert bool, contrast float64) string {
	b := img.Bounds()
	height := max(1, min(180, int(math.Round(float64(b.Dy())/float64(b.Dx())*float64(width)*0.5))))
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
			brightness = max(0, min(1, (brightness-0.5)*contrast+0.5))
			if invert {
				brightness = 1 - brightness
			}
			out.WriteByte(palette[min(len(palette)-1, int(brightness*float64(len(palette)-1)))])
		}
		out.WriteByte('\n')
	}
	return out.String()
}
