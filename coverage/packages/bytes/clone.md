# Independent review of `bytes.Clone`

**Verdict: FAITHFUL**

Evidence SHA-256: `a91fa715e7a8fb1b2c7ca2ae29e637feb473d8ff55ad018777b3643bf7cb35b3`

Reviewed author commit `01c772c` against pinned Go 1.27.0 `src/bytes/bytes.go:1391–1396`. A nil slice remains nil. A nonnil empty slice remains nonnil and has zero capacity. Populated results retain the visible bytes and use fresh backing. The native isolation probe verifies changes in each direction and checks the source spare backing.

Six literal fixtures passed against native Go and compiled Lean, alongside three native alias probes. Collapsing nonnil empty input into nil changed one fixture. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/bytes/copy-replace/review` from the repository root.

The author checked 25,653 boundaries and 60,000 sampled cases at seeds 17, 2718, and 65537. Exact capacity claims apply to zero-capacity results and replacement allocation. Allocator-dependent append capacity and allocation failure remain outside scope.
