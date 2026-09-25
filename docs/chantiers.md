# Chantiers

Ce qui est ouvert, et ce qui attend une décision.

Ce fichier existe pour qu'une conversation neuve reprenne sans rien
redécouvrir et sans rouvrir un débat déjà tranché. Les raisons des choix
faits sont dans `architecture.md`, le vocabulaire dans `glossaire.md`, la
conception de la collection dans `dex.md`, la première activité dans
`oreille.md`. Ici il n'y a que ce qui reste à faire, et les décisions
qu'il ne faut pas rouvrir.

## Où on en est

`harmony` tient : hauteurs, intervalles, ensembles, gammes, tonalités,
accords, fonctions, systèmes, tétracordes, et les 35 modes. Testé,
benché, vert.

`naming` sépare la langue (français, anglais) et la notation (signes
par défaut, mots en option). Chaque mode a un nom systématique et des
alternatives : le registre parlé français et les alias.

`analysis` identifie les accords de façon déterministe, sans pondération.

`dex` a son corps et sa persistance JSON. Restent `Cooling`,
`Discoverable` et `Components`.

`synth` joue : conversion en fréquence, moteur polyphonique sans
allocation, discipline des buffers audio, et un histogramme des délais
à cases fixes.

`games/keyboard` a deux sources : le clavier MIDI, qui saute les ports
Through quand aucun n'est demandé, et la séquence rejouée, qui joue des
notes datées sur sa propre goroutine avec la gigue d'un doigt.

`games/keys` fait sonner un clavier MIDI, environ 12,5 ms de plancher
mesuré, aucun flam audible.

`games/ear` est le tracer bullet : les sept modes du système naturel à
reconnaître, en français ou en anglais, en signes ou en mots, avec le
dex persisté dans `~/.config/gohar/dex.json`. Il traverse toutes les
couches et il tient au premier playtest.

## Les décisions du dex

**Le dex reçoit des faits, il ne juge pas.** Décider qu'une réponse
est juste, qu'une couleur était ouverte, qu'un trajet harmonique tient,
appartient aux règles du jeu, qui émet ensuite les faits
correspondants. Le dex ne lit aucune cible. `Target` vit dans
`harmony`, avec un `Target.Resolve` qui rend des faits et pas un
verdict : c'est là que le jeu trouve de quoi juger.

**Une erreur compte pour du beurre.** Elle se corrige, elle ne
s'apprend pas : ce qui s'apprend, c'est la correction produite par le
joueur. Aucun fait ne dit qu'une tentative a échoué, et `FactNamed`
n'est envoyé que sur une identification. Dans `ear`, une erreur ne
produit rien : on montre la correction en comparant les deux modes, et
le mode revient plus loin sur une autre tonique.

**`FactSounded` ne touche que `Last`.** `Count` compte des
démonstrations, pas des passages à l'oreille.

**`FactChosen` marque aussi la production.** Qui a choisi a joué :
`Used` et `Produced[tonic]`, sinon un mode toujours improvisé et jamais
exigé n'entre jamais dans la grille des douze toniques.

**La couleur est ouverte tant que la cible ne nomme pas le mode.** Une
tétrade ne donne au mieux qu'une fonction et une indication de couleur
(un m7♭5 implique le mineur). Ce sont les extensions et les degrés
caractéristiques qui disent qu'une couleur a été choisie : sur un slot
`G7`, le lydien dominante n'est attribué que si le ♯11 a sonné. C'est
le jeu qui en juge.

**Ce qui ouvre une entrée.** `FactHeard` (sonner et nommer),
`FactNamed`, `FactProduced`, `FactChosen`. `FactSounded` sur une notion
jamais rencontrée marque `Overheard` et fait apparaître une silhouette,
pas une découverte : c'est le musicien qui joue une couleur parce
qu'elle sonne classe, bien avant d'en connaître le nom. Le jeu peut
alors la nommer, et ce `FactHeard` est la découverte.

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
- [ ] `Notion.Components` : un mode rend-il visible toute sa gamme mère ?
      Laissé vide en attendant. Les silhouettes ne viennent pour
      l'instant que de ce qui a été entendu sans être nommé.
- [ ] `KindSystem` : un système est une entrée à part entière, avec ses
      propres marques. Acté, pas écrit.
- [ ] distinguer le système comme gamme mère des modes et le système
      comme référence de l'harmonie tonale. Le mineur naturel n'est pas
      un citoyen de seconde zone.
- [ ] écrire dans `dex.md` le principe « apprendre au joueur à se passer
      du jeu », et les décisions ci-dessus.

## L'harmonie

- [ ] le pont `analysis` → `FactSounded` : repérer les degrés
      caractéristiques d'un mode sur une fenêtre de jeu. Choix à faire
      sur la taille de la fenêtre et l'ambiguïté entre modes voisins.
      C'est ce qui fera exister les silhouettes et `FactChosen`.
- [ ] passe ligne à ligne sur `Mode.Function`. Le catalogue porte
      aujourd'hui 17 toniques, 9 dominantes, 5 sous-dominantes et 2 sans
      fonction, remplies à partir des tétrades. La dictée d'Arnaud
      donnait tout tonique sauf les deux sans tétrade et le mixolydien
      ♭2 ♭5. Les deux premiers degrés du mineur harmonique et du majeur
      harmonique portent `Tonic | Dominant` à dessein : ce sont aussi
      des avatars de dominante sur pédale de tonique.
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

- [ ] comportement sous charge, Ebitengine et ramasse-miettes. Le
      raisonnement est écrit dans `architecture.md`. L'outil est prêt :
      la touche `H` de `ear` affiche n, p50, p99 et le maximum de
      l'histogramme. Reste à le lire après plusieurs minutes avec le
      clavier vectoriel animé, jamais la moyenne, et à corréler avec
      `GODEBUG=gctrace=1`.
- [ ] calibration chez le joueur. Promise dès le premier jour, et les
      chiffres mesurés ici sont ceux d'une machine, pas une promesse.
- [ ] timbre : quelques partiels, le jour où le frottement des sinus
      purs cessera d'être beau.
- [ ] niveau de sortie bas, conséquence de la marge prise sur le gain.
      Réglage, pas conception.

Deux portes délibérément ouvertes et non planifiées, décrites dans
`architecture.md` : un lecteur `/dev/snd/midiC*D*` en pur Go pour sortir
du C++, et la cible navigateur. Les garder ouvertes ne coûte qu'une
chose, maintenir la règle des deux surfaces : oto n'est importé que par
`synth/device.go`, gomidi que par `games/keyboard/midi.go`.

## Les jeux

La suite de `ear`, dans l'ordre de `oreille.md` :

- [ ] le clavier vectoriel : deux octaves, la note qui sonne animée
      pendant la question, le mode entier surligné seulement à la
      révélation, sinon le clavier donne la réponse.
- [ ] le test sous charge, une fois le clavier animé (voir l'audio).
- [ ] WASM : `midi.go` derrière un build tag, `webmididrv` plus tard,
      le dex dans `localStorage`, et l'écran « cliquer pour commencer »
      qui existe déjà. La police des signes est embarquée, rien à faire
      de ce côté.
- [ ] juger à l'oreille le tempo (350 ms par note) et le registre (la
      gamme entre la3 et sol♯4, la pédale deux octaves dessous).
- [ ] les paliers suivants : distracteurs proches (lydien contre
      ionien), puis les autres systèmes.
- [ ] afficher les alternatives d'un mode (alias, registre parlé) le
      jour où l'activité quitte le système naturel, qui n'en a pas.

Les autres jeux :

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

- [ ] vérifier la section History du README : le lien vers
      `gohar-archive` et ce qu'elle dit de l'ancien module.
- [ ] tags préfixés le jour de la publication : `harmony/v0.1.0`,
      `synth/v0.1.0`.
- [ ] `oreille.md` : le découpage a vieilli, les étapes 1 à 4 sont
      faites et l'étape 4 utilise `text/v2`, pas `DebugPrint`.

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