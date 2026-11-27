# Independent review of `encoding/binary.Decode`

**Verdict: FAITHFUL**

Evidence SHA-256: `f0bd58c0841f50b6f1f2765fddaa65558b958aa6e47f58d3713e66bcbd08dbfc`

Short-buffer preflight takes precedence over an invalid ordinary scalar destination, while a defined byte scalar reaches invalid-type handling. Oversized buffers update both bool or int8 entries and then panic. Ordinary and defined byte slices decode only the necessary prefix and report its length. The callback-free controls reject clipping the bool and int8 source and rolling back partial writes. Pinned Go and Lean agreed on 68,644 author cases, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Encoding/Binary/PublicBuffer.lean) · [Comparison](../../../difftest/encoding/binary/public-buffer/check.go)

Ten independent Lean rows and direct native controls passed. The registered scope uses disjoint storage, built-in byte orders, and the supported primitive forms and defined byte forms. Exact allocator capacity is outside scope. [Reviewer control](../../../difftest/encoding/binary/public-buffer/review/main.go)
