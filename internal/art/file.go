package art

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
)

func RenderFile(path string, width int, invert bool, contrast float64) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, (8<<20)+1))
	if err != nil {
		return "", err
	}
	if len(data) > 8<<20 {
		return "", fmt.Errorf("image exceeds 8 MiB")
	}
	config, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	if config.Width <= 0 || config.Height <= 0 || int64(config.Width)*int64(config.Height) > 12_000_000 {
		return "", fmt.Errorf("image exceeds 12 million pixels")
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	return RenderAdjusted(img, width, invert, contrast), nil
}
