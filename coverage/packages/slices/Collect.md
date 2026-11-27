# Independent review of `slices.Collect`

**Verdict: FAITHFUL**

Evidence SHA-256: `56b170aff0502376d1fdf5e0d49e9a446ecf69c24e3eaa8418d9c3c9833e23a1`

Reviewed 4c5474f against Go 1.27.0. Collect invokes AppendSeq with a nil initial slice. An empty producer leaves that slice nil. Nonempty input preserves the exact order and yields are acknowledged before a later producer panic propagates.

The shared control passes five explicit model rows and native checks on named string slices. It tests nil and nonnil emptiness, backing reuse versus growth, acknowledgment, exact panic identity, and original backing after panic. Deferring all writes until producer completion would leave the old spare untouched in the panic fixture.

The complete independent replay passed 12 boundaries and 60,000 samples with seeds 17, 2718 and 65537. Exact capacity of newly allocated storage and dynamic reentrant producer behavior remain outside scope.

Run the shared reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/slices/append-collect/review/main.go
```
