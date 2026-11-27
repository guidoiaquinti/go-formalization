# Independent review of `math/rand/v2.NewChaCha8`

**Verdict: FAITHFUL**

Reviewed author commit `56f4df4` against pinned Go 1.27.0 `src/math/rand/v2/chacha8.go:24–28 and internal/chacha8rand/chacha8.go:61–77`. A new stream serializes the seed in little-endian words and starts at position zero with no pending bytes. The first word and complete golden transcript exercise the initial generated block.

Eighteen pinned Go and Lean fixtures passed. The control checks both Uint64 and Read against the SHA-256 transcript published in Go's pinned `math/rand/v2/chacha8_test.go`, then exercises mixed calls and refill boundaries. Starting with a pending byte changed all 18 fixtures. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/math/rand/v2/chacha-stream/review` from the repository root.

Author evidence passed 1,800 boundaries and 60,000 samples at seeds 17, 2718, and 65537. All 32-byte seeds. Generated state and bytes are checked through initialized streams. Concurrency and resource allocation failure are outside scope.
Evidence SHA-256: `b4c1862cf21ea4fadf5f4abd1d82551620cbfe34b20423eec805ebe7e4220715`
