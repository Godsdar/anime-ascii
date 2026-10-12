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
	"net/http"
	"strconv"
	"time"

	"github.com/Godsdar/anime-ascii/internal/art"
)

func routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		plain(w)
		fmt.Fprint(w, "ANIME ASCII\n\nGET /art?style=makima&mood=happy&frame=0\nStyles: makima, reze\nMoods: happy, sleepy, surprised\nPOST /convert?width=80 (raw PNG/JPEG/GIF body)\nGET /fast and /slow: concurrency demo\n")
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
	style, mood := q.Get("style"), q.Get("mood")
	if style == "" {
		style = "makima"
	}
	if mood == "" {
		mood = "happy"
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
	output, err := art.Portrait(style, mood, frame)
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
	output, err = art.Colorize(output, color, frame, preview)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if preview {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// Only escaped characters and generated color spans reach this template.
		page := template.Must(template.New("preview").Parse(`<!doctype html><html lang="en"><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Anime ASCII</title><style>body{margin:0;min-height:100vh;display:grid;place-items:center;background:#1a1b26;color:#e2e5f0;font-family:monospace}main{padding:32px;border:1px solid #414868;border-radius:20px;background:#24283b;box-shadow:0 20px 80px #0005}p{color:#a9b1d6}pre{font-size:clamp(10px,1.8vw,16px);line-height:1.25;font-family:Menlo,Consolas,monospace}a{color:#b4befe}nav{display:flex;gap:20px}</style><main><p>ANIME ASCII / TOKYO NIGHT</p><pre>{{.}}</pre><nav><a href="/preview?style=makima">Makima</a><a href="/preview?style=reze">Reze</a></nav><p>go run ./cmd/art -animate</p></main></html>`))
		page.Execute(w, template.HTML(output))
		return
	}
	plain(w)
	fmt.Fprint(w, output)
}

func convert(w http.ResponseWriter, r *http.Request) {
	width := 80
	if r.URL.Query().Has("width") {
		var err error
		width, err = strconv.Atoi(r.URL.Query().Get("width"))
		if err != nil || width < 10 || width > 160 {
			http.Error(w, "width must be an integer between 10 and 160", http.StatusBadRequest)
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
	fmt.Fprint(w, art.Render(img, width))
}
