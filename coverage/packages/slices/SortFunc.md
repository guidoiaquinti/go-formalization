# Independent review of `slices.SortFunc`

**Verdict: FAITHFUL**

Evidence SHA-256: `518d1ba7ea31a560c44e6c483ace94c95555accffca4e9267cc928516ca79acc`

A comparator panic on the second comparison retained the first swap, kept the panic object, and stopped the ordered callback trace. A nil comparator succeeded for a singleton and panicked on a pair. This rejects rollback, eager nil checking, and extra callbacks. Pinned native Go and Lean agreed on 66,563 author cases for this declaration, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Slices/PublicSort.lean) · [Comparison](../../../difftest/slices/public-sort/check.go)

The reviewer replay adds four hand-computed Lean rows and separate direct pinned Go controls. The scope covers finite int64 and binary64 inputs, scripted comparisons and finite iterators. Arbitrary callback mutation, reentrancy and exact allocation capacity are outside scope. [Reviewer control](../../../difftest/slices/public-sort/review/main.go)
