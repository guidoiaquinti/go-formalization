# Independent review of `regexp/syntax.Inst.op`

**Verdict: FAITHFUL**

Evidence SHA-256: `ff397601aaff27e4597055a3638d1064030461b0dba9b64e2b59aca3b7c3550d`

The private normalizer converts exactly the specialized rune opcodes 8, 9 and 10 to the ordinary rune opcode 7. It leaves unrecognized values unchanged and a nil receiver panics. The pinned Go and Lean replay passed 101,765 fixtures, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Regexp/Syntax/EmptyWidth.lean) · [Comparison](../../../difftest/regexp/syntax/empty-width/check.go)

The independent command `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/regexp/syntax/empty-width/review` matched seven native and Lean fixtures and rejected nine altered outcomes across the group. The cases include negative nonsentinels, newline and word contexts, high Arg bits, synthetic instruction fields, and nil receivers. [Reviewer control](../../../difftest/regexp/syntax/empty-width/review/main.go)
