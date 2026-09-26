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
- **Affichage** : pour l'instant du texte et trois boutons. Le clavier
  vectoriel qui anime la note qui sonne reste à faire.

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
  entier se surligne à la révélation. À tenir pour le clavier
  vectoriel.
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
- Go Regular n'a pas ♭ ♮ ♯. Une coupe de DejaVu Sans de 4 Ko, embarquée
  (`games/ear/fonts`), sert de police de secours.
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

Reste :

5. **Le clavier vectoriel** : `ebiten/v2/vector`, deux octaves,
   animation de la note qui sonne, révélation du mode. Moyen.
6. **Charge** : lire l'histogramme après plusieurs minutes pendant que
   le clavier s'anime. Faible, une fois 5 fait.
7. **WASM** : build tag sur `midi.go`, stockage dans `localStorage`.
   Plus tard.
