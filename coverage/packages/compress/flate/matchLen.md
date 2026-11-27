# Independent review of `compress/flate.matchLen`

**Verdict: FAITHFUL**

The pinned Go 1.27.0 implementation checks each eight-byte load against the right slice's capacity. A mismatch in a successfully loaded word returns before the tail transition, even when the differing offset is beyond the original visible length. The little-endian XOR places earlier differing bytes in lower bits, so the trailing-zero count divided by eight gives the first differing byte in that word.

When all preceding words match, `b[n:]` requires the consumed prefix to fit the original visible view, including when the left tail is empty. The subsequent `b[:len(a)]` can extend that tail through the remaining capacity and must succeed before any tail mismatch is observed. Lean models both checks and the earlier word-load checks. The claim covers legal slice geometry and distinguishes a slice bounds failure from a successful result.

The pinned native Go and Lean replay passed 6,084 fixed and 60,000 sampled cases on darwin/arm64 using seeds 17, 2718 and 65537. The three registered theorems cover empty left input, reslicing an initially empty visible view into a matching two-byte tail, and failure after an equal word when its eight-byte extent exceeds the visible length. The scoped Lean axiom audit found only `propext`, `Classical.choice` and `Quot.sound`.

The reviewer control uses 13 literal outcomes on both compiled Lean and the pinned Go overlay. It distinguishes capacity from visible length, and exercises the last byte of each of two blocks, matching words followed by a zero-length tail, and a partial tail with insufficient capacity despite an available first mismatch. Limiting reads to visible length fails 6 controls. Omitting the visible-length check at tail transition fails 2. Checking bounds only as each byte is read fails 3.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/compress/flate/match-length/review
```

Evidence SHA-256: `40299f6809f95e9faecedd494267fdb56a886f0084bc7dc2421d430339d69aef`
