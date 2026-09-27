# Les voicings

Comment gohar pense la disposition des accords au clavier. Ce document
tient les décisions prises avec Arnaud, les questions laissées
ouvertes, et le lien avec le dex et les grilles. Il ne tient pas de
cours : les règles viennent de la pratique de l'école Maury, telle
qu'Arnaud l'a apprise d'Étienne Guéreau et de Mathias Berger-Forestier.

Le vocabulaire des shells, drops et positions sans fondamentale A et B
est celui de la méthode Berklee. Il n'est pas celui de gohar.

## La position

Un voicing se caractérise par sa **position** : la suite des degrés de
l'accord, de la basse vers le soprane. CMaj7(add9, add13) en position
1-7-9-3-13 donne do, si, ré, mi, la.

- **Les numéros n'ont pas d'altération.** G7(♭9, ♭13) en position
  1-7-3-13-9 : c'est « la treizième et la neuvième de l'accord », qui
  se trouvent abaissées sur cet accord. La position est indépendante de
  la couleur, qui vient de l'accord.
- **Les numéros suivent le chiffrage.** Un C6 a une sixte, pas une
  treizième. Changer d'avis sur la lecture d'un accord ambigu oblige
  donc à réévaluer sa position : la position se calcule à partir des
  hauteurs jouées **et** d'une lecture de l'accord, et une lecture
  ambiguë donne plusieurs positions, sans en choisir une en douce.
- **Un renversement est une position** dont la basse n'est pas la
  fondamentale. La position fondamentale est un cas particulier.
- **La basse peut être tenue par un autre.** Pour donner une consigne
  au pianiste, la position se dicte sans la basse. Pour analyser
  l'ensemble, elle l'inclut, en précisant qu'elle est tenue par le
  bassiste.

## La réalisation

Une position fixe un ordre, pas des intervalles exacts. Une même
position admet plusieurs **réalisations** : 1-3-7-5 peut se jouer avec
une dixième plutôt qu'une tierce en bas. Ce qui borne cette liberté,
c'est la taille de la main.

Dans l'aigu, la limite la plus évidente en piano solo est la mélodie.

## Les règles de tessiture

Dans le grave, l'oreille fixe pour chaque intervalle une note la plus
grave en dessous de laquelle il devient trouble. Par exemple, do3-mi3
est l'une des tierces majeures les plus graves qu'on fera entendre, et
do2-mi3 l'une des dixièmes les plus graves.

- La règle porte sur **la note la plus grave de l'intervalle**, et vaut
  en théorie pour **toutes les paires de voix**, pas seulement les
  voisines : une septième trop basse contre la fondamentale reste
  trouble même si la quinte qui les sépare est propre.
- **L'octave n'a pas de limite** : doubler la basse à l'octave
  inférieure (do2-do3-mi3) n'est pas « fautif ».
- **Les règles sont indicatives.** On peut les enfreindre, mais
  délibérément, avec une intention claire sur le son voulu, plutôt que
  « parce que je ne connais pas d'autre position ». Dans gohar, une
  réalisation hors tessiture est **annotée**, jamais refusée : c'est le
  même principe que le dex, qui constate et ne juge pas.
- **La table de départ** est celle que Berklee appelle *low interval
  limits*, à corriger d'après le premier chapitre du cours « Les bases
  de l'harmonisation » d'Étienne Guéreau.

Les octaves se comptent avec do4 pour le do du milieu et la4 = 440 Hz
(voir `glossaire.md`).

## Le noyau des positions

Les positions serrées (1-3-5-7 et ses renversements serrés) sont
« sans intérêt » : on les fait d'instinct, et apprendre à changer de
position, c'est justement en sortir.

Le noyau, dans l'ordre où on l'apprend :

1. **1-3-7 et 1-7-3**, à trois sons, sans la quinte.
2. **1-3-7-5 et 1-7-3-5**, les mêmes avec la quinte ajoutée.

Les positions des accords 6 et m6 sont distinctes de celles des
tétrades. Le II-V-I mineur se termine sans doute sur un m6 en 1-3-6-5,
à vérifier dans le cours de Guéreau (qui inclut la quinte sur les
accords de tonique).

## Les enchaînements

Dans la pratique, on apprend des **couples de positions**, mais ces
couples dérivent de la **règle de moindre mouvement** pour la conduite
des voix. Le premier couple est 1-3-7 qui alterne avec 1-7-3 sur un
II-V-I. En do :

| Accord | Position | Notes |
|---|---|---|
| Dm7 | 1-3-7 | ré, fa, do |
| G7 | 1-7-3 | sol, fa, si |
| CMaj7 | 1-3-7 | do, mi, si |

Seule la basse saute. Sur une basse qui descend par quintes, la tierce
d'un accord devient la septième du suivant, et la septième glisse d'un
demi-ton pour devenir la tierce : la position alterne d'elle-même.

gohar sait donc **dériver** les couples par le calcul, et le dex
**retient** ceux que le joueur connaît. Le calcul dit pourquoi, le dex
dit ce qui est su.

## Ce que le dex en retient

- **Une position de tétrade est une notion élémentaire** au sens du
  dex (voir `dex.md`) : un apprenant qui débute l'harmonisation veut
  les cocher, un pianiste aguerri les utilise comme il respire.
- **Elle se suit par tétrade et par fondamentale** : Maj7, m7, 7,
  m7♭5, dim7 et mMaj7, sur douze fondamentales. Une position du noyau
  donne un tableau de 72 cases, pas 72 entrées : le dex montrera
  « 1-7-3 : produite sur m7 dans 9 tonalités sur 12 », comme il montre
  la production d'un mode tonique par tonique. Aujourd'hui `Produced`
  se détaille par tonique ; il faudra qu'il se détaille par un contexte
  propre à la notion.
- **Une progression voicée depuis une position de départ est une
  entrée.** Le II-V-I qui part de 1-3-7 et celui qui part de 1-7-3 sont
  deux entrées sœurs : deux compétences distinctes pour les mains.
- **C'est d'abord une notion de la main.** La production compte, et
  l'emploi (la marque « Employée ») prend ici tout son sens : un
  pianiste qui place spontanément du 1-7-3 sur les changes d'une grille
  emploie la position en situation. C'est le terrain des grilles iReal
  (voir `chantiers.md`).

## Ouvert, à trancher sur le terrain

Ce sont des décisions pédagogiques, que les élèves, les playtests et
l'expérience trancheront. Elles vivront donc dans des **données**
(niveaux, curriculum), jamais dans le code, qui ne fournit que la
mécanique : positions, réalisations, annotations de tessiture, calcul
du moindre mouvement, suivi dans le dex.

- **Entendre une position** à l'oreille : pas indispensable pour
  l'instant, peut-être plus tard.
- **La suite du curriculum** après le noyau : a priori, une extension à
  la place de la quinte (1-7-3-13 sur la dominante, par exemple).
- **Les voicings avec extensions dans le dex** : seulement s'ils vivent
  dans des enchaînements courants, comme le II-V-I au sus4(♭9) du cours
  de Guéreau, dont on ne sait pas encore quoi faire.

## Le nom des accords

En français, on **écrit** les chiffrages avec les lettres américaines
et on les **dit** avec les syllabes latines : « CmMaj7(add9) » se dit
« do mineur majeur 7 add 9 ». Le rendu d'un chiffrage en tient compte
(voir `glossaire.md`), et c'est ce qui compte pour l'accessibilité : une
synthèse vocale doit dire ce qu'un musicien dirait.
