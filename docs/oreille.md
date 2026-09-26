# Première activité d'oreille

Le tracer bullet : une activité étroite qui traverse toutes les couches,
`synth`, `keyboard`, `harmony`, `naming` en français et en anglais,
`dex`, et un premier affichage. Ce n'est pas une production. Ce qu'elle
doit rendre, c'est la liste de ce qui coince.

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
- **Réponse** : au clic sur un nom. Le clavier MIDI sert à rejouer ou à
  explorer librement pendant la question.
- **Après une erreur** : la bonne réponse est révélée, puis les deux
  modes sont joués l'un après l'autre sur la même tonique. C'est la
  correction qu'on montre, pas un fait qu'on retient.
- **Reprise** : le mode raté revient deux ou trois questions plus loin,
  sur une autre tonique, une seule fois par série. Ce qui s'apprend,
  c'est la correction produite par le joueur.
- **Session** : dix questions, plus au plus une reprise par erreur. Un
  seul `Report` à la fin, les `Change` affichés (découvertes, marques
  gagnées).
- **Langue** : français ou anglais, par `-lang` ou la touche `L`.
- **Notation** : signes par défaut (si♭, phrygien ♮6), mots en option
  (si bémol, phrygien bécarre 6), par `-notation` ou la touche `N`.
- **Affichage** : trois boutons et un clavier vectoriel de do3 à do6,
  qui s'allume sous ce qui sonne, de n'importe quelle source, et
  s'éteint en un quart de seconde. Ce qui sort de la plage (la pédale,
  un clavier MIDI réglé sur une autre octave) est ramené dedans et
  affiché atténué. À la révélation, les touches du mode se colorent en
  entier (clair sur les blanches, foncé sur les noires), la tonique en
  bleu, avec le nom de chaque note orthographié dans le mode (mi♭ et
  fa♯ en sol mineur harmonique) et toujours en signes, faute de place.
  Après une erreur, les notes du seul bon mode sont en vert, celles du
  seul mode choisi en rouge, les communes en gris : ce qui les sépare
  saute aux yeux. Repris de la démo gohareact de l'ancien gohar.

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
  entier se surligne à la révélation.
- **Jouer une séquence à l'heure.** `Engine.NoteOn` n'ordonnance rien :
  `at` sert à mesurer, pas à différer, et piloter la gamme depuis
  `Update` donnerait une gigue d'une frame (16,7 ms). Résolu par
  `keyboard.Sequence`, une `Source` qui rejoue des notes datées sur sa
  propre goroutine : sa gigue est celle d'un doigt sur le clavier MIDI,
  le buffer du périphérique.
- **Deux sources, un moteur, un affichage.** Le rappel d'une `Source`
  doit rendre la main tout de suite. Résolu par `onKey`, seul chemin de
  toute source vers le moteur et l'affichage.

## Ce que le premier playtest a appris

- Sous Linux, le port « Midi Through » arrive en tête de liste et ne
  joue rien. `OpenMIDI` le saute quand aucun port n'est demandé.
- Go Regular n'a pas ♭ ♮ ♯, ni le double dièse 𝄪 et le double bémol 𝄫,
  qui sont des signes à part entière. Une coupe de Noto Music de 2 Ko,
  embarquée (`games/ear/fonts`), sert de police de secours, et `naming`
  écrit les doubles avec leur propre signe.
- Un layout de 640 × 360 agrandi par Ebitengine pixellise tout. Le jeu
  dessine maintenant à la résolution réelle de la fenêtre, en gardant
  ses coordonnées logiques (`canvas.go`).
- Le français écrivait les notes en mots et l'anglais en signes. La
  notation est devenue un choix du `Namer`, indépendant de la langue,
  pour les deux langues à la fois : c'est aussi ce qui ouvre la porte à
  une synthèse vocale.

## Ce qu'il ne faut pas fermer, pour WASM

- Le jeu doit tourner sans clavier MIDI. `midi.go` passera derrière un
  build tag (`!js`), avec `webmididrv` plus tard.
- Le dex ne sait que se sérialiser en JSON. Où le ranger est l'affaire
  du jeu : un fichier au bureau (`store.go`), `localStorage` dans le
  navigateur.
- Le navigateur ne démarre l'audio qu'après un geste de l'utilisateur :
  l'écran « cliquer pour commencer » existe déjà.
- La police des signes est embarquée, pas lue sur le système.

## Le découpage

Fait :

1. **Persistance du dex**, et suppression de `Fact.Correct`.
2. **Séquence rejouée** : `keyboard.Sequence`.
3. **Histogramme des délais** dans `synth.Engine`, à cases fixes, sans
   allocation sur la goroutine audio.
4. **Le jeu en texte** : `games/ear`, la boucle complète, dex compris,
   texte en `text/v2` avec la police de secours. La touche `H` affiche
   déjà l'histogramme.
5. **Le clavier vectoriel** : `piano.go`, trois octaves, animation et
   révélation.
6. **Charge** : dix minutes de jeu, 2 208 événements, p99 et maximum à
   11 ms, le ramasse-miettes sans effet sur le son. Détail dans
   `architecture.md`.

7. **Timbres 8 bits** : impulsions à 12,5 et 25 %, carré, triangle,
   bruit, enveloppe ADSR. Rendu adouci par défaut, `-authentic` pour
   le grain brut des consoles. `-timbre` choisit le son.

Reste :

8. **Soundfont** : go-meltysynth dans `synth/soundfont`, derrière la
   même interface `Instrument`, puis refaire la mesure 6. Le lecteur
   est étudié (voir `chantiers.md`) ; il manque un piano SF2 réduit.
9. **WASM** : build tag sur `midi.go`, stockage dans `localStorage`.

## La feuille de route

Une fois le tracer bullet bouclé, `ear` devient le jeu d'oreille à
part entière.

**L'abstraction, faite.** Les activités ont toutes la même forme :
faire sonner quelque chose, proposer des réponses, émettre des faits.
Seuls changent ce qui sonne, les réponses et les notions.

- `Question` est une donnée : la tonique, les choix sous forme de
  notions du dex (le jeu sait les nommer dans toute langue et toute
  notation), l'index de la bonne réponse. Ce qui sonne, c'est le bon
  choix, et l'activité ne garde pas d'autre secret.
- `Activity` est ce qui diffère : planifier une série, reprendre une
  question après une erreur (autre tonique pour un mode, même tonalité
  pour un degré), l'énoncé, ce qui sonne, la correction, ce que montre
  le clavier, les faits d'une bonne réponse.
- `Series` est générique : l'erreur qui compte pour du beurre, la
  reprise unique, le rapport. Ce sont les principes du jeu, pas d'une
  activité.
- Les modes sont la première activité (`modes.go`), les tétracordes
  la deuxième (`tetrachords.go`).

**Les niveaux sont des données** : une activité et ses paramètres
(quels degrés, quel système, quelle distance entre distracteurs,
combien de questions). Un défi composé par le joueur n'est qu'un niveau
dont il a choisi les paramètres.

**Répondre en jouant** : la première note jouée après la question est
la réponse, ce qui sonne pendant la question elle-même étant ignoré.
Une bonne réponse jouée est aussi une production.

**Le menu, fait** : toutes les activités jouables d'emblée. Imposer un
ordre d'étude frustrerait un joueur qui a déjà des bases ; le
dévoilement attendra un meilleur terrain d'essai. Les activités y sont
dans l'ordre de la progression, une touche chiffrée chacune. En fin de
série, Espace relance la même activité et Entrée revient au menu (pas `M`,
qui tombe sur la virgule d'un clavier AZERTY).

**Des réglages pour le joueur**, dans un fichier à part du dex :
l'ambitus et le tempo en premier.

**Une progression par compétences d'oreille** :

- **Degrés** (functional ear training) : une tonique posée, puis une
  note de la gamme, et le joueur répond son degré (1 à 7) plutôt que son
  nom. Quelques degrés d'abord, puis toute la gamme, puis d'autres
  gammes, et au plus difficile l'échelle chromatique, de l'unisson à
  l'octave. Poser la tonique demande une cadence, donc des accords, et
  la tonique reste fixe pendant une série. Le degré n'est pas une
  notion du dex : il sert à entraîner des **intervalles** depuis la
  tonique, qui en sont (élémentaires, voir `dex.md`). Et une série
  jouée qui trouve tous les degrés d'une gamme crédite la **production
  de cette gamme** sur sa tonique : la mineure harmonique, c'est
  l'éolien ♮7.
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
  les deux ne tiennent pas côte à côte. Les autres systèmes viendront
  par leurs tétracordes avant leurs modes.
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
