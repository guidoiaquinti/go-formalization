# Independent review of `encoding/binary.ensure`

**Verdict: FAITHFUL**

Evidence SHA-256: `0cc08c73669f7bb6fa1941c4e77e1091ca7cfd71f34d34fdfcafe36c26bdb166`

Go 1.27.0 grows the input when necessary, extends the resulting slice by the requested count, and returns a suffix beginning at the old length. Existing spare bytes become visible before zero-filled storage, including when the input detaches. Input and returned suffix nilness coincide on zero growth. A negative count panics before modifying storage. The author replay passed all 66,103 cases, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Encoding/Binary/Ensure.lean) · [Comparison](../../../difftest/encoding/binary/ensure/check.go)

Nine independent Lean rows passed. Native offset views confirmed two-way mutation between returned slices in both reused and detached storage, proper original backing aliasing, and unchanged surrounding sentinel bytes. Clearing newly exposed cells failed four native rows, omitting existing spare bytes on allocation failed two, and copying the returned suffix failed four. The model records suffix contents and whether allocation is required. The native control checks alias identity. Allocator capacity and allocation failure are outside scope. [Reviewer control](../../../difftest/encoding/binary/ensure/review/main.go)
