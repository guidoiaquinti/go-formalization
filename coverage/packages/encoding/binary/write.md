# Independent review of `encoding/binary.Write`

**Verdict: FAITHFUL**

Evidence SHA-256: `51dc8adfd43b0cd15cef4444d3c04d3ca5b9ba4754ca0d11d29551ef0cc45f82`

Pinned Go calls the writer once with the encoded payload. For ordinary nonempty byte slices, the supplied payload shares the original backing. For a slice of defined byte elements, the writer receives a distinct encoded slice. The count is ignored even when a writer scripts a short successful result, while a returned error is propagated by identity. Supported zero-width data still makes one writer call and rejected types make none. All 68,644 author cases and nine reviewer Lean rows plus direct native controls passed. [Model](../../../lean/GoLean/Encoding/Binary/PublicStream.lean) · [Comparison](../../../difftest/encoding/binary/public-stream/check.go)

The control rejects copying the ordinary byte slice, aliasing the defined-element slice, substituting a short-write error, and skipping the empty payload call. The scope uses built-in byte orders and deterministic callbacks. [Reviewer control](../../../difftest/encoding/binary/public-stream/review/main.go)
