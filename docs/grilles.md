# L'analyse des grilles

Comment gohar lit une grille d'accords : ce que fait chaque accord, les
cadences qui les relient, la tonalité qu'ils finissent par installer,
et la gamme qui va sur chacun.

Le fil est celui d'*En Harmonie* (Dericq et Guéreau,
Outre Mesure), tome 1, chapitres 8 à 10, et « le livre » désigne ce
tome-là sauf mention du tome 2. Le code vit dans `harmony/analysis` ;
`charts/cmd/analyse` affiche une grille annotée dans le terminal, et
`charts/cmd/corpus` compare l'analyse à l'app iReal sur des playlists
entières.

Chaque partie suit le même mouvement, pour qu'on sache toujours où on
en est : ce que dit le livre, comment gohar le modélise, ce qui est
difficile, ce qu'on a décidé. Quand un temps n'a rien à dire, on le
saute.

## 1. À quoi sert l'analyse

### Dans le livre

Le livre donne trois raisons d'analyser un morceau, et gohar les
reprend dans cet ordre :

- **jouer juste** : savoir quelle gamme ou quel mode va sur chaque
  accord, ce qui rejoint le dex ;
- **modifier l'harmonie** en connaissance de cause, en sachant ce que
  fait chaque accord avant d'en changer un ;
- **mémoriser** la grille, parce qu'une suite de cadences se retient
  mieux qu'une suite d'accords : *Autumn Leaves*, ce sont deux II-V-I
  et un turnaround, pas trente-deux mesures à apprendre par cœur.

Le public visé est l'autodidacte, avec des questions pratiques
immédiates. L'analyse lui dit ce que fait chaque accord ; elle ne lui
raconte pas l'histoire du morceau.

### Dans gohar : cinq principes

- **L'analyse propose, elle ne tranche pas.** Le mieux reste d'écouter,
  de relever et de connaître l'original. À défaut, gohar suggère ce qui
  se déduit « en toute logique », et le présente comme tel.
- **Toutes les lectures valables sont rendues.** Em7 A7 est un II-V de
  ré, ou un III-VI7 de do : deux façons exactes de lire la même paire
  d'accords, et l'analyse donne les deux, dans un ordre fixé par des
  règles explicites. L'ordre dit laquelle est proposée d'abord, jamais
  laquelle est vraie. Il n'y a pas de score, pas de probabilité : les
  mêmes accords donnent toujours la même analyse.
- **De droite à gauche.** C'est la méthode du livre, et c'est aussi
  l'algorithme : on repère les accords qui appartiennent à la
  tonalité, on prend l'accord d'arrivée, puis on remonte les accords
  qui le préparent. L'arrivée départage ce que rien d'autre ne
  départage : Fm7 B♭7 est un II-V de mi♭ s'il va sur E♭maj7, et un
  IVm7 ♭VII7 de do s'il va sur Cmaj7.
- **La basse compte.** Elle distingue une cadence parfaite d'une
  imparfaite, un C/E (le I avec sa tierce à la basse) d'un Em7 (le
  III), un Dm9 d'un G13sus4 (les mêmes notes, sur une basse ré ou
  sol).
- **La durée compte.** Elle sépare un accord de passage d'une plage
  modale, une tonicisation d'une modulation : le Dm tenu deux mesures
  dans le pont de *Black Orpheus* n'est pas le Dm7 d'un II-V. L'analyse
  travaille donc en temps, pas en cases iReal.

### Les cinq étages

Une grille se lit à cinq hauteurs, du morceau entier à l'accord seul.
Ce tableau est la carte du document : les parties qui suivent prennent
les étages de bas en haut, de l'accord au morceau, parce que c'est
l'ordre dans lequel on comprend.

| Étage | Ce qu'il porte | Exemple |
|---|---|---|
| Le morceau | La structure, le rythme harmonique, la tonalité | Tenderly : ABAC, 32 mesures, un accord par mesure, mi♭ majeur |
| La plage | Tonale, modale ou atonale | *One Finger Snap* : modale mesures 1 à 12, tonale ensuite |
| La région | Une tonalité installée : la principale, ou une modulation | *Black Orpheus* : la mineur, do majeur mesures 6 à 12, la mineur |
| Le bloc | Une cadence, avec la tonalité qu'elle annonce | Gm7♭5 C7(♭9), qui annonce fa mineur |
| L'accord | Son degré, sa fonction, sa provenance, sa couleur proposée | D♭7 : ♭VII7, emprunté à la♭ mineur mélodique |

La région se lit après coup, de droite à gauche. Elle a un pendant qui
se lit au présent, de gauche à droite : la tonique pressentie, ce que
l'oreille attend accord par accord. Les deux ont leur place dans la
partie sur la tonalité.

## 2. L'accord : ce qu'il fait

Le plus petit étage, et le plus sûr : ce qu'on lit sur un accord seul
ne dépend presque de rien d'autre.

### Le chiffrage et la basse

#### Dans le livre

Un chiffrage iReal dit une fondamentale, une qualité et parfois une
basse : Dm7, G7, C/E. Le livre en lit plus que ce qui est écrit. Il
relève les chiffrages trompeurs, comme le Dm7 « souvent chiffré à
tort » là où sonne un B♭/D (le I de si♭ avec sa tierce à la basse), ou
le Cm7 qui est un A♭maj7/C. Et il tient à la basse : G7 Cmaj7 est une
cadence parfaite, G7/B Cmaj7 une imparfaite, parce que la basse ne
saute plus de quinte.

#### Dans gohar

L'analyse lit chaque accord comme une fondamentale, un motif de notes et
une basse, et garde les lectures qui se recouvrent, à rendre toutes :

- Em7 A7 : un II-V de ré, ou un III-VI7 de do ;
- Fm7 B♭7 : un II-V de mi♭, ou un IVm7 ♭VII7 de do, selon l'arrivée ;
- Dm7 avant D♭dim7 Cm7 en si♭ : un B♭/D, le I sur sa tierce, malgré le
  chiffrage ;
- D♭7/G : aussi un G7(♭9, ♭5), le même accord vu de l'autre côté du
  triton (Debussy s'en amuse dans *La plus que lente*) ;
- G13sus4 et Dm9 : les mêmes notes sur une basse sol ou ré ;
- F♯dim7 avant Gm7 : un accord de passage, ou un D7(♭9) sans sa
  fondamentale.

Ces cas sont gardés comme tests.

### Le degré

#### Dans le livre

Le livre chiffre chaque accord par son degré dans la tonalité du
passage, en chiffres romains, et compte ce degré dans la gamme de la
tonalité, même en mineur : dans ses tableaux VI V I, G∆9 est le VI de
si mineur harmonique, pas un ♭VI. Un II-V qui vise un autre accord que
la tonique se chiffre entre crochets, relativement à sa cible : en
mi♭, Gm7♭5 C7(♭9) avant Fm7 s'imprime « II V » sous un crochet marqué
fa mineur.

#### Dans gohar

Deux lectures, rendues toutes les deux, parce qu'elles servent à des
choses différentes.

**Sur la tonalité du passage**, chaque accord par son degré dans la
région où il sonne : en mi♭, le même Gm7♭5 C7♭9 est IIIm7♭5 VI7, le
III-VI d'un III-VI-II-V-I. C'est la lecture qu'affiche `analyse` sous
les accords.

**En crochets**, comme le livre les imprime : le II-V se chiffre dans
la tonalité qu'il annonce, II V de fa mineur. `analyse` la donne sur la
ligne des blocs, qui tient lieu de crochet, quand elle diffère de la
première. Les fiches du livre se comparent à celle-là. Une plagale ne
se lit pas en crochets : son IVm7 ♭VII7 se chiffre sur la tonique
qu'il conclut.

Les règles d'écriture :

- un degré diatonique s'écrit sans qualité (II, V), un degré emprunté
  avec (IVm7 en majeur) ;
- une septième de dominante ailleurs que sur le V garde sa qualité
  même diatonique, parce que ce n'est pas *la* dominante de la gamme :
  en la mineur, G7 est VII7 ;
- une fondamentale hors de la gamme est abaissée de préférence (♭II,
  ♭III, ♭VI, ♭VII) et haussée sous la quarte et la quinte (♯IV) ;
- un accord de passage suit sa basse, un sus4 garde sa qualité (V7sus4),
  les renversements s'écrivent I/3, I/5 ;
- une dominante secondaire se chiffre par son degré, VI7 pour celle du
  II, que d'autres écrivent V7/II : gohar garde la relation, l'écriture
  est un choix d'affichage ;
- une dominante chromatique aussi, ♭II7 vers le I, ♯IV7 ou ♭III7
  ailleurs ;
- un 7alt s'écrit 7alt, le locrien ♭4 dit en toutes lettres (VI7alt,
  pas VI7♭5).

#### Ce qu'on a décidé

Compter les degrés dans la gamme de la tonalité, y compris en mineur,
comme le livre : en fa mineur, D♭maj7 est VI, pas ♭VI. L'écriture 7alt
est un choix fait faute de source qui le tranche, à corriger par un
expert.

### La provenance

#### Dans le livre

Un accord vient d'une gamme, et cette gamme est ce qu'on joue dessus.
Le livre l'écrit pour chaque emprunt : dans Tenderly, D♭7 est
« emprunté à la♭ mineur mélodique », et du F7 il dit que c'est « un
IIe degré altéré pouvant correspondre à différents modes ».

#### Dans gohar

La provenance se calcule sans table, parmi toutes les gammes nommées qui
contiennent l'accord, extensions comprises : C7 vient de sept gammes,
C7(♭9) de trois. Toutes sont rendues, dans un ordre qui met d'abord la
tonalité de la région, puis celle qu'annonce le bloc, puis les autres.
Dans Tenderly, D♭7 est emprunté à la♭ mineur mélodique, Gm7♭5 C7(♭9) à
fa mineur harmonique, Bdim7 à do mineur harmonique.

Une **variante** de cadence est un squelette de degrés avec une
provenance par degré : les formes de l'anatole en fa que donne le livre
(majeure, mineures, mélangées) ne sont pas des progressions distinctes,
c'est I VI II V où chaque degré prend sa tétrade dans l'une des gammes
permises.

Repères : `Degrees`, `Bracketed`, `harmony.Provenances`.

## 3. La cadence : le bloc, et ce qu'il annonce

Une cadence, c'est un ou deux accords qui en préparent un autre, et
l'arrivée qu'ils préparent. C'est l'unité de l'analyse : le livre lit
une grille cadence par cadence, et gohar aussi.

### Dans le livre

**Chaque X7 pose une question**, que le livre formule ainsi : est-il la
dominante de la tonalité, la dominante secondaire de l'accord suivant,
ou sa dominante chromatique ? Il peut aussi être un X7sus4 de passage,
une septième d'espèce (le blues), ou, tenu longtemps, une plage modale.
Répondre à cette question pour chaque X7, c'est déjà la moitié de
l'analyse.

Le livre range les cadences par familles :

| Cadences | Formes |
|---|---|
| À deux accords | parfaite (V-I à l'état fondamental), imparfaite (un renversement), demi-cadence (…-V), rompue (V-…, V-VI par exemple) |
| II-V-I | majeur IIm7-V7-Imaj7 ; mineur harmonique IIm7♭5-V7(♭9, ♭13)-Im(maj7) ; mineur mélodique IIm7-V7(9, ♭13)-Im(maj7) ; mixte IIm7♭5-V7(9, ♭13)-Imaj7 |
| Plagales | IVmaj7-I ; IV7-I et IV7-Im (mineur mélodique) ; IVm7-I (mineur harmonique) ; IVm(maj7)-I (majeur harmonique) ; IV-IVm-I ; IV7-I7, « bluesy » |
| Modales ♭VII-I | ♭VII7 (éolien), ♭VIImaj7 (dorien vers Im7, mixolydien vers I7), ♭VIIm7 (phrygien) ; préparées II-♭VII7-I, IV-♭VII7-I, IVm7-♭VII7-I (tome 2) |
| Avatars du V | V7, ♭II7, VIIdim7 |

Et les accords de sous-dominante, d'après le tableau du chapitre 9 :
IIm7, IIm7♭5, II7, ♭II7 ; IVmaj7 ou IV6, IVm7 ou IVm6, IVm7♭5, IV7,
♯IVdim7 ; ♭VImaj7, VIm7♭5, ♭VI7 ; ♭VII7.

**Le I qui dure** a ses couleurs (triade, maj7, 6, le mouvement oblique
Imaj7-Imaj7(♯5)-I6), ses prolongements (IVmaj7, IVm(maj7), IVm6/I,
V7sus4, ♭VII7, ♭VIImaj7, I-II-III-V), la ligne gospel
I-I7/3-IV-♯IVdim7-I/5, et les lignes chromatiques du Im (Im, Im(maj7),
Im7, Im6 depuis la fondamentale ; Im, Im(♭6), Im6, Im7 depuis la
quinte).

Deux définitions du livre, pour la suite :

- **La demi-cadence** (« … – V ») « consiste à enchaîner un accord sur
  la fonction de dominante. Cette cadence, au caractère suspensif, se
  termine sur une tension qui laisse entrevoir une suite », « un point
  d'interrogation dans une phrase » : le D7 de la mesure 8 d'*It Don't
  Mean A Thing*.
- **La cadence rompue** (« V – … »), ou cadence évitée quand il y a
  modulation, « consiste à enchaîner l'accord de dominante avec un
  autre accord que l'accord de tonique prioritairement attendu, créant
  un effet de surprise » : C7 Dm7, ou B7 Cmaj7 dans *Nardis*.

Sur les dominantes chromatiques, le livre donne la règle et sa raison :
« tout accord X7 peut être substitué à un autre accord X7, quelle que
soit sa fonction », le G7 par le D♭7 qui partage son triton, et il le
démontre par l'accord de sixte augmentée, où chaque voix rejoint
l'accord de tonique par demi-ton. Il dit aussi comment retrouver ce qui
est sous les substitutions : « l'analyse doit s'effectuer de la droite
vers la gauche », en prenant l'accord d'arrivée et en le préparant par
des dominantes en quintes. Sous le G♭7 F7 E7 E♭7 A♭maj7 de *Sophisticated
Lady*, on retrouve C7 F7 B♭7 E♭7.

### Dans gohar : les préparations

Préparer un accord, c'est ajouter ou modifier un ou deux accords devant
lui. L'anglais dit *approach*, et c'est le mot du code. Pour chaque
accord, l'analyse dit comment il prépare le suivant, parmi ces types :

| Type | Chiffrage | Exemple vers Dm7 en do |
|---|---|---|
| Dominante secondaire | V7 de…, écrit V7/II ou VI7 | A7 |
| Dominante chromatique | ♭II7 de…, un X7 un demi-ton au-dessus | E♭7 |
| Accord diminué | un demi-ton sous l'arrivée | C♯dim7 |
| X7sus4 | retarde son propre X7 | A7sus4 A7 |
| Sous-dominante secondaire | II-V de…, IIm7 ou IIm7♭5 | Em7 A7 |
| Sous-dominante chromatique | ♭VIm7 ou ♭VIm7♭5 devant le ♭II7 | B♭m7 E♭7, ou B♭m7 A7 |
| Plagale | le IV de toute qualité, ou le ♭VII7, devant une tonique | Gm7 Dm, C7 Dm |
| IVm7-♭VII7 de… | la plagale préparée | Gm7 C7 Dm |

Il manque à ce tableau l'accord parallèle (même qualité, une seconde à
côté : E♭m7 ou C♯m7 vers Dm7), qui demande la mélodie pour être sûr.

**Les relations se composent.** La sous-dominante chromatique est un II
de… appliqué à une dominante chromatique, la sous-dominante secondaire
un II de… appliqué à une dominante secondaire. Le modèle n'a que
quelques relations élémentaires qu'on enchaîne, et le nom composé sert à
l'affichage. Chez Ellington, G♭7 F7 E7 E♭7 A♭maj7 est une chaîne de V7
de… dont un accord sur deux est remplacé par son ♭II7 pour faire
descendre la basse. Ce qu'on en tire est une **forme sous-jacente**, pas
une histoire : dans *I Thought About You*, chaque X7 reçoit sa cible et
la chaîne de II-V par quintes réapparaît sous la grille écrite.

**Le X7sus4**, sans tierce, n'a pas de triton et n'est pas une
dominante. C'est un accord de sous-dominante, « un retard de la
sous-dominante », dit le tome 2. En cadence, il prolonge ou remplace
le II et se résout sur son X7 ; il se chiffre V7sus4, jamais V, et un
II-V qui tombe dessus n'est pas mis entre crochets, le sus étant un II
qui se cache sous la fondamentale du V suivant (Dm7 G7 C7sus4 C7 F6
dans *My Lucky Star*, un III-VI-II-V de fa). Sans résolution et tenu
longtemps, il installe une plage modale.

**La dominante chromatique** se dérive sans table : toute cible a deux
dominantes, celle qui est une quinte au-dessus et celle qui est un
demi-ton au-dessus, qui partagent le même triton. On dit « dominante
chromatique » plutôt que « substitution tritonique », parce que c'est
le nom du livre pour l'accord, la substitution étant l'opération.

**L'accord diminué** a deux emplois. Un demi-ton sous l'arrivée, c'est
son V7(♭9) sans fondamentale : C♯dim7 prépare Dm7, et comme chaque
diminué a quatre fondamentales, il prépare aussi toute cible un
demi-ton au-dessus de l'une de ses notes. Entre deux accords, c'est un
**accord de passage** sur une basse chromatique, que seule l'analyse
d'une grille voit, puisqu'il faut trois accords et leurs basses ; il
est alors nommé d'après sa basse, haussé en montant, abaissé en
descendant. En montant (I ♯Idim7 II, *Mean to Me*), les deux emplois
coïncident et les deux lectures sont rendues ; en descendant (B♭/D
D♭dim7 Cm7), il ne reste que le passage.

**La plagale** conclut autant qu'un V-I : le F/C C final de *My Way*
est un gros amen sur do. Mais elle attire moins : un V-I peut faire
changer de tonalité, un IV-I ne le fait pas. Le ♭VII7-I en est le faux
nez mineur : avec le IV à la basse, B♭7 devient Fm6 en do. Quand un
♭VImaj7 la précède, ♭VImaj7 ♭VII7 Imaj7, c'est la cadence modale de do
éolien du tome 2, une cellule (voir plus bas).

### Dans gohar : les blocs

Un bloc est une cadence : **[II] [sus4] V → cible**. Le V est le seul
accord obligatoire, et prépare la cible comme dominante, dominante
chromatique ou diminué ; une plagale fait aussi un bloc, sa
sous-dominante à la place du V, avec son IVm7 pour II quand c'est le
♭VII7. La cible n'est pas dans le bloc : il la désigne.

Le découpage se fait de droite à gauche. Chaque accord appartient à un
seul bloc, et un bloc peut viser un accord d'un autre : A7 Dm7 G7 Cmaj7
donne [Dm7 G7] → Cmaj7 et [A7] → Dm7, le V7/II puis le II-V-I, et le
double rôle du bebop en découle (Dm7 cible d'un bloc et II du suivant).
Un II et un V qui ne préparent pas l'accord suivant font un bloc sans
cible, le **II-V sans résolution** du livre, qui annonce quand même sa
tonalité ; un V seul qui ne prépare rien n'en est pas un.

**La tonalité annoncée** est une tonique et les gammes de la pratique
tonale (majeure, mineures naturelle, harmonique, mélodique) qui
contiennent toute la préparation : un IIm7 garde le majeur et le mineur
mélodique, un IIm7♭5 le mineur harmonique, un V7(♭13) les mineurs. La
tierce de la cible départage ensuite, si la préparation le permet : Gm7
C7 → Fm7 annonce fa mineur mélodique, Gm7♭5 C7 → F reste fa mineur
harmonique, le cas mixte de *What Is This Thing Called Love*, où le
bloc porte fa mineur et le I porte fa majeur. Une dominante chromatique
et son II ne sont pas des degrés de la tonalité : la cible décide
seule. Quand les couleurs écrites n'entrent dans aucune gamme, les
tétrades seules sont essayées, puis, pour le V, son seul triton : une
dominante altérée (C7alt, C7♭5) prend ses couleurs hors de la tonalité
qu'elle annonce, et reste le V de fa (Gm7♭5 C7alt Fm7♭5 dans *One
Finger Snap*).

**Le Ier degré temporaire** du livre, l'arrivée d'un bloc emprunté
(*Misty* : A♭maj7 préparé par B♭m7 E♭7), est une tonicisation. Le
premier pas d'une **marche IIm7-V7**, dont le V7 devient le II du bloc
suivant sur la même fondamentale, n'est pas mis entre crochets : Cm7 F7
Fm7 B♭7 en mi♭ (Tenderly mesures 13 à 16) se lit VI II7 II V, comme le
livre. F7 est un « IIe degré altéré », pas le V d'un si♭ qui ne vient
jamais ; le bloc annonce quand même si♭ au-dessus des accords.

### Ce qui est difficile

**Le bloc seul ne suffit pas.** Rien dans G7 Am7 ne dit si c'est le
♭VII7 Im de la mineur, une plagale, ou le V VI de do, une cadence
rompue. Rien dans G7 Dm7 ne dit si c'est le IV7 Im dorien de ré mineur
ou le V de do qui revient sur son II. Seule la tonique qu'on entend au
moment du V tranche, et le bloc est lu avant qu'on l'entende. C'est
pour cela que l'analyse écoute deux fois (voir « La tonalité »).

**Le V-VI n'a pas toujours son II.** Dans *That Old Feeling*, Fm7 G7
Am7 : G7 va sur son VI sans II devant lui. Un V seul qui ne prépare
rien n'est pas un bloc, mais celui dont la tonique est entendue ne va
pas nulle part, il la promet.

**La demi-cadence** est une phrase qui s'arrête sur le V, et le D7 de la
mesure 8 d'*It Don't Mean A Thing* se résout pourtant au retour du A.
Ce n'est donc pas un bloc sans cible, c'est une affaire de structure :
la même section revient avec une fin différente, la question sur le V
puis la réponse sur le I (voir « Le morceau »).

**Le IV7-I7 du blues** n'est pas une plagale : dans un blues, le I7 est
le fond, et rien ne conclut. Le livre range pourtant IV7-I7 parmi les
plagales, avec *Black Coffee*, où le Imaj7 « peut être modifié en I7
afin de donner une couleur bluesy » : un morceau tonal qui emprunte la
couleur du blues. Ce cas-là reste ouvert.

### Ce qu'on a décidé

- **Toutes les cadences du tableau sont codées**, sauf le IV7-I7, la
  demi-cadence, et les modales ♭VIImaj7-I et ♭VIIm7-I qui, sans leur
  mode, ne se distinguent pas des mouvements conjoints de n'importe
  quelle grille tonale (Em7 Fmaj7 en do, Dm7 Em7).
- **La cadence rompue** se note « … » : sur un bloc qui a son II et
  dont le V va ailleurs que sur sa cible, sur une cellule dont le V
  évite le I promis, et sur le V-VI relu à la seconde écoute. La
  distinction du livre entre rompue et évitée (selon qu'il y a
  modulation) n'est pas faite : la ligne des toniques montre déjà si la
  nouvelle tonique s'installe.
- **La lecture en crochets** comme dans le livre, et la lecture sur la
  tonalité du passage, toutes les deux rendues (voir « L'accord »).
- **Le dorien reste plagal** quand le mineur est la tonique : Fm7 B♭7
  dans *Mas Que Nada* est le I IV7 de fa mineur, et un II-V rejoué plus
  de quatre mesures est un vamp, pas une cadence.

Repères : `Approaches` et `harmony.ApproachKind`, `Blocks` et
`Block.Announced`, `PassingChords`, `Reread`.

## 4. Entre les cadences : liens, cellules, marches

Les cadences ne se suivent pas au hasard. Le livre nomme ce qui les
relie, et gohar le lit en plus des blocs.

### Les II-V consécutifs

Des II-V qui s'enchaînent sans se résoudre sur la tonique annoncée :
enchaîner la sous-dominante et la dominante d'une tonalité sans en
donner la résolution crée une rupture, et la conclusion en est d'autant
plus attendue, un « rebond » d'une cadence à l'autre, cher au be-bop
(*En Harmonie*, tome 1, chapitre 8 §3.3). Le livre nomme le lien par le
pas qu'on trouve d'une mesure à l'autre entre un accord de chacun
(`Links`) :

| Pas | Exemple | Thème |
|---|---|---|
| ½ ton | Dm7 G7 \| E♭m7 A♭7, Dm7 G7 \| D♭m7 G♭7 | *Butch And Butch* : C♯m7 F♯7 Cm7 F7 Bm7 E7 |
| ton | Dm7 G7 \| Em7 A7, et du V au II suivant : Dm7 G7 \| Am7 D7, Dm7 G7 \| Fm7 B♭7 | *Satin Doll* |
| cycle des quintes | Dm7 G7 \| Cm7 F7 : Cm7 est le I que la première cadence aurait pu conclure, et le II de la suivante | *Confirmation* : Em7♭5 A7 Dm7 G7 Cm7 F7 |

Les II-V se lisent sur les préparations (II→ puis son V), pas sur les
blocs : dans *Satin Doll*, Dm7 G7 | Dm7 G7 | Em7 A7 monte d'un ton,
quoi qu'on entende du G7 qui revient sur Dm7. Un II-V rejoué ne fait
pas de lien mais ne rompt pas la chaîne, et un pas que le livre ne
nomme pas n'en fait pas non plus. `analyse` écrit le pas devant le II :
« ½ II→ », « step II→ » (un ton), « 5th II→ » (le cycle des quintes).

Le cycle des quintes l'emporte quand les fondamentales continuent de
tomber de quinte en quinte à travers les deux II-V : E A D G C F dans
*Confirmation*. Quand le second II-V rompt la chaîne, le lien se lit
entre les II : dans *Autumn Leaves*, D7 tombe d'une quinte sur Gm7,
mais G♭7 rompt le cycle, et Am7♭5 D7 | Gm7 G♭7 descend d'un ton. Règle
à nous, tirée de l'écoute.

Le II peut être celui d'une dominante chromatique, un demi-ton au-dessus
d'elle : dans *Autumn Leaves*, Gm7 G♭7 | Fm7 E7 descend d'un ton. Le
second maillon est détourné, E7 allant sur Am7♭5 dont il est le V, et
Fm7 est alors, par rapport à la, un ♭VIm7 plutôt qu'un II. L'analyse
garde « II », la lecture de la chaîne que l'oreille suit : c'est le
rebond d'un II-V qui ne conclut pas.

Les dominantes s'enchaînent de même sans leurs II, chacune V de la
suivante : D7 G7 C7 F7, le pont des rhythm changes (*Anthropology*),
par le cycle des quintes, ou E7 E♭7 D7 D♭7 par demi-tons, chacune
dominante chromatique de la suivante. `analyse` écrit « 5th V→ » ou
« ½ V→ » devant la dominante qui continue la chaîne. Un V qui va sur un
II (G7 Cm7) ne continue pas une chaîne de dominantes : c'est le cycle
des II-V.

### Les cellules

Des formules de quelques accords que les standards reprennent, et
qu'un musicien entend d'un bloc (`Cells`). *En Harmonie* les présente
avec les enchaînements fréquents (tome 1, chapitre 8 §3.3) :

- **l'anatole**, I VI II V, « connu en France sous le nom d'"anatole" »
  (la cellule, à ne pas confondre avec la forme anatole) « et dans les
  pays anglo-saxons "rhythm changes" » : B♭ Gm7 Cm7 F7 dans *I Got
  Rhythm*. Il forme « un enchaînement cyclique suivant le cycle des
  quintes », qu'on joue là où la durée n'est pas définie, une
  introduction ou une coda. « Rencontrée en majeur et en mineur, cette
  progression harmonique existe sous de nombreuses formes grâce à
  divers emprunts et substitutions » : Fm D♭maj7 Gm7♭5 C7 en fa mineur
  harmonique, Cm Am7♭5 Dm7♭5 G7 dans *Softly, As In A Morning
  Sunrise*, et avec des dominantes secondaires, Cmaj7 A7 Dm7 G7 ou
  Cmaj7 A7 D7 G7 (chapitre 9) ;
- **le III-VI-II-V-I**, « simple variante de l'anatole, fréquemment
  rencontrée en début ou fin de morceau lorsque l'on veut jouer deux
  fois de suite l'anatole sans pour autant rejouer le degré I. Le IIIe
  degré est substitué au Ier » : Fmaj7 Dm7 Gm7 C7 Am7 Dm7 Gm7 C7 dans
  *Have You Met Miss Jones*.

Le tome 2 en ajoute une, parmi les cadences modales (chapitre 2 §5.2) :

- **la cadence modale de do éolien**, ♭VImaj7 ♭VII7 I, qui monte vers
  le I par tons : A♭maj7 B♭7 Cm7 dans son propre mode, A♭maj7 B♭7
  Cmaj7 quand elle résout en majeur. Ce second cas est un emprunt, qui
  remplace le II-V d'un II-V-I « afin de dynamiser l'enchaînement par
  un nouveau mouvement de basses et des couleurs étrangères à la
  tonalité ». Les musiques de jeux vidéo en ont fait un classique,
  *Final Fantasy* plus encore que *Super Mario*, dont la fanfare de fin
  de niveau lui a pourtant donné son surnom : la « cadence Mario ».
  Sans son ♭VImaj7, le ♭VII7-I reste une plagale mineure déguisée, un
  bloc : c'est le ♭VImaj7 qui fait la cadence modale. Dans le corpus,
  elle résout le plus souvent en mineur : Amaj7 B7 C♯m7 dans le *Guile's
  Theme* de *Street Fighter II*, G A7 Bm dans l'*Aria di mezzo
  carattere* de *Final Fantasy VI*. La « cadence Mario » y est plus
  rare : F G7 A dans la *Route 209* de *Pokémon*.

Les règles de lecture sont les nôtres. Une cellule se lit sur les
fondamentales, depuis la tonique que pointe le V : I (ou III), VI (sur
la sixte majeure, ou mineure en mineur), II, V. Le premier accord est
un accord de tonique pour l'anatole, un accord mineur pour le
III-VI-II-V ; le VI a une tierce, de n'importe quelle qualité ; le II
est mineur ou de dominante, un accord majeur y étant l'arrivée d'une
cadence (Em7 A7 Dmaj7 G7 est un II-V-I de ré, pas un III-VI-II-V de
do) ; le dernier est une dominante. Un accord tenu plus longtemps
compte une fois. La cadence éolienne se lit de même, sur trois
fondamentales qui montent par tons : un accord majeur, une dominante,
un accord de tonique, mineur (la cadence dans son propre mode) ou
majeur (l'emprunt, la « cadence Mario »). Elle peut partager son I
avec l'anatole qui en part.

Les X7 d'une cellule peuvent être des substitutions tritoniques :
« tout accord X7 peut être substitué à un autre accord X7, quelle que
soit sa fonction » (chapitre 9, p. 141). Le livre analyse un tel
passage « de la droite vers la gauche », pour « retrouver les cadences
initiales » (p. 142-143, *Sophisticated Lady*, *I Thought About You*).
`Cells` fait de même : elle lit les quatre accords tels qu'écrits,
puis avec un X7 remplacé par son jumeau, puis deux, puis trois, et
garde la première lecture qui fait une cellule. Un VI restitué doit
être le V du II. C E♭7 A♭7 D♭7 est l'anatole C A7 D7 G7, et Em7 E♭7
Dm7 G7 (*Blue In Green*, *Too Young*) un III-VI-II-V ; C D7 Dm7 G7
(*Take The A Train*) n'en est pas un, D7 y étant le II7 et non le
jumeau d'un VI. `analyse` écrit « anatole ♭II », nom retenu en
attendant mieux.

Pas de cellule, en revanche, quand son V mène à une dominante au
milieu d'une mesure : la formule se dissout dans une chaîne de
dominantes, où « chaque accord est emprunté à une tonalité différente »
(Siron, p. 357). Dans *Body And Soul*, Dm7 G7 | C7 B7 B♭7 descend de
dominante en dominante jusqu'à E♭m, et dans *'Round Midnight*, Fm7♭5
B♭7 | E♭m7 A♭7 D♭7 continue sur G♭7. Une dominante sur le premier
temps d'une mesure est une arrivée : le I7 d'un blues. Que le milieu de
la mesure les distingue est une règle à nous, tirée de l'écoute ; elle
écarte 21 cellules du corpus.

Une cellule demande aussi un rythme harmonique constant : ses quatre
accords durent autant les uns que les autres. Dans *Autumn Leaves*,
Am7♭5 | D7 | Gm7 G♭7 est un II-V-I suivi d'un I qui s'en va ; le
III-VI-II-V est le suivant, Gm7 G♭7 | Fm7 E7. Règle à nous, elle aussi
tirée de l'écoute : sans rythme constant, ce n'est pas une formule mais
une simple descente de quintes, et une cellule ne ralentit pas. Elle
écarte près d'un quart des cellules du corpus, le plus souvent un II-V
par mesure suivi d'un II et d'un V d'une mesure chacun (Em7 A7 | Dm7 |
G7).

Le livre nomme la variante avec son I, que le V promet. Une cellule
s'entend pourtant à sa forme, que le V tienne sa promesse ou non : dans
*Anthropology*, Fm7 B♭7 E♭7 A♭7 est un III-VI-II-V de ré♭, et A♭7 va
sur Dm7. C'est une cadence rompue, que le livre note « V – … » :
« enchaîner l'accord de dominante avec un autre accord que l'accord de
tonique prioritairement attendu » (tome 1, chapitre 8). Le III en est
une aussi : la médiante tient lieu du I dans la cellule qui suit, mais
le V n'y résout pas. Dans *Have You Met Miss Jones*, le C7 du premier
anatole va sur Am7 ; dans *Anthropology*, F7 va sur Dm7. La cellule le
note (`Resolves`), et `analyse` écrit les cellules sous les toniques,
« anatole ───── », avec des points de suspension quand le V évite le I
promis : « III-VI-II-V… ───── ». La cadence éolienne s'y écrit
« ♭VI-♭VII-I ───── ».

### Les accords parallèles

Deux notions se croisent ici. **Des accords parallèles** sont au
moins trois accords de même tétrade, séparés par un intervalle qui
peut varier. **Une marche harmonique** répète un motif, un accord ou un
II-V, à intervalle constant : les II-V consécutifs en sont. Les deux
ensemble font **une marche d'accords parallèles**, dont les X7 qui
descendent par quintes sont le cas le plus trivial, puisqu'ils suivent
la résolution qu'appelle le triton.

Les accords parallèles se lisent sur la grille seule, sous une mélodie
qui ne les suit pas. Les mouvements conjoints sont les plus courants ;
un accord mineur plaqué sur chaque note d'une même tétrade en est un
autre cas. Le pont de *Stolen Moments* (Dm D♯m | Em Fm | F♯m Fm |
Em E♭m) fait monter et descendre par demi-tons des accords doriens
sous un ostinato sur un intervalle de tierce : c'est une couleur,
chaque accord vient de son propre mode, et le fond ne bouge pas.
**Un accord qui harmonise une note de la mélodie**, bref et de même
qualité une seconde à côté de l'arrivée, est un autre cas, qui demande
la mélodie pour être sûr : sans elle, gohar ne peut que le conjecturer.

Des accords parallèles peuvent aussi se préparer l'un l'autre : les X7
de *Sophisticated Lady*, les B♭m7 Bm7 B♭m7 Bm7 d'*Along Came Betty*,
où chaque Bm7 est le II de E7. Les blocs et les chaînes de dominantes
les lisent déjà ; reste à décider ce que l'analyse montre de la
marche par-dessus. Pas encore codé.

Repères : `Links`, `Cells`.

## 5. La tonalité


> « The most difficult thing in jazz is to find the one. »
> (Hal Galper)

Cette blague de musicien dit tout : trouver la tonique, comme trouver
le premier temps de la mesure, est tout aussi difficile que trouver
l'âme sœur. Pour l'analyse, c'est le travail le plus laborieux de tous,
et celui où l'oreille garde le dernier mot. Cette partie raconte comment
on s'y prend : ce qu'on entend en avançant dans la grille, accord par
accord, puis ce qu'on conclut une fois arrivé au bout.

### Dans le livre

Le livre ne cherche pas la tonalité dans l'armure, il la lit dans les
accords. Sa méthode va de droite à gauche : « il faut dissocier les
accords réels (appartenant à la tonalité), les accords modifiés et les
accords ajoutés », prendre l'accord d'arrivée, puis remonter les
accords qui le préparent. L'armure, il la recommande à l'apprenant qui
« se creuse la tête » sur une partition, sans plus. Et l'un des dix
commandements de l'harmoniste de la BEPA, dus à Étienne Guéreau, dit
le reste : « Tu ne suivras pas bêtement les indications du Real
Book. »

Entre les accords de la tonalité, il distingue ce qui la quitte de ce
qui n'en sort pas. **L'emprunt** fait venir un accord ou une cadence
d'ailleurs sans changer de tonalité : Dm7 G7(♭13) Cmaj7 emprunte son
G7(♭13) à do mineur. Une cadence vers un autre degré en fait un « Ier
degré temporaire » : dans *There Will Never Be Another You*, Dm7♭5 G7♭9
Cm7 tonicise le VI de mi♭ sans quitter mi♭. **La modulation** change de
tonalité : « une cadence prépare généralement la modulation, celle-ci
étant confirmée si la durée est assez longue pour que la nouvelle
tonalité soit installée » (tome 1, chapitre 10 §1.8, p. 159). Le livre
ne chiffre pas cette durée ; ses fiches la montrent : quatre mesures de
do dans *Tune Up*, sept dans *Black Orpheus*. Siron affine en trois
niveaux, tonicisation, modulation transitoire, modulation vraie (voir
« La modulation »).

Enfin le livre met en garde, dans le tome 2, contre l'analyse modale
appliquée à tort : Dm7 G7 Cmaj7 n'est pas ré dorien, sol mixolydien et
do ionien, mais une cadence parfaite en do.

### Dans gohar, en avançant : la tonique pressentie

L'oreille qui avance dans une grille n'attend pas la fin pour savoir où
elle est. À chaque accord, elle a une tonique en tête, avec ce qui a
sonné et rien d'autre. C'est la **tonique pressentie**, une notion à
nous, le pendant au présent de la région que la fiche montre après
coup. Au troisième accord de Tenderly, deux mesures de E♭maj7 A♭7 ont
installé mi♭ : E♭m7 s'entend comme la tonique qui change de couleur,
pas comme le II de ré♭.

#### Trois toniques, et une mémoire

**La tonique de fond** est celle qui est installée : les degrés se
comptent sur elle, même quand une cadence tonicise un autre degré (Dm7♭5
G7 Cm7 en mi♭ se lit II V VI). **La tonicisation** fait de l'accord où
une dominante se résout un « degré I temporaire » (Siron, *La partition
intérieure*, p. 341), le temps de cet accord. **La région** est une
tonique secondaire qu'une cadence a installée pour plus longtemps (voir
« Ce qui l'installe ») ; elle dure tant que les accords suivants
tiennent en elle, ou sont sur sa tonique avec sa tierce (Gm7 puis
Gm(maj7) sur sol mineur, une septième qui se promène), ou préparent un
accord qui y tient ; et elle ne change pas le fond. Chacune est un
ensemble de tonalités sur une même tonique (fa mineur harmonique ou
mélodique, quand rien ne tranche), et le fond est vide au début d'un
morceau, avant la première tonique.

À part, **la tonique de départ** est gardée en mémoire, pour reconnaître
le retour au fond de départ après un pont qui a modulé, si fréquent sur
une forme AABA.

#### Ce qui fait une tonique

*En Harmonie* le dit en une phrase : l'accord de tonique « peut tout
aussi bien être Xmaj7 que Xm7 », mais « ne comportera jamais le IVe
degré en tant que note dans l'accord » (tome 1, chapitre 8 §2.1,
p. 101). C'est sa tierce qui dit son mode, majeur ou mineur, et ses
extensions n'y changent rien : un Cmaj7(♯11) ou un Fm6/9(♯11) sont des
toniques, leurs extensions se confondant avec des appoggiatures des
notes réelles (tome 2, chapitre 6 §3). Il peut être renversé.

Pouvoir être une tonique n'est pas l'être, et gohar distingue deux cas.
Une triade, un maj7, un 6, un m6 ou un m(maj7), à l'état fondamental ou
avec sa tierce à la basse, sont des toniques partout. Un m7, ou une
tonique renversée sur sa quinte, ne le sont que là où une cadence se
résout sur eux, parce que hors de ce contexte ils sont d'abord autre
chose :

- **Le m7** est le plus souvent un II. Les grilles écrivent pourtant la
  tonique mineure m7 bien plus souvent qu'on ne la joue (m6, m(maj7),
  m(maj9), pour adoucir la septième). Un m7 est une tonique quand une
  cadence se résout dessus, plagale comprise, et qu'il n'est pas
  lui-même le II d'un bloc. *Softly, As In A Morning Sunrise* est ainsi
  en do mineur, et la cadence éolienne Amaj7 B7 C♯m7 installe do♯
  mineur dans le *Guile's Theme*. Le C7 Gm7 de *Honeysuckle Rose* n'en
  fait pas une tonique : Gm7 est aussitôt le II de Gm7 C7 F. Un
  turnaround vers un premier accord en m7 non plus : l'Am7 qui ouvre
  *Fly Me To The Moon* est un VI, tant que E7 n'y revient pas à la
  mesure 8.
- **La tonique renversée sur sa quinte** est une tonique quand un V se
  résout sur elle : c'est la cadence imparfaite, C7 Fmaj7/C (tome 1,
  p. 103), sa basse étant un « retard de la dominante » (tome 2,
  p. 139). Sans V devant elle, elle se lit comme un accord sur une
  pédale : dans le F/C C qui clôt *My Way*, F/C est le IV de do
  au-dessus d'une pédale de tonique.

Avec toute autre basse, c'est la basse qui tient : le E♭maj7/F de *The
Look Of Love* est une pédale de fa.

#### Ce qui l'installe

Une **cadence qui se résout** tonicise sa cible. Elle n'ouvre une
région, une modulation au sens d'*En Harmonie*, que si elle a sa
sous-dominante et que son I termine un groupe de deux ou quatre mesures
de la section, ou en ouvre une. La sous-dominante, c'est celle d'un
II-V-I ou d'un Vsus V I ; une plagale est une sous-dominante suivie de
sa tonique, et toute cadence définit la tonalité ; un V seul est la
dominante secondaire de Siron, qui tonicise (p. 341).

La place est un réglage à nous, calibré sur les fiches du livre. Siron
ne donne ici qu'un exemple, la halte d'une modulation transitoire
« à l'aide d'une cadence complète à la fin d'un groupe de deux ou quatre
mesures » (p. 379). gohar retient la seconde mesure de chaque paire et
la troisième de chaque groupe de quatre, là où une phrase de huit
mesures conclut en cadence forte, sur sa septième (p. 390), et là où
arrivent les modulations de *Tune Up* et de *Black Orpheus* ; pas la
première mesure d'un groupe de quatre à l'intérieur d'une section, où
tombent les crochets de *There Will Never Be Another You* (mesures 21 et
25). La première mesure d'une section ouvre en revanche une région,
d'après le critère de place que Siron donne pour la modulation vraie :
elle est « plus volontiers située au début ou à la fin d'un cycle de
mesures ou d'une phrase mélodique » (p. 380). Le E♭maj7 qui ouvre le
pont de *My Funny Valentine* ouvre ainsi mi♭ majeur, dans le mode de
l'accord d'arrivée même quand le II-V était mineur (Fm7♭5 B♭7). Et le I d'une cadence qui
repart aussitôt comme le II de la suivante n'ouvre rien : Cm7♭5 F7♭9
B♭m7 E♭7 A♭maj7, au pont de *Lullaby Of Birdland*, est un III-VI-II-V
de la♭, qui tonicise le cycle des quintes. Sans mesures, en direct,
toute cadence qui a sa sous-dominante ouvre une région.

Seul un accord qui a une triade majeure ou mineure peut être tonicisé.
Hors d'un contexte modal, un m7♭5 est un II mineur : c'est ainsi que
s'enseigne le II-V-I, m7♭5 sur le II en mineur là où le majeur joue m7,
et Siron ne donne pour tonalités d'emprunt d'une dominante secondaire
que celles des degrés II à VI, jamais le VII, le m7♭5 du majeur
(p. 342). Le Fm7♭5 de la mesure 9 de *Tenderly*, où mène Gm7♭5 C7♭9,
est le II de mi♭ mineur.

Au début, le **premier accord** installe le fond s'il peut être une
tonique, sinon la première cadence résolue. C'est un indice faible,
beaucoup de standards commençant sur un II ou un IV : ce fond est
l'hypothèse de l'oreille qui avance, et la tonalité du morceau se
conclut au bout (voir « La tonalité, par le premier et le dernier
accord »). Un **blues** reconnu a sa tonique pour fond dès la première
mesure.

Une **plagale fait tout ce que fait un V-I** : IV puis I, ou ♭VII7
puis I, elle confirme une tonique déjà là, ramène au fond de départ, ou
ouvre une région. *En Harmonie* ne réserve ce rôle à aucune cadence :
une cadence est « un enchaînement d'accords caractéristiques tendant à
marquer le passage d'une phrase à l'autre, tout en définissant la
tonalité » (tome 1, chapitre 8 §3, p. 103). Ce qui sépare un Dm7 Am7 de
passage d'un IV-I qui installe la mineur, c'est la durée (voir « La
modulation »). Ce qu'elle annonce, l'oreille
l'attend : après le D♭7 de Tenderly, on attend mi♭.

La **cadence éolienne**, ♭VImaj7 ♭VII7 I, est une cadence modale, qui
« doit obligatoirement faire entendre les DCN et DCA du mode » (tome 2).
Dans le *Guile's Theme* de *Street Fighter II*, Amaj7 B7 C♯m7 revient
quatre fois, suivi chaque fois d'une mesure de repos sur C♯m7 : le
morceau est en do♯ mineur, et le seul II-V-I, en mi, tonicise le
relatif majeur dans la troisième section. Elle n'est jamais relue comme
un V qui va sur son VI (voir « Une seconde écoute », plus bas).

Ne change rien au fond : un accord diatonique, un emprunt sur la même
tonique (Im7, IVm, ♭VII7), une préparation qui ne se résout pas. Un
accord dont la fondamentale est la tonique du fond se lit comme un **I
emprunté**, même quand il est le II d'un II-V qui ne se résout pas :
E♭m7 A♭7 en mi♭ est Im7 IV7, et le bloc annonce toujours ré♭. Le livre
écrit ces mesures « I IV », la qualité empruntée notée une fois et plus
ensuite ; l'analyse l'écrit chaque fois.

#### La modulation

*En Harmonie* la définit en deux conditions : « Une cadence prépare
généralement la modulation, celle-ci étant confirmée si la durée est
assez longue pour que la nouvelle tonalité soit installée » (tome 1,
chapitre 10 §1.8, p. 159). gohar la lit au présent, à deux niveaux. La
**région** est libérale : une cadence à la bonne place l'ouvre, et
appeler modulation une tonicisation appuyée est une analyse que
beaucoup de musiciens feraient. Le **fond**, lui, est strict : il ne
bascule que sur une modulation vraie.

Une région devient le fond quand c'est une **modulation vraie** : entendue
depuis la cadence qui l'a ouverte, elle tient la première mesure d'une
section et au moins la moitié de celle-ci (voir « Transitoire ou
vraie »). Le pont de *Body and Soul* installe ré ; le do d'*All The
Things You Are* reste une région, et les degrés restent comptés en
la♭ : c'est la double analyse de Siron, les degrés dans la tonalité et
les fonctions dans la région (p. 379). La ligne des crochets, elle,
lit une région confirmée dans sa tonalité, comme le livre imprime ses
modulations : le B♭maj7 Gm7 de *Tune Up* y est I VI. Une région dure
tant qu'aucun accord du fond ne revient : un accord qui n'appartient ni
au fond ni à la région est une tonicisation à l'intérieur de celle-ci,
le Gm7 C7 du pont de *Body and Soul*, en ré.

Sans sections, rien ne dit où tombe une modulation dans la forme, et
une région reste une région : le fond ne bascule pas. Une règle de
repli a longtemps tenu là, une mesure d'accords stables ou une deuxième
cadence ; elle ne servait plus qu'aux suites d'accords sans mesures des
tests, et elle est tombée.

Une cadence qui n'ouvre pas de région ne fait que toniciser, et laisse
la région en place (A7 Dm7 dans la région de do, dans *Black Orpheus*).
Les cadences vers les tons voisins et leurs relatifs « sont fréquemment
utilisées comme modulations transitoires » (Siron, *La partition
intérieure*, p. 385) : une région secondaire qui ne met pas la tonalité
du morceau en danger.

**Trois niveaux.** Siron en distingue trois (*La partition intérieure*,
p. 341, 342 et 378 à 380). La **tonicisation** « ne porte que sur un
accord » : une dominante fait de l'accord où elle se résout un « degré
I temporaire », et les degrés restent comptés dans la tonalité (le Dm7
qui suit A7 en do). La **modulation transitoire** « introduit de
manière plus prolongée les altérations d'une autre tonalité », sans que
la tonalité du morceau soit en danger : « l'oreille ne cesse de garder
un œil sur son parfum ». Elle sert de halte, comme le do d'*All The
Things You Are*, ou de tension, quand s'enchaînent des cadences
incomplètes : dans Em7 A7 | E♭m7 A♭7 | Dm7 G7 | CΔ, chaque II-V est une
modulation transitoire. La **modulation vraie** : un nouveau centre,
installé.

*En Harmonie* découpe autrement : sa modulation (une cadence, et une
durée « assez longue ») est plus large que la tonicisation de Siron et
plus étroite que sa transitoire. Il met entre crochets un II-V-I tenu
deux mesures, là où Siron appellerait transitoire un II-V qui ne se
résout pas. gohar suit le livre, dont les fiches sont l'oracle : la
tonicisation de chaque accord (`Tonicised`), la région (`Region`),
qui est une modulation au sens d'*En
Harmonie*, et le fond (`Ground`). `analyse` écrit la région entre
parenthèses et la tonicisation entre crochets : « (C) », « [Dm] ». Le
passage de la région au fond est la modulation vraie de Siron (voir
plus haut, et « Transitoire ou vraie »).

**Ce qui fait une modulation vraie.** Siron la juge à quatre critères
(p. 380) : la durée, « importante pour distinguer la sensation de
modulation vraie d'une modulation transitoire » ; la mémoire auditive,
« la première tonalité entendue a toujours un énorme poids » ; la forme,
une modulation vraie étant « plus volontiers située au début ou à la fin
d'un cycle de mesures ou d'une phrase mélodique. Au milieu d'une phrase
harmonique, l'oreille entend plutôt une modulation transitoire » ; la
proximité, « la sensation de modulation est plus forte si les tonalités
sont éloignées », alors que « l'enchaînement rapproché de centres tonaux
éloignés détruit la sensation d'une véritable modulation ». La distance
se compte en crans sur le cycle des quintes, un relatif partageant
l'armure de sa tonalité (p. 383). gohar les
mesure pour chaque **zone tonale**, un passage entendu autour d'une
autre tonique que la première, région ou fond installé, compté depuis
la cadence qui y mène (`TonalAreas`). Aucune source ne les pondère ;
gohar en fait peser deux, la place avec la durée (« Transitoire ou
vraie ») et la distance (« Les centres éloignés »), la mémoire pas
encore. `analyse` les affiche toutes. *Tune Up* : do mesures 5 à 8,
quatre mesures, quittant ré, la première tonalité, à deux crans ; si♭
mesures 9 à 12, quatre mesures, à deux crans de do ; transitoires
toutes deux.

Une zone ne compte que si elle dure : « une cadence prépare
généralement la modulation, celle-ci étant confirmée si la durée est
assez longue pour que la nouvelle tonalité soit installée » (*En
Harmonie*, tome 1, chapitre 10 §1.8, p. 159). Le livre ne donne pas de
longueur ; gohar retient deux mesures au moins, la cadence comprise.
Les fiches du livre sont l'oracle de ce seuil : ses modulations (do et
si♭ dans *Tune Up*, do dans *Black Orpheus*) doivent être des zones,
sur les mêmes mesures, et ses crochets hors modulation (*Tenderly*,
*There Will Never Be Another You*) n'en ouvrir aucune. Gm7 C7 | Fm7,
aux mesures 30 et 31 de *There Will Never Be Another You*, le livre le
met entre crochets : une tonicisation, qui dure moins de deux mesures.
Le pont d'*All The Things You Are*, F♯m7 | B7 | Emaj7, que Siron range
parmi les modulations (p. 387), en dure trois. Une zone est donc une
modulation au sens d'*En Harmonie* ; la transitoire de Siron, qui
commence dès un II-V, en compte davantage.

Trois exemples de Siron ont rejoint les fiches, pour servir d'oracle à
la modulation vraie : la « respiration secondaire » en do d'*All The
Things You Are* (5.11.1), une modulation transitoire, que l'analyse
entend comme une zone des mesures 6 à 8 ; les cadences incomplètes
enchaînées Em7 A7 | E♭m7 A♭7 | Dm7 G7 | CΔ (5.11.2), des modulations
transitoires de tension, plus fines que les zones, que le test note
sans les exiger ; la section A d'*In a Sentimental Mood* (5.11.5),
ambiguë entre ré mineur et fa majeur, où l'analyse entend fa sur les
deux dernières mesures.

D'autres les ont suivis. Le premier est *Pithecanthropus Erectus*
(5.11.6), mesures 9 à 15 : une modulation dans la région de la
sous-dominante mineure, de fa mineur vers sol♭, par des accords-pivots
que Siron lit dans les deux tonalités. Les deux derniers sont des
**régions d'ambiguïté tonale** (5.11.7). La première est une
cellule-anatole très chromatique, CΔ A7 | A♭Δ♯5 G7♭5 | F♯7 FΔ♯5 |
A♭7♭5 D♭7♭5, « région tonalement floue due aux nombreuses altérations
jamais résolues ». La seconde est le début de *Grand Central*, trois
II-V qui ne résolvent pas, vers la♭, sol♭ et mi, entre deux Fm. Un
troisième manquait : le pont de *Jordu*, une cascade d'accords 7 sur le
cycle des quintes, de G7 à D♭7 puis de F7 à G7 (5.8.37), où « chaque
degré tend à devenir interchangeable avec son voisin » (p. 357). Dans
ces régions, dit Siron, « il devient alors difficile de parler de
véritables modulations » (p. 381) : le test y accepte une zone, jamais
une modulation vraie.

Deux cadences coltraniennes complètent le tableau (p. 533). *Countdown*
remplace le Dm7 G7 | CΔ de *Tune Up* par Dm7 E♭7 | A♭Δ B7 | EΔ G7 | CΔ,
« un carrousel de tonalités » entre do, la♭ et mi, trois tonalités
qui n'ont que trois notes communes. *Giant Steps* n'est « composé que de
cadences dans 3 tonalités », si, sol et mi♭, à une tierce majeure les
unes des autres ; l'avant-dernière cadence, en mi♭, lui donne « plus de
poids ». Siron n'y chiffre pas de degrés : il met chaque cadence entre
crochets avec sa tonalité. Aucune de ces cadences ne dure assez pour
installer sa tonalité ; une zone qui en enjambe plusieurs contredit le
livre.

**Transitoire ou vraie.** gohar tient une zone pour une modulation
vraie quand elle tient la première mesure d'une section et au moins la
moitié de celle-ci. Siron situe la modulation vraie « plus volontiers
[…] au début ou à la fin d'un cycle de mesures ou d'une phrase
mélodique », la transitoire « au milieu d'une phrase harmonique », et
la durée y pèse (p. 380) ; le seuil de la moitié est à nous. La
mémoire, son troisième critère, ne pèse encore rien ; la distance, le
quatrième, empêche pour l'instant une zone de se former (voir plus
bas), sans rendre une modulation plus vraie. Les témoins : le pont de *Body and Soul*, de ré♭ à ré, et le
deuxième A de *Joy Spring*, un demi-ton plus haut, les deux modulations
abruptes de Siron (p. 382), sont vraies ; le pont d'*In a Sentimental
Mood*, en ré♭ dès la mesure 17 (*En Harmonie*, p. 159), aussi ; le do
d'*All The Things You Are*, sa « respiration secondaire », est
transitoire. Les grilles de l'app mènent d'ailleurs à ces ponts par un
II-V (Em7 A7 avant le Dmaj7 de *Body and Soul*, A♭m7 D♭7 avant le
G♭maj7 de *Joy Spring*) : la modulation abrupte de Siron se lit ici
comme une modulation par cadence, et gohar n'a pas besoin d'un autre
mécanisme pour l'entendre. Le pont de *Lullaby Of Birdland* et celui de
*My Funny Valentine* sortent aussi en modulations vraies, ce que
l'écoute confirme. `analyse` l'écrit au bout de chaque zone :
« true modulation » ou « transitory ».

**Les centres éloignés.** Les tonalités qu'une dominante emprunte sont,
« dans une harmonie peu chromatique », les voisines de la tonalité
(Siron, p. 342) : les degrés II à VI. Un II-V qui résout sur l'accord de
tonique d'une tonalité éloignée, à deux crans ou plus sur le cycle des
quintes, fait entendre un autre centre. Ce centre coupe la zone où il
sonne, sans en former une : une tonicisation « ne porte que sur un
accord » (p. 379). Et « l'enchaînement rapproché de centres tonaux
éloignés détruit la sensation d'une véritable modulation » (p. 380) :
un passage de deux mesures au plus, entre deux centres éloignés aussi
courts que lui, n'installe rien. *Giant Steps* n'est « composé que de
cadences dans 3 tonalités », si, sol et mi♭, à une tierce majeure les
unes des autres (p. 533) : gohar n'y entend plus aucune zone, là où il
entendait mi♭ de la mesure 2 à la mesure 9, par-dessus les cadences en
sol et en si. Les deux crans et les deux mesures sont à nous. Le pont
de *Grand Central*, quatre mesures de F♯m7 B7, reste une modulation
vraie : l'écoute y entend bien le centre bouger.

Le **retour au fond de départ** est asymétrique : une seule cadence sur
la tonique de départ le réinstalle, parfaite ou plagale, et même son
accord de tonique seul (le F/C C de *My Way*). Quitter demande plus
d'indices que revenir.

En direct, le fond bascule au moment où la modulation vraie est
acquise ; après coup, la zone commence au bloc qui y menait, et la
ligne des crochets la lit dans sa tonalité (Tune Up, mesure 7 : Cmaj7
est I, pas ♭VIImaj7).

#### Local, avec une mémoire

La contrainte du direct tient : chaque calcul ne regarde qu'un nombre
borné d'accords, et la tonique pressentie est un état porté d'un accord
au suivant, un résumé de ce qui a sonné. Seule la lecture d'une grille
entière va plus loin, parce qu'elle le peut : une grille qui boucle est
entendue comme son deuxième chorus, qui part de la fin du premier, avec
pour chez-soi la tonalité du morceau, désormais connue. Un turnaround
en fin de grille prépare donc le premier accord, et Tune Up commence en
ré. Le direct n'a que la première écoute, où une cadence à travers la
boucle n'a pas encore sonné.

#### Tenderly, au présent

| Mesure | Accord | Tonique pressentie | Ce que l'oreille entend |
|---|---|---|---|
| 1 | E♭maj7 | mi♭ | I |
| 2 | A♭7 | mi♭ | IV7, une plagale qui annonce mi♭ mineur mélodique |
| 3 | E♭m7 | mi♭ | la tonique change de couleur : I emprunté |
| 4 | A♭7 | mi♭ | IV7 encore, en parallèle avec la mesure 2 |
| 5 | Fm7 | mi♭ | II |
| 6 | D♭7 | mi♭ | ♭VII7, plagale mineure : on attend mi♭ |
| 7 | E♭maj7 | mi♭ | I |
| 8 | Gm7♭5 C7♭9 | mi♭ | II V de fa mineur : on attend Fm |
| 9 | Fm7♭5 | mi♭ | pas Fm, mais le II de mi♭ mineur |

Ce tableau sert de test, validé en attendant l'avis d'une oreille plus
experte.

### Dans gohar, au bout : la tonalité du morceau

Une fois la grille lue, il reste à conclure, comme le livre : par le
premier accord et par le dernier. Le dernier se lit sur les
**phrases**, là où le morceau s'arrête ; le premier sur la première
cadence, qui dit s'il est la tonique là où il sonne.

#### Les phrases

Une phrase va d'un repos au suivant. Elle se termine quand elle se
pose, ou quand le morceau s'arrête.

« La cadence conclusive est le point d'arrivée d'une phrase
harmonique », et dans une musique carrée, les cadences conclusives
marquent les groupes de mesures (Siron, *La partition intérieure*,
p. 390). Quand la grille a des mesures, une phrase ne se pose donc que
là où une section de la forme conclut (voir « La structure ») :

- **une tonique atteinte en milieu de section** confirme la tonique en
  cours, et ne termine aucune phrase : le E♭maj7 de la mesure 3 de
  *Let's Cool One* ;
- **une section qui finit sur son II-V**, son I tombant sur la première
  mesure de la section suivante, se termine ouverte sur son V : c'est
  la demi-cadence d'*En Harmonie* (tome 1, chapitre 8), le point
  d'interrogation. Le I ouvre la phrase suivante, et à l'écoute d'*A
  Fine Romance*, les paroles le confirment.

Cette règle remplace la nôtre, qui posait une phrase sur toute tonique
amenée par une cadence et qui revenait sur l'ouverture, ou tenait plus
d'une mesure et plus longtemps que sa préparation. Sur le corpus, la
moitié de ses repos tombaient en plein milieu d'une section, et un
quart sur le premier temps de la suivante.

Sans mesures, pour un morceau joué en direct dont la forme n'est pas
encore lue, l'ancienne règle reste en repli : une phrase se pose sur un
accord de tonique qu'une cadence amène, et qui revient sur l'accord
d'ouverture (*How Insensitive*, un long soupir de Dm à Dm, quatorze
mesures plus loin) ou tient plus d'une mesure et plus longtemps que les
accords qui y mènent (le Gm6 d'*Autumn Leaves*). Une tonique de passage
ne pose rien (le B♭maj7 d'*Autumn Leaves*), ni un IV, si long soit-il
(le E♭maj7 de *Cherokee*), ni un m7 après l'ouverture (le Cm7 de *There
Will Never Be Another You*, son VI). Le repos dure autant que la
tonique tient.

#### La maison, une notion tombée

Le morceau partait de là où sa première phrase conclut, la « maison » :
*Autumn Leaves* de sol mineur, *Fly Me To The Moon* de la mineur. La
règle du premier et du dernier accord l'a remplacée (ci-dessous), et
le premier accord ne vaut que si sa première cadence le confirme : le
Am7 de *Fly Me To The Moon*, dont la première cadence va à do, est un
VI. Une règle sœur a tenu ici aussi : un morceau qui « s'ouvre au
repos », sa tonique tenue plus d'une mesure et sa première cadence y
revenant, y restait où qu'il s'arrête. Elle faisait lire *In a
Sentimental Mood* en ré mineur pour une mauvaise raison ; le livre dit
que le thème « est en Ré mineur pour se terminer dans la tonalité de
son relatif Fa majeur » (*En Harmonie*, tome 1, chapitre 10, p. 159),
et c'est la prédominance qui le lit en ré mineur. Il reste de la maison
un repli, quand aucune phrase ne conclut, et l'installation du fond à
la première phrase conclusive dans l'écoute au présent (`Home`), à
retirer quand des témoins diront ce qu'elle protège.

#### La tonalité, par le premier et le dernier accord

*En Harmonie* la lit ainsi (tome 1, chapitre 8 §1.2, p. 99 et 100) :
le premier et le dernier accord, turnaround exclu, la confirment quand
ils désignent la même tonalité (*Blame It On My Youth*, *Angel Eyes*).
Quand ils diffèrent, c'est « la prédominance de l'une ou l'autre des
deux tonalités durant le morceau » qui tranche : *My Funny Valentine*
commence sur Cm, finit sur E♭6, et il est en do mineur. Ni le livre ni
Siron ne disent ce qu'est la prédominance ; gohar prend celle qu'on
entend le plus longtemps : le temps des zones tonales pour la leur, le
reste du morceau pour la tonalité où il commence. Une tonicisation,
plus courte qu'une zone, compte pour la tonique qui l'entoure.

Le premier accord ne vaut que s'il est la tonique là où il est, et le
livre n'en dit pas plus. gohar le tient pour tonique quand la première
cadence du morceau se résout sur lui : *My Funny Valentine* s'ouvre sur
Cm6, et Dm7♭5 G7♭9 revient à Cm7. *Just Friends* s'ouvre sur Cmaj7,
mais sa première cadence va à sol : do est son IV, et seul le dernier
accord parle. De même l'Am de *Blue Skies*, le Fm7 d'*All The Things
You Are*, le Em7 de *Tune Up*.

*In a Sentimental Mood* « est en Ré mineur pour se terminer dans la
tonalité de son relatif Fa majeur » (tome 1, chapitre 10 §1.8,
p. 159) : premier et dernier accord diffèrent, et ré mineur, entendu
le plus longtemps, est la tonalité du morceau. *It Don't Mean A Thing*
s'ouvre et s'entend en sol mineur, et s'arrête sur B♭6 : sol mineur.
*Lullaby Of Birdland* s'ouvre en fa mineur et conclut ses sections sur
A♭maj7 : la règle le lit en fa mineur, là où l'app déclare la♭, et
l'écoute le confirme, ses A sombres en fa mineur, son pont lumineux en
la♭.

#### Le dernier accord, là où s'arrête le morceau

Le dernier accord est celui où le morceau s'arrête, à son dernier
chorus. Quand la grille le dit, on la suit : après la coda, au « Fine », ou sur
l'accord que marque le symbole de fin du lecteur (le `U` des grilles
iReal), tenu sous un point d'orgue. Sans rien de tout cela, l'app
ajoute à la fin un accord de son cru, la tonique de la tonalité qu'elle
déclare : c'est l'armure jouée à l'oreille, et l'analyse ne la lit pas.
La fin marquée l'emporte sur la forme : *Somewhere* conclut sa dernière
section sur A♭, mais la grille l'arrête deux mesures plus tôt, sur E♭.
Elle ne vaut que sur un accord de tonique, ou sur une tonique qu'une
cadence installe : une douzaine de grilles finissent sur un accord 7
(le B♭7 de *Manteca*, le F7♯11 de *Bud Powell*), et un m7 final est
aussi souvent un II laissé en suspens (*Wave*, *Triste*) qu'une tonique.

Une grille muette s'arrête comme finit sa dernière section :

- **sur sa conclusion**, quand elle en a une. Le turnaround qui suit
  ramène au premier accord et ne se joue pas à la fin : *Lullaby Of
  Birdland* s'ouvre sur fa mineur et s'arrête sur A♭maj7, avant que
  Gm7♭5 C7 ne ramène à Fm ; il est en la♭. *All The Things You Are*
  n'est tranché que par son dernier A♭maj7 ;
- **à travers la boucle**, quand elle finit ouverte sur son V : la
  dernière fois, ce V se résout sur le premier accord, et le morceau
  s'arrête là. C'est la demi-cadence d'*A Fine Romance*, prise à la fin
  du morceau : il n'y a plus de phrase suivante, et le I qui l'aurait
  ouverte devient la fin. *Yesterdays* finit sur A7 et s'arrête sur le
  Dm de la première mesure ; *Sugar* finit sur G7 et s'arrête sur Cm7.

Quand la dernière section ne fait ni l'un ni l'autre, et pour un
morceau joué en direct, sans forme lue, l'ancienne règle reste en
repli : le morceau s'arrête sur la dernière tonique entendue. Une
tonique déjà installée, celle de l'ouverture, du repos précédent ou de
la première cadence, revient sans cadence et dans toute position, son
second renversement compris (le E♭6 final de *'Round Midnight*, le D6
de *Chega De Saudade*, le B♭maj7/F d'avant G7 Cm7 F7 dans *Someday My
Prince Will Come*). Une autre demande un II-V-I ou une plagale, et
aucun morceau ne s'arrête sur le IV du repos précédent.

Une section conclut sur une tonique qui sonne dans ses trois dernières
mesures, même atteinte avant : *Sweet Sue* arrive sur G6 à la mesure 5
et le tient jusqu'au bout de la section. Il faut que ce soit l'accord
de tonique qui tienne : dans *Yesterdays*, le Dm Dm(maj7) de la
mesure 5 devient Dm7 à la mesure 6, une ligne qui descend vers Bm7♭5
E7, et la tonique est finie avant les dernières mesures. Un m7 suivi
de son V, ou de son Vsus, reste un II : le Fm7 B♭7sus de la mesure 35
de *Star Eyes* mène au E♭6 final.

Le B♭maj7 de la seconde section de *Yesterdays* tombe dans ses trois
dernières mesures, à la mesure 14 de ses 16. L'oreille n'y entend
pourtant aucune fin de phrase : après sa longue descente de quinte en
quinte, le morceau s'appuie un instant sur ce VI, et la section finit
ouverte sur A7, qui ramène au Dm de la première mesure. Siron en donne
la raison. Une dominante secondaire fait de sa destination un « degré I
temporaire » (p. 341). Dans une harmonie peu chromatique, les tonalités
qu'elle emprunte sont les voisines : le relatif, la dominante, la
sous-dominante et leurs relatifs (p. 342). Et une tonicisation « ne
porte que sur un accord » (p. 379). Ne conclut donc pas un accord amené
par un V seul, sur un temps faible de la section (ailleurs que sur
l'avant-dernière mesure), tenu une mesure au plus, et voisin de la
tonique où s'ouvre la section suivante : un cran au plus sur le cycle
des quintes, relatifs compris. La mesure est un réglage à nous.

Chaque condition compte : *Rosetta* descend le même genre de cycle
jusqu'à F6, mais sur la mesure 15, la forte, et y conclut avant que
Bm7♭5 E7 ne mène au Am de son pont ; *Lover Man* se pose sur Fmaj7 à la
mesure 16, VI du Am de son pont, mais par Gm7 C7, un II-V ; *Spain* se
pose deux mesures sur Bm7, par F♯7 seul, avant de revenir à Gmaj7 :
si mineur est voisin de sol, relatif de sa dominante, mais tenu trop
longtemps pour n'être qu'un appui. La règle ne s'appliquait d'abord qu'au VI d'une
tonique mineure ; étendue aux tonalités voisines, elle rend *Alfie* à
si♭, et le corpus ne perd rien.

Aucune phrase ne se pose sur le IV de la tonique d'ouverture, quand une
cadence revient à celle-ci dans le chorus, et aucun morceau ne s'y
arrête, même quand la grille y marque sa fin. *Unforgettable* passe de
sol à Cmaj7 comme un blues passe du I au IV : c'est sa phrase redite à
la quarte, pas une modulation. Sa grille s'arrête sur le Cmaj7 de la
mesure 31, puis Am7 D7 ramène à sol, où le morceau s'arrête. La
condition sur la cadence compte : *Somewhere* s'ouvre sur B♭, mais ce
B♭ devient B♭7 et aucune cadence n'y revient ; son E♭ n'est pas un IV. Siron
donne au IV ce statut : une « deuxième tonique », une « détente
secondaire qui peut parfois entrer en conflit avec le degré I », et la
cinquième mesure du blues est « une sorte de modulation transitoire à
la sous-dominante » (*La partition intérieure*, p. 384). Une région où
l'on se repose en chemin, pas une tonalité d'arrivée.

*Chega De Saudade* s'ouvre sur ré mineur et s'arrête sur D6, et
l'analyse lit ré majeur, là où il s'arrête.

Quand le premier accord est la tonique là où il sonne, et que la
tonalité du morceau est une autre, `analyse` donne les deux : *I Love
Paris* part de do mineur et s'entend en do majeur (« heard in C,
setting out from Cm »). Le premier accord compte comme pour la tonalité
du morceau (*En Harmonie*, tome 1, chapitre 8 §1.2, p. 99) : la
première cadence du morceau doit y résoudre. *Fly Me To The Moon*
s'ouvre sur Am7, mais sa première cadence va à do : ce Am7 est un VI,
et la ligne n'en dit rien.

#### La tierce picarde

Elle ne rend pas majeur un morceau mineur. Héritée de la musique
d'église, où un accord majeur frotte moins sous la résonance d'un grand
orgue, elle majorise la tonique sur le dernier accord seulement : la
tonique est mineure là où on l'entend d'abord et là où on l'entend en
dernier avant la fin, et majeure sur l'accord final. *'Round Midnight*
clôt son deuxième A sur E♭6, puis son dernier A repart sur E♭m : tierce
picarde, le morceau reste mineur et l'analyse le signale. *Chega De
Saudade* tient ré majeur toute sa seconde moitié : pas une tierce
picarde, un morceau autant majeur que mineur. De même *I Love Paris* et
*Black And Tan Fantasy*, dont la seconde partie est en majeur.

*Once Upon A Summertime* et *Maybe September* tiennent une section
entière en majeur, mais reviennent au mineur avant leur dernier accord :
tierce picarde, et des morceaux mineurs, comme la grille les déclare.
*Somewhere* est un cas ouvert : sa grille ne donne la tonique qu'en
E♭m au pont et en E♭ sur la fin, et l'analyse y entend une tierce
picarde.

#### Ce que le verdict pèse

La tonalité d'un morceau n'est pas un fait mais un verdict, et
l'analyse montre sur quoi il repose : exactement ce que la règle
d'*En Harmonie* pèse (tome 1, chapitre 8 §1.2), et rien d'autre.
`analyse` l'écrit sous l'en-tête :

```
how the tonality is heard:
  first chord    Cm6, bar 1: Cm, its first cadence resolves on it
  last chord     E♭6, bar 35: E♭, the turnaround left out
  heard longest  Cm 27 bars, E♭ 9 bars
```

Le premier accord, et s'il est la tonique là où il sonne ; le dernier,
turnaround exclu, avec sa tierce picarde s'il en a une ; et quand ils
divergent, combien de mesures chacune de leurs deux toniques est
entendue. *My Funny Valentine* se lit ainsi en do mineur, et *Love Me
Or Leave Me* en fa mineur, 14 mesures contre 7 à la♭, là où l'oreille
entend la♭ l'emporter en chemin : ce que la durée seule ne rend pas.
Un blues n'a qu'une ligne, sa forme.

### Ce qui est difficile

**Le plafond des grilles seules.** *In a Sentimental Mood* et *Lullaby
Of Birdland* ont le même profil : une tonique mineure posée d'entrée et
reposée souvent, une fin sur le relatif majeur. La prédominance les lit
tous deux en mineur ; l'app déclare le second en la♭. *It Don't Mean A
Thing* s'arrête sur B♭6, et la mélodie pose d'entrée 1 3 5 de sol
mineur, où la prédominance le lit aussi. Sur des grilles seules, on ne
fera guère mieux ; la mélodie pèsera en plus, quand un format la
portera. C'est d'ailleurs le premier des dix commandements : « Tu
partiras de la mélodie. » Et la grille n'a jamais été la musique.
Chailley le dit des « feuilles bleues que vendent les crieurs de rue »
: « Que dirait-on d'un chef de jazz qui exigerait de son orchestre de
jouer exactement ce qui est imprimé » sur elles, et « que diraient les
musicologues du XXIIIe siècle, en l'absence de disques, s'ils devaient
juger le jazz du XXe siècle d'après lesdites feuilles bleues ? » (*40
000 ans de musique*, p. 263). Une grille iReal est une feuille bleue :
l'analyse en tire ce qu'elle peut, et l'oreille garde le jugement
final, comme au temps où un morceau n'était définitif « que lorsqu'il
avait été "éprouvé", c'est-à-dire soumis à la "preuve" d'une exécution
d'essai » (p. 137).

**L'ouverture sur un accord qui n'est pas la tonique.** *Only Trust
Your Heart* s'ouvre sur Fmaj7♯11, le IV lydien de do, et l'analyse part
de fa avant d'entendre do. Le ♯11 ne tranche pas : neuf standards du
corpus s'ouvrent sur un maj7♯11, trois l'ont pour tonique, trois pour
IV, et le tome 2 confirme que le lydien va très bien sur le I. Le
rythme harmonique ne tranche pas non plus. Reste la mélodie.

**Ce qui est posé** dépend du rythme harmonique : une mesure de Fmaj7
suivie d'une mesure d'autre chose n'est pas posée là où tous les
accords durent une mesure. On a essayé de le mesurer à la carrure, la
fin d'un groupe de quatre mesures : elle pèse trop, elle pose l'E♭maj7
de la mesure 4 de *Jordu*, qui fait de do mineur sa référence dès la
mesure 2 et y revient sans cesse. Le seuil reste « plus d'une mesure et
plus long que sa préparation ».

**Le relatif.** *Autumn Leaves* tonicise d'abord si♭, puis se pose en
sol mineur, où la première phrase s'arrête et où le morceau s'arrête :
le relatif mineur, que l'analyse trouve parce qu'elle attend le repos. *Corcovado*
(la mineur ou do) reste ambigu : la tradition tranche parfois là où
l'oreille hésite.

**L'accord renversé sur sa quinte.** Le F/C de *My Way* n'est pas une
tonique, mais le I sur lequel un V se résout en reste une, renversé ou
non, comme le livre le montre avec ses pédales de dominante (Dm9/G G7
C6/9/G dans *My Romance*). L'analyse fait la différence : voir « Ce qui
fait une tonique ».

**Le m7 d'ouverture.** L'Am7 de *Fly Me To The Moon*, le Dm7 de *Satin
Doll* : un II ou un VI, pas une tonique, tant qu'aucune cadence mineure
n'y résout. Le turnaround qui y ramène à la fin de la grille ne suffit
pas.

### Ce qu'on a décidé

**L'armure n'est pas lue.** C'est un attribut de la partition écrite,
et au mieux le plus faible des indices : les fiches du livre se lisent
à l'identique sans elle. La tonalité que l'analyse conclut n'en dépend
pas, et `analyse` signale quand la tonalité déclarée par la grille iReal
diffère : Tune Up, déclaré en si♭, commence et finit en ré. La
tonalité déclarée se trompe parfois, et pour de bonnes raisons : les
thèmes modaux sont écrits sans armure, un morceau mineur est déclaré
dans le relatif majeur qui a la même armure, et les grilles de jeux
vidéo, relevées par des étudiants, sont moins sûres.

**Dans quelle tonalité compter les degrés** est pourtant une décision à
prendre, puisque les accords seuls ne suffisent pas toujours. C'est
donc à qui lit la grille de la prendre. `analyse` compte par défaut
dans la tonalité que l'analyse conclut : elle juge comme un analyste,
en connaissance de l'harmonie. Sur demande, elle compte dans la
tonalité déclarée par l'app, un indice parmi d'autres mais celui qu'a
choisi l'auteur de la grille, ou dans celle qu'on lui impose. Une fois
la tonalité fixée, on ne la corrige pas : *Lullaby Of Birdland*
compté en la♭ sur demande garde fa mineur pour tonicisation, là où
l'analyse l'entend en fa mineur.

**Une seconde écoute.** Certains blocs ne se lisent bien qu'une fois la
tonique entendue. Lus seuls, G7 Am7 est le ♭VII7 Im de la mineur et G7
Dm7 le IV7 Im de ré mineur mélodique, deux plagales. Là où l'on entend
do, G7 est son V : il va sur son VI, la cadence rompue, ou revient vers
son II, le II-V rejoué de *Satin Doll* (Dm7 G7 | Dm7 G7). L'analyse
écoute donc deux fois : une première fois pour installer la tonique,
une seconde avec ces blocs relus sur elle. Un V qui revient vers le II
d'où il vient se relit aussi hors de la tonique (Em7 A7 | Em7 A7, le
II-V de ré rejoué dans *Satin Doll*), sauf quand l'alternance dure plus
de quatre mesures, un vamp. Là où le mineur est la tonique, la plagale
reste : Fm7 B♭7 dans *Mas Que Nada* est le I IV7 dorien de fa mineur.
Elle reste aussi dans une cadence éolienne : Amaj7 B7 C♯m7 n'est pas
un V de mi qui va sur son VI.
La tonalité du morceau ne change pas, et la seconde écoute est aussi
directe que la première : elle ne demande que la tonalité donnée au
départ et ce qui a sonné.

**Les seuils** sont dans le code, à passer en données quand un second
jeu en aura besoin : deux mesures pour une zone, deux crans pour un
centre éloigné, la moitié d'une section pour une modulation vraie, une
mesure pour un appui. Aucune source ne les chiffre ; ils se règlent de
façon empirique, sur le corpus et les fiches.

Repères : `Sense` (la tonique pressentie), `Grounds` (les régions après
coup), `Phrases`, `Home`, `ReadTune` et `Tune`, `Picardy`, `Reread` et
`Hear` (la seconde écoute). Le drapeau `-key` d'`analyse` choisit la
tonalité où compter les degrés : `heard` par défaut, `declared`, ou une
tonalité comme l'app l'écrit (`F`, `A-`). Sur le corpus, l'analyse
tombe d'accord avec la référence pour l'essentiel des grilles jugées :
le rapport de `corpus` donne le compte du jour.

## 6. Le morceau : blues, plages, structure

L'étage du haut, et celui où le moins est fait : ce qu'on y lit
aujourd'hui, c'est le blues à sa forme, et les plages modales à leur
durée. La structure est le chantier suivant.

### Le blues

#### Dans le livre

Le I7 du blues n'est pas une dominante : c'est une septième d'espèce,
la couleur du degré, qui ne prépare rien. Rien dans son son ne le dit,
c'est sa place qui le dit.

#### Dans gohar

Le blues se reconnaît donc à sa forme : douze mesures jouées une ou deux
fois, ou vingt-quatre en temps doublé, et un squelette volontairement
lâche. Un accord sur la tonique à la mesure 1, quelle que soit sa
qualité (le I7, le Imaj7 du Bird blues, le Im7 du blues mineur), le IV
à la mesure 5, la tonique à la mesure 11. Les mesures 7 à 10, où les
variantes divergent, ne sont pas regardées. Reconnu, le blues a sa
tonique pour fond dès la première mesure ; sans cela, F7 B♭7
installait le IV.

#### Ce qui est difficile

Les formes plus longues qui portent le squelette par hasard (*If I
Loved You*, *I Remember You*) sont écartées, au prix d'un seul blues
manqué, *West Coast Blues*, écrit sur 36 mesures en 3/4. Et les blues
d'une autre forme (*Freddie Freeloader*, *Doxy*, *Watermelon Man*)
attendent qu'on en dresse le catalogue.

### Les plages

#### Dans le livre

Le tome 2 distingue l'harmonie fonctionnelle, où les accords ont un
rôle (préparation, tension, résolution), de l'harmonie modale, où ils
n'ont « plus de rôles tonals », « mais plus un rôle de couleur ». Une **plage modale** est un passage
où « aucun mouvement harmonique n'est rencontré », de longueur libre :
les huit mesures de ré dorien de *So What*, les cinq modes de *Flamenco
Sketches* « sans indication pour le nombre de mesures ». Un même thème
peut mêler les deux, comme *One Finger Snap*, modal dans son A et tonal
dans son B.

#### Dans gohar

Une plage est **tonale** quand les accords s'enchaînent autour d'un
centre et ont une fonction : c'est l'essentiel de ce document. Elle est
**modale** quand un accord tenu plusieurs mesures installe un mode et
non un centre, comme le X7sus4 de *Maiden Voyage*, qui par sa durée,
sans résolution sur son X7, n'a pas de fonction. Elle est **atonale**
quand on passe d'un mode à l'autre au gré des accords (*Pee Wee*) :
chaque accord reçoit sa provenance, aucun ne reçoit de degré.

**La plage modale se reconnaît, son mode ne se lit pas.** Un accord
tenu quatre mesures ou plus, sans cadence qui y mène ni qui en sorte,
fait une plage : les trois de *So What*, les huit de *Maiden Voyage*.
Un turnaround vers le premier accord n'est pas une cadence qui y mène.
Dans une plage, pas de degré ni de tonique pressentie ; `analyse` écrit
« modal » sous l'accord. Un morceau dont les plages font au moins la
moitié est un morceau modal, et le corpus le met à part.

#### Ce qui est difficile

Une grille qui n'écrit que des tétrades ne dit pas la couleur : Dm7 ne
distingue pas ré dorien de ré éolien, et c'est la mélodie, ou ce qu'on
sait du morceau, qui dit que *So What* est dorien. L'analyse le laisse
aux provenances de l'accord et ne choisit pas. *Nardis* est le cas
exemplaire : mi phrygien à l'oreille, le Fmaj7 donnant la ♭2 et le B7
posant mi, et l'analyse entend do sur les retours à Cmaj7, parce qu'elle
ne lit pas encore les cadences modales à deux accords du tome 2
(D♭maj7 Cm7 phrygien, D7 Cmaj7 lydien). La modalité se travaillera sur
de vraies grilles modales, qui portent leurs couleurs.

### La structure

#### Dans le livre

Chaque fiche du livre commence par la forme : AABA, ABAC, AB, le blues,
la forme anatole. Le rythme harmonique, le nombre d'accords par mesure,
en fait partie. Et deux notions dépendent de la structure : la
demi-cadence, une section qui s'arrête sur le V, et le turnaround,
« utilisé en fin de cycle » pour « relancer le thème ».

#### Chez Siron

*La partition intérieure* donne à la structure son cadre. Le temps
musical est fait de niveaux emboîtés, temps, mesures, groupes de
mesures, et au-dessus les formes, sections et morceaux (p. 147). Les
standards groupent leurs mesures en « carrure », des fragments
symétriques de 4, 8 ou 16 mesures : « La forme habituelle des morceaux
de jazz respecte donc la carrure (phrases de 8 mesures, ou plus
rarement de 4 mesures). Le blues possède une carrure particulière (3
phrases de 4 mesures) » (p. 146 et 147).

Dans une musique carrée, **les cadences conclusives marquent les
groupes de mesures** (p. 390). Dans une phrase de 8 mesures, la
cadence conclusive aboutit le plus souvent en mesure 7, la « forte »,
parfois en mesure 8, la « faible ». L'accord de conclusion ouvre alors
la **cadence-boucle**, le turnaround, qui se résout sur le premier
accord du groupe suivant : « une anacrouse de la phrase harmonique
suivante » (p. 353). Dans une forme AABA, c'est elle qui distingue le
plus souvent les A entre eux.

La forme song AABA a des sections de 8 mesures, 32 en tout, et le pont
contraste « avec souvent un changement de tonalité ». Ses variantes
changent la fin des A, la longueur (*Alone Together*, 14 + 14 + 8 + 8)
ou l'ordre (ABAC, AABC). Le blues instrumental est une forme carrée de
12 mesures en trois phrases de 4 (p. 408 et 488).

#### Dans gohar

`Sections` trouve les sections d'une grille par ses seuls accords. Les
marques de section d'une grille iReal ([A], [B]) ne sont pas lues :
comme la tonalité, la structure se reconnaît, et les annotations de la
grille ne pourront servir qu'à affiner, en option.

**Une section, c'est un passage qui revient.** Chaque mesure devient ce
qui y sonne, ses accords, leur place, leur fondamentale et leur
tétrade, et deux passages sont les mêmes quand leurs mesures le sont,
à une transposition près. Tout passage de 4 mesures ou plus qui revient
est une récurrence. Dans *It Don't Mean A Thing*, les mesures 9 à 12
reprennent les mesures 1 à 4 : une section commence en 9.

**La forme grandit depuis la première mesure.** La première mesure
commence une section ; on place d'abord la plus longue récurrence qui
commence là où une section est connue, et chaque occurrence placée dit
où d'autres commencent. Sans cet ancrage, la plus longue récurrence
l'emporterait même mal placée : dans *Autumn Leaves*, la fin du A1 et
tout le A2, douze mesures, reviennent à l'identique dans le B et le
début du C. Une récurrence transposée doit toujours commencer sur une
section connue : de courtes chaînes de dominantes se ressemblent à la
quinte, et découperaient la forme au hasard.

**La grille dit la longueur d'une section.** Quand une occurrence est
suivie aussitôt de la suivante, comme A1 de A2, leur écart est la
longueur de la section, même si leurs dernières mesures diffèrent :
c'est la cadence-boucle de Siron. Chaque occurrence couvre alors cette
longueur, et pas plus : dans *Billy Boy*, le A' dont les deux dernières
mesures préparent le pont reste un A de 16 mesures, et le pont une
section à part. L'écart entre la dernière occurrence et la fin de la
grille ne compte pas : une grille qui commence par une levée de deux
mesures coupe d'autant son dernier A. Ce qui ne revient
nulle part est une section à part, sauf une ou deux mesures avant la
première section, qui sont une levée (notée « - »).

**La carrure ne fait que départager.** C'est la seule connaissance d'un
style que la forme utilise, pour que le mécanisme reste valable hors du
jazz. Une récurrence à deux mesures d'un multiple de 8 en prend la
longueur ; une récurrence de moins de 8 mesures ne peut pas commencer
hors de la grille de 8, et celle qui commence dessus couvre 8 mesures
quand elle en a la place, ses dernières mesures étant une cadence-boucle
qui diffère (les A de *There Is No Greater Love*, identiques sur cinq
mesures seulement). Une plage d'au moins 16 mesures, multiple de 8, qui
ne revient nulle part se coupe en sections de 8 : le B et le C
d'*Autumn Leaves*, et un thème de 16 mesures d'un seul tenant comme
*Blue Bossa*, que les grilles marquent en deux sections de 8.
La grille se cale sur la levée ou l'introduction de chaque morceau. Une
longueur que la grille dit elle-même n'est jamais arrondie : les A de
14 mesures d'*Alone Together* restent des A de 14.

#### Ce qui est difficile

Sur 1350 standards, `Sections` lit un AABA dans un tiers d'entre eux
et un ABAC dans un sur sept. Restent :

- **les grilles sans reprise** : un blues, dont les trois phrases ne
  reviennent pas, reste une seule section ; un thème de 16 ou 32
  mesures sans reprise se coupe en sections de 8, ce qui est juste
  pour *Blue Bossa* ou *Stella By Starlight*, mais cache les reprises
  de *Just In Time* ou *How Insensitive*, qui diffèrent de plus d'une
  mesure ;
- **les niveaux** : « A16 A16 », c'est souvent un ABAC lu au niveau des
  moitiés, et la séquence du pont du Rhythm Changes (D7 G7 puis C7 F7)
  vit au niveau inférieur ; la forme n'en rend qu'un ;
- **la demi-cadence**, qui s'entendra au contraste entre deux fins
  d'une même section, la première sur le V, la seconde sur le I, la
  question et sa réponse. Un II-V en dernière mesure qui renvoie au
  début de la section suivante n'est pas une demi-cadence, c'est une
  cadence-boucle.

**Comment finit une section.** `Conclusions` donne, pour chaque
section, sa cadence conclusive : la dernière cadence qui se résout sur
une tonique dans ses trois dernières mesures, forte si elle aboutit sur
l'avant-dernière mesure, faible sinon, et l'accord où commence la
cadence-boucle qui la suit. Une section qui s'arrête sur son V, une
demi-cadence, ou qui traverse ses dernières mesures sans se résoudre,
n'en a pas. `analyse` l'affiche en tête de grille, avec la forme.

L'analyse de la tonalité n'utilise pas encore la forme : elle en fera
ses phrases, et la place des cadences conclusives dira où la musique se
pose. `charts/cmd/forms` lit la forme de playlists entières, comme
`corpus` la tonalité.

Repères : `Sections`, `SectionsWith`, `Conclusions`, `Blues`, `Modal`, `IsModal`.

## 7. Ce que l'analyse rend

### La fiche

#### Dans le livre

Le chapitre 10 analyse des standards sur une fiche toujours bâtie de la
même façon : la forme et le nombre de mesures, la tonalité de départ et
d'arrivée, les modulations avec leurs mesures, les emprunts avec leur
gamme d'origine, les cadences repérées, et pour chaque accord la gamme
à jouer. Le tout annoté sur la grille elle-même, en chiffres romains
sous les accords et en crochets au-dessus.

#### Dans gohar

L'analyse produit la même fiche, et l'affiche sur la grille. Trois
supports, dans cet ordre : **le terminal** (`charts/cmd/analyse`), banc
d'essai qui grandit à chaque étape ; **une fenêtre Ebitengine**, une
fois l'analyse validée, avec une police de Real Book, du marqueur noir
sur fond blanc et les chiffrages en indices et exposants ; **une page
web**, le moteur compilé en WASM, pour distribuer et faire connaître le
travail, dans la lignée de l'ancien gohareact.

Dans le terminal, la grille se lit ligne par ligne : au-dessus des
accords, ce que chacun fait au suivant (V→, II→, IV→) et les blocs avec
la tonalité qu'ils annoncent ; sous les accords, leur degré ; sous les
degrés, la tonique que l'oreille entend là où elle change ; et
dessous, les cellules. La légende s'affiche avec le drapeau `-legend`.

### Les couleurs proposées

#### Dans le livre

La gamme de la cible colore l'accord qui la prépare : vers un accord
mineur, le V7 prend ♭9 ou ♭13 ; vers un majeur, 9 et 13 ; la dominante
chromatique prend ♯11, comme le F13(♯11) qui va sur Em11 dans *But
Beautiful*.

#### Dans gohar

Ces règles régulières deviennent des propositions : A7(♭13) vers Dm7,
C13 vers Fmaj7, F13(♯11) vers Em11. Le I majeur est proposé en maj7 le
plus souvent, ou en 6, ou avec la neuvième majeure par-dessus la
septième quand la mélodie est sur la fondamentale. Ce sont des
propositions « en toute logique » : l'original, quand on le connaît,
l'emporte.

Que ces extensions soient des couleurs et non des fautes tient à
l'histoire. La dissonance ne s'entend pas « dans l'acception arbitraire
et inexacte des traités d'harmonie, qui appellent ainsi tout accord
autre que l'accord parfait à 3 sons, mais dans son sens réel :
agrégation non réductible à une perception globale de consonance,
elle-même extensible et variable en fonction de l'évolution du
langage » (Chailley, *40 000 ans de musique*, p. 151-152). La
consonance s'est étendue par « notes étrangères progressivement
assimilées », Debussy introduisant la onzième naturelle et Ravel la
stabilisant (p. 154) ; les extensions du jazz sont l'étape suivante du
même mouvement, et le tome 2 d'*En Harmonie* les entend comme des
appoggiatures des notes réelles. Ce qui reste à rebrancher ici, c'est
le lien avec la tonique pressentie : ce que l'on propose sur une
tonicisation, sur une région, sur un II-V qui ne résout pas (voir
`chantiers.md`).

### L'attente et la surprise

Ce n'est pas encore une chose que l'analyse rend, c'est l'objectif qui
la justifie : un jeu où l'on entend en direct.

Un musicien à l'oreille entraînée entend « on dirait qu'on est en ré
majeur » tant que rien ne le dément, et sursaute quand une couleur belle
et inattendue le détrompe. L'analyse en direct doit faire la même
chose, au même instant : c'est un objectif, tant qu'on n'a pas prouvé
qu'il est impossible. Les **lectures provisoires** sont celles que ce
qui a sonné permet, un II-V annonçant son arrivée avant qu'elle ne
sonne. **La surprise** est l'écart entre l'arrivée attendue et ce qui
arrive : une cadence rompue, un emprunt, une dominante chromatique qui
repart ailleurs. Ce n'est pas une erreur mais une couleur que la
théorie sait nommer, la différence entre « faux » et « monstrueux ».
Pour un jeu, c'est la récompense idéale, qui salue une prise de risque
réussie plutôt que la conformité.

La surprise repose sur l'analyse sans la contraindre : c'est l'analyse
qui dit ce qu'on attendait et ce qui est arrivé, et la surprise ne fait
que mesurer l'écart. Elle a déjà son cas de test à l'envers, le pivot
diminué de Tenderly (mesure 12, un Bdim7 qui peut résoudre sur mi♭
mineur comme sur do, et c'est Cm7 qui arrive) : un mouvement bien
écrit, que *Someday My Prince Will Come* fait aussi, où la surprise ne
doit pas se déclencher. Le reste, la jauge de tension et sa mécanique,
est dans `chantiers.md`.

### Hors périmètre, et ouvert

**Accord modifié ou ajouté.** Savoir si un A7 remplace un Am7 de la
grille d'origine ou s'y ajoute relève presque du travail de l'historien
du jazz ; ceux qui veulent cette culture la trouveront auprès de
professeurs qui en sont des puits. L'analyse dit ce que fait chaque
accord, jamais d'où il vient, et l'édition d'une grille se contente de
réanalyser la nouvelle version.

**Ce qui demande la mélodie**, tant qu'on ne l'a pas : repérer à coup
sûr un accord parallèle qui harmonise une note, ou refuser une
dominante chromatique quand la mélodie tient la ♯11 du X7 (elle
deviendrait la fondamentale du nouvel accord, *Midnight Sun*). Ces
règles iront dans les propositions de réharmonisation, le jour où une
grille aura sa mélodie, avec deux autres commandements pour
garde-fous : « Tu ne réharmoniseras pas comme un sauvage » et « Tu ne
feras pas étalage de ton savoir au détriment du thème ». **L'analyse
rythmique**, que le livre écarte aussi.

**Encore ouvert** :

- les thèmes ambigus entre relatifs, *Corcovado* (la mineur ou do) en
  tête : la tradition tranche parfois là où l'oreille hésite ;
- l'attente d'un diminué, ses quatre toniques possibles, pour le pivot ;
- les seuils de la modulation et de la plage modale, en données, à
  régler sur les morceaux de référence puis à l'oreille ;
- la longueur d'une chaîne de préparations avant qu'elle ne soit plus
  une préparation mais une région ;
- l'ordre exact des lectures quand plusieurs valent, règle par règle,
  confronté aux fiches.

## Annexe : les morceaux de référence

Les morceaux qui ont fixé une règle, à garder comme tests. Les parties
du document y renvoient.

### Les fiches du livre

Le livre analyse des morceaux du corpus iReal ; ses fiches, transcrites
à la main dans `charts/ireal/testdata/fiches`, sont l'oracle. Les degrés
en crochets concordent à 83 sur 83, avec ou sans la tonalité que l'app
déclare, qui donne Tune Up en si♭ et le joue sur 32 mesures avec deux
fins.

| Morceau | La fiche | Ce qu'il fixe |
|---|---|---|
| Tune Up | AA', 16 mesures, ré majeur, modulations en do (5 à 8) et en si♭ (9 à 12) | la modulation à chaque phrase, clé de construction du morceau ; 13 sur 13 |
| Black Orpheus | AB, 32 mesures, la mineur, do majeur 6 à 12 | une région de do qui résiste à A7♭9 Dm7 ; 20 sur 20 |
| There Will Never Be Another You | ABAC, 32 mesures, mi♭ majeur, sans modulation | Cm7 tenu une mesure tonicise le VI sans l'installer |
| Tenderly | ABAC, 32 mesures, mi♭ majeur, huit emprunts, plagales | la tonique pressentie, le I emprunté, le pivot diminué, la marche IIm7-V7 lue VI II7 |

Le tome 2 donne aussi des fiches, plus courtes, sous l'angle modal (un
mode par accord) : *Someday My Prince Will Come*, *The Days Of Wine And
Roses*, *Body And Soul*, *Fall*, *Very Early*, *Peace*, *Re: Person I
Knew*. L'analyse tombe d'accord sur leur tonalité, sauf *Peace*, et
*Fall*, sans centre tonal, est écarté.

### Les autres cas tranchés

| Morceau | Ce qu'il fixe |
|---|---|
| Along Came Betty | « le cul entre deux chaises » : deux tonalités à un demi-ton qui se chevauchent, des toniques tenues une seule mesure, pas de modulation ; l'analyse montre l'hésitation sans la trancher |
| Body And Soul | ré♭ majeur, le début entendu en mi♭ mineur : *En Harmonie* le nomme mi♭ dorien, « la sensation de Mi♭ mineur » l'emportant au début (tome 2), et l'analyse l'entend partir de mi♭ mineur |
| Giant Steps | trois centres, mais des tonicisations de moins d'une seconde à ce tempo : pas de modulation |
| Black Orpheus, le pont | Dm (une triade) tenu deux mesures : une modulation temporaire, ou une tonicisation appuyée, lecture juste de ce cas limite |
| Autumn Leaves | le relatif qui tonicise d'abord : si♭ passé, sol mineur où se pose la première phrase |
| How Insensitive | une première phrase de quatorze mesures, qui revient sur son accord d'ouverture |
| Just Friends | un morceau qui s'ouvre sur son IV, confirme son I au milieu des sections, et ne se pose qu'à la dernière |
| Lullaby Of Birdland, All The Things You Are | la dernière tonique entendue, turnaround exclu, fait la tonalité |
| In a Sentimental Mood, Blue Skies | commencer sur le relatif mineur et finir sur le majeur : en fa (*En Harmonie*, p. 159), en do |
| Yesterdays, Virgo, Unforgettable | le VI où *Yesterdays* s'appuie sans conclure, et sa fin à travers la boucle ; le IV où la grille s'arrête avant le turnaround |
| Someday My Prince Will Come | la tonique installée qui revient renversée sur sa quinte |
| My Way | F/C, le IV sur pédale de tonique, et le retour par une plagale |
| Sugar, Fly Me To The Moon | un turnaround ne dit rien de la tonalité ; *Sugar* finit ouvert sur G7 et s'arrête à travers la boucle, Fly Me s'ouvre sur le VI et s'arrête en do |
| My Lucky Star | le II tonicisé quatre mesures, sans modulation |
| Softly, Summertime | la tonique mineure écrite m7 |
| 'Round Midnight | la tierce picarde, sur le dernier accord seulement |
| Chega De Saudade | le majeur homonyme installé pour de bon : pas une tierce picarde ; le morceau est en ré, mineur et majeur à la fois, mis à part du corpus |
| Satin Doll | le II-V rejoué, Dm7 G7 \| Dm7 G7, un V qui revient sur son II et non un vamp dorien |
| Anthropology | l'anatole et le III-VI-II-V, dont le V va sur le III sans y résoudre |
| Guile's Theme | la cadence éolienne qui résout dans son mode, en mineur, et installe sa tonique : do♯ mineur, le seul II-V-I tonicisant le relatif majeur |
| Pokémon, Route 209 | la cadence éolienne qui résout en majeur, la « cadence Mario » |
| Stolen Moments | la marche d'accords parallèles |
| Nardis | mi phrygien, le Fmaj7 pour ♭2 et le B7 qui pose mi : un mode que l'analyse n'entend pas encore |
| So What, Maiden Voyage | la plage modale reconnue, son mode laissé à la grille qui ne le dit pas |
| Sonnymoon for Two, Chasin' the Trane, Blues For Alice | le blues reconnu à sa forme |

### Le corpus

`charts/cmd/corpus` compare, sur une ou plusieurs playlists, la
tonalité que l'analyse entend à celle que l'app déclare, en groupant
les écarts par relation (relatif, quinte, quarte, homonyme, autre) avec
des indices pour trancher : le nombre de cadences résolues sur une
tonique, la fin sur la tonique entendue, le blues, la tierce picarde,
les plages modales, la grille de jeu vidéo.

La tonalité de référence est celle de l'app, sauf pour les grilles
d'une liste **vérifiée à l'oreille**, `charts/ireal/testdata/keys.txt`
(`My Lucky Star | F`, dans l'orthographe de l'app, `F-` pour fa
mineur) : l'app se trompe parfois, et ce qu'on a vérifié est la donnée
qui vaut. Le rapport compte ces grilles à part. Une erreur courante
vient de l'armure : sol mineur et si♭ majeur s'écrivent tous deux avec
deux bémols à la clef, et *It Don't Mean A Thing*, que la mélodie pose
d'entrée sur 1 3 5 de sol mineur, est déclaré en si♭.

Les grilles auxquelles il manque ce qui dit une tonalité sont **mises à
part**, pas jugées : les morceaux modaux et ceux où aucune cadence ne
se résout sur une tonique, que le rapport reconnaît, et une liste
relue à la main, `charts/ireal/testdata/set-aside.txt`, avec la raison
de chaque titre : les thèmes modaux à accords courts dont la grille
n'écrit pas les couleurs (*Speak No Evil*, *Infant Eyes*, *Nefertiti*,
*Afro Blue*), et les blues d'une forme que `Blues` ne connaît pas
(*Freddie Freeloader*, *Doxy*, *Watermelon Man*), un morceau que rien
n'oblige à trancher (*Chega De Saudade*, en ré, mineur et majeur à la
fois), et les thèmes que ni l'un ni l'autre ne connaissons, sans rien
pour vérifier l'analyse. On résiste à la tentation de gérer des grilles
auxquelles il manque l'information : on travaille sur de vraies
données, et cette liste est celle des grilles à réécrire dans un format
qui porte l'information.

Le rapport compte les grilles où l'analyse tombe d'accord avec la
tonalité de référence, celle de l'app ou celle vérifiée à l'oreille,
et range les écarts par familles, qui pointent les questions ouvertes.
Les écarts ne sont pas tous des erreurs de l'analyse : sur plusieurs
grilles vérifiées, c'est l'app qui se trompait. Le chiffre du jour est
celui du rapport, pas de cette page, qui serait fausse dès la livraison
suivante.

Le rapport se termine par **la forme** : les sections que `Sections`
trouve par les accords, comparées aux marques de section de la grille
([A], [B]), que l'analyse ne lit pas. La comparaison se fait lettres
renommées dans l'ordre d'apparition, pour qu'un AABA marqué avec
d'autres lettres reste un AABA ; une section transposée y compte comme
une autre, puisque c'est ce que la grille marque (le pont de *So What*
est son A un demi-ton plus haut). Chaque grille est rangée en même
forme, mêmes sections sous d'autres lettres, sections à deux mesures
près (une levée, une cadence-boucle), même forme à un autre niveau
(A16 B16 marqué, ABAC trouvé : les moitiés d'un côté, les phrases de
l'autre), autre forme, ou pas de marques. Une levée est mise de côté
avant de comparer : la grille de *I Should Care* compte sa levée dans
son premier A, la détection la met à part.

Il se termine par **les phrases** : là où l'analyse se repose
aujourd'hui (`Phrases`), comparé à là où les sections trouvées
concluent (`Conclusions`), la règle de Siron appelée à remplacer la
nôtre. Avec `-phrases`, il liste les grilles où les deux divergent,
mesure par mesure. Comme
l'armure, les marques sont souvent absentes ou approximatives : un
désaccord ne dit pas qui a tort. Avec `-forms`, le rapport liste les
grilles en désaccord, la forme marquée au-dessus de la forme trouvée.
