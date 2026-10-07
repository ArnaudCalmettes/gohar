# Chapitre 1 : le clavier

Le premier chapitre du cours des grands débutants (voir `../debutants.md`).
Son objectif : connaître son clavier, trouver une note en do et en C,
avec ses altérations. Les leçons suivent le modèle du bonhomme professeur
décrit dans `../debutants.md`.

Il compte quatre leçons :

| Leçon | Ce qu'elle apporte |
|---|---|
| 1.1 | Les groupes de touches noires ; do et fa ; do ré mi fa sol |
| 1.2 | La gamme de do majeur : les sept notes blanches, avec la et si ; mi et si à droite des groupes, comme do et fa à gauche |
| 1.3 | Les dièses : une note qu'on fait monter d'un cran ; mi♯ et si♯, les premières notes enharmoniques |
| 1.4 | Les bémols : une note qu'on fait descendre d'un cran ; deux noms par touche noire ; fa♭ et do♭ ; l'échelle chromatique |

Il se clôt par une leçon qui ramène au jeu, « Premières grilles, sans
tempo », dont l'activité est une grille jouée hors tempo, au palier 0
(voir plus bas). Ce n'est pas tout à fait le mode
déchiffrage : il faut viser juste, et donner la bonne note du premier
coup. La grille est une vraie grille du jeu : une section A entière de
*Satin Doll*, ou le blues, ou les deux. Leurs fondamentales demandent ce
que le chapitre vient d'apprendre : A♭ et D♭ dans *Satin Doll*, B♭ dans
le blues en F. Le blues transposé en C en donne une version sans
altérations : C, F, D, G et A, que des touches blanches. Le jeu le
transpose à la volée, sans grille à part.

Le chapitre se joue sans batterie : le rythme n'arrive qu'au chapitre 2.
S'il faut en glisser un avant-goût, l'usage le dira.

## Les leçons

### 1.1 Premiers pas

La leçon 1.1 mobilise l'oreille le plus tôt possible : le joueur chante
ce qu'il joue. Elle est écrite dans le jeu : ses étapes dans
`games/walk/lessons/chapter1.go`, ses bulles dans `games/walk/locales`.
Ce qui suit la résume.

1. Le bonhomme se présente : avant de jouer ensemble, on fait
   connaissance avec le clavier.
2. Il montre les paires de touches noires, partout sur le clavier, et
   fait appuyer sur les deux touches d'une paire, ensemble ; de même
   pour les groupes de trois. Le joueur prend en main la géographie du
   clavier, deux et trois doigts posés, avant le premier nom de note.
   Ce sont les deux demandes de la leçon où l'on ne se trompe pas de
   bonne foi : l'easter egg y veille.
3. Do et fa, ensemble et dès le départ : do à gauche d'une paire, fa à
   gauche d'un trio, chacun en do ré mi et en lettres (C, F). Des
   débutants les confondent au début ; les apprendre l'un contre
   l'autre, par leur groupe, coupe court à la confusion. Le bonhomme
   les fait trouver, mélangés, dans les deux notations.
4. Do ré mi fa sol, à partir du do, joués après lui en les chantant,
   puis en redescendant, sol fa mi ré do.
   Cinq notes, pas sept : elles tiennent sous une main, et C D E F G se
   suivent dans l'alphabet, ce qui aide à retenir les lettres.
5. Les motifs, do-ré-do, do-mi-do, do-fa-do, do-sol-do, répétés après
   lui sans faute, d'une main ou des deux. Avec deux mains, le do reste
   sous la gauche, et la droite n'a qu'à se décaler de touche blanche
   en touche blanche, d'un motif au suivant. Chanter les motifs fait
   retenir, en même temps que les notes, la seconde, la tierce, la
   quarte et la quinte, que le joueur retrouvera nommées au chapitre 3 ;
   le conseil, donné à l'étape 4, reviendra dans l'activité.

**L'activité de la leçon 1.1** reprend ses motifs, C-D-C, C-E-C, C-F-C,
C-G-C, sans tempo, de deux façons :

- **montrés** : le bonhomme joue le motif, le joueur le répète ;
- **dits en lettres** : « C, E, C », sans rien jouer ni allumer.

La leçon a nommé les notes dans les deux notations ; passé ce moment de
découverte, le jeu ne demande plus que les lettres. Le clavier à
l'écran garde ses deux étiquettes, « do » au-dessus de « C » : le
joueur qui pense en do ré mi s'y retrouve, et apprend les lettres en
les lisant à côté.

Chaque façon passe d'abord les motifs dans l'ordre (D, puis E, puis F,
puis G), puis mélangés. Viennent ensuite les deux façons mélangées
entre elles, avec la suite entière de la leçon, C D E F G, et sa
descente, G F E D C, chacune une fois montrée et une fois dite. Le conseil de chanter revient deux fois. Le
tirage change à chaque partie.

Comme dans la leçon, un motif se rejoue sans faute. Montré, un raté le
fait rejouer par le bonhomme. Dit, un raté allume toutes les touches
du motif, jusqu'à ce qu'il soit rejoué depuis le début.

Elle est écrite dans le jeu (`games/walk/lessons/activities.go`) ; on
l'ouvre d'une leçon faite, par « S'entraîner », sous la leçon.

### 1.2 La gamme de do majeur

La leçon 1.2 complète la gamme, et mobilise d'abord la mémoire
auditive et lexicale du joueur : c'est la suite chantée qui situe la et
si, plutôt que leur place parmi les touches noires. Elle est écrite
dans le jeu, comme la 1.1 ; ce qui suit la résume.

1. Le bonhomme reprend où la leçon 1.1 s'est arrêtée : une dernière
   révision de C D E F G, demandés en lettres, avant de les dire
   acquis.
2. La gamme entière, de do au do du dessus, en montant puis en
   descendant, répétée après lui en chantant. La suite s'installe, dans
   l'oreille et dans les mots. La consigne précise, entre parenthèses,
   qu'elle ne se joue pas en rythme : les doigtés viendront plus tard,
   et le passage du pouce est hors de portée d'un grand débutant. Un
   seul doigt suffit.
3. Les mêmes notes en lettres, C D E F G A B, avec la mise en garde :
   après G, l'alphabet repart de A. Puis il les demande, dites, sans
   les jouer.

   Un gag en guise de chute, la sensible avant son nom. Sur le dernier
   B, le jeu garde la note : le son reste quand le joueur lâche la
   touche, sans que le clavier à l'écran la montre tenue. Le bonhomme
   s'impatiente, il manque quelque chose ; le C au-dessus clignote.
   Le joueur le joue, la note tenue s'arrête, et le bonhomme, soulagé,
   remercie. Le mot de sensible ne vient qu'au chapitre 4 ; le joueur
   en a déjà senti l'effet.
4. Une astuce pour deux notes : do et fa sont à gauche des groupes de
   deux et de trois touches noires ; mi et si sont à droite de ces
   mêmes groupes. Les quatre touches s'allument, le temps de la bulle,
   sans exercice : s'attarder sur do fa mi si, qui ne forment pas une
   phrase, n'aurait rien de musical.
5. Les motifs, en miroir de ceux de 1.1 : do-si-do, do-la-do,
   do-sol-do, répétés après lui en chantant. Cette fois, le do est à
   la main droite et l'autre note à la main gauche, qui descend de
   touche blanche en touche blanche. Avec le sol, ces motifs couvrent
   un tétracorde complet, sol la si do.
6. Pour finir, les sept notes à trouver, une fois chacune, dans un
   ordre tiré à chaque leçon, chacune demandée en lettres ou en do ré
   mi. C'est là que mi et si, comme les autres, se retrouvent sous les
   doigts.

Dans le jeu, la note tenue est une étape à part, « Laisser en
suspens » (`Hang`) : une suite demandée comme « Répéter après lui »,
dont la dernière note reste en l'air jusqu'à la note qui la résout.

**L'activité de la leçon 1.2** suit le modèle de celle de 1.1. Elle met
ensemble tous les motifs de 1.1 et de 1.2, et sa séquence longue est la
gamme entière, de C au C du dessus. Dit en lettres, C-G-C est le même
dans les deux leçons, en montant comme en descendant.

### 1.3 Les dièses

La leçon 1.3 nomme les touches noires par leurs dièses. Elle donne au
joueur le sens de l'orientation, avec le vocabulaire qui va avec : on
**monte** vers la droite du clavier, on descend vers la gauche. Elle
est écrite dans le jeu ; ce qui suit la résume.

1. Les touches noires ont des noms, elles aussi : un dièse, c'est une
   note qu'on a fait monter d'un cran, vers la droite. Puis le signe,
   ♯, et les deux notations : do♯ et C♯.
2. Do puis do♯, fa puis fa♯, répétés après lui : la note monte d'un
   cran. Le mot de demi-ton attend le chapitre 3.
3. Les cinq touches noires à trouver par leur dièse, en do ré mi puis
   en lettres, dans un ordre tiré à chaque leçon.
4. Deux notes à deviner : « Et mi♯ ? », puis « Et si♯ ? » Le joueur a
   trois chances, sans qu'on les lui compte : chaque raté sonne la
   casserole, sans rien montrer, avec un rappel, un dièse monte la
   note d'un cran. Au troisième, le bonhomme donne la réponse et montre
   la touche, à jouer : mi♯ est sur le fa, si♯ sur le do, faute de
   touche noire entre les deux. Trouvées ou données, il nomme la chose :
   deux noms pour une même touche, des notes **enharmoniques**. L'astuce
   de 1.2 l'annonçait : mi et si sont collés à droite d'un groupe de
   touches noires, sans touche noire après eux.
5. La séquence longue, dite en lettres : C C♯ D D♯ E F F♯ G G♯ A A♯ B C.
   Elle monte touche par touche et prépare la basse en chromatisme du
   chapitre 3.

Dans le jeu, la devinette est une étape à part, « Faire deviner »
(`Guess`) : une note demandée sans avoir été montrée, avec un nombre de
chances avant la réponse.

**L'activité de la leçon 1.3** ne reprend pas les motifs de 1.1 et 1.2.
Elle demande des séquences de trois notes, dites en lettres : une
touche blanche, sa voisine noire, la touche blanche suivante, en
montant (C C♯ D) comme en descendant (D C♯ C), pour chaque touche
noire. Un premier tour les passe toutes, dans un ordre tiré ; un second
y glisse les pièges, mi♯ et si♯ : E E♯ F♯, B B♯ C♯, C♯ D♯ E♯ et
G♯ A♯ B♯.

### 1.4 Les bémols

La leçon 1.4 est le miroir de 1.3 : les bémols, en descendant. Elle est
écrite dans le jeu ; ce qui suit la résume.

1. Un bémol, c'est une note qu'on a fait descendre d'un cran, vers la
   gauche. Puis le signe, ♭, et les deux notations : ré♭ et D♭.
2. Ré puis ré♭, si puis si♭, répétés après lui.
3. Les cinq touches noires à trouver par leur bémol, en do ré mi puis
   en lettres, dans un ordre tiré à chaque leçon.
4. Do♯, puis ré♭, à trouver l'un après l'autre : c'est la même touche.
   Chaque touche noire a deux noms, et le bonhomme le prend sur un ton
   léger : les musiciens aiment bien compliquer les choses. Do♯ et ré♭
   sont enharmoniques, comme mi♯ et fa.
5. Deux notes à deviner, comme à la leçon 1.3 : fa♭, puis do♭, trois
   chances chacune, la casserole à chaque raté avec le rappel, un bémol
   descend la note d'un cran. Au bout de trois essais, la réponse : fa♭
   est sur le mi, do♭ sur le si.
6. La séquence longue, qui descend cette fois, dite en lettres :
   C B B♭ A A♭ G G♭ F E E♭ D D♭ C.
7. En conclusion, le nom de la chose : les douze touches jouées une à
   une, en montant à la leçon 1.3 et en descendant ici, forment
   l'**échelle chromatique**.

**L'activité de la leçon 1.4** suit le modèle de celle de 1.3, en
bémols : une touche blanche, sa voisine noire nommée par son bémol, la
touche blanche suivante, en montant (C D♭ D) comme en descendant
(D D♭ C), dites en lettres. Le second tour y glisse fa♭ et do♭ :
F F♭ E♭, C C♭ B♭, G♭ F♭ E♭ et D♭ C♭ B♭.

### Premières grilles, sans tempo

La leçon qui clôt le chapitre explique ce qu'est une grille, et ce
qu'on attend du joueur. Elle est écrite dans le jeu ; ce qui suit la
résume.

1. Un accord, tel qu'une grille l'écrit : E♭7, à la main des grilles.
2. Sa fondamentale : E♭ est celle de E♭7, comme de E♭6, de E♭m ou de
   E♭maj7(♯11), écrits côte à côte. C'est la note du début, avec son
   altération ; ce que veut dire le reste du symbole, on le découvrira
   plus loin.
3. Quelques fondamentales à trouver, l'accord écrit : E♭7, F7, B♭m7,
   A♭maj7. Les altérations viennent d'être apprises ; ces accords
   préparent le blues et *Satin Doll*.
4. La grille : la musique derrière une chanson, ses accords écrits
   mesure par mesure. Une ligne du blues en C s'affiche, comme dans le
   jeu : | C7 | F7 | C7 | C7 |.
5. Ce qu'on attend du joueur : la basse, la fondamentale de chaque
   accord, une note par mesure. Le bonhomme joue la ligne à la
   contrebasse, chaque fondamentale tenue toute la mesure, quatre
   temps, la mesure qu'il joue grisée. Il claque des doigts sur 2 et 4 :
   on entend que c'est lent exprès.
6. Un conseil, en « si » faute de pouvoir connaître le clavier : un
   clavier de deux octaves descend de deux octaves en appuyant deux
   fois sur « Oct − », et la basse sonne alors là où elle doit.
7. Au joueur, sur la même ligne : la contrebasse sous ses doigts, la
   grille qui l'attend, mesure par mesure, et la casserole à une fausse
   fondamentale, la ligne restant sur sa mesure. C'est le palier 0 en
   petit.

Dans le jeu, l'accord écrit est une étape à part, « Écrire » (`Write`),
et la ligne jouée une autre, « Jouer une ligne » (`PlayLine`), par le
bonhomme ou par le joueur. Le clavier à l'écran descend d'une octave
quand le joueur joue la basse.

## Le palier 0

C'est l'activité qui clôt le chapitre. Un palier plus simple que le
premier de *Walk with me* : le joueur pose la fondamentale de chaque
accord, sans exigence de temps.

Il ne tolère pas l'erreur pour autant. Le mode déchiffrage de *Walk with
me* attend patiemment la bonne note et ne compte pas les autres ; le
palier 0, lui, **attend sa note, à l'octave près, et une note ratée
compte comme une cible manquée**. C'est ce qui apprend à viser.

Après un raté, la grille n'avance pas : elle fait sonner un bruit de
casserole et montre la bonne touche sur le clavier à l'écran. Chaque
nouvelle fausse note sonne de même, jusqu'à ce que la bonne soit jouée ;
alors seulement la grille passe à l'accord suivant. La casserole est un
gag plus qu'une sanction : le raté s'entend, il fait sourire, et la
correction vient tout de suite, la bonne touche sous les yeux.
