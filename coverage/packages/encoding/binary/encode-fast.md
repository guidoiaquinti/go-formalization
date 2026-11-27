# Independent review of `encoding/binary.encodeFast`

**Verdict: FAITHFUL**

Evidence SHA-256: `58c01684d6ae44be0d8c37086a11130b1f97d967f014fff4d0306ef5580a9531`

Go 1.27.0 uses a bounded copy for byte slices. Signed byte and bool slices iterate and can panic after writing a prefix. Wide encoders keep completed writes and leave the incomplete word and spare cells untouched. The fast encoder accepts supported scalar values and ignores rejected dynamic types. Independent raw-bit cases preserved negative zero and a signaling NaN payload. All 68,640 author cases passed, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Encoding/Binary/FastCodec.lean) · [Comparison](../../../difftest/encoding/binary/fast-codec/check.go)

Nine independent native and Lean fixtures passed. Using a panicking iteration for byte copy, rolling back partial writes, and canonicalizing floating-point bits were rejected. Built-in byte orders, nonnil scalar pointers and disjoint buffers define the scope. [Reviewer control](../../../difftest/encoding/binary/fast-codec/review/main.go)
