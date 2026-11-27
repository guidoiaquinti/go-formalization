# Independent review of `encoding/binary.dataSize`

**Verdict: FAITHFUL**

Evidence SHA-256: `f04b368000079aebd3148265e3316f77a408edeb26873e342bdb814d5d948cdf`

Private dataSize accepts reflected arrays and slices by combining their length with the encoded element width. Unsupported element shapes fail before multiplying, even at length zero. Both supported and unsupported struct layouts matched with the native cache cleared and immediately warmed. Invalid reflect.Value returns minus one. The pinned Go, Lean and layout-oracle replay matched all 61,814 layouts, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Encoding/Binary/ReflectSize.lean) · [Comparison](../../../difftest/encoding/binary/reflect-size/check.go)

Ten independent native and Lean layouts passed. The control rejected using the native padded width, accepting unsupported element types at zero length, recursively dereferencing pointers, and ignoring the permitted top-level indirection. The model describes returned size or panic and leaves cache timing and concurrency outside scope. Encoded widths must fit the native int. [Reviewer control](../../../difftest/encoding/binary/reflect-size/review/main.go)
