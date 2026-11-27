# Independent review of `go/constant.Bytes`

**Verdict: FAITHFUL**

Reviewed author commit `56599df` against Go 1.27.0 `src/go/constant/value.go:724–750`. Emits the shortest little-endian byte representation of the magnitude, with empty bytes for zero. Unknown and Boolean inputs panic.

Nine independent literal fixtures passed against the pinned native declarations and compiled Lean. Using big-endian output changed five fixtures. The control runs all eight altered rules and checks the exact rejection counts. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/go/constant/integer/review` from the repository root.

Author evidence passed 3,024 boundaries and 60,000 sampled rows at seeds 17, 2718, and 65537. Scope covers canonical integers, Unknown, and concrete Boolean wrong-kind checks. Inexact low-word results record the pinned implementation. Incidental allocation capacity and allocation failure remain outside scope.
Evidence SHA-256: `fac1122f6d218b9204a27854deb6cdf97ffd6be51df9c6cfa6ab00dec5e2d57d`
