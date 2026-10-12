package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Result struct {
	URL      string
	Status   int
	Duration time.Duration
	Err      error
}

func check(client *http.Client, url string) Result {
	start := time.Now()
	resp, err := client.Get(url)
	if err != nil {
		return Result{URL: url, Duration: time.Since(start), Err: err}
	}
	defer resp.Body.Close()
	_, err = io.Copy(io.Discard, resp.Body)
	return Result{URL: url, Status: resp.StatusCode, Duration: time.Since(start), Err: err}
}

func main() {
	timeout := flag.Duration("timeout", 3*time.Second, "таймаут каждой проверки")
	flag.Parse()
	if *timeout <= 0 {
		fmt.Fprintln(os.Stderr, "таймаут должен быть больше нуля")
		os.Exit(2)
	}
	urls := flag.Args()
	if len(urls) == 0 {
		urls = []string{"http://127.0.0.1:8080/slow", "http://127.0.0.1:8080/fast"}
	}
	client := &http.Client{Timeout: *timeout}
	results := make(chan Result)
	start := time.Now()

	for _, url := range urls {
		go func() {
			results <- check(client, url)
		}()
	}

	failed := false
	for range urls {
		result := <-results
		if result.Err != nil {
			fmt.Printf("ОШИБКА %s: %v\n", result.URL, result.Err)
			failed = true
			continue
		}
		fmt.Printf("%d %s (%s)\n", result.Status, result.URL, result.Duration.Round(time.Millisecond))
		if result.Status >= 400 {
			failed = true
		}
	}
	fmt.Printf("Всего: %s\n", time.Since(start).Round(time.Millisecond))
	if failed {
		os.Exit(1)
	}
}
