# Independent review of `slices.stableOrdered`

**Verdict: FAITHFUL**

Evidence SHA-256: `9a642256fa317d820e78340d1648cba84d5227c763737049ec9e7b2bf1483725`

The ordered helper sorted the reported prefix using the original signed values. A suffix sentinel remained unchanged, and fixed replay distinguished natural numeric comparison from grouping by the tag key. All 62,162 author cases for this declaration matched pinned native Go and Lean, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Slices/Stable.lean) · [Comparison](../../../difftest/slices/stable/check.go)

Six paired native and Lean controls rejected six altered outcomes for stability, prior mutations on panic and nil callback timing. The registered scope uses valid bounded ranges and signed values with comparison-invisible tags for stable identity tests. Arbitrary callback side effects are outside scope. [Reviewer control](../../../difftest/slices/stable/review/main.go)
