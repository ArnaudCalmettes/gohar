# Chapitre 1 : le clavier

Le premier chapitre du cours des grands débutants (voir `../debutants.md`).
Son objectif : connaître son clavier, trouver une note en do et en C,
avec ses altérations. Les leçons suivent le modèle du bonhomme professeur
décrit dans `../debutants.md`.

Il compte quatre leçons :

| Leçon | Ce qu'elle apporte |
|---|---|
| 1.1 | Les groupes de touches noires ; do et fa ; do ré mi fa sol |
| 1.2 | La et si : les sept notes blanches ; mi et si à droite des groupes, comme do et fa à gauche |
| 1.3 | Les dièses : une note qu'on fait monter d'un cran ; mi♯ et si♯, les premières notes enharmoniques |
| 1.4 | Les bémols : une note qu'on fait descendre d'un cran ; deux noms par touche noire ; fa♭ et do♭ ; l'échelle chromatique |

Il se clôt par une activité qui ramène au jeu : une grille jouée hors
tempo, au palier 0 (voir plus bas). Ce n'est pas tout à fait le mode
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
`games/walk/lessons/scripts.go`, ses bulles dans `games/walk/locales`.
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
4. Do ré mi fa sol, à partir du do, joués après lui en les chantant.
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
entre elles, avec la suite entière de la leçon, C D E F G, une fois
montrée et une fois dite. Le conseil de chanter revient deux fois. Le
tirage change à chaque partie.

Comme dans la leçon, un motif se rejoue sans faute. Montré, un raté le
fait rejouer par le bonhomme. Dit, un raté allume toutes les touches
du motif, jusqu'à ce qu'il soit rejoué depuis le début.

Elle est écrite dans le jeu (`games/walk/lessons/activities.go`) ; on
l'ouvre d'une leçon faite, par « S'entraîner ».

### 1.2 La et si

**La leçon 1.2** complète la gamme, et mobilise d'abord la mémoire
auditive et lexicale du joueur : c'est la suite chantée qui situe la et
si, plutôt que leur place parmi les touches noires.

1. *Dire* : « Do ré mi fa sol, tu les as. On continue. »
2. *Répéter après lui* : do ré mi fa sol la si do, en montant puis en
   descendant, en chantant. La suite entière s'installe, dans l'oreille
   et dans les mots.
3. *Dire*, puis *demander* : les mêmes notes en lettres, C D E F G A B.
   Le bonhomme prévient : « Attention, après G, l'alphabet repart de A. »
   Puis : « joue C D E F G A B ».

   Un gag en guise de chute, la sensible avant son nom. Sur le dernier
   si, le jeu garde la note enfoncée : le son reste, sans que le clavier
   à l'écran la montre tenue. Le bonhomme a l'air tendu : « Humpf ! Ça
   peut pas finir comme ça ; il manque quelque chose… » Le do au-dessus
   clignote sur le clavier. Le joueur le joue : « Ah ! Merci ! Ça va
   beaucoup mieux. » Le mot de sensible ne vient qu'au chapitre 4 ;
   le joueur en a déjà senti l'effet.
4. *Montrer* : une astuce pour deux notes. Do et fa sont à gauche des
   groupes de deux et de trois touches noires ; mi et si sont à droite
   de ces mêmes groupes.
5. *Demander* : un petit exercice qui renforce ces quatre touches,
   « joue do fa mi si do », puis « joue C F E B C », puis d'autres
   enchaînements des quatre.
6. *Les motifs*, en miroir de ceux de 1.1 : do-si-do, do-la-do,
   do-sol-do, en « répète après moi », en chantant. Cette fois, le do
   est à la main droite et l'autre note à la main gauche, qui descend
   de touche blanche en touche blanche. Avec le sol, ces motifs
   couvrent un tétracorde complet, sol la si do.

L'activité de la leçon 1.2 met ensemble tous les motifs de 1.1 et de
1.2, et sa séquence longue est la gamme entière : do ré mi fa sol la si,
et C D E F G A B.

### 1.3 Les dièses

**La leçon 1.3** nomme les touches noires par leurs dièses. Elle donne
au joueur le sens de l'orientation, avec le vocabulaire qui va avec : on
**monte** vers la droite du clavier, on descend vers la gauche.

1. *Dire* : « Les touches noires ont des noms aussi. Un dièse, c'est une
   note qu'on a fait monter d'un cran, vers la droite. » Puis le signe
   et les deux notations : do♯ et C♯.
2. *Jouer*, puis *répéter après lui* : do puis do♯, fa puis fa♯. La
   note monte d'un cran. Le mot de demi-ton attend le chapitre 3.
3. *Demander* : les cinq touches noires par leur dièse, do♯ ré♯ fa♯ sol♯
   la♯, en do ré mi puis en lettres, mélangées.
4. *Demander*, à deviner : « Et mi♯ ? » Le joueur a trois chances, sans
   qu'on les lui compte : chaque raté sonne la casserole, avec « Non, il
   n'est pas là. Rappelle-toi : un dièse, c'est une note qu'on a fait
   monter d'un cran. » Puis de même pour si♯. Au bout de trois essais
   infructueux, le bonhomme donne la réponse : « Mi♯ est sur le fa, et
   si♯ est sur le do. » Trouvées ou données, il nomme la chose : mi♯ et
   fa sont des notes **enharmoniques**, deux noms pour une même touche.
   L'astuce de 1.2 l'annonçait : mi et si sont collés à droite d'un
   groupe de touches noires, sans touche noire après eux.
5. *Demander*, la séquence longue :
   « joue C C♯ D D♯ E F F♯ G G♯ A A♯ B C ». Elle monte touche par
   touche et prépare la basse en chromatisme du chapitre 3.

L'activité de la leçon 1.3 ne reprend pas les motifs de 1.1 et 1.2.
Elle demande des séquences de notes diésées, dites en do ré mi et en
lettres, et y glisse de temps en temps un mi♯ ou un si♯.

### 1.4 Les bémols

**La leçon 1.4** est le miroir de 1.3 : les bémols, en descendant.

1. *Dire* : « Un bémol, c'est une note qu'on a fait descendre d'un cran,
   vers la gauche. » Puis le signe et les deux notations : ré♭ et D♭.
2. *Jouer*, puis *répéter après lui* : ré puis ré♭, si puis si♭.
3. *Demander* : les cinq touches noires par leur bémol, ré♭ mi♭ sol♭ la♭
   si♭, en do ré mi puis en lettres, mélangées.
4. *Dire*, puis *demander* : « joue do♯ », puis « joue ré♭ » ; c'est la
   même touche. Chaque touche noire a deux noms, et le bonhomme le
   prend sur un ton léger : « Les musiciens aiment bien compliquer les
   choses. » Do♯ et ré♭ sont enharmoniques, comme mi♯ et fa.
5. *Demander*, à deviner, comme à la leçon 1.3 : fa♭, puis do♭, trois
   chances chacun, la casserole à chaque raté avec « Non, il n'est pas
   là. Rappelle-toi : un bémol, c'est une note qu'on a fait descendre
   d'un cran. » Au bout de trois essais, la réponse : « Fa♭ est sur le
   mi, et do♭ est sur le si. »
6. *Demander*, la séquence longue, qui descend cette fois :
   « joue C B B♭ A A♭ G G♭ F E E♭ D D♭ C ».
7. *Dire*, en conclusion : les douze notes jouées touche par touche, en
   montant à la leçon 1.3 et en descendant ici, forment l'**échelle
   chromatique**.

L'activité de la leçon 1.4 demande des séquences de notes bémolisées,
dites en do ré mi et en lettres, et y glisse de temps en temps un fa♭ ou
un do♭.

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
