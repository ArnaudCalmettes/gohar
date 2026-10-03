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
  go.work                    use ./charts ./dex ./games ./harmony ./synth
  docs/                      architecture, glossaire, chantiers, et la
                             conception du dex, de l'oreille, des
                             voicings et de l'analyse des grilles

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
    provenance.go            NamedScale, Provenances : d'où vient un accord
    tetrachord.go            Tetrachord, lecture tétracordale d'une gamme
    progression.go           Step, Progression, Match
    phrase.go                Phrase, Direction
    target.go                Target, Approach, Resolution
    approach.go              ApproachKind, ApproachOf : les préparations

    naming/                  orthographe, locales, noms des 35 modes,
                             des intervalles et des gammes nommées,
                             symboles d'accords (ChordStyle)
    analysis/                identification déterministe, moteur,
                             et l'analyse d'une grille : suite
                             d'accords (Changes), préparations,
                             passages, blocs et tonalités annoncées,
                             II-V consécutifs (Links), cellules
                             (Cells), seconde écoute (Reread),
                             phrases et tonalité du morceau
                             (ReadTune, Tune), tierce picarde,
                             plages modales, blues, degrés, tonique
                             pressentie et modulations (Sensed),
                             zones tonales (TonalAreas), forme par
                             récurrences (Sections), cadences
                             conclusives et demi-cadences
                             (Conclusions), pédales (Pedals)

  dex/         go.mod        collection du joueur, dépend de harmony
    notion.go                identité d'une notion, forme persistée
    dex.go                   marques, entrées, faits, questions
    persist.go               lecture et écriture JSON

  synth/       go.mod        synthèse et sortie audio, dépend d'oto
                             et de go-meltysynth
    tuning.go                numéro de touche vers fréquence
    engine.go                voix, enveloppe, io.Reader
    timbre.go                formes d'onde, timbres 8 bits
    queue.go                 Instrument, file de commandes, notes datées
    clock.go                 échantillons vers horloge murale
    mixer.go                 plusieurs instruments, une sortie
    soundfont.go             fichiers SF2, Sampler
    clip.go                  un son enregistré en WAV, Clip
    histogram.go             histogramme des délais
    device.go                ouverture d'oto et discipline des buffers

  charts/      go.mod        grilles venues d'autres logiciels
    ireal/                   URL iReal Pro : playlist, jetons, mesures,
                             dépliage de la forme, chiffrages, temps,
                             pont vers analysis.Changes
    cmd/analyse/             une grille et son analyse dans le terminal
    cmd/corpus/              l'analyse de playlists entières, comparée
                             à la tonalité que l'app déclare
    cmd/forms/               la forme de playlists entières, lue par
                             analysis.Sections

  games/       go.mod        Ebitengine, ark, MIDI
    keyboard/                port des touches, seul endroit qui voit gomidi
    ear/                     ear trainer : menu, degrés, tétracordes, modes
    walk/                    Walk with me : métronome, swing, temps
                             attendus, marqueur, walking bass générée,
                             batterie, enregistrement, grille, clavier
                             et bonhomme à l'écran
    keys/                    clavier jouable, mesure de latence bout en bout
    otolatency/              sonde de la seule moitié audio
    latency/                 sonde historique, par ebiten/v2/audio
```

Le dex est commun à tous les jeux et n'en connaît aucun. Sa conception
est dans `dex.md`.

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
utilisé que par la sonde de mesure `games/latency`.

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

### Les notes datées

Ce que joue le joueur part tout de suite : `NoteOn` s'applique au début
du prochain tampon, et `at` ne sert qu'à mesurer le retard. Ce que joue
le programme, un métronome ou une ligne de basse, ne peut pas suivre ce
chemin : la boucle d'Ebiten tourne toutes les 16,7 ms, et un temps
envoyé depuis elle tomberait n'importe où dans cet intervalle. Une
gigue de cet ordre s'entend sur un claquement de doigts.

D'où `ScheduleOn` et `ScheduleOff`, que fournit la `queue` commune à
tous les instruments. Le programme les appelle en avance, une fenêtre
d'une centaine de millisecondes avant le temps, et la goroutine audio
applique chaque commande sur son échantillon : elle remplit le tampon
jusqu'à cet échantillon, applique la commande, puis continue.

Il faut pour cela savoir à quel instant correspond chaque échantillon.
C'est le rôle de l'horloge du paquet (`clock`). Elle n'entend pas le
haut-parleur : elle voit seulement à quels instants le pilote appelle
`Read`. Elle en tire une estimation lissée, qui ne suit chaque écart
qu'au soixante-quatrième. Le pilote appelle par rafales, plusieurs
tampons d'un coup pour remplir son ring buffer, et une lecture isolée
ne dit pas grand-chose ; la moyenne, elle, suit la lente dérive entre
le quartz de la carte son et celui du système. Au-delà de 100 ms
d'écart (un décrochage, un portable mis en veille), l'horloge saute au
lieu de lisser.

Le son sort plus tard que l'échantillon écrit, du temps des files du
dessous. Ce retard est constant pour un réglage donné : il décale toutes
les notes datées de la même quantité, sans toucher à leur espacement.
C'est la latence de sortie, que la calibration mesure. Une note datée
dans le passé part tout de suite, en retard ; une fenêtre d'avance plus
large que le pire hoquet de la boucle de jeu évite d'en arriver là.

### Le mélangeur

Un `Device` ne prend qu'une source, et un `Engine` ne joue qu'un timbre.
Un jeu qui fait entendre une contrebasse, une batterie, un claquement de
doigts et le piano du joueur passe donc par un `Mixer`, qui additionne ses entrées,
chacune avec son gain, et écrête la somme comme le fait un `Engine` sur
un accord trop dense. Ses tampons sont alloués une fois pour toutes ;
une lecture plus longue qu'eux se fait en plusieurs passes.

Chaque instrument garde sa propre horloge. Le mélangeur les lit l'un
après l'autre dans le même appel, à quelques microsecondes d'écart :
leurs estimations coïncident, et deux notes datées du même instant sur
deux instruments tombent sur le même échantillon.

### Les soundfonts

Un fichier SF2 contient des instruments enregistrés, échantillonnés
note par note : une contrebasse, un piano, des kits de batterie. Le
`Sampler` en joue un, désigné comme en General MIDI par une banque et un
numéro de programme ; la banque 128 contient les kits, où chaque touche
est une percussion. Il passe par la même `queue` que l'`Engine`, et donc
par les mêmes notes datées et la même mesure des délais.

Le moteur est go-meltysynth (licence MIT, rien d'autre que la
bibliothèque standard). Il n'est importé que par `soundfont.go`. Ce
n'est pas une surface au sens de la section suivante, puisqu'il ne
touche pas le monde réel : il calcule des échantillons, et suivra le
jeu dans le navigateur sans rien changer.

Trois choix :

- **Le fichier dans `synth`, pas dans un sous-paquet.** Le `Sampler` a
  besoin de la `queue`, que rien n'exporte ; un sous-paquet aurait
  demandé d'exporter la file, ses commandes et son rendu pour un seul
  client.
- **Un bloc de 16 échantillons.** meltysynth calcule par blocs et
  n'applique une commande qu'au bloc suivant : une note datée tombe à un
  bloc près, un tiers de milliseconde, sous la gigue du pilote.
  L'`Engine` tombe sur l'échantillon, le `Sampler` à un bloc près.
- **Ni réverbération ni chorus**, qui coûtent plus cher que toutes les
  voix réunies, et dont une basse dans un mixage n'a pas besoin. À
  revoir avec le piano.

Un preset absent du fichier est une erreur. meltysynth, lui, se replie
en silence sur son premier preset : un jeu qui demande une contrebasse
et reçoit un clavecin sonnerait faux sans dire pourquoi.

Le parsing lit tous les échantillons en mémoire et alloue d'autant ; il
se fait une fois, avant le jeu, et plusieurs `Sampler` partagent le même
`SoundFont`. Ensuite, ni `Read`, ni `NoteOn` n'allouent.

### Les sons enregistrés

Certains sons n'existent dans aucune soundfont : le claquement de
doigts de *Walk with me*, par exemple. Le `Clip` joue un seul fichier
WAV, une percussion : chaque frappe le relance depuis le début, à sa
vélocité, et il s'éteint seul. La touche n'y change rien, le relâchement
non plus. Quatre prises peuvent se chevaucher.

Le WAV doit être échantillonné à 48 kHz, comme tout le paquet : une
autre fréquence est une erreur plutôt qu'un rééchantillonnage. PCM 16 ou
24 bits, ou flottants 32 bits, en mono ou en stéréo. Le décodeur ne lit
que les blocs « fmt » et « data », et saute les métadonnées.

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
| `charts` | `harmony` et `naming`, pour lire les chiffrages, et `analysis`, pour lui passer une grille (`Changes`) et l'afficher ; le reste de la lecture d'un format n'utilise que la bibliothèque standard |
| `games` | tout |

Un type du noyau peut avoir un `String()` lisible par un musicien, dans
une notation fixe, celle des grilles de jazz : `V`, `♭II7`,
`Dm7`. C'est ce qu'affichent les tests quand ils échouent et les outils
en ligne de commande. Ce que `naming` garde pour lui, c'est tout ce qui
dépend d'un choix : la langue, la notation (do ou C), l'orthographe
d'une note dans une tonalité.

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
par **un seul fichier**, hors des sondes de mesure (`games/latency`,
`games/otolatency`).

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

## Les sources

La nomenclature des modes vient de l'école de Bernard Maury. Trois
billets publiés sur Zeste de Savoir servent de source à `analysis` : la
représentation des accords, leur reconnaissance, et la game loop de
l'improvisateur. Le troisième explore des possibles, ce ne sont pas des
exigences. L'analyse des grilles suit le tome 1 d'*En Harmonie*
(Dericq et Guéreau), chapitres 8 à 10 (voir `grilles.md`).
