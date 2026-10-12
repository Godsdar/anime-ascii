# anime-ascii

Small ASCII chibi portraits of Makima, Reze and Denji for your terminal. Each character fits in **16 columns and 9 rows**, with clear eyes and a restrained Tokyo Night palette.

![Small Tokyo Night chibi portraits](docs/images/anime-gallery.png)

The default artwork is drawn directly as ASCII. Hair, eyes and clothing use separate colors. The gallery renders captured CLI output at 3× resolution; [open the vector preview](docs/images/anime-gallery.svg) for sharp zooming.

## Install

You need Go 1.26+ to install from source. The application uses the Go standard library.

```sh
GOBIN="$HOME/.local/bin" go install github.com/Godsdar/anime-ascii/cmd/anime-ascii@latest
export PATH="$HOME/.local/bin:$PATH"
```

The installed command works from any directory without Go or a server.

```sh
anime-ascii makima
anime-ascii reze
anime-ascii denji
```

| Character | Size | Details |
| --- | --- | --- |
| Makima | 16 × 9 | Braids, ringed eyes, shirt and tie. |
| Reze | 16 × 9 | Bob haircut, green eyes and choker. |
| Denji | 15 × 8 | Spiky hair, open eyes and shirt collar. |

## Colors

Tokyo Night is the default theme. Chibi portraits use stable colors for hair, eyes and clothing. `-theme rainbow` and `-theme pastel` are also available.

```sh
anime-ascii reze -color always
anime-ascii makima -color never
anime-ascii denji -theme pastel
```

`-color auto` enables color in a terminal, unless `NO_COLOR` is set or `TERM=dumb`. Ordinary file and pipe output is plain text. Use `-color always` to preserve colors when piping to `less -R`.

## Optional image mode

The larger portraits converted from illustrations remain available through `-mode image`. This mode supports widths of 32, 40, 48, 64, 80, 100, 160, 200 and 240 columns. Its automatic size stays compact; explicit larger sizes may need scrolling.

```sh
anime-ascii makima -mode image
anime-ascii reze -mode image -width 80 -color always
anime-ascii denji -mode image -width 240 -color always | less -RS
```

`-width` controls image conversion and image-mode portraits. Chibi sketches keep their original small size. `-animate` is available in image mode for the moving gradient.

Vector image-mode previews: [Makima](docs/images/makima-terminal.svg), [Reze](docs/images/reze-terminal.svg), [Denji](docs/images/denji-terminal.svg).

The converted image-mode portraits use illustrations by [eastman](https://wall.alphacoders.com/big.php?i=1284789), [SpringHaru](https://twiman.net/user/1348570965135949830/illust/1970415339994521708), and [masoq](https://twitter.com/masoq095/status/1645434593502437376). Exact source crops and conversion settings are recorded in [artwork sources](internal/art/portraits/SOURCES.md).

## Convert your own image

PNG, JPEG and GIF files work locally:

```sh
anime-ascii -image picture.png -width 80 -invert -contrast 1.5 -color always
anime-ascii -image picture.png -width 80 > picture.txt
```

`-invert` makes dark pixels produce dense characters. `-contrast` accepts 0.1 to 3. Local CLI files and HTTP uploads are limited to 8 MiB and 12 million pixels. Output width is 10 to 240 characters, with a maximum height of 180 lines. GIF conversion uses the first frame; transparency blends onto black.

For a source crop or stronger line preservation:

```sh
go run ./cmd/convert -image picture.jpg -crop 10,20,500,500 -width 80 -invert -edges
```

## Optional server

```sh
git clone https://github.com/Godsdar/anime-ascii.git
cd anime-ascii
go run ./cmd/server
```

The server listens on `127.0.0.1:8080`. Local CLI use does not need it. To use HTTP explicitly:

```sh
anime-ascii reze -server http://127.0.0.1:8080
curl -fsS 'http://127.0.0.1:8080/art?style=makima&color=tokyonight'
```

| Route | Response |
| --- | --- |
| `GET /art?style=makima` | Small plain-text chibi portrait. |
| `GET /preview?style=reze` | Colored browser preview. |
| `GET /art?style=denji&mode=image&width=80` | Converted image-mode portrait. |
| `POST /convert?width=80&invert=true&contrast=1.5` | ASCII from raw image bytes. |
| `GET /fast`, `GET /slow` | Responses after 200 ms and one second for the URL-checking demo. |

Both portrait routes accept `style`, `mode` (`chibi` by default), `color` (`none`, `tokyonight`, `rainbow`, `pastel`), and `frame`. `width` applies to image mode. `/art` defaults to plain text; `/preview` defaults to Tokyo Night.

## Development

```sh
mkdir -p bin
go build -o bin/anime-ascii ./cmd/anime-ascii
go build -o bin/ascii-server ./cmd/server
go test -race ./...
go vet ./...
```

The CLI is intended for macOS and Linux and has been tested on macOS. Chibi sketches are in `internal/art/chibi.go`; image conversion and color handling are in `internal/art/`. The separate `cmd/check` program demonstrates concurrent URL checks.

Characters: Makima, Reze and Denji from Chainsaw Man, by Tatsuki Fujimoto.
