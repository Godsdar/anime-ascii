# Artwork sources

These face crops were converted from existing illustrations on October 12, 2026. The source raster files are not included. Resolutions: 32, 40, 48, 64, 80, 100, 160, 200 and 240 columns. The server defaults to 40 columns; the CLI selects a compact size that fits the terminal.

| Character | Artist | Source | Crop and conversion |
| --- | --- | --- | --- |
| Makima | eastman | [Wallpaper Abyss](https://wall.alphacoders.com/big.php?i=1284789) | 1920×1536 source; crop `650,330,1270,950`; invert; contrast 1.4; preserve edges. |
| Reze | SpringHaru | [Twiman illustration](https://twiman.net/user/1348570965135949830/illust/1970415339994521708) | 894×1200 source; crop `220,270,735,785`; invert; contrast 1.5; preserve edges. |
| Denji (human) | masoq | [Original face study](https://twitter.com/masoq095/status/1645434593502437376) | 1097×1200 source; crop `110,730,465,1095`; invert; contrast 1.4; preserve edges. |

## Eye preservation

The raster conversion uses `cmd/convert -edges`. Dark samples retain 65% of the cell weight so fine eyelids do not disappear into a pale average. After quantization, annotated pupil centers are projected onto the character grid and preserved as `@`. This keeps each pupil represented even in the smallest supported portrait. The CLI and browser preview apply gold or green accents to eye regions; the rest of each portrait keeps its selected gradient.

Pupil annotations are relative to the crops above:

| Character | Left pupil | Right pupil |
| --- | --- | --- |
| Makima | `(178, 220)` | `(450, 220)` |
| Reze | `(71, 202)` | `(313, 140)` |
| Denji | `(80, 161)` | `(211, 134)` |

Characters: Makima, Reze and Denji from Chainsaw Man, by Tatsuki Fujimoto. Artwork credit belongs to the source artists.
