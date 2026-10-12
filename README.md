# anime-ascii

ASCII portraits of Makima and Reze from Chainsaw Man for your Unix terminal, in the Tokyo Night palette. A Go HTTP server supplies the artwork; the CLI adds ANSI truecolor and optional blinking eyes.

![Makima and Reze ASCII portraits in Tokyo Night colors](docs/images/anime-gallery.png)

The previews render actual captured CLI output at 3× resolution. For sharp zooming, open the vector versions: [portrait gallery](docs/images/anime-gallery.svg) and [Makima terminal](docs/images/makima-terminal.svg).

## Quick start

You need **Go 1.26+**. The project uses only the Go standard library.

```sh
git clone https://github.com/Godsdar/anime-ascii.git
cd anime-ascii
go run ./cmd/server
```

Keep the server running. In a second terminal, from the same directory:

```sh
go run ./cmd/art -style makima -animate
```

Press **Ctrl+C** to stop and return to your shell. Use a terminal that supports ANSI escape sequences and 24-bit color.

![Makima ASCII portrait in a high-resolution terminal preview](docs/images/makima-terminal.png)

## Choose a character

| Character | Style | Details |
| --- | --- | --- |
| Makima | `makima` | Braided hair, ringed eyes, shirt and tie. Default. |
| Reze | `reze` | Short hair, choker with a pull ring, sleeveless top. |

```sh
go run ./cmd/art -style makima -color always
go run ./cmd/art -style reze -color always
go run ./cmd/art -style reze -mood sleepy -animate
```

Use `-mood happy`, `-mood sleepy`, or `-mood surprised` to change the expression.

## Colors and pipes

The default theme is `tokyonight`: blue (`#7aa2f7`), purple (`#bb9af7`), cyan (`#7dcfff`), and pink (`#f7768e`). The previews use Tokyo Night background colors as well. The CLI colors the text and keeps your terminal background. `-theme rainbow` and `-theme pastel` remain available.

| Option | Behavior |
| --- | --- |
| `-color auto` | Enables color for a character-device stdout unless `NO_COLOR` is set or `TERM=dumb`. Default. |
| `-color always` | Forces ANSI truecolor output. |
| `-color never` | Prints plain text. |

Redirecting ordinary output to a file or pipe leaves out the color codes:

```sh
go run ./cmd/art -style makima > makima.txt
go run ./cmd/art -style reze | less
```

Use `-animate` in an interactive terminal. It uses screen-control sequences even when colors are disabled.

## Turn your own image into ASCII

The converter accepts PNG, JPEG, and GIF. It maps image brightness to characters and can apply the same color themes to the result.

```sh
go run ./cmd/art -image picture.png -width 100 -color always
go run ./cmd/art -image picture.png -width 100 > picture.txt
```

Uploads are limited to **8 MiB** and **12 million pixels**. Output width is 10 to 160 characters, with a maximum height of 80 lines. Transparent pixels blend onto black; GIF conversion uses the first frame. Very tall images are compressed when they reach the height limit.

## Use curl

Get a colored portrait without the CLI:

```sh
curl -fsS 'http://127.0.0.1:8080/art?style=makima&color=tokyonight'
curl -fsS 'http://127.0.0.1:8080/art?style=reze&mood=happy&color=tokyonight'
```

Convert an image by sending its bytes as the request body:

```sh
curl -fsS --data-binary @picture.png \
  'http://127.0.0.1:8080/convert?width=100' > picture.txt
```

## HTTP API

| Route | Response |
| --- | --- |
| `GET /` | Endpoint reference. |
| `GET /art` | Plain-text Makima portrait by default. |
| `GET /preview?style=makima` | A single colored frame in the browser. |
| `POST /convert?width=80` | ASCII from a raw PNG, JPEG, or GIF request body. |
| `GET /fast` | A response after 200 ms for the concurrency demo. |
| `GET /slow` | A response after one second for the concurrency demo. |

`/art` and `/preview` accept `style`, `mood`, and a nonnegative integer `frame`. The `color` parameter accepts `none`, `tokyonight`, `rainbow`, or `pastel`. `/art` defaults to `none`; `/preview` defaults to `tokyonight`.

Successful `/art` and `/convert` responses use `text/plain; charset=utf-8`. Invalid parameters or images return **400**; oversized uploads or image dimensions return **413**.

The server listens on `127.0.0.1:8080` by default. To use another port:

```sh
go run ./cmd/server -addr 127.0.0.1:8090
go run ./cmd/art -server http://127.0.0.1:8090 -style makima
```

## Build and check

```sh
mkdir -p bin
go build -o bin/ascii-server ./cmd/server
go build -o bin/art ./cmd/art
go test -race ./...
go vet ./...
```

Run `./bin/ascii-server` in one terminal and `./bin/art -style makima -animate` in another. The compiled binaries run without a Go installation.

The CLI is intended for macOS and Linux. It has been tested on macOS.

## Source map

| Path | Purpose |
| --- | --- |
| `cmd/server/` | HTTP server and request handlers. |
| `cmd/art/` | Terminal client, flags, and signal handling. |
| `internal/art/` | Character templates, color gradients, and image conversion. |
| `cmd/check/` | A separate concurrent URL-checking demo. |

The portraits are built from text templates. Image conversion works with files you provide; it does not generate new illustrations.

Characters: Makima and Reze from Chainsaw Man, by Tatsuki Fujimoto.
