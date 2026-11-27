# Independent review of `internal/bytealg.IndexRabinKarp`

**Verdict: FAITHFUL**

Evidence SHA-256: `17e0c0a3d12bd726ba611339c8dee7ed24f077942380d403012dbd267717672c`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0fd301d29af2c90f49f0950c53b897b9297cf27ad3b9215f28cbcb77913a4265`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6b21bf921a4fefdc5dd948971185cc713fab477845029f93816f983b32beb7f4`.

Reviewed 1816c15 against Go 1.27.0. The forward search computes an initial raw-byte window, then rolls using the incoming byte and the weighted outgoing byte. An equal hash triggers exact byte confirmation. The scan returns the first match.

The shared independent control checks ten pinned native and Lean cases, computes hash coefficients at arbitrary precision, and reduces the results modulo 2^32. It exercises differing first and last hits, overlapping matches, an empty separator, longer-separator bounds panic, malformed raw bytes, and both hash-collision directions. Rolling past a collision to a genuine match distinguishes byte confirmation from accepting hash equality alone.

The complete independent replay passed 597 boundaries and 60,000 samples using seeds 17, 2718 and 65537, through both native generic instantiations. The claim covers exact arithmetic and raw byte outcomes, including the bounds panic class. Timing and exact panic wording are outside scope.

Run the shared reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/internal/bytealg/rabin-karp/review/main.go
```
