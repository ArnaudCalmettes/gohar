# Chantiers

Ce qui est ouvert, et ce qui attend une décision.

Ce fichier existe pour qu'une conversation neuve reprenne sans rien
redécouvrir et sans rouvrir un débat déjà tranché. Les raisons des choix
faits sont dans `ARCHITECTURE.md`, le vocabulaire dans `GLOSSAIRE.md`, la
conception de la collection dans `DEX.md`. Ici il n'y a que ce qui reste
à faire.

## Où on en est

`harmony` tient : hauteurs, intervalles, ensembles, gammes, tonalités,
accords, fonctions, systèmes, tétracordes, et les 35 modes nommés dans
`naming`. Testé, benché, vert.

`analysis` identifie les accords de façon déterministe, sans pondération.

`dex` a ses types et sa forme persistée. Son corps est bouchonné.

`synth` joue : conversion en fréquence, moteur polyphonique sans
allocation, et la discipline des buffers audio.

`games/keys` fait sonner un clavier MIDI, environ 12,5 ms de plancher
mesuré, aucun flam audible. C'est la validation que le projet est
faisable.

## Le principe : le dex reçoit des faits

Le dex enregistre des faits, il ne juge pas. Décider qu'une réponse est
juste, qu'une couleur était ouverte, qu'un trajet harmonique tient,
appartient aux règles du jeu, qui émet ensuite les faits correspondants.
`Apply` applique des faits sans interpréter la cible.

## Les quatre questions du dex, tranchées

**1. Une réponse fausse ne rafraîchit pas la date.** `Last` ne bouge
pas : le joueur n'a rien démontré, et `FactSounded` couvre déjà la
fraîcheur. Rien n'est jamais retiré.

**2. `FactSounded` ne touche que `Last`.** `Count` sert à la répétition
espacée et compte des démonstrations, pas des passages à l'oreille.

**3. `FactChosen` marque aussi la production.** Qui a choisi a joué :
`Used` et `Produced[tonic]`, sinon un mode toujours improvisé et jamais
exigé n'entre jamais dans la grille des douze toniques.

**4. La couleur est ouverte tant que la cible ne nomme pas le mode.**
Une tétrade ne donne au mieux qu'une fonction et une indication de
couleur (un m7♭5 implique le mineur). Ce sont les extensions et la
présence des degrés caractéristiques qui disent qu'une couleur a été
choisie. Donc sur un slot `G7`, un choix de mode n'est attribué que si
ses degrés caractéristiques ont sonné (le ♯11 pour le lydien dominante),
sinon aucun choix n'a eu lieu. Ce jugement est celui du jeu, pas du dex
(voir le principe) : le jeu émet `FactChosen` ou non.

Conséquence dans le code : `Fact.Target` est supprimé, le dex ne lit
plus aucune cible. `Target` vit déjà dans `harmony`, avec un
`Target.Resolve` qui rend des faits et pas un verdict : c'est là que
le jeu trouve de quoi juger.

**Une erreur compte pour du beurre.** Elle se corrige, elle ne
s'apprend pas : ce qui s'apprend, c'est la correction produite par le
joueur. Aucun fait ne dit qu'une tentative a échoué, donc
`Fact.Correct` disparaît et `FactNamed` n'est envoyé que sur une
identification.

**Ce qui ouvre une entrée.** `FactHeard` (sonner et nommer),
`FactNamed`, `FactProduced`, `FactChosen`. `FactSounded` sur
une notion jamais rencontrée marque `Overheard` et fait apparaître une
silhouette, pas une découverte : c'est le musicien qui joue une couleur
parce qu'elle sonne classe, bien avant d'en connaître le nom. Le jeu
peut alors la nommer, et ce `FactHeard` est la découverte.

## Le dex

- [x] le corps : `New`, `Apply`, `Look`, `Tonics`, `Collection`,
      `Visible`, et `Entry.Overheard` pour les silhouettes. Restent
      `Cooling`, `Discoverable` et `Components`.
- [ ] dans les jeux, mettre en scène la silhouette : « ce pokémon
      s'appelle phrygien bécarre 6, tu veux le capturer ? »
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
- [ ] écrire dans `DEX.md` le principe « apprendre au joueur à se passer
      du jeu ».

## L'harmonie

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

## L'audio

- [ ] comportement sous charge, Ebitengine et ramasse-miettes. Le
      raisonnement est écrit dans `ARCHITECTURE.md`, la mesure reste à
      faire. Lire le pire cas après plusieurs minutes, jamais la
      moyenne, et corréler avec `GODEBUG=gctrace=1`.
- [ ] histogramme des délais dans `keys` plutôt qu'un simple pire cas :
      un pic isolé et un étalement régulier n'ont pas la même cause.
      Utile le jour du test sous charge, pas avant.
- [ ] calibration chez le joueur. Promise dès le premier jour, et les
      chiffres mesurés ici sont ceux d'une machine, pas une promesse.
- [ ] timbre : quelques partiels, le jour où le frottement des sinus
      purs cessera d'être beau.
- [ ] niveau de sortie bas, conséquence de la marge prise sur le gain.
      Réglage, pas conception.

Deux portes délibérément ouvertes et non planifiées, décrites dans
`ARCHITECTURE.md` : un lecteur `/dev/snd/midiC*D*` en pur Go pour sortir
du C++, et la cible navigateur. Les garder ouvertes ne coûte qu'une
chose, maintenir la règle des deux surfaces : oto n'est importé que par
`synth/device.go`, gomidi que par `games/keyboard/midi.go`.

## Les jeux

- [ ] la première activité d'oreille, simple, celle qui ramène au sujet.
      Règles et découpage dans `OREILLE.md`.
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
      un degré de référence.
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

- [ ] un README, avec un mot sur `gohar-archive` et l'ancien module
      mono-dépôt dont les versions restent résolvables par le proxy.
- [ ] tags préfixés le jour de la publication : `harmony/v0.1.0`,
      `synth/v0.1.0`.

## Pour une conversation neuve

Sur l'harmonie, Arnaud est la source. Il a étudié à l'école de Bernard
Maury et la nomenclature des modes vient de là. Ne pas nommer un mode
par enharmonie, ne pas confondre un intervalle et un accord, dire
« neuvième mineure » et « accord diminué 7 », « dominante chromatique »
plutôt que « substitution tritonique ». Le reste du vocabulaire est dans
`GLOSSAIRE.md`.

Trois billets publiés sur Zeste de Savoir servent de source à
`analysis` : la représentation des accords, leur reconnaissance, et la
game loop de l'improvisateur. Le troisième explore des possibles, ce ne
sont pas des exigences.