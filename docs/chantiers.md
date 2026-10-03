# Chantiers

Ce qui est ouvert, et ce qui attend une décision.

Ce fichier existe pour qu'une conversation neuve reprenne sans rien
redécouvrir et sans rouvrir un débat déjà tranché. Les raisons des choix
sont dans `architecture.md`, le vocabulaire dans `glossaire.md`, la
collection dans `dex.md`, l'ear trainer dans `oreille.md`, les voicings
dans `voicings.md`, l'analyse des grilles dans `grilles.md`, *Walk with
me* dans `walk.md`. Ici, seulement ce qui reste à faire, les décisions à
ne pas rouvrir, et un bref état des lieux. Une entrée terminée s'en va :
ce qu'elle a tranché vit dans la doc de son sujet.

Mettre une doc à jour n'est pas un chantier : ça fait partie de la
tâche qui la fait mentir.

## Où on en est

- **`harmony`** : hauteurs, intervalles, gammes, tonalités, accords,
  fonctions, systèmes, tétracordes, les 35 modes, les préparations d'un
  accord à l'autre. Testé, benché, vert.
- **`naming`** : la langue (français, anglais) et la notation (signes,
  mots), les noms des intervalles, des gammes et des modes, et les
  symboles d'accords (`ChordStyle`).
- **`analysis`** : la reconnaissance déterministe des accords, et
  l'analyse d'une grille à la manière d'*En Harmonie* : blocs, degrés,
  tonique pressentie et modulations, cellules, forme, fins de section,
  demi-cadences, pédales, tonalité du morceau. Les fiches du livre et
  de Siron concordent ; sur le corpus, 87 % d'accord avec l'app ou
  l'oreille.
- **`charts`** : la lecture des grilles iReal Pro, de l'URL jusqu'aux
  accords en ordre de jeu, et trois commandes, `analyse` (une grille
  annotée et réorthographiée), `corpus` (l'analyse contre l'app sur des
  playlists entières) et `forms` (la forme).
- **`dex`** : le corps, la persistance JSON, `Components`.
- **`synth`** : un moteur polyphonique sans allocation, des timbres
  8 bits, la sortie audio et un histogramme des délais.
- **`games`** : `keyboard` (MIDI et séquences rejouées), `keys` (un
  clavier qui sonne), et `ear`, l'entraînement de l'oreille (degrés,
  tétracordes, modes du système naturel), avec le dex persisté.

## Les décisions à ne pas rouvrir

**Le dex** (détail dans `dex.md`) : il reçoit des faits et ne juge pas,
une erreur compte pour du beurre, `FactChosen` marque aussi la
production, la couleur reste ouverte tant que la cible ne nomme pas le
mode, et ce qui sonne sans être nommé apparaît en silhouette.

**Le nommage.** Les signes par défaut dans les deux langues, les mots
en option du `Namer` pour l'accessibilité. Le nom systématique d'un mode
par défaut (phrygien ♮3), les noms raffinés en alternatives
(`ModeAlternatives`).

**Les symboles d'accords** (`ChordStyle`) : Cmaj7, Cm7, Cm7♭5, Cdim7,
C6⁄9 par défaut, comme *En Harmonie* ; C♮7 ou CΔ7, C-7, Cø, C°7, C6/9
en option. Le bécarre ne se lit bien qu'en exposant : sur une ligne de
texte, C♮9 se lirait comme un do avec une neuvième bécarre. Les
extensions naturelles montent sur la septième (C13 : 9e et 13e ; Cm13 :
9e, 11e et 13e ; sur une tierce majeure, la 11e reste à part), le reste
va entre parenthèses sans espace, C7(♭9,♯11), et le mineur-majeur
s'écrit Cm(maj7,9).

**L'orthographe des grilles** : une sortie de l'analyse, jamais une
entrée (règles P0 à P3 dans `grilles.md`).

## Le dex

- [ ] la silhouette mise en scène dans les jeux (« ce pokémon s'appelle
      phrygien bécarre 6 ») : attend le pont `analysis` → `FactSounded`.
- [ ] `Cooling`, les intervalles de répétition espacée : la courbe se
      juge à l'oreille.
- [ ] `Discoverable` : demande un recensement de toutes les notions.
- [ ] `KindSystem` : une gamme est une entrée à part entière, distincte
      de ses modes. Acté, pas écrit ; `harmony.NamedScale` la désigne.
- [ ] la gamme comme référence tonale, distincte du système comme gamme
      mère : la mineure naturelle est une gamme de référence au même
      titre que la majeure.

## L'harmonie

- [ ] le pont `analysis` → `FactSounded` : repérer les degrés
      caractéristiques d'un mode sur une fenêtre de jeu (taille de la
      fenêtre, modes voisins).
- [ ] un champ `Tetrad` dans le catalogue des modes, extensions en
      motif : le catalogue actuel deviendrait l'oracle des tests.
- [ ] dire un chiffrage dans la langue (« do mineur majeur 7 add 9 »).
- [ ] le catalogue de progressions, pour que `ProgressionID` désigne
      quelque chose ; il servirait aussi à juger les détours d'une
      réharmonisation. Le contenu v0 est dans `grilles.md`.
- [ ] l'accord parallèle comme type de préparation : demande la
      mélodie.
- [ ] la basse chromatique sous d'autres accords que le diminué (*It
      Never Entered My Mind*), à côté de `PassingChords`.

## Les voicings

Les décisions sont dans `voicings.md`.

- [ ] les types `Position` et réalisation dans `harmony`, et le calcul
      de la position d'une réalisation.
- [ ] la table de tessiture, des *low interval limits* de Berklee
      corrigées à l'oreille ; des annotations, jamais de refus.
- [ ] le moindre mouvement, qui dérive les couples de positions d'une
      progression.
- [ ] dans le dex, le couple (tétrade, fondamentale) comme détail de
      production d'une position.
- [ ] les voicings sur une grille : la marque « Employée » constate une
      position placée spontanément sur les changes.

## La lecture des grilles

- [ ] un format de grille ouvert, mieux conçu que celui d'iReal, qui
      porterait aussi la mélodie et les couleurs modales.

## L'analyse des grilles

Les règles et les décisions sont dans `grilles.md`.

- [ ] *Yesterday's Gardenias* entendu en fa♯ ; il est en si♭ majeur.
      La règle du turnaround mis de côté y prend la mauvaise fin
      (F♯maj7 mesure 32, puis F9sus4 vers B♭maj7).
- [ ] les 80 désaccords « autre » du corpus, jamais triés.
- [ ] l'ouverture sur un maj7 qui n'est pas la tonique (*Only Trust
      Your Heart*) : ni le ♯11 ni le rythme harmonique ne tranchent
      (voir « Ce qui est difficile », section 5 de `grilles.md`).
- [ ] la modalité sur de vraies grilles modales : ce que devient un
      accord dans une plage, la modulation parallèle (*On Green Dolphin
      Street*), la lecture du mode quand les tétrades suffisent
      (*Nardis*, mi phrygien que l'analyse entend en do). Les grilles
      propres d'*Infant Eyes* et de *Naima* serviront de référence.
- [ ] les cadences modales à deux accords du tome 2 (p. 43), II-I et
      VII-I pour chaque mode ; reste à les distinguer d'un mouvement
      conjoint de grille tonale.
- [ ] les accords parallèles (*Stolen Moments*) : ce qu'on en montre
      quand les blocs lisent déjà les accords.
- [ ] le turnaround, reconnu à sa place dans la structure.
- [ ] le rythme harmonique : mesurer « posé » au pas harmonique du
      passage. Deux essais n'ont pas fait mieux que le seuil actuel.
- [ ] la mémoire, dernier critère de Siron pour la modulation vraie
      (*Love Me Or Leave Me*, voir le tableau dans `grilles.md`).
- [ ] la modulation « confirmée », et nos seuils en données (zone,
      centre éloigné, modulation vraie, appui), à régler sur le corpus.
- [ ] la forme donnée plutôt que détectée, quand elle est connue.
- [ ] la fiche *Django* (Siron 5.11.4), faute de savoir lire son niveau
      de modulation.
- [ ] la mélodie, qui sépare *In a Sentimental Mood* de *Lullaby Of
      Birdland* et lève le doute de *Somewhere* dès ses deux premières
      mesures : demande un format qui la porte.
- [ ] les couleurs proposées à l'apprenant, à rebrancher sur `Sensed`.
- [ ] le niveau de jeu tiré d'une grille : quelles mécaniques d'abord.
- [ ] la jauge de tension, l'attente et la surprise (le pivot diminué de
      *Tenderly*, la pédale de dominante).
- [ ] la grille annotée dans une fenêtre Ebitengine, puis en WASM dans
      une page web.

## En attente d'un cas

Rien n'y est commencé, rien ne presse. Une entrée remonte dans sa
section quand une grille, un jeu ou un besoin la demande.

- une passe sur `Mode.Function`, quand les modes serviront à autre chose
  qu'à l'oreille (seul le double emploi `Tonic | Dominant` des premiers
  degrés harmoniques est verrouillé par un test) ;
- un D.S. dans une reprise pas encore terminée ;
- `irealbook://`, l'ancien schéma non brouillé ;
- un parseur de chiffrages général, pour ce qu'on tape soi-même ;
- afficher les alternatives d'un mode, hors du système naturel ;
- les doigtés ;
- les dominantes sur pédale de tonique, reconnues aux extensions d'un
  m(maj7) ;
- la cadence évitée distinguée de la rompue ;
- un registre parlé anglais pour les modes.

## Les sources à dépouiller

Déjà dépouillés : *En Harmonie* (tome 1, chapitres 8 à 10 ; tome 2,
chapitres 2, 5 et les analyses modales) et Chailley, *40 000 ans de
musique* (sans règle d'harmonie). Les dix commandements de l'harmoniste
de la BEPA, dus à Étienne Guéreau, se citent librement : 

  1. Tu partiras de la mélodie. 
  2. Tu ne suivras pas bêtement les indications du Real Book. 
  3. Tu alterneras les positions et les renversements. 
  4. Tu varieras les nuances et les registres. 
  5. Tu veilleras au toucher et à l'équilibre sonore. 
  6. Tu ne réharmoniseras pas comme un sauvage. 
  7. Tu utiliseras différentes techniques de réharmonisation. 
  8. Tu ne feras pas étalage de ton savoir au détriment du thème. 
  9. Tu seras attentif à la cohérence de l'ensemble. 
  10. Tu transposeras.

- [ ] *En Harmonie*, tome 2, chapitre 6 : l'accord de dominante sur
      tonique et l'accord appoggiaturé.
- [ ] Siron, *La partition intérieure* : le rythme harmonique et la
      carrure, les formes, la mélodie, les blues.
- [ ] Baudoin, *Jazz mode d'emploi* : les blues que `Blues` ne connaît
      pas (*Freddie Freeloader*, *Doxy*, *Watermelon Man*), et les
      réharmonisations.
- [ ] Siskind, *Jazz Piano Fundamentals*, pour *Walk with me*.

## L'audio

- [ ] les allocations de la boucle de jeu, à réduire quand on touche au
      rendu (une collecte toutes les une à deux secondes, sans risque
      pour le son).
- [ ] la calibration chez le joueur, promise dès le premier jour.
- [ ] les timbres 8 bits dans les préférences du joueur (aujourd'hui
      `-timbre` et `-authentic`).
- [ ] la soundfont : GeneralUser GS est retenue (`make sounds`) ;
      restent un meilleur piano, Salamander réduit avec Polyphone par
      exemple, et une mesure sous charge.

Deux portes ouvertes, décrites dans `architecture.md` : un lecteur MIDI
en pur Go et la cible navigateur. Elles ne coûtent que la règle des
deux surfaces (oto dans `synth/device.go`, gomidi dans
`games/keyboard/midi.go`).

## Les jeux

- [ ] *Walk with me* (`walk.md`). Le premier jalon tourne : les
      fondamentales à la main gauche sur un blues jazz, sur la
      batterie, marquées sur la grille, avec la phase sans tempo, le
      bonhomme et une démo qui marche d'après Siskind. Restent :
      - régler à l'oreille les fenêtres, les seuils du bonhomme et le
        gain du kit ;
      - la correction de latence (`oreille.md`) : un premier
        enregistrement (`-record`) montre les notes en avance d'une
        trentaine de millisecondes, régulièrement ;
      - le commentaire du bonhomme en fin de run ;
      - trancher le modèle de paliers avant le deuxième, la basse en
        deux ;
      - d'autres grilles ;
      - un catalogue de patterns pour la basse de référence (voir « La
        basse de référence » dans `walk.md`).
- [ ] `ear`, la suite d'`oreille.md` : la réponse jouée (ce que joue le
      joueur s'allume, la séquence non), les réglages, les niveaux
      paramétrables, les paliers suivants des modes (les autres
      systèmes) et des degrés (tonique mobile, à la basse, d'autres
      gammes, l'échelle chromatique).
- [ ] `ear` en WASM : `midi.go` derrière un build tag, le dex dans
      `localStorage`, l'audio démarré au clic du menu.
- [ ] le shoot'em up, jeu d'harmonie déguisé, sur rail, sans esquive ;
      une grille iReal comme niveau.
- [ ] les quatre pistes du billet sur la game loop de l'improvisateur,
      en jeux distincts : « le canevas, c'est la grille » (Chailley,
      *40 000 ans de musique*, p. 261-263 et 306).

Pistes sur les cibles, pas encore assez nettes pour en faire des
règles ; elles relèvent du jeu, jamais du dex :

- [ ] **la cible fonctionnelle** : un empilement de contraintes
      optionnelles (fonction, tétrade, mode), relatif à une tonalité.
- [ ] **les fondamentales admises pour une fonction** : V7, ♭II7 et
      VIIdim7 pour la dominante, d'après le livre (voir `grilles.md`).
- [ ] **le choix de tétrade comme fait** : jouer D♭7 plutôt que G7
      sous un slot fonctionnel.
- [ ] **les ancres et les détours** d'une réharmonisation, acceptés
      quand ils suivent une progression du catalogue.

## Le dépôt

- [ ] des tags préfixés le jour de la publication : `harmony/v0.1.0`,
      `synth/v0.1.0`.
