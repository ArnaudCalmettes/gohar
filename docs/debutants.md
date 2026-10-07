# Les grands débutants

Cette page débroussaille un cours pour ceux qui ont un clavier mais ne
savent pas encore où sont les notes, ce qu'est une gamme, un accord ou
une tonalité. Rien n'est écrit dans le code : ce sont des pistes et des
décisions de principe, avec les questions qui restent ouvertes.

Le cours pourra devenir un programme à part, ou un tutoriel que chaque
jeu intègre quand il en a besoin (voir « La forme »).

## Le but

Amener le débutant à jouer, et non à savoir. Le cours ne vise pas un
examen de théorie : il vise *Walk with me*, et chaque leçon débloque
aussitôt une façon d'y jouer.

La première destination est modeste : poser les fondamentales de *Satin
Doll*, en rythme, à un tempo lent. *Satin Doll* s'y prête : presque
uniquement des II-V majeurs, autour de C, deux accords par mesure, le
même geste de la main gauche d'une tonalité à la voisine.

## Les principes

- **Apprendre ce qu'on applique tout de suite.** Chaque notion arrive au
  moment où une grille en a besoin, et pas avant. Le débutant n'a pas
  besoin de comprendre l'harmonie pour poser des fondamentales : il lui
  suffit de situer les notes. Il n'a pas besoin des degrés pour jouer
  une basse en deux, fondamentale et quinte : il lui suffit de connaître
  la quinte.
- **Pas de portée.** Le jeu repose sur des grilles : le joueur lit des
  chiffrages, jamais des notes sur une portée. C'est une barrière très
  courante de l'apprentissage du piano dont on s'affranchit.
- **Les deux notations à la fois.** Un francophone apprend do ré mi, et
  les grilles écrivent C D E. S'il n'apprend que la première, il est
  perdu devant une grille ; il apprend donc les deux dès le départ.
- **Jouer délibérément.** Le débutant apprend à viser une note, pas à
  essayer des touches jusqu'à tomber sur la bonne. Une note ratée compte
  comme ratée.
- **La musique ne s'arrête jamais dans les menus.** Tant que le joueur
  navigue, l'orchestre tourne : c'est le `jam` des menus de *Walk with
  me*. Dans un exercice, une fausse note ne coupe rien, et il n'y a pas
  d'écran d'échec.

## Le *ramp up*

Le cours a deux niveaux. Un **chapitre** a un objectif, un savoir qui
débloque une façon de jouer. Il se découpe en **leçons**, assez courtes
pour tenir en une séance, qui y mènent pas à pas. « Connaître son
clavier » est l'objectif du chapitre 1 ; trouver do et fa, puis do ré mi
fa sol, c'est la leçon 1.1.

Les chapitres, et ce que chacun débloque :

| Chapitre | Ce que le joueur sait ensuite | Ce qu'il joue aussitôt |
|---|---|---|
| 1. Le clavier | Trouver une note, en do et en C, avec ses altérations | Le palier 0 : les fondamentales, une par accord |
| 2. Le temps | Compter une mesure à quatre temps, tenir un accord toute sa durée | Le palier 1 : les fondamentales sur le temps, liées |
| 3. Les premiers intervalles | Compter un intervalle en demi-tons et le nommer, de la seconde à la quatorzième ; majeurs, mineurs, justes, augmentés, diminués | La basse en deux : fondamentale et quinte, puis fondamentale et chromatisme |
| 4. La gamme majeure | Ses sept notes, en C d'abord | Des lignes dans la gamme : le début de la walking bass |
| 5. Les degrés | La tierce, la quinte, la septième, la neuvième, comptées dans la gamme | Les triades, les tétrades, puis les voicings : un autre palier, celui de la main droite |

Les accords viennent tard, et c'est voulu. Leur formule demande les
intervalles et la gamme :

- triade = fondamentale + tierce + quinte ;
- tétrade = fondamentale + tierce + quinte + septième, plus les
  extensions.

Parler de tierce, de quinte ou de septième, c'est parler des degrés de
la gamme : il faut au moins la gamme majeure avant de plaquer une
triade.

Chaque chapitre a son fichier, avec ses leçons et l'activité qui le
clôt :

- `debutants/chapitre-1.md`, le clavier ;
- `debutants/chapitre-2.md`, le temps ;
- `debutants/chapitre-3.md`, les premiers intervalles.

## La gamme dans les douze tons

La gamme majeure ne s'apprend pas en une leçon qu'on termine. Elle
devient un fil continu, qui court à côté du reste du cursus : C d'abord,
puis les tonalités une à une, de la moins altérée à la plus altérée. La
difficulté d'une tonalité, c'est son nombre d'altérations.

**L'ordre alterne les dièses et les bémols**, une altération de plus à
chaque paire : C, puis G et F, puis D et B♭, puis A et E♭, et ainsi de
suite jusqu'aux tonalités à six altérations. Le joueur se construit ainsi
le cycle des quintes, des deux côtés à la fois, et le jeu le lui montre
à mesure qu'il le remplit.

**La règle s'explique dès le départ**, parce qu'elle fait d'une liste de
douze gammes une seule idée :

- en montant d'une quinte, on crée une sensible : on hausse le 7e degré
  de la nouvelle gamme. De C à G, le 7e degré de G est haussé : F♯ ;
- en descendant d'une quinte, on détruit la sensible de la gamme
  précédente : on abaisse cette note, qui tombe sur le 4e degré de la
  nouvelle gamme. De C à F, le 7e degré de C, abaissé, devient le 4e
  degré de F : B♭.

**Posséder une gamme dans une tonalité**, c'est savoir la jouer :

- de façon régulière, sans se tromper, du premier coup ;
- pas forcément vite ;
- dans les deux sens, en montant et en descendant ;
- d'une main ou des deux : le MIDI ne dit pas laquelle, et le jeu ne le
  demande pas.

Le dex marque alors la gamme produite dans cette tonalité.

Le dex sait déjà la suivre : la production d'une notion se compte par
tonique (voir `dex.md`). Ce fil dit aussi, à chaque instant, **dans
quelles tonalités le jeu peut transposer** une grille ou un exercice :
celles dont le joueur possède la gamme. Une tonalité tiédit si le joueur
ne la pratique plus, et le jeu la lui repropose : c'est le chantier de
la répétition espacée du dex (`Cooling`).

## Le microdex

Le dex est une collection d'objets sonores, qu'on peut reconnaître,
produire et dont un jeu constate la présence (voir `dex.md`). Une grande
partie de ce qu'apprend le débutant n'en est pas :

| Notion | Dans le dex ? |
|---|---|
| L'emplacement d'une note | Non : on la produit, mais on ne la reconnaît pas à l'oreille, sauf avec l'oreille absolue |
| Le temps, la pulsation | Non : ce n'est pas de l'harmonie |
| L'intervalle | Oui, déjà, en notion élémentaire |
| La gamme majeure | Oui : le chantier `KindSystem`, acté |
| La tonique, « la maison » | Un exercice plutôt qu'une notion, comme le degré |

D'où un **microdex**, le registre des acquis du débutant, chacun
débloquant une façon de jouer :

- **les notes**, une entrée par note dans chaque notation : « sait
  trouver ré », « sait trouver D » ;
- **le temps**, par tempo : « tient le temps à 80 » ;
- **les intervalles** et **la gamme majeure**, qui vivent déjà dans le
  dex et s'y lisent.

**Le microdex est à part du dex**, sans appartenir à un jeu : c'est un
dex spécial, pour les notions de base et les prérequis, que tous les
jeux partagent. Un jeu y lit ce que son joueur sait déjà, pour lui
proposer la leçon qui lui manque ; une leçon y marque ce qu'elle a
appris. Il n'a ni fraîcheur ni dévoilement : un prérequis est acquis ou
ne l'est pas encore. La gamme dans les douze tons n'en fait pas partie :
elle appartient au dex.

## Les deux notations

Le clavier à l'écran affiche « do / C » sur ses touches. Les exercices
demandent tantôt l'une, tantôt l'autre ; une grille, elle, n'écrit que
des lettres. Le joueur fait le lien par la pratique. Les étiquettes du
clavier peuvent s'effacer à mesure que le microdex montre qu'il l'a
fait.

## Les grilles

Le cours se joue d'emblée sur de vraies grilles, *Satin Doll* et le
blues (voir `debutants/chapitre-1.md`). Une petite grille, de quatre ou
huit mesures, ne s'écrit que quand une leçon en a besoin, libre de
droits comme toutes les grilles de gohar : celle de la leçon 2.3, quatre
accords par mesure, en est une. Le II-V-I arrivera quand le joueur
saura ce que sont un II, un V et un I.

## Rendre le cours amusant

- **Un succès immédiat** : « joue seulement les touches noires », et
  tout sonne avec l'orchestre. C'est une gamme pentatonique, qu'on
  nommera plus tard ; elle donne le droit de jouer avant de savoir.
- **Le bonhomme comme partenaire** : il joue une phrase, le joueur la
  répète ; il claque des doigts quand ça tourne, comme dans *Walk with
  me*.
- **Des mélodies que tout le monde connaît**, dans le domaine public
  (*Au clair de la lune*, l'*Ode à la joie*), pour placer les notes.
  Des fanfares originales dans l'esprit des jeux vidéo, jamais de vraies
  musiques de jeu, qui ne sont pas libres.
- **La découverte avant le nom** : le joueur trouve un son, le jeu le
  nomme ensuite. C'est le mécanisme d'`Overheard` dans le dex, appliqué
  au débutant.
- **Des séances courtes**, de trois à cinq minutes par leçon.
- ***Au clair de la lune* en neuvièmes mineures**, un gag éprouvé du
  chapitre des intervalles. Le bonhomme joue d'abord la première phrase
  en octaves, do à gauche et do à droite : ça sonne plein, rassurant.
  Puis, arrivé à la neuvième mineure, le demi-ton écarté d'une octave,
  il la rejoue sur cet intervalle : en do à gauche, en do♯ à droite. On
  entend d'un coup ce qu'est une dissonance, et tout le monde rit, les
  adultes aussi. Le joueur peut ensuite tenir une des deux mains,
  pendant que le bonhomme joue l'autre. La chanson est dans le domaine
  public.
- **La sensible par l'attente.** *Shave and a Haircut*, la formule de
  fin que tout le monde connaît (« tagada tsoin tsoin » en français),
  traditionnelle et libre : sol, sol-sol-la, sol, … si, do. Le bonhomme
  s'arrête sur le si et ne joue pas le do : le joueur sent la tension,
  la note qui manque. C'est lui qui joue le do, et « ah ! ça va mieux ».
  Le si est la sensible : la note qui tire vers la tonique, un demi-ton
  au-dessous. Le joueur en a déjà senti l'effet à la leçon 1.2, à la
  fin de la gamme ; ce gag-ci la nomme et la fait **retenir**. Il a sa
  place au chapitre 4, juste avant la règle du cycle des quintes (voir
  « La gamme dans les douze tons »), qui parle de créer et de détruire
  une sensible. À ce stade, le joueur sent déjà que les degrés de la
  gamme jouent des rôles différents : s'arrêter sur le 5e degré pose une
  question, sur le 3e c'est une virgule, sur le 7e ça perturbe. Il le
  sent sans qu'on le lui ait expliqué ; les degrés ne s'expliquent
  qu'au chapitre 5.
- **Un easter egg, pour qui fait exprès de rater.** Seulement sur les
  demandes où l'on ne peut pas se tromper de bonne foi, comme « appuie
  sur les deux touches noires d'une paire » : là, dix casseroles d'affilée,
  c'est forcément exprès. Le bonhomme s'agace par paliers : un « Hé… »
  d'abord, puis « Tu le fais exprès ? », et à la dixième : « Bon, j'en
  ai marre ! Je peux pas travailler dans ces conditions. » Il part à
  pied vers la droite, sort de l'écran, et un GAME OVER s'affiche. Une
  touche, et il revient (« Bon, d'accord. Mais c'est la dernière
  fois ! ») : la leçon reprend où elle en était, rien n'est perdu. Une
  étape le permet ou non ; le compte repart de zéro à la première bonne
  note.

  Le joueur qui a ri une fois recommence. Le gag change donc à chaque
  fois, de pire en pire, pour qu'il finisse par trouver plus long de
  désobéir que d'écouter :
  1. la première fois, la sortie et le GAME OVER ;
  2. la deuxième, il s'assoit dans l'herbe, les mains derrière lui :
     « Bon, je me pose là et j'attends que tu veuilles bien faire ce
     que je t'ai demandé. » Il se relève à la bonne réponse ;
  3. la troisième, il se met en lotus, au son d'un bol chantant : « Je
     suis calme. Je suis très calme. » Il ne réagit plus à rien, et
     n'en sort qu'à la bonne réponse : « Ah. Merci. » Toutes les fois
     suivantes, le lotus encore.

  Le compte des gags vaut pour toute la session : sortir de la leçon
  et y revenir ne le remet pas à zéro.

  Dans le jeu (`games/walk/gags.go`), une étape « Demander » le
  permet par son drapeau `Teasing` : à la leçon 1.1, la paire et le
  groupe de trois. « Hé… » vient à la quatrième casserole d'affilée,
  « Tu le fais exprès ? » à la septième, la sortie à la dixième. Le
  bonhomme sort à droite de profil, en marchant ; GAME OVER s'affiche,
  sur un petit air de piano dans les aigus : quatre tritons qui
  descendent d'un demi-ton, le dernier en trémolo. Une touche, de
  l'ordinateur ou du clavier MIDI, et il revient de la droite, du même
  pas, tourné vers la gauche ; une touche encore, et l'étape reprend.
  Comme ses autres répliques, ses mots attendent que la touche soit
  relâchée. Assis puis en lotus, il reste à sa place et la leçon
  continue : les touches comptent, les casseroles sonnent. D'une pose à
  l'autre, une image tenue un instant : de debout, accroupi, une main
  derrière lui au sol ; d'une pose assise, vers l'autre ou pour se
  relever, une main au sol devant lui. Relevé, il claque des doigts. Le bol
  chantant est synthétisé (`games/walk/band/sounds/bowl.go`).

## Plus loin

La tonalité, au sens de l'harmonie, n'a pas encore de place dans le
*ramp up* : c'est le point où le cours rejoint le cursus de *Walk with
me*, et l'analyse, qui la détecte déjà. Sa place se décidera une fois
les cinq premiers chapitres construits, et éprouvés avec de vrais
débutants.

## La forme

Un paquet de leçons, chacune une scène (`games/scene`) avec ses phrases
(`games/lang`). Deux usages :

- **intégré** : un jeu pousse une leçon quand il manque un acquis à son
  joueur ; *Walk with me* propose « Lire un chiffrage, trois minutes »
  avant le premier run ;
- **autonome** : un programme qui enchaîne les chapitres dans l'ordre
  du *ramp up*, et leurs leçons.

## Le bonhomme professeur

Les leçons, c'est le bonhomme de *Walk with me* qui les donne, dans ses
bulles, en tutoyant, comme un copain musicien : jamais un juge.

**Apprendre en faisant.** Une leçon est une suite d'étapes courtes, qui
finissent par un geste au clavier plutôt que par « appuie sur Entrée ».
Une règle la tient : jamais plus de deux bulles sans que le joueur joue
quelque chose.

| Étape | Ce que fait le bonhomme | Ce qui fait avancer |
|---|---|---|
| Dire | Une bulle d'explication | Le temps de lecture, ou une touche |
| Montrer | Il allume des touches sur le clavier à l'écran, et le dit | Sa bulle lue ; les touches restent allumées pour l'étape suivante, qu'il prépare |
| Jouer | Il joue une note ou une phrase, qu'on entend et qu'on voit | La fin de la phrase |
| Demander | « Trouve le do », « joue la quinte de F » | La bonne réponse ; un raté, c'est la casserole et la bonne touche montrée |
| Répéter après lui | Il joue, le joueur reproduit ; ou il la dit seulement, « C, E, C » | La phrase rejouée juste |
| Faire deviner | Il demande une note qu'il n'a pas montrée, « Et mi♯ ? » | La bonne réponse ; un raté, c'est la casserole et un rappel, et au bout de trois, la réponse donnée |
| Écrire | Il écrit des accords, côte à côte ou en ligne de grille, et dit ce que c'est | Sa bulle lue ; ce qu'il a écrit reste |
| Jouer une ligne | Une ligne de grille, jouée à la basse par lui, ou par le joueur, mesure par mesure | La fin de sa phrase ; ou la ligne jouée jusqu'au bout |
| Laisser en suspens | Il dit une phrase ; sa dernière note reste en l'air, la note qui la résout clignote | La note qui la résout |
| Jouer ensemble | L'orchestre tourne, il annonce ce qu'on fait | Un nombre de réussites d'affilée |

« Demander » est le palier 0 ; « Jouer ensemble », une mini-partie sur
une petite grille. R fait réentendre ou relire la dernière étape ; une
touche saute ce que le joueur sait déjà.

**L'activité de la leçon.** Une leçon se donne une fois ; elle ouvre une
activité, sur laquelle le joueur s'entraîne aussi longtemps et aussi
souvent qu'il veut, avant de passer à la leçon suivante. L'activité de
la leçon 1.1 (`debutants/chapitre-1.md`) en donne le modèle. Dans le
jeu, une leçon faite qui a une activité ouvre un dernier choix :
« S'entraîner », ou « Revoir la leçon ».

**Le passage à la leçon suivante**, c'est le joueur qui le décide, quand
il se sent prêt. Une leçon donnée est marquée faite, d'une coche verte ;
la suivante apparaît, prête à être lancée quand il le voudra. Une leçon
faite se rejoue à volonté, et son activité reste ouverte.

**Chanter est un conseil, jamais une obligation.** Le jeu n'écoute pas
et n'écoutera pas : pas de micro. Certains joueurs sont mal à l'aise à
l'idée de chanter, d'autres jouent au casque avec du monde autour ; ce
qu'on se permet dans une vraie leçon de piano serait intrusif dans un
jeu. Le bonhomme le conseille, et le reconseille souvent, d'une leçon à
l'autre, sans jamais le vérifier.

**L'écran** ressemble à celui de la partie : le bonhomme à la même
place, à gauche, et le clavier en bas. La grille n'apparaît que quand
la leçon en a besoin. Le clavier peut se réduire à deux ou trois
octaves, pour être plus gros et plus lisible.

**Le clavier du joueur.** Le MIDI ne dit pas combien de touches compte
un clavier. Le plus simple est de le demander en jouant, et c'est une
bonne première étape : « joue la touche la plus grave de ton clavier,
puis la plus aiguë ». Le jeu garde l'étendue dans les réglages, et
l'affine s'il entend une note au-delà. Le nom du clavier, qui contient
parfois son nombre de touches (« Keystation 49 »), ne s'y fie pas : rien
ne l'oblige.

**Du texte**, pour l'instant : enregistrer des voix est hors de portée,
et deux langues en doublent le coût. Une synthèse vocale, plus tard,
peut-être ; le navigateur en offre une.

**Le code.** Les leçons s'écrivent d'abord en Go : leurs conditions
(« un do », « la quinte de F ») sont du code, et un format de données
n'a de sens qu'une fois les types d'étapes stabilisés. Elles vivent
d'emblée dans un paquet à elles, `games/walk/lessons`, tant que seul
*Walk with me* s'en sert. Ce paquet ne connaît pas *Walk with
me* : il décrit les étapes et juge les notes, et le jeu lui prête une
scène (une interface : dire dans une bulle, allumer des touches, jouer
une phrase, lancer l'orchestre). Les leçons restent ainsi extractibles
vers les autres jeux, comme prévu dans « La forme », et se testent sans
écran.

Le paquet a pris forme ainsi :

- **la scène** (`Stage`) : dire une phrase, par son identifiant, que le
  jeu traduit ; allumer des touches, désignées par leur note, à toutes
  les octaves ; faire jouer une phrase au bonhomme ; répondre à un raté,
  par la casserole et la bonne touche montrée. L'orchestre viendra avec
  « Jouer ensemble » ;
- **les cibles** (`Target`) : une note à n'importe quelle octave, jugée
  sur la dernière touche enfoncée, pour qu'un do tenu sous la main
  gauche ne compte pas comme une faute ; un groupe de touches noires
  enfoncées ensemble, la paire ou le trio, qui reste « en cours » tant
  que toutes ses touches ne sont pas là, et compte comme raté si on les
  relâche avant ;
- **les étapes** : Dire, Montrer, Jouer, Demander, Répéter après lui,
  Faire deviner, Laisser en suspens, Écrire et Jouer une ligne. Répéter se joue sans faute : une fausse note
  fait sonner la casserole, le bonhomme rejoue la phrase, et le joueur
  la reprend du début. Une gamme jouée de si à do contient « do sol
  do », elle ne le répète pas. Dite plutôt que jouée, une fausse note
  allume toutes les touches de la phrase, jusqu'à ce qu'elle soit
  rejouée en entier ;
- **le déroulé** (`Runner`) : le jeu lui transmet les touches, la fin
  des phrases du bonhomme et les bulles lues, et il passe à l'étape
  suivante quand celle en cours est finie. R rejoue l'étape en cours.

## Dans le jeu

**Le menu.** Pour l'instant, le cours vit dans *Walk with me*, sous
« Apprendre » au menu du titre. L'écran liste les chapitres ouverts,
puis les leçons du chapitre choisi : celles qui sont faites, d'une
coche verte, et la suivante ; les autres restent cachées. Le catalogue
des chapitres et la règle qui ouvre les leçons sont dans
`games/walk/lessons`, la progression du joueur dans le fichier
`course.json` de ses réglages.

**L'écran d'une leçon** (`games/walk/lesson.go`). Le bonhomme se tient
à gauche, face au joueur, les deux mains dans les poches (voir
« Le bonhomme » dans `walk.md`). Sa bulle court sur plusieurs lignes à
droite, et sa queue part du point de la bulle le plus proche de sa
tête, qu'elle vise. Le clavier de la partie, en bas, est réduit à trois
octaves, de C3 à C6, chaque touche blanche nommée « do » au-dessus de
« C » ; les touches montrées s'éclairent en bleu pâle. Tout le clavier
sonne en piano, comme partout hors d'une grille : la contrebasse sous
G3 n'a pas de sens pour une leçon, et un do qui devient contrebasse une
octave plus bas déroute. L'orchestre se tait : le chapitre 1 se
joue sans tempo, et les phrases du bonhomme s'entendent seules.

**Les touches.** Quand une bulle attend d'être lue, n'importe quelle
touche la tourne, de l'ordinateur ou du clavier MIDI ; celle du clavier
ne sonne pas alors, et ne s'allume pas, pour qu'on ne la prenne pas
pour une réponse. R rejoue l'étape, Échap revient au cours sans marquer
la leçon faite. La ligne d'état dit ce qu'on attend : une touche pour
continuer, ou des notes à jouer.

**Les réactions.** La casserole, la cloche du kit, sonne avec la fausse
note, tout de suite : c'est le son de la note qui n'était pas attendue.
Une bonne réponse a son signe à elle : le bonhomme sort une main de sa
poche et claque des doigts, avec un mot d'approbation, jamais deux fois
de suite le même. Lui ne réagit qu'une fois la touche relâchée, ou un
temps après qu'elle a été enfoncée, qu'il s'agisse de saluer une bonne
réponse, de s'agacer ou de dire un mot après une fausse note :
répondre à une touche encore tenue coupe la parole au joueur. Après une fausse note, il remontre la consigne, les
bonnes touches ou la phrase rejouée, sur la mesure suivante, comme si
le joueur jouait en rythme : une note par temps, quatre temps par
mesure, la phrase commencée sur le 1. Un motif de trois notes raté sur
sa troisième tombe sur le 3 ; le 4 laisse entendre la casserole, et le
bonhomme reprend sur le 1, comme le claquement d'une bonne réponse
tombe sur le 4. Une étape
que le joueur termine au clavier ne passe pas aussitôt la main non
plus : la suivante attend que ses touches soient relâchées, puis le
temps d'une noire, pour ne pas le presser.

**L'avancement.** Les quatre leçons du chapitre 1 et leurs activités
sont écrites, ainsi que la leçon qui le clôt ; restent son activité, le
palier 0, et les chapitres suivants. Une fois une leçon dans le jeu,
sa doc la résume, étape par étape, sans recopier ses bulles : le texte
exact vit dans les fichiers de phrases, où on le retouche.

## Les questions ouvertes

- La place du microdex dans le dépôt : un paquet à côté du dex, dans
  son module, ou un module à lui ; et son nom.
- Ce qui clôt les chapitres suivants. Le chapitre 1 se clôt par une
  activité qui ramène au jeu, le palier 0 ; reste à voir si chaque
  chapitre se clôt de même par le palier qu'il débloque, et comment le
  montrer au joueur.
- Le tempo du palier 1 au départ, et si le jeu l'accélère de lui-même,
  comme le conseil du bilan de *Walk with me*.
- Le point où le cours s'arrête et où le cursus de *Walk with me* prend
  le relais (voir « Plus loin »).
