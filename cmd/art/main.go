package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Godsdar/anime-ascii/internal/art"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "art:", err)
		os.Exit(1)
	}
}

func run() error {
	server := flag.String("server", "http://127.0.0.1:8080", "ASCII server URL")
	style := flag.String("style", "makima", "makima, reze")
	colorMode := flag.String("color", "auto", "auto, always, never")
	theme := flag.String("theme", "tokyonight", "tokyonight, rainbow, pastel")
	imagePath := flag.String("image", "", "PNG/JPEG/GIF file to convert")
	width := flag.Int("width", 80, "image output width: 10..160")
	invert := flag.Bool("invert", false, "dense characters for dark image pixels")
	contrast := flag.Float64("contrast", 1, "image contrast multiplier: 0.1..3")
	animate := flag.Bool("animate", false, "animate the color gradient; Ctrl+C to stop")
	flag.Parse()
	if *colorMode != "auto" && *colorMode != "always" && *colorMode != "never" {
		return fmt.Errorf("color must be auto, always or never")
	}
	if *theme != "tokyonight" && *theme != "rainbow" && *theme != "pastel" {
		return fmt.Errorf("theme must be tokyonight, rainbow or pastel")
	}
	useColor := *colorMode == "always"
	if *colorMode == "auto" {
		info, err := os.Stdout.Stat()
		_, noColor := os.LookupEnv("NO_COLOR")
		useColor = err == nil && info.Mode()&os.ModeCharDevice != 0 && !noColor && os.Getenv("TERM") != "dumb"
	}
	if flag.NArg() != 0 {
		return fmt.Errorf("unexpected arguments; use -help for options")
	}
	if *width < 10 || *width > 160 {
		return fmt.Errorf("width must be between 10 and 160")
	}
	if math.IsNaN(*contrast) || *contrast < 0.1 || *contrast > 3 {
		return fmt.Errorf("contrast must be between 0.1 and 3")
	}
	if *animate && *imagePath != "" {
		return fmt.Errorf("animation supports portraits only")
	}
	endpoint, err := url.Parse(*server)
	if err != nil || (endpoint.Scheme != "http" && endpoint.Scheme != "https") || endpoint.Host == "" {
		return fmt.Errorf("server must be an http:// or https:// URL")
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if *animate {
		// Alternate screen preserves shell contents; every exit restores the cursor.
		fmt.Print("\x1b[?1049h\x1b[?25l")
		defer fmt.Print("\x1b[0m\x1b[?25h\x1b[?1049l")
	}
	client := &http.Client{Timeout: 10 * time.Second}
	for frame := 0; ; frame = (frame + 1) % 64 {
		method := http.MethodGet
		var body io.Reader
		var file *os.File
		q := url.Values{}
		if *imagePath != "" {
			file, err = os.Open(*imagePath)
			if err != nil {
				return err
			}
			body = file
			method = http.MethodPost
			endpoint.Path = "/convert"
			q.Set("width", strconv.Itoa(*width))
			q.Set("invert", strconv.FormatBool(*invert))
			q.Set("contrast", strconv.FormatFloat(*contrast, 'f', -1, 64))
		} else {
			endpoint.Path = "/art"
			q.Set("style", *style)
			q.Set("frame", strconv.Itoa(frame))
		}
		endpoint.RawQuery = q.Encode()
		req, err := http.NewRequestWithContext(ctx, method, endpoint.String(), body)
		if err != nil {
			if file != nil {
				file.Close()
			}
			return err
		}
		if file != nil {
			req.Header.Set("Content-Type", "application/octet-stream")
		}
		resp, err := client.Do(req)
		if file != nil {
			file.Close()
		}
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return err
		}
		data, err := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
		resp.Body.Close()
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("server returned %s: %s", resp.Status, data)
		}
		if useColor {
			colored, err := art.Colorize(string(data), *theme, frame, false)
			if err != nil {
				return err
			}
			data = []byte(colored)
		}
		if *animate {
			fmt.Print("\x1b[H")
		}
		if _, err := os.Stdout.Write(data); err != nil {
			return err
		}
		if !*animate {
			return nil
		}
		timer := time.NewTimer(200 * time.Millisecond)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil
		case <-timer.C:
		}
	}
}
