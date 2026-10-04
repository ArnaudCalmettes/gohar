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
5. En fin de run, ce qui a tenu et ce qui a lâché, par situation.
6. Les marques partent au dex, qui décide de la suite.

**La phase de déchiffrage, sans tempo** : la grille attend le joueur et
avance quand il a joué ce qu'on attend, comme le mode « attente » des
jeux de piano. Les marques de note restent, celles de temps
disparaissent. Elle dure autant qu'il le souhaite.

## Le clavier

- Deux zones, partagées par défaut au **sol2** (MIDI 55, G3 en notation
  américaine), le partage réglable (`-split`). Pas de mains croisées :
  le partage suffit toujours.
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

## La basse de référence

- Au premier passage d'une grille, le jeu joue une vraie ligne de
  basse ; au second, le joueur prend la main. Le générateur et le
  marqueur partagent leurs règles : une ligne générée doit recevoir de
  bonnes marques, et les tests le vérifient.
- Le générateur (`generator.go`, joué par `-demo`) suit les cinq règles
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
  (`walker.go`), qui suffit à valider que la démarche suit le jeu avant
  de commander le moindre dessin :
  - l'aisance est une moyenne glissante des arrivées, à peu près les
    huit dernières ; il saute à chaque arrivée posée, et ne trébuche
    sur une manquée que si l'aisance est déjà retombée ;
  - deux images, comme un cycle de sprite, jamais d'entre-deux : une
    animation fluide colle mal au pas. Sur les temps forts, les jambes
    écartées ; sur les temps faibles, les jambes qui se croisent ;
  - le claquement de doigts à la manière jazz, du coude ;
  - les seuils sont des valeurs de départ, à régler en jouant.
- **Plus tard, il parle**, avec les mots du jeu (des marques, pas un
  jugement) :
  - pendant le jeu, quand un motif se dessine dans les marques : « tu
    presses » sur une série de notes en avance, « tu traînes » sur une
    série en retard, un mot quand ça tourne ou quand le joueur se
    rattrape. Chaque bulle attend un moment avant de pouvoir revenir,
    pour qu'il ne radote pas. Elles supposent la latence calibrée :
    sinon, elle passerait pour de la précipitation ;
  - à la fin du run, ce qui a tenu et ce qui a lâché, par situation
    (voir « La boucle de jeu », point 5).
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

- Ebitengine, comme `games/ear`, et Ark si un ECS se justifie.
- Une écriture « marqueur », noire sur fond blanc, comme sur un Real
  Book : les accords en **MuseJazz Text**, la police de MuseScore
  (licence SIL OFL 1.1), qui a ♭ ♮ ♯, ø et °. Elle n'a pas de crénage,
  son fichier source ayant été perdu : les paires délicates se règlent
  à la main.
- Le jeu choisit son `ChordStyle`, et les noms d'accords suivent
  l'orthographe de l'analyse (voir « L'orthographe entendue » dans
  `grilles.md`).
- **La mise en page** : la grille sur les deux tiers droits de l'écran,
  quatre mesures par ligne ; le bonhomme dans le tiers gauche ; le
  clavier en dessous, sur toute la largeur.

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
  kit n'en ayant. Le snap est versionné dans `games/walk/sounds` ;
  GeneralUser, 32 Mo, ne l'est pas : `make sounds` la télécharge dans
  le cache de l'utilisateur, d'un commit figé, empreinte vérifiée.

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

La suite est dans `chantiers.md` : la calibration, les bulles du
bonhomme, l'affichage de grilles plus longues, puis les paliers.

## Les questions ouvertes

- Le partage du clavier face aux voicings A/B (voir « Le clavier »).
- La liste des situations et des motifs à déclarer dans le dex.
- Le retour sur une note fausse, sans casser le groove.
- Le lore : le nom de l'univers, le ton, les titres des grilles.
