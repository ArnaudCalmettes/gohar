# Walk with me

Le document de conception du jeu, une puce par décision. Ce qui vient
d'une source le dit, avec sa page ; le reste est à nous. Pas encore de
code.

Les sources : Jeremy Siskind, *Jazz Piano Fundamentals* (Unit 8,
« Playing Bass in Two », p. 115 à 118 ; Book 2, Unit 10, « Drop-Two
Voicings & Bass in Four », p. 204 à 207) ; Jacques Siron, *La partition
intérieure* (§5.6, p. 324 et 325 ; §9.2, p. 690 à 694) ; Jacques
Chailley, *40 000 ans de musique*.

## Le pitch

- Un jeu de piano jazz qui apprend à remplacer le bassiste à la main
  gauche, puis à poser des pêches à la main droite, pour accompagner une
  chanteuse ou un soufflant. C'est ce qu'on attend du pianiste en duo :
  jouer la basse à la main gauche pendant qu'il accompagne de la droite
  (Siskind, p. 115 et 204).
- On joue **à vue sur une grille arbitraire** : la progression porte sur
  des situations harmoniques, jamais sur des morceaux appris par cœur.
- Le canevas, c'est la grille (Chailley, p. 306, « improvisateurs sur
  canevas ») : le joueur reçoit une tradition et la modèle à sa guise.
  Une même trame admet de nombreuses réalisations (Siron, p. 324).

## Les piliers

1. **Jouer à vue.** Une grille n'est maîtrisée que si ses situations le
   sont. Le piège inverse existe aussi : tenir le joueur toujours en
   déséquilibre rendrait le jeu frustrant (voir « Le répertoire »).
2. **Une seule nouveauté à la fois.** Une difficulté qui monte fait
   redescendre les autres.
3. **On marque, on ne juge pas.** En retard, hors de l'accord, approche
   résolue : le jeu prend acte, il ne note pas le joueur. Le swing ne se
   juge pas : le feel s'enseigne, s'aide et se fête.
4. **On entend ce qu'on joue.** Une vraie basse, un vrai piano, de vrais
   claquements de doigts.

## Le joueur

- Un pianiste qui lit les chiffrages et connaît ses accords, mais ne
  sait pas accompagner.
- Il joue sur un clavier MIDI, idéalement 88 touches.

## La boucle de jeu

1. Le jeu choisit une grille et un palier ; le dex dit quelles
   situations travailler.
2. Au premier passage, le jeu fait entendre une vraie ligne de basse
   (voir « La basse de référence »).
3. Un décompte, puis la grille défile au tempo, les chiffrages affichés
   comme sur un pupitre, et le joueur prend la main.
4. Chaque note reçoit une marque, et le juice suit : il encourage quand
   ça tourne, il signale discrètement quand ça décroche.
5. En fin de run, ce qui a tenu et ce qui a lâché, par situation.
6. Les marques partent au dex, qui décide de la suite.

**La phase d'entraînement, sans tempo** : la grille attend le joueur et
avance quand il a joué ce qu'on attend, comme le mode « attente » des
jeux de piano. Les marques de note restent, celles de temps
disparaissent. Elle dure autant qu'il le souhaite.

## Le clavier

- Deux zones, partagées par défaut au **sol2** (MIDI 55, G3 en notation
  américaine), le partage réglable. Pas de mains croisées : le partage
  suffit toujours.
- La main gauche sonne comme une contrebasse, la droite comme un piano.
- **Le registre de la basse.** Celui des quatre cordes de la
  contrebasse, du mi le plus grave du piano (E1) au sol (G2), guide la
  ligne, écrite une octave au-dessus de ce qui sonne (Siskind, p. 116 et
  204). « Prenez les 8vb au sérieux » : beaucoup de pianistes jouent
  leurs basses trop haut. Descendre sous le mi est permis, mais jamais
  sous le do le plus grave (C1) : les trois ou quatre notes les plus
  graves sonnent trop bas (FAQ, p. 117).
- **Le registre de la main droite.** Les voicings A/B à une main gardent
  leur note la plus grave entre C3 et C4 (Siskind, p. 118). À
  réconcilier avec le partage par défaut : un voicing qui commence entre
  C3 et F♯3 tomberait dans la zone de la basse.

## La ligne de basse, d'après les sources

### Progression et ligne

- La basse est la « deuxième mélodie » des musiques harmoniques. La
  **progression de basse** est la succession des notes graves des
  accords, celle que donnent les chiffrages ; la **ligne de basse** est
  sa réalisation par le bassiste, le plus souvent improvisée (Siron,
  p. 324). Le jeu donne la progression, le joueur réalise la ligne.
- **La fondamentale est la note-cible par excellence** : elle ancre le
  rythme harmonique et les changements d'accords. Une ligne simple pose
  les fondamentales sur les points forts, à chaque changement (Siron,
  p. 692 ; Siskind, p. 115).
- **L'octave est libre** : l'harmonie ne donne pas de sens à la hauteur
  d'octave de la basse, une quinte descendante vaut une quarte
  ascendante, et le registre est le choix du bassiste (Siron, p. 325).
- **Legato**, sans silence entre les notes : sur la contrebasse, un
  pizzicato résonne longtemps (Siskind, p. 116, 117 et 204).

### La basse en deux

- Des blanches, sur les temps 1 et 3 : le **two beat**, hérité des
  anciens styles où la basse doublait la main gauche du pianiste stride
  (Siron, p. 690). Les bassistes la jouent aux tempos détendus ou au
  début d'un morceau, avant la walking bass (Siskind, p. 115).
- Les formules (Siskind, p. 115) : un accord par blanche, sa
  fondamentale ; un accord par mesure, la fondamentale au temps 1 et au
  temps 3 la **quinte**, la **tierce**, ou une **voisine chromatique**
  de la fondamentale suivante.
- Les mélanger, une seule rend la ligne prévisible. Alterner montée et
  descente sur le cycle des quintes. Éviter la même note deux fois de
  suite, ou changer d'octave pour masquer la répétition (p. 116).
- Sur un accord tenu plus d'une mesure, on peut quitter la fondamentale
  au premier temps : « fondamentale, seconde, tierce, fondamentale »,
  ou un walk up si les accords suivent le cycle des quintes (p. 117).

### La walking bass

- Une note par temps, colonne vertébrale de la section rythmique, autant
  par l'avance rythmique que par l'assise harmonique ; popularisée par
  Walter Page chez Count Basie (Siron, p. 690).
- Les règles de Siskind (p. 204 à 206), selon le rythme harmonique :
  1. **quatre accords par mesure** : la fondamentale de chacun, ou la
     basse écrite après la barre ;
  2. **deux accords par mesure** : la fondamentale, puis la quinte, la
     tierce, ou la note à un demi-ton de la fondamentale suivante ;
  3. **un accord par mesure, par le cycle des quintes** : le **walk
     up** (un ton puis trois demi-tons montants, G A A♯ B vers C), le
     **walk down** (quatre degrés descendants du mode, D C B A vers G),
     ou **la triade et deux demi-tons descendants** (E♭ G B♭ A vers A♭) ;
  4. **un même accord sur plusieurs mesures** : l'arpège de la triade
     (1 3 5 3), la montée de la gamme vers la quinte du temps 1 suivant,
     ou « trois montées, deux descentes » ;
  5. **un accord par mesure, hors du cycle** : une basse en deux
     doublée (1 1 5 5), une liaison par degrés conjoints et voisines
     chromatiques, ou un saut sur une note de l'accord, puis des degrés
     conjoints, des voisines inférieures et des **enclosures
     chromatiques** vers la cible.
- **Les notes** (Siron, p. 692 et 693) : la fondamentale ; les autres
  notes de la tétrade, qui renversent l'accord quand elles tombent sur
  un temps fort ; les arpèges de la tétrade (1 3 7 5, 1 7 5 3) ; les
  gammes d'accord. Le bassiste est plus « conservateur » que le
  soliste : dorien, mixolydien, ionien sur un II-V-I, là où le soliste
  joue dorien, superlocrien, lydien. La tension entre la ligne et la
  mélodie du soliste est recherchée.
- **Les approches** (Siron, p. 693 et 694), sur le 4e temps, par
  fragments de gamme ou par sauts de quintes descendantes :

  | Mouvement vers la cible | Exemples, vers C |
  |---|---|
  | degré adjacent diatonique | ♭7 → 1 en mineur, 9 → 1 |
  | sensible diatonique ou chromatique | 7Δ → 1, ♭9 → 1 |
  | quinte descendante, quarte ascendante | 5 → 1 |
  | combinaisons | ♭9 7Δ 5 → 1, 7Δ ♭9 → 1 |
  | vers une autre note que la fondamentale | ♭9 de 3 → 3, 5 de 3 → 3 |

  Les chromatismes servent souvent de sensibles supérieures ou
  inférieures. La direction et la force mélodiques comptent (p. 694).
- **Les variations rythmiques** : anticipation, skip, note fantôme,
  triolets. Un équilibre entre la stabilité de la pulsation et les
  relances : sans surprise la ligne devient mécanique, surchargée elle
  perd sa fluidité (Siron, p. 690 et 691).
- **Les plages modales** : la tonique sert de note-cible au début de
  chaque groupe de mesures (Siron, p. 694).
- Siskind conseille de travailler les lignes sur des II-V-I, la main
  droite en voicings anticipés, en descendant par tons pour couvrir six
  tonalités, puis les six autres (p. 206 et 207).

## La main droite, d'après les sources

- Les **voicings A/B à une main**, à jouer à la droite pendant que la
  gauche tient la basse (Siskind, p. 117 et 118) : trois sons, la tierce
  et la septième plus la neuvième ou la quinte.
  - **Type A** : 3, 7, 9 de bas en haut.
  - **Type B** : 7, 3, 5 de bas en haut.
  - Alterner A et B donne la conduite la plus douce sur un II-V-I.
  - Le son le plus grave entre C3 et C4.
- Ensuite, les drop 2 en rythmes anticipés de la Unit 10 (p. 207).
- Le chantier des voicings s'y branche (`voicings.md` : `Position`,
  tessiture, moindre mouvement).

## Les paliers

Ce qui suit est notre lecture des sources.

- **Main gauche** :
  1. les fondamentales, à chaque changement d'accord ;
  2. la basse en deux : la fondamentale, puis la quinte ou la tierce ;
  3. la voisine chromatique au temps 3 ;
  4. la walking bass sur le cycle des quintes : walk up, walk down,
     triade et demi-tons ;
  5. les accords tenus : arpège, montée vers la quinte, trois-deux ;
  6. la liaison libre : conjointe, par sauts, par enclosures ;
  7. les variations rythmiques, les pédales, les réharmonisations.
- **Main droite** : les voicings A/B, posés sur les temps puis en
  anticipation, puis les drop 2.
- **Les deux mains** : un palier à part entière, très dur. Savoir
  marcher à gauche et pêcher à droite ne suffit pas : on revient en
  fondamentales seules, ou en basse en deux, pour réunir les mains.

## Les axes de progression

- Quatre axes indépendants :
  - **la ligne** (les paliers ci-dessus) ;
  - **le tempo** ;
  - **le rythme harmonique** (un accord par mesure, deux, quatre) : le
    plus local des axes, il change d'une mesure à l'autre dans une même
    grille, et c'est lui qui dit quelle règle de Siskind s'applique ;
  - **les situations** : ce que l'analyse reconnaît (II-V-I majeur,
    II-V mineur, dominante chromatique, cycle de dominantes, cellule
    anatole, région qui s'ouvre, blues, accord tenu, plage modale), puis
    les substitutions.
- Plus **les mains** : gauche seule, droite seule, les deux.
- Les axes avancent ensemble, mais un palier n'en fait monter qu'un. Une
  nouveauté dure fait redescendre les autres : aborder les substitutions
  fait ralentir et repasser en fondamentales ou en basse en deux ;
  réunir les deux mains aussi.
- À trancher : le palier comme vecteur (ligne, mains, rythme harmonique,
  situations, tempo), et la règle qui dit quels axes redescendent quand
  un autre monte.

## Ce qu'on attend, temps par temps

- L'attente dit ce qu'on espère sur chaque temps ; la marque prend acte,
  après coup, de ce que le joueur a joué. La **note-cible** garde le
  sens de Siron : la fondamentale sur un point fort.
- Les temps 1 et 3 sont **forts**, les temps 2 et 4 **faibles**.
- **La règle commune** : la fondamentale suivante tombe au temps fort où
  son accord arrive, et les autres temps l'approchent. Siron et Siskind
  concordent.
- Basse en deux : au temps 1, la fondamentale ; au temps 3, une note de
  l'accord ou une voisine de la fondamentale suivante.
- Walking bass :
  - temps 1 : la fondamentale de préférence, ou une autre note de
    l'accord sur un accord tenu ;
  - temps 2 et 3 : des notes de l'accord ou de la gamme du moment ;
  - temps 4 : une approche de la note-cible suivante, nommée d'après
    la table de Siron (degré adjacent, sensible, quinte, combinaison).
- **Une note hors de la gamme qui résout est une approche
  chromatique** : bonne sur un temps faible, une tension sur un temps
  fort.
- **Une approche ne se marque qu'après coup** : une note du temps 4
  n'est une approche que si le temps 1 suivant arrive sur sa cible. La
  marque attend d'un temps.
- **La même note deux fois de suite** reçoit une marque discrète, sauf à
  l'octave (Siskind, p. 116).
- **La classe de hauteur compte, pas l'octave**, dans les bornes du
  registre.
- La gamme du moment vient de l'analyse : le fond, la région, la
  tonicisation (`Sensed`), la tonalité qu'annonce un bloc ; et pour la
  basse, la gamme « conservatrice » de Siron.
- Le temps : une tolérance en millisecondes, réglée de façon empirique.
  Elle suppose la latence calibrée (`oreille.md`).
- **Les marques du palier 1** (`Marker`) : chaque note de la zone basse
  est rapportée au temps le plus proche et reçoit deux marques.
  - Le temps : **à l'heure** (moins de 30 ms d'écart), **en avance** ou
    **en retard** (jusqu'à 100 ms), **entre deux temps** au-delà, et
    alors aucun temps ne la compte. Valeurs de départ, à régler à
    l'oreille.
  - La hauteur, rapportée à l'accord du temps : **fondamentale** (la
    basse écrite d'un accord renversé), **note de l'accord**, **hors de
    l'accord**.
  - Chaque arrivée d'accord, une fois sa fenêtre refermée : **posée**
    si une fondamentale l'a prise, **manquée** sinon.
  - **Doublée** : deux notes sur un même temps. Sur une arrivée, elle
    vaut manquée : le joueur doit faire entendre la note qu'il veut,
    et devant deux notes, le marqueur refuse de deviner laquelle
    (« In the face of ambiguity, refuse the temptation to guess », *The
    Zen of Python*).
  - Entre deux arrivées, répéter la fondamentale ou se taire ne reçoit
    aucune marque : au palier 1, c'est l'exercice.
  - Une autre note de l'accord ne pose pas une arrivée pour un
    débutant, à qui on demande les fondamentales. Plus tard, une règle
    (`Inversions`) l'autorisera, quand le joueur saura quoi faire de
    ce renversement.
- On juge ce qui sonne, jamais ce qui s'écrit (`oreille.md`, d'après
  Chailley, p. 162-163).
- Les règles de chaque ligne sont des données : une contrainte de run
  est une règle de plus, pas du code de plus.

## Les contraintes de run

- Des règles en plus, le temps d'un run : « uniquement des approches
  chromatiques », « basse en deux sur les A, walking bass sur le pont »,
  « reste sous le do 3 », « jamais la fondamentale sur le temps 1, sauf
  à l'arrivée d'une cadence ».
- Elles forcent à varier les chemins au lieu d'en répéter un.

## Les grilles et le répertoire

- Des grilles builtin, l'import des exports iReal du joueur, et des
  exercices générés : une cellule, ou un II-V-I, transposé dans les
  douze tons (comme les exercices de Siskind, p. 206).
- Les builtin suivent l'idée de Charles Cornell pour un de ses cours :
  une suite d'accords n'est pas protégée, une mélodie l'est. On ne
  reprend que les accords, transposés hors du ton d'origine, sous des
  titres assez évidents pour être reconnus (« Nothing That You Are »,
  « Falling Leaves »). À faire relire avant publication. Le lore se
  situe dans un univers parallèle où ces standards portent ces titres.
- À côté du jeu à vue, le joueur peut apprendre des standards célèbres,
  y revenir et les posséder : un terrain de confort pour rouler entre
  deux défis.

## La basse de référence

- Au premier passage d'une grille, le jeu joue une vraie ligne de
  basse ; au second, le joueur prend la main.
- Le générateur applique les règles de Siskind, mesure par mesure,
  selon le rythme harmonique que lit l'analyse ; il mélange les
  formules et évite les répétitions.
- Le générateur et le marqueur partagent les mêmes règles : une ligne
  générée doit recevoir de bonnes marques, et c'est un test tout trouvé.

## La jauge de tension

- Le terrain d'essai de la tension et de la surprise, laissées en dehors
  de l'analyse jusqu'ici.
- Le cas d'école : au second passage, le joueur quitte le chiffrage
  quelques mesures pour une descente chromatique. On reconnaît le motif
  au bout de trois ou quatre notes, on l'encourage, et on explose de
  joie s'il le fait atterrir sur la tonique.
- **Reconnaître en cours de route** : des reconnaisseurs en ligne
  (descente ou montée chromatique, marche, pédale, cycle de quintes,
  walk up, walk down, enclosure), qui se déclarent sur le début du
  motif.
- **La tension monte** quand le joueur s'éloigne du chiffrage : notes
  étrangères à l'accord, plus encore sur un temps fort, ligne
  chromatique qui dure. **Elle retombe à l'atterrissage**, sur la
  tonique ou la fondamentale d'un temps fort ; la joie est à la mesure
  de la tension.
- **Un motif qui s'effondre sans atterrir** reçoit une marque discrète,
  sans casser le groove.
- **L'usure** : la récompense baisse à chaque répétition d'un motif au
  même endroit de la grille ; le même motif ailleurs s'use moins, un
  motif neuf pas du tout. C'est l'équilibre de Siron entre stabilité et
  relances, rendu jouable.
- À trancher : la forme de la courbe, le poids de chaque écart, la
  vitesse d'usure, et ce qu'on entend par « même endroit » (même
  mesure, ou même situation harmonique).

## Le bonhomme

- Un personnage qui marche en rythme : le jeu s'appelle *Walk With Me*.
  Sa démarche dit comment ça tourne, à la place d'une jauge abstraite ;
  il montre l'énergie du moment, jamais un score.
- **Son pas** suit le `Metronome`, un pas par temps, calé sur les
  temps et jamais sur les images. **Sa démarche** suit l'aisance du
  joueur sur les dernières mesures, avec de l'inertie : une fausse note
  isolée ne le fait pas trébucher, une série oui. **Ses bulles** suivent
  la jauge de tension, usure comprise.
- Les états, du plus bas au plus haut :
  1. **il cherche le tempo** : pas hésitant, regard vers le joueur ;
     l'état de la phase sans tempo et du décompte ;
  2. **il marche** : pas régulier, l'état par défaut ;
  3. **il est dedans** : rebond sur les temps, la tête suit ;
  4. **ça tourne depuis un moment** : il claque des doigts sur 2 et 4.
- Les réactions, brèves, par-dessus l'état en cours : **accroché** (une
  pédale, une descente reconnue), il tend l'oreille, petite bulle
  (« Cool ! ») ; **atterrissage**, il saute, grosse bulle (« Yeah ! ») ;
  **décroché**, il perd le pas et se rattrape, sans moquerie.
- Il vit en périphérie, en bas de l'écran par exemple, et ne masque
  jamais un chiffrage.
- **Pour le MVP**, un bonhomme en bâtons dessiné au trait par le moteur,
  avec les mêmes états et animations : il suffit à valider que la
  démarche suit le jeu et que le juice fonctionne, avant de commander le
  moindre dessin.
- **Pour le graphiste**, plus tard : un cycle de marche par état, quatre
  images au moins ; des animations brèves (tendre l'oreille, sauter,
  trébucher et se rattraper) ; une bulle de BD extensible, les
  exclamations écrites par le moteur pour les traduire et les varier. À
  fixer : la taille du sprite et la palette (noir sur blanc, dans
  l'esprit d'un Real Book ?).

### Les autres marcheurs

- Une foule de personnages qui « marchent avec lui » : un chat, une
  jeune fille, un papy, un personnage de cartoon. Plus ils sont
  nombreux, absurdes et variés, mieux c'est.
- **Plus le combo dure, plus il y a de monde**, et le juice est
  démultiplié. Ils entrent sur un premier temps, mieux au début d'une
  section.
- **Quand le joueur décroche, ils décrochent aussi**, progressivement,
  les derniers arrivés partant les premiers. S'il se rattrape dans la
  mesure, ils restent.
- **Chacun a sa démarche**, qui peut porter une leçon : le papy en basse
  en deux, la jeune fille qui marche, le chat en grace notes. Quand le
  joueur change de ligne, celui qui lui ressemble se met en avant.
- **Une collection** : ils se dévoilent au fil du jeu, certains avec une
  situation (le premier backdoor atterri, la première pédale tenue).
- Un nombre plafonné à l'écran, en périphérie. Un contrat de sprite
  commun, documenté, pour que des contributeurs ajoutent les leurs. Un
  bon cas pour Ark : beaucoup d'entités qui partagent les mêmes
  systèmes.

## Le carnet

- Tout ce que joue le joueur s'enregistre : notes horodatées, placées
  sur la grille (mesure, temps), avec leurs marques.
- Les passages qui ont fait réagir la jauge sont mis de côté : ses
  bonnes idées et ses découvertes. On peut les réécouter et les
  exporter en MIDI.

## L'affichage

- Ebiten, comme `games/ear`, et Ark si un ECS se justifie.
- Une écriture « marqueur », noire sur fond blanc, comme sur un Real
  Book.
- Les polices de MuseScore, sous licence SIL OFL 1.1 : **MuseJazz Text**
  pour les lettres, les chiffres et les qualités, **MuseJazz** pour les
  glyphes musicaux (♮, ♭, ♯, 𝄫, 𝄪, aux points de code SMuFL). `naming`
  écrit les signes Unicode ; le rendu les remplace par les glyphes
  SMuFL.
- MuseJazz Text n'a pas de crénage, son fichier source ayant été perdu
  (forum MuseScore) : les paires aux jonctions des deux polices (B♭7,
  F♮7) se règlent à la main. Elle a ♭ ♮ ♯, ø et ° : pour le premier
  jalon, elle suffit seule (`games/walk/fonts`).
- Le jeu choisit son `ChordStyle` : la septième majeure peut s'y écrire
  « ♮7 », puisque la police la met en exposant (voir les décisions du
  nommage dans `chantiers.md`). Les noms d'accords suivent l'orthographe
  d'`analyse` (voir « L'orthographe entendue » dans `grilles.md`).

## Les aides

- La note à jouer affichée, puis seulement indiquée, puis rien. Elles
  s'effacent avec la progression, et se donnent comme des indices.
- À trancher : le retour sur une note fausse, sans casser le groove.

## Le son

- **Le métronome, plus humain qu'un clic** : un claquement de doigts sur
  2 et 4, comme le bonhomme quand ça tourne. Avec la contrebasse qui
  marche, c'est la première chose qu'un joueur entendra : c'est lui qui
  pose l'ambiance. La fonction métrique est partagée par la basse et la
  batterie, un tandem « souvent alchimique » (Siron, p. 694).
- Plus tard : batterie et piano quand le joueur travaille la main
  gauche, batterie seule quand il joue à deux mains.
- La basse et le piano du joueur passent par le synthé, une zone du
  clavier chacun.
- Le soufflant ou la chanteuse : la mélodie, quand un format de grille
  la portera.
- Les échantillons : un piano (Salamander, CC-BY, réduit avec
  Polyphone), une contrebasse, une charleston et un claquement de
  doigts, tous sous licence claire, joués par go-meltysynth (voir « Les
  briques »).
- Un `Metronome` (`games/walk`) tient la carte du temps musical
  (tempo, temps, mesures), que lisent le défilement, le marqueur et le
  son. Il programme le son à l'avance, la boucle d'Ebiten à 60 Hz étant
  trop grossière ; les notes du joueur, datées par le MIDI et corrigées
  de la latence calibrée, se placent sur la même carte.
- Le métronome ne tient que la pulsation, sans style. Les temps forts
  relèvent d'un cadre stylistique : le 1 et le 3 du jazz à quatre temps
  vivent provisoirement dans `Metronome.Strong`.
- **Sous la noire.** Les pêches anticipées sur le « et » du 4 (Siskind,
  p. 207) demandent que la carte descende à la croche : `AtBeats` et
  `Beats` comptent en temps fractionnaires, et `DueEvery` sert les
  croches d'une fenêtre comme `Due` en sert les temps.
- **Le swing**, un décorateur du métronome (`Swing`), réglable par la
  place du « et » dans le temps : 0,5 pour des croches égales, 0,66
  pour le swing ternaire « d'école », 0,75 pour le swing très serré des
  débuts du jazz à La Nouvelle-Orléans. Il déforme le temps à
  l'intérieur de chaque temps, dans les deux sens : de la position vers
  l'instant pour le son, de l'instant vers la position pour le
  marqueur, qui doit reconnaître un « et » joué swing comme un « et ».
  Les temps eux-mêmes ne bougent jamais.
- **L'humanisation** appartient au musicien, pas à la pulsation : un
  décorateur de l'instrument au moment de programmer ses notes. Chaque
  musicien a son placement (devant, sur ou derrière le temps) et sa
  dispersion, et ses vélocités varient. Le marqueur juge toujours
  contre la grille du métronome, jamais contre les notes humanisées de
  l'accompagnement. Le tirage part d'une graine fixe, pour des tests
  reproductibles.
- **Les sons retenus à l'écoute** : une seule soundfont, GeneralUser GS
  (v2.0.3, licence libre y compris dans un logiciel), et un claquement
  de doigts en WAV (newagesoup sur Freesound, CC0), aucun des kits
  essayés n'en ayant.
  - La contrebasse (0:32) : ronde, sans l'attaque percussive de FluidR3
    ni l'inégalité des pizzicati bruts de VSCO-2.
  - Le kit Jazz (128:32) : la charleston au pied (touche 44) et la ride,
    au choix du joueur entre la 1 (touche 51) et la 2 (touche 59), qui
    ne diffèrent que par la hauteur.
  - Plus tard, le joueur pourra charger ses propres soundfonts.
  - Le snap est versionné dans `games/walk/sounds`, avec sa provenance
    (`CREDITS.md`), et embarqué dans l'exécutable. GeneralUser, 32 Mo,
    n'est pas versionnée : `make sounds` la télécharge dans le cache de
    l'utilisateur, d'un commit figé, empreinte vérifiée. Sa licence
    permettrait de l'embarquer ; le jour où le jeu se distribuera en
    binaire, on pourra en extraire les seuls presets utiles.
- **L'accompagnement doit groover**, pas seulement tomber juste : la
  charleston et le snap jouant la même chose sur 2 et 4 sonnent
  redondants. Le motif de ride classique (1, 2 et, 3, 4 et, accents sur
  2 et 4) y pourvoit, sur des croches swing. Les instruments entrent
  l'un après l'autre : la charleston au décompte, la basse à la mesure
  1, le reste ensuite.

## La progression

- Déléguée au dex : fraîcheur, marques, dévoilement (`dex.md`). Le jeu
  déclare ses situations comme des notions du dex, et sans doute ses
  motifs (walk up, walk down, enclosure), reconnus quand le joueur les
  emploie.
- La répétition espacée porte sur les **situations**, pas sur les
  grilles. Le corpus fournit les grilles : un joueur faible sur les II-V
  mineurs en reçoit qui en sont pleines.
- Le jeu doit rester amusant : des runs sur des difficultés déjà
  surmontées, pour le plaisir de rouler. Le dosage entre défi et confort
  fait partie du modèle. Branche le chantier `Cooling` du dex.

## Les briques

1. **Les temps attendus** : pour chaque temps, l'accord, la gamme du
   moment, temps fort ou faible, l'accord suivant, le rythme harmonique
   de la mesure. Pur, testable sur le corpus. Fait pour le premier
   jalon (`Expect`, dans `games/walk`) : l'accord, s'il arrive sur ce
   temps, le suivant, temps fort ou faible, le nombre d'accords de la
   mesure (0 pour un accord tenu depuis une mesure précédente). La
   gamme du moment viendra de l'analyse, la coda plus tard.
2. **Le marqueur** : des notes horodatées et les temps attendus en
   entrée, des marques en sortie. Pur, sans horloge. Fait pour le
   palier 1 (`Marker`, voir « Ce qu'on attend, temps par temps »),
   testé sur des notes datées à la main.
3. **Le synthé** :
   - les **notes programmées** : fait, `ScheduleOn` et `ScheduleOff`
     (voir « Les notes datées » dans `architecture.md`) ;
   - le **mélangeur** : fait, `synth.Mixer` ;
   - les **soundfonts** : fait, `synth.Sampler` (voir « Les
     soundfonts » dans `architecture.md`) ; restent l'audition des
     candidats et une mesure sous charge avec plusieurs instruments ;
   - les **zones du clavier**.
4. **La coquille du jeu** (`games/walk`) : défilement, entrée,
   décompte, affichage, le bonhomme. Première livraison faite : la
   grille en MuseJazz Text, quatre mesures par ligne, la mesure jouée
   grisée et un curseur ; le décompte en grand ; les claquements sur 2
   et 4 ; la main gauche en contrebasse, la droite en piano ; `-demo`
   pour la basse de référence. Restent les marques, la phase sans tempo
   et le bonhomme.
5. **La progression**, par le dex.

## Le premier jalon jouable

- Les fondamentales, main gauche seule, sur un blues jazz lu dans un
  export iReal (le IV en mesure 2, le VI7, un II-V, un turnaround), au
  claquement de doigts sur 2 et 4, la contrebasse et le métronome en
  soundfont.
- La phase d'entraînement sans tempo.
- Les marques de temps et de fondamentale, sans progression.
- Le bonhomme en bâtons.
- Ensuite, dans l'ordre : la basse en deux, la jauge de tension, le
  carnet, la basse de référence.

## Ce que gohar a déjà

- La grille dans le temps : `Changes`, mesures et durées.
- L'analyse : blocs, cadences, cellules, liens (le cycle des quintes),
  tonique pressentie, régions, plages modales, pédales ; de quoi nommer
  les situations et choisir la règle de chaque mesure.
- L'orthographe des accords et les symboles (`ChordStyle`).
- Le temps réel : le moteur MIDI, `keyboard.Source`,
  `keyboard.Sequence`, le synthé et la mesure de latence.
- Le stockage local de `games/ear`.
- Les degrés et les gammes de `harmony`.

## Les questions ouvertes

- Le partage du clavier : sol2 par défaut, alors que les voicings A/B
  descendent jusqu'à C3.
- La tolérance de temps et la calibration (`oreille.md`).
- La liste des situations et des motifs à déclarer dans le dex.
- Le retour sur une note fausse.
- Le lore : le nom de l'univers, le ton, les titres des grilles.
