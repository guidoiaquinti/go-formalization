# Independent review of `encoding/binary.AppendUvarint`

**Verdict: FAITHFUL**

Evidence SHA-256: `88ff19b80220ebbf73071a4abce64593bd0f319c09c73a7d4e023c8dcde13e26`

Pinned Go 1.27.0 `src/encoding/binary/varint.go:41–47` appends one marked low-seven-bit group at a time, shifts right by seven, and appends the unmarked last byte. The model's resulting payload agrees with this order. The Go loop writes all groups that fit into the old spare capacity before the first growth. The model preserves that visible portion of the original backing even when the return uses new storage.

The author checked 1,070 boundary cases and 60,000 samples using seeds 17, 2718 and 65537. Seven independent literal native and Lean fixtures cover nil, exact fit, split growth with one or two spare bytes, encoded width changes and the ten-byte unsigned maximum. In the offset-view fixture for 128, the return is [9,128,1], while the original array is [7,9,128]. Replacing the loop with one bulk append of the complete payload leaves that spare byte unchanged and fails four cases.

Allocation size after growth and allocation failure remain outside scope.
