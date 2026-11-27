# Independent review of `slices.SortedFunc`

**Verdict: FAITHFUL**

Evidence SHA-256: `29b9d81106a7bcdf473278f95abece2ad61aab700edd17bbc5b85bf998a9b63d`

Exact producer panic identity was preserved after three acknowledged yields, with no comparison. A successful native comparison verified the iterator had already returned. Independent model rows verified comparator panic after collection and no published result. These controls reject sorting while producing and skipping producer failures. Pinned native Go and Lean agreed on 66,563 author cases for this declaration, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Slices/PublicSort.lean) · [Comparison](../../../difftest/slices/public-sort/check.go)

The reviewer replay adds four hand-computed Lean rows and separate direct pinned Go controls. The scope covers finite int64 and binary64 inputs, scripted comparisons and finite iterators. Arbitrary callback mutation, reentrancy and exact allocation capacity are outside scope. [Reviewer control](../../../difftest/slices/public-sort/review/main.go)
