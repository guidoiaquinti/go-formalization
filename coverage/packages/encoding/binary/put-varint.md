# Independent review of `encoding/binary.PutVarint`

**Verdict: FAITHFUL**

Evidence SHA-256: `8fe71bd6177212e208ec6b501ab57573b1671ab92aa18ffa140f2c7f7aa0da67`

I traced Go 1.27.0 `src/encoding/binary/varint.go`. The writer stores one encoded byte at a time. A short destination raises the index-bounds panic after the writes that fit, while prefix and suffix sentinels stay unchanged. The model's payload and write prefix preserve these observations. The signed wrapper folds the sign before entering the same unsigned writer.

The author passed 1,070 boundaries and 60,000 seeded inputs using 17, 2718 and 65537. Nine independent pinned darwin/arm64 and Lean fixtures passed, covering zero writable length, positive 63/64, negative -64/-65, machine extrema, capacity beyond the visible slice, and two completed writes before a panic. A length check before the first write would fail 5 fixtures for this method. Omitting the signed complement failed four cases, including the negative boundary values.

The comparison observes the index-bounds panic class and exact modified bytes. Exact error text, instruction layout and other unexecuted platforms are outside this review.
