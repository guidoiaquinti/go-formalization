# Independent review of `math/rand/v2.ChaCha8.Read`

**Verdict: FAITHFUL**

Reviewed author commit `56f4df4` against pinned Go 1.27.0 `src/math/rand/v2/chacha8.go:54–71`. Read drains pending bytes first, uses little-endian output, and saves the unused suffix of a partial word. Native Go returns the requested length and nil for every reviewed call.

Eighteen pinned Go and Lean fixtures passed. The control checks both Uint64 and Read against the SHA-256 transcript published in Go's pinned `math/rand/v2/chacha8_test.go`, then exercises mixed calls and refill boundaries. Ignoring pending bytes changed sixteen fixtures. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/math/rand/v2/chacha-stream/review` from the repository root.

Author evidence passed 1,800 boundaries and 60,000 samples at seeds 17, 2718, and 65537. Nonnil receivers initialized by NewChaCha8 or Seed and reached by sequential Read and Uint64 calls. A zero-value receiver before explicit initialization, states restored by UnmarshalBinary, allocation capacities after growth, and concurrency are outside scope.
Evidence SHA-256: `e82d056ed4142888335e04513375e2708ec6646e75635bc6d72928dac71ff4c0`
