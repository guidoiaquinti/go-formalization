# Independent review of `encoding/binary.Append`

**Verdict: FAITHFUL**

Evidence SHA-256: `2b7f2bd81d76dfef906558bb18a9b0f2a2ca55df6953ef17342a5c17624ad1f9`

The original prefix survives both reuse and allocation. Reuse writes exactly the appended payload and preserves remaining spare bytes. Growth detaches while leaving the old backing unchanged. Unsupported input returns nil and an error without changing backing. A zero-byte append preserves the distinction between nil and a nonnil empty view. Pinned Go and Lean agreed on 68,644 author cases, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Encoding/Binary/PublicBuffer.lean) · [Comparison](../../../difftest/encoding/binary/public-buffer/check.go)

Ten independent Lean rows and direct native controls passed. The registered scope uses disjoint storage, built-in byte orders, and the supported primitive forms and defined byte forms. Exact allocator capacity is outside scope. [Reviewer control](../../../difftest/encoding/binary/public-buffer/review/main.go)
