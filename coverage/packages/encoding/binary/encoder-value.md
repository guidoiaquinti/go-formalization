# Independent review of `encoding/binary.encoder.value`

**Verdict: FAITHFUL**

Evidence SHA-256: `ab683e27756761af4a9bbc894e1e0d7b46943efc962d3e4b91b727dc63a4e090`

Go clears blank fields against capacity before advancing. A failed clear leaves its bytes and offset intact. A subsequent narrow field can fail at logical length after an earlier clear into spare cells. An incomplete complex write preserves the first component and stops at the unchanged offset of the second component. The pinned native replay passed 63,136 author cases per method, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Encoding/Binary/CoderValue.lean) · [Comparison](../../../difftest/encoding/binary/coder-value/check.go)

Seven independent native and Lean fixtures passed and rejected six altered outcomes for blank fields, length versus capacity, partial complex effects and finite float bits. The checks include negative zero, a nonzero offset, retained spare bytes, and an ignored unsupported kind. The scope uses settable reflected values, finite floating components and built-in byte orders. [Reviewer control](../../../difftest/encoding/binary/coder-value/review/main.go)
