# Independent review of `slices.stableCmpFunc`

**Verdict: FAITHFUL**

Evidence SHA-256: `7ab9e3813e702d572cfe7d83ca49db6d97739e3dc3f4dccd3d60d19e16c49044`

The comparator helper preserved equal-key order and retained a completed insertion swap when the next callback panicked. A nil callback required no invocation for a singleton and panicked on a two-element input. All 62,162 author cases for this declaration matched pinned native Go and Lean, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Slices/Stable.lean) · [Comparison](../../../difftest/slices/stable/check.go)

Six paired native and Lean controls rejected six altered outcomes for stability, prior mutations on panic and nil callback timing. The registered scope uses valid bounded ranges and signed values with comparison-invisible tags for stable identity tests. Arbitrary callback side effects are outside scope. [Reviewer control](../../../difftest/slices/stable/review/main.go)
