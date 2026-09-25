# Glossaire

Ce document fixe le vocabulaire du projet. Un mot, un sens, dans tout le
dépôt : code, commentaires, tests, documentation, messages d'erreur.

Le domaine musical a un vocabulaire ambigu de naissance. « Note » désigne
selon le contexte un signe écrit, une hauteur sonore ou une touche
enfoncée. « Degré » désigne un rang dans une gamme et un chiffrage
d'accord. « Ton » désigne une tonalité et un intervalle de deux
demi-tons. La bibliothèque précédente confondait les trois. Les tableaux
ci-dessous tranchent.

## Types du noyau

Le noyau ne connaît aucun nom, aucune orthographe, aucune langue. Il
manipule des nombres et des ensembles.

| Type | Définition | Représentation |
|---|---|---|
| `Pitch` | Une hauteur absolue, identifiée à son numéro MIDI. Do central = 60. | `int8` |
| `Semitones` | Un écart entre deux hauteurs, en demi-tons. Signé. | `int8` |
| `PitchClass` | Une hauteur modulo l'octave. Do = 0, do♯ = 1, si = 11. Sans orthographe : do♯ et ré♭ sont **la même valeur**. | `uint8`, domaine 0 à 11 |
| `PitchSet` | Un ensemble de `PitchClass`. Masque de douze bits, bit *n* pour la classe *n*. | `uint16` |
| `Degrees` | Une distance en degrés de gamme. Position contre écart, comme `Pitch` contre `Semitones`. | `int8` |
| `Interval` | Un écart entre deux notes : un nombre de degrés **et** un nombre de demi-tons. Les deux sont nécessaires. | struct |
| `Tetrachord` | Une cellule de quatre notes, désignée par ses trois pas. Le littéral hexadécimal est la désignation : `0x131` est le tétracorde 1 3 1. | `uint16` |
| `ScalePattern` | La structure d'une gamme, indépendante de sa tonique. Masque de douze bits. | `uint16` |
| `ChordPattern` | La structure d'un accord, déployée sur deux octaves pour distinguer une neuvième mineure d'une seconde mineure. Masque de vingt-quatre bits. | `uint32` |

`PitchClass` et `PitchSet` partagent une représentation mais pas un rôle.
Le premier est un élément, le second un ensemble. Ne pas les faire
converger.

`ScalePattern` et `PitchSet` partagent une représentation et sont
convertibles, mais un pattern est ancré sur une tonique implicite en
position 0, alors qu'un `PitchSet` est absolu. La conversion est
explicite dans les deux sens et porte un nom.

## Pourquoi `Interval` porte deux nombres

Ni l'un ni l'autre ne suffit. Six demi-tons, c'est une quarte augmentée
ou une quinte diminuée, et le compte de demi-tons ne tranche pas. Trois
degrés, c'est une quarte de qualité inconnue. Ensemble ils sont exacts.

Une version antérieure de ce glossaire proscrivait `Interval` en
affirmant que ce choix était de l'orthographe. C'était faux : il est
porté par le degré, qui est structurel. Un accord dont la quarte est
haussée n'est pas le même qu'un accord dont la quinte est abaissée,
quel que soit le nom qu'on leur donne.

Ce qui appartient à `naming`, c'est le **mot**. Le noyau calcule
`Interval.Alteration()` et rend un nombre ; `naming` en fait
« augmentée », « diminuée », « mineure » ou « majeure », avec la
partition entre intervalles justes et imparfaits pour son compte.

## Notions structurelles

| Terme | Sens retenu | À ne pas confondre avec |
|---|---|---|
| `Degree` | Rang d'une note dans un pattern de gamme, compté à partir de 1. Borne haute : la taille du pattern, pas sept. Une pentatonique a des degrés 1 à 5. | `Degrees`, qui est une distance et non un rang |
| `Tonic` | Hauteur de référence d'une gamme. Champ de `Scale`. | `Root` |
| `Root` | Hauteur de référence d'un accord. Champ de `Chord`. | `Tonic` |
| `Scale` | Couple tonique + `ScalePattern`. | `ScalePattern` seul |
| `Chord` | Couple fondamentale + `ChordPattern`. | `ChordPattern` seul |
| `Phrase` | Une forme mélodique : suite d'écarts signés depuis la première note. Jamais repliés. | `Progression`, qui replie |
| `Direction` | Un pas de contour : `Up`, `Down`, `Level`. `Level` et non `Same`, qui voisinerait avec `Match.Same`. | |
| `Tonality` | Un contexte tonal : tonique + pattern **heptatonique**. N'existe que là où invoquer une tonalité a du sens, typiquement pour identifier une cadence ou épeler par degré. | `Scale`, qui accepte n'importe quel pattern |

`Tonality` appartient au noyau et non à `naming`, parce qu'elle ne porte
aucun nom : c'est une tonique et un pattern, deux entiers. `naming` s'en
sert pour épeler, `analysis` pour raisonner en fonctions et en cadences.
Les deux l'importent du noyau.

Son constructeur **refuse** un pattern qui n'a pas sept notes. Une gamme
par tons ou diminuée n'a pas de tonalité, et ce n'est pas un cas
dégradé à rattraper.

L'**absence** de `Tonality` est un état légitime, porté par sa valeur
zéro et reconnu par `IsZero()`. `analysis` doit pouvoir répondre qu'il
n'a inféré aucun contexte tonal sans que ce soit une erreur, et
l'interface ne doit jamais présenter le défaut d'épellation comme une
tonalité reconnue : la valeur zéro n'est pas do majeur.

`Tonic` et `Root` désignent la même idée dans deux contextes. Les
distinguer coûte un mot et évite qu'on écrive `scale.Root` en pensant à
une gamme et `chord.Tonic` en pensant à un accord.

## Types de `naming`

L'orthographe vit exclusivement dans `naming`. Le noyau l'ignore.

| Type | Définition |
|---|---|
| `SpelledNote` | Une `PitchClass` associée à une lettre et une altération. Mi♯ et fa sont deux `SpelledNote` distinctes pour la même `PitchClass`. |
| `Accidental` | Une altération écrite : bécarre, dièse, bémol, double dièse, double bémol. Valeur entière en demi-tons, domaine -2 à +2. |
| `Locale` | Table de noms de lettres et de patterns pour une langue. Passé en paramètre, jamais global. |

Le mot **alteration** ne figure nulle part. C'est un gallicisme :
l'anglais dit *accidental* pour le signe écrit.

## Règle d'épellation

L'épellation prend **toujours** un contexte tonal. Elle est totale :
elle réussit toujours et ne retourne pas d'erreur.

`naming` ne devine rien. Il ne cherche pas de fondamentale, il ne
reconnaît pas de pattern, il ne choisit pas de tonalité. C'est
`analysis` qui infère le contexte et le lui fournit. Une responsabilité
par paquet.

1. **Épellation par degré.** Dans une gamme à sept notes, chaque lettre
   apparaît exactement une fois. Fa♯ majeur donne fa♯ sol♯ la♯ si do♯
   ré♯ mi♯, avec un mi♯ et pas un fa. Sol♯ majeur donne un fa double
   dièse : le cas est légitime et doit être couvert par les tests.
2. **Hauteurs hors du contexte.** Une classe absente de la gamme s'épelle
   en altérant le degré le plus proche. À égalité entre deux degrés, on
   hausse celui du dessous, comme le fait la convention par défaut. Le
   moteur ne connaît pas la direction du mouvement, donc elle n'entre
   pas dans la règle.
3. **Contexte absent.** Do naturel majeur. C'est une **convention
   d'épellation de dernier recours**, pas une tonalité : elle ne se
   nomme pas, ne s'affiche pas, et ne doit jamais être présentée au
   joueur comme un contexte reconnu.

Le contexte ne peut pas osciller librement : voir la note sur les deux
constantes de temps dans ARCHITECTURE.md.

## Analyse

| Terme | Sens retenu |
|---|---|
| `Snapshot` | Les hauteurs qui sonnent à un instant. Des hauteurs et non des classes, pour que la basse puisse trancher. |
| `Reading` | Une identification d'un instantané : fondamentale, pattern normalisé, tétrade, extensions. Pas de score : une lecture existe ou n'existe pas. |
| `Key` | Une touche du clavier physique, enfoncée ou relâchée. |

Attention : `Key` en anglais désigne aussi la tonalité. Dans ce dépôt il
désigne **uniquement** la touche physique. La tonalité se dit
`Tonality`, la gamme se dit `Scale`, la hauteur de référence se dit
`Tonic`. Le mot `Key` ne doit jamais apparaître dans `naming`.

## Phrase contre Progression

Les deux sont des formes relatives et une seule règle les sépare.

Une **progression** rebase chaque accord par le mouvement le plus
court : une quinte montante devient une quarte descendante, parce qu'un
auditeur entend le plus court chemin entre deux fondamentales et que
l'octave où le pianiste les voice ne dit rien.

Une **phrase** fait l'inverse. Monter d'une septième majeure et
descendre d'une seconde mineure tombent sur la même classe et ne sont
pas la même mélodie : l'une saute, l'autre marche. Les écarts sont donc
signés, jamais repliés, et peuvent dépasser l'octave.

Corollaire sur `Match.Shift` : pour une progression c'est une distance
entre classes, pour une phrase une distance entre hauteurs. Une mélodie
chantée une octave plus haut est transposée de douze, pas de rien.

## Ce qu'un slot attend

Un slot de jeu ne demande pas un accord, il demande d'**atteindre une
cible**. La différence n'est pas de la sévérité, elle est structurelle :
un substitut tritonique n'est la dominante de do que par rapport à do,
et le même ré♭7 est la dominante de sol♭ ailleurs. La fonction ne se lit
donc pas sur l'accord seul, et une table tétrade vers fonction est
nécessaire sans jamais suffire.

| Terme | Sens retenu |
|---|---|
| `Target` | Ce qu'un slot demande d'atteindre : une fondamentale, éventuellement une tétrade, et la tonalité qui sert de juge |
| `Approach` | Ce que le joueur a produit dans le slot : la suite d'accords sur le chemin |
| `Resolution` | Ce que le chemin a été **constaté** faire. Des faits, jamais un score |

Un accord isolé est le cas à un pas d'un mouvement, pas un type à part.

Trois choses se dérivent de la théorie, sans table :

- **La dominante** se reconnaît à son triton entre tierce et septième.
  Deux fondamentales à un triton l'une de l'autre le partagent, donc le
  substitut n'est pas un cas particulier à lister, c'est le même objet
  vu d'ailleurs. Vérifié par balayage : exactement deux fondamentales
  sont dominantes de chaque cible.
- **Le retard** est deux accords pour un degré : même fondamentale, la
  quarte ou la seconde descend sur la tierce d'un accord suivant.
- **Les extensions** se valident contre le mode du degré, jamais dans
  l'absolu.

Ce qui **ne** se dérive pas : jusqu'où un chemin a le droit de
s'éloigner. Un ii-V à la place d'un V, une chaîne de dominantes
secondaires, rien dans la théorie transcrite ici ne dit où ça s'arrête.
Ces règles viendront d'un musicien et se placeront au-dessus de
`Resolution`, jamais dedans.

## Mots proscrits

| Mot | Pourquoi | À la place |
|---|---|---|
| `Note` | Trois sens concurrents. | `Pitch`, `SpelledNote` ou `Key` selon le cas. |
| `Alt`, `Alteration` | Gallicisme. | `Accidental` |
| `Ambitus` | Latin dans du code anglais. | `Range` |
| `Tone` | Ambigu : tonalité ou deux demi-tons. | `Semitones(2)` ou `Scale`. |
| `CurrentLocale` | Singleton mutable. | Paramètre explicite. |

## Points encore ouverts

- La lecture tétracordale des modes des systèmes altérés. Sept notes
  forcent la coupe après la quatrième, ce qui est musicalement fondé
  pour le système naturel et l'essentiel de la mélodique, mais
  mécanique pour 18 des 35 modes : les formes produites n'ont pas de
  nom et certaines ne couvrent pas même une quarte. À reprendre
  quand la pratique à l'oreille dira ce qu'on en fait.

## Points tranchés

- `SpelledNote`, et non `WrittenNote` : *spelling* est un acte oral, le
  mot ne présuppose aucune portée.
- `Tonality` pour le contexte tonal, `Key` pour la touche physique.
- `Tonality` vit dans le noyau, heptatonique par construction. Son
  absence est sa valeur zéro, testée par `IsZero()` : ni pointeur ni
  booléen d'accompagnement.
- Conversions `ScalePattern` ↔ `PitchSet` : `ScalePattern.At(tonic)`
  pour sortir de l'espace des patterns, `NewScalePattern(set, tonic)`
  pour y entrer.
- `Interval` revient dans le noyau, avec `Degrees` pour la distance en
  degrés. `Extension` disparaît, remplacé par lui.
- `Function` est un champ de bits : un accord majeur est tonique,
  sous-dominante et dominante selon sa place. Il vit dans `harmony` et
  non dans `naming` : la fonction harmonique est de la structure, pas
  de l'orthographe, et `Target` doit pouvoir l'interroger.
- `System` vit dans `harmony` : désigner un mode par sa gamme mère et
  son degré est une opération de structure, pas de nommage.
- Les tétracordes se **désignent par leurs pas** : « le tétracorde
  1 1 3 ». Seuls cinq ont un nom, ceux que la pratique utilise : majeur
  `2 2 1`, mineur `2 1 2`, phrygien `1 2 2`, lydien `2 2 2`, harmonique
  `1 3 1`. On résiste à la tentation de nommer le reste, plutôt que
  d'emprunter un vocabulaire à une autre culture ou à une filiation
  douteuse.
- L'identification des accords est **déterministe**. Les règles de
  construction tranchent avant toute comparaison ; ce qui ne
  correspond à aucune tétrade n'est pas un accord. Pas de score.
