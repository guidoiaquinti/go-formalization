# Working on go-formalization

Follow the contributor guide for source selection, file layout, test thresholds, registry fields and review format. [Contributor guide](CONTRIBUTING.md)

- Use Go for repository tooling, differential checks and their tests. [Coverage CLI](scripts/coverage/)
- Before editing, read the Go and Lean pins, run `make inventory TARGET=<requested-target>` (for example, `make inventory TARGET=math/bits`), and inspect the existing models and evidence. [Go version](.go-version) · [Lean toolchain](lean/lean-toolchain)
- Run `make coverage-update` and `make verify` from the repository root. Report the outcome, tested domain or case count, seeds where applicable, final inventory state and any missing verification.
- Leave passing evidence at `AWAITING_REVIEW` until a separate reviewer records a verdict and fingerprint. Keep repository review notes with the relevant evidence. Reserve [FINDINGS.md](FINDINGS.md) for confirmed or potential Go bugs and label each entry by status for human triage.
