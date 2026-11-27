# Independent review of `bytes.Reader.Len`

**Verdict: FAITHFUL**

Evidence SHA-256: `a3d9879a83259b9375239954edbc69c503bd037b57912a39b022b800ff45c606`

Pinned Go 1.27.0 compares the signed reader cursor with the original slice length before subtracting. On public nonnegative positions the subtraction is bounded by the original slice length and fits in int, while positions at or beyond end produce zero. The routine leaves the cursor and previous rune unchanged and a nil receiver panics. All 2,278 boundaries and 60,000 samples at seeds 17, 2718 and 65537 matched. [Model](../../../lean/GoLean/Bytes/ReaderLen.lean) · [Comparison](../../../difftest/bytes/reader-len/check.go)

Nine independent fixtures passed for pinned Go and Lean, including signed int64 maximum and an offset that would vanish under a conversion to int32. The native sequence preserved UnreadRune after repeated Len calls and confirmed that a negative seek does not produce a negative cursor. The review excludes synthetic negative internal positions. [Reviewer control](../../../difftest/bytes/reader-len/review/main.go)
