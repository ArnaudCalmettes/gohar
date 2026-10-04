# Le format des grilles

gohar lit aujourd'hui ses grilles dans le format d'iReal Pro. C'est la
porte d'entrée naturelle, puisque des milliers de grilles y sont
saisies, mais ce n'est pas un format pour un projet libre, et il lui
manque l'essentiel pour une analyse fine. Cette page dit ce qui manque,
ce qu'on a choisi à la place, et comment gohar s'en sert.

## Ce qui manque à iReal

- **Un format fermé et brouillé.** Une grille iReal voyage dans une URL
  dont le contenu est volontairement mélangé (`charts/ireal/scramble.go`
  le remet dans l'ordre). Rien n'en décrit la grammaire : on la
  reconstitue en lisant des grilles.
- **Des chiffrages à lui** : `^7` pour maj7, `-7` pour m7, `h7` pour le
  demi-diminué, loin de ce qu'un musicien écrit.
- **Une armure souvent fausse.** La tonalité déclarée par l'app se
  trompe régulièrement ; c'est une des raisons pour lesquelles l'analyse
  détecte la tonalité au lieu de la lire (voir `grilles.md`).
- **Des accords mis de côté.** Les petits accords d'iReal, une
  dominante chromatique sur un turnaround par exemple, n'ont pas de
  statut clair : le lecteur les écarte.
- **Pas de mélodie, pas de couleur modale.** Une grille modale dont la
  couleur n'est pas écrite ne dit pas son mode ; c'est la liste des
  grilles « mises à part » de `grilles.md`, auxquelles il manque
  l'information.
- **Des grilles qui ne sont pas les nôtres.** Une grille iReal saisie
  par quelqu'un est son travail. Les playlists du corpus restent donc
  hors du dépôt (`testdata/local`), et aucune ne peut servir de niveau
  dans un jeu publié.

## Les choix : ChordPro et ABC

Plutôt que d'inventer un format, gohar adopte deux formats ouverts, en
texte, répandus et outillés.

- **ChordPro** pour la grille. C'est le format libre des partitions
  « paroles et accords », lu par de nombreuses applications. Ses
  versions récentes ont des **sections de grille** qui portent
  exactement ce qu'une grille de jazz demande : les mesures, les temps,
  les reprises, les fins, les mesures répétées.
- **ABC** pour la mélodie. C'est une notation texte née dans la musique
  folk, courte, faite pour s'écrire à la main, et que des bibliothèques
  JavaScript dessinent en partition dans un navigateur. ChordPro la
  connaît déjà : un fichier ChordPro peut contenir des blocs ABC, qu'il
  confie à un outil externe pour les dessiner.

gohar lit et écrit un **profil** de ChordPro, le sien, sans viser tout
le format : de quoi écrire une grille de jazz, sa forme, ses accords,
plus tard sa mélodie. Un fichier de gohar reste un fichier ChordPro
valide, qui s'ouvre ailleurs.

## Ce que ChordPro offre, repris tel quel

- **Les métadonnées** : `{title}`, `{composer}` (répétable pour
  plusieurs compositeurs), `{key}`, `{time}`, `{tempo}`, `{copyright}`,
  et `{meta: nom valeur}` pour tout le reste, avec des noms libres.
- **Les sections** : une grille par section, `{start_of_grid
  label="A" shape="4x4"}` … `{end_of_grid}`.
- **Dans une grille** :
  - les barres `|`, `||` et `|.` ;
  - les reprises `|:` et `:|`, les fins `|1` et `:|2` ;
  - `%` et `%%` pour répéter une ou deux mesures ;
  - `.` pour un temps où rien ne change, `/` pour refrapper l'accord,
    `~` pour plusieurs accords dans un même temps.
- **La mélodie** : `{start_of_abc label="A"}` … `{end_of_abc}`.
- **Les extensions** : les directives préfixées `x_`, que ChordPro
  réserve à cet usage.

## Le profil « grille de jazz » de gohar

1. **Un morceau par fichier**, en `.cho`, encodé en UTF-8. Un corpus
   est un dossier.
2. **La tonalité et les sections peuvent s'écrire**, avec `{key}` et
   les étiquettes de section, mais l'analyse les **détecte** toujours :
   ce qui est écrit reste une annotation, lue en option pour affiner,
   jamais pour remplacer la détection. La même règle vaut pour le `K:`
   qu'ABC rend obligatoire.
3. **La place des accords dans la mesure**, avec les raccourcis d'usage
   des grilles de jazz, en lecture comme en écriture : un accord seul
   remplit la mesure (`| F7 |`), deux accords la partagent en deux
   (`| Gm7 C7 |`). Seule une coupe inégale s'écrit temps par temps, un
   point pour un temps où rien ne change (`| F7 . . D7 |`). La spec de
   ChordPro montre des mesures d'un seul accord sans points, mais ne dit
   pas comment lire deux accords sans points : l'usage des musiciens
   tranche.
4. **Le chiffrage est celui des musiciens**, lu par `naming` : `Bbmaj7`,
   `Am7b5`, `D7(b9)`, `C7alt`, `F/A`, `N.C.`. Les dièses et les bémols
   se lisent en ASCII (`b`, `#`) comme en Unicode (`♭`, `♯`), et
   s'écrivent toujours en ASCII : c'est ce que tous les logiciels
   ChordPro comprennent, et ce qui se tape le plus facilement à la main.
   L'Unicode reste pour l'écran. Le 7alt garde le sens qu'il a dans
   gohar : le locrien ♭4, C7(♭5, ♭9, ♭10, ♭13).
5. **Les accords optionnels**, ceux qu'un joueur ajoutera une fois plus
   avancé (les petits accords d'iReal), s'écrivent entre parenthèses,
   joints par `~` à l'accord qu'ils ornent : `| F6 D7~(Ab7) |`. Pour
   l'analyse, l'accord qu'ils ornent dure comme s'ils n'étaient pas là.
   Un autre logiciel les affiche tels quels, ce qui reste lisible.
6. **La forme se joue dans l'ordre du fichier**, section après section,
   avec leurs reprises et leurs fins. Une section qui revient peut se
   **rappeler** au lieu d'être recopiée, comme on l'écrit couramment pour
   un AABA : `{x_play: A}` rejoue la section d'étiquette A à cet endroit.
   **La coda** vient après `{x_coda}` : le chorus reboucle avant elle,
   et elle ne se joue qu'une fois, au dernier tour, pour conclure, comme
   iReal la joue. Les renvois (D.S., D.C.) s'écrivent dépliés : la
   section rejouée est rappelée, ou recopiée.
7. **La grille a sa propre licence**, distincte du `{copyright}` de la
   composition : `{meta: chart_author …}` et `{meta: chart_license …}`.
   `meta` étant standard, il est préféré à une directive `x_`.
8. **La mélodie** vient dans un bloc ABC par section, de même étiquette
   que sa grille, aligné mesure par mesure (voir plus bas).
9. **Les couleurs modales** viendront dans une directive `x_`, par
   section ou par accord, à dessiner le jour où l'analyse modale en aura
   besoin.

Notre blues, écrit dans ce profil :

```
{title: 12 Bar Blues}
{time: 4/4}
{meta: style Medium Swing}
{meta: chart_author gohar}
{meta: chart_license CC0-1.0}

{start_of_grid label="A" shape="4x4"}
| F7 | Bb7 | F7 | F7 |
| Bb7 | Bb7 | F6 | D7 |
| Gm7 | C7 | F6 D7~(Ab7) | G7 C7~(Gb7) |.
{end_of_grid}
```

Un AABA, avec son rappel :

```
{start_of_grid label="A" shape="4x4"}
| Bbmaj7 G7 | Cm7 F7 | Dm7 G7 | Cm7 F7 |
| Fm7 Bb7 | Ebmaj7 Ab7 | Dm7 G7 | Cm7 F7 |.
{end_of_grid}
{x_play: A}
{start_of_grid label="B" shape="4x4"}
| D7 | % | G7 | % |
| C7 | % | F7 | % |.
{end_of_grid}
{x_play: A}
```

Le rappel est une extension de gohar : un autre logiciel ChordPro ignore
`{x_play}`, et n'affiche alors le A qu'une fois.

## La mélodie

ChordPro ne représente pas la mélodie : il la délègue. Le bloc ABC est
recopié tel quel, et un outil externe le dessine. Pour gohar, qui veut
**analyser** la mélodie, ce bloc doit devenir une donnée :

- **un lecteur d'ABC**, sur le sous-ensemble utile : les hauteurs, les
  durées, les silences, les barres, les reprises, les chiffrages entre
  guillemets (`"F7"A2 c2`) ;
- **une règle d'alignement** : un bloc par section, de même étiquette
  que sa grille, mesure par mesure, avec les mêmes reprises ;
- **le `K:` d'ABC est une armure écrite**, pas une tonalité : l'analyse
  la traite comme l'armure d'iReal, une annotation.

La mélodie lèvera beaucoup des ambiguïtés que laisse une grille seule :
la note que le thème appuie sur un accord dit sa couleur, une tonalité
hésitante se tranche sur les notes du thème, et les grilles modales
« mises à part » trouvent enfin l'information qui leur manquait.

## Les droits

- **Une grille d'accords** saisie par nous, accords seuls, ne pose en
  général pas de problème : un enchaînement d'accords n'est pas protégé
  comme une mélodie. Ce sont ces grilles, sous licence libre, que le
  dépôt publie et que les jeux embarquent.
- **Une mélodie** est protégée tant que l'œuvre n'est pas dans le
  domaine public. Les mélodies de standards relevées dans des partitions
  restent dans un corpus **privé**, hors du dépôt, comme
  `testdata/local` aujourd'hui.

## Ce que fait gohar aujourd'hui

- **`charts/chordpro`** lit et écrit le profil, et en tire les accords
  que l'analyse lit (`Song.Changes`), avec les règles de
  `charts/ireal` : un accord écrit de nouveau continue, une mesure
  partagée l'est au temps près, arrondie au temps suivant. Une nuance :
  iReal compare l'écriture (`Bb^7` puis `Bb^` font deux accords), le
  profil compare l'accord (`Bbmaj7` puis `Bbmaj7` continue, quelle que
  soit la façon de l'avoir écrit).
- **Les extensions entre parenthèses** se lisent dans n'importe quel
  ordre, ajoutées à l'accord qui les précède : `G7sus4(b9)` comme
  `G7sus(b9)`, `C7(#11,b9)` comme `C7(b9,#11)`.
- **La conversion depuis iReal** (`chordpro.FromIReal`, et la commande
  `cmd/tochordpro` pour des playlists entières) écrit la forme **telle
  qu'elle se joue**, une fois : les reprises, les fins et les sauts d'un
  D.S. ou d'un D.C. sont dépliés, une section commence à chaque marque
  de répétition, et la coda vient après `{x_coda}`. Les accords
  sont ceux de la ligne de temps d'iReal, de sorte qu'une grille
  convertie donne à l'analyse exactement les mêmes accords que la
  grille d'origine, coda comprise. Un test le vérifie sur le blues et
  sur une coda, et un autre sur tout le corpus privé, là où il se
  trouve. Les petits accords d'iReal deviennent des accords optionnels.
- **Ce que la conversion ne fait pas encore** : garder les reprises et
  les fins au lieu de les déplier, et les changements de mesure. Une
  grille dont la mesure change est laissée de côté, et le dit : une
  quarantaine sur 1 678 dans le corpus privé, des musiques de jeux pour
  les trois quarts.

## Les étapes

1. Ce profil.
2. `charts/chordpro` : le lecteur et l'exporteur du profil, vers
   `analysis.Changes`, avec des tests qui se lisent en noms d'accords.
3. Un convertisseur d'iReal vers ChordPro, pour passer un corpus privé
   dans le nouveau format.
4. Le blues, *Tune Up* et *Autumn Leaves* en ChordPro, accords seuls,
   sous licence libre, lus par *Walk with me* (`games/walk/grids`). Le
   jeu montre les accords comme la grille les écrit (`Song.Spelled`).
5. `charts/abc`, le sous-ensemble d'ABC, et l'alignement sur la grille.
6. La mélodie dans l'analyse : les ambiguïtés qu'elle lève (voir
   `chantiers.md`).
