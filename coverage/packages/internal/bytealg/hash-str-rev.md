# Independent review of `internal/bytealg.hashStrRev`

**Verdict: FAITHFUL**

Evidence SHA-256: `1b8f4454b00171cc06b79db34d825b1191f72ef5dfc0fefaa20b4d2a4379506d`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1aca9226c166804b1cf1083b2806762077c54a949807e2e8ee578a9fc7169b48`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `79685a400d7ac7c9fbe24cad5e7339ce714ab5f65a973e6f45ebf48d57d9b985`.

Reviewed 1816c15 against Go 1.27.0. The reverse hash uses the same byte recurrence and wrapped power, traversing the input from its final byte to its first.

The shared independent control checks ten pinned native and Lean cases, computes hash coefficients at arbitrary precision, and reduces the results modulo 2^32. It exercises differing first and last hits, overlapping matches, an empty separator, longer-separator bounds panic, malformed raw bytes, and both hash-collision directions. Rolling past a collision to a genuine match distinguishes byte confirmation from accepting hash equality alone.

The complete independent replay passed 597 boundaries and 60,000 samples using seeds 17, 2718 and 65537, through both native generic instantiations. The claim covers exact arithmetic and raw byte outcomes, including the bounds panic class. Timing and exact panic wording are outside scope.

Run the shared reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/internal/bytealg/rabin-karp/review/main.go
```
