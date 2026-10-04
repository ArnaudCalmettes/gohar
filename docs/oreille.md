# L'ear trainer

`games/ear`, le jeu d'entraînement de l'oreille : quelque chose sonne
sur une pédale de tonique, et le joueur le nomme. Les règles ci-dessous
valent pour toutes les activités ; la feuille de route décrit la suite.

## Les règles, décidées

- **Contenu** : les sept modes du système naturel. On adaptera une fois
  le jeu en main.
- **Distracteurs** : trois choix par question, et des distracteurs
  lointains d'abord : au moins deux notes d'écart avec la réponse à
  tonique égale (le lydien n'est pas proposé face à l'ionien). Les
  distracteurs proches viendront en palier suivant.
- **Son** : une pédale de tonique, puis la gamme montante sur la pédale.
  La couleur se lit contre la pédale.
- **Tonique** : tirée au hasard à chaque question. `Recognized` n'est
  pas détaillé par tonique, on entraîne l'oreille relative d'emblée.
- **Réponse** : au clic sur un bouton ou par sa touche chiffrée. Le
  clavier MIDI sert à rejouer ou à explorer librement pendant la
  question.
- **Après une erreur** : la bonne réponse est révélée, puis les deux
  modes sont joués l'un après l'autre sur la même tonique, l'écran
  nommant celui qui sonne (« On entend : dorien »). `R` rejoue la
  correction entière. C'est la correction qu'on montre, pas un fait
  qu'on retient.
- **Reprise** : le mode raté revient deux ou trois questions plus loin,
  sur une autre tonique, une seule fois par série. Ce qui s'apprend,
  c'est la correction produite par le joueur.
- **Session** : dix questions, plus au plus une reprise par erreur. Un
  seul `Report` à la fin, les `Change` affichés (découvertes, marques
  gagnées).
- **Langue** : français ou anglais, par `-lang` ou la touche `L`.
- **Notation** : signes par défaut (si♭, phrygien ♮6), mots en option
  (si bémol, phrygien bécarre 6), par `-notation` ou la touche `N`.
- **Affichage** : une rangée de boutons numérotés et le clavier à
  l'écran de do3 à do6 (`screen.Piano`), qui s'allume sous ce qui
  sonne. À la révélation, les touches du mode se colorent, la tonique
  en bleu, avec le nom de chaque note orthographié dans le mode (mi♭ et
  fa♯ en sol mineur harmonique). Après une erreur, les notes du seul bon
  mode sont en vert, celles du seul mode choisi en rouge, les communes
  en gris : ce qui les sépare saute aux yeux. La touche `P` masque le
  clavier pendant la question, pour travailler à l'oreille seule.
- **Boutons** : toute la rangée s'écrit de la même façon, numéro et nom
  côte à côte en grande police si tout tient, sinon en petite, sinon
  le numéro au-dessus du nom.

## Les faits émis

- Réponse juste, du premier coup ou en reprise : `FactNamed` et
  `FactHeard` sur le mode, que le jeu nomme en confirmant. Le dex ne
  sait pas qu'il y a eu une erreur avant, et c'est voulu.
- Réponse fausse : rien. Une erreur compte pour du beurre : elle se
  corrige, elle ne s'apprend pas. Pas même de `FactHeard` pour la
  comparaison. Un mode ne se découvre donc qu'en le reconnaissant.
- Rien sur l'exploration libre au clavier : ce serait du `FactSounded`,
  qui attend le pont avec `analysis`.
- `Tonic` renseigné partout, même s'il ne sert pas à `Recognized`.

## Les pièges repérés

- **Le clavier ne doit pas donner la réponse.** Surligner les touches
  du mode avant la réponse revient à l'afficher. Pendant la question on
  n'anime que la note qui sonne et les touches enfoncées ; le mode
  entier se surligne à la révélation. Pour les degrés, la note qui
  sonne est déjà la réponse : en attendant que seul ce que joue le
  joueur s'allume pendant la question (voir `chantiers.md`), la touche
  `P` y pare.
- **Jouer une séquence à l'heure** depuis `Update` donnerait une gigue
  d'une image (16,7 ms) : `keyboard.Sequence` rejoue les notes datées
  sur sa propre goroutine (voir aussi « Les notes datées » dans
  `architecture.md`).
- **Deux sources, un moteur, un affichage** : le rappel d'une `Source`
  rend la main tout de suite, et `onKey` est le seul chemin de toute
  source vers le moteur et l'affichage.

- **Le papier fait paraître compliqué ce que l'oreille trouve
  simple.** Chailley raconte une lecture de Messiaen à Royaumont : les
  lecteurs, sur la partition, « Quelle complexité rythmique ! » ; les
  auditeurs, sans le papier, « Rien n'est plus simple. C'est du 4/4 en
  rubato » ; d'où « les illusions nées de la complexité sur papier »
  (*40 000 ans de musique*, p. 162-163). Le jeu à l'oreille doit poser
  ses questions sur ce qui sonne, jamais sur ce qui s'écrit ; et une
  grille chargée d'extensions peut s'entendre comme trois accords.

## Ce que le premier playtest a appris

- Sous Linux, le port « Midi Through » arrive en tête de liste et ne
  joue rien : `OpenMIDI` le saute quand aucun port n'est demandé.
- Go Regular n'a pas ♭ ♮ ♯ 𝄪 𝄫 : une coupe de Noto Music de 2 Ko,
  embarquée (`games/ear/fonts`), sert de police de secours.
- La notation, signes ou mots, est un choix du `Namer`, indépendant de
  la langue : c'est aussi ce qui ouvre la porte à une synthèse vocale.
- Ebitengine désigne les touches par leur place sur un clavier
  américain : sur un AZERTY, `M` tombe sur la virgule et `Q` sur le A.
  Le retour au menu est donc sur Entrée, et seul Échap quitte le jeu.

## Ce qu'il ne faut pas fermer, pour WASM

- Le jeu doit tourner sans clavier MIDI. `midi.go` passera derrière un
  build tag (`!js`), avec `webmididrv` plus tard.
- Le dex ne sait que se sérialiser en JSON. Où le ranger est l'affaire
  du jeu : un fichier au bureau (`games/settings`), `localStorage` dans
  le navigateur.
- Le navigateur ne démarre l'audio qu'après un geste de l'utilisateur :
  le clic dans le menu en tient lieu.
- La police des signes est embarquée, pas lue sur le système.

## Ce qui reste du premier jalon

- **Soundfont** : brancher le `Sampler` dans `ear` avec un piano SF2
  réduit, puis refaire la mesure de charge.
- **WASM** : build tag sur `midi.go`, stockage dans `localStorage`.

## La feuille de route

**Les activités ont toutes la même forme** : faire sonner quelque
chose, proposer des réponses, émettre des faits. `Question` est une
donnée (la tonique, les choix sous forme de notions du dex, la bonne
réponse) ; `Activity` est ce qui diffère (la série, la reprise,
l'énoncé, ce qui sonne, la correction, les faits) ; `Series` porte les
principes du jeu, l'erreur qui compte pour du beurre, la reprise
unique, le rapport.

**Les niveaux sont des données** : une activité et ses paramètres
(quels degrés, quel système, quelle distance entre distracteurs,
combien de questions). Un défi composé par le joueur n'est qu'un niveau
dont il a choisi les paramètres.

**Répondre en jouant** : la première note jouée après la question est
la réponse, ce qui sonne pendant la question elle-même étant ignoré.
Une bonne réponse jouée est aussi une production.

**Le menu** : toutes les activités jouables d'emblée, dans l'ordre de
la progression. Imposer un ordre d'étude frustrerait un joueur qui a
déjà des bases ; le dévoilement attendra un meilleur terrain d'essai.

**Des réglages pour le joueur**, dans un fichier à part du dex
(`games/settings`) : l'ambitus et le tempo en premier.

**Une progression par compétences d'oreille** :

- **Degrés, faits** (functional ear training). L'exercice du débutant
  au clavier : la main gauche tient la tonique, l'index droit tombe au
  hasard dans l'octave au-dessus, et on chante la gamme jusqu'à la note
  pour trouver son numéro.
  - *Son* : une pédale de tonique dans le même registre que la note (la
    tonique à la basse est un palier au-dessus), puis une note de la
    gamme majeure. Pas d'accords : la pédale suffit à poser la tonique,
    les cadences viendront avec le volet harmonie.
  - *Énoncé* : il nomme la gamme (« Gamme : ré majeur ») plutôt que la
    tonique.
  - *Réponses* : les sept degrés, toujours à leur place. Un premier
    niveau réduit serait de trop : dix questions sur les sept degrés
    font très bien l'affaire.
  - *Tonique* : fixe pendant la série, tirée au hasard.
  - *Correction* : après la réponse, bonne ou non, la gamme rentre à la
    tonique la plus proche, de la tonique jusqu'à la note pour le
    tétracorde inférieur, de la note jusqu'à l'octave pour le
    supérieur. La pédale redescend à la basse pour éviter l'unisson au
    départ de la tonique. Après une erreur, le chemin du degré choisi
    suit.
  - *Dex* : le degré n'est pas une notion ; ce sont les
    **intervalles** depuis la tonique qui en sont (élémentaires, voir
    `dex.md`). Une série jouée qui trouve tous les degrés d'une gamme
    créditera la **production de cette gamme** sur sa tonique, avec la
    réponse jouée.
  - *Paliers suivants* : la tonique mobile, la tonique à la basse,
    d'autres gammes, et au plus difficile l'échelle chromatique.
- **Tétracordes, faits** pour le système naturel : l'étape avant les
  modes. Les quatre formes (majeur, mineur, phrygien, lydien) sont
  proposées à chaque question, toujours à la même place : quatre
  choix n'appellent pas de distracteurs, et des places fixes laissent
  la main apprendre où vit chaque réponse pendant que l'oreille
  travaille. Le son est celui des modes, pédale comprise, mais sans
  octave ajoutée : un tétracorde s'arrête sur sa quatrième note, tenue.
  Le bouton porte le qualificatif seul (« phrygien »), la phrase le
  nom entier (« le tétracorde phrygien »).
- **Modes**, en continuant sur le système naturel. Au palier
  au-dessus, fait, les sept modes sont proposés à chaque question,
  chacun à la place de son degré : la touche 6 est le mode du sixième
  degré de la gamme mère. L'ordre des degrés plutôt que du plus clair
  au plus sombre, parce que l'enjeu est de retrouver instantanément la
  gamme mère d'un mode, et que ça commence par une mémorisation stricte
  des degrés. Les distracteurs proches viennent avec, sans règle pour
  les choisir. Sur une ligne, le numéro passe au-dessus du nom quand
  les deux ne tiennent pas côte à côte (voir Boutons, plus haut). Les
  autres systèmes viendront par leurs tétracordes avant leurs modes.
- **À la fin**, tout mélanger, ou mieux, laisser le joueur composer ses
  propres défis.

Cette progression répond à `Notion.Components`, écrit : les
composants d'un mode sont ses deux tétracordes (un seul quand ils sont
égaux, comme pour l'ionien). Un mode apparaîtra en silhouette quand ses
tétracordes sont connus, et la progression sortira du dévoilement du
dex plutôt que d'une liste de niveaux. Il manque pour cela le
recensement dont `Discoverable` a besoin.

**Ensuite, le volet harmonie** : triades, tétrades, tétrades avec
extensions, renversements, cadences et progressions.
