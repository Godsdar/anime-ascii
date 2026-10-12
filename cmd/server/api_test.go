package main

import (
	"bytes"
	"image"
	"image/png"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRoutes(t *testing.T) {
	handler := routes()
	for _, tc := range []struct {
		path string
		code int
	}{
		{"/art", 200},
		{"/art?style=reze&width=32", 200},
		{"/preview?style=makima&width=40", 200},
		{"/art?style=denji&width=240", 200},
		{"/art?style=makima&width=160", 200},
		{"/art?mode=image&width=161", 400},
		{"/art?width=oops", 400},
		{"/preview?style=denji&width=240", 200},
		{"/art?style=makima&color=tokyonight", 200},
		{"/preview?style=reze&color=tokyonight", 200},
		{"/art?color=unknown", 400},
		{"/art?style=reze&frame=7", 200},
		{"/art?style=unknown", 400},
		{"/art?frame=-1", 400},
		{"/art?frame=oops", 400},
		{"/missing", 404},
	} {
		r := httptest.NewRequest("GET", tc.path, nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		if w.Code != tc.code {
			t.Errorf("%s: got %d, want %d", tc.path, w.Code, tc.code)
		}
	}
}

func TestConvert(t *testing.T) {
	var data bytes.Buffer
	if err := png.Encode(&data, image.NewGray(image.Rect(0, 0, 20, 20))); err != nil {
		t.Fatal(err)
	}
	handler := routes()
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest("POST", "/convert?width=20", &data))
	if w.Code != 200 || w.Body.String() != strings.Repeat(strings.Repeat(" ", 20)+"\n", 10) {
		t.Fatalf("unexpected conversion: %d %q", w.Code, w.Body.String())
	}
	for _, path := range []string{"/convert", "/convert?width=0", "/convert?width=241", "/convert?width=oops", "/convert?invert=oops", "/convert?contrast=NaN", "/convert?contrast=4"} {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, httptest.NewRequest("POST", path, strings.NewReader("not an image")))
		if w.Code != 400 {
			t.Errorf("%s: got %d, want 400", path, w.Code)
		}
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest("POST", "/convert", strings.NewReader(strings.Repeat("x", (8<<20)+1))))
	if w.Code != 413 {
		t.Fatalf("oversized upload: got %d, want 413", w.Code)
	}
}

func TestDefaultPortraitStaysTiny(t *testing.T) {
	for _, style := range []string{"makima", "reze", "denji"} {
		for _, suffix := range []string{"", "&width=240"} {
			w := httptest.NewRecorder()
			routes().ServeHTTP(w, httptest.NewRequest("GET", "/art?style="+style+suffix, nil))
			if w.Code != 200 {
				t.Fatal(w.Body.String())
			}
			lines := strings.Split(strings.TrimSuffix(w.Body.String(), "\n"), "\n")
			if len(lines) > 10 {
				t.Fatal("default art must not become a large portrait")
			}
			for _, line := range lines {
				if len(line) > 20 {
					t.Fatal("default art must fit a small terminal")
				}
			}
		}
	}
}
