# Independent review of `internal/strconv.lower`

**Verdict: FAITHFUL**

Evidence SHA-256: `9d0006e6d1e96a8eb0a8815a604cbbd06e2149db6ef91e30a6c344b59a6ad948`

The helper unconditionally sets bit five and preserves every other bit. It also changes punctuation, NUL, and high bytes where that bit is clear. ASCII letter-only folding and toggling bit five fail the independent observations. All 256 byte values matched pinned native Go and Lean. [Model](../../../../lean/GoLean/Strconv/IntegerBits.lean) · [Comparison](../../../../difftest/internal/strconv/integer-bits/check.go)

The independent command `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/internal/strconv/integer-bits/review` passed 15 Lean observations and ten pinned native observations, rejecting 15 altered results across the group. [Reviewer control](../../../../difftest/internal/strconv/integer-bits/review/main.go)
