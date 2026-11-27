.PHONY: verify build proof diff contracts coverage coverage-update coverage-check inventory lake-update

PINNED_GO = GOWORK=off GOTOOLCHAIN=go$(shell cat .go-version) go
# Build repository tools for the host even in a cross-compilation shell.
NATIVE_GO = GOOS= GOARCH= GOENV=off GOEXPERIMENT= GOFLAGS= CGO_ENABLED=0 GO386= GOAMD64= GOARM= GOARM64= GOMIPS= GOMIPS64= GOPPC64= GORISCV64= GOWASM= $(PINNED_GO)
COVERAGE = $(NATIVE_GO) run ./scripts/coverage
TARGET_FLAG = $(if $(TARGET),--target "$(TARGET)")

verify: coverage-check proof $(if $(filter-out math,$(TARGET)),,contracts)
	@$(NATIVE_GO) test ./scripts/coverage
	@$(COVERAGE) diff $(TARGET_FLAG)

build:
	cd lean && lake build

proof: build
	cd lean && lake env lean Audit.lean

diff:
	@$(COVERAGE) diff-build $(TARGET_FLAG)

contracts:
	@$(NATIVE_GO) run ./difftest/math/contracts/check.go -replay
	@$(NATIVE_GO) run -gcflags=command-line-arguments=-d=ssa/intrinsics/off -gcflags=math=-d=ssa/intrinsics/off ./difftest/math/contracts/check.go -replay

lake-update:
	@$(COVERAGE) lake-update

coverage:
	@$(COVERAGE) report

coverage-update:
	@$(COVERAGE) update

coverage-check:
	@$(COVERAGE) check

inventory:
	@$(COVERAGE) inventory $(TARGET_FLAG) $(if $(filter 1,$(EXTERNAL_BLOCKERS)),--external-blockers)
