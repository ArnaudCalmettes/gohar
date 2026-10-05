# go test ./... ne traverse pas les frontières de modules, d'où ce
# fichier plutôt qu'une commande à retenir.

MODULES := harmony dex synth games charts

.PHONY: test vet fmt bench all

all: fmt vet test

test:
	@for m in $(MODULES); do echo "== $$m"; (cd $$m && go test ./...); done

vet:
	@for m in $(MODULES); do (cd $$m && go vet ./...); done

fmt:
	@for m in $(MODULES); do (cd $$m && gofmt -l -w .); done

bench:
	cd harmony && go test -bench . -benchmem ./...

# Les sons trop gros pour le dépôt, téléchargés une fois dans le cache
# de l'utilisateur, là où les jeux les cherchent (os.UserCacheDir sous
# Linux). Le commit est figé et l'empreinte vérifiée : un fichier qui
# change en amont ne passe pas en silence.
CACHE ?= $(or $(XDG_CACHE_HOME),$(HOME)/.cache)/gohar
GENERALUSER_URL := https://raw.githubusercontent.com/mrbumpy409/GeneralUser-GS/684543d5e5efaef08d02be50dcda8d552478fa60/GeneralUser-GS.sf2
GENERALUSER_SHA256 := 9575028c7a1f589f5770fccc8cff2734566af40cd26ed836944e9a5152688cfe

.PHONY: sounds
sounds: $(CACHE)/GeneralUser-GS.sf2

$(CACHE)/GeneralUser-GS.sf2:
	@mkdir -p $(CACHE)
	curl -fL --progress-bar -o $@.part $(GENERALUSER_URL)
	@echo "$(GENERALUSER_SHA256)  $@.part" | sha256sum -c --quiet
	@mv $@.part $@
	@echo "GeneralUser GS v2.0.3 dans $@"

# La soundfont de Walk with me, réduite à ce qu'il joue : la
# contrebasse, le piano, et quatre touches du kit Jazz (la charleston au
# pied, la ride, les deux wood blocks de la calibration). Elle est
# versionnée et embarquée dans le jeu ; à refaire quand le jeu joue un
# son de plus. Voir synth/sf2 et « Le son » dans docs/walk.md.
WALK_SF2 := games/walk/sounds/walk.sf2

.PHONY: slim
slim: $(CACHE)/GeneralUser-GS.sf2
	cd synth && go run ./cmd/sfslim -keep 0:32 -keep 0:0 -keep 128:32/44,51,76,77 -o ../$(WALK_SF2) $<

# Walk with me dans le navigateur, dans build/web : le programme en
# WebAssembly, sans ses tables de débogage (un demi-Mo de moins), le
# chargeur de Go qui va avec, la page, et le synthé du fil audio,
# FluidSynth en AudioWorklet (js-synthesizer, voir synth/web). make serve
# sert le tout sur http://localhost:8080 ; localhost compte comme un
# contexte sûr, ce que la Web MIDI exige. Voir games/walk/web/index.html,
# et .github/workflows/pages.yml pour la version en ligne.
WEB := build/web

# Les trois scripts de js-synthesizer, du paquet npm, version figée,
# empreintes vérifiées, téléchargés une fois.
JSSYNTH_URL := https://cdn.jsdelivr.net/npm/js-synthesizer@1.11.0
JSSYNTH := \
	externals/libfluidsynth-2.4.6.js:42ece3b53ef3289a3b30e0e1ebb4951e9da9987ae6ca975c57a0aabc91237f66 \
	dist/js-synthesizer.worklet.min.js:93ba8c10322273523fa7de8c4b91d402692db8f24a52d7ca3bd31b99b5d06510 \
	dist/js-synthesizer.min.js:296f3cf76304b3650d7fb24a00263e2225e21d612fc0882234de7b1090da2b53

.PHONY: wasm serve jssynth latency-web
jssynth:
	@mkdir -p $(WEB)
	@for f in $(JSSYNTH); do \
		path=$${f%%:*}; sum=$${f##*:}; out=$(WEB)/$$(basename $$path); \
		[ -f $$out ] || curl -fsSL -o $$out $(JSSYNTH_URL)/$$path; \
		echo "$$sum  $$out" | sha256sum -c --quiet || exit 1; \
	done

wasm: jssynth
	cd games && GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o ../$(WEB)/walk.wasm ./walk
	cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" $(WEB)/
	cp games/walk/web/index.html $(WEB)/
	@ls -lh $(WEB)/walk.wasm

serve: wasm
	cd $(WEB) && python3 -m http.server 8080

# La sonde de latence du navigateur : le même synthé, joué au clavier
# MIDI avec la soundfont du jeu, sans Go. Ouvrir
# http://localhost:8080/latency.html ; voir games/walk/web/latency.html.
latency-web: jssynth
	cp games/walk/web/latency.html games/walk/sounds/walk.sf2 $(WEB)/
	cd $(WEB) && python3 -m http.server 8080

# La latence audio ne se mesure pas en test : elle se joue et s'écoute.
.PHONY: latency
latency:
	cd games && go run ./latency
