# Independent review of `archive/zip.checksumReader.Stat`

**Verdict: FAITHFUL**

Reviewed author commit `669f20f` against Go 1.27.0 `src/archive/zip/reader.go`. Returns a live view of the backing file header regardless of the sticky read error or byte count. Missing file faults. The control rejects returning the sticky read error on all five rows. Native checks that reader progress and the error are unchanged.

Pinned Go and Lean matched 256 fixed rows and 60,000 samples at seeds 17, 2718 and 65537. Five independent literal rows matched both. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/archive/zip/entry-stat/review` after building `godiff-zip-entry-stat`.

The scope is metadata on nonnil outer receivers, including absent nested entry or file pointers. Lean abstracts pointer identities and the nil-dereference panic class. Native checks exact pointer identity and live views.

Evidence SHA-256: `d0cf70fa72233889706fb6a8b16f5162e789630a8aff06d2c05bc5e6eee21d46`
