# Independent review of `encoding/binary.decodeFast`

**Verdict: FAITHFUL**

Evidence SHA-256: `c5d36dcbe779ff6a588503a0a65897dde41edd22f697079e0a44620960a0df1c`

Go 1.27.0 treats byte slices as a bounded copy. Bool and signed byte destination slices iterate over the supplied source bytes and may panic at the first destination overflow. Wide slices keep earlier assignments on a short source. Unsupported values return false, including supported scalar types supplied by value. Raw float bits survive decoding, including negative zero and signaling NaN. All 68,640 author cases passed, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Encoding/Binary/FastCodec.lean) · [Comparison](../../../difftest/encoding/binary/fast-codec/check.go)

Nine independent native and Lean fixtures passed. Using bounded copy behavior for signed bytes and bools failed two cases. The controls also rejected rolling back a completed wide assignment and canonicalizing float bits. Built-in byte orders, nonnil scalar pointers and disjoint buffers define the scope. [Reviewer control](../../../difftest/encoding/binary/fast-codec/review/main.go)
