# Independent review of `go/constant.makeInt`

**Verdict: FAITHFUL**

Reviewed author commit `56599df` against Go 1.27.0 `src/go/constant/value.go:285–290`. Uses the compact representation when signed 64-bit range permits. Larger positive and negative magnitudes remain arbitrary precision, and the native observer confirms retention of the original big.Int for large values.

Nine independent literal fixtures passed against the pinned native declarations and compiled Lean. Forcing the signed low word into compact storage changed three fixtures. The control runs all eight altered rules and checks the exact rejection counts. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/go/constant/integer/review` from the repository root.

Author evidence passed 3,024 boundaries and 60,000 sampled rows at seeds 17, 2718, and 65537. Scope covers canonical integers, Unknown, and concrete Boolean wrong-kind checks. Inexact low-word results record the pinned implementation. Incidental allocation capacity and allocation failure remain outside scope.
Evidence SHA-256: `03f68d2e43dfb4aede971e84068c8078034fca2343991732bf062e1d8d5cfd76`
