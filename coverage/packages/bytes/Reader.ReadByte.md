# Independent review of `bytes.Reader.ReadByte`

**Verdict: FAITHFUL**

Evidence SHA-256: `d47efbfaa527ed1d8b8a62f66ac5bc898b7958303dee2f45c0e877d9f7b7863e`

Go 1.27.0 clears the previous-rune marker before checking the cursor. At and beyond end it returns zero and io.EOF without moving the cursor. Before end it returns the raw byte and advances by one. A nil receiver panics. The model preserves the original storage and its spare bytes. The direct overlay matched all 63,859 fixtures, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Bytes/ReaderByte.lean) · [Comparison](../../../difftest/bytes/reader-byte/check.go)

Nine independent paired native and Lean fixtures passed. Retaining the rune marker on EOF failed five fixtures. The public control read a multibyte rune, then a failed ReadByte, and checked that UnreadRune rejected the stale marker while UnreadByte still rewound one raw byte. Synthetic negative cursors are outside scope. [Reviewer control](../../../difftest/bytes/reader-byte/review/main.go)
