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

Reste :

7. **Soundfont** : étudier un lecteur SF2 en pur Go (allocations sur
   la goroutine audio, paquet à part puisque `synth` promet de ne
   dépendre que de la bibliothèque standard, licence et poids de la
   soundfont pour WASM), puis l'intégrer et refaire la mesure 6.
8. **WASM** : build tag sur `midi.go`, stockage dans `localStorage`.

## La feuille de route

Une fois le tracer bullet bouclé, `ear` devient le jeu d'oreille à
part entière.

**Une abstraction d'abord.** Les activités prévues ont toutes la même
forme : faire sonner quelque chose, proposer des réponses, émettre des
faits. Seuls changent ce qui sonne, les réponses et les notions. Une
interface `Activity`, extraite de `Series`, avant la deuxième
activité, pour ne pas écrire trois jeux copiés-collés.

**Des réglages pour le joueur**, dans un fichier à part du dex :
l'ambitus et le tempo en premier.

**Une progression par compétences d'oreille** :

- **Degrés** (functional ear training) : une tonique posée, puis une
  note de la gamme, et le joueur répond son degré (1 à 7) plutôt que son
  nom. Quelques degrés d'abord, puis toute la gamme, puis d'autres
  gammes. Poser la tonique demande une cadence, donc des accords.
  Question ouverte : un degré relatif est-il une notion du dex ?
- **Tétracordes**, ceux du système naturel d'abord : l'étape avant les
  modes.
- **Modes**, en continuant sur le système naturel. Les autres systèmes
  viendront par leurs tétracordes avant leurs modes.
- **À la fin**, tout mélanger, ou mieux, laisser le joueur composer ses
  propres défis.

Cette progression répond à `Notion.Components` : les composants d'un
mode sont ses deux tétracordes. Un mode apparaîtrait en silhouette
quand ses tétracordes sont connus, et la progression sortirait du
dévoilement du dex plutôt que d'une liste de niveaux.

**Ensuite, le volet harmonie** : triades, tétrades, tétrades avec
extensions, renversements, cadences et progressions.
