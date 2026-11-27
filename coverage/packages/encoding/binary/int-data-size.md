# Independent review of `encoding/binary.intDataSize`

**Verdict: FAITHFUL**

Evidence SHA-256: `3cac5d163e305c43f52a9b6c332c4adf0a6b7e8a1569f70fe263896ee04ce1ea`

Go 1.27.0 uses exact dynamic types for its fast sizes. Typed nil pointers to supported scalars still receive the scalar width. Defined pointer types, defined slice types, and slices of defined byte elements miss the switch and return zero with a nil byte result. A type alias retains the original type identity. The byte result from an actual byte slice preserves nilness, length, capacity and backing storage. The author replay passed all 61,440 native and Lean rows, including 60,000 sampled inputs at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Encoding/Binary/IntDataSize.lean) · [Comparison](../../../difftest/encoding/binary/int-data-size/check.go)

Eight independent Lean rows passed, alongside native type-switch and offset-slice controls. Incorrect nil-pointer rejection failed two controls, matching defined types by their underlying kind failed four, and clipping capacity or copying the backing each failed two. The mutation checks included a nonnil zero-length byte view with spare capacity. Lean preserves the original modeled view and the native check verifies write-through identity. Slice sizes whose byte count cannot fit a native int are outside scope. [Reviewer control](../../../difftest/encoding/binary/int-data-size/review/main.go)
