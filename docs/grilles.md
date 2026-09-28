# L'analyse des grilles

Comment gohar lit une grille : tonalités, cadences, emprunts,
modulations, et la gamme qui va sur chaque accord.

La progression suit le tome 1 d'*En Harmonie* (Dericq et Guéreau,
Outre Mesure), chapitres 8 à 10, et « le livre » désigne ici ce
tome-là. Le code vit dans `harmony/analysis` ; `charts/cmd/analyse`
affiche une grille annotée dans le terminal, et `charts/cmd/corpus`
compare l'analyse à l'app sur des playlists entières.

## À quoi elle sert

Le livre donne trois raisons d'analyser un morceau, et gohar les
reprend dans cet ordre : **jouer juste**, la gamme ou le mode qui va
sur chaque accord, ce qui rejoint le dex ; **modifier l'harmonie** en
connaissance de cause, en sachant ce que fait chaque accord ;
**mémoriser** la grille, puisqu'une suite de cadences se retient mieux
qu'une suite d'accords.

Le public visé est l'autodidacte, avec des questions pratiques
immédiates. L'analyse lui dit ce que fait chaque accord ; elle ne lui
raconte pas l'histoire du morceau.

## Les principes

- **L'analyse propose, elle ne tranche pas.** Le mieux reste d'écouter,
  de relever et de connaître l'original. À défaut, gohar suggère ce qui
  se déduit « en toute logique », et le présente comme tel.
- **Déterministe, sans score.** Toutes les lectures valables sont
  rendues, dans un ordre fixé par des règles explicites, qui dit
  laquelle est proposée d'abord, jamais laquelle est vraie. Em7-A7 est
  un II-V de ré ou un III-VI7 de do : deux représentations exactes de
  la même progression, sans hiérarchie de niveau.
- **De droite à gauche.** On repère les accords réels de la tonalité,
  on prend l'accord d'arrivée, puis on remonte ses préparations. C'est
  aussi l'algorithme, et l'arrivée départage ce que rien d'autre ne
  départage : Fm7 B♭7 est un II-V de mi♭ ou un IVm7-♭VII7 de do selon
  l'accord qui suit.
- **La basse compte.** Elle distingue une cadence parfaite d'une
  imparfaite, un I/3 d'un III, un Dm9 d'un G13sus4.
- **La durée compte.** Elle sépare une plage modale d'un accord de
  passage, une modulation d'une tonicisation. L'analyse travaille donc
  en temps, pas en cases iReal.

## Les étages

| Étage | Ce qu'il porte | Exemple |
|---|---|---|
| Le morceau | La structure, le rythme harmonique, la tonalité | Tenderly : ABAC, 32 mesures, un accord par mesure, mi♭ majeur |
| La plage | Tonale, modale ou atonale | *One Finger Snap* : modale mesures 1 à 12, tonale ensuite |
| La région | Une tonalité installée : la principale, ou une modulation | *Black Orpheus* : la mineur, do majeur mesures 6 à 12, la mineur |
| Le bloc | Une cadence, avec la tonalité qu'elle annonce | Gm7♭5 C7(♭9), qui annonce fa mineur |
| L'accord | Son degré, sa fonction, sa provenance, sa couleur proposée | D♭7 : ♭VII7, emprunté à la♭ mineur mélodique |

La région se lit après coup, de droite à gauche. Elle a un pendant qui
se lit au présent, de gauche à droite : la tonique pressentie, ce que
l'oreille attend accord par accord.

## Le morceau

### Sa tonalité

Deux questions, que l'analyse sépare : d'où le morceau part, et où il
s'arrête.

**La maison** (`Home`) est là où se pose la première phrase. Une
phrase n'a pas de longueur fixe : elle court jusqu'à se poser sur un
accord de tonique qu'une cadence amène, et qui revient sur l'accord
d'ouverture (*How Insensitive*, un long soupir de Dm à Dm, quatorze
mesures plus loin) ou tient plus d'une mesure et plus longtemps que les
accords qui y mènent (le Gm6 d'*Autumn Leaves*). Une tonique de passage
ne pose rien (le B♭maj7 d'*Autumn Leaves*), ni un IV, si long soit-il
(le E♭maj7 de *Cherokee*), ni un m7 après l'ouverture (le Cm7 de *There
Will Never Be Another You*, son VI). Un morceau peut aussi s'ouvrir au
repos, sur sa tonique tenue plus d'une mesure, sauf quand la première
phrase se pose ensuite une quinte au-dessus : *Just Friends* s'ouvre
sur Cmaj7, son IV, et se pose sur Gmaj7. Tant que rien ne s'est posé,
le premier accord n'est qu'une hypothèse.

**La tonalité du morceau** (`Tune`) est la dernière tonique entendue,
là où s'arrête la dernière phrase. On joue un standard jusqu'à elle, et
pas au-delà : le turnaround qui suit ramène au premier accord, vers
lequel il pointe forcément, et ne dit rien de la tonalité. *Lullaby Of
Birdland* part de fa mineur et s'arrête sur A♭maj7, avant que Gm7♭5 C7
ne ramène à Fm : il est en la♭. *All The Things You Are* n'est tranché
que par son dernier A♭maj7. Pour arrêter le morceau hors du fond, il
faut un II-V-I : un V seul ne fait que traverser (*Yesterdays* passe
par B♭maj7 dans son cycle de dominantes, et reste en ré mineur), sauf
sur le tout dernier accord de la grille. Un m7 ne l'arrête que s'il
est la tonique du fond : le F♯7 Fm7 de la fin de *Sugar* tonicise le
IV. Certaines grilles n'écrivent pas leur dernière tonique et
s'arrêtent sur le IV de la maison, avant un turnaround qui y revient :
*Unforgettable* s'arrête sur Cmaj7 avant Am7 D7, et ce IV est sauté.

**Un morceau qui s'ouvre au repos et que sa première cadence ramène
chez lui** a installé sa maison avant d'en partir, et il y reste où
qu'il s'arrête. *In a Sentimental Mood* tient Dm deux mesures (Dm,
Dm(maj7), Dm7, Dm6), y revient par A7, et reste en ré mineur bien qu'il
conclue par Gm7 C7♭9 Fmaj7, dans son relatif majeur ; c'est aussi ce
que dit la tradition. *Blue Skies* s'ouvre sur la même ligne depuis
Am, mais sa première cadence va à C6 : il est en do.

Quand la maison et la fin diffèrent, `analyse` donne les deux : *Fly Me
To The Moon* part de la mineur et s'arrête en do, *Lullaby Of Birdland*
part de fa mineur et s'arrête en la♭.

**La tierce picarde** ne rend pas majeur un morceau mineur. Héritée de
la musique d'église, où un accord majeur sonne avec moins de partiels
qui frottent sous la résonance d'un grand orgue, elle majorise la
tonique sur le dernier accord : la tonique, mineure jusque-là et
exclusivement, se majorise sur l'accord de fin. Le morceau reste
mineur, et l'analyse la signale (`Picardy`). Une section peut se clore
de la même façon avant, si le mineur revient ensuite : *'Round
Midnight* finit son deuxième A comme il finit, sur E♭6, et son dernier
A repart sur E♭m. Tout autre majeur homonyme est une modulation :
*Chega De Saudade* est en ré mineur pour ses deux premières parties, en
ré majeur pour les deux dernières, s'arrête sur D6, et est en ré
majeur. De même
*I Love Paris* et *Black And Tan Fantasy*, dont la seconde partie est
en majeur. Un morceau qui finit sur sa propre tonique garde le mode où
il finit, même s'il s'est ouvert au repos dans l'autre.

**L'armure n'est pas lue.** C'est un attribut de la partition écrite,
utile pour ne pas écrire des altérations partout, et au mieux le plus
faible des indices d'une tonalité. Le livre la recommande à l'apprenant
qui se creuse la tête sur une partition ; pour un analyste qui ne se
fonde que sur ses propres relevés, elle n'existe pas. « Tu ne suivras
pas bêtement les indications du Real Book. » Les fiches du livre se
lisent à l'identique sans elle. Le champ de tonalité d'une grille iReal
n'est qu'un point de comparaison, et l'analyse signale quand il se
trompe : Tune Up, déclaré en si♭, commence et finit en ré. Il se trompe
souvent pour de bonnes raisons : les thèmes modaux sont écrits sans
armure, comme en do majeur, et les grilles de jeux vidéo, relevées par
des étudiants et parfois modales, sont moins sûres.

### Le blues

Son I7 est une septième d'espèce, pas une dominante, et rien dans son
son ne le dit : c'est sa place qui le dit. Le blues se reconnaît donc à
sa forme (`Blues`), douze mesures jouées une ou deux fois, ou
vingt-quatre en temps doublé, et à un squelette volontairement lâche :
un accord sur la tonique à la mesure 1, quelle que soit sa qualité (le
I7, le Imaj7 du Bird blues, le Im7 du blues mineur), le IV à la mesure
5, la tonique à la mesure 11. Les mesures 7 à 10, où les variantes
divergent, ne sont pas regardées. Reconnu, le blues a sa tonique pour
fond dès la première mesure ; sans cela, F7 B♭7 installait le IV.

Sur le corpus, 53 grilles sont reconnues, toutes des blues. Les formes
plus longues qui portent le squelette par hasard (*If I Loved You*,
*I Remember You*) sont écartées, au prix d'un seul blues manqué, *West
Coast Blues*, écrit sur 36 mesures en 3/4.

### Sa structure et ses plages

Le rythme harmonique est le nombre d'accords par mesure. La structure
(AABA, ABAC, AB, blues, forme anatole) se lit dans les sections de la
grille et se vérifie en comparant les progressions entre elles
(`Progression.Compare`).

Une plage est **tonale** quand les accords s'enchaînent autour d'un
centre et ont une fonction : c'est l'essentiel de ce document. Elle est
**modale** quand un accord tenu plusieurs mesures installe un mode et
non un centre, comme le X7sus4 « d'espèce » de *Maiden Voyage*, qui
par sa durée, sans résolution sur son X7, n'a pas de fonction ; le
critère est une durée sans cadence, au-delà d'un seuil. Elle est
**atonale** quand on passe d'un mode à l'autre au gré des accords (*Pee
Wee*) : chaque accord reçoit sa provenance, aucun ne reçoit de degré,
et c'est la valeur zéro de `Tonality`, un état légitime. Un même
morceau peut mêler les trois. Les plages modales ne sont pas encore
codées : *So What* (ré dorien, mi♭ dorien, ré dorien) en sera le test,
et son pont est aujourd'hui chiffré ♭IIm7 en ré mineur, ce qui n'a pas
de sens.

## Les cadences

### Les préparations

Préparer un accord, c'est ajouter ou modifier un ou deux accords devant
lui. L'anglais dit *approach* : `ApproachKind` nomme le type d'une
préparation, et `ApproachOf` dit comment un accord prépare le suivant.
L'unité de l'analyse est une **arrivée** et ce qui la prépare.

| Type | Chiffrage | Exemple vers Dm7 en do | Codé |
|---|---|---|---|
| Dominante secondaire | V7 de…, écrit V7/II ou VI7 | A7 | oui |
| Dominante chromatique | ♭II7 de…, un X7 un demi-ton au-dessus | E♭7 | oui |
| Accord diminué | un demi-ton sous l'arrivée | C♯dim7 | oui |
| X7sus4 | retarde son propre X7 | A7sus4 A7 | oui |
| Sous-dominante secondaire | II-V de…, IIm7 ou IIm7♭5 | Em7 A7 | oui |
| Sous-dominante chromatique | ♭VIm7 ou ♭VIm7♭5 devant le ♭II7 | B♭m7 E♭7, ou B♭m7 A7 | oui |
| Plagale | le IV de toute qualité, ou le ♭VII7, devant une tonique | Gm7 Dm, C7 Dm | oui |
| IVm7-♭VII7 de… | la plagale préparée | Gm7 C7 Dm | oui |
| Accord parallèle | même qualité, une seconde à côté | E♭m7 ou C♯m7 | non |

**Les relations se composent.** La sous-dominante chromatique est un II
de… appliqué à une dominante chromatique, la sous-dominante secondaire
un II de… appliqué à une dominante secondaire. Le modèle n'a que
quelques relations élémentaires qu'on enchaîne, et le nom composé sert à
l'affichage. Chez Ellington, G♭7 F7 E7 E♭7 A♭maj7 est une chaîne de V7
de… dont un accord sur deux est remplacé par son ♭II7 pour faire
descendre la basse. Ce qu'on en tire est une **forme sous-jacente**, pas
une histoire : dans *I Thought About You*, chaque X7 reçoit sa cible et
la chaîne de II-V par quintes réapparaît sous la grille écrite.

**Chaque X7 pose une question**, que le livre formule : est-il la
dominante de la tonalité, la dominante secondaire de l'accord suivant,
ou sa dominante chromatique ? Il peut aussi être un X7sus4 de passage,
une septième d'espèce (le blues), ou, tenu longtemps, une plage modale.

**Le X7sus4**, sans tierce, n'a pas de triton et n'est pas une
dominante. En cadence, c'est une position d'attente : il prolonge ou
remplace le II avec la fonction de sous-dominante, et se résout sur son
X7. Sans résolution et tenu longtemps, il installe une plage modale.

**La dominante chromatique** partage son triton avec la dominante un
triton plus loin : toute cible a deux dominantes, que `Resolution`
dérive sans table. Le livre en donne la raison par l'accord de sixte
augmentée : en cinq étapes, G7 devient D♭7, et chaque voix rejoint
l'accord de tonique par demi-ton, la règle de moindre mouvement de
`voicings.md` poussée à son maximum. On dit « dominante chromatique »
plutôt que « substitution tritonique ».

**L'accord diminué** a deux emplois. Un demi-ton sous l'arrivée, c'est
son V7(♭9) sans fondamentale : C♯dim7 prépare Dm7, et comme chaque
diminué a quatre fondamentales, il prépare aussi toute cible un
demi-ton au-dessus de l'une de ses notes. Entre deux accords, c'est un
**accord de passage** sur une basse chromatique, que seule l'analyse
d'une grille voit, puisqu'il faut trois accords et leurs basses
(`PassingChords`) ; il est alors nommé d'après sa basse, haussé en
montant, abaissé en descendant. En montant (I ♯Idim7 II, *Mean to Me*),
les deux emplois coïncident et les deux lectures sont rendues ; en
descendant (B♭/D D♭dim7 Cm7), il ne reste que le passage.

**La plagale** conclut autant qu'un V-I : le F/C C final de *My Way*
est un gros amen sur do. Mais elle attire moins : V-I peut faire
changer de tonalité avec une force d'attraction maximale, IV-I ne le
fait pas. Le ♭VII7-I en est le faux nez mineur : avec le IV à la
basse, B♭7 devient Fm6 en do. Ce que la plagale fait de la tonalité en
découle (voir « La tonique pressentie »).

### Les blocs

Un bloc est une cadence : **[II] [sus4] V → cible** (`Blocks`). Le V
est le seul accord obligatoire, et prépare la cible comme dominante,
dominante chromatique ou diminué ; une plagale fait aussi un bloc, sa
sous-dominante à la place du V, avec son IVm7 pour II quand c'est le
♭VII7. La cible n'est pas dans le bloc : il la désigne.

Le découpage se fait de droite à gauche. Chaque accord appartient à un
seul bloc, et un bloc peut viser un accord d'un autre : A7 Dm7 G7 Cmaj7
donne [Dm7 G7] → Cmaj7 et [A7] → Dm7, le V7/II puis le II-V-I, et le
double rôle du bebop en découle (Dm7 cible d'un bloc et II du suivant).
Un II et un V qui ne préparent pas l'accord suivant font un bloc sans
cible, le **II-V sans résolution** du livre, qui annonce quand même sa
tonalité ; un V seul qui ne prépare rien n'en est pas un.

**La tonalité annoncée** (`Block.Announced`) est une tonique et les
gammes de la pratique tonale (majeure, mineures naturelle, harmonique,
mélodique) qui contiennent toute la préparation : un IIm7 garde le
majeur et le mineur mélodique, un IIm7♭5 le mineur harmonique, un
V7(♭13) les mineurs. La tierce de la cible départage ensuite, si la
préparation le permet : Gm7 C7 → Fm7 annonce fa mineur mélodique, Gm7♭5
C7 → F reste fa mineur harmonique, le cas mixte de *What Is This Thing
Called Love*, où le bloc porte fa mineur et le I porte fa majeur. Une
dominante chromatique et son II ne sont pas des degrés de la tonalité :
la cible décide seule. La majeure harmonique contient aussi certaines
préparations, mais c'est une gamme d'emprunt, pas une tonalité qu'on
annonce.

Une chaîne de II-V est une suite de blocs, chacun avec sa tonalité
annoncée, et la **marche** est une relation entre eux : « II-V en mi
mineur, puis en ré mineur, puis en do mineur qui se résout sur do
majeur ». Le livre appelle **Ier degré temporaire** l'arrivée d'un bloc
emprunté (Misty : A♭maj7 préparé par B♭m7 E♭7) : c'est une
tonicisation.

Le premier pas d'une **marche IIm7-V7**, dont le V7 devient le II du
bloc suivant sur la même fondamentale, n'est pas mis entre crochets :
Cm7 F7 Fm7 B♭7 en mi♭ (Tenderly mesures 13 à 16, There Will Never Be
Another You mesures 12 à 16) se lit VI II7 II V, comme le livre. F7 est
un « IIe degré altéré », pas le V d'un si♭ qui ne vient jamais ; le bloc
annonce quand même si♭ au-dessus des accords.

### Les accords parallèles

Ils n'ont pas de fonction, et sont de deux sortes. **Un accord qui
harmonise une note de la mélodie**, bref et de même qualité une seconde
à côté de l'arrivée, demande la mélodie pour être sûr : sans elle,
gohar ne peut que le conjecturer. **Une marche d'accords parallèles**,
elle, se lit sur la grille seule : au moins trois accords de même
qualité reliés par des mouvements conjoints, sous une mélodie qui ne
les suit pas. Le pont de *Stolen Moments* (Dm D♯m | Em Fm | F♯m Fm |
Em E♭m) fait monter et descendre par demi-tons des accords doriens
sous un ostinato sur un intervalle de tierce : c'est une couleur,
chaque accord vient de son propre mode, et le fond ne bouge pas. Un
accord qui appartient à un bloc n'en fait pas partie : les X7 de
*Sophisticated Lady* se préparent l'une l'autre, et dans *Along Came
Betty*, B♭m7 Bm7 B♭m7 Bm7 a la forme d'une marche mais chaque Bm7 est
le II de E7. Pas encore codé.

## La provenance

La gamme d'où vient un accord, c'est ce que le joueur improvise dessus
(`harmony.Provenances`). Elle se calcule sans table parmi toutes les
gammes nommées qui contiennent l'accord, extensions comprises : C7
vient de sept gammes, C7(♭9) de trois. Dans Tenderly, D♭7 est emprunté
à la♭ mineur mélodique, Gm7♭5 C7(♭9) à fa mineur harmonique, Bdim7 à do
mineur harmonique. L'ordre de présentation met d'abord la tonalité de
la région, puis celle qu'annonce le bloc, puis les autres, toutes
rendues : le livre dit du F7 de Tenderly que c'est « un IIe degré
altéré pouvant correspondre à différents modes ».

Une **variante** de cadence est un squelette de degrés et une
provenance par degré : les formes de l'anatole en fa (majeure, mineures,
mélangées) ne sont pas des progressions distinctes, c'est I-VI-II-V où
chaque degré prend sa tétrade dans l'une des gammes permises.

## Le chiffrage

Deux lectures, simultanées. **Sur la tonalité du passage**
(`Degrees`), chaque accord par son degré dans la région où il sonne :
en mi♭, Gm7♭5 C7♭9 avant Fm7♭5 est IIIm7♭5 VI7, le III-VI d'un
III-VI-II-V-I. **En crochets** (`Bracketed`), comme le livre les
imprime : un II-V se chiffre relativement à sa cible, dans la tonalité
qu'il annonce, et la même mesure est II V de fa mineur. Sans le
crochet, un « II V » se lirait comme un II-V de la tonique ; c'est
pourquoi `analyse` donne la première lecture sous les accords et la
seconde sur la ligne des blocs, qui sert de crochet, quand leurs degrés
diffèrent. Les fiches du livre se comparent à la seconde. Une plagale
ne se lit pas en crochets : son IVm7 ♭VII7 se chiffre sur la tonique
qu'il conclut.

Les règles communes : un degré diatonique s'écrit sans qualité, un
degré emprunté avec (II, mais IVm7), et le degré se compte dans la gamme
de la tonalité (en fa mineur, D♭maj7 est VI, pas ♭VI). Une fondamentale
hors de la gamme est abaissée de préférence (♭II, ♭III, ♭VI, ♭VII) et
haussée sous la quarte et la quinte (♯IV) ; un accord de passage suit
sa basse. Les renversements s'écrivent I/3, I/5. Un V seul, un diminué
seul et un accord de passage se chiffrent sur la tonalité du passage,
avec leur qualité : VI7 pour la dominante secondaire de II, écrite
aussi V7/II (gohar garde la relation, l'écriture est un choix
d'affichage). Une dominante chromatique se chiffre par son degré (♭II7
vers le I, ♯IV7 ou ♭III7 ailleurs).

## La tonique pressentie

À chaque accord, la tonique que l'oreille attend, avec ce qui a sonné
et rien d'autre (`Sense`). Au troisième accord de Tenderly, deux
mesures de E♭maj7 A♭7 ont installé mi♭ : E♭m7 s'entend comme la tonique
qui change de couleur, pas comme le II de ré♭. C'est le pendant, au
présent, de la région : la fiche montre les régions, le direct la
tonique pressentie, et l'écart entre ce qu'elle attendait et ce qui
arrive est la surprise.

### Deux toniques

**La tonique de fond** est installée : les degrés se comptent sur elle,
même quand une cadence tonicise un autre degré (Dm7♭5 G7 Cm7 en mi♭ :
II V VI). **La tonique locale** est celle qu'une cadence vient de
toniciser ; elle dure tant que les accords suivants tiennent en elle ou
préparent un accord qui y tient, et ne change pas le fond. Chacune est
un ensemble de tonalités sur une même tonique, comme une tonalité
annoncée, et le fond est vide au début d'un morceau, avant la première
tonique, ou dans une plage atonale. À part, la **tonique de départ** est
gardée en mémoire pour reconnaître le retour à la maison après un pont
qui a modulé, si fréquent sur une forme AABA.

### Ce qui fait une tonique

Un accord de tonique est une triade, un maj7, un 6, un m6 ou un
m(maj7), à l'état fondamental ou avec sa tierce à la basse. Avec sa
quinte à la basse, c'est une quarte et sixte sur une pédale : dans le
F/C C de *My Way*, le fa n'est qu'une broderie au-dessus du do. Un m7
est presque toujours une sous-dominante, et on ne module pas pour
s'installer en éolien. Mais les grilles écrivent la tonique mineure m7
bien plus souvent qu'on ne la joue (m6, m(maj7), m(maj9) pour adoucir
la septième) : un m7 est donc une tonique quand une cadence mineure se
résout dessus et qu'il n'est pas lui-même le II d'un bloc. *Softly, As
In A Morning Sunrise* est ainsi en do mineur. Un turnaround vers un
premier accord en m7 n'en fait pas une tonique : l'Am7 qui ouvre *Fly
Me To The Moon* est un VI.

### Ce qui l'installe

Une **cadence qui se résout** fait de sa cible une tonique locale. Au
début, le **premier accord** installe le fond s'il peut être une
tonique, sinon la première cadence résolue ; c'est l'indice le plus
faible, beaucoup de standards commençant sur un II ou un IV, et après
coup les accords d'avant la première tonique lui appartiennent. Ce fond
n'est qu'une hypothèse : quand la première phrase se pose (voir « Sa
tonalité »), sa tonique devient la maison et le fond. Un **blues**
reconnu a sa tonique pour fond dès la première mesure.

Une **plagale conclut sans ouvrir** : elle confirme une tonique déjà là
(le fond, la tonique de départ, la tonique locale), ramène à la maison
et compte comme une cadence qui confirme une modulation, mais n'ouvre
jamais seule une tonique locale. Ce qu'elle annonce, l'oreille
l'attend : après D♭7 dans Tenderly, on attend mi♭.

Ne change rien au fond : les accords diatoniques, les emprunts sur la
même tonique (Im7, IVm, ♭VII7), les préparations qui ne se résolvent
pas. Un accord dont la fondamentale est la tonique du fond se lit comme
un **I emprunté**, même quand il est le II d'un II-V qui ne se résout
pas : E♭m7 A♭7 en mi♭ est Im7 IV7, et le bloc annonce toujours ré♭.
Le livre écrit ces mesures « I IV » : il note la qualité empruntée une
fois et plus ensuite, l'analyse l'écrit chaque fois.

### La modulation

**L'emprunt** fait venir un accord ou une cadence d'une autre tonalité
sans quitter la sienne : Dm7 G7(♭13) Cmaj7 emprunte son G7(♭13) à do
mineur. **La modulation** change de tonalité pour de bon. Le livre n'en
donne pas une règle unique, et ses exemples montrent trois indices qui
se combinent : la cible est-elle un degré de la tonalité (Dm7♭5 G7♭9
Cm7 dans *There Will Never Be Another You* tonicise le VI), combien de
temps la nouvelle tonique tient (Tune Up, un I tenu deux mesures), et
combien de cadences la confirment (*Black Orpheus*, do majeur sur sept
mesures). gohar les lit au présent, avec un parti pris : **être
libéral**.
Appeler modulation une tonicisation appuyée est une analyse que
beaucoup de musiciens feraient ; ne pas voir bouger les repères tonaux
serait une faiblesse.

Une tonique locale devient le fond quand elle tient **plus d'une mesure
d'accords stables**, ceux qui tiennent en elle sans rien préparer
(Dmaj7 tenu deux mesures, ou B♭maj7 Gm7 dans Tune Up), ou quand **une
deuxième cadence** la vise tant qu'elle dure. On ne module ni vers un
accord de sous-dominante ni pour un seul accord : le I doit pouvoir
être une tonique et ne pas être aussitôt le II d'un autre bloc. Une
cadence vers un degré de la tonique locale qui ne peut pas en être une
(A7 Dm7 quand do est local, dans Black Orpheus) ne l'interrompt pas. Le
relatif n'est pas une modulation plus faible, seulement plus facile.

Le **retour à la maison** est asymétrique : une seule cadence sur la
tonique de départ la réinstalle, parfaite ou plagale, et même son
accord de tonique seul. Quitter demande plus de preuves que revenir.

En direct, le fond bascule au moment où l'indice est rempli ; après
coup (`Grounds`), la région commence au bloc qui y menait, et c'est sur
elle que les degrés se comptent (Tune Up, mesure 7 : Cmaj7 est I, pas
♭VIImaj7). Les seuils (plus d'une mesure, deux cadences) sont dans le
code, à passer en données quand un second jeu en aura besoin. La
modulation **confirmée**, après une phrase entière, qui distinguerait
une modulation passagère d'une vraie région, n'est pas encore codée.

### Local, avec une mémoire

La contrainte du direct (chaque calcul ne regarde qu'un nombre borné
d'accords) tient : la tonique pressentie est un état porté d'un accord
au suivant, un résumé de ce qui a sonné. Seule la lecture d'une grille
entière va plus loin, parce qu'elle le peut : une grille qui boucle est
entendue comme son deuxième chorus, qui part de la fin du premier, avec
pour chez-soi la tonalité du morceau, désormais connue. Un turnaround en fin de
grille prépare donc le premier accord, et Tune Up commence en ré. À la
première écoute, une cadence à travers la boucle n'a pas encore sonné
quand son premier accord sonne ; le direct n'a que cette première
écoute.

### Tenderly, au présent

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
| 9 | Fm7♭5 | mi♭, locale fa mineur | pas Fm, mais le II de mi♭ mineur |

Ce tableau sert de test, validé en attendant l'avis d'une oreille plus
experte.

## L'attente et la surprise

Un musicien à l'oreille entraînée entend « on dirait qu'on est en ré
majeur » tant que rien ne le dément, et sursaute quand une couleur belle
et inattendue le détrompe. L'analyse en direct doit faire la même
chose, au même instant : c'est un objectif, tant qu'on n'a pas prouvé
qu'il est impossible. Les **lectures provisoires** sont celles que ce
qui a sonné permet, un II-V annonçant son arrivée avant qu'elle ne
sonne. **La surprise** est l'écart entre l'arrivée attendue et ce qui
arrive, et elle se qualifie (cadence rompue, emprunt, dominante
chromatique qui repart ailleurs) : ce n'est pas une erreur mais une
couleur que la théorie sait nommer, la différence entre « faux » et
« monstrueux ». Pour un jeu, c'est la récompense idéale, qui salue une
prise de risque réussie plutôt que la conformité.

**La surprise se gradue.** Dans Tenderly, mesure 12, B♭7 fait attendre
mi♭ mineur ; Bdim7 (si ré fa la♭, soit Ddim7) est le VIIdim7 de mi♭
mineur et prolonge l'attente, mais aussi celui de do, et c'est Cm7 qui
arrive. Un tel **pivot** laisse les deux résolutions ouvertes : au
moment du dim7, l'attente devrait porter ses quatre toniques possibles,
et l'arrivée sur l'une d'elles est plus douce qu'une arrivée que rien
n'annonçait. Ce mouvement très courant, bien écrit plutôt que
monstrueux (*Someday My Prince Will Come* le fait aussi), est un test à
l'envers : la surprise ne doit pas s'y déclencher.

**La tension** monte tant que l'oreille est tenue loin d'une tonique et
retombe quand elle y arrive : une jauge qui se remplit. Le [B] de
Tenderly en est le modèle, avec ses cadences qui ne se résolvent pas et
qui ne trouvent la maison qu'au retour du [A]. Elle se déduit de la
tonique pressentie, de façon déterministe (un compte, pas une
probabilité) : elle monte à chaque cadence qui ne se résout pas, à
chaque attente déçue, à chaque accord qui ne tient pas dans le fond ;
elle retombe en partie sur une tonique locale, entièrement au retour à
la maison.

## Les couleurs proposées

La gamme de la cible colore sa préparation, selon des règles
régulières qui deviennent des propositions : vers un accord mineur, un
V7 avec ♭13 ou ♭9 (A7(♭13) vers Dm7) ; vers un accord majeur ou de
dominante, 9 et 13 (C13 vers Fmaj7) ; la dominante chromatique avec
♯11 (F13(♯11) vers Em11 dans *But Beautiful*) ; le I majeur en maj7 le
plus souvent, ou, quand la mélodie est sur la fondamentale, en 6 ou
avec la neuvième majeure par-dessus la septième. Ce sont des
propositions « en toute logique » : l'original, quand on le connaît,
l'emporte.

## Le catalogue v0

Des données, relues ligne à ligne comme la table des qualités iReal.

| Cadences | Formes | Codées |
|---|---|---|
| À deux accords | parfaite (V-I à l'état fondamental), imparfaite (un renversement), demi-cadence (…-V), rompue (V-VI) | parfaite et imparfaite ; la demi-cadence et la rompue se définissent par une absence, et attendent |
| II-V-I | majeur IIm7-V7-Imaj7 ; mineur harmonique IIm7♭5-V7(♭9, ♭13)-Im(maj7) ; mineur mélodique IIm7-V7(9, ♭13)-Im(maj7) ; mixte IIm7♭5-V7(9, ♭13)-Imaj7 | oui |
| Plagales | IVmaj7-I ; IV7-I et IV7-Im (mineur mélodique) ; IVm7-I (mineur harmonique) ; IVm(maj7)-I (majeur harmonique) ; IV-IVm-I ; IV7-I7 (« bluesy ») | toutes, sauf IV7-I7 (le I7 n'est pas un accord de tonique) |
| Modales ♭VII-I | ♭VII7 (éolien), ♭VIImaj7 (mixolydien), ♭VIIm7 (phrygien) ; préparées II-♭VII7-I, IV-♭VII7-I, IVm7-♭VII7-I | ♭VII7, lu comme une plagale mineure, et IVm7-♭VII7-I |
| Avatars du V | V7, ♭II7, VIIdim7 | oui |

**Les sous-dominantes**, d'après le tableau du chapitre 9 : IIm7,
IIm7♭5, II7, ♭II7 ; IVmaj7 ou IV6, IVm7 ou IVm6, IVm7♭5, IV7, ♯IVdim7 ;
♭VImaj7, VIm7♭5 (lu comme le ♯VIm7♭5 du mineur mélodique), ♭VI7 ;
♭VII7.

**Les cellules** tournent au lieu d'aboutir, ce qui les distingue des
cadences. L'anatole I-VI-II-V, en majeur, en mineur ou mélangée, se
joue sur une durée non définie (intro, coda), à ne pas confondre avec
la forme anatole des *rhythm changes*. III-VI-II-V-I remplace le I par
le III. Les chaînes de II-V vont par demi-ton, par ton ou par cycle des
quartes, où le I potentiel devient le II suivant. Le turnaround prépare
le retour au début sur une ou deux mesures et traverse la barre
finale. Pas encore reconnues : elles se liront sur la ligne des degrés.

**Le I qui dure** a ses couleurs (triade, maj7, 6, le mouvement oblique
Imaj7-Imaj7(♯5)-I6), ses prolongements (IVmaj7, IVm(maj7), IVm6/I,
V7sus4, ♭VII7, ♭VIImaj7, I-II-III-V), la ligne gospel
I-I7/3-IV-♯IVdim7-I/5, et les lignes chromatiques du Im (Im, Im(maj7),
Im7, Im6 depuis la fondamentale ; Im, Im(♭6), Im6, Im7 depuis la
quinte).

## Les lectures qui se recouvrent

Cas où plusieurs lectures sont rendues, à garder comme tests : Em7-A7
(II-V de ré, ou III-VI7 de do) ; Fm7 B♭7 (II-V de mi♭, ou IVm7-♭VII7
de do, selon l'arrivée) ; Dm7 avant D♭dim7 Cm7 en si♭, qui est B♭/D
(« souvent chiffré à tort Dm7 », dit le livre), comme A♭maj7/C souvent
chiffré Cm7 ; D♭7/G, qui est aussi G7(♭9, ♭5) (Debussy, *La plus que
lente*) ; G13sus4 et Dm9 sur une basse sol ; F♯dim7 avant Gm7, accord de
passage ou D7(♭9) sans fondamentale.

## La fiche et son affichage

Sur le modèle des analyses du chapitre 10, l'analyse d'une grille
produit la structure et le nombre de mesures, le rythme harmonique, la
tonalité de départ et de fin, les plages et les modulations avec leurs mesures, les
emprunts avec leur gamme d'origine, les cadences et les préparations
repérées, et pour chaque accord la ou les gammes à jouer. Elle
s'affiche sur la grille elle-même, annotée comme dans le livre. Trois
supports, dans cet ordre : **le terminal**, banc d'essai qui grandit à
chaque étape ; **une fenêtre Ebitengine**, une fois l'analyse validée,
avec une police de Real Book, du marqueur noir sur fond blanc et les
chiffrages en indices et exposants ; **une page web**, le moteur
compilé en WASM, pour distribuer et faire connaître le travail, dans
la lignée de l'ancien gohareact.

## Les morceaux de référence

### Les fiches du livre

Le livre analyse des morceaux du corpus iReal ; ses fiches, transcrites
à la main dans `charts/ireal/testdata/fiches`, sont l'oracle. Les degrés
en crochets concordent à 83 sur 83, avec ou sans la tonalité que l'app
déclare, qui donne Tune Up en si♭ et le joue sur 32 mesures avec deux
fins.

| Morceau | La fiche | Ce qu'il fixe |
|---|---|---|
| Tune Up | AA', 16 mesures, ré majeur, modulations en do (5 à 8) et en si♭ (9 à 12) | la modulation à chaque phrase, clé de construction du morceau ; 13 sur 13 |
| Black Orpheus | AB, 32 mesures, la mineur, do majeur 6 à 12 | la modulation par deuxième cadence, qui résiste à A7♭9 Dm7 ; 20 sur 20 |
| There Will Never Be Another You | ABAC, 32 mesures, mi♭ majeur, sans modulation | Cm7 tenu une mesure tonicise le VI sans l'installer |
| Tenderly | ABAC, 32 mesures, mi♭ majeur, huit emprunts, plagales | la tonique pressentie, le I emprunté, le pivot diminué, la marche IIm7-V7 lue VI II7 |

### Les autres cas tranchés

| Morceau | Ce qu'il fixe |
|---|---|
| Along Came Betty | le cul entre deux chaises : deux tonalités à un demi-ton qui se chevauchent, des toniques tenues une seule mesure, pas de modulation ; l'analyse montre l'hésitation sans la trancher |
| Giant Steps | trois centres, mais des tonicisations de moins d'une seconde à ce tempo : pas de modulation |
| Black Orpheus, le pont | Dm (une triade) tenu deux mesures : une modulation temporaire, ou une tonicisation appuyée, lecture juste de ce cas limite |
| Autumn Leaves | le relatif qui tonicise d'abord : si♭ passé, sol mineur où se pose la première phrase |
| How Insensitive | une première phrase de quatorze mesures, qui revient sur son accord d'ouverture |
| Just Friends | un morceau qui s'ouvre sur son IV, et se pose sur son I par une plagale |
| Lullaby Of Birdland, All The Things You Are | la dernière tonique entendue, turnaround exclu, fait la tonalité |
| In a Sentimental Mood, Blue Skies | la maison installée avant de partir, et ce qui la distingue d'une simple ouverture |
| Yesterdays, Unforgettable | une tonique traversée par un V seul, un IV où la grille s'arrête sans écrire sa tonique |
| My Way | la quarte et sixte F/C, et le retour par une plagale |
| Sugar, Fly Me To The Moon | un turnaround ne dit rien de la tonalité ; Fly Me part de la mineur et s'arrête en do |
| Softly, Summertime | la tonique mineure écrite m7 |
| 'Round Midnight | la tierce picarde, sur le dernier accord seulement |
| Chega De Saudade | le majeur homonyme installé pour de bon : une modulation, le morceau est en ré |
| Stolen Moments | la marche d'accords parallèles |
| So What | la plage modale, à venir |
| Sonnymoon for Two, Chasin' the Trane, Blues For Alice | le blues reconnu à sa forme |

### Le corpus

`charts/cmd/corpus` compare, sur une ou plusieurs playlists, la
tonalité que l'analyse entend à celle que l'app déclare, en groupant
les écarts par relation (relatif, quinte, quarte, homonyme, autre) avec
des indices pour trancher : le nombre de cadences résolues, la fin sur
la tonique entendue, le blues, la tierce picarde, la grille de jeu
vidéo. Sur les 1678 grilles des deux playlists, les deux tombent
d'accord pour 1300 (77 %), l'app se trompant parfois (Chega De
Saudade, déclaré en ré mineur). Les écarts ne sont pas tous des erreurs de
l'analyse, et leurs familles pointent les questions ouvertes.

## Hors périmètre

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
grille aura sa mélodie. **L'analyse rythmique**, que le livre écarte
aussi.

## Ouvert

- **Les thèmes ambigus entre relatifs**, *Corcovado* (la mineur ou do)
  en tête : la tradition tranche parfois là où l'oreille hésite.
- **L'attente d'un diminué**, ses quatre toniques possibles, pour le
  pivot.
- **Les seuils** de la modulation, en données, et la durée d'une plage
  modale, à régler sur les morceaux de référence puis à l'oreille.
- **La longueur d'une chaîne** de préparations avant qu'elle ne soit
  plus une préparation mais une région.
- **L'ordre exact des lectures** quand plusieurs valent, règle par
  règle, confronté aux fiches.
