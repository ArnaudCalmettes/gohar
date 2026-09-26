# Architecture

## Principe directeur

Le noyau est une bibliothèque de calcul pur. Il n'a pas d'entrées-sorties,
donc **il n'a pas de ports**. Lui en donner reviendrait à mettre une
interface devant une opération bitwise.

Les ports appartiennent au seul endroit qui touche le monde extérieur :
les jeux, qui lisent un périphérique et font entendre du son. Le moteur
d'analyse n'en a pas non plus, pour les raisons données plus bas.

Cette asymétrie est délibérée. Ne pas la corriger par souci de symétrie.

## Arborescence

Un dépôt, plusieurs modules. Le découpage suit le graphe de
dépendances, pas le rangement.

```
gohar/
  go.work                    use ./dex ./games ./harmony ./synth
  docs/                      glossaire, architecture

  harmony/     go.mod        théorie musicale, zéro dépendance
    pitch.go                 Pitch, Semitones, Transposable
    interval.go              Degrees, Interval
    pitchclass.go            PitchClass
    pitchset.go              PitchSet
    scale.go                 Degree, ScalePattern, Scale
    tonality.go              Tonality
    chord.go                 ChordPattern, Normalize, Chord
    function.go              Function, FunctionOf
    system.go                System, ModeOf : gamme mère et degré d'un mode
    tetrachord.go            Tetrachord, lecture tétracordale d'une gamme
    progression.go           Step, Progression, Match
    phrase.go                Phrase, Direction

    naming/                  orthographe, locales, les 35 modes
    analysis/                identification déterministe, moteur

  dex/         go.mod        collection du joueur, dépend de harmony
    notion.go                identité d'une notion, forme persistée
    dex.go                   marques, entrées, faits, questions

  synth/       go.mod        synthèse et sortie audio, dépend d'oto
    tuning.go                numéro de touche vers fréquence
    engine.go                voix, enveloppe, mélange, io.Reader
    device.go                ouverture d'oto et discipline des buffers

  games/       go.mod        Ebitengine, ark, MIDI
    keyboard/                port des touches, seul endroit qui voit gomidi
    keys/                    clavier jouable, mesure de latence bout en bout
    otolatency/              sonde de la seule moitié audio
    latency/                 sonde historique, par ebiten/v2/audio
```

Le dex est commun à tous les jeux et n'en connaît aucun. Sa conception
est dans `docs/DEX.md`.

## Pourquoi plusieurs modules

Un seul `go.mod` à la racine mettrait Ebitengine dans le graphe de
dépendances de quiconque importe `harmony` pour transposer trois
accords. Les frontières de modules sont donc là où les dépendances
changent, et nulle part ailleurs.

`harmony` n'a aucune dépendance de production. `synth` en a une seule,
oto, et la raison est plus bas : il possède la sortie audio parce que
personne d'autre ne peut la régler correctement. La synthèse elle-même
reste un `io.Reader` sur du PCM, que la bibliothèque standard suffit à
écrire, et le moteur de jeu qui l'affiche est l'affaire du jeu. C'est la
même discipline que pour les ports : une bibliothèque ne connaît pas la
techno qui l'appelle.

Le `go.work` à la racine fait compiler l'ensemble en développement sans
`replace` à maintenir. Deux conséquences à connaître : les étiquettes
de version sont préfixées, `harmony/v0.1.0` et `synth/v0.1.0` étant
deux tags distincts du même dépôt, et `go test ./...` ne traverse pas
les frontières de modules, donc il faut viser chaque module.

## Ce que le synthé ne partage pas avec l'harmonie

La conversion d'une hauteur en fréquence vit dans `synth`, jamais dans
`harmony`. C'est un bon marqueur de frontière : `Pitch` existe
précisément pour s'affranchir du tempérament, et lui donner un diapason
déferait l'abstraction pour laquelle il a été construit.

## Le son ne sort pas par Ebitengine

`synth` possède la sortie audio et parle à oto directement. Ebitengine
garde le graphisme et les entrées, et le paquet `ebiten/v2/audio` n'est
pas utilisé du tout.

Ce n'est pas un contournement, c'est la seule façon d'atteindre une
latence jouable. La raison est précise et vaut d'être écrite, parce que
la question « pourquoi ne pas utiliser l'audio d'Ebitengine ? » se
reposera.

### Les couches, du programme au haut-parleur

Quatre files d'attente séparent une décision de jouer du son qu'on
entend. Chacune a son défaut, et chacun de ces défauts est calé pour de
la lecture de fichier.

| Couche | Réglée par | Défaut |
|---|---|---|
| ring buffer du lecteur | `oto.Player.SetBufferSize` | **0,5 s** |
| buffer du périphérique | `oto.NewContextOptions.BufferSize` | **2 × 1024 frames, soit 42,7 ms** |
| serveur de son (PipeWire) | quantum du graphe | 1024 frames, soit 21,3 ms |
| bus USB et convertisseur | rien, c'est le matériel | trames de 1 ms |

Le ring buffer est rempli d'un coup au `Play()`, donc un son demandé ensuite
arrive derrière tout ce qui y est déjà. Le buffer de périphérique vient
du `else` de `driver_unix.go` dans oto, qui pose 1024 frames de période
quand on ne lui donne rien.

### Ce qu'Ebitengine ne permet pas

`audio.NewContext(sampleRate)` ne prend que la fréquence
d'échantillonnage et ne passe jamais `BufferSize` à oto. Les 42,7 ms de
périphérique deviennent donc un plancher inaccessible.
`audio.Player.SetBufferSize` existe et fonctionne, mais il ne règle que
le ring buffer au-dessus, pas la couche qui coûte cher.

### La règle

**Aucun buffer ne reste à son défaut. Tous sont posés explicitement.**

Pas par méfiance de principe : parce que chaque défaut de cette pile est
un bon choix pour jouer un fichier et un mauvais pour un instrument. Une
nappe de fond gagne à avoir une demi-seconde d'avance ; un clavier la
paie en mollesse.

### Mesures du 25 septembre 2026

Dell XPS 14 9440 sous Linux, PipeWire, Focusrite Scarlett Solo 4e
génération en USB, profil `output:analog-stereo`.

| Chemin | Plancher | Flam audible |
|---|---|---|
| `ebiten/v2/audio`, ring buffer à 30 ms | 71,5 ms | oui, franc |
| oto direct, périphérique et ring buffer à 5 ms | 12,5 ms | non |

Gigue de 2,5 ms sur le second, sans dérive ni décrochage. Le matériel
n'a jamais été en cause : le Scarlett tient 5 ms sans broncher, et le
quantum de PipeWire n'avait aucun effet tant que la couche au-dessus
demandait des périodes de 21,3 ms.

Le protocole de mesure tient dans `games/otolatency`. Le chiffre affiché
reste un plancher, puisque le matériel ajoute le sien sans qu'on puisse
le voir. La mesure qui tranche est le flam : on tape une touche fort, le
clic mécanique arrive par l'air et le son arrive après, et l'écart entre
les deux s'entend ou ne s'entend pas.

### Ce qui reste non mesuré

Sous charge. La sonde ne fait rien d'autre qu'attendre, alors qu'un jeu
dessine soixante fois par seconde et alloue. C'est là qu'on saura si
5 ms tient.

Le sujet est instruit, pas traité, et il n'est pas urgent : avant de
corriger un craquement toutes les trente secondes, il faut un jeu qui
mérite qu'on y joue trente secondes. Ce qu'on sait déjà, pour ne pas
refaire le raisonnement :

**Le mark assist est le vrai risque, pas les pauses.** Les pauses
stop-the-world de Go tiennent en quelques dizaines à quelques centaines
de microsecondes et passent sous un buffer de 5 ms. En revanche, une
goroutine qui alloue pendant une collecte se voit facturer du travail de
marquage sur place. Une goroutine qui n'alloue rien ne l'est jamais.
C'est la raison d'être de l'absence d'allocation dans `Engine.Read`, qui
est une protection et pas un réflexe de performance.

**Le coût du marquage suit le nombre de pointeurs, pas la mémoire.** Un
ECS par archétypes range les composants en tableaux contigus de
structures sans pointeurs, ce qui raccourcit chaque phase de marquage en
plus de réduire les allocations. Deux bénéfices distincts.

**Le levier suivant n'est pas `GOGC`.** C'est `debug.SetMemoryLimit` sur
une valeur confortable avec `GOGC=off` : le ramasse-miettes cesse alors
de suivre le taux d'allocation et ne se déclenche qu'en approchant de la
limite. Deux lignes dans un `main`, pour une empreinte bornée et connue.

**Un plafond qu'aucun réglage ne lève.** La goroutine audio est une
goroutine ordinaire, pas un thread en priorité temps réel comme celui
d'un DAW. C'est le meilleur argument pour garder quelques millisecondes
de marge plutôt que de viser le minimum : 5 ms est un compromis, pas un
score à battre.

**Comment trancher le jour venu.** Lire le pire cas après plusieurs
minutes de jeu soutenu, jamais la moyenne, qui cache exactement le
défaut qu'un musicien entend. Puis corréler avec `GODEBUG=gctrace=1` :
un pic qui tombe sur une ligne de trace désigne le ramasse-miettes, un
pic sans corrélation désigne l'ordonnanceur ou l'USB.

**Tranché le 26/09/2026.** Dix minutes de `ear` avec le clavier animé,
à 120 fps, en jouant des accords au MIDI par-dessus les séquences :
2 208 événements, p50 à 5 ms, p99 et maximum à 11 ms. Le maximum n'a
pas bougé entre 324 et 2 208 événements : c'est un plafond structurel,
à peu près deux périodes de lecture d'oto, et non un accident. La trace
compte 572 collectes, avec des pauses stop-the-world de 0,46 ms au plus,
et aucun pic ne s'y corrèle. Le ramasse-miettes ne se fait pas sentir.

Le jeu alloue pourtant beaucoup : une collecte toutes les une à deux
secondes en fin de session, le tas oscillant entre 14 et 28 Mo. La
goroutine audio n'en paie rien, puisqu'elle n'alloue pas. C'est du
gaspillage côté boucle de jeu, à réduire par opportunisme, pas un
risque pour le son.

L'autre moitié du trajet. Du doigt vers le programme, c'est-à-dire le
clavier maître, l'USB MIDI et la bibliothèque qui le lit. Probablement
petit, mais non mesuré, et c'est la somme des deux que le joueur sent.

La calibration reste au programme quoi qu'il arrive : ces chiffres sont
ceux d'une machine, et rien ne dit qu'ils ressemblent à ceux du joueur.

## Répartition entre `harmony` et `naming`

`harmony` porte **toutes les opérations d'usage courant**. Retrouver la
gamme mère et le degré d'un mode, comparer deux gammes, connaître la
fonction d'une tétrade : rien de tout cela ne doit exiger un objet de
nommage.

Quand un doute se présente, `naming` duplique un peu de logique plutôt
que d'obliger `harmony` à passer par lui. Le critère se lit simplement :
une opération qui n'a rien à voir avec un nom, une langue ou une
orthographe appartient à `harmony`, même si elle est apparue d'abord à
côté d'un nom. C'est ce qui a fait descendre `Function` puis `System`.

## Règles de dépendance

Une seule, mais stricte : **les flèches vont vers le noyau, jamais
l'inverse.**

| Paquet | Peut importer |
|---|---|
| `harmony` | rien du dépôt, et de la bibliothèque standard uniquement `math/bits`, `iter`, `slices`, `fmt`, `strconv` |
| `naming` | `harmony` |
| `analysis` | `harmony` |
| `dex` | `harmony` |
| `synth` | rien du dépôt, et d'externe uniquement oto |
| `keyboard` | rien du dépôt, et d'externe uniquement gomidi |
| `games` | tout |

Le noyau n'a **aucune méthode `String()` de présentation**. La conversion
en texte lisible appartient à `naming`. Un `String()` de débogage est
toléré s'il affiche la représentation brute (le masque en binaire, la
valeur numérique) et jamais un nom de note.

Aucun dot-import nulle part. Si un dot-import paraît nécessaire, c'est
que la frontière de paquet est au mauvais endroit.

Aucune variable de paquet mutable. Pas de `CurrentLocale`, pas
d'équivalent. Tout ce qui est configurable se passe en paramètre ou se
porte par une structure construite explicitement.

## Frontière de module

Ebitengine et la pile MIDI ne doivent jamais apparaître dans le `go.mod`
du noyau. Quelqu'un qui importe la bibliothèque pour transposer trois
accords ne doit pas tirer une pile graphique.

D'où le module `games/`, qui porte seul ces dépendances.

La lecture du MIDI vit donc dans `games/keyboard`, tranché en l'écrivant
et détaillé plus bas.

## Deux surfaces, et deux portes qu'elles gardent ouvertes

Deux bibliothèques externes touchent le monde réel : oto pour le son,
`gitlab.com/gomidi/midi/v2` pour les touches. Chacune n'est importée que
par **un seul fichier**.

| Bibliothèque | Unique point d'entrée | Devant |
|---|---|---|
| oto | `synth/device.go` | `synth.Device` |
| gomidi | `games/keyboard/midi.go` | `keyboard.Source` |

Ce n'est pas de la propreté décorative, c'est ce qui garde deux
possibilités ouvertes à coût faible, et il faut le maintenir même quand
ce sera tentant de faire autrement.

**Le navigateur.** Les trois couches savent y aller : Ebitengine compile
en `js/wasm`, oto a un pilote Web Audio, gomidi a `webmididrv` qui
attaque l'API Web MIDI en pur Go. Le choix se fait par balise de
compilation. Une bibliothèque appelée depuis dix endroits rendrait cette
cible théorique ; appelée depuis un fichier, elle reste atteignable.

**Sortir du C++.** `rtmididrv` passe par RtMidi, qui est du C++ et traîne
`libstdc++`. Sous Linux, un clavier MIDI USB est aussi un simple
périphérique caractère, `/dev/snd/midiC*D*`, qu'on ouvre avec `os.Open`
et dont on lit des octets MIDI bruts ; le décodage tient en une
cinquantaine de lignes et n'a besoin d'aucun cgo. On y perdrait le
partage du clavier avec un autre logiciel, le branchement à chaud,
l'énumération et les autres plateformes, donc ce n'est pas fait. Mais
c'est une implémentation de `keyboard.Source` de plus, pas une refonte.

À noter pour éviter un faux débat : **cgo est déjà là sous Linux**, parce
qu'Ebitengine y passe par X11 et OpenGL et qu'oto y passe par ALSA. La
bibliothèque MIDI n'introduit pas cette contrainte, elle en hérite. Sur
Windows et macOS, Ebitengine et oto sont passés à purego.

Un `keyboard.Source` n'est pas un port inventé pour faire joli, au
contraire de ceux écartés juste en dessous : il y a réellement plusieurs
sources de touches. Un clavier MIDI en est une, une séquence enregistrée
rejouée en mode entraînement en est une autre, et ce n'est pas un
bouchon de test mais un vrai mode d'un vrai jeu.

## Il n'y a pas de port dans `analysis`

Prévu, puis abandonné en concevant le moteur. La raison mérite d'être
écrite, parce que la tentation de les réintroduire reviendra.

`Clock` devait rendre les fenêtres testables sans dormir. Mais
l'appelant conduit : un jeu Ebitengine appelle `Advance` une fois par
frame avec l'heure de la frame, et le moteur ne dort jamais, ne tique
jamais, ne lance aucune goroutine. Un test qui veut sauter quatre
secondes passe une heure quatre secondes plus tard. Une interface qui
n'existerait que pour être simulée ne vaut pas mieux que le paramètre
qu'elle remplace.

`EventSource` devait produire les note-on et note-off. Mais c'est le jeu
qui lit son périphérique, ou qui rejoue une séquence figée en mode
entraînement, et qui appelle `NoteOn` et `NoteOff`. Le port existe, il
est un cran plus haut : il appartient à `games/`, où le MIDI et le
rejeu sont deux implémentations. Le moteur n'a aucune raison de savoir
laquelle.

Fabriquer ces deux interfaces pour faire ressembler `analysis` à un
hexagone aurait été exactement le culte du cargo que ce document écarte
plus haut pour le noyau.

## L'identification est déterministe, pas pondérée

Une première version de `analysis` comparait chaque pattern à chaque
fondamentale avec des poids réglables et rendait un classement. Mauvaise
forme.

L'harmonie a des règles de construction, et ces règles tranchent la
plupart de ce qu'une pondération aurait deviné. Qu'un ré dans un accord
de do soit une seconde ou une neuvième découle de la présence d'une
tierce, et aucun réglage n'améliore le fait de le savoir.

D'où l'ordre : `ChordPattern.Normalize()` applique les règles, on isole
la tétrade, on la cherche dans une table par égalité. Ce qui ne
correspond à rien **n'est pas un accord**, et le dire vaut mieux qu'un
classement plausible de réponses fausses. C'est aussi ce qu'une couche
de score ne sait pas faire : elle produit toujours un premier.

Ce qui reste réellement incertain est étroit : quelle note tenue est la
fondamentale, et que faire quand un accord symétrique en admet
plusieurs. Résolu par des règles d'ordre énoncées, pas par des nombres.

Corollaire sur les tétrades : les variantes sans quinte ont leur propre
entrée dans la table. Un voicing sans quinte n'est pas un accord auquel
il manque une note, c'est l'accord tel qu'on le joue, et en jazz la
quinte juste d'un accord de dominante est même déconseillée.

Corollaire sur l'hystérésis : sans score, il n'y a pas de marge à
franchir. Une lecture encore valide est conservée, ce qui suffit à figer
l'affichage sous un accord diminué 7 dont quatre fondamentales sont
également correctes.

## Découpage du moteur d'analyse

La reconnaissance se sépare en deux couches.

**Couche pure.** Un `Snapshot` de hauteurs vers des lectures ordonnées.
Pas d'état, pas d'horloge, pas de goroutine. Elle prend des hauteurs et
non des classes : la basse tranche entre des lectures que les classes
seules ne séparent pas, les mêmes quatre notes étant un do majeur avec
sixte ajoutée sur do, et un la mineur septième sur la. Elle se couvre
entièrement par tests en table.

**Couche temporelle.** Une machine à état qui agrège les touches sur une
fenêtre, appelle la couche pure, et applique la temporisation et la
stickiness pour que l'affichage n'oscille pas. Elle reçoit l'instant de
l'appelant, sans horloge à elle.

Ne pas mélanger les deux. La première est du calcul et se valide par
test, la seconde est du réglage et se valide à l'oreille. Une assertion
peut prouver que le moteur attend deux cents millisecondes ; elle ne
peut pas prouver que deux cents soit le bon nombre.

## Deux constantes de temps

`analysis` infère le contexte tonal et le fournit à `naming`, qui épelle
sans rien deviner. Si ce contexte oscille, l'orthographe oscille avec
lui et le joueur voit fa♯ devenir sol♭ à l'écran sans avoir rien changé
à son jeu.

Le moteur porte donc deux inerties distinctes, et il ne faut pas les
confondre :

- la reconnaissance d'accord, rapide, de l'ordre de la seconde, parce
  qu'un accord change souvent ;
- l'inférence de tonalité, lente, de l'ordre de la dizaine de secondes,
  parce qu'une tonalité change rarement et qu'une note étrangère isolée
  ne doit jamais la remettre en cause.

La seconde ne doit bouger que sur accumulation de preuves contraires,
pas sur un seul instantané contredisant le contexte courant.

## Conventions de test

`testify`, `assert` et `require`. `require` quand la suite du test n'a
pas de sens si l'assertion échoue, `assert` sinon.

Tests en table avec sous-tests nommés par ce qu'ils affirment, pas par
un numéro de cas. Le nom du sous-test doit se lire comme une phrase
quand il apparaît dans la sortie de `go test -v`.

Bancs sur le noyau, systématiquement. Les opérations du noyau sont des
instructions machine, et une allocation qui s'y glisse ne se voit pas à
la lecture. `-benchmem` fait partie du protocole habituel, pas d'une
enquête exceptionnelle.

Le noyau ne doit rien allouer sur ses chemins chauds. C'est une propriété
à vérifier par banc, pas une intention à écrire en commentaire.

## Choix de conception hérités

Quatre pièges de la version précédente, notés pour qu'on ne les
réintroduise pas.

**Le décalage par une hauteur.** `1 << int(p)` avec `p` de type hauteur
absolue panique en Go dès que `p` est négatif, et déborde en silence
au-delà de 31. La séparation `Pitch` / `Semitones` rend l'erreur
inexprimable, à condition que les opérations d'ensemble n'acceptent que
`Semitones` ou `PitchClass`, jamais `Pitch`.

**Le correctif d'indice en dur.** Le calcul d'intervalles d'accord
patchait la sortie à l'indice 3 pour les cas irréguliers, ce qui suppose
la position du septième degré dans le résultat. Les cas irréguliers se
décident depuis le pattern, jamais en corrigeant une sortie déjà
construite.

**La pondération à la place des règles.** Un système de poids produit
toujours un classement, donc des réponses plausibles là où il ne devrait
y en avoir aucune. Les contraintes de construction sont du savoir, pas
du réglage : elles s'appliquent avant, et elles peuvent refuser.

**Le tampon de sortie fourni par l'appelant.** Le motif
`Into(out []T) ([]T, error)` avec vérification de capacité est antérieur
aux itérateurs. Il porte deux valeurs d'erreur et un helper de
vérification pour un gain qui ne se mesure plus. Le noyau expose des
`iter.Seq`, l'appelant collecte s'il en a besoin.