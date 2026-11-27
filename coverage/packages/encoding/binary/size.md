# Independent review of `encoding/binary.Size`

**Verdict: FAITHFUL**

Evidence SHA-256: `40b89881719bb90ba9a89240e28bd33ead60cc4340772b5d836ad14a2531c06e`

Public Size dereferences once at the top level. A live pointer to a slice gets the encoded size of the slice. A pointer to a pointer and a typed nil scalar pointer return minus one. The independent tests ruled out recursive dereference and rejecting every top-level pointer. The pinned Go, Lean and layout-oracle replay matched all 61,814 layouts, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Encoding/Binary/ReflectSize.lean) · [Comparison](../../../difftest/encoding/binary/reflect-size/check.go)

Ten independent native and Lean layouts passed. The control rejected using the native padded width, accepting unsupported element types at zero length, recursively dereferencing pointers, and ignoring the permitted top-level indirection. The model describes returned size or panic and leaves cache timing and concurrency outside scope. Encoded widths must fit the native int. [Reviewer control](../../../difftest/encoding/binary/reflect-size/review/main.go)
