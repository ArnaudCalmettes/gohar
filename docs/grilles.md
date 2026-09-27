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
  hiérarchie de niveau.
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
devant lui. L'unité de l'analyse est donc une **arrivée** et ce qui la
prépare, chaque accord de préparation ayant un type.

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
  chromatique. En montant (I ♯Idim7 II ♯IIdim7 III, *Mean to Me*), les
  deux emplois coïncident et les deux lectures sont rendues. En
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
contiennent l'accord sur ce degré, et pas seulement celles de la
tonique. Tenderly le montre : D♭7 est emprunté à la♭ mineur mélodique,
Gm7♭5 C7(♭9) à fa mineur harmonique, Bdim7 à do mineur harmonique.

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

- Un degré diatonique s'écrit sans qualité, un degré emprunté avec :
  II, mais IVm7.
- Les renversements s'écrivent I/3, I/5.
- Un II-V se chiffre **relativement à sa cible**, avec un crochet et une
  flèche vers elle, comme dans le livre.
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

Les fiches se transcrivent à la main depuis le livre, dans des données
de test, pas dans le code.

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
