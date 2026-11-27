# Independent review of `math/rand/v2.ChaCha8.Seed`

**Verdict: FAITHFUL**

Reviewed author commit `56f4df4` against pinned Go 1.27.0 `src/math/rand/v2/chacha8.go:31–35`. Seed replaces the generated state, resets the stream position and clears the read buffer. The control resets mixed streams with pending bytes to the published test seed, also from a zero-valued receiver.

Eighteen pinned Go and Lean fixtures passed. The control checks both Uint64 and Read against the SHA-256 transcript published in Go's pinned `math/rand/v2/chacha8_test.go`, then exercises mixed calls and refill boundaries. Retaining the old pending bytes changed eight fixtures. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/math/rand/v2/chacha-stream/review` from the repository root.

Author evidence passed 1,800 boundaries and 60,000 samples at seeds 17, 2718, and 65537. Nonnil receivers, with differential replay on seeded reachable states and a native zero-value reset probe. Concurrency and resource allocation failure are outside scope.
Evidence SHA-256: `ff850b7d183df7c5e65d8fd21eb7aa962ccc47f103ba9add08c5e3e83706f811`
