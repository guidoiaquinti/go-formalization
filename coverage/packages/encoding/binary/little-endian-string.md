# Independent review of `encoding/binary.littleEndian.String`

**Verdict: FAITHFUL**

Evidence SHA-256: `c04cef849f59afd89161433c0f0493efe70ca9e8a1aebb329e5392f0010cc45c`

The pinned Go 1.27.0 method in `src/encoding/binary/binary.go` has one receiver value and returns the constant `LittleEndian`. The Lean theorem is the same exact constant. The shared independent control calls all six methods directly and checks all six compiled Lean rows. The `String` value also matched quoted string formatting on darwin/arm64. Substituting LittleEndian for NativeEndian fails both native-display rows, and dropping the package prefix fails all three GoString rows.

The claim concerns display tokens. Reader, writer and architecture byte-order behavior are outside scope.
