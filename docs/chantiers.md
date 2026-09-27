# Chantiers

Ce qui est ouvert, et ce qui attend une décision.

Ce fichier existe pour qu'une conversation neuve reprenne sans rien
redécouvrir et sans rouvrir un débat déjà tranché. Les raisons des choix
faits sont dans `architecture.md`, le vocabulaire dans `glossaire.md`, la
conception de la collection dans `dex.md`, l'ear trainer dans
`oreille.md`, les voicings dans `voicings.md`, l'analyse des grilles
dans `grilles.md`. Ici il n'y a que ce qui reste à faire, et les
décisions qu'il ne faut pas rouvrir.

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

`analysis` identifie les accords de façon déterministe, sans pondération.

`dex` a son corps, sa persistance JSON et `Components`. Restent
`Cooling` et `Discoverable`.

`synth` joue : conversion en fréquence, moteur polyphonique sans
allocation, discipline des buffers audio, et un histogramme des délais
à cases fixes.

`charts/ireal` lit les grilles d'iReal Pro, de l'URL jusqu'aux accords
dans l'ordre de jeu : playlist, jetons sans perte, mesures, dépliage de
la forme, chiffrages lus en accords de `harmony`, durées en temps.
Vérifié sur 1678 grilles réelles, qui restent hors du dépôt, et sur
quatre fiches d'analyse du livre *En Harmonie*.

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
- [ ] l'inverse : lire un chiffrage écrit, pour les grilles.
- [ ] catalogue de progressions, pour que `ProgressionID` désigne
      quelque chose. Servirait aussi à juger les détours d'une
      réharmonisation (voir les jeux). Le contenu v0 est dans
      `grilles.md` : cadences, cellules, préparations, variantes en
      squelette de degrés et provenance par degré.
- [x] la provenance d'un accord : `harmony.Provenances`, toutes les
      gammes nommées qui le contiennent, sur toutes les toniques, avec
      le degré de sa fondamentale et le mode à jouer. Vérifiée sur les
      emprunts de Tenderly.
- [ ] les types de préparation de `grilles.md`, `harmony.ApproachKind`
      (*approach* en anglais) : dominante, dominante chromatique,
      diminué (dominante sans fondamentale), sus4 et II de… faits
      (`ApproachOf`, et les crochets II-V des fiches du livre vérifiés
      avec) ; restent les accords parallèles, en dernier.
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

## Les grilles

Lire les grilles iReal Pro pour qu'un pianiste travaille la sienne avec
le dex : suggestions d'harmonisation, revue de tous les avatars
cadentiels qu'il connaît sur chaque « change », voicings. Débroussaillé
le 27/09.

Fait :

- [x] **L'URL et la playlist** : les champs lus autour de la grille,
      ce qui couvre les anciennes dispositions ; le débrouillage.
- [x] **Les jetons**, sans perte : joints, ils redonnent la grille.
      Rustines trouvées sur les vraies grilles : les qualités tapées
      librement entre étoiles (`A*m7*`), les blancs parasites (un saut
      de ligne enregistré dans une grille), `T12` pour 12/8.
- [x] **Les mesures** : barres, sections, métrique, fins, signes, et
      chaque accord dans sa case. `LZ` porte une case et `Kcl` deux.
      Les cases vides de mise en page (après une reprise, après la
      barre finale) ne font pas de mesures.
- [x] **Le dépliage** en ordre de jeu : reprises (« 3x » compris), fins
      prises dans l'ordre écrit, D.C. et D.S. al Coda ou al Fine avec
      la dernière fin au retour, « x » et « r » résolus. Tolère les
      huit bizarreries de saisie du corpus (reprise jamais fermée, fins
      dans le désordre) et ne boucle jamais. Vérifié sur des grilles
      connues (All The Things You Are 36 mesures jouées, Autumn Leaves
      32, Anthropology 32).
- [x] **Les chiffrages** lus en accords de `harmony` : une table des
      62 qualités de l'app, déjà sous la forme que `Normalize` produit
      (un test y veille), plus les qualités tapées à la main, réécrites
      dans l'orthographe de l'app (« m7 », « maj7 », « 7+ ») ou lues
      dans une petite table à part (le diminué 7 ♮14). Sur 60 740
      accords réels, un seul reste illisible : une faute de frappe.
      Table relue et validée : la quinte écrite sauf altération, 2 lu
      comme sus2, h seul comme m7♭5, 11 comme 7sus4 add9, 7susadd3
      comme un 7 avec onzième, 7alt comme l'accord pandiatonique du
      locrien ♭4, 7(♭9, ♭10, ♭5, ♭13).
- [x] **Des cases aux temps** : `Chart.Timeline`, les accords avec leur
      durée, le début de chaque mesure, et la grille lue comme un cycle
      (`Next`). La coda, jouée au dernier chorus seulement, est hors du
      cycle. Un accord commence sur le temps où tombe sa case quand les
      cases se partagent la mesure, arrondi au temps suivant : en 3/4
      sur quatre cases, la case 2 est le temps 3 ; en 5/4, le temps 4
      (Take Five, vérifié à l'écoute). Le temps est la pulsation : une
      noire en 4/4 et 3/4, une noire pointée en 6/8 et 12/8.
- [x] **Les fiches de référence** du livre dans
      `charts/ireal/testdata/fiches` : Tune Up, Black Orpheus, There
      Will Never Be Another You, Tenderly. Les quatre concordent mesure
      par mesure avec les grilles du corpus, à la tétrade près.

Les tests commités n'utilisent que des grilles fabriquées ; les exports
de l'app vont dans `charts/ireal/testdata/local/`, ignoré par git, où un
test les lit s'il y en a.

Reste :

- [ ] les directions que le dépliage ignore encore : « D.C. al 2nd
      ending » et consorts, un D.S. écrit dans une reprise pas encore
      terminée.
- [ ] `irealbook://`, l'ancien schéma non brouillé : refusé tant qu'on
      n'en a pas vu un vrai.
- [ ] l'analyse d'une grille, conçue dans `grilles.md` : étages
      (morceau, plage, région, bloc, accord), analyse de droite à
      gauche, lectures multiples, fiche de sortie. Les fiches du livre
      sont l'oracle : le test compare aujourd'hui les accords, il
      comparera les tonalités, modulations, emprunts et cadences à
      mesure que l'analyse les produira.
- [ ] les seuils de l'analyse, en données : durée d'une plage modale,
      indices d'une modulation.
- [ ] l'ordre de l'analyse, dans `harmony/analysis` comme une couche
      au-dessus de l'identification des accords : la suite d'accords
      (`analysis.Changes` : basses, durées, bouclage) et son pont depuis
      iReal, faits ; les préparations sur toute la suite et les chaînes
      remontées depuis chaque arrivée (`Approaches`, `Chains`), faites,
      et visibles avec `charts/cmd/analyse` ; le passage
      (`PassingChords`), fait ; les blocs et leurs tonalités annoncées
      (`Blocks`), faits ; les relations entre blocs (marches de II-V,
      cycle des quartes) ; degrés et tonalités annoncées en gammes
      précises (`Degrees`, `Block.Announced`), faits, 70 degrés sur 83
      comme le livre ; la tonique pressentie, lue de gauche à droite
      (tonique de fond et tonique locale, le I emprunté), conçue dans
      `grilles.md` ; l'affichage annoté en ASCII dans le terminal,
      enrichi à chaque étape ; les étages hauts dans la mesure où la
      fiche en a besoin. Puis le direct, avec l'attente
      et la surprise (voir `grilles.md`).
- [ ] la grille annotée dans une fenêtre Ebitengine (police de Real
      Book, chiffrages en indices et exposants, réglable), une fois le
      cœur validé.
- [ ] le moteur d'analyse en WASM dans une page web, pour distribuer et
      faire connaître le travail. Le moins prioritaire, à ne pas
      perdre de vue.
- [ ] les voicings sur une grille : la marque « Employée » du dex
      constate une position placée spontanément sur les changes.
- [ ] un parseur de chiffrages général dans `naming`, pour ce qu'on
      tape soi-même, quand un jeu en aura besoin.
- [ ] un format de grille ouvert, mieux conçu que celui d'iReal, une
      fois la structure et le dépliage en place.
- [ ] une grille iReal comme niveau du shoot'em up.

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
