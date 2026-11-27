# Independent review of `archive/zip.fileListEntry.stat`

**Verdict: FAITHFUL**

Reviewed author commit `669f20f` against Go 1.27.0 `src/archive/zip/reader.go`. A duplicate precedes file and directory selection. Its fresh error includes the raw name, including malformed bytes and NUL. A directory returns the original entry. A regular file returns a live pointer view, and a missing file faults. The control rejects directory-first selection and dropping the raw duplicate name on one row each.

Pinned Go and Lean matched 256 fixed rows and 60,000 samples at seeds 17, 2718 and 65537. Five independent literal rows matched both. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/archive/zip/entry-stat/review` after building `godiff-zip-entry-stat`.

The scope is metadata on nonnil outer receivers, including absent nested entry or file pointers. Lean abstracts pointer identities and the nil-dereference panic class. Native checks exact pointer identity and live views.

Evidence SHA-256: `f976aa141b326a528456c49fc2cc6ce98c95a43ed6a69320816db84c774e6319`
