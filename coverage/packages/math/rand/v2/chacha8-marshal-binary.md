# Independent review of `math/rand/v2.ChaCha8.MarshalBinary`

**Verdict: FAITHFUL**

Reviewed author commit `56f4df4` against pinned Go 1.27.0 `src/math/rand/v2/chacha8.go:112–115 and internal/chacha8rand/chacha8.go:133–143`. The payload records the current seed and big-endian draw position, with a length-prefixed pending buffer when present. The fixtures decode both fresh and mixed-call states independently.

Eighteen pinned Go and Lean fixtures passed. The control checks both Uint64 and Read against the SHA-256 transcript published in Go's pinned `math/rand/v2/chacha8_test.go`, then exercises mixed calls and refill boundaries. Writing the draw position little-endian changed all eighteen fixtures. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/math/rand/v2/chacha-stream/review` from the repository root.

Author evidence passed 1,800 boundaries and 60,000 samples at seeds 17, 2718, and 65537. Nonnil receivers initialized by NewChaCha8 or Seed and reached by sequential Read and Uint64 calls. A zero-value receiver before explicit initialization, states restored by UnmarshalBinary, allocation capacities after growth, and concurrency are outside scope.
Evidence SHA-256: `d4f9c463b4872cfd4f97bba04bc07c0224d8be2efa2f9cfb48dcf0f04b79db2e`
