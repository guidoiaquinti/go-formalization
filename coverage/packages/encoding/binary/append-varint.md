# Independent review of `encoding/binary.AppendVarint`

**Verdict: FAITHFUL**

Evidence SHA-256: `06993614e3d61d74f12522f23494b3eba82edc63621c9426da0614eee0e1b991`

Pinned Go 1.27.0 `src/encoding/binary/varint.go:91–97` converts the signed input to its fixed-width word, shifts left, complements for negative input, and calls the gradual unsigned append. The model uses the same fixed-width sign and shift. It preserves original backing writes made before a later append grows the destination.

The author checked 1,070 boundary cases and 60,000 samples using seeds 17, 2718 and 65537. Seven independent literal native and Lean fixtures cover nil, exact fit, signed threshold 63/64, negative one, minimum int64, and partial growth. In the signed 64 fixture, one spare byte goes from 55 to 128, the returned buffer grows and contains [9,128,1], and the preceding sentinel stays 7. Encoding the payload separately and then bulk-appending fails four cases. Omitting the negative complement fails two cases.

Allocation size after growth and allocation failure remain outside scope.
