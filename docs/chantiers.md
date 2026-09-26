# Chantiers

Ce qui est ouvert, et ce qui attend une décision.

Ce fichier existe pour qu'une conversation neuve reprenne sans rien
redécouvrir et sans rouvrir un débat déjà tranché. Les raisons des choix
faits sont dans `architecture.md`, le vocabulaire dans `glossaire.md`, la
conception de la collection dans `dex.md`, la première activité dans
`oreille.md`. Ici il n'y a que ce qui reste à faire, et les décisions
qu'il ne faut pas rouvrir.

Mettre une doc à jour n'est pas un chantier : ça fait partie de la
tâche qui la fait mentir.

## Où on en est

`harmony` tient : hauteurs, intervalles, ensembles, gammes, tonalités,
accords, fonctions, systèmes, tétracordes, et les 35 modes. Testé,
benché, vert.

`naming` sépare la langue (français, anglais) et la notation (signes
par défaut, mots en option). Chaque mode a un nom systématique et des
alternatives : le registre parlé français et les alias.

`analysis` identifie les accords de façon déterministe, sans pondération.

`dex` a son corps, sa persistance JSON et `Components`. Restent
`Cooling` et `Discoverable`.

`synth` joue : conversion en fréquence, moteur polyphonique sans
allocation, discipline des buffers audio, et un histogramme des délais
à cases fixes.

`games/keyboard` a deux sources : le clavier MIDI, qui saute les ports
Through quand aucun n'est demandé, et la séquence rejouée, qui joue des
notes datées sur sa propre goroutine avec la gigue d'un doigt.

`games/keys` fait sonner un clavier MIDI, environ 12,5 ms de plancher
mesuré, aucun flam audible.

`games/ear` est sorti du tracer bullet : un menu ouvre les quatre
tétracordes et les modes du système naturel (trois choix, ou les sept
à la place de leur degré), en français ou en anglais, en signes ou en
mots, avec le dex persisté dans `~/.config/gohar/dex.json`.

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
- [ ] `KindSystem` : un système est une entrée à part entière, avec ses
      propres marques. Acté, pas écrit.
- [ ] distinguer le système comme gamme mère des modes et le système
      comme référence de l'harmonie tonale. Le mineur naturel n'est pas
      un citoyen de seconde zone.

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
      `Reading`.
- [ ] catalogue de progressions, pour que `ProgressionID` désigne
      quelque chose. Servirait aussi à juger les détours d'une
      réharmonisation (voir les jeux).
- [ ] un registre parlé anglais, s'il en existe un qui mérite d'être
      proposé en alternative. Aujourd'hui seul le français en a un.

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
- [ ] la feuille de route de `oreille.md` : l'abstraction, le menu et
      les tétracordes sont faits ; restent les degrés et leur réponse
      jouée, les réglages et les niveaux paramétrables.
- [ ] WASM : `midi.go` derrière un build tag, `webmididrv` plus tard,
      le dex dans `localStorage`, et l'écran « cliquer pour commencer »
      qui existe déjà. La police des signes est embarquée, rien à faire
      de ce côté.
- [ ] juger à l'oreille le tempo (350 ms par note) et le registre (la
      gamme entre la3 et sol♯4, la pédale deux octaves dessous).
- [ ] les paliers suivants : les autres systèmes. Les distracteurs
      proches sont couverts par le palier des sept modes.
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
      `Mode.Function` devient bloquante.
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

## Pour une conversation neuve

Sur l'harmonie, Arnaud est la source. Il a étudié à l'école de Bernard
Maury et la nomenclature des modes vient de là. Ne pas nommer un mode
par enharmonie, ne pas confondre un intervalle et un accord, dire
« neuvième mineure » et « accord diminué 7 », « dominante chromatique »
plutôt que « substitution tritonique ». Le reste du vocabulaire est dans
`glossaire.md`.

Trois billets publiés sur Zeste de Savoir servent de source à
`analysis` : la représentation des accords, leur reconnaissance, et la
game loop de l'improvisateur. Le troisième explore des possibles, ce ne
sont pas des exigences.

Méthode de travail : Arnaud compile et teste de son côté, les
livraisons se font en archive à extraire à la racine du dépôt, et avant
toute opération touchant beaucoup de fichiers il envoie l'état de son
arbre.