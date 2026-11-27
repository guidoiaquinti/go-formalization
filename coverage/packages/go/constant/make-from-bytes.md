# Independent review of `go/constant.MakeFromBytes`

**Verdict: FAITHFUL**

Reviewed author commit `56599df` against Go 1.27.0 `src/go/constant/value.go:756–780`. Interprets raw bytes as an unsigned little-endian magnitude and then compacts when possible. The native differential mutates the input after construction to verify ownership.

Nine independent literal fixtures passed against the pinned native declarations and compiled Lean. Using big-endian input changed four fixtures. The control runs all eight altered rules and checks the exact rejection counts. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/go/constant/integer/review` from the repository root.

Author evidence passed 3,024 boundaries and 60,000 sampled rows at seeds 17, 2718, and 65537. Scope covers canonical integers, Unknown, and concrete Boolean wrong-kind checks. Inexact low-word results record the pinned implementation. Incidental allocation capacity and allocation failure remain outside scope.
Evidence SHA-256: `273814a85063412c41e1571d8fcface0b7c03d25f52e182b93227a8b4598d310`
