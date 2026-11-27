# Independent review of `regexp/syntax.Op.String`

**Verdict: FAITHFUL**

Reviewed author commit `f56c74f` against Go 1.27.0 `src/regexp/syntax/op_string.go`. The contiguous public range runs from one through nineteen. The separate pseudo value is 128, and unsigned decimal covers the remaining bytes. The fixtures reject omitting the pseudo spelling on one row and signed decimal on two.

Pinned Go, Lean and the author oracle agree on every byte. Ten literal reviewer cases passed on pinned Go and Lean. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/regexp/syntax/opcode-strings/review` after building `godiff-regexp-syntax-opcode-strings`.

Evidence SHA-256: `17436278c4ccb3cc2e7b811f786c2285f76455de9ee9a09729f88ab1002f1864`
