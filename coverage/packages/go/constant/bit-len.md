# Independent review of `go/constant.BitLen`

**Verdict: FAITHFUL**

Reviewed author commit `56599df` against Go 1.27.0 `src/go/constant/value.go:667–680`. Counts magnitude bits, including 64 for MinInt64. Unknown has length zero and the Boolean case panics.

Nine independent literal fixtures passed against the pinned native declarations and compiled Lean. Counting an extra sign bit changed three fixtures. The control runs all eight altered rules and checks the exact rejection counts. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/go/constant/integer/review` from the repository root.

Author evidence passed 3,024 boundaries and 60,000 sampled rows at seeds 17, 2718, and 65537. Scope covers canonical integers, Unknown, and concrete Boolean wrong-kind checks. Inexact low-word results record the pinned implementation. Incidental allocation capacity and allocation failure remain outside scope.
Evidence SHA-256: `6d3297763509130efac70c98f9e56b73fe03a310d2bf46ac937b15ff80f424b1`
