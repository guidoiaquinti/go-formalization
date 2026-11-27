# Independent review of `slices.symMergeCmpFunc`

**Verdict: FAITHFUL**

Evidence SHA-256: `c8b45aa3f6f12c02648799c10b5d77078035d307a0dcc7317974fc42f252df3a`

The comparator helper retained the relative order of equal-key tagged elements. A panic during the first recursive merge preserved the prior central rotation and stopped the exact callback trace. All 71,666 author cases for this declaration matched pinned native Go and Lean, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Slices/Stable.lean) · [Comparison](../../../difftest/slices/stable/check.go)

Six paired native and Lean controls rejected six altered outcomes for stability, prior mutations on panic and nil callback timing. The registered scope uses valid bounded ranges and signed values with comparison-invisible tags for stable identity tests. Arbitrary callback side effects are outside scope. [Reviewer control](../../../difftest/slices/stable/review/main.go)
