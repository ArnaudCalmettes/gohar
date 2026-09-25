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
- **Langue** : français ou anglais, par option.
- **Affichage** : un clavier vectoriel qui anime la note qui sonne.

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
  `at` sert à mesurer, pas à différer. Piloter la séquence depuis
  `Update` donnerait une gigue d'une frame (16,7 ms). Proposition :
  une `keyboard.Source` qui rejoue une séquence datée sur sa propre
  goroutine, avec des timers. C'est la « séquence rejouée » que
  `keyboard.go` annonce déjà comme une vraie source, et sa gigue est
  celle d'un doigt sur le clavier MIDI : le tampon du périphérique.
- **Deux sources, un moteur, un affichage.** Le rappel d'une `Source`
  doit rendre la main tout de suite. Chaque événement va au moteur
  directement, et vers le jeu par une file que `Update` vide.
- **La persistance n'existe pas.** `Notion.String` donne la forme
  écrite d'une notion, mais rien n'écrit un `Dex`. Sans elle, la
  collection repart de zéro à chaque lancement et on ne valide rien.

## Ce qu'il ne faut pas fermer, pour WASM

- Le jeu doit tourner sans clavier MIDI. `midi.go` passera derrière un
  build tag (`!js`), avec `webmididrv` plus tard.
- La persistance s'écrit sur un `io.Writer` et se lit sur un
  `io.Reader`. Le fichier côté bureau, `localStorage` côté navigateur,
  c'est l'affaire du jeu.
- Le navigateur ne démarre l'audio qu'après un geste de l'utilisateur :
  prévoir un écran « cliquer pour commencer », utile aussi au bureau.

## Le découpage

1. **Persistance du dex** : `MarshalJSON` et `UnmarshalJSON`, clés en
   `Notion.String`, dates en RFC 3339, test d'aller-retour. Et
   suppression de `Fact.Correct` : aucun fait ne dit qu'une tentative a
   échoué, `FactNamed` n'est envoyé que sur une identification. Faible.
2. **Séquence rejouée** : une `keyboard.Source` qui joue une liste
   d'événements datés. Faible à moyen.
3. **Histogramme des délais** dans `synth.Engine`, à cases fixes, sans
   allocation sur la goroutine audio. Faible.
4. **Le jeu en texte** : `games/ear`, Ebitengine, noms cliquables,
   `ebitenutil.DebugPrint` pour tout le reste. La boucle complète, dex
   compris. Moyen.
5. **Le clavier vectoriel** : `ebiten/v2/vector`, deux octaves,
   animation de la note qui sonne, révélation du mode. Moyen.
6. **Charge** : histogramme affiché par une touche de debug pendant que
   le clavier s'anime. Faible, une fois 3 et 5 faits.
7. **WASM** : build tags, stockage, écran de démarrage. Plus tard.