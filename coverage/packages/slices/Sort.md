# Independent review of `slices.Sort`

**Verdict: FAITHFUL**

Evidence SHA-256: `4e0a054552a89422fae68386dd5a2c5b573d150dbaf71eeff57f418287be3e67`

The direct binary64 control checked exact signaling and quiet NaN payload order, numerical order, and distinct signed zero encodings. It rejects dropping, canonicalizing or sorting NaN words by unsigned bit pattern. Pinned native Go and Lean agreed on 66,563 author cases for this declaration, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Slices/PublicSort.lean) · [Comparison](../../../difftest/slices/public-sort/check.go)

The reviewer replay adds four hand-computed Lean rows and separate direct pinned Go controls. The scope covers finite int64 and binary64 inputs, scripted comparisons and finite iterators. Arbitrary callback mutation, reentrancy and exact allocation capacity are outside scope. [Reviewer control](../../../difftest/slices/public-sort/review/main.go)
