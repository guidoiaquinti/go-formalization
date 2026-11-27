# Independent review of `encoding/binary.PutUvarint`

**Verdict: FAITHFUL**

Evidence SHA-256: `50559deef5422cb4e3f7f0991db91751a22f9dc3151d64b53d1d04edb906afdf`

I traced Go 1.27.0 `src/encoding/binary/varint.go`. The writer stores one encoded byte at a time. A short destination raises the index-bounds panic after the writes that fit, while prefix and suffix sentinels stay unchanged. The model's payload and write prefix preserve these observations. The signed wrapper folds the sign before entering the same unsigned writer.

The author passed 1,070 boundaries and 60,000 seeded inputs using 17, 2718 and 65537. Nine independent pinned darwin/arm64 and Lean fixtures passed, covering zero writable length, positive 63/64, negative -64/-65, machine extrema, capacity beyond the visible slice, and two completed writes before a panic. A length check before the first write would fail 6 fixtures for this method. Treating backing capacity as writable length also failed six unsigned cases.

The comparison observes the index-bounds panic class and exact modified bytes. Exact error text, instruction layout and other unexecuted platforms are outside this review.
