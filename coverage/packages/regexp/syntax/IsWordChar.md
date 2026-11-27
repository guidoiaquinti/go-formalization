# Independent review of `regexp/syntax.IsWordChar`

**Verdict: FAITHFUL**

Evidence SHA-256: `3a91d224cbcd47912f6fba111a0aab043eaeb208ef00d0e117f87bd26ba02aa9`

The entry recognizes only ASCII letters, digits and underscore. The independent controls reject omitting underscore and allowing the Unicode Kelvin sign. The pinned Go and Lean replay passed 101,765 fixtures, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Regexp/Syntax/EmptyWidth.lean) · [Comparison](../../../difftest/regexp/syntax/empty-width/check.go)

The independent command `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/regexp/syntax/empty-width/review` matched seven native and Lean fixtures and rejected nine altered outcomes across the group. The cases include negative nonsentinels, newline and word contexts, high Arg bits, synthetic instruction fields, and nil receivers. [Reviewer control](../../../difftest/regexp/syntax/empty-width/review/main.go)
