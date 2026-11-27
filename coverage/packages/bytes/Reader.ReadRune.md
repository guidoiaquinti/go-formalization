# Independent review of `bytes.Reader.ReadRune`

**Verdict: FAITHFUL**

Evidence SHA-256: `102ec5774a070cea5815dca40f567a12eff3729921b9b77fb590519226c57e8e`

Go 1.27.0 clears the saved marker at exhaustion and returns zero rune, zero width and the exact EOF sentinel. The live path saves the input cursor, decodes ASCII in one byte or uses UTF-8 decoding, and advances by the returned width. Both paths preserve the original backing slice. Nil receiver dereference and error identity were checked in the native overlay. All 64,099 author inputs matched, including 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Bytes/ReaderRune.lean) · [Comparison](../../../difftest/bytes/reader-rune/check.go)

Nine independent paired native and Lean rows passed. Treating every replacement rune as width one failed one row. Retaining the old marker at EOF failed three. A public sequence distinguished a literal replacement rune from malformed input, rewound the literal rune, and checked that EOF invalidated the marker. Synthetic negative cursors and concurrent changes to the backing slice are outside scope. [Reviewer control](../../../difftest/bytes/reader-rune/review/main.go)
