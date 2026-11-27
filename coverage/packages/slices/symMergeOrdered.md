# Independent review of `slices.symMergeOrdered`

**Verdict: FAITHFUL**

Evidence SHA-256: `c9db45ef1be136bd78fd26a2a2e7db6b644faf9089b4adedba90406b72e3c846`

The ordered helper merged two nonempty sorted runs while keeping elements outside the chosen bounds intact. The fixed fixture uses different raw and comparator order, so calling the comparator variant in its place produces a different result. All 71,666 author cases for this declaration matched pinned native Go and Lean, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Slices/Stable.lean) · [Comparison](../../../difftest/slices/stable/check.go)

Six paired native and Lean controls rejected six altered outcomes for stability, prior mutations on panic and nil callback timing. The registered scope uses valid bounded ranges and signed values with comparison-invisible tags for stable identity tests. Arbitrary callback side effects are outside scope. [Reviewer control](../../../difftest/slices/stable/review/main.go)
