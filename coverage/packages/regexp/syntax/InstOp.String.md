# Independent review of `regexp/syntax.InstOp.String`

**Verdict: FAITHFUL**

Reviewed author commit `f56c74f` against Go 1.27.0 `src/regexp/syntax/prog.go`. The generated instruction name table starts with InstAlt at zero and ends with InstRuneAnyNotNL at ten. All other byte values return the empty string. The independent fixtures reject decimal fallback on seven rows and truncating the last table name on one.

Pinned Go, Lean and the author oracle agree on every byte. Ten literal reviewer cases passed on pinned Go and Lean. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/regexp/syntax/opcode-strings/review` after building `godiff-regexp-syntax-opcode-strings`.

Evidence SHA-256: `b050feeb4a9041e929119846cf37c6b9e5902b7798fbfb665643fbc3e92cfde1`
