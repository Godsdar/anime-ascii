package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/Godsdar/anime-ascii/internal/art"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "convert:", err)
		os.Exit(1)
	}
}

func run() error {
	input := flag.String("image", "", "PNG/JPEG/GIF source file")
	width := flag.Int("width", 100, "output columns: 10..160")
	invert := flag.Bool("invert", false, "render dark source pixels as dense characters")
	contrast := flag.Float64("contrast", 1, "contrast multiplier: 0.1..3")
	crop := flag.String("crop", "", "optional source rectangle: left,top,right,bottom")
	flag.Parse()
	if *input == "" || flag.NArg() != 0 || *width < 10 || *width > 160 || math.IsNaN(*contrast) || *contrast < 0.1 || *contrast > 3 {
		return fmt.Errorf("use -image FILE, -width 10..160 and -contrast 0.1..3")
	}
	file, err := os.Open(*input)
	if err != nil {
		return err
	}
	defer file.Close()
	config, _, err := image.DecodeConfig(file)
	if err != nil {
		return err
	}
	if config.Width <= 0 || config.Height <= 0 || int64(config.Width)*int64(config.Height) > 12_000_000 {
		return fmt.Errorf("image exceeds 12 million pixels")
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}
	if *crop != "" {
		parts := strings.Split(*crop, ",")
		if len(parts) != 4 {
			return fmt.Errorf("crop must be left,top,right,bottom")
		}
		var coords [4]int
		for i, part := range parts {
			coords[i], err = strconv.Atoi(part)
			if err != nil {
				return fmt.Errorf("crop must contain four integers")
			}
		}
		if coords[0] >= coords[2] || coords[1] >= coords[3] {
			return fmt.Errorf("crop must have positive width and height")
		}
		rect := image.Rect(coords[0], coords[1], coords[2], coords[3])
		if !rect.In(img.Bounds()) {
			return fmt.Errorf("crop is outside the source image")
		}
		cropped := image.NewRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
		draw.Draw(cropped, cropped.Bounds(), img, rect.Min, draw.Src)
		img = cropped
	}
	_, err = fmt.Print(art.RenderAdjusted(img, *width, *invert, *contrast))
	return err
}
