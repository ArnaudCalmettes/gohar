# signs.ttf

The glyphs Go Regular lacks: ♭ ♮ ♯ (U+266D to U+266F), and the double
sharp 𝄪 and double flat 𝄫 (U+1D12A and U+1D12B), which are signs in
their own right rather than two single signs side by side. Cut from
Noto Music 2.003, with no other change:

```sh
pyftsubset noto-music-music-400-normal.woff \
    --unicodes=U+266D-266F,U+1D12A-1D12B --output-file=signs.ttf
```

Noto Music is licensed under the SIL Open Font License 1.1, reproduced
in `OFL.txt`, which allows embedding and redistributing the font, cut
or not.

To add a glyph, cut again with a wider range.
