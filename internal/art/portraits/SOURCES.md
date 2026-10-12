# Artwork sources

These ASCII portraits were converted from existing illustrations on October 12, 2026. The source raster images are not included in this repository.

| Character | Artist | Source page | Conversion |
| --- | --- | --- | --- |
| Makima | eastman | [Wallpaper Abyss](https://wall.alphacoders.com/big.php?i=1284789) | 1920×1536 source; crop `400,0,1510,1500`; 100 columns; inverted brightness; contrast 1.5. |
| Reze | SpringHaru | [Illustration on Twiman](https://twiman.net/user/1348570965135949830/illust/1970415339994521708) | 894×1200 source; crop `50,0,870,1150`; 100 columns; inverted brightness; contrast 1.5. |

The conversion uses `cmd/convert`, area-averaged luminance, and the ASCII density ramp ` .:-=+*#%@`. Character proportions account for terminal cells being taller than they are wide.

Characters: Makima and Reze from Chainsaw Man, by Tatsuki Fujimoto. Artwork credit belongs to the source artists.
