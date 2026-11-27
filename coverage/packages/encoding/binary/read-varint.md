# Independent review of `encoding/binary.ReadVarint`

**Verdict: FAITHFUL**

Evidence SHA-256: `8a2cced67dc66f476d202a6c0442a67f63a86393ee2f4f7cb03b1224d8c4219c`

Go 1.27.0 performs signed zigzag decoding on the accumulated unsigned result whether or not the stream succeeded. The wrapper preserves the underlying error and read count. The independent native check observed both signs of partial value, exact custom and standard error identity, the effect of the tenth continuation's low bit on signed magnitude, and the valid ten-byte encoding of the minimum signed integer. All 65,124 author comparison inputs passed, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Encoding/Binary/VarintRead.lean) · [Comparison](../../../difftest/encoding/binary/varint-read/check.go)

Ten independent native and Lean fixtures passed. Replacing every signed result accompanied by an error with zero failed six rows. Both tenth-continuation fixtures disallowed an eleventh call. The native control provides exact error identity, while Lean represents categories. Arbitrary user reader internals and panics are outside the scripted reader domain. [Reviewer control](../../../difftest/encoding/binary/varint-read/review/main.go)
