# Independent review of `encoding/binary.Read`

**Verdict: FAITHFUL**

Evidence SHA-256: `baf8328957f1c763c09ebff5f2eaf6cd6226706887c04abc0f5a51a8064451a1`

Pinned Go requests the entire encoded width initially and exactly the remaining bytes after each successful short read. The destination stays unchanged on an incomplete read. Direct partial EOF becomes UnexpectedEOF, while a same-text custom error and wrapped EOF retain their identities. A terminal error returned with the final requested byte is suppressed by ReadFull. An ordinary scalar invokes the reader before falling back to an invalid-type error. A defined byte scalar fails before I/O, and a supported zero-width slice makes no reader call. All 68,644 author cases and nine reviewer Lean rows plus direct native controls passed. [Model](../../../lean/GoLean/Encoding/Binary/PublicStream.lean) · [Comparison](../../../difftest/encoding/binary/public-stream/check.go)

The controls reject fixed-size requests after partial success, converting wrapped EOF, mutating the destination on an incomplete read, and invoking the reader for a zero-width value. The model uses the registered deterministic reader and primitive value scope. [Reviewer control](../../../difftest/encoding/binary/public-stream/review/main.go)
