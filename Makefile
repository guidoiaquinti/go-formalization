.PHONY: verify build proof diff coverage coverage-update coverage-check inventory

PINNED_GO = GOWORK=off GOTOOLCHAIN=go$(shell cat .go-version) go
# Build repository tools for the host even in a cross-compilation shell.
NATIVE_GO = GOOS= GOARCH= GOENV=off GOEXPERIMENT= GOFLAGS= CGO_ENABLED=0 GO386= GOAMD64= GOARM= GOARM64= GOMIPS= GOMIPS64= GOPPC64= GORISCV64= GOWASM= $(PINNED_GO)
COVERAGE = $(NATIVE_GO) run ./scripts/coverage

verify: coverage-check proof
	@$(NATIVE_GO) test ./scripts/coverage
	@$(COVERAGE) diff

build:
	cd lean && lake build

proof: build
	cd lean && lake env lean Audit.lean

diff: build
	@$(COVERAGE) diff $(if $(TARGET),--target "$(TARGET)")

coverage:
	@$(COVERAGE) report

coverage-update:
	@$(COVERAGE) update

coverage-check:
	@$(COVERAGE) check

inventory:
	@$(COVERAGE) inventory $(if $(TARGET),--target "$(TARGET)")
