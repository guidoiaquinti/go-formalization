# Independent review of `regexp/syntax.Inst.MatchEmptyWidth`

**Verdict: FAITHFUL**

Evidence SHA-256: `b59381683d484712833e898d3ca836b293ba9cf2315ed29335adbe586d429029`

The matcher uses the low byte of Arg and accepts one recognized flag. Text sentinels are exactly -1. It ignores the instruction opcode and panics for a nil receiver or unknown argument. The pinned Go and Lean replay passed 101,765 fixtures, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Regexp/Syntax/EmptyWidth.lean) · [Comparison](../../../difftest/regexp/syntax/empty-width/check.go)

The independent command `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/regexp/syntax/empty-width/review` matched seven native and Lean fixtures and rejected nine altered outcomes across the group. The cases include negative nonsentinels, newline and word contexts, high Arg bits, synthetic instruction fields, and nil receivers. [Reviewer control](../../../difftest/regexp/syntax/empty-width/review/main.go)
