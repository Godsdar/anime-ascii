# Artwork sources

These ASCII portraits were converted from existing illustrations on October 12, 2026, at 80, 100, 160, 200 and 240 columns. The server defaults to 160 columns. The source raster images are not included in this repository.

| Character | Artist | Source page | Conversion |
| --- | --- | --- | --- |
| Makima | eastman | [Wallpaper Abyss](https://wall.alphacoders.com/big.php?i=1284789) | 1920×1536 source; crop `400,0,1510,1500`; 80/100/160/200/240 columns; inverted brightness; contrast 1.5. |
| Reze | SpringHaru | [Illustration on Twiman](https://twiman.net/user/1348570965135949830/illust/1970415339994521708) | 894×1200 source; crop `50,0,870,1150`; 80/100/160/200/240 columns; inverted brightness; contrast 1.5. |
| Denji (Chainsaw Man) | Official Reze Arc character visual | [Movie Natalie gallery](https://natalie.mu/eiga/gallery/news/616805/2555644) | 1280×1708 source; crop `0,0,1280,910`; 80/100/160/200/240 columns; inverted brightness; contrast 1.5. |

The conversion uses `cmd/convert`, area-averaged luminance, and the ASCII density ramp ` .:-=+*#%@`. Character proportions account for terminal cells being taller than they are wide.

Characters: Makima, Reze and Denji from Chainsaw Man, by Tatsuki Fujimoto. Artwork credit belongs to the source artists.
