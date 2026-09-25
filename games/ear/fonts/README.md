# signs.ttf

Les glyphes que Go Regular n'a pas : ♭ ♮ ♯ (U+266D à U+266F) et
l'ellipse (U+2026). Découpés dans DejaVu Sans 2.37, sans autre
modification :

```sh
pyftsubset DejaVuSans.ttf --unicodes=U+266D-266F,U+2026 --output-file=signs.ttf
```

DejaVu Sans est sous la licence Bitstream Vera, avec les ajouts DejaVu
dans le domaine public. Elle permet de redistribuer et de modifier la
police, à condition que le nom d'une version modifiée ne contienne ni
« Bitstream » ni « Vera », ce qui est le cas ici. Le texte complet est
sur https://dejavu-fonts.github.io/License.html.

Pour ajouter un glyphe, refaire le découpage avec la plage élargie.
