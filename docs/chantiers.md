# Chantiers

Ce qui est ouvert, et ce qui attend une décision.

Ce fichier existe pour qu'une conversation neuve reprenne sans rien
redécouvrir et sans rouvrir un débat déjà tranché. Les raisons des choix
faits sont dans `architecture.md`, le vocabulaire dans `glossaire.md`, la
conception de la collection dans `dex.md`, l'ear trainer dans
`oreille.md`, les voicings dans `voicings.md`, l'analyse des grilles
dans `grilles.md`. Ici il n'y a que ce qui reste à faire, les
décisions qu'il ne faut pas rouvrir, et, pour s'y retrouver, un bref
état de ce qui est fait.

Mettre une doc à jour n'est pas un chantier : ça fait partie de la
tâche qui la fait mentir.

## Où on en est

`harmony` tient : hauteurs, intervalles, ensembles, gammes, tonalités,
accords, fonctions, systèmes, tétracordes, et les 35 modes. Pour
l'analyse des grilles : les gammes nommées et la provenance d'un
accord, et les préparations d'un accord à l'autre. Testé, benché, vert.

`naming` sépare la langue (français, anglais) et la notation (signes
par défaut, mots en option). Il nomme les intervalles jusqu'à la
septième et les gammes nommées (majeure, mineure naturelle, mineure
harmonique, mineure mélodique, majeure harmonique) sur leur tonique,
et écrit les symboles d'accords dans le style choisi (`ChordStyle`).
Chaque mode a un nom systématique et des alternatives : le registre
parlé français et les alias.

`analysis` identifie les accords de façon déterministe, sans pondération,
et analyse une grille à la manière d'*En Harmonie* : préparations,
passages, blocs et tonalités qu'ils annoncent, degrés sur la tonique
installée et en crochets, tonique pressentie et modulations, blues.
Les cellules suivent les formes du livre : l'anatole, et le
III-VI-II-V-I en majeur, résolu sur la tonique entendue ; les II-V
contigus sont des marches. Les fiches du livre et de Siron concordent
sur leurs degrés et leurs modulations, et l'analyse tombe d'accord avec
la tonalité déclarée par l'app, ou vérifiée à l'oreille, sur
l'essentiel du corpus (le rapport
de `corpus` donne le compte du jour), sans la lire. `analyse` affiche
ce que le verdict pèse, le premier et le dernier accord et, quand ils
divergent, la durée de leurs toniques : le verdict s'explique. Il
réécrit les accords de la grille dans la tonalité entendue (voir
« L'orthographe entendue ») ; `-smells` liste ce qui sent encore.

`dex` a son corps, sa persistance JSON et `Components`. Restent
`Cooling` et `Discoverable`.

`synth` joue : conversion en fréquence, moteur polyphonique sans
allocation, discipline des buffers audio, et un histogramme des délais
à cases fixes.

`charts/ireal` lit les grilles d'iReal Pro, de l'URL jusqu'aux accords
dans l'ordre de jeu : playlist, jetons sans perte, mesures, dépliage de
la forme, chiffrages lus en accords de `harmony`, durées en temps.
Vérifié sur des playlists réelles, qui restent hors du dépôt, et sur
les fiches d'analyse d'*En Harmonie* et de Siron. Deux commandes :
`charts/cmd/analyse` affiche une grille annotée dans le terminal,
`charts/cmd/corpus` compare l'analyse à l'app sur des playlists
entières, les grilles auxquelles il manque l'information mises à part
(`charts/ireal/testdata/set-aside.txt`).

`games/keyboard` a deux sources : le clavier MIDI, qui saute les ports
Through quand aucun n'est demandé, et la séquence rejouée, qui joue des
notes datées sur sa propre goroutine avec la gigue d'un doigt.

`games/keys` fait sonner un clavier MIDI, environ 12,5 ms de plancher
mesuré, aucun flam audible.

`games/ear` est sorti du tracer bullet : un menu ouvre les degrés de la
gamme majeure, les quatre tétracordes et les modes du système naturel
(trois choix, ou les sept à la place de leur degré), en français ou en
anglais, en signes ou en mots, avec le dex persisté dans
`~/.config/gohar/dex.json`.

## Les décisions du dex

Elles sont dans `dex.md`, qui les tient : le dex reçoit des faits et ne
juge pas, une erreur compte pour du beurre, `FactChosen` marque aussi la
production, la couleur est ouverte tant que la cible ne nomme pas le
mode, et ce qui sonne sans être nommé apparaît en silhouette.

## Les décisions du nommage

**Les signes par défaut, dans les deux langues.** Les mots (si bémol,
phrygien bécarre 6) sont une option du `Namer`, pas de la `Locale`,
pour l'accessibilité : une synthèse vocale lit des mots, pas des
signes. Un jeu peut tenir deux namers, signes à l'écran et mots pour un
lecteur d'écran.

**Les symboles d'accords en réglage.** `ChordStyle` écrit la qualité
d'un accord : Cmaj7, Cm7, Cm7♭5, Cdim7 par défaut, comme *En
Harmonie* ; C♮7 ou CΔ7, C-7, Cø, C°7 en option. Le bécarre ne se lit
bien qu'en exposant, sur une grille gravée : sur une ligne de texte,
C♮9 se lirait comme un do avec une neuvième bécarre. Les extensions
naturelles empilées montent sur la septième (C9 : la 9e ; C13 : la 9e
et la 13e, sans 11e ; Cm11 : la 9e et la 11e ; Cm13 : la 9e, la 11e et
la 13e). Sur une tierce majeure, la 11e reste toujours à part :
C13(11). Le reste va entre parenthèses sans espace, la plus grave
d'abord : C7(♭9,♯11). Pas d'espace, parce que sur une grille une espace sépare
deux accords. Le mineur-majeur s'écrit Cm(maj7), et ses extensions
entrent dans la même parenthèse sans monter : Cm(maj7,9,♯11). Le 7alt s'écrit 7alt.

**Le nom systématique par défaut, sans exception.** La base naturelle
suivie de chaque degré altéré : phrygien ♮3, lydien ♯5. Les noms
raffinés (phrygien majeur, phrygien dominante, lydien augmenté,
phrygien sixte majeure) sont des alternatives, `ModeAlternatives`,
pour la rigueur et pour qu'un joueur retrouve un mode sous le nom qu'il
connaît.

## Le dex

- [ ] dans les jeux, mettre en scène la silhouette : « ce pokémon
      s'appelle phrygien bécarre 6, tu veux le capturer ? » Attend le
      pont `analysis` → `FactSounded`.
- [ ] `Cooling` : les intervalles de répétition espacée, jamais décidés.
      La bonne courbe pour un clavier se juge à l'oreille, pas dans une
      formule recopiée d'un logiciel de fiches.
- [ ] `Discoverable` : demande un recensement de toutes les notions, qui
      n'existe nulle part.
- [ ] `KindSystem` : une gamme est une entrée à part entière, avec ses
      propres marques, distincte des entrées de ses modes (l'ionien,
      l'éolien). Acté, pas écrit. La désignation existe :
      `harmony.NamedScale` désigne les cinq gammes nommées, comme
      `System` et un degré désignent un mode.
- [ ] distinguer le système comme gamme mère des modes et la gamme
      comme référence de l'harmonie tonale. La gamme mineure naturelle
      n'est pas un citoyen de seconde zone : dans un contexte tonal,
      c'est une gamme de référence au même titre que la majeure, même
      si c'est aussi le mode du sixième degré de celle-ci.

## L'harmonie

- [ ] le pont `analysis` → `FactSounded` : repérer les degrés
      caractéristiques d'un mode sur une fenêtre de jeu. Choix à faire
      sur la taille de la fenêtre et l'ambiguïté entre modes voisins.
      C'est ce qui fera exister les silhouettes et `FactChosen`.
- [ ] champ `Tetrad harmony.ChordPattern` dans le catalogue, extensions
      en motif, et `naming` réduit au rendu du chiffrage. Le catalogue
      actuel devient l'oracle du test plutôt que la donnée. Le métier ne
      doit pas parser une chaîne quand il a une représentation exacte
      sous la main.
- [ ] rendre un chiffrage depuis une lecture : `C7♯9` à partir d'une
      `Reading`. Écrit en lettres américaines dans toutes les langues,
      dit dans la langue (« do mineur majeur 7 add 9 ») : deux rendus,
      dont un seul dépend de la locale.
- [ ] catalogue de progressions, pour que `ProgressionID` désigne
      quelque chose. Servirait aussi à juger les détours d'une
      réharmonisation (voir les jeux). Le contenu v0 est dans
      `grilles.md` : cadences, cellules, préparations, variantes en
      squelette de degrés et provenance par degré.
- [ ] l'accord parallèle comme type de préparation
      (`harmony.ApproachKind`), le seul du tableau de `grilles.md` qui
      manque ; il demande la mélodie pour être sûr (voir « L'analyse
      des grilles » pour la marche d'accords parallèles, qui ne la
      demande pas).
- [ ] la ligne de basse chromatique sous d'autres accords que le
      diminué (dominantes chromatiques et renversements, It Never
      Entered My Mind) : une autre lecture, à côté de `PassingChords`.
- [ ] un registre parlé anglais, s'il en existe un qui mérite d'être
      proposé en alternative. Aujourd'hui seul le français en a un.

## Les voicings

Les décisions sont dans `voicings.md`.

- [ ] les types `Position` et réalisation dans `harmony`, et le calcul
      de la position d'une réalisation selon une lecture d'accord.
- [ ] la table de tessiture : partir des *low interval limits* de
      Berklee, puis la corriger d'après le cours « Les bases de
      l'harmonisation » d'Étienne Guéreau. Les annotations portent sur
      toutes les paires de voix, jamais de refus.
- [ ] le calcul du moindre mouvement, qui dérive les couples de
      positions d'une progression.
- [ ] dans le dex, un détail de production propre à la notion : le
      couple (tétrade, fondamentale) pour une position.
- [ ] vérifier la fin du II-V-I mineur en m6, 1-3-6-5.

## La lecture des grilles iReal

Lire les grilles iReal Pro pour qu'un pianiste travaille la sienne avec
le dex : suggestions d'harmonisation, revue de tous les avatars
cadentiels qu'il connaît sur chaque « change », voicings.

Fait, dans `charts/ireal` :

- **L'URL et la playlist**, les anciennes dispositions comprises, et le
  débrouillage.
- **Les jetons**, sans perte : joints, ils redonnent la grille, qualités
  tapées librement (`A*m7*`), blancs parasites et `T12` compris.
- **Les mesures** : barres, sections, métrique, fins, signes, chaque
  accord dans sa case (`LZ` une case, `Kcl` deux).
- **Le dépliage** en ordre de jeu : reprises (« 3x »), fins dans l'ordre
  écrit, D.C. et D.S. al Coda, al Fine et al Nth ending, écrits
  n'importe où dans le commentaire, « x » et « r ». Tolère les
  bizarreries de saisie du corpus et ne boucle jamais ; seul un « D.C.
  on cue » reste non suivi.
- **Les chiffrages** lus en accords de `harmony` : la table des 62
  qualités de l'app, relue et validée, et les qualités tapées à la main.
  Sur 60 740 accords réels, un seul reste illisible, une faute de
  frappe.
- **Des cases aux temps** (`Chart.Timeline`) : un accord commence sur
  le temps où tombe sa case, arrondi au temps suivant (Take Five,
  vérifié à l'écoute) ; la pulsation est la noire, ou la noire pointée
  en 6/8 et 12/8. La coda, jouée au dernier chorus seulement, est hors
  du cycle.
- **Les fiches de référence** du livre dans
  `charts/ireal/testdata/fiches` : Tune Up, Black Orpheus, There Will
  Never Be Another You, Tenderly, qui concordent mesure par mesure avec
  les grilles du corpus, à la tétrade près.

Les tests commités n'utilisent que des grilles fabriquées ; les exports
de l'app vont dans `charts/ireal/testdata/local/`, ignoré par git, où un
test les lit s'il y en a.

Reste :

- [ ] un format de grille ouvert, mieux conçu que celui d'iReal.

## L'orthographe entendue

La première cause d'analyses fausses, chez un humain, est une notation
erronée : un ♭VII de do écrit A♯7 envoie l'apprenti analyste chercher
un accord qui mène à si. « Tu ne suivras pas bêtement les indications
du Real Book », dit le deuxième commandement.

gohar n'est pas trompé par la graphie, parce que le noyau ne la lit
pas : il raisonne sur des hauteurs, et A♯7 s'analyse comme B♭7. Le
jour où l'on analysera un flux MIDI, qui ne porte aucune orthographe,
le problème sera le même : nommer des accords à partir de touches. Dans
les deux cas, l'orthographe est une **sortie** de l'analyse, pas une
entrée. On entend un ♭VII7, donc on écrit B♭7. Sur une grille, cela
donne un correcteur : marquer les graphies qui s'écartent de ce que
gohar entend, et proposer la bonne, pour apprendre à l'apprenant à
faire de même.

Deux niveaux sont acquis par construction :

- **La fondamentale.** `naming` écrit dans un contexte : en si♭ majeur,
  il n'écrira pas ré♯, sauf choix délibéré de l'analyse, comme le
  ♯IIdim7 d'une basse qui monte.
- **Les altérations.** Le constructeur des `ChordPattern` normalise :
  un accord noté ♯9, ♯11, ♯5 ressort en (♭5, ♭10, ♭13). Proposer des
  couleurs est une autre affaire, et leur notation suivra les mêmes
  règles.

`analyse` écrit désormais chaque grille ainsi, et la graphie de l'app
ne sert plus à rien d'autre qu'à la lecture des accords. Les règles,
qui sont les nôtres, par ordre de priorité :

- **P0, des mouvements lisibles.** Une basse qui bouge d'une quinte
  bouge d'une quinte juste, d'une tierce d'une tierce, d'un pas d'une
  seconde. Un demi-ton peut aussi rester sur sa lettre, chromatique, et
  une note étrangère à la zone suit le sens de la ligne : dièse en
  montant, bémol en descendant (E7 E♭7 D7 D♭7 C, C C♯dim7 Dm7).
- **P1, le nom suit le chiffrage.** Une fondamentale prend la lettre
  de son degré dans la zone, le fond où se comptent les degrés : un
  ♭II se nomme comme un ♭II. Une fondamentale diatonique est le cadre
  où se lisent les mouvements et garde sa lettre ; P0 décide des
  autres, les fondamentales chromatiques et les accords qui préparent
  une résolution, nommés depuis leur cible. P1 ne fait alors que
  départager.
- **P2, la grille la plus simple.** La tonique de chaque zone s'écrit
  avec le moins d'altérations à la clef, ou comme son enharmonique si
  une armure peut l'écrire, selon ce qui donne le moins de smells,
  puis le moins d'altérations. *Crepuscule With Nellie* passe ainsi de
  sol♯ mineur à la♭ mineur.
- **P3, une tonicisation se nomme comme l'accord qu'elle tonicise** :
  [F♯m] sous F♯m7.
- **Une basse tenue**, une pédale, s'écrit une fois par son degré dans
  la zone et garde ce nom tant qu'elle est tenue.

**L'usage simplifie ensuite les noms affichés**, jamais le chiffrage ni
les mouvements :

- E♯, F♭, B♯ et C♭ passent sur la lettre voisine, sauf la sensible
  haussée d'un mineur ;
- une double altération aussi, et c'est une **enharmonie tolérée**,
  que `analyse -smells` liste à part : « D♯7/G for D♯7/F𝄪 bar 12 » ;
- un accord qui divise l'octave en parts égales, dim7 ou triade
  augmentée, posé sur une de ses notes, s'écrit sur sa basse : Cdim7/A
  devient Adim7, C+/E devient E+.

Une fondamentale libre, chromatique ou qui prépare une résolution,
n'est d'ailleurs jamais écrite autrement qu'à l'usage.

**Le smell** est une fondamentale ou une basse affichée avec une double
altération, ou en E♯, F♭, B♯, C♭ que l'usage garde. Il compte une fois
par accord, et `analyse -smells` en donne le total et la liste. Le
choix de la zone (P2) compte les smells sur les noms affichés, et les
altérations sur l'écriture par degré : do♯ majeur affiche son mi♯ en
fa, mais écrit toujours sept dièses. Il désigne d'ordinaire un
endroit où l'analyse s'est trompée de zone ou de chiffrage, ou une
chaîne qui a fait le tour des tonalités.

Reste **le mode** : sur une grille modale, écrire les accords de façon
à rendre compte des modes entendus. *Nardis* est le cas d'école, parce
que ses tétrades disent tout : Fmaj7 donne la ♭2, B7 pose mi, donc mi
phrygien, donc do majeur pour gamme-mère hors des dominantes, donc les
extensions à proposer pour le faire sonner. C'est le terrain des
cadences modales à deux accords (voir « L'analyse des grilles »).

La limite, ce sont les grilles où l'information manque : les versions
mal fichues d'*Infant Eyes* ou de *Naima* (locrien, phrygien contre
lydien, pédales) ne se liront pas de but en blanc. Pour ces deux-là, on
dispose de grilles propres, à l'orthographe corrigée par un professeur,
et de leurs voicings : elles serviront de référence, la grille
fautive en entrée, la propre en sortie attendue.

- [x] la fondamentale et la basse des accords orthographiées par degré
      dans la tonalité d'arrivée, puis simplifiées par l'usage, dans
      `analyse`.
- [x] l'endroit où rompre une chaîne qui ne se referme pas sur les
      lettres : à la frontière des zones, là où le chiffrage change de
      tonique. *Lush Life* et *Yesterday's Gardenias* ne sentent plus.
      La rupture G♭m7 → F♯m7 de *Grand Central*, mesure 17, tombe à
      cette frontière : ♭IIm7 en fa mineur, puis Im7 en fa♯ mineur.
- [x] le rendu des symboles d'accords entiers dans `naming`, la qualité
      comprise (`ChordStyle`), et ses réglages dans `analyse`.
- [ ] le signalement, en simple remarque, des graphies qui s'écartent
      de ce qui est entendu : « écrit A♯7, entendu ♭VII7 ».
- [ ] la lecture du mode quand les tétrades suffisent, *Nardis* en
      tête, et les extensions qui en découlent.
- [ ] un verdict de confiance sur les degrés et les modes, comme celui
      de la tonalité contre l'app, sur des grilles dont on connaît la
      bonne lecture, *Infant Eyes* et *Naima* comprises.
- [ ] la proposition de graphie, quand le verdict le permet.

## L'analyse des grilles

Conçue dans `grilles.md`, qui en tient les règles et les décisions,
et codée dans `harmony/analysis`. Le but du moment : charger une grille
iReal et en afficher une analyse exemplaire au sens d'*En Harmonie*,
de quoi la montrer à ses auteurs.

Fait :

- **La suite d'accords** (`Changes` : basses, durées, bouclage, coda)
  et son pont depuis iReal.
- **Les préparations** d'un accord au suivant (`Approaches`) et **le
  passage** (`PassingChords`), nommé d'après sa basse.
- **Les blocs**, [II] [sus4] V → cible, et les tonalités qu'ils
  annoncent en gammes précises (`Blocks`, `Block.Announced`).
- **Les degrés**, en deux lectures simultanées : sur la tonique
  installée (`Degrees`) et en crochets comme le livre (`Bracketed`).
- **La tonique pressentie**, lue de gauche à droite (`Sense`) : tonique
  de fond, tonique locale, ce qu'une cadence attend, le I emprunté.
- **La modulation**, installée en direct (`Sense`) et datée après coup
  (`Grounds`), le retour à la maison sur une cadence ou sur le seul
  accord de tonique ; Tune Up et Black Orpheus concordent entièrement
  avec leurs fiches.
- **Les phrases** (`Phrases`), qui vont d'un repos au suivant ; **la
  maison** (`Home`), là où se pose la première, et **la tonalité du
  morceau** (`Tune`), là où il s'arrête à son dernier chorus : la fin
  que marque la grille, sinon la conclusion de sa dernière section,
  turnaround exclu, sinon à travers la boucle ; jamais l'armure ; une grille qui boucle
  entendue comme son deuxième chorus ; la tonique mineure écrite m7 ;
  la tierce picarde (`Picardy`).
- **Les cadences plagales** : le IV de toute qualité, et le ♭VII7 avec
  son IVm7 (`PlagalApproach`), qui concluent sans ouvrir de tonique.
- **Le blues**, reconnu à sa forme (`Blues`), sa septième d'espèce lue
  comme sa tonique.
- **La plage modale** (`Modal`), un accord tenu quatre mesures sans
  cadence, reconnue sans que son mode soit lu : la grille ne le dit
  pas. Le corpus met les morceaux modaux à part (`IsModal`).
- **L'affichage** dans le terminal (`charts/cmd/analyse`) et **le
  rapport sur le corpus** (`charts/cmd/corpus`).

La suite, dans l'ordre :

- [ ] la modalité, sur de vraies grilles modales qui portent leurs
      couleurs : ce que devient un accord dans une plage, et le format
      de grille qui les écrit.
      *Nardis* en est un cas : mi phrygien à l'oreille, le Fmaj7 en
      donnant la ♭2 (son degré caractéristique) et le B7 posant mi comme
      tonique. L'analyse entend do, sur les arrêts et les retours à
      Cmaj7.
      Le livre en pose les repères (tome 2, chapitre 2 §5) : « la
      cadence modale a pour but d'installer la couleur d'un mode, alors
      que la cadence tonale va installer ou confirmer une tonalité », et
      une plage modale est un passage où « aucun mouvement harmonique
      n'est rencontré », de longueur libre (*So What*, *Flamenco
      Sketches*). D'autres cas à regarder : la modulation parallèle,
      qui change de mode sans changer de tonique (*On Green Dolphin
      Street*, E♭maj7 E♭m7 F7/E♭ Emaj7/E♭), et l'ostinato dorien de
      l'introduction de *Stolen Moments* (Cm7 Dm7/C E♭maj7/C sur une
      pédale de do).
- [ ] les cadences modales du tableau du livre (tome 2, p. 43), à deux
      accords, II-I et VII-I pour chaque mode, en do : D7 Cmaj7 et Bm7
      Cmaj7 (lydien), D♭maj7 Cm7 et B♭m7 Cm7 (phrygien), Dm7 Cm7 et
      B♭maj7 Cm7 (dorien), Dm7 C7 et B♭maj7 C7 (mixolydien), Dm7♭5 Cm7
      et B♭7 Cm7 (éolien). Aujourd'hui seul le ♭VII7-I se lit, comme
      une plagale. Le Fmaj7 Em7 de *Nardis* est le ♭IImaj7-Im7
      phrygien. Une cadence modale « doit obligatoirement faire
      entendre les DCN et DCA du mode », ce que les deux accords font
      d'eux-mêmes (le fa♯ de D7 sur Cmaj7). Reste à savoir quand deux
      accords sont une cadence modale plutôt qu'un mouvement conjoint
      de grille tonale (Dm7 Cm7 dans un II-V de si♭).
- [ ] le IIm7♭5 → I : dans *I'm Old Fashioned*, Gm7♭5 Fmaj7 « peut
      être perçu comme une cadence plagale mineure, cet accord se
      confondant avec un B♭m6 (IVm) » (tome 2, p. 34). Vérifier ce que
      les préparations en font.
- [ ] les cadences du catalogue qui restent, chacune avec ce qui lui
      manque : le ♭VIImaj7-I et le ♭VIIm7-I (les modes de la grille). Le V seul, sans II, qui va
      ailleurs que sur sa tonique n'est pas un bloc : les degrés disent
      déjà V III ou V IV7, et rien ne demande plus.
- [x] les II-V consécutifs (`Links`) : ½ ton, ton, cycle des quintes,
  d'après *En Harmonie* (tome 1, chapitre 8 §3.3), et les chaînes de
  dominantes (par quintes ou par demi-tons). Les relations entre blocs
  sont faites.
- [x] le V qui revient sur son II (Dm7 G7 | Dm7 G7 dans *Satin Doll*),
  lu seul comme une plagale IV7 → m, se relit sur la tonique entendue
  (`Reread`, une seconde écoute) : un II-V rejoué, sauf quand le mineur
  est la tonique (le I IV7 dorien de *Mas Que Nada*, de *It Ain't
  Necessarily So*) ou que l'alternance dure plus de quatre mesures (un
  vamp). 564 blocs relus dans 237 standards, sans changer la tonalité
  d'aucun.
- [ ] les accords parallèles (*Stolen Moments*) : au moins trois
  accords de même tétrade, à intervalle variable ; à intervalle
  constant, une marche d'accords parallèles. Reste à décider ce qu'on
  en montre quand les blocs lisent déjà les accords (voir
  `grilles.md`).
- [ ] distinguer la cadence évitée, où la tonalité change, de la
  rompue. Sans urgence : « … » suffit pour l'instant, la ligne des
  toniques montrant si la nouvelle tonique s'installe.
- [x] les cellules (`Cells`) : l'anatole et le III-VI-II-V-I, d'après
  *En Harmonie* (tome 1, chapitres 8 et 9).
- [x] la cadence rompue V-VI (G7 Am7 en do), lue seule comme un
  ♭VII7-Im de la mineur, se relit de même sur la tonique entendue au V
  (`Reread`).
- [x] la demi-cadence (`Conclusion.Half`) : une section qui ne conclut
  pas et s'arrête sur le V d'un bloc, d'après *En Harmonie* (tome 1,
  chapitre 8). Le contraste entre deux fins d'une même section a été
  écarté (voir « Ce qu'on a décidé » des cadences dans `grilles.md`).
  Le turnaround, qui se reconnaît à sa place dans la structure, reste
  à faire ; ce n'est pas une cellule (une anatole, souvent, ou un
  turnaround sur pédale en introduction).
- [x] les cellules par substitutions tritoniques, lues sur les accords
  qu'elles remplacent, comme *En Harmonie* retrouve les cadences
  initiales (tome 1, chapitre 9) : « anatole ♭II », nom retenu en
  attendant mieux. Au passage, un accord majeur n'est plus un II : Em7
  A7 Dmaj7 G7 est un II-V-I de ré (327 fausses cellules en moins sur
  le corpus).
- [ ] *Yesterday's Gardenias* est entendu en fa♯, et non en si♭ : la
  grille finit sur F♯maj7, mesure 32. À regarder.
- [ ] *Peace* (Horace Silver) : le livre lui donne un « centre tonal
  autour de Si♭ », l'analyse entend ré♭. Les autres fiches du tome 2
  concordent (*Body And Soul* compris, voir « Les autres cas tranchés »
  dans `grilles.md`), et *Fall*, sans centre tonal, est écarté.
- [x] les pédales à la basse (`Pedals`) : de tonique, de dominante ou
  sur un autre degré, générales ou passagères, notées comme le livre
  (« B♭ ped. ») entre les accords et les degrés (voir « La pédale »
  dans `grilles.md`). La double pédale n'est pas reprise : c'est un
  conseil d'arrangement, sans objet sur une grille d'accords seuls.
- [x] le I renversé sur sa quinte sur lequel un V se résout reste le I
  (d'après *En Harmonie*) : le livre prolonge la pédale de
  dominante sur l'accord de tonique « entendu renversé sur sa 5te »,
  Fm9/B♭ B♭7 E♭maj9/B♭ (tome 2, §1.5), Dm9/G G7 C6/9/G dans *My
  Romance*, que les tests reprennent (`fifthTonic`). Essayé :
  admettre tout accord renversé sur sa quinte comme tonique gagne 8
  grilles du corpus (*The Look Of Love*, *Sail Away*, *Re: Person I
  Knew*…) et en perd 6 : un IVm sur pédale de tonique devient une
  tonique (A♭m/E♭ dans *I See Your Face Before Me*). Il faut s'en
  tenir au I sur lequel le V se résout.
- [ ] la modulation « confirmée », et les seuils en données.
- [ ] l'ouverture sur un maj7 qui n'est pas la tonique : *Only Trust
  Your Heart* s'ouvre sur Fmaj7♯11, le IV lydien de do, et l'analyse
  part de fa avant d'entendre do. Deux pistes, à ne pas trancher trop
  tôt pour ne pas casser d'autres grilles : le ♯11 écrit dit lydien,
  pas de tonique ; ou un retour à l'ouverture ne fait un repos que si
  elle est tenue plus d'une mesure. À vérifier sur les familles quarte
  et relatif du corpus. Sur les 1350 standards, neuf s'ouvrent sur un
  maj7♯11 : trois l'ont pour tonique (*Jackie-ing*, *Zoltan*, *Afro
  Centric*), trois pour IV (*Only Trust Your Heart*, *Jinrikisha*,
  *Leaving*), trois autre chose (*Spain*, où c'est le ♭VI de si
  mineur). Le ♯11 ne tranche donc pas. Le rythme harmonique non plus :
  le turnaround Gm7 C7, deux temps chacun, prépare le Fmaj7♯11 d'une
  mesure, deux fois plus long. Reste la mélodie.
  Le livre le confirme : le lydien « peut tout à fait être utilisé sur
  l'accord du Ier degré » et « est fréquemment utilisé pour conclure un
  thème » (tome 2, p. 33 et 42 : *Make Someone Happy*, *Night and
  Day*).
- [ ] le rythme harmonique : « posé » est relatif à la densité des
  changements. Une mesure de Fmaj7 suivie d'une mesure d'autre chose
  n'est pas posée là où les accords durent une mesure ; elle l'est
  là où ils changent à chaque temps. Les seuils d'`held` et de la
  modulation (« plus d'une mesure ») comptent en mesures absolues, et
  devraient se mesurer au pas harmonique du passage.
  Essayé : « deux fois plus long que la préparation » seul pose le
  Fmaj7 d'une mesure de *Stella By Starlight* (mesure 13, après B♭m7
  E♭7) et perd 12 grilles. Conjugué à la carrure (fin d'un groupe de
  quatre depuis le début de la section), il n'en perd plus que 5, mais
  la carrure pèse trop : elle pose l'E♭maj7 de la mesure 4 de *Jordu*,
  alors que c'est Cm6, dès la mesure 2, qui fait référence. En
  attente : ni l'un ni l'autre ne fait mieux que le seuil actuel.
- [x] la tonalité d'analyse au choix : `Sense` la reçoit, `analyse`
  compte par défaut dans celle qu'il entend, `-key declared` dans
  celle de l'app, `-key F` dans celle qu'on impose.
- [x] le chiffrage en mineur, compté dans la gamme mineure comme le
  fait *En Harmonie* (tome 1, chapitre 8), une septième de dominante
  hors du V écrite avec son 7 (VII7), et le 7alt écrit 7alt.
- [x] ce que le verdict pèse (`ReadTune`), affiché par `analyse` : le
  premier accord, le dernier, et la durée de leurs toniques quand ils
  divergent (voir « Ce que le verdict pèse » dans `grilles.md`). Il
  remplace le comptage des preuves (`Candidates`), qui ne décidait pas.
- [ ] la mémoire, dernier critère de Siron pour la modulation vraie, et
  ce que la durée seule ne tranche pas : *Love Me Or Leave Me*, une
  phrase sur deux en fa mineur, l'autre en la♭, s'entend en la♭. Les
  témoins à écouter d'abord : sur le corpus, 43 grilles ont un premier
  et un dernier accord qui divergent, et 7 se jouent à moins d'un contre
  deux en durée. Pour chacune, laquelle des deux tonalités on entend,
  et ce qui l'emporte, la première entendue ou l'arrivée :

  | Morceau | Premier | Dernier | Analyse |
  |---|---|---|---|
  | *Lullaby Of Birdland* | fa mineur, 17 mesures | la♭, 15 | fa mineur, confirmé à l'oreille |
  | *Love Me Or Leave Me* | fa mineur, 14 | la♭, 7 | fa mineur ; l'oreille dit la♭ |
  | *How Deep Is The Ocean* | do mineur, 18,5 | mi♭, 10 | do mineur ; l'app dit mi♭ |
  | *So In Love* | fa mineur, 46 | la♭, 23 | fa mineur |
  | *Bess You Is My Woman* | si♭, 28,5 | ré mineur, 23,5 | si♭ |
  | *Glad To Be Unhappy* | sol mineur, 8 | fa, 16 | fa |
  | *The Summer Wind* | fa, 14,5 | la, 10 | fa |

- [ ] *Somewhere* : sa grille ne donne la tonique qu'en mi♭ mineur au
  pont et en mi♭ sur la fin, et l'analyse y entend une tierce picarde ;
  à trancher à l'oreille.
- [ ] `Home` et `homeAt`, la « maison » de l'ancien modèle, qui
  installent encore le fond à la première phrase conclusive dans
  `Sense` et servent de repli quand rien ne conclut : à retirer quand
  des témoins diront ce qu'ils protègent. Le seul connu est *Just
  Friends*, qui s'ouvre sur Cmaj7, son IV, et ne se pose qu'à la fin
  sur G6 : à écouter.
- [ ] les couleurs proposées, à rebrancher sur `Sensed` : ce qu'on dit à
  l'apprenant sur une tonicisation (la gamme du fond, celle de l'accord
  tonicisé, les deux ?), sur une région transitoire, sur un II-V qui ne
  résout pas. Le premier livrable d'un outil d'analyse pour apprenants.
- [ ] la forme donnée plutôt que détectée : des sections passées à
  `Sense` quand elles sont connues (un standard AABA de 32 mesures),
  pour le direct comme pour les grilles dont la forme se sait.
- [ ] la fiche *Django* (Siron 5.11.4), laissée de côté faute de
  savoir lire son niveau de modulation.
- [ ] le niveau de jeu tiré d'une grille : quelles mécaniques d'abord,
  parmi ce que l'analyse sait (les cellules à reconnaître, la cadence
  qui arrive, la tonique qui bouge, la section qui se referme).
- [ ] les seuils à nous, à régler de façon empirique, aucune source ne
  les chiffrant : deux mesures pour une zone, deux crans pour un centre
  éloigné, la moitié d'une section pour une modulation vraie, la place
  des régions dans les sections, une mesure pour un appui.
- [ ] la mélodie : ce qui sépare *In a Sentimental Mood* (ré mineur)
  de *Lullaby Of Birdland*, au même profil, et pose d'entrée sol
  mineur dans *It Don't Mean A Thing*.
  Demande un format de grille qui porte la mélodie.
- [ ] la jauge de tension, puis le direct avec l'attente et la
      surprise (voir `grilles.md`), dont le pivot diminué de Tenderly
      comme test à l'envers.
      La pédale de dominante en est une source : elle « crée une
      tension qui ne trouvera sa résolution » qu'au retour de la
      tonique (*En Harmonie*, tome 2, chapitre 5).
- [ ] la grille annotée dans une fenêtre Ebitengine (police de Real
      Book, chiffrages en indices et exposants, réglable), une fois le
      cœur validé.
- [ ] le moteur d'analyse en WASM dans une page web, pour distribuer et
      faire connaître le travail. Le moins prioritaire, à ne pas
      perdre de vue.
- [ ] les voicings sur une grille : la marque « Employée » du dex
      constate une position placée spontanément sur les changes.
- [ ] une grille iReal comme niveau du shoot'em up.

## En attente d'un cas

Ce qui n'avance que le jour où une grille, un jeu ou un besoin le
demande. Ce ne sont pas des chantiers : rien n'y est commencé, et rien
ne presse. Une entrée remonte dans sa section quand son cas se présente.

- passe ligne à ligne sur `Mode.Function`, le jour où les modes
  serviront à autre chose qu'à la reconnaissance à l'oreille. La
  fonction d'un mode dit comment on l'emploie, et peut différer de
  celle de sa tétrade : le dorien est une tonique, son m7 seul un
  II. Le premier degré du mineur harmonique (éolien ♮7), celui du
  majeur harmonique (ionien ♭6) et celui du majeur double
  harmonique (ionien ♭2 ♭6, « Xmaj7 ou ♭II7/I » dans *En
  Harmonie*, tome 2) portent `Tonic | Dominant` : ce sont aussi des
  avatars de dominante sur pédale de tonique. Seul ce double emploi
  est verrouillé par un test.
- un D.S. écrit dans une reprise pas encore terminée, s'il s'en
  présente un.
- `irealbook://`, l'ancien schéma non brouillé : refusé tant qu'on
  n'en a pas vu un vrai.
- un parseur de chiffrages général dans `naming`, pour ce qu'on
  tape soi-même, quand un jeu en aura besoin.
- afficher les alternatives d'un mode (alias, registre parlé) le
  jour où l'activité quitte le système naturel, qui n'en a pas.
- les doigtés, le jour où l'on travaillera les mains : des règles
  simples pour le cas général, et les exceptions en données
  d'expert.
- les dominantes sur pédale de tonique, à reconnaître quand un
  m(maj7) est chiffré avec des extensions qui décrivent les modes
  correspondants. Proche de la lecture des modes par les tétrades
  (*Nardis*).

## Les sources à dépouiller

Sources à consulter quand une question les touche, plutôt que de
trancher sans source.

- [x] *En Harmonie*, tome 1, chapitres 8 à 10 : Le chiffrage en
  mineur y est tranché (chapitre 8), et le X7sus4 comme accord de
  sous-dominante (chapitre 9, « Modifier un enchaînement »).
- [x] *En Harmonie*, tome 1, chapitre 8 §3.3, « Les enchaînements
  harmoniques fréquents » : les II-V consécutifs, codés (`Links`). La
  suite du chapitre (anatole, III-VI-II-V-I, p. 116 et 117) aussi
  (`Cells`).
- [x] Étienne Guéreau, cours « Les bases de l'harmonisation » : relu,
  reporté dans `voicings.md` sans en reprendre le contenu (un cours
  payant, cité comme exemple). Ses limites d'intervalles servent à
  vérifier la table de tessiture, à l'oreille.
- [x] Les dix commandements de l'harmoniste de la BEPA, dus à Étienne
  Guéreau, à citer librement là où ils appuient une décision :
  1. Tu partiras de la mélodie.
  2. Tu ne suivras pas bêtement les indications du Real Book.
  3. Tu alterneras les positions et les renversements.
  4. Tu varieras les nuances et les registres.
  5. Tu veilleras au toucher et à l'équilibre sonore.
  6. Tu ne réharmoniseras pas comme un sauvage.
  7. Tu utiliseras différentes techniques de réharmonisation.
  8. Tu ne feras pas étalage de ton savoir au détriment du thème.
  9. Tu seras attentif à la cohérence de l'ensemble.
  10. Tu transposeras.

  Le premier, le deuxième, le sixième et le huitième sont cités dans
  `grilles.md`, le troisième et le dixième dans `voicings.md`.
- [x] Jacques Chailley, *40 000 ans de musique* (dans le projet) : une
  histoire culturelle, sans règle d'harmonie ; rien sur la cadence, la
  tonique ni la modulation, et les notions que Siron lui emprunte
  (sensibilisation, naturalisation) sont dans son *Traité historique
  d'analyse harmonique* et *Expliquer l'harmonie ?*, à se procurer.
  Cité dans `grilles.md` pour les feuilles bleues (p. 263), le morceau
  « éprouvé » par l'oreille (p. 137) et l'extension de la consonance
  (p. 151-154) ; dans `oreille.md` pour Royaumont (p. 162-163) ; ici
  pour l'improvisateur (p. 261-263, 306).
- [x] *En Harmonie*, tome 2, chapitre 5, « Les pédales » jusqu'à
  l'ostinato : la pédale simple et double, de tonique et de dominante,
  générale ou passagère, sa notation (« X ped. » plutôt que la barre
  oblique), le sus4 comme sous-dominante sur pédale de dominante
  (§1.5), la fonction des accords préservée sur la pédale (§1.6), le
  turnaround sur pédale (§2.5), l'ostinato distingué de la pédale. Ce
  qu'on en tire est reporté dans les chantiers de l'analyse.
- [x] *En Harmonie*, tome 2, chapitre 2 §5, « Les modes naturels dans
  les thèmes » : chaque mode dans des thèmes, les cadences modales
  (§5.2), les plages modales (§5.3), les deux approches de la modalité
  (§5.4), la rencontre des cadences modales et tonales (§5.5). Ce
  qu'on en tire est reporté dans les chantiers de l'analyse.
- [x] *En Harmonie*, tome 2, « Récapitulatif des modes » et « Analyses
  modales » : les tétracordes, le rétrograde inversé de chaque mode,
  le tableau des cinq systèmes (tétrade, extensions, degrés
  caractéristiques), et des fiches (tonalité, forme, un mode par
  accord) pour *Someday My Prince Will Come*, *The Days Of Wine And
  Roses*, *Body And Soul*, *Fall*, *Very Early*, *Peace* et *Re:
  Person I Knew*. Le catalogue des modes concorde avec le tableau, à
  deux écarts près, notés dans « L'harmonie ».
- [ ] *En Harmonie*, tome 2, chapitre 6, l'accord de dominante sur
  tonique (la double fonction de l'ionien ♭6 et de l'éolien ♮7) et
  l'accord appoggiaturé : reconnaissance et chiffrage.
- [ ] Jacques Siron, *La partition intérieure* : le rythme harmonique
  et la carrure (ce qui fait qu'une harmonie est posée), les formes et
  leurs sections (la demi-cadence en fin de section), le rapport de la
  mélodie aux accords, le blues et ses variantes. À vérifier sur le
  livre, d'après ce qu'on en connaît.
- [ ] Philippe Baudoin, *Jazz mode d'emploi* : les blues que `Blues`
  ne connaît pas encore (*Freddie Freeloader*, *Doxy*, *Watermelon
  Man*), et les réharmonisations, pour le catalogue des cellules.

## L'audio

- [ ] les allocations de la boucle de jeu : une collecte toutes les une à
      deux secondes (voir `architecture.md`, tranché le 26/09). Sans
      risque pour le son, à réduire quand on touche au rendu :
      `reveal()` refait ses namers à chaque frame, par exemple.
- [ ] calibration chez le joueur. Promise dès le premier jour, et les
      chiffres mesurés ici sont ceux d'une machine, pas une promesse.
- [ ] les timbres 8 bits dans les préférences du joueur, avec le rendu
      authentique ou adouci : pour l'instant `-timbre` et `-authentic`,
      le rendu adouci étant le défaut.
- [ ] niveau de sortie bas, conséquence de la marge prise sur le gain.
      Réglage, pas conception.

Deux portes délibérément ouvertes et non planifiées, décrites dans
`architecture.md` : un lecteur `/dev/snd/midiC*D*` en pur Go pour sortir
du C++, et la cible navigateur. Les garder ouvertes ne coûte qu'une
chose, maintenir la règle des deux surfaces : oto n'est importé que par
`synth/device.go`, gomidi que par `games/keyboard/midi.go`.

## Les jeux

- [ ] *Walk with me* : remplacer le bassiste à la main gauche, poser
      des pêches à la main droite, à vue sur une grille arbitraire. Le
      document de conception est dans `walk.md`. Le premier jalon : les
      fondamentales, main gauche seule, sur un blues en fa, au
      métronome, avec la phase sans tempo. Il demande au synthé des
      notes programmées à l'échantillon près et un mélangeur, et un
      `MetronomeSystem`. À côté : dépouiller *Jazz Piano Fundamentals*
      de Jeremy Siskind, et trancher le modèle de paliers.

La suite de `ear`, dans l'ordre de `oreille.md` :

- [ ] un rendu plus joli d'une touche enfoncée : l'enfoncement de 2 px
      passe pour un MVP mais a l'air bon marché.
- [ ] la soundfont, dans `synth/soundfont` : go-meltysynth (MIT, rien
      d'autre que la bibliothèque standard, n'alloue pas en rendu, SF2
      seulement) enveloppé dans la `queue` commune. Il faut un piano
      réduit à quelques Mo, sous licence claire (Salamander, CC-BY),
      préparé avec Polyphone. Puis nouvelle mesure sous charge.
- [ ] la feuille de route de `oreille.md` : l'abstraction, le menu, les
      degrés aux boutons et les tétracordes sont faits ; restent la
      réponse jouée, les réglages et les niveaux paramétrables.
- [ ] la réponse jouée devra régler ce que le piano montre : les touches
      qui s'allument sous la note de la question donnent la réponse. Ce
      que joue la séquence ne s'allume pas pendant la question, ce que
      joue le joueur, si.
- [ ] WASM : `midi.go` derrière un build tag, `webmididrv` plus tard,
      le dex dans `localStorage`. Le geste qui démarre l'audio est le
      clic dans le menu. La police des signes est embarquée, rien à
      faire de ce côté.
- [ ] juger à l'oreille le tempo (350 ms par note) et le registre (la
      gamme entre la3 et sol♯4 ; la pédale deux octaves dessous, sauf
      pendant la question d'un degré où elle est dans la même octave).
- [ ] les paliers suivants des modes : les autres systèmes. Les
      distracteurs proches sont couverts par le palier des sept modes.
- [ ] les paliers suivants des degrés : la tonique mobile, la tonique à
      la basse, d'autres gammes, et au plus difficile l'échelle
      chromatique.

Les autres jeux :

- [ ] le shoot'em up bullet hell qui est un jeu d'harmonie déguisé, sur
      rail, sans esquive.
- [ ] les quatre pistes du billet sur la game loop de l'improvisateur,
      triées en jeux distincts plutôt qu'empilées dans un seul. Le
      billet a sa phrase d'ouverture chez Chailley : « Une "composition"
      n'est qu'une improvisation qui, jugée particulièrement réussie, a
      été fixée dans la mémoire », et l'interprète d'avant le papier,
      « un créateur qui reçoit une tradition et la modèle à sa guise »,
      « c'est encore ce que font nos orchestres de jazz » (*40 000 ans
      de musique*, p. 261-263) ; le jazz, « prééminence des interprètes,
      improvisateurs sur canevas » (p. 306). Le canevas, c'est la
      grille : le level design d'un niveau tiré d'une grille revient à
      de la composition.

Pistes ouvertes sur les cibles, sans urgence et pas encore assez nettes
pour en faire des règles. Elles relèvent du jeu, jamais du dex.

- [ ] **cible fonctionnelle.** Une cible n'est pas forcément une
      tétrade : en arrangement ou en réharmonisation, c'est souvent une
      fonction. La cible devient un empilement de contraintes
      optionnelles (fonction, puis tétrade, puis mode), et est ouvert
      tout ce qu'elle ne nomme pas. Une fonction est relative : « la
      dominante de quelque chose », donc la cible porte une tonalité ou
      un degré de référence. Voir ce que `Target.Resolve` porte déjà.
- [ ] **fondamentales admises pour une fonction.** Un slot « dominante »
      accepte d'office la dominante chromatique (D♭7 pour aller en C),
      et sans doute l'accord diminué 7 sur la sensible. Notion distincte
      de `Mode.Function`, qui dit seulement qu'un mode peut tenir un
      rôle. Si la fonction devient un critère de validité, la passe sur
      `Mode.Function` devient bloquante. Le livre répond pour
      l'essentiel : V7, ♭II7 et VIIdim7 pour la dominante, et la liste
      des sous-dominantes (voir `grilles.md`).
- [ ] **choix de tétrade comme fait.** Sous un slot purement
      fonctionnel, jouer D♭7 plutôt que G7 est un choix d'harmonisation,
      comme un choix de couleur. Faut-il un fait pour ça, et les
      substitutions sont-elles des notions du dex ?
- [ ] **ancres et détours.** Une réharmonisation peut ignorer
      délibérément des cibles, à trois conditions : une logique
      audible, rien qui coince avec la mélodie, et une résolution sur
      une cible attendue (tonique ou substitution de tonique). Les
      cibles se séparent alors en ancres, à atteindre, et en cibles
      indicatives, contournables. La mélodie et l'ancre se vérifient
      de façon déterministe. La logique audible est la difficile :
      piste préférée, un catalogue de progressions (dominantes en
      chaîne, dominantes chromatiques en cascade, marche de basse
      chromatique, structure constante…), un détour étant accepté s'il
      en suit une. Rejoint le catalogue de progressions de la section
      harmonie.

## Le dépôt

- [ ] tags préfixés le jour de la publication : `harmony/v0.1.0`,
      `synth/v0.1.0`.
