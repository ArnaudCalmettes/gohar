# go test ./... ne traverse pas les frontières de modules, d'où ce
# fichier plutôt qu'une commande à retenir.

MODULES := harmony dex synth games

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

# La latence audio ne se mesure pas en test : elle se joue et s'écoute.
.PHONY: latency
latency:
	cd games && go run ./latency
