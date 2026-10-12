package art

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestRenderLocalFile(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 20, 28))
	for y := range 28 {
		for x := range 20 {
			img.SetGray(x, y, color.Gray{Y: uint8(x * 12)})
		}
	}
	path := filepath.Join(t.TempDir(), "portrait.png")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(file, img); err != nil {
		file.Close()
		t.Fatal(err)
	}
	file.Close()
	got, err := RenderFile(path, 160, true, 1.5)
	if err != nil || got != RenderAdjusted(img, 160, true, 1.5) {
		t.Fatalf("local file must use the same image conversion: %v", err)
	}
	if err := os.WriteFile(path, []byte("invalid image"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := RenderFile(path, 160, false, 1); err == nil {
		t.Fatal("invalid image must fail")
	}
}
