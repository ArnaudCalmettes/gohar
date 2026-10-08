# Walk with me

Le document de conception du jeu, une puce par décision. Ce qui vient
d'une source le dit, avec sa page ; le reste est à nous. Le code est
dans `games/walk`, et ce qui tourne est résumé dans « Le premier
jalon ».

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
   déséquilibre rendrait le jeu frustrant (voir « Les grilles et le
   répertoire »).
2. **Une seule nouveauté à la fois.** Une difficulté qui monte fait
   redescendre les autres.
3. **On marque, on ne juge pas.** En retard, hors de l'accord, approche
   résolue : le jeu prend acte, il ne note pas le joueur. Le swing ne se
   juge pas : le feel s'enseigne, s'aide et se fête.
4. **On entend ce qu'on joue.** Une vraie basse, un vrai piano, une
   vraie batterie, de vrais claquements de doigts.

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
5. En fin de run, ce qui a marché et ce qui est à consolider, par
   situation (voir « Le bilan »).
6. Les marques partent au dex, qui décide de la suite.

**La phase de déchiffrage, sans tempo** : la grille attend le joueur et
avance quand il a joué ce qu'on attend, comme le mode « attente » des
jeux de piano. Les marques de note restent, celles de temps
disparaissent. Elle dure autant qu'il le souhaite.

## Le clavier

- Deux zones, partagées par défaut au **sol2** (MIDI 55, G3 en notation
  américaine), le partage réglable (`-split`). Pas de mains croisées :
  le partage suffit toujours.
- Pendant une grille, la main gauche sonne comme une contrebasse, la
  droite comme un piano. Ailleurs, dans les menus où l'on improvise sur
  la ligne de basse du `jam` comme dans les leçons, tout le clavier est
  un piano.
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
- **Les marques du palier 1** (`Marker`) : chaque note de la zone basse
  est rapportée au temps le plus proche.
  - Le temps : **sur le temps**, **en avance** ou **en retard**,
    **entre deux temps** au-delà, et alors aucun temps ne la compte.
    Les fenêtres sont des fractions de temps, pour s'élargir quand le
    tempo ralentit (les valeurs sont dans `FirstPalier`).
  - La hauteur, rapportée à l'accord du temps : **fondamentale** (la
    basse écrite d'un accord renversé), **note de l'accord**, **hors de
    l'accord**.
  - Chaque arrivée d'accord : **posée** si une fondamentale l'a prise,
    **manquée** sinon. Sur une mesure où l'accord se prolonge, le temps
    1 attend aussi une note, n'importe laquelle de l'accord : un
    bassiste y joue quand même.
  - **Doublée** : deux notes sur un même temps valent une arrivée
    manquée. Devant deux notes, le marqueur refuse de deviner laquelle
    était voulue (« In the face of ambiguity, refuse the temptation to
    guess », *The Zen of Python*).
  - Une autre note de l'accord ne pose pas une arrivée : au palier 1,
    on demande les fondamentales. Une règle (`Inversions`)
    l'autorisera plus tard.
- On juge ce qui sonne, jamais ce qui s'écrit (`oreille.md`, d'après
  Chailley, p. 162-163).
- Les règles de chaque ligne sont des données : une contrainte de run
  est une règle de plus, pas du code de plus.

## Les contraintes de run

- Des règles en plus, le temps d'un run : « uniquement des approches
  chromatiques », « basse en deux sur les A, walking bass sur le pont »,
  « reste sous le do 3 », « jamais la fondamentale sur le temps 1, sauf
  à l'arrivée d'une cadence », « jamais deux fois la même mesure sur la
  même fondamentale ».
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
- **Les grilles d'aujourd'hui**, en ChordPro (`games/walk/grids`, voir
  `formats.md`), accords seuls, saisis pour gohar sous licence libre :
  - le blues de douze mesures en F, le premier ;
  - *Satin Doll*, en C : des II-V partout, à un ton ou un demi-ton l'un
    de l'autre, puis en F et en G pour le pont. Le terrain idéal pour
    apprendre à varier ses chemins dans plusieurs tonalités ;
  - *Tune Up*, d'après l'analyse d'*En Harmonie* (t. 1, p. 160) : trois
    II-V-I majeurs un ton plus bas l'un que l'autre ;
  - *Autumn Leaves*, en G mineur : les II-V-I du relatif majeur, puis
    ceux du mineur, sur 32 mesures.
  Elles portent pour l'instant leurs vrais titres, dans leurs tons
  d'origine.

## La basse de référence

- Au premier passage d'une grille, le jeu joue une vraie ligne de
  basse ; au second, le joueur prend la main. Le générateur et le
  marqueur partagent leurs règles : une ligne générée doit recevoir de
  bonnes marques, et les tests le vérifient.
- Le générateur (le paquet `games/walk/bass`, joué par `-demo`) suit les cinq règles
  de Siskind (voir « La walking bass »), avec trois écarts :
  - sur un accord tenu, la gamme peut aussi descendre vers la quinte du
    dessous, pour qu'un walk down puisse se poursuivre ;
  - la basse en deux doublée (1 1 5 5) devient 1 1 5 et une approche,
    la seconde quinte sautant une octave ;
  - la mesure jouée la dernière fois sur la même fondamentale n'est pas
    rejouée, si autre chose convient.
- **Les règles de la ligne** : de préférence les quatre cordes,
  jusqu'au do 3 au plus haut (sans cette marge, un fa en haut des
  quatre cordes ne laisse qu'une descente) ; jamais la même note deux
  fois de suite (p. 116) ; pas de saut de plus d'une quinte, l'octave
  seulement sur la fondamentale, jusqu'au do 1 (p. 117). Les deux
  dernières ne viennent pas des sources : la démo sautait de grands
  intervalles à des endroits peu naturels.
- **Une marche se poursuit** deux fois sur trois : après une mesure qui
  avance par degrés, la suivante continue dans ce sens si elle le peut.
- **La démo conclut** : une grille qui boucle sans dire où elle finit
  s'arrêterait sur une dominante. La démo pose alors une dernière note,
  la fondamentale de l'accord où le chorus reboucle.
- **Un modèle écarté** : des notes-cibles sur les temps forts, et des
  approches vers n'importe quelle note sur les temps faibles (Siron,
  p. 693 et 694). Il ouvre des centaines de mesures possibles, et les
  poids qui devaient les départager appelaient chacun une correction de
  plus.
- **La suite** : un catalogue de patterns par situation, deux ou trois
  pour chacune, plus quelques lignes idiomatiques, choisis à partir de
  lignes de bassiste enregistrées avec `-record`. Par exemple, sur le
  blues en fa, un fa qui descend d'une octave en passant par la quinte
  sur la barre de mesure :

      | F7          | F7          | B♭7
      | F E♭ D D♭   | C B♭ A F    | B♭

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
- **Professeur**, dans les leçons du cours des débutants, il change de
  pose : de face, les deux mains dans les poches ; la droite en sort
  pour claquer des doigts à une bonne réponse (`figure/teacher.go`, et
  « Dans le jeu » dans `debutants.md`). Quand l'easter egg l'a poussé à bout, il sort
  de l'écran, de profil, et revient de même, tourné vers la gauche ; les
  fois suivantes, il s'assoit dans l'herbe, puis se met en lotus
  (`gags.go`). Toutes ses poses partagent les proportions du croquis
  (`figure/walker.go`), et ses bulles, dans la partie comme dans les
  leçons, le même dessin (`figure/bubble.go`), la queue visant sa tête
  où qu'elle soit.
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
- Il vit en périphérie, à gauche de la grille, et ne masque jamais un
  chiffrage.
- **Pour l'instant, un bonhomme en bâtons** dessiné par le moteur
  (`figure/walker.go`), qui suffit à valider que la démarche suit le
  jeu avant de commander le moindre dessin :
  - l'aisance est une moyenne glissante des arrivées, à peu près les
    huit dernières ; il saute à chaque arrivée posée, et ne trébuche
    sur une manquée que si l'aisance est déjà retombée ;
  - deux images, comme un cycle de sprite, jamais d'entre-deux : une
    animation fluide colle mal au pas. Sur les temps forts, les jambes
    écartées ; sur les temps faibles, les jambes qui se croisent ;
  - le claquement de doigts à la manière jazz, du coude ;
  - les seuils sont des valeurs de départ, à régler en jouant.
- **Il parle**, avec les mots du jeu (des marques, pas un jugement),
  au-dessus de sa tête, écrit à la main de la grille et tenu une
  seconde et demie (`figure/coach.go`). Pas de bulle fermée : un trait
  sous la phrase et un trait qui descend vers lui, à la manière d'une
  BD dessinée en bâtons comme lui :
  - **rien quand le joueur joue dans les temps** : le claquement de
    doigts le dit déjà, c'est le juice ;
  - « Tu presses » ou « Détends-toi » quand les cinq dernières notes
    tombent en moyenne 30 ms ou plus en avance ; « Tu traînes » ou « Ça
    traîne » quand elles tombent en retard. C'est une dérive qu'il
    relève, pas une faute : un bon joueur presse sans sortir de la
    fenêtre « sur le temps », et une première version, qui comptait les
    notes hors de cette fenêtre, ne parlait presque jamais une fois la
    latence calibrée ;
  - un mot quand ça se met à tourner, au moment où il commence à
    claquer des doigts, puis toutes les huit arrivées posées d'affilée
    tant que ça tourne, à peu près une fois par chorus de blues ; et
    quand le joueur se rattrape, cinq notes revenues à 15 ms ou moins
    du temps en moyenne après une remarque : « Cool ! », « Yeah ! », « Continue
    comme ça ! », « Ça swingue ! », « Groovy ! », « Super ! », « Ça
    joue ! » ; en anglais « Keep it up! », « Swingin'! », « I dig
    that! »… ;
  - jamais deux fois le même mot de suite, et huit temps de silence
    après chaque bulle, pour qu'il ne « radote » pas ; la remarque
    suivante pèse des notes neuves ;
  - seulement quand le joueur tient la basse : en démo, il se tait ;
  - tout cela suppose la latence calibrée : sinon, elle passerait pour
    de la précipitation ;
  - les seuils (cinq notes, 30 ms, huit arrivées, huit temps) sont des
    valeurs de départ, à régler en jouant.
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

## Le bilan

À la fin d'un run joué jusqu'au bout, et seulement alors, le jeu montre
ce qui a marché et ce qui est à consolider. Un run arrêté avec Espace
n'a pas de bilan, une démo non plus : le bonhomme n'a rien à se dire.

- **Le principe** est celui du marqueur et du dex : le bilan prend acte
  et ne juge pas. Il dit les comptes tels qu'ils sont, « 22 sur 24 »,
  et ne relance pas le joueur.
- **Les situations** sont, pour l'instant, la place du changement dans
  la grille :
  - un accord qui arrive en début de mesure ;
  - un accord qui arrive au milieu de la mesure, quand deux accords se
    la partagent : c'est la plus difficile, le moment où la main doit
    changer de fondamentale à mi-chemin ;
  - le temps 1 d'un accord qui dure.
- **La carte de la grille** montre la grille entière, ses rangées
  resserrées si elle est longue, chaque mesure teintée en vert si elle
  a marché, en orange si elle est à consolider, sur tous les chorus
  ensemble. Un musicien y voit d'un coup d'œil où ça a coincé : les
  deux mesures du turnaround, le II-V qui module.
- **Le seuil** : une situation ou une mesure a marché quand trois temps
  attendus sur quatre au moins ont été posés. C'est une valeur de
  départ, à régler en jouant.
- **Le temps**, en un mot plutôt qu'en millisecondes : « Bien en
  rythme », « Tendance à presser », « Tendance à traîner » ou « À
  stabiliser ». Le jeu mesure l'écart moyen des notes au temps et leur
  régularité autour de cet écart (l'écart type), la latence calibrée
  déduite. Au-delà de 20 ms d'écart moyen, le joueur presse ou traîne ;
  en deçà, au-delà de 30 ms de dispersion, le placement est à
  stabiliser. Des valeurs de départ, comme les seuils du bonhomme. Les
  millisecondes restent dans l'enregistrement. Pas de « tu » dans ce
  verdict : « tu presses » sonne comme un reproche, « tendance à
  presser » comme un constat.
- **Le conseil**, une seule priorité, comme un professeur la donnerait
  en fin de cours :
  - quand des changements lâchent et que le temps est bon : travailler
    la formule qui a le plus lâché, à ce tempo, « Conseil : travaille
    les anatoles à ce tempo ». Les formules sont celles que l'analyse
    trouve dans la grille : les anatoles, les III-VI-II-V-I, les
    cadences éoliennes (`analysis.Cells`), puis les II-V-I
    (`analysis.Blocks`), et les II-V dont le V ne résout pas, comme le
    premier Dm7 G7 de *Satin Doll* ; un temps manqué compte pour la formule dont
    son accord fait partie. Hors de toute formule, le conseil nomme la
    situation : les changements en milieu de mesure, par exemple ;
  - quand les changements et le temps lâchent ensemble : ralentir le
    tempo ;
  - quand seul le temps lâche : rester à ce tempo, le temps de le
    stabiliser ;
  - quand tout a marché : monter le tempo, de 10 sous 140 à la noire,
    de 5 au-delà.
- **Plus tard**, le bilan suivra le cursus : ce qu'il compte découlera
  de ce que mesure chaque palier.

N'importe quelle touche, du clavier d'ordinateur ou du piano, ramène à
la partie.

## L'affichage

- Ebitengine, comme `games/ear`, et Ark si un ECS se justifie.
- Une écriture « marqueur », noire sur fond blanc, comme sur un Real
  Book : les accords en **MuseJazz Text**, la police de MuseScore
  (licence SIL OFL 1.1), qui a ♭ ♮ ♯, ø et °. Elle n'a pas de crénage,
  son fichier source ayant été perdu : les paires délicates se règlent
  à la main.
- Le jeu choisit son `ChordStyle`. Les noms d'accords doivent suivre
  l'orthographe de l'analyse (voir « L'orthographe entendue » dans
  `grilles.md`), zone par zone. **Pour l'instant, ils gardent celle de
  la grille** : une seule tonalité ne suffit pas à épeler une grille qui
  module, et le C du Cmaj7 de *Tune Up*, entendu en D, s'écrirait B♯.
  La réécriture par zones existe dans `charts/cmd/analyse` ; il reste à
  la passer dans la bibliothèque (voir `chantiers.md`).
- **La mise en page** : la grille sur les deux tiers droits de l'écran,
  quatre mesures par ligne ; le bonhomme dans le tiers gauche ; le
  clavier en dessous, sur toute la largeur.
- **Trois rangées à la fois.** Une grille plus longue tourne ses pages
  d'une rangée à la fois, comme on lit un Real Book : la rangée jouée en
  deuxième position, celle d'avant au-dessus, la suivante en dessous,
  pour lire en avance. La page glisse d'une rangée à l'autre, en un
  tiers de seconde qui ralentit à l'arrivée, plutôt que de sauter :
  l'œil suit la grille au lieu de la chercher. Une rangée à moitié
  sortie est coupée net au bord de la page. Au repos, la molette ou le
  doigt glissé sur la grille la font défiler, pour la lire en entier
  avant de la jouer ; un chevron au-dessus ou en dessous dit qu'il
  reste des rangées cachées de ce côté. Une mesure à deux accords les écrit plus petits,
  et alors toute la grille avec eux : une seule taille d'accord par
  grille, pour que la page se lise d'un œil égal.
- **Les marques de sections.** Chaque section porte sa lettre, encadrée
  à gauche de sa première mesure, comme dans un Real Book : celle que
  la grille écrit, son label, rappels compris (A1, A2, B, A1 pour *Satin
  Doll*). Le jeu ne la déduit pas des accords : une grille est un
  niveau de jeu, qu'un joueur pourra un jour écrire lui-même, et le jeu
  la montre telle qu'elle est écrite. La forme trouvée par l'analyse
  (`analysis.Sections`) reste l'affaire des outils d'analyse. Une
  grille sans label, comme le blues, n'en montre aucune.
- **Les exposants.** La fondamentale et le type de tétrade restent sur
  la ligne (D7, Am7, Cmaj7) ; la quinte altérée et les extensions
  montent en exposant : Am7<sup>♭5</sup>, D7<sup>(♭13)</sup>. Restent
  encore sur la ligne, à ajuster quand le cas se présentera : le
  mineur-majeur avec ses extensions, Cm(maj7,9), l'alt du 7alt, et les
  extensions écrites à la place de la septième, C9, C13.
- **Le tempo est celui de la grille** : le `{tempo}` que le fichier
  écrit, sinon celui de son style (`{meta: style Medium Swing}`), que
  le jeu traduit par une table à lui (Ballad 70, Medium Swing 120,
  Medium Up Swing 160, Up Tempo Swing 200), sinon 120. Le joueur s'en
  écarte pour s'entraîner, de 5 en 5 entre 60 et 240, et l'écran
  rappelle celui de la grille à côté (« 100 BPM (grille : 120) ») ;
  l'écart vaut jusqu'au changement de grille, sans être gardé d'une
  séance à l'autre. Le drapeau `-bpm` impose un tempo à toutes les
  grilles. La progression du joueur en tempo (« Tune Up à 135, monte à
  140 ») viendra plus tard.
- **Les boutons**, pour une souris ou un doigt : la partie se joue
  aussi sur un écran tactile, si peu de place qu'il y ait. Trois
  boutons, que l'écran montrait déjà, et la grille qui défile :
  - le chevron à gauche du titre ramène au titre (Échap) ;
  - le bloc du titre et du tempo, souligné au repos, ouvre les
    réglages de la grille (Tab) ;
  - le cadre du mode, en haut à droite, lance et arrête (Espace).

  Dans un palier, seul le retour reste. Le bilan se ferme d'une touche
  ou d'un toucher n'importe où.
- **Les réglages de la grille**, au repos : la grille, son tempo, le
  mode, a tempo ou déchiffrage, et la démo. ↑ et ↓ choisissent, ← et →
  changent, comme un toucher à gauche ou à droite du milieu d'une
  valeur ; la grille choisie vaut pour la séance. T et D restent des
  raccourcis dans la partie.
- **La démo**, dans les réglages ou d'un D : elle s'allume ou s'éteint,
  avec le tempo seulement. Face à une grille dont il ne voit pas comment elle
  doit sonner, le joueur met la démo à 160 pour en saisir la logique et
  trouver des chemins, puis la travaille lui-même. Sans clavier MIDI,
  elle est allumée d'office : le bonhomme joue, le joueur écoute. Le
  nombre de chorus viendra au même endroit.
- **La ligne d'état dit les touches telles qu'elles sont** : « Espace :
  jouer » à l'arrêt, « Espace : arrêter » pendant un run ; T et D
  nomment ce vers quoi ils basculent (« T : déchiffrage », « D : démo
  ON ») ; « Tab : réglages » ouvre ceux de la grille.

## L'écran titre

Une ouverture jouée par le moteur du jeu, sans vidéo :

1. le décompte en claquements de doigts sur 2 et 4, en gros plan sur le
   seul bout du bras du bonhomme, la main qui claque ;
2. pendant le décompte, la caméra recule ; elle a fini de reculer au
   moment où la basse et la ride entrent : le bonhomme est à sa place,
   à gauche, et marche en claquant des doigts ;
3. le titre et le menu s'affichent, sur le même blues de douze mesures
   improvisé par la basse de référence.

C'est purement cosmétique, et cela demande des transitions entre
scènes : le premier écran titre se contente de la position finale.

- **Ce qu'on voit** : le bonhomme grisé à gauche, à la place qu'il
  occupe pendant la partie, qui marche en claquant des doigts ; le
  titre et le menu au milieu, « Apprendre », « Jouer une grille »,
  « Options » et « Quitter », choisis avec les flèches et validés par
  Entrée, ou d'un clic, d'un doigt sur un écran tactile : les menus se
  parcourent ainsi sur un téléphone. « Apprendre » ouvre le cours des grands débutants
  (`debutants.md`), en premier parce qu'un nouveau venu lit le menu de
  haut en bas.
- **Ce qu'on entend** : à 160 à la noire, plus enlevé que la partie
  (120 par défaut), car le titre invite, il n'enseigne pas.
  - D'abord deux mesures de décompte, les seuls claquements du
    bonhomme sur 2 et 4, pendant qu'il marche déjà : un avant-goût de
    l'ouverture à venir, sans le gros plan.
  - Puis la basse de référence, la ride et le charley sur 2 et 4,
    comme pendant la partie, avec les claquements.
  - Le clavier sonne déjà, pour jouer par-dessus.
- **Un chorus après l'autre** : la ligne est tirée à neuf à chaque
  tour de grille. Elle ne se raccorde pas encore d'un chorus au
  suivant : la dernière note du chorus ne prépare pas la première du
  suivant.
- **D'une scène à l'autre** : la musique du titre, le `jam`, vit hors
  des scènes, dans le contexte partagé (voir « L'orchestre »). Elle passe sans coupure du titre
  aux options et à la calibration, et retour, la basse sur le même
  blues et le même temps ; seule l'orchestration change d'un écran à
  l'autre :
  - au titre, la ride, le charley et les claquements ;
  - dans les options et dans le cours, la basse et le charley sur 2 et
    4 seulement, le bonhomme qui marche dans le rythme sans claquer des
    doigts ;
  - sous la mesure battue de la calibration, la basse seule, plus douce.

  Jouer une grille l'arrête : la partie lance son propre décompte, à
  son tempo.
  Échap ramène de la partie au titre, où le blues reprend avec son
  décompte de claquements, et du titre quitte le jeu.
- **Les options** : la langue, qui change tout de suite, et la
  calibration (voir `architecture.md`), qui n'est jamais imposée : le
  jeu se joue sans. Les flèches gauche et droite changent une valeur ;
  un clic à gauche du milieu de l'option la baisse, à droite il la
  monte. La langue est gardée d'une séance à l'autre (`walk.json`,
  parmi les réglages communs) ; elle sert de valeur par défaut au
  drapeau `-lang`, de sorte qu'un drapeau passé au programme l'emporte
  encore. Le tempo n'y est plus : il est celui de chaque grille.
- **On passe toujours par le titre**, même avec des drapeaux sur la
  ligne de commande. Quand les options auront repris tous les
  drapeaux, en passer au programme voudra dire qu'on le teste, et un
  drapeau permettra d'aller droit à la scène qui intéresse.

## Les aides

- La note à jouer affichée, puis seulement indiquée, puis rien. Elles
  s'effacent avec la progression, et se donnent comme des indices.

## Le son

- **Le métronome, plus humain qu'un clic** : la batterie. Avec la
  contrebasse qui marche, c'est elle qui pose l'ambiance ; la fonction
  métrique est partagée par la basse et la batterie, un tandem
  « souvent alchimique » (Siron, p. 694).
- **L'accompagnement doit groover**, et les instruments entrent l'un
  après l'autre :
  - au décompte, la charleston au pied, deux mesures comptées « 1, 3,
    1, 2, 3, 4 » : le temps d'amener les mains de la barre d'espace au
    clavier ;
  - ensuite, le motif de ride classique (chaque temps et le « et » de 2
    et 4, en croches swing, accentué sur 2 et 4), la charleston sur 2
    et 4, et la contrebasse de la démo ;
  - **le juice s'entend** : le claquement de doigts ne vient sur 2 et 4
    que quand le bonhomme claque des doigts. Le joueur l'entend arriver
    avant même de regarder l'écran, et l'entend partir quand il
    décroche.
- La basse et le piano du joueur passent par le synthé, une zone du
  clavier chacun.
- Plus tard : batterie et piano quand le joueur travaille la main
  gauche, batterie seule à deux mains ; le soufflant ou la chanteuse
  pour la mélodie, quand un format de grille la portera.
- **Le temps musical** est une carte (`games/tempo`), que lisent le
  défilement, le marqueur et le son, programmé à l'avance (voir « Les
  notes datées » dans `architecture.md`). Elle ne tient que la
  pulsation : les temps forts relèvent d'un style, et vivent
  provisoirement dans `Metronome.Strong`. Elle descend sous la noire,
  pour les pêches anticipées sur le « et » du 4 (Siskind, p. 207).
- **Le swing** déforme le temps à l'intérieur de chaque temps, réglable
  par la place du « et » : 0,5 pour des croches égales, 0,66 pour le
  swing ternaire « d'école », 0,75 pour le swing serré des débuts du
  jazz à La Nouvelle-Orléans. Dans les deux sens : le son place un
  « et » swing, le marqueur reconnaît un « et » joué swing comme un
  « et ». Les temps eux-mêmes ne bougent jamais.
- **L'humanisation**, à venir, appartient au musicien, pas à la
  pulsation : chacun son placement (devant, sur ou derrière le temps),
  sa dispersion, ses vélocités. Le marqueur juge toujours contre la
  pulsation, jamais contre l'accompagnement humanisé.
- **Les sons retenus à l'écoute** : une seule soundfont, GeneralUser GS
  (v2.0.3, licence libre y compris dans un logiciel), pour la
  contrebasse (0:32, ronde, sans l'attaque percussive de FluidR3) et le
  kit Jazz (128:32 : la charleston au pied et la ride 1) ; et un
  claquement de doigts en WAV (newagesoup sur Freesound, CC0), aucun
  kit n'en ayant. GeneralUser, 32 Mo, n'est pas versionné : `make
  sounds` la télécharge dans le cache de l'utilisateur, d'un commit
  figé, empreinte vérifiée.
- **La soundfont réduite**, embarquée dans le jeu avec le snap
  (`games/walk/band/sounds`) : `make slim` tire de GeneralUser `walk.sf2`,
  2,5 Mo environ au lieu de 32, réduite à ce que le jeu joue : la
  contrebasse (0:32), le piano (0:0) et huit touches du kit Jazz
  (128:32 : la charleston au pied, la ride 1, les deux wood blocks de
  la calibration ; la cloche, casserole des fausses notes des leçons ;
  la grosse caisse, la caisse claire et la charleston fermée du
  chapitre 2 du cours). C'est `synth/sf2` qui découpe : il garde les presets
  demandés, leurs instruments et leurs échantillons, et pour un kit les
  seules zones des touches frappées. Un test vérifie, sur GeneralUser
  quand il est là, que chaque note gardée sonne comme avant, échantillon
  pour échantillon. Sa licence permet de la modifier et de l'utiliser
  dans un logiciel ; elle prévient que l'origine de certains
  échantillons est incertaine, ce qui ne concerne qu'un produit
  commercial (voir `games/walk/band/sounds/CREDITS.md`). Le jeu joue la soundfont
  embarquée ; `-sf2` en essaie une autre, GeneralUser complète par
  exemple, pour choisir un autre son. Un son de plus dans le jeu, la
  casserole par exemple, demande de refaire `make slim`.

## L'orchestre

Derrière le joueur joue un petit orchestre : une contrebasse, une
batterie réduite à la ride et à la charleston, et les claquements de
doigts du bonhomme. Il accompagne la partie, mais aussi les menus, où
il ne s'arrête jamais d'un écran à l'autre. Le faire jouer juste, à la
milliseconde, sur une boucle de jeu qui ne tourne que soixante fois par
seconde, demande quatre étages, chacun avec un seul métier.

### Le synthé : jouer à l'échantillon près

Tout en bas, les instruments du paquet `synth` reçoivent deux sortes
d'ordres :

- les ordres **immédiats** (`NoteOn`, `NoteOff`), appliqués au début du
  prochain tampon audio : c'est le chemin des touches du joueur, qui
  doivent sonner sans attendre ;
- les ordres **datés** (`ScheduleOn`, `ScheduleOff`), que la goroutine
  audio applique sur l'échantillon qui correspond à leur date (voir
  « Les notes datées » dans `architecture.md`).

Les ordres datés existent parce que la boucle d'Ebiten tourne toutes
les 16,7 ms : une ride envoyée depuis elle tomberait n'importe où dans
cet intervalle, et cette gigue s'entend. On programme donc chaque son
un peu à l'avance, avec sa date exacte.

### Les parties : qui joue quoi sur un temps

Un musicien de l'orchestre, la ride par exemple, sait ce qu'il joue sur
chaque temps sans regarder l'horloge : « ding, ding-da, ding,
ding-da », appuyé sur 2 et 4. Le jeu le modélise ainsi (`parts.go`) :

- une **partie** est une fonction pure : on lui dit où l'on est dans la
  mesure, quelle note la ligne de basse a sur ce temps, si le bonhomme
  claque des doigts ; elle répond par des **coups** : un son, une
  vélocité, et une place dans le temps, 0 sur le temps, deux tiers pour
  le « et » swingué ;
- un **arrangement** est une liste de parties qui jouent ensemble.

Les parties d'aujourd'hui :

| Partie | Ce qu'elle joue |
|---|---|
| le décompte à la charleston | 1 et 3 dans la première mesure, chaque temps dans la seconde : « 1, 3, 1, 2, 3, 4 » |
| le décompte en claquements | les claquements seuls, sur 2 et 4 |
| la ride | chaque temps, plus fort sur 2 et 4, et le « et » swingué après 2 et 4 |
| la charleston | fermée au pied sur 2 et 4 |
| les claquements | sur 2 et 4, quand le bonhomme claque des doigts |
| la basse | la note de la ligne, quand l'orchestre en a une |
| la basse douce | la même, plus discrète |
| la ride de fin | un coup, sur la dernière note de la démo |

Et les arrangements qu'elles composent :

| Arrangement | Parties | Où |
|---|---|---|
| complet | ride, charleston, claquements, basse | la partie, l'écran titre |
| léger | charleston, basse | les options |
| basse seule | basse douce | sous la mesure battue de la calibration |
| décompte | décompte à la charleston | avant la partie |
| ouverture | décompte en claquements | avant la musique du titre |
| fin | ride de fin, basse | la dernière tonique de la démo |

Parce qu'elles sont pures, les parties se testent comme un musicien
les décrirait : « sur le temps 2, la charleston, la ride appuyée, son
"et" swingué, le claquement et la basse ». Et un nouvel arrangement, une
ambiance de plus pour un nouvel écran, tient en une ligne.

### L'orchestre : jouer les coups

Le `band` (`sound.go`) possède les instruments : la contrebasse et le
piano de la soundfont, le kit Jazz, le claquement enregistré, ou les
voix 8 bits sans soundfont. Il ne décide de rien : `band.play` reçoit
les coups d'un temps et les traduit en ordres datés.

- **La basse est legato**, comme le demande Siskind : chaque note tient
  jusqu'à la suivante, relâchée à la date même où l'autre attaque. Il
  n'y a qu'une voix de basse pour tout le jeu, d'où la note tenue que
  l'orchestre retient, et qu'il relâche quand la musique s'arrête.
- **Une cymbale sonne jusqu'à ce qu'on la refrappe** : chaque coup
  relâche d'abord le précédent, pour que la ride ne s'empile pas sur
  elle-même.
- **Les touches du joueur** passent à côté : elles sonnent tout de suite,
  la contrebasse sous le split pendant une grille, le piano au-dessus
  et partout ailleurs, depuis la goroutine MIDI, sans attendre la boucle
  de jeu.

### Le chef d'orchestre : quand jouer

Reste à savoir quand appeler les parties. Un chef tient un métronome
(un instant de départ et la durée d'un temps) et le numéro du prochain
temps à programmer. À chaque image, il demande quels temps tombent dans
les 100 prochaines millisecondes, et programme chacun une seule fois,
même si une image arrive en retard. Les temps négatifs sont le
décompte.

Il y a aujourd'hui deux chefs :

- **la partie**, pendant un run : l'arrangement complet, précédé du
  décompte, et ce qui est propre au jeu : les notes de la démo
  transmises au marqueur à l'instant où elles sonnent, la fin au bout
  des chorus, la dernière tonique ;
- **le `jam`** (`music.go`), la musique des menus.

### Le `jam` : une musique qui ne s'arrête pas

Le `jam` est un chef d'orchestre avec sa basse : un métronome à 160, et
un bassiste qui tire une ligne neuve à chaque tour de grille, avec le
générateur de la basse de référence.

**Il vit hors des scènes**, dans le contexte partagé du jeu. Son
métronome est un instant absolu, qui ne dépend de personne : peu
importe donc quelle scène le fait jouer. À chaque image, la scène du
dessus l'appelle avec son arrangement : complet au titre, léger dans
les options, basse seule sous la calibration. Quand on passe d'un écran
à l'autre, la nouvelle scène reprend au temps exact où l'ancienne
s'était arrêtée ; et comme les 100 ms suivantes étaient déjà
programmées, rien ne manque à l'oreille. La basse ne s'interrompt
jamais ; seule l'orchestration change, sur le temps.

- **La calibration** ne connaît ni le `jam` ni l'orchestre : le jeu lui
  prête le métronome du `jam` et une fonction qui le fait jouer. Sa
  mesure battue tombe ainsi sur les temps de la musique, et part sur le
  premier temps de la mesure suivante.
- **Jouer une grille l'arrête** : la partie lance son propre décompte, à son
  tempo. En revenant au titre, un `jam` neuf repart avec son ouverture
  en claquements. Une leçon l'arrête aussi ; en revenant à la liste des
  leçons, où le bonhomme marche sans claquer des doigts, le `jam` neuf
  repart sans ouverture, l'orchestre tout de suite.

### Ce qui reste à faire

- **Les deux chefs se ressemblent**, et la boucle qui programme les
  temps est écrite deux fois. Les fusionner maintenant, ce serait
  deviner leur interface commune ; le besoin qui la dessinera viendra
  avec une partie qui démarre sans couper la musique du titre, ou qui
  enchaîne les grilles.
- **Le piano** entrera comme une partie de plus, à deux conditions que
  le modèle respecte déjà : un coup peut tomber ailleurs que sur le
  temps (les anticipations, le rythme de Charleston), et une partie
  devra recevoir l'accord du moment, que le repère d'un temps (`cue`)
  pourra porter à côté de la note de basse.

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

## Le premier jalon

Il tourne, dans `games/walk` : les fondamentales, main gauche seule, sur
un blues jazz lu dans un export iReal, sur la batterie et la contrebasse
en soundfont ; les marques de temps et de fondamentale sur la grille ;
la phase de déchiffrage sans tempo ; le bonhomme en bâtons ; la basse de
référence ; l'enregistrement (`-record`) ; le français et l'anglais.

Depuis, le jeu a gagné un écran titre et ses options, la calibration de
la latence, une musique qui ne s'arrête pas entre les écrans, les
bulles du bonhomme, trois grilles en ChordPro (le blues, *Tune Up*,
*Autumn Leaves*) choisies depuis la partie avec leur tempo, la grille
qui tourne ses pages, les accords avec leurs exposants, et le bilan de
fin de run.

La suite est dans `chantiers.md`, avant tout les paliers.

## Les questions ouvertes

- Le partage du clavier face aux voicings A/B (voir « Le clavier »).
- La liste des situations et des motifs à déclarer dans le dex.
- Le retour sur une note fausse, sans casser le groove.
- Le lore : le nom de l'univers, le ton, les titres des grilles.
