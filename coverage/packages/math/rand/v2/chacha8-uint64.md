# Independent review of `math/rand/v2.ChaCha8.Uint64`

**Verdict: FAITHFUL**

Reviewed author commit `56f4df4` against pinned Go 1.27.0 `src/math/rand/v2/chacha8.go:38–46 and internal/chacha8rand/chacha8.go:49–97`. The scalar model preserves the output word interleaving, four-block transition, and reseed after 124 emitted words. Uint64 advances the generated word index while retaining pending Read bytes.

Eighteen pinned Go and Lean fixtures passed. The control checks both Uint64 and Read against the SHA-256 transcript published in Go's pinned `math/rand/v2/chacha8_test.go`, then exercises mixed calls and refill boundaries. Clearing pending bytes changed sixteen fixtures. Emitting the four reserved reseed words changed six. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/math/rand/v2/chacha-stream/review` from the repository root.

Author evidence passed 1,800 boundaries and 60,000 samples at seeds 17, 2718, and 65537. Nonnil receivers initialized by NewChaCha8 or Seed and reached by sequential Read and Uint64 calls. A zero-value receiver before explicit initialization, states restored by UnmarshalBinary, allocation capacities after growth, and concurrency are outside scope.
Evidence SHA-256: `3f4ebb08e10a34367a8d610e8b031b7d0458395578960f09f5f5cbc11b8b9306`
