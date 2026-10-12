# anime-ascii

Makima, Reze and Denji illustrations converted into ASCII for your Unix terminal, in the Tokyo Night palette. The standalone Go command includes every portrait and adds ANSI truecolor. It runs without a server; an optional HTTP mode is also available.

![Makima, Reze and Denji ASCII portraits in Tokyo Night colors](docs/images/anime-gallery.png)

The portraits were converted from illustrations by [eastman](https://wall.alphacoders.com/big.php?i=1284789), [SpringHaru](https://twiman.net/user/1348570965135949830/illust/1970415339994521708), and a [Denji face study by masoq](https://twitter.com/masoq095/status/1645434593502437376). The previews render actual captured CLI output at 3× resolution. The gallery shows compact 40-column portraits; vector detail previews use 240 columns. For sharp zooming, open the vector versions: [portrait gallery](docs/images/anime-gallery.svg), [Makima terminal](docs/images/makima-terminal.svg), [Reze terminal](docs/images/reze-terminal.svg), and [Denji terminal](docs/images/denji-terminal.svg).

## Install and run

You need **Go 1.26+** to install from source. The project uses only the Go standard library.

```sh
GOBIN="$HOME/.local/bin" go install github.com/Godsdar/anime-ascii/cmd/anime-ascii@latest
export PATH="$HOME/.local/bin:$PATH"

anime-ascii makima
anime-ascii reze
```

The installed binary runs from any directory without Go or an HTTP server. It includes Makima, Reze and Denji at every supported resolution. Put `~/.local/bin` on your shell's `PATH` to keep the command available in new terminals.

Use a terminal that supports ANSI escape sequences and 24-bit color. By default, the CLI chooses a compact size from 32 to 48 columns based on the terminal width and height. It leaves room for the prompt. In a standard 80×24 terminal, Makima and Reze use 40 columns and 20 rows. Use `-width 40` to choose a compact size explicitly or `-width 240` for the most detail. `less -RS` preserves colors and lets you scroll vertically and horizontally.

[Open the optional 240-column Makima detail](docs/images/makima-terminal.svg).

## Choose a character

| Character | Style | Details |
| --- | --- | --- |
| Makima | `makima` | Converted portrait by eastman. Default. |
| Reze | `reze` | Converted illustration by SpringHaru. |
| Denji | `denji` | Human portrait with both eyes open, converted from a face study by masoq. |

```sh
anime-ascii makima -color always
anime-ascii reze -color always
anime-ascii reze -width 240 | less
anime-ascii denji -width 240 -color always
```

The facial expressions come from the source images. The portraits are static. Optional `-animate` moves the color gradient; it requires a terminal tall enough to hold the entire portrait. Ctrl+C returns to your shell.

## Readable eyes

Every built-in portrait is a close face crop with both eyes visible. Thin dark contours survive downsampling; the source pupil centers remain distinct `@` characters at every supported width. Eye regions use Tokyo Night gold or green accents so the pupils stand apart from the surrounding gradient.

The smallest supported portrait is 32 columns. Automatic sizing keeps portraits compact without shrinking below that minimum. Denji uses his human form so his eyes remain visible.

## Resolution

The default `-width 0` automatically selects a small portrait. Explicit widths are `32`, `40`, `48`, `64`, `80`, `100`, `160`, `200`, or `240`. Each version is converted directly from the source image. The height grows with the width to preserve the portrait proportions.

| Character | Compact (40 columns) | Detailed (240 columns) |
| --- | --- | --- |
| Makima | 40 × 20 | 240 × 120 |
| Reze | 40 × 20 | 240 × 120 |
| Denji | 40 × 21 | 240 × 123 |

```sh
anime-ascii makima -width 240 -color always
anime-ascii denji -width 240 -color always | less -R
```

[Open the optional 240-column Denji detail](docs/images/denji-terminal.svg).

## Colors and pipes

The default theme is `tokyonight`: blue (`#7aa2f7`), purple (`#bb9af7`), cyan (`#7dcfff`), and pink (`#f7768e`). The previews use Tokyo Night background colors as well. The CLI colors the text and keeps your terminal background. `-theme rainbow` and `-theme pastel` remain available.

| Option | Behavior |
| --- | --- |
| `-color auto` | Enables color for a character-device stdout unless `NO_COLOR` is set or `TERM=dumb`. Default. |
| `-color always` | Forces ANSI truecolor output. |
| `-color never` | Prints plain text. |

Redirecting ordinary output to a file or pipe leaves out the color codes:

```sh
anime-ascii makima > makima.txt
anime-ascii reze | less
```

Use `-animate` in an interactive terminal. It uses screen-control sequences even when colors are disabled.

## Turn your own image into ASCII

The converter accepts PNG, JPEG, and GIF. It maps image brightness to characters and can apply the same color themes to the result. `-invert` makes dark pixels produce dense characters, which works well for illustrations on a white background. `-contrast` accepts values from 0.1 to 3 (default 1).

```sh
anime-ascii -image picture.png -width 100 -invert -contrast 1.5 -color always
anime-ascii -image picture.png -width 100 > picture.txt
```

Uploads are limited to **8 MiB** and **12 million pixels**. Output width is 10 to 240 characters, with a maximum height of 180 lines. Transparent pixels blend onto black; GIF conversion uses the first frame. Very tall images are compressed when they reach the height limit.

To convert locally without a server:

```sh
go run ./cmd/convert -image picture.jpg -width 100 -invert -contrast 1.5 > picture.txt
```

Add `-crop left,top,right,bottom` to select a region in source pixels. The bundled portraits use this same Go converter; exact settings and source credits are in [artwork sources](internal/art/portraits/SOURCES.md).

## Optional HTTP server

The CLI renders locally by default. Start the server only when you want an HTTP API or the browser preview:

```sh
git clone https://github.com/Godsdar/anime-ascii.git
cd anime-ascii
go run ./cmd/server
```

In another terminal:

```sh
anime-ascii reze -server http://127.0.0.1:8080 -width 240
```

`-server` explicitly switches the CLI to HTTP mode. Built-in portraits and local image conversion work without it.

## Use curl

Get a colored portrait without the CLI:

```sh
curl -fsS 'http://127.0.0.1:8080/art?style=makima&color=tokyonight'
curl -fsS 'http://127.0.0.1:8080/art?style=reze&color=tokyonight'
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

`/art` and `/preview` accept `style`, portrait `width` (32/40/48/64/80/100/160/200/240; default 40), and a nonnegative integer `frame` for the color gradient. The `color` parameter accepts `none`, `tokyonight`, `rainbow`, or `pastel`. `/art` defaults to `none`; `/preview` defaults to `tokyonight`.

Successful `/art` and `/convert` responses use `text/plain; charset=utf-8`. Invalid parameters or images return **400**; oversized uploads or image dimensions return **413**.

`/convert` also accepts `invert=true` and `contrast=1.5`.

The server listens on `127.0.0.1:8080` by default. To use another port:

```sh
go run ./cmd/server -addr 127.0.0.1:8090
anime-ascii -server http://127.0.0.1:8090 -style makima
```

## Build and check

```sh
mkdir -p bin
go build -o bin/ascii-server ./cmd/server
go build -o bin/anime-ascii ./cmd/anime-ascii
go test -race ./...
go vet ./...
```

Run `./bin/anime-ascii makima` directly. Use `./bin/ascii-server` only for the optional HTTP API. The compiled binaries run without a Go installation.

The CLI is intended for macOS and Linux. It has been tested on macOS.

## Source map

| Path | Purpose |
| --- | --- |
| `cmd/server/` | HTTP server and request handlers. |
| `cmd/anime-ascii/` | Standalone command, optional HTTP client, flags, and signal handling. |
| `internal/art/` | Embedded converted portraits, color gradients, and image conversion. |
| `cmd/convert/` | Local image-to-ASCII conversion with crop, inversion, and contrast controls. |
| `cmd/check/` | A separate concurrent URL-checking demo. |

The bundled portraits are image conversions embedded as plain ASCII text. Source raster files are not included. Image conversion also works with files you provide; it does not generate new illustrations.

Characters: Makima, Reze and Denji from Chainsaw Man, by Tatsuki Fujimoto.
