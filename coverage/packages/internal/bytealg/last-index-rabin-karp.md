# Independent review of `internal/bytealg.LastIndexRabinKarp`

**Verdict: FAITHFUL**

Evidence SHA-256: `b55084681980d3df5e2217c417db757408f90c6fbc703e5f475a3b70563066b3`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a4c960e120bf24b664bf3d3a12ec54d1d84c4f62024e119d94725ccbf612cdc7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e44dd1bf6ba5f6cdfbaeb80e591e96323a2b6a4e2125a1f6dd8dc77031b32c77`.

Reviewed 1816c15 against Go 1.27.0. The reverse search computes its initial window at the end and rolls toward the beginning. The model uses reversed sequences and converts the resulting offset back to the original coordinates. Equal hashes still require byte confirmation.

The shared independent control checks ten pinned native and Lean cases, computes hash coefficients at arbitrary precision, and reduces the results modulo 2^32. It exercises differing first and last hits, overlapping matches, an empty separator, longer-separator bounds panic, malformed raw bytes, and both hash-collision directions. Rolling past a collision to a genuine match distinguishes byte confirmation from accepting hash equality alone.

The complete independent replay passed 597 boundaries and 60,000 samples using seeds 17, 2718 and 65537, through both native generic instantiations. The claim covers exact arithmetic and raw byte outcomes, including the bounds panic class. Timing and exact panic wording are outside scope.

Run the shared reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/internal/bytealg/rabin-karp/review/main.go
```
