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

# La latence audio ne se mesure pas en test : elle se joue et s'écoute.
.PHONY: latency
latency:
	cd games && go run ./latency
