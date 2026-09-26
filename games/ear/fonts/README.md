# signs.ttf

Les glyphes que Go Regular n'a pas : ♭ ♮ ♯ (U+266D à U+266F), et le
double dièse 𝄪 et le double bémol 𝄫 (U+1D12A et U+1D12B), qui sont des
signes à part entière et pas deux signes simples accolés. Découpés dans
Noto Music 2.003, sans autre modification :

```sh
pyftsubset noto-music-music-400-normal.woff \
    --unicodes=U+266D-266F,U+1D12A-1D12B --output-file=signs.ttf
```

Noto Music est sous licence SIL Open Font License 1.1, reproduite dans
`OFL.txt` : elle permet d'embarquer et de redistribuer la police,
découpée ou non.

Pour ajouter un glyphe, refaire le découpage avec la plage élargie.
