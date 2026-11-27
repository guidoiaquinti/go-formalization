# Independent review of `bytes.ReplaceAll`

**Verdict: FAITHFUL**

Evidence SHA-256: `dea8a2bea525a7bb326d37e06e41163dae2935bb96d706833d72ad5b6e2c54c4`

Reviewed author commit `01c772c` against pinned Go 1.27.0 `src/bytes/bytes.go:1206–1208`. The wrapper requests unlimited replacement. Its direct native fixtures cover multi-match replacements, empty-pattern insertion inside malformed input, a complete encoded replacement rune, and exact allocation capacity. The native alias probe uses source subslices for both patterns and replacements and checks isolation in both directions.

Six literal fixtures passed against native Go and compiled Lean, alongside three native alias probes. Delegating with limit one changed three fixtures. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/bytes/copy-replace/review` from the repository root.

The author checked 25,653 boundaries and 60,000 sampled cases at seeds 17, 2718, and 65537. Exact capacity claims apply to zero-capacity results and replacement allocation. Allocator-dependent append capacity and allocation failure remain outside scope.
