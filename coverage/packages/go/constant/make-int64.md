# Independent review of `go/constant.MakeInt64`

**Verdict: FAITHFUL**

Reviewed author commit `56599df` against Go 1.27.0 `src/go/constant/value.go:397`. Every signed word is stored as the compact int64 representation, including negative inputs and MinInt64.

Nine independent literal fixtures passed against the pinned native declarations and compiled Lean. Treating a negative signed word as large unsigned changed four fixtures. The control runs all eight altered rules and checks the exact rejection counts. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/go/constant/integer/review` from the repository root.

Author evidence passed 3,024 boundaries and 60,000 sampled rows at seeds 17, 2718, and 65537. Scope covers canonical integers, Unknown, and concrete Boolean wrong-kind checks. Inexact low-word results record the pinned implementation. Incidental allocation capacity and allocation failure remain outside scope.
Evidence SHA-256: `d88d0605cb4c48092290cae0908df4eee887c97de6ec693bb8a4645f4fe5b9c4`
