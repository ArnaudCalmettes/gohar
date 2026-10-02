# Walk with me

Le document de conception du jeu, une puce par décision. Ce qui vient
d'une source le dit ; le reste est à nous. Pas encore de code.

## Le pitch

- Un jeu de piano jazz qui apprend à remplacer le bassiste à la main
  gauche, puis à poser des pêches à la main droite, pour accompagner une
  chanteuse ou un soufflant.
- On joue **à vue sur une grille arbitraire** : la progression porte sur
  des situations harmoniques, jamais sur des morceaux appris par cœur.
- Le canevas, c'est la grille (Chailley, *40 000 ans de musique*,
  p. 306, « improvisateurs sur canevas ») : le joueur reçoit une
  tradition et la modèle à sa guise.

## Les piliers

1. **Jouer à vue.** Une grille n'est maîtrisée que si ses situations le
   sont.
2. **Une seule nouveauté à la fois.** Une difficulté qui monte fait
   redescendre les autres.
3. **On marque, on ne juge pas.** En retard, hors de l'accord, approche
   résolue : le jeu prend acte, il ne note pas le joueur.
4. **On entend ce qu'on joue.** Une vraie basse, un vrai piano, une
   vraie charleston.

## Le joueur

- Un pianiste qui lit les chiffrages et connaît ses accords, mais ne
  sait pas accompagner.
- Il joue sur un clavier MIDI, idéalement 88 touches.

## La boucle de jeu

1. Le jeu choisit une grille et un palier ; le dex dit quelles
   situations travailler.
2. Au premier passage, le jeu fait entendre une vraie ligne de basse.
3. Un décompte, puis la grille défile au tempo, les chiffrages affichés
   comme sur un pupitre, et le joueur prend la main.
4. Chaque note reçoit une marque, et le juice suit : il encourage quand
   ça tourne, il signale discrètement quand ça décroche.
5. En fin de run, ce qui a tenu et ce qui a lâché, par situation.
6. Les marques partent au dex, qui décide de la suite.

## La phase d'entraînement, sans tempo

- La grille attend le joueur : elle avance quand il a joué ce qu'on
  attend, comme le mode « attente » des jeux de piano.
- Les marques de note restent, celles de temps disparaissent.
- Elle dure aussi longtemps qu'il le souhaite, et il en sort à la
  demande : le temps de prendre ses marques, sans pression.

## Le clavier

- Deux zones, partagées au **sol2** par défaut (MIDI 55, G3 en notation
  américaine), le partage réglable.
- La main gauche sonne comme une contrebasse : la walking bass vit pour
  l'essentiel entre le mi grave de la contrebasse (MIDI 28) et le do du
  milieu.
- La main droite sonne comme un piano : les rootless y tiennent leurs
  secondes et leurs tierces serrées. Les seuils exacts des intervalles
  graves restent à sourcer.
- Pas de mains croisées : le partage suffit toujours.

## Deux parcours, puis leur réunion

- **Main gauche, la basse** : fondamentales, two-feel, walk, puis les
  raffinements.
  - Fondamentales : en rondes, en blanches quand la mesure a deux
    accords, dans les temps.
  - Two-feel : blanches systématiques ; la fondamentale sur le temps
    1, la quinte (typiquement) ou une note qui mène à l'accord suivant
    sur le temps 3.
  - Walk : noires ; notes d'approche vers la fondamentale suivante.
  - Raffinements : varier les chemins, grace notes ponctuelles,
    réharmoniser, pédales.
- **Main droite, les pêches** : accords rootless en position serrée,
  posés sur les temps ou en anticipation. Branche le chantier des
  voicings (`voicings.md` : `Position`, tessiture, moindre mouvement).
- **Les deux mains** : un palier à part entière, très dur. Savoir
  marcher à gauche et pêcher à droite ne suffit pas : on revient en
  fondamentales seules, ou en two-feel, pour réunir les deux mains.

## Les axes de progression

- Quatre axes indépendants :
  - **la ligne** (fondamentales, two-feel, walk, raffinements) ;
  - **le tempo** ;
  - **le rythme harmonique** (un accord par mesure, puis deux) : le
    plus local des axes, il varie d'une mesure à l'autre dans une même
    grille ;
  - **les situations** : ce que l'analyse reconnaît (II-V-I majeur,
    II-V mineur, dominante chromatique, cycle de dominantes, cellule
    anatole, région qui s'ouvre, blues), puis les substitutions.
- Plus **les mains** : gauche seule, droite seule, les deux.
- **Une seule nouveauté à la fois.** Les axes avancent ensemble, mais
  un palier n'en fait monter qu'un.
- **Une nouveauté dure fait redescendre les autres.** Aborder les
  substitutions fait ralentir et repasser en fondamentales ou en
  two-feel ; réunir les deux mains aussi.
- À trancher : le palier comme vecteur (ligne, mains, rythme
  harmonique, situations, tempo), et la règle qui dit quels axes
  redescendent quand un autre monte.

## Les marques, temps par temps

- Les temps 1 et 3 sont **forts**, les temps 2 et 4 **faibles**.
- Fondamentales : la fondamentale sur le temps où l'accord arrive.
- Two-feel : sur le temps 1, la fondamentale ; sur le temps 3, une
  note de l'accord, la quinte de préférence, ou une note qui mène à
  l'accord suivant.
- Walk :
  - temps 1 : une note de l'accord, la fondamentale de préférence ;
  - temps 2 et 3 : des notes de l'accord ou de la gamme du moment ;
    le temps 3 reste fort ;
  - temps 4 : une approche de la fondamentale suivante.
- Les approches : chromatique par-dessus ou par-dessous, conjointe dans
  la gamme, par la quinte de la cible (l'approche par la dominante).
  À sourcer chez Siskind.
- **Une note hors de la gamme qui résout est une approche
  chromatique** : bonne sur un temps faible, pas sur un temps fort.
- **Une approche ne se marque qu'après coup** : une note du temps 4
  n'est une approche que si le temps 1 suivant arrive sur sa cible. La
  marque attend d'un temps.
- La gamme du moment vient de l'analyse : le fond, la région, la
  tonicisation (`Sensed`), et la tonalité qu'annonce un bloc.
- Le temps : une tolérance en millisecondes, réglée de façon
  empirique, comme les seuils de l'analyse. Elle suppose la latence
  calibrée (chantier de `oreille.md`).
- On juge ce qui sonne, jamais ce qui s'écrit (`oreille.md`, d'après
  Chailley, p. 162-163).
- Les règles de chaque ligne sont des données : une contrainte de run
  est une règle de plus, pas du code de plus.

## Les contraintes de run

- Des règles en plus, le temps d'un run : « uniquement des approches
  chromatiques », « two-feel sur les A, walk sur le pont », « reste
  sous le do 3 », « jamais la fondamentale sur le temps 1, sauf à
  l'arrivée d'une cadence ».
- Elles forcent à varier les chemins au lieu d'en répéter un.

## Le répertoire

- À côté du jeu à vue, le joueur peut apprendre des standards célèbres
  (dans leurs versions builtin), y revenir et les posséder.
- Ces grilles connues servent aussi de terrain de confort, pour rouler
  sur du terrain sûr entre deux défis.

## La basse de référence

- Au premier passage d'une grille, le jeu joue une vraie ligne de
  basse ; au second, le joueur prend la main.
- Le générateur de lignes et le marqueur partagent les mêmes règles :
  une ligne générée doit recevoir de bonnes marques, et c'est un test
  tout trouvé.

## La jauge de tension

- Le terrain d'essai de la tension et de la surprise, laissées en
  dehors de l'analyse jusqu'ici.
- Le cas d'école : au second passage, le joueur quitte le chiffrage
  quelques mesures pour une descente chromatique. On reconnaît le motif
  au bout de trois ou quatre notes, on l'encourage, et on explose de
  joie s'il le fait atterrir sur la tonique.
- **Reconnaître en cours de route** : des reconnaisseurs en ligne sur la
  ligne de basse (descente ou montée chromatique, marche, pédale, cycle
  de quintes, encadrement d'une cible), qui se déclarent sur le début
  du motif.
- **La tension monte** quand le joueur s'éloigne du chiffrage : notes
  étrangères à l'accord, plus encore sur un temps fort, ligne
  chromatique qui dure.
- **Elle retombe à l'atterrissage**, sur la tonique ou sur la
  fondamentale de l'accord, sur un temps fort. La joie est à la mesure
  de la tension accumulée.
- **Un motif qui s'effondre sans atterrir** reçoit une marque discrète,
  sans casser le groove.
- **L'usure** : on garde la trace des couples (motif, endroit de la
  grille) déjà joués, et la récompense baisse à chaque répétition au
  même endroit. Le même motif ailleurs s'use moins, un motif neuf pas
  du tout.
- À trancher : la forme de la courbe de tension, le poids de chaque
  écart, la vitesse d'usure, et ce qu'on entend par « même endroit »
  (même mesure, ou même situation harmonique).

## Le bonhomme

- Un personnage en pixel art qui marche en rythme : le jeu s'appelle
  *Walk With Me*. Sa démarche dit comment ça tourne, à la place d'une
  jauge abstraite.
- **Son pas** suit le `MetronomeSystem` : un pas par temps, l'animation
  calée sur les temps, jamais sur les images.
- **Sa démarche** suit l'aisance du joueur sur les dernières mesures
  (marques de temps et de note), avec de l'inertie : une fausse note
  isolée ne le fait pas trébucher, une série oui.
- **Ses bulles** suivent la jauge de tension : un motif reconnu, un
  atterrissage.
- **L'usure** vaut pour lui aussi : un motif répété au même endroit le
  fait réagir moins fort (un hochement plutôt qu'un saut), et les
  exclamations tournent pour ne pas se répéter.
- Les états, du plus bas au plus haut :
  1. **il cherche le tempo** : pas hésitant, regard vers le joueur ;
     c'est l'état de la phase sans tempo et du décompte ;
  2. **il marche** : pas régulier, l'état par défaut ;
  3. **il est dedans** : rebond sur les temps, la tête suit ;
  4. **ça tourne depuis un moment** : il claque des doigts sur 2 et 4.
- Les réactions, brèves, par-dessus l'état en cours :
  - **accroché** (une pédale, une descente chromatique reconnue) : il
    tend l'oreille, petite bulle, « Cool ! », « Wooh ! » ;
  - **atterrissage** : il saute, grosse bulle, « Yeah ! »,
    « Monstrueux ! » ;
  - **décroché** : il perd le pas et se rattrape, sans moquerie. On
    marque, on ne juge pas.
- Il vit en périphérie, en bas de l'écran par exemple, et ne masque
  jamais un chiffrage : le joueur lit la grille à vue.
- Il montre l'énergie du moment, jamais un score.
- On le prototype avec des rectangles qui rebondissent, en attendant
  les images.

### Pour le graphiste

- Un cycle de marche par état, quatre images au moins : hésitant,
  régulier, rebond, claquement de doigts sur 2 et 4.
- Des animations brèves : tendre l'oreille, sauter, trébucher puis se
  rattraper.
- Une bulle de BD extensible ; les exclamations sont écrites par le
  moteur dans la police du jeu, pour les traduire et les varier sans
  redessiner.
- À fixer ensemble : la taille du sprite et la palette (noir sur blanc,
  comme le reste, dans l'esprit d'un Real Book ?).

### Les autres marcheurs

- Une foule de personnages secondaires qui peuvent « marcher avec
  lui » : un chat, une jeune fille, un papy, un personnage de cartoon.
  Plus ils sont nombreux, absurdes et variés, mieux c'est.
- **Plus le combo dure, plus il y a de monde** dans la marche, et le
  juice est démultiplié.
- **Ils entrent en musique** : sur le premier temps d'une mesure, mieux
  encore au début d'une section.
- **Quand le joueur décroche, les autres marcheurs aussi**, mais
  progressivement : ils trébuchent, ralentissent, et les derniers
  arrivés partent les premiers. S'il se rattrape dans la mesure, ils
  restent : on ménage le joueur.
- **Chacun a sa démarche**, et elle peut porter une leçon : le papy en
  two-feel, la jeune fille qui marche, le chat en grace notes. Quand le
  joueur change de ligne, celui qui lui ressemble se met en avant.
- **Une collection** : les marcheurs se dévoilent au fil du jeu, et
  certains arrivent avec une situation (le premier backdoor atterri, la
  première pédale tenue).
- **La lisibilité d'abord** : un nombre plafonné à l'écran, en
  périphérie, et jamais devant la grille.
- **Un contrat de sprite commun** : les mêmes cycles et les mêmes
  réactions pour tous, documentés, pour que des contributeurs puissent
  ajouter les leurs.
- C'est un bon cas pour Ark : beaucoup d'entités qui partagent les
  mêmes systèmes (pas, démarche, bulles, entrées et sorties).

## Le carnet

- Tout ce que joue le joueur s'enregistre : notes horodatées, placées
  sur la grille (mesure, temps), avec leurs marques.
- Les passages qui ont fait réagir la jauge de tension sont mis de côté
  : ce sont ses bonnes idées et ses découvertes.
- On peut les réécouter sur la grille et les exporter en MIDI.

## L'affichage

- Ebiten, comme `games/ear`, et Ark si un ECS se justifie.
- Une écriture « marqueur », noire sur fond blanc, comme sur un Real
  Book.
- Les polices de MuseScore, toutes deux sous licence SIL OFL 1.1 :
  **MuseJazz Text** pour les lettres, les chiffres et les qualités,
  **MuseJazz** pour les glyphes musicaux (♮, ♭, ♯, 𝄫, 𝄪, aux points de
  code SMuFL). `naming` écrit les signes Unicode ordinaires ; le rendu
  du jeu les remplace par les glyphes SMuFL.
- MuseJazz Text n'a pas de crénage, son fichier source ayant été perdu
  (forum MuseScore) : les paires qui comptent, aux jonctions des deux
  polices (B♭7, F♮7), se règlent à la main dans le rendu.
- Reste à vérifier où trouver ø et °.
- **Un chiffrage réglable.** Par défaut, la septième majeure s'écrit
  « ♮7 » : bécarre, naturel et majeur vont ensemble, et le joueur s'y
  habitue. C'est un réglage du rendu des chiffrages dans `naming`,
  chantier déjà ouvert.
- Les noms d'accords suivent l'orthographe d'`analyse` (`spell.go`) :
  la lettre du degré, des mouvements lisibles, l'usage pour
  l'affichage.

## Les aides

- La note à jouer affichée, puis seulement indiquée, puis rien.
- Elles s'effacent avec la progression, et se donnent comme des indices.
- À trancher : à quoi ressemble le retour sur une note fausse, sans
  casser le groove.

## Le son

- Au début, un métronome, plus humain qu'un clic : un claquement de
  doigts sur 2 et 4, comme le bonhomme quand ça tourne. Il faut lui
  trouver un échantillon sous licence claire, comme pour les autres
  instruments.
- Plus tard : batterie et piano quand le joueur travaille la main
  gauche, batterie seule (en guise de métronome) quand il joue à deux
  mains. On construit au fil de ce que les autres chantiers débloquent.
- La basse et le piano du joueur passent par le synthé, une zone du
  clavier chacun.
- Le soufflant ou la chanteuse : la mélodie, quand un format de grille
  la portera.
- Ce qu'il manque à `synth` :
  - **des notes programmées dans le futur**, appliquées à l'échantillon
    près par la goroutine audio : aujourd'hui `NoteOn` part à la
    lecture suivante du tampon, et un métronome piloté par un timer
    aurait de la gigue ;
  - **un mélangeur**, puisqu'un `Engine` joue un seul timbre et que
    `synth.Open` ne prend qu'une source ;
  - **les soundfonts**, pour la charleston, la contrebasse et le piano,
    par go-meltysynth dans `synth/soundfont`, déjà au programme de
    `chantiers.md` (« Les jeux ») : licence MIT, rien d'autre que la
    bibliothèque standard, pas d'allocation au rendu, SF2 seulement,
    enveloppé dans la `queue` commune. Reste à trouver, comme pour le
    piano (Salamander, CC-BY, réduit avec Polyphone), une contrebasse et
    une charleston sous licence claire, et un claquement de doigts pour
    le métronome, puis à mesurer sous charge avec
    plusieurs instruments.
- Un `MetronomeSystem` synchronise tout : il tient la carte du temps
  musical (tempo, temps et mesures), que lisent le défilement, le
  marqueur et le son.
  - Il programme le son à l'avance : la boucle d'Ebiten tourne à 60 Hz,
    trop grossière pour déclencher une note à l'heure, et le synthé
    reçoit des événements datés qu'il joue à l'échantillon près.
  - Les notes du joueur, datées par le MIDI et corrigées de la latence
    calibrée, sont placées sur la même carte.

## Les grilles

- Des grilles builtin, et l'import des exports iReal du joueur.
- Les builtin suivent l'idée de Charles Cornell pour un de ses cours :
  une suite d'accords n'est pas protégée, une mélodie l'est. On ne
  reprend que les accords, transposés hors du ton d'origine, sous des
  titres assez évidents pour être reconnus (« Nothing That You Are »,
  « Falling Leaves »). À faire relire avant publication.
- Le lore du jeu se situe dans un univers parallèle au nôtre, où ces
  standards portent ces titres-là.
- Des exercices générés : une cellule, transposée dans les douze tons.

## La progression

- Elle est déléguée au dex : fraîcheur, marques, dévoilement
  (`dex.md`). Le jeu déclare ses situations comme des notions du dex.
- La répétition espacée porte sur les **situations**, pas sur les
  grilles.
- Le corpus fournit les grilles : un joueur faible sur les II-V
  mineurs en reçoit qui en sont pleines.
- Mais le jeu doit rester amusant : des runs sur des difficultés déjà
  surmontées, pour le plaisir de rouler, récompensent le joueur et
  ménagent son estime de lui. Le dosage entre défi et confort fait
  partie du modèle de progression.
- Branche le chantier `Cooling` du dex (les intervalles, jamais
  décidés).

## Les briques

1. **Les temps attendus** : pour chaque temps, l'accord, la gamme du
   moment, temps fort ou faible, l'accord suivant. Pur, testable sur le
   corpus.
2. **Le marqueur** : des notes horodatées et les temps attendus en
   entrée, des marques en sortie. Pur, sans horloge, testé avec
   `keyboard.Sequence`.
3. **Le synthé** : notes programmées, mélangeur, puis soundfonts et
   zones du clavier.
4. **La coquille du jeu** (`games/walk`) : défilement, entrée,
   décompte, affichage.
5. **La progression**, par le dex.

## Le premier jalon jouable

- Les fondamentales, main gauche seule, sur un blues en fa généré, au
  métronome.
- La phase d'entraînement sans tempo.
- Les marques de temps et de fondamentale, sans progression.
- Ensuite, dans l'ordre : la jauge de tension, le carnet, la basse de
  référence.

## Ce que gohar a déjà

- La grille dans le temps : `Changes`, mesures et durées.
- L'analyse : blocs, cadences, cellules, liens, tonique pressentie,
  régions, zones ; de quoi nommer les situations.
- L'orthographe des accords dans `analyse`.
- Le temps réel : le moteur MIDI, `keyboard.Source`,
  `keyboard.Sequence`, le synthé et la mesure de latence.
- Le stockage local de `games/ear`.
- Les degrés et les gammes de `harmony`.

## Les questions ouvertes

- La tolérance de temps et la calibration (`oreille.md`).
- La liste des situations à déclarer dans le dex.
- Le retour sur une note fausse.
- Le lore : le nom de l'univers, le ton, les titres des grilles.

## Les sources

- Jeremy Siskind, *Jazz Piano Fundamentals* : à dépouiller pour la
  basse marchée (notes d'approche, two-feel) et les rootless. À défaut,
  des règles simples et bien connues : la basse reste un instrument
  simple.
- *En Harmonie* et Siron, pour ce qu'ils disent des notes de passage
  et des approches, si ça tient.
- Chailley, pour le canevas.

## Les pièges

- Ne pas confondre jouer à vue et réciter : la tentation d'un mode qui
  rejoue toujours la même grille. Mais il faut aussi laisser le joueur
  prendre ses marques et apprendre des standards célèbres : le tenir
  toujours en déséquilibre rendrait le jeu frustrant.
- Ne pas juger le swing : on marque le temps et la note. Mais on veut
  quand même donner un feel au joueur : des anticipations à la main
  droite, des grace notes. Le feel s'enseigne, s'aide et se fête ; il ne
  se sanctionne pas.
- Ne pas empiler les difficultés : une seule nouveauté par palier.
- Le registre de la main gauche, tenu par le partage du clavier.
- La jauge de tension qui récompense toujours le même tour : l'usure
  est là pour ça.
