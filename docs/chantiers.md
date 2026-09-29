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
harmonique, mineure mélodique, majeure harmonique) sur leur tonique.
Chaque mode a un nom systématique et des alternatives : le registre
parlé français et les alias.

`analysis` identifie les accords de façon déterministe, sans pondération,
et analyse une grille à la manière d'*En Harmonie* : préparations,
passages, blocs et tonalités qu'ils annoncent, degrés sur la tonique
installée et en crochets, tonique pressentie et modulations, blues.
Les fiches du livre concordent à 83 degrés sur 83, et l'analyse tombe
d'accord avec la tonalité déclarée par l'app, ou vérifiée à l'oreille,
sur 82 % du corpus (1272 des 1547 grilles jugées, les modales et une
liste relue à la main mises à part), sans la lire. Chaque tonique que
le morceau donne est un candidat avec ses preuves, que `analyse`
affiche : le verdict s'explique.

`dex` a son corps, sa persistance JSON et `Components`. Restent
`Cooling` et `Discoverable`.

`synth` joue : conversion en fréquence, moteur polyphonique sans
allocation, discipline des buffers audio, et un histogramme des délais
à cases fixes.

`charts/ireal` lit les grilles d'iReal Pro, de l'URL jusqu'aux accords
dans l'ordre de jeu : playlist, jetons sans perte, mesures, dépliage de
la forme, chiffrages lus en accords de `harmony`, durées en temps.
Vérifié sur 1678 grilles réelles, qui restent hors du dépôt, et sur
quatre fiches d'analyse du livre *En Harmonie*. Deux commandes :
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
- [ ] passe ligne à ligne sur `Mode.Function`. La règle est tranchée :
      la fonction d'un mode se dérive de sa tétrade, plus les fonctions
      qu'un expert ajoute. Le premier degré du mineur harmonique
      (éolien ♮7) et celui du majeur harmonique (ionien ♭6) portent
      ainsi `Tonic | Dominant` : ce sont aussi des avatars de dominante
      sur pédale de tonique. Reste à relire les ajouts d'expert mode par
      mode.
- [ ] champ `Tetrad harmony.ChordPattern` dans le catalogue, extensions
      en motif, et `naming` réduit au rendu du chiffrage. Le catalogue
      actuel devient l'oracle du test plutôt que la donnée. Le métier ne
      doit pas parser une chaîne quand il a une représentation exacte
      sous la main.
- [ ] test « un mode ne porte jamais moins de fonctions que sa tétrade ».
      Les ajouts d'expert n'enlèvent jamais un rôle.
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

- [ ] un D.S. écrit dans une reprise pas encore terminée, s'il s'en
      présente un.
- [ ] `irealbook://`, l'ancien schéma non brouillé : refusé tant qu'on
      n'en a pas vu un vrai.
- [ ] un parseur de chiffrages général dans `naming`, pour ce qu'on
      tape soi-même, quand un jeu en aura besoin.
- [ ] un format de grille ouvert, mieux conçu que celui d'iReal.

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
  morceau** (`Tune`), là où s'arrête la dernière, turnaround exclu, sauf quand le morceau s'ouvre au repos et que sa
  première cadence y revient ; jamais l'armure ; une grille qui boucle
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
- [ ] les cadences du catalogue qui restent, chacune avec ce qui lui
      manque : la demi-cadence (les fins de section), la rompue (la
      surprise), le V seul et ce qu'il fait attendre (la surprise), le
      ♭VIImaj7-I et le ♭VIIm7-I (les modes de la grille), IV7-I7 (le
      blues).
- [x] les II-V consécutifs (`Links`) : ½ ton, ton, cycle des quintes,
  d'après *En Harmonie* (tome 1, chapitre 8 §3.3), et les chaînes de
  dominantes (par quintes ou par demi-tons). Les relations entre blocs
  sont faites.
- [ ] le V qui revient sur son II (Dm7 G7 | Dm7 G7 dans *Satin Doll*)
  est lu comme une plagale IV7 → m, et annonce ré mineur : c'est un
  II-V rejoué. Mais la même ligne est le vamp dorien Im7 IV7 quand le
  m7 est la tonique (Fm7 B♭7 dans *Mas Que Nada*, Gm7 C7 dans *It
  Ain't Necessarily So*), et les accords seuls ne les séparent pas.
  Essayé : « un V entre deux fois son II n'est pas une plagale » gagne
  4 grilles et perd ces 2 ; exiger en plus que le V se résolve ailleurs
  fait pire. En attente des cadences modales du tome 2 d'*En Harmonie*
  (chapitre 2 §5.2).
- [ ] la marche d'accords parallèles (*Stolen Moments*).
- [x] les cellules (`Cells`) : l'anatole et le III-VI-II-V, d'après
  *En Harmonie* (tome 1, chapitres 8 et 9).
- [ ] la suite du catalogue des cellules (le turnaround, l'anatole
  réharmonisé par substitutions tritoniques, I ♭III7 ♭VI7 ♭II7).
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
- [x] la tonalité du morceau par accumulation de preuves (`Candidates`),
  lues et affichées, sans décider : le comptage confirme `Tune` plus
  qu'il ne le dépasse (voir « Le verdict et ses preuves » dans
  `grilles.md`). Sur des grilles seules, on ne fera guère mieux.
- [ ] la mélodie comme preuve : ce qui sépare *In a Sentimental Mood*
  (ré mineur) de *Lullaby Of Birdland* (la♭), au même profil de
  preuves, et pose d'entrée sol mineur dans *It Don't Mean A Thing*.
  Demande un format de grille qui porte la mélodie.
- [ ] la jauge de tension, puis le direct avec l'attente et la
      surprise (voir `grilles.md`), dont le pivot diminué de Tenderly
      comme test à l'envers.
- [ ] la grille annotée dans une fenêtre Ebitengine (police de Real
      Book, chiffrages en indices et exposants, réglable), une fois le
      cœur validé.
- [ ] le moteur d'analyse en WASM dans une page web, pour distribuer et
      faire connaître le travail. Le moins prioritaire, à ne pas
      perdre de vue.
- [ ] les voicings sur une grille : la marque « Employée » du dex
      constate une position placée spontanément sur les changes.
- [ ] une grille iReal comme niveau du shoot'em up.

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
- [ ] Jacques Chailley, *40 000 ans de musique* (dans le projet) :
  l'histoire, à citer pour la tierce picarde et ce que la doc en dit.
- [ ] *En Harmonie*, tome 2, chapitre 5, les pédales : la quarte et
  sixte, les X/Y sur pédale (Fmaj7/G, E♭maj7/F) traités au cas par cas,
  le sus4 et sa fonction (§1.5, §1.6), le turnaround sur pédale (§2.5).
- [ ] *En Harmonie*, tome 2, chapitre 2 §5 : cadences modales (le
  ♭VIImaj7-I écarté, Dm7 G7 Am7 lu éolien), plages modales, rencontre
  des cadences modales et tonales. La base sourcée du chantier modal.
- [ ] *En Harmonie*, tome 2, « Analyses modales » (p. 91) : des fiches
  de référence pour les thèmes modaux, comme celles du tome 1 pour le
  tonal.
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
- [ ] afficher les alternatives d'un mode (alias, registre parlé) le
      jour où l'activité quitte le système naturel, qui n'en a pas.

Les autres jeux :

- [ ] les doigtés, le jour où l'on travaillera les mains : des règles
      simples pour le cas général, et les exceptions en données
      d'expert.
- [ ] le shoot'em up bullet hell qui est un jeu d'harmonie déguisé, sur
      rail, sans esquive.
- [ ] les quatre pistes du billet sur la game loop de l'improvisateur,
      triées en jeux distincts plutôt qu'empilées dans un seul.

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
