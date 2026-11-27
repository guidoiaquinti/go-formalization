# Independent review of `internal/strconv.isPowerOfTwo`

**Verdict: FAITHFUL**

Evidence SHA-256: `90e5a6b9d337e21aceeb6a2793ca0b8dc5d9e6de384cc4f9d7094e01d6179d56`

The helper recognizes zero and native words with a single set bit. Minimum signed int therefore succeeds, while its immediate neighbors fail. Positivity checks and a sign-bit shortcut fail the independent fixture. The author replay matched 60,258 native 64-bit words and 60,130 additional modeled 32-bit words at seeds 17, 2718 and 65537. Native 32-bit execution is outside this review. [Model](../../../../lean/GoLean/Strconv/IntegerBits.lean) · [Comparison](../../../../difftest/internal/strconv/integer-bits/check.go)

The independent command `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/internal/strconv/integer-bits/review` passed 15 Lean observations and ten pinned native observations, rejecting 15 altered results across the group. [Reviewer control](../../../../difftest/internal/strconv/integer-bits/review/main.go)
