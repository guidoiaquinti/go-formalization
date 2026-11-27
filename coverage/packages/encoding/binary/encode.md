# Independent review of `encoding/binary.Encode`

**Verdict: FAITHFUL**

Evidence SHA-256: `99a9cfc23d7dfaf39ec75f5947148930854b1f5cd3c6a66f7e13e93449331de5`

Short input leaves the destination and its spare backing unchanged. A defined byte scalar writes exactly one byte and preserves the suffix. These fixtures reject writes before a failed preflight and overwriting spare cells. Pinned Go and Lean agreed on 68,644 author cases, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Encoding/Binary/PublicBuffer.lean) · [Comparison](../../../difftest/encoding/binary/public-buffer/check.go)

Ten independent Lean rows and direct native controls passed. The registered scope uses disjoint storage, built-in byte orders, and the supported primitive forms and defined byte forms. Exact allocator capacity is outside scope. [Reviewer control](../../../difftest/encoding/binary/public-buffer/review/main.go)
