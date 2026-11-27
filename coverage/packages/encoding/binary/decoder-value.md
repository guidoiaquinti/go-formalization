# Independent review of `encoding/binary.decoder.value`

**Verdict: FAITHFUL**

Evidence SHA-256: `fd62db03d1ee431e77344d91a39f5800eeb1a60c8ba42ca0db596589f1316692`

Go skips blank fields by advancing the reflected size without touching or checking the backing. A narrow read checks logical length and a wide read checks capacity. The destination gets a complex value only once both components are available. On a failing second component, the cursor retains the completed first read while the prior complex value remains intact. The pinned native replay passed 63,136 author cases per method, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Encoding/Binary/CoderValue.lean) · [Comparison](../../../difftest/encoding/binary/coder-value/check.go)

Seven independent native and Lean fixtures passed and rejected six altered outcomes for blank fields, length versus capacity, partial complex effects and finite float bits. The checks include negative zero, a nonzero offset, retained spare bytes, and an ignored unsupported kind. The scope uses settable reflected values, finite floating components and built-in byte orders. [Reviewer control](../../../difftest/encoding/binary/coder-value/review/main.go)
