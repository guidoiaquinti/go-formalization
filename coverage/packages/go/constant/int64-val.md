# Independent review of `go/constant.Int64Val`

**Verdict: FAITHFUL**

Reviewed author commit `56599df` against Go 1.27.0 `src/go/constant/value.go:505–515`. Returns the exact compact value with true, observes the pinned signed low word with false for large values, and yields zero and false for Unknown. The fixture also checks a Boolean panic.

Nine independent literal fixtures passed against the pinned native declarations and compiled Lean. Returning zero for inexact integer results changed three fixtures. The control runs all eight altered rules and checks the exact rejection counts. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/go/constant/integer/review` from the repository root.

Author evidence passed 3,024 boundaries and 60,000 sampled rows at seeds 17, 2718, and 65537. Scope covers canonical integers, Unknown, and concrete Boolean wrong-kind checks. Inexact low-word results record the pinned implementation. Incidental allocation capacity and allocation failure remain outside scope.
Evidence SHA-256: `ea91fa38d75d5af1a252e5f8346efd98f7417cb9508a3df423808f503ddb6d44`
