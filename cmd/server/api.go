package main

import (
	"bytes"
	"fmt"
	"html/template"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/Godsdar/anime-ascii/internal/art"
)

func routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		plain(w)
		fmt.Fprint(w, "ANIME ASCII\n\nGET /art?style=makima&frame=0\nStyles: makima, reze, denji\nPOST /convert?width=80 (raw PNG/JPEG/GIF body)\nGET /fast and /slow: concurrency demo\n")
	})
	mux.HandleFunc("GET /art", portrait)
	mux.HandleFunc("GET /preview", portrait)
	mux.HandleFunc("POST /convert", convert)
	mux.HandleFunc("GET /fast", delayedReply(200*time.Millisecond))
	mux.HandleFunc("GET /slow", delayedReply(time.Second))
	return mux
}

func plain(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
}

func portrait(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	style := q.Get("style")
	if style == "" {
		style = "makima"
	}
	frame := 0
	if q.Has("frame") {
		var err error
		frame, err = strconv.Atoi(q.Get("frame"))
		if err != nil || frame < 0 {
			http.Error(w, "frame must be a nonnegative integer", http.StatusBadRequest)
			return
		}
	}
	width := 40
	if q.Has("width") {
		var err error
		width, err = strconv.Atoi(q.Get("width"))
		if err != nil {
			http.Error(w, "width must be 32, 40, 48, 64, 80, 100, 160, 200 or 240", http.StatusBadRequest)
			return
		}
	}
	mode := q.Get("mode")
	if mode == "" {
		mode = "chibi"
	}
	var output string
	var err error
	if mode == "chibi" {
		output, err = art.Chibi(style)
	} else if mode == "image" {
		output, err = art.Portrait(style, width)
	} else {
		err = fmt.Errorf("mode must be chibi or image")
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	color := q.Get("color")
	preview := r.URL.Path == "/preview"
	if color == "" {
		color = "none"
		if preview {
			color = "tokyonight"
		}
	}
	if mode == "chibi" {
		output, err = art.ColorizeChibi(output, color, frame, style, preview)
	} else {
		output, err = art.ColorizePortrait(output, color, frame, style, preview)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if preview {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// Only escaped characters and generated color spans reach this template.
		page := template.Must(template.New("preview").Parse(`<!doctype html><html lang="en"><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Anime ASCII</title><style>body{margin:0;min-height:100vh;display:grid;place-items:center;background:#1a1b26;color:#e2e5f0;font-family:monospace}main{max-width:calc(100vw - 100px);overflow:auto;padding:32px;border:1px solid #414868;border-radius:20px;background:#24283b;box-shadow:0 20px 80px #0005}p{color:#a9b1d6}pre{font-size:18px;line-height:1.25;font-family:Menlo,Consolas,monospace}a{color:#b4befe}nav{display:flex;gap:20px}</style><main><p>ANIME ASCII / TOKYO NIGHT</p><pre>{{.Art}}</pre><nav><a href="/preview?style=makima">Makima</a><a href="/preview?style=reze">Reze</a><a href="/preview?style=denji">Denji</a><a href="?style={{.Style}}&mode=image&width=80">image mode</a></nav><p>anime-ascii {{.Style}}</p></main></html>`))
		page.Execute(w, struct {
			Art   template.HTML
			Style string
		}{template.HTML(output), style})
		return
	}
	plain(w)
	fmt.Fprint(w, output)
}

func convert(w http.ResponseWriter, r *http.Request) {
	invert := false
	contrast := 1.0
	q := r.URL.Query()
	if q.Has("invert") {
		var err error
		invert, err = strconv.ParseBool(q.Get("invert"))
		if err != nil {
			http.Error(w, "invert must be true or false", http.StatusBadRequest)
			return
		}
	}
	if q.Has("contrast") {
		var err error
		contrast, err = strconv.ParseFloat(q.Get("contrast"), 64)
		if err != nil || math.IsNaN(contrast) || contrast < 0.1 || contrast > 3 {
			http.Error(w, "contrast must be between 0.1 and 3", http.StatusBadRequest)
			return
		}
	}
	width := 160
	if r.URL.Query().Has("width") {
		var err error
		width, err = strconv.Atoi(r.URL.Query().Get("width"))
		if err != nil || width < 10 || width > 240 {
			http.Error(w, "width must be an integer between 10 and 240", http.StatusBadRequest)
			return
		}
	}
	r.Body = http.MaxBytesReader(w, r.Body, 8<<20)
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read image (maximum 8 MiB)", http.StatusRequestEntityTooLarge)
		return
	}
	config, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		http.Error(w, "send a valid PNG, JPEG or GIF as the request body", http.StatusBadRequest)
		return
	}
	if config.Width <= 0 || config.Height <= 0 || int64(config.Width)*int64(config.Height) > 12_000_000 {
		http.Error(w, "image must contain between 1 and 12 million pixels", http.StatusRequestEntityTooLarge)
		return
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		http.Error(w, "could not decode image", http.StatusBadRequest)
		return
	}
	plain(w)
	fmt.Fprint(w, art.RenderAdjusted(img, width, invert, contrast))
}
