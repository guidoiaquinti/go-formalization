# Independent review of `regexp/syntax.EmptyOpContext`

**Verdict: FAITHFUL**

Evidence SHA-256: `37b94f1a1c023a7cffd6de81c7e43cf8ea71fb41fbabdcc98d24fd78f5a8e57e`

Any negative rune contributes its beginning or ending text and line context, whereas newline only contributes line context. It adds exactly one of word boundary or no word boundary. The pinned Go and Lean replay passed 101,765 fixtures, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Regexp/Syntax/EmptyWidth.lean) · [Comparison](../../../difftest/regexp/syntax/empty-width/check.go)

The independent command `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/regexp/syntax/empty-width/review` matched seven native and Lean fixtures and rejected nine altered outcomes across the group. The cases include negative nonsentinels, newline and word contexts, high Arg bits, synthetic instruction fields, and nil receivers. [Reviewer control](../../../difftest/regexp/syntax/empty-width/review/main.go)
