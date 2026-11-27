# Independent review of `encoding/binary.decoder.int16`

**Verdict: FAITHFUL**

Evidence SHA-256: `08d4b92308e17808093574f5856b0287c542c0f667fcd142aa971e6c3a3060fc`

Go 1.27.0 calls the selected ByteOrder method with a 2-byte, two-index view at the current offset. The view is allowed to extend beyond the logical slice length up to its capacity and aliases the original storage. Offset advancement follows a successful callback. A delegated panic preserves writes already made through the view and leaves the offset unchanged. A bounds failure takes precedence over invoking a nil ByteOrder. The signed wrapper reinterprets the returned word at the original width. A high-bit-only word became the exact signed minimum. The pinned Go, Lean, and oracle replay passed 61,995 inputs for this declaration, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Encoding/Binary/CoderWord.lean) · [Comparison](../../../difftest/encoding/binary/coder-word/check.go)

Five independent rows for this declaration passed. The shared control rejected checking logical length instead of capacity, committing offset before delegation, rolling back delegate writes, looking up the nil-interface panic before constructing the slice, and capping the view at its width. The signed decoder controls also rejected zero extension. The scope uses nonnil constructed coder receivers, the built-in byte orders, a nil interface, and the registered deterministic delegates. Other delegate behavior is outside this model. [Reviewer control](../../../difftest/encoding/binary/coder-word/review/main.go)
