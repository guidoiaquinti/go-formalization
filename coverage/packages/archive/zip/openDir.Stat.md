# Independent review of `archive/zip.openDir.Stat`

**Verdict: FAITHFUL**

Reviewed author commit `669f20f` against Go 1.27.0 `src/archive/zip/reader.go`. Delegates to its stored entry, preserving the duplicate error, directory pointer, or live header view. A missing entry faults even when a separate file is present. The control rejects substituting a header on that row. Native checks that the read offset is unchanged.

Pinned Go and Lean matched 256 fixed rows and 60,000 samples at seeds 17, 2718 and 65537. Five independent literal rows matched both. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/archive/zip/entry-stat/review` after building `godiff-zip-entry-stat`.

The scope is metadata on nonnil outer receivers, including absent nested entry or file pointers. Lean abstracts pointer identities and the nil-dereference panic class. Native checks exact pointer identity and live views.

Evidence SHA-256: `5c92e12c443d63300b11932f79d9f8c4825733042792cff84b02b916ec21822c`
