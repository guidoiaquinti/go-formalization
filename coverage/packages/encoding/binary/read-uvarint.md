# Independent review of `encoding/binary.ReadUvarint`

**Verdict: FAITHFUL**

Evidence SHA-256: `e429e862089b766d45c97d020ac6244f0952ba7de80bfa66691652da02fe5e85`

Go 1.27.0 returns the accumulated value at the first reader error. Only direct EOF after a successful byte is translated to UnexpectedEOF. The native review checked exact error identity for EOF, UnexpectedEOF, a distinct error with the same EOF text, and wrapped EOF. It also checked that both overflow branches and both public methods share the same overflow error object. A terminating tenth byte greater than one leaves the accumulated value alone. A continuing tenth byte contributes its low bit at bit 63 before overflow. Both paths stop after exactly ten calls. All 65,124 author comparison inputs passed, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Encoding/Binary/VarintRead.lean) · [Comparison](../../../difftest/encoding/binary/varint-read/check.go)

Ten independent native and Lean fixtures passed. Matching wrapped EOF, adding an invalid tenth terminator, and discarding the tenth continuation contribution each failed one fixture. Attempting an eleventh call after the tenth continuation failed both corresponding fixtures. Lean represents error categories and read counts. The native control provides error identity and callback-budget checks. Arbitrary user reader internals and panics are outside the scripted reader domain. [Reviewer control](../../../difftest/encoding/binary/varint-read/review/main.go)
