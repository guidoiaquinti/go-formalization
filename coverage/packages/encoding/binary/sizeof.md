# Independent review of `encoding/binary.sizeof`

**Verdict: FAITHFUL**

Evidence SHA-256: `ee7116e04c0b007655cdca420ee42fab2bb0a70f2c74d902d4e2d5f2be6a62ca`

Private sizeof adds encoded field widths including blank fields, without native alignment padding. It rejects unsupported elements even in zero-length arrays. A nil reflect.Type raises the native nil dereference panic, distinct from a valid but unsupported type. The pinned Go, Lean and layout-oracle replay matched all 61,814 layouts, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Encoding/Binary/ReflectSize.lean) · [Comparison](../../../difftest/encoding/binary/reflect-size/check.go)

Ten independent native and Lean layouts passed. The control rejected using the native padded width, accepting unsupported element types at zero length, recursively dereferencing pointers, and ignoring the permitted top-level indirection. The model describes returned size or panic and leaves cache timing and concurrency outside scope. Encoded widths must fit the native int. [Reviewer control](../../../difftest/encoding/binary/reflect-size/review/main.go)
