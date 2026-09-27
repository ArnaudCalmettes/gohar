# L'analyse des grilles

Comment gohar lit une grille : tonalités, cadences, emprunts,
modulations, et la gamme qui va sur chaque accord.

La progression suit le tome 1 d'*En Harmonie* (Dericq et Guéreau,
Outre Mesure), chapitres 8 à 10. Les exemples cités par titre en
viennent.

## À quoi elle sert

Le livre donne trois raisons d'analyser un morceau, et gohar les
reprend dans cet ordre :

1. **Jouer juste** : la gamme ou le mode qui va sur chaque accord.
   C'est la sortie la plus utile au joueur, et elle rejoint le dex, où
   gammes et modes sont des entrées.
2. **Modifier l'harmonie** en connaissance de cause : savoir ce que fait
   chaque accord, pour le remplacer ou le préparer autrement.
3. **Mémoriser** la grille : une suite de cadences se retient mieux
   qu'une suite d'accords.

Le public visé est l'autodidacte, avec des questions pratiques
immédiates. L'analyse lui dit ce que fait chaque accord ; elle ne lui
raconte pas l'histoire du morceau.

## Les principes

- **L'analyse propose, elle ne tranche pas.** Le mieux reste d'écouter,
  de relever et de connaître l'original. En l'absence de ces
  informations, gohar suggère ce qui se déduit « en toute logique », et
  le présente comme tel.
- **Déterministe, sans score.** Comme l'identification des accords :
  toutes les lectures valables sont rendues, dans un ordre fixé par des
  règles explicites. L'ordre dit laquelle est **proposée d'abord** dans
  ce contexte, jamais laquelle est vraie. Em7-A7 est un II-V de ré ou un
  III-VI7 de do : deux façons valables de se représenter la même
  progression, plus ou moins pertinentes selon le contexte, sans
  hiérarchie de niveau. Plusieurs outils peuvent ainsi expliquer le
  même phénomène, et chacun reste un calcul exact : l'analyse n'est pas
  une affaire d'interprétation, elle rend toutes les lectures vraies.
- **De droite à gauche.** Le livre le dit en toutes lettres : on repère
  les accords réels de la tonalité, on prend l'accord d'arrivée, puis on
  remonte ses préparations. C'est aussi l'algorithme. L'arrivée
  départage des lectures que rien d'autre ne départage : Fm7 B♭7 est un
  II-V de mi♭ ou un IVm7-♭VII7 de do, et c'est l'accord qui suit qui
  choisit.
- **La basse compte.** Elle distingue une cadence parfaite d'une
  imparfaite, un I/3 d'un III, un Dm9 d'un G13sus4 (mêmes notes, seule la
  basse change). Une lecture d'accord sans sa basse est incomplète.
- **La durée compte.** Elle sépare une plage modale d'un accord de
  passage, une modulation d'un emprunt. L'analyse travaille donc en
  temps, pas en cases iReal (voir `chantiers.md`).

## Les étages

Une analyse se lit du plus large au plus fin.

| Étage | Ce qu'il porte | Exemple |
|---|---|---|
| Le morceau | La structure, le rythme harmonique, la tonalité de départ et de fin | Tenderly : ABAC, 32 mesures, un accord par mesure, mi♭ majeur |
| La plage | Tonale, modale ou atonale | *One Finger Snap* : modale mesures 1 à 12, tonale ensuite |
| La région | Une tonalité installée : la principale, ou une modulation | *Black Orpheus* : la mineur, do majeur mesures 6 à 12, la mineur |
| Le bloc | Une cadence ou une préparation, avec la tonalité qu'elle annonce | Gm7♭5 C7(♭9), qui annonce fa mineur |
| L'accord | Son degré, sa fonction, sa provenance, sa couleur proposée | D♭7 : ♭VII7, emprunté à la♭ mineur mélodique |

La région se lit après coup, de droite à gauche. Elle a un pendant qui
se lit au présent, de gauche à droite : la tonique pressentie (voir
plus bas).

### Le morceau

La tonalité se déduit de l'armure et des cadences de début et de fin.
L'armure ne dit ni majeur ou mineur, ni s'il y a des modulations. Pour
une grille iReal, le champ de tonalité de l'app tient lieu d'armure et
précise même le mode (« A- ») : c'est un indice de départ, que les
cadences confirment ou démentent.

Le rythme harmonique est le nombre d'accords par mesure. La structure
(AABA, ABAC, AB, blues, forme anatole) se lit dans les sections de la
grille et se vérifie en comparant les progressions entre elles
(`Progression.Compare`).

### Les plages

- **Tonale** : les accords s'enchaînent autour d'un centre tonal et ont
  une fonction. C'est l'harmonie fonctionnelle, l'essentiel de ce
  document.
- **Modale** : un accord tenu plusieurs mesures installe un mode et non
  un centre tonal. C'est le X7sus4 « d'espèce » de *Maiden Voyage* : par
  sa durée, sans résolution sur le X7, il n'a pas de fonction. Critère
  mécanique : une durée sans cadence, au-delà d'un seuil.
- **Atonale** : pas de centre tonal, on passe d'un mode à l'autre au gré
  des accords (*Pee Wee*). Chaque accord reçoit sa provenance, aucun ne
  reçoit de degré. Dans `harmony`, c'est la valeur zéro de `Tonality` :
  un état légitime, pas un échec.

Un même morceau peut mêler les trois.

### Emprunt ou modulation

**L'emprunt** fait venir un accord ou une cadence d'une autre tonalité
sans quitter la sienne : Dm7 G7(♭13) Cmaj7 emprunte son G7(♭13) à do
mineur. **La modulation** change de tonalité pour de bon : une cadence
la prépare, et il faut que l'oreille ait le temps de l'entendre.

Le livre n'en donne pas une règle unique, et ses exemples montrent
trois indices qui se combinent :

1. **La cible est-elle un degré de la tonalité ?** Dans *There Will
   Never Be Another You*, Dm7♭5 G7(♭9) mène à Cm7, le VI de mi♭ : c'est
   une tonicisation, pas une modulation.
2. **Combien de temps la nouvelle tonique tient-elle ?** Tune Up change
   de tonalité toutes les quatre mesures, avec un II-V-I complet et un I
   tenu deux mesures : le livre parle de modulations.
3. **Combien de cadences la confirment ?** Dans *Black Orpheus*, do
   majeur a les mêmes notes que la mineur, mais plusieurs cadences
   l'installent sur sept mesures : c'est une modulation.

gohar calcule les trois indices et propose un classement. Les seuils
sont des **paramètres, en données**, et près d'un seuil les deux
lectures sont rendues.

### Les blocs

Un bloc est une cadence : **[II] [sus4] V → cible** (`analysis.Blocks`).
Le V est le seul accord obligatoire : il prépare la cible comme
dominante, dominante chromatique ou diminué. Le sus4 est le X7sus4 de
même fondamentale, le II celui qui prépare le V (ou son sus4). La cible
n'est pas dans le bloc : le bloc la désigne.

- **Découpage, de droite à gauche.** Chaque accord appartient à un
  seul bloc, et un bloc peut viser un accord d'un autre : A7 Dm7 G7
  Cmaj7 donne [Dm7 G7] → Cmaj7 et [A7] → Dm7, le V7/II puis le
  II-V-I. Le double rôle du bebop en découle : dans Em7 A7 Dm7 G7
  Cmaj7, Dm7 est la cible d'un bloc et le II du suivant.
- **Sans résolution.** Un II et un V qui ne prépare pas l'accord
  suivant font un bloc sans cible, le II-V sans résolution du livre,
  qui annonce quand même sa tonalité, une quinte sous le V. Un V seul
  qui ne prépare rien n'est pas un bloc.
- **La tonalité annoncée** est une tonique et des gammes, une
  `Tonality` par gamme (`Block.Announced`) : les gammes de la
  pratique tonale (majeure, mineure naturelle, harmonique, mélodique)
  qui contiennent toute la préparation, le II, le sus4 et le V. Un IIm7 garde ainsi le majeur
  et le mineur mélodique, un IIm7♭5 le mineur harmonique, un V7(♭13)
  les mineurs. La tierce de la cible départage ensuite, et seulement
  si la préparation le permet : Gm7 C7 → Fm7 annonce fa mineur
  mélodique, Gm7♭5 C7 → F reste fa mineur harmonique (le cas mixte).
  Une dominante chromatique et son II ne sont pas des degrés de la
  tonalité : la cible décide seule. La majeure harmonique contient
  aussi certaines préparations, mais c'est une gamme d'emprunt, pas
  une tonalité qu'on annonce ; elle reste dans la provenance.

Un bloc porte la tonalité qu'il **annonce**, et l'accord d'arrivée
porte la sienne. Dans *What Is This Thing Called Love*, le II-V porte
« fa mineur » et le I sur lequel il se résout porte fa majeur : c'est
le cas mixte.

Une chaîne de II-V est une suite de blocs, chacun avec sa tonalité
annoncée, et la marche est une relation entre eux : « II-V en mi
mineur, puis II-V en ré mineur, puis II-V en do mineur qui se résout sur
do majeur ».

Le livre appelle **Ier degré temporaire** l'arrivée d'un bloc emprunté
(Misty : A♭maj7 préparé par B♭m7 E♭7). C'est une **tonicisation** :
faire jouer à un accord le rôle de tonique secondaire.

## Les préparations

Préparer un accord cible, c'est ajouter ou modifier un ou deux accords
devant lui. L'anglais dit *approach* : dans le code, `ApproachKind`
nomme le type d'une préparation, et `ApproachOf` dit comment un accord
prépare le suivant. L'unité de l'analyse est donc une **arrivée** et
ce qui la prépare, chaque accord de préparation ayant un type.

| Type | Chiffrage | Exemple vers Dm7 en do |
|---|---|---|
| X7sus4 | retarde son propre X7 | A7sus4 A7 |
| Dominante secondaire | V7 de…, écrit V7/II ou VI7 | A7 |
| Dominante chromatique | ♭II7 de…, un X7 un demi-ton au-dessus | E♭7 |
| Accord diminué | un demi-ton sous l'arrivée | C♯dim7 |
| Accord parallèle | même qualité, une seconde à côté | E♭m7 ou C♯m7 |
| Sous-dominante secondaire | II-V de…, IIm7 ou IIm7♭5 | Em7 A7 |
| IVm7-♭VII7 de… | cadence modale préparée | Gm7 C7 |
| Sous-dominante chromatique | ♭VIm7 ou ♭VIm7♭5 devant le ♭II7 | B♭m7 E♭7, ou B♭m7 A7 |

**Les relations se composent.** La sous-dominante chromatique est un
II de… appliqué à une dominante chromatique ; la sous-dominante
secondaire, un II de… appliqué à une dominante secondaire. Le modèle
n'a donc que quelques relations élémentaires, qu'on enchaîne : le nom
composé sert à l'affichage, pas à la structure. Chez Ellington, G♭7 F7
E7 E♭7 A♭maj7 est une chaîne de V7 de… (C7 F7 B♭7 E♭7) dont un accord
sur deux est remplacé par son ♭II7 pour faire descendre la basse
chromatiquement.

**Ce qu'on en tire est une forme sous-jacente**, pas une histoire. Dans
*I Thought About You*, chaque X7 reçoit sa cible et la chaîne de II-V
par quintes réapparaît sous la grille écrite. C'est le résultat normal
de l'analyse de droite à gauche.

C'est ce que modélisent déjà `Target` et `Approach` dans `harmony` : une
cible, et le chemin qui y mène. La liste des approches s'élargit à ces
types. Le commentaire de `Resolution` disait que les règles sur
« jusqu'où un chemin peut s'éloigner » viendraient d'un musicien : ce
tableau en est la première version.

### Chaque X7 pose une question

Le livre la formule : est-il la dominante de la tonalité, la dominante
secondaire de l'accord suivant, ou sa dominante chromatique ? Ce sont
les lectures qu'une analyse rend pour un accord de septième de
dominante. Il peut aussi être un X7sus4 de passage, ou, tenu longtemps,
une plage modale.

### Le X7sus4

Sans tierce, il n'a pas de triton et n'est donc pas une dominante. En
cadence, c'est une **position d'attente** : il prolonge ou remplace le
II avec la fonction de sous-dominante, et se résout sur son X7 (tout X7
peut être précédé de son sus4, selon le tempo). Sans résolution et tenu
longtemps, il exprime une couleur et installe une plage modale.

### La dominante chromatique

Un X7 et celui situé un triton plus loin partagent leur triton, donc
toute cible a deux dominantes : `Resolution` le dérive déjà sans table.
Le livre en donne la raison profonde par l'accord de sixte augmentée de
l'harmonie classique : en cinq étapes, G7 devient D♭7, et à l'arrivée
**chaque voix rejoint l'accord de tonique par demi-ton**. C'est la règle
de moindre mouvement de `voicings.md` poussée à son maximum : quatre
conjonctions chromatiques sur quatre voix.

On dit « dominante chromatique » plutôt que « substitution
tritonique » (voir `glossaire.md`).

### L'accord diminué

Trois accords diminués seulement, puisque chacun a quatre fondamentales
possibles : `harmony` le calcule par symétrie. Deux emplois :

- **Dominante sans fondamentale.** Placé un demi-ton sous l'accord
  d'arrivée, le diminué est son V7(♭9) sans fondamentale : C♯dim7 porte
  la tierce, la quinte, la septième et la ♭9 de A7(♭9) et prépare Dm7.
  C'est la basse suivante qui désigne la fondamentale sous-entendue.
- **Accord de passage** entre deux accords diatoniques, pour une basse
  chromatique. Ce ne sont pas deux bits d'un même calcul : la dominante
  se voit sur deux accords (`ApproachOf`), le passage demande trois
  accords et leurs basses, que seule l'analyse d'une grille a. En
  montant (I ♯Idim7 II ♯IIdim7 III, *Mean to Me*), les deux emplois
  coïncident et les deux lectures sont rendues. En
  descendant, entre I/3 et II (B♭/D D♭dim7 Cm7), le diminué n'est pas la
  dominante de sa cible : il ne reste que le passage.

### Les accords parallèles

Sans fonction : ils harmonisent une note de la mélodie. Or une grille
iReal n'a pas la mélodie. gohar les repère par leur forme (même qualité
ou presque, une seconde à côté de l'arrivée, souvent brefs) et le
signale comme une conjecture.

## La provenance

La gamme d'où vient un accord, c'est ce que le joueur improvise dessus.
Elle se calcule, sans table, parmi **toutes les gammes nommées** qui
contiennent l'accord, et pas seulement celles de la tonique :
`harmony.Provenances`. Toutes les notes comptent, extensions
comprises : C7 vient de sept gammes, C7(♭9) de trois. Tenderly le
montre : D♭7 est emprunté à la♭ mineur mélodique, Gm7♭5 C7(♭9) à fa
mineur harmonique, Bdim7 à do mineur harmonique.

Ordre de présentation : d'abord la tonalité de la région, puis celle
qu'annonce le bloc, puis les autres. Plusieurs provenances sont
fréquentes et toutes sont rendues : le livre dit du F7 de Tenderly que
c'est « un IIe degré altéré pouvant correspondre à différents modes ».

Une **variante** de cadence est un squelette de degrés et une
provenance par degré. Les formes de l'anatole en fa (majeur, mineur
harmonique, mineur mélodique, et leurs mélanges) ne sont pas des
progressions distinctes : c'est I-VI-II-V, où chaque degré prend sa
tétrade dans l'une des gammes permises. Le catalogue stocke le
squelette et les gammes, `harmony` calcule les tétrades.

## Le chiffrage

Les degrés (`analysis.Degrees`) :

- Un degré diatonique s'écrit sans qualité, un degré emprunté avec :
  II, mais IVm7. Le degré se compte dans la gamme de la tonalité : en
  fa mineur, D♭maj7 est VI, pas ♭VI.
- Une fondamentale hors de la gamme est un degré abaissé de
  préférence (♭II, ♭III, ♭VI, ♭VII), haussé sous la quarte et la
  quinte (♯IV) ; un accord de passage suit sa basse, haussé en montant
  (♯Idim7, ♯Vdim7), abaissé en descendant (♭IIIdim7).
- Les renversements s'écrivent I/3, I/5.
- Un II-V se chiffre **relativement à sa cible**, dans la tonalité
  qu'il annonce, avec un crochet et une flèche vers elle, comme dans le
  livre. Un V seul, un diminué seul et un accord de passage se
  chiffrent dans la tonalité du morceau, avec leur qualité : VI7 pour
  la dominante secondaire de II.
- La tonalité du morceau est pour l'instant celle que donne l'app ; les
  régions et les modulations viendront avec les étages hauts.
- Une dominante secondaire s'écrit V7/II ou VI7. gohar stocke la
  relation (« V7 de Dm7 ») et rend l'une ou l'autre écriture : c'est un
  choix d'affichage.
- Une dominante chromatique se chiffre par son degré (♭II7 vers le I,
  ♯IV7 ou ♭III7 ailleurs), et la relation « ♭II7 de… » est gardée.

## Les couleurs proposées

La gamme de la cible colore sa préparation. La réalisation du livre
suit des règles régulières, qui deviennent des propositions :

- V7 vers un accord mineur : ♭13 ou ♭9 (A7(♭13) vers Dm7, E7(♭9) vers
  Am7).
- V7 vers un accord majeur ou de dominante : 9 et 13 (C13 vers Fmaj7,
  G13 vers Cmaj7).
- Dominante chromatique : ♯11 (F13(♯11) vers Em11 dans *But
  Beautiful*).
- Le I majeur : maj7 le plus souvent. Quand la mélodie est sur la
  fondamentale, deux façons d'éviter le frottement de la septième
  majeure : un accord 6, ou la neuvième majeure ajoutée par-dessus la
  septième.

Ce sont des propositions « en toute logique ». L'original, quand on le
connaît, l'emporte.

## Le catalogue v0

Des données, relues ligne à ligne, comme la table des qualités iReal.

### Les cadences

- **À deux accords** : parfaite (V-I, deux accords à l'état
  fondamental), imparfaite (V-I, au moins un renversement), demi-cadence
  (…-V), rompue (V-…, par exemple V-VI). La demi-cadence et la rompue se
  définissent par une absence : il faut la tonalité pour les voir.
- **Plagales** : IVmaj7-I ; IV7-I (mineur mélodique) ; IVm7-I (mineur
  harmonique) ; IV7-Im (mineur mélodique) ; IVm(maj7)-I (majeur
  harmonique) ; IV7-I7 (couleur « bluesy ») ; IV-IVm-I.
- **II-V-I** : majeur IIm7-V7-Imaj7 ; mineur harmonique IIm7♭5-V7(♭9,
  ♭13)-Im(maj7) ; mineur mélodique IIm7-V7(9, ♭13)-Im(maj7) ; mixte
  IIm7♭5-V7(9, ♭13)-Imaj7.
- **Modales ♭VII-I** : ♭VII7 (éolien), ♭VIImaj7 (mixolydien), ♭VIIm7
  (phrygien, plus rare) ; et préparées : II-♭VII7-I, IV-♭VII7-I,
  IVm7-♭VII7-I.
- **Les avatars du V** : V7, ♭II7, VIIdim7.

### Les sous-dominantes

Le tableau récapitulatif du chapitre 9 : IIm7, IIm7♭5, II7, ♭II7 ;
IVmaj7 ou IV6, IVm7 ou IVm6, IVm7♭5, IV7, ♯IVdim7 ; ♭VImaj7, VIm7♭5
(lu comme le ♯VIm7♭5 du mineur mélodique), ♭VI7 ; ♭VII7.

### Les cellules

Elles tournent au lieu d'aboutir, ce qui les distingue des cadences.

- **L'anatole** I-VI-II-V, en majeur, en mineur et mélangée. Cyclique,
  elle se joue sur une durée non définie : intro, coda. À ne pas
  confondre avec la forme anatole (les *rhythm changes*), qui est une
  structure de morceau.
- **III-VI-II-V-I**, variante où le III remplace le I (VI7 possible).
- **Les chaînes de II-V** par demi-ton, par ton, par cycle des quartes.
  Dans ce dernier cas, le I potentiel devient le II suivant (Cm7).
- **Le turnaround** : la fin qui prépare le retour au début, sur une ou
  deux mesures. Il traverse la barre finale, et le morceau ne commence
  pas toujours sur le I : l'analyse doit savoir que la forme boucle.

### Le I qui dure

- Ses couleurs : triade (rare), maj7, 6, et le mouvement oblique
  Imaj7-Imaj7(♯5)-I6.
- Ses prolongements : IVmaj7, IVm(maj7), IVm6/I, V7sus4, ♭VII7,
  ♭VIImaj7, I-II-III-V.
- La ligne gospel I-I7/3-IV-♯IVdim7-I/5.
- Les lignes chromatiques du Im : Im, Im(maj7), Im7, Im6 depuis la
  fondamentale, au soprano ou à la basse ; Im, Im(♭6), Im6, Im7 depuis
  la quinte.

## Les lectures qui se recouvrent

Cas où plusieurs lectures sont rendues, à garder comme tests :

- Em7-A7 : II-V de ré, ou III-VI7 de do.
- Fm7 B♭7 : II-V de mi♭, ou IVm7-♭VII7 de do. L'arrivée choisit.
- Dm7 avant D♭dim7 Cm7 en si♭ : Dm7, ou B♭/D (« souvent chiffré à tort
  Dm7 », dit le livre). De même A♭maj7/C souvent chiffré Cm7.
- Un X7 sur la basse à la quarte augmentée : D♭7/G est aussi G7(♭9, ♭5)
  (Debussy, *La plus que lente*).
- G13sus4 et Dm9 sur une basse sol.
- F♯dim7 avant Gm7 : accord de passage, ou D7(♭9) sans fondamentale.

## La fiche

Ce que l'analyse d'une grille produit, sur le modèle des analyses du
chapitre 10 :

- la structure et le nombre de mesures ;
- le rythme harmonique ;
- la tonalité de départ et de fin ;
- les plages, et les modulations avec leurs mesures ;
- les emprunts, chacun avec sa gamme d'origine ;
- les cadences et les préparations repérées, en degrés relatifs à leur
  cible ;
- pour chaque accord, la ou les gammes à jouer.

La fiche s'affiche sur la grille elle-même, annotée comme dans le
livre : degrés sous les accords, crochets et flèches des II-V,
tonalités au-dessus des blocs. Trois supports, dans cet ordre :

1. **Le terminal**, en texte : c'est le banc d'essai, qui grandit à
   chaque étape de l'analyse.
2. **Une fenêtre Ebitengine**, le jour où l'analyse est validée et où
   l'on joue en direct : une police de Real Book ou de MuseScore, du
   marqueur noir sur fond blanc, les chiffrages avec leurs indices et
   leurs exposants, le tout réglable. Cosmétique, donc après le cœur.
3. **Une page web**, le moteur compilé en WASM : on donne une grille,
   on reçoit l'analyse annotée. Pour distribuer et faire connaître le
   travail, dans la lignée de l'ancien gohareact ; le moins pressé.

## L'attente et la surprise

Un musicien à l'oreille entraînée entend « on dirait qu'on est en ré
majeur » tant que rien ne le dément, et sursaute quand une couleur belle
et inattendue le détrompe. L'analyse en direct doit faire la même
chose, au même instant : c'est un objectif, tant qu'on n'a pas prouvé
qu'il est impossible.

- **Les lectures provisoires** sont celles que ce qui a sonné permet :
  un II-V annonce son arrivée avant qu'elle ne sonne.
- **La surprise** est l'écart entre l'arrivée attendue et ce qui
  arrive, et elle se qualifie : cadence rompue, emprunt, dominante
  chromatique qui repart ailleurs. Ce n'est pas une erreur, c'est une
  couleur que la théorie sait nommer, et c'est ce qui fait la
  différence entre « faux » et « monstrueux ».
- **Pour un jeu**, c'est la récompense idéale : elle salue une prise de
  risque réussie, pas la conformité.

Une seule exigence en découle pour tout le code d'analyse : **chaque
calcul reste local**, il ne regarde qu'un nombre borné d'accords autour
de lui. L'analyse d'une grille et l'analyse en direct sont alors le
même code, sur une suite complète ou sur une suite qui s'allonge. Les
étages qui demandent de la durée (tonalité, modulation) restent
provisoires plus longtemps.

## La tonique pressentie

À chaque accord, la tonique que l'oreille attend, avec ce qui a sonné
et rien d'autre. Au troisième accord de Tenderly, deux mesures de
E♭maj7 A♭7 ont installé mi♭ : E♭m7 s'entend comme la tonique qui a
changé de couleur, pas comme le II de ré♭. Le II-V E♭m7 A♭7 ne
toniciserait ré♭ que s'il y arrivait.

C'est le pendant, au présent, de la région. Les deux lectures
coexistent : la fiche montre les régions, le direct montre la tonique
pressentie, et l'écart entre ce qu'elle attendait et ce qui arrive est
la surprise.

### Deux toniques

- **La tonique de fond** : celle qui est installée. Les degrés se
  comptent sur elle, comme le livre le fait même quand une cadence
  tonicise un autre degré (Dm7♭5 G7 Cm7 en mi♭ : II V VI).
- **La tonique locale** : celle qu'une cadence vient de toniciser, le
  temps de cette cadence. Elle ne change pas le fond.

Chacune est un ensemble de tonalités sur une même tonique, comme la
tonalité annoncée d'un bloc (mi♭ majeur, ou mi♭ majeur et mineur
mélodique quand rien ne tranche), ou vide quand rien n'est installé :
au début d'un morceau sans armure, dans une plage atonale.

### Ce qui installe une tonique

Du plus fort au plus faible :

1. **Une cadence qui se résout** : la cible devient tonique locale.
2. **La durée** : une tonique locale qui tient, confirmée par d'autres
   cadences, devient le fond. Ce sont les trois indices du livre
   (cible hors de la tonalité, durée, cadences qui confirment), lus au
   présent au lieu d'après coup ; leurs seuils sont en données.
3. **L'armure** : pour une grille iReal, le champ de tonalité de l'app
   donne un fond de départ, que la suite confirme ou dément.
4. **Le premier accord**, faute d'armure, s'il peut être un accord de
   tonique (maj7, 6, m6, m7, m(maj7)). C'est l'indice le plus faible :
   beaucoup de standards commencent sur un II ou un IV.

Ce qui ne change rien au fond : les accords diatoniques, les emprunts
sur la même tonique (Im7, IVm, ♭VII7), et les préparations qui ne se
résolvent pas.

### Le I emprunté

Un accord dont la fondamentale est la tonique de fond se lit comme un
I emprunté, même quand il est le II d'un II-V : E♭m7 dans E♭m7 A♭7 en
mi♭ est I, emprunté à l'éolien, et A♭7 redevient IV7. Le II-V reste un
bloc qui annonce ré♭ ; c'est sa lecture en degrés qui change.

### Tenderly, au présent

| Mesure | Accord | Tonique pressentie | Ce que l'oreille entend |
|---|---|---|---|
| 1 | E♭maj7 | mi♭ (armure, premier accord) | I |
| 2 | A♭7 | mi♭ | IV7, ou V de ré♭ : on attend peut-être ré♭ |
| 3 | E♭m7 | mi♭ | ré♭ n'arrive pas ; la tonique change de couleur : I emprunté |
| 4 | A♭7 | mi♭ | IV7 encore, en parallèle avec la mesure 2 |
| 5 | Fm7 | mi♭ | II |
| 6 | D♭7 | mi♭ | ♭VII7 : on attend mi♭ |
| 7 | E♭maj7 | mi♭, confirmée | I, par la cadence ♭VII7-I |
| 8 | Gm7♭5 C7♭9 | mi♭, locale fa mineur | II V de fa mineur : on attend Fm |
| 9 | Fm7♭5 | mi♭, locale mi♭ mineur | surprise : pas Fm, mais le II de mi♭ mineur |

### Local, avec une mémoire

La contrainte du direct (chaque calcul ne regarde qu'un nombre borné
d'accords) tient : la tonique pressentie est un état porté d'un accord
au suivant, un résumé de ce qui a sonné, pas un retour en arrière. Le
même code lit une grille complète en la parcourant de gauche à droite.

La tonalité du morceau s'en déduit : c'est le fond de l'accord final,
et quand un autre fond a dominé le morceau (le relatif majeur
d'*Autumn Leaves*), les deux sont rendus.

### Tranché

- **Le relatif qui tonicise d'abord.** *Autumn Leaves* en sol mineur
  commence par Cm7 F7 B♭maj7 : la première cadence installe si♭, et sol
  mineur n'arrive qu'après. C'est un cas ambigu notoire, et les deux
  lectures sont acceptables de loin : en direct, le fond passe de si♭ à
  sol mineur quand sol mineur s'installe ; après coup, les deux sont
  rendues, et la tonalité du morceau est celle de l'accord final, sol
  mineur.
- **La tonique de départ est gardée en mémoire**, à part du fond, pour
  reconnaître le retour à la maison après un pont qui a modulé : c'est
  très souvent le cas sur une forme AABA.
- **Le tableau de Tenderly** est validé en attendant l'avis d'une
  oreille plus experte : il sert de test.

### À trancher

- Quand une tonique locale tient-elle assez pour devenir le fond ?
  Tune Up et Black Orpheus serviront à régler les seuils.

## Les tests de référence

Le livre analyse des morceaux qui sont dans le corpus iReal. Sa fiche
devient l'oracle : on compare celle de gohar à celle des auteurs.

- **Tune Up** : AA', 16 mesures, un accord par mesure, ré majeur,
  modulations en do (5 à 8) et en si♭ (9 à 12).
- **Black Orpheus** : AB, 32 mesures, un ou deux accords par mesure, la
  mineur, modulation en do majeur (6 à 12) puis retour.
- **There Will Never Be Another You** : ABAC, 32 mesures, mi♭ majeur,
  II-V vers le VI et vers le IV, ♭VII7-I, sans modulation.
- **Tenderly**, mesures 1 à 16 : ABAC, 32 mesures, un accord par
  mesure, mi♭ majeur sans modulation, huit emprunts avec leur gamme,
  plagales, ♭VII7-Imaj7, Bdim7 de passage (aussi G7(♭9) sans
  fondamentale), II7 en marche IIm7-V7.

Les fiches sont transcrites à la main dans
`charts/ireal/testdata/fiches`, en données de test, pas dans le code.
Elles montrent déjà ce que l'analyse devra encaisser : l'app donne
Tune Up en si♭ et le joue sur 32 mesures avec deux fins, là où le livre
l'analyse en ré majeur sur 16. Le champ de tonalité n'est qu'un
indice.

## Hors périmètre

- **Accord modifié ou ajouté.** Savoir si un A7 remplace un Am7 de la
  grille d'origine ou s'y ajoute relève presque du travail de
  l'historien du jazz. Les joueurs visés ont d'abord des questions
  pratiques, et ceux qui veulent cette culture la trouveront auprès de
  professeurs qui en sont des puits. L'analyse dit ce que fait chaque
  accord, jamais d'où il vient ; l'édition d'une grille se contente de
  réanalyser la nouvelle version.
- **Ce qui demande la mélodie**, tant qu'on ne l'a pas : repérer à
  coup sûr un accord parallèle, ou refuser une dominante chromatique
  quand la mélodie tient la ♯11 du X7 (elle deviendrait la fondamentale
  du nouvel accord, *Midnight Sun*). Ces règles iront dans les
  propositions de réharmonisation, le jour où une grille aura sa
  mélodie.
- **L'analyse rythmique**, que le livre écarte aussi.

## Ouvert

- Les seuils : durée d'une plage modale, indices d'une modulation. À
  régler sur les tests de référence, puis à l'oreille.
- La longueur d'une chaîne de préparations avant qu'elle ne soit plus
  une préparation mais une région.
- L'ordre exact des lectures quand plusieurs valent : à écrire règle
  par règle, et à confronter aux fiches du livre.
- **Le II-V sans résolution, contre le livre.** Les degrés concordent
  avec les fiches du livre à 70 sur 83, et les écarts sont de deux
  sortes. Les modulations (Tune Up, Black Orpheus) attendent les
  régions. Les autres sont des IIm7 V7 qui ne se résolvent pas, là où
  le livre entend autre chose. Deux cas, à reprendre une fois la
  tonalité du morceau détectée :
  - **Le II qui est un I emprunté** (cas limite, ouvert). E♭m7 A♭7
    dans Tenderly, mesures 3-4 : nous lisons un II-V de ré♭ qui ne se
    résout pas, le livre lit I emprunté à l'éolien, puis IV7, en
    parallèle avec E♭maj7 A♭7 des mesures 1-2. La clé : le « II » a
    pour fondamentale la tonique pressentie, installée par les deux
    premières mesures. La règle en découle (voir « Le I emprunté ») et
    reste compatible avec le direct. À vérifier en la codant : qu'elle
    ne casse rien, et qu'elle permet d'affiner proprement l'analyse.
  - **La marche IIm7-V7.** Cm7 F7 avant Fm7 (Tenderly mesures 13-14,
    There Will Never Be Another You mesures 12-13) : le livre lit VI,
    puis II7. Ici le II n'est pas sur la tonique ; la lecture du livre
    viendra plutôt des cadences du catalogue ou des relations entre
    blocs.
  - Écartée : départager par la qualité du II (m7 contre m7♭5). Elle
    colle aux quatre cas du corpus mais ne repose sur aucune raison
    musicale solide.