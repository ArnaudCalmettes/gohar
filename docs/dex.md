# Le dex

Le dex est la collection de notions musicales que le joueur se
constitue. Il est **commun à tous les jeux** et n'en connaît aucun.

Ce document fixe ce qu'il contient, ce que les jeux lui envoient et ce
qu'ils peuvent lui demander. Il a précédé le code, comme le glossaire a
précédé la bibliothèque, et il le suit depuis : quand l'un change,
l'autre aussi.

## Le principe qui commande tout le reste

**On célèbre ce que le joueur connaît, on ne lui montre pas ce qui lui
reste.** Apprendre la musique est l'affaire d'une vie et le total n'a
aucun sens.

Conséquence technique, à tenir fermement parce que la tentation
reviendra : **le dex ne calcule jamais de dénominateur global**. Pas de
total, donc pas de pourcentage, pas de barre de progression, pas de
« 47 sur 312 ». Un compteur d'acquis se dit seul, « 23 notions », jamais
sur fond d'inventaire.

Une exception, et une seule : on peut compter à l'intérieur d'un
ensemble que **la musique elle-même borne**. Les sept degrés d'un
système, les douze toniques d'une notion. « Plus que trois modes dans le
mineur harmonique » motive, parce que sept est un fait théorique et non
un objectif administratif.

Le catalogue de progressions, lui, n'est borné que par nos ajouts. Y
afficher un compte serait un mensonge qui se rétracte le jour où on
l'étoffe.

Le mot « dex » est gardé pour la référence qu'il évoque. Le ton, lui,
n'en reprend rien : aucune injonction à tout attraper, aucune formule
qui parle de ce qui manque en dehors du cas borné ci-dessus.

## Apprendre au joueur à se passer du jeu

Un principe qui vaut pour tous les jeux du projet, posé dans le billet
sur la game loop de l'improvisateur : **on ne veut surtout pas rendre
un musicien prisonnier de son écran pour jouer**. Le but n'est pas que
le joueur revienne au jeu, c'est qu'il n'en ait plus besoin, et qu'il
reconnaisse et emploie une couleur au clavier, là où aucun jeu ne lui
dit ce qu'il entend.

Le dex sert ce but, il ne le remplace pas. C'est le joueur qui
s'évalue, presque tout le temps : il sait ce qu'il reconnaît et ce qui
lui échappe encore. Le dex garde ce qu'il a montré et ne lui en rend un
bilan honnête que quand il le demande. Il ne le juge pas, ne le relance
pas, et ne fabrique aucune raison de jouer qui ne soit pas la musique
elle-même.

## Ce qu'est une entrée

Une notion que le joueur peut **reconnaître**, **produire**, et dont un
jeu peut constater la présence. Les trois à la fois.

Ça admet les modes, les tétracordes, les tétrades, les progressions du
catalogue. Ça écarte les notions qui se comprennent sans se jouer : on
ne joue pas « une sous-dominante » comme on joue un lydien dominante.

Le dex est donc une collection d'**objets sonores**, pas de savoirs.
C'est une limite réelle, et il ne faut pas la forcer : une partie de la
théorie n'y logera jamais.

Une gamme n'est pas une entrée de plus : c'est un mode. La mineure
harmonique, c'est l'éolien ♮7, premier mode de son système, et le
mineur naturel, c'est l'éolien.

### Les notions élémentaires

L'**intervalle** est une entrée, mais élémentaire : un petit pas qu'on
récompense chez le grand débutant et qu'un musicien cesse de compter.
Un II-V-I altéré fait sonner des dizaines d'intervalles, et dire que le
joueur y a entendu une tierce mineure n'aurait aucun sens. D'où la
règle :

**Une notion élémentaire ne reçoit jamais `Sonné`.** Elle n'est marquée
que par une activité qui la demande. Ses marques cessent de croître
d'elles-mêmes quand le joueur passe à ce qu'elle construit, et elle
n'apparaît pas dans l'état des lieux : une tierce acquise ne se révise
pas, les gammes qu'elle compose, si.

Un degré, lui, n'est pas une notion : c'est un exercice qui entraîne
des intervalles dans le contexte d'une gamme.

Chaque forme est une entrée de plein droit. Le ii-V-I majeur, le
ii-V-I mineur et le ii-V suspendu sont trois entrées, pas une avec des
variantes. Ça remplit plus de pages et ça évite un modèle à deux
étages.

Une notion se désigne comme la bibliothèque la désigne, jamais par une
chaîne libre : sans quoi le même ii-V-I finirait en trois exemplaires
qui s'ignorent.

## Les quatre marques

Elles sont **indépendantes**, pas des paliers.

| Marque | Ce qu'elle dit |
|---|---|
| Rencontrée | La notion a sonné et a été nommée au joueur |
| Reconnue | Le joueur l'identifie à l'oreille |
| Produite | Le joueur la pose au clavier, dans telle tonique |
| Employée | Le joueur l'a choisie à bon escient alors qu'il avait le choix |

Reconnaître et produire sont deux compétences distinctes, et l'une ne
suppose pas l'autre : on peut tenir la seconde sans la première, par
l'exploration au clavier. Un dex qui les empilerait en niveaux mentirait
sur ce chemin et refuserait de créditer ce que le joueur sait faire.

### La production se détaille par tonique, la reconnaissance non

L'oreille est relative : reconnaître une forme ne dépend pas de la
tonalité, c'est la définition même. Y compter les toniques n'apprendrait
rien.

Le clavier n'a pas cette invariance. Douze toniques, douze
topographies, douze doigtés. C'est même le cœur de l'entraînement :
poser la chose partout est ce qui muscle l'oreille relative.

La production garde donc **par tonique** sa date et son compte de
renforcements. L'agrégat par entrée se calcule à la demande ; l'inverse
ne s'invente pas.

### Employée

Cette marque n'exige pas de mesurer l'improvisation libre, seulement que
**le joueur ait eu le choix**. Un slot qui demande une dominante en
laissant la couleur ouverte, le joueur pose un lydien dominante : ça se
coche. Le même mode imposé par le jeu ne coche rien, si bien joué
soit-il.

Juger que la couleur était ouverte revient au jeu, pas au dex : le dex
ne lit aucune cible, il reçoit un `FactChosen` ou n'en reçoit pas. La
couleur est ouverte tant que la cible ne nomme pas le mode. Une tétrade
ne donne au mieux qu'une fonction et une indication de couleur (un
m7♭5 implique une sous-dominante en mineur) ; ce sont les extensions et les degrés
caractéristiques qui disent qu'une couleur a été choisie. Sur un slot
`G7`, le lydien dominante n'est attribué que si le ♯11 a sonné.

Qui a choisi a joué : `FactChosen` marque aussi la production sur sa
tonique. Sans quoi un mode toujours improvisé et jamais exigé
n'entrerait jamais dans la grille des douze toniques.

## La fraîcheur

Aucune marque ne se retire jamais. **Un acquis reste acquis.** Un fait
négatif qui dégraderait une entrée serait punitif, et le dex ne fait
aucun reproche.

Ce qui passe, c'est le temps. Chaque marque porte sa date de dernier
renforcement et son compte, ce dernier étant indispensable au rappel
espacé : une notion revue une fois il y a un mois et une notion revue
huit fois ne demandent pas le même intervalle.

Une entrée ne se dégrade donc pas, elle **tiédit**, et le dex propose au
lieu d'accuser.

La fraîcheur n'est visible **que sur demande**. Le joueur qui vient
chercher un état des lieux accepte une réponse franche ; le même bilan
affiché de force serait une remontrance. D'où deux vues distinctes, et
non une vue avec un bouton :

- **la collection**, ouverte pour le plaisir : ce que je possède, ce que
  j'ai croisé. Aucune notion de temps ;
- **l'état des lieux**, qu'on va chercher : ce qui refroidit, depuis
  quand, dans quelles toniques.

## Le dévoilement

Une entrée jamais croisée ne s'affiche pas. On ne court pas après ce
qu'on ne voit pas, et un débutant devant trente-cinq cases vides lit
d'abord tout ce qu'il ignore.

Une entrée devient visible en silhouette quand **ses composants sont
déjà connus**. Le ii-V-I mineur apparaît dès que le joueur possède le
mineur 7 bémol 5, la dominante et le mineur majeur 7. Avant, il n'existe
pas pour lui.

Une notion peut aussi apparaître en silhouette parce qu'elle a **sonné
sans être nommée** : le musicien qui joue une couleur parce qu'elle
sonne classe, bien avant d'en connaître le nom. C'est `Overheard`, qui
n'est pas une des quatre marques : la notion n'entre pas dans la
collection, elle devient visible. Le jeu peut alors la nommer (« ce mode
s'appelle phrygien ♮6, tu veux le capturer ? »), et c'est cette
nomination qui fait la découverte.

Une famille apparaît de même une fois entamée, et montre alors ses
places. Le compte « plus que trois » arrive ainsi au moment où il
motive, à quelqu'un déjà dedans, et jamais à l'accueil.

Le dex se dévoile donc par ce que le joueur sait, pas selon un ordre de
leçons décidé d'avance. C'est l'encyclopédie qu'on se construit dans son
ordre à soi, rendue mécanique.

## Ce que les jeux envoient

Un **compte rendu** par activité terminée, pas un flux. Il porte le jeu,
l'activité, l'instant de fin, et une liste de faits.

Ni score, ni réussite globale, ni difficulté, ni combo : le vocabulaire
d'un jeu n'entre pas ici, sinon le dex cesse d'être un bien commun.

**Le dex reçoit des faits, il ne juge pas.** Décider qu'une réponse est
juste, qu'une couleur était ouverte, qu'un trajet harmonique tient,
appartient aux règles du jeu, qui émet ensuite les faits
correspondants. Le dex décide seulement ce que chaque fait marque.

| Fait | Ce qu'il porte | Ce qu'il marque |
|---|---|---|
| Entendu | la notion, la tonique | Rencontrée |
| Nommé | la notion, la tonique | Reconnue |
| Produit | la notion, la tonique | Produite sur cette tonique |
| Choisi | la notion, la tonique | Employée, et Produite sur cette tonique |
| Sonné | la notion | rien : rafraîchit ce qui est déjà marqué, ou fait apparaître une silhouette |

La première marque posée sur une notion est sa **découverte**, et c'est
le moment que le jeu peut célébrer.

**Sonné** est le fait de fraîcheur, et il répond à un piège du rappel
espacé : les notions ne sont pas indépendantes. Réviser un ii-V-I fait
sonner un mineur 7, une dominante et un majeur 7, et sans ce fait le
joueur se verrait proposer demain ce qu'il vient de pratiquer sans le
savoir. C'est au moteur de reconnaissance de le produire, jamais à
l'auteur du jeu, qui rapporterait ce qu'il a programmé plutôt que ce qui
a retenti.

**Aucun fait ne dit qu'une tentative a échoué**, ni qu'une performance
était bonne. Le premier est punitif, le second est un jugement
déguisé : ce n'est pas au dex d'évaluer le joueur, c'est au joueur de
s'évaluer, et il n'obtient un bilan honnête que quand il le demande.

**Une erreur compte pour du beurre.** Elle se corrige, elle ne
s'apprend pas : ce qui s'apprend, c'est la correction que le joueur
produit ensuite. Un jeu n'envoie donc `Nommé` que sur une
identification, et la correction réussie arrive comme n'importe quelle
autre, sans que le dex sache qu'une erreur l'a précédée.

## La forme écrite

Le dex se sérialise en JSON, et c'est tout ce qu'il sait de son
stockage : un fichier au bureau, le stockage local dans un navigateur,
c'est l'affaire du jeu. Une notion s'y écrit sous sa forme stable
(`mode:0/4`), jamais sous les valeurs numériques de ses champs, et une
marque jamais posée ne s'écrit pas. La lecture est stricte : ce qui ne
se relit pas exactement est une erreur, jamais une devinette, parce
qu'une collection mal relue déplacerait des marques en silence.

## Ce que les jeux demandent

| Question | Réponse | Usage |
|---|---|---|
| Que sait-il | l'état d'une ou plusieurs entrées | les deux vues |
| Qu'est-ce qui refroidit | entrées et toniques dont le rappel est dû, par urgence | composer une révision |
| Que peut-il rencontrer | entrées jamais croisées dont les composants sont connus | composer de la découverte |
| Qu'est-ce qui est visible | ce que le joueur a le droit de voir, silhouettes comprises | le dévoilement |

Cette seconde interface mérite autant de soin que la première : c'est
elle qui rend le dex utile plutôt que décoratif.

Le **dosage** entre révision et découverte appartient au jeu. Le dex
répond aux deux questions et ignore la part que chacun en fait. Sans
quoi les jeux piocheraient surtout dans le tiède, le joueur ne
rencontrerait plus rien de neuf, et la collection cesserait d'être une
exploration.

## Ce qui reste à décider

- Le catalogue de progressions, à dicter comme l'ont été les 35 modes.
  Sans urgence : rien n'oblige plus à le borner.
- Les extensions d'accord. Une neuvième mineure sur une dominante est
  exactement le genre de notion qu'on aimerait collectionner, mais elle
  n'est pas dénombrable et ne passe pas le critère d'admission. Il y a
  peut-être une entrée à trouver du côté de « telle couleur sur telle
  fonction ».
- Le nom du paquet dans le dépôt, qui n'est pas forcément le mot montré
  au joueur.
