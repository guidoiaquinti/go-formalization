# Independent review of `archive/zip.FileHeader.SetMode`

**Verdict: FAITHFUL**

Reviewed author commit `1578735` against Go 1.27.0 `src/archive/zip/struct.go:331`. Creator becomes Unix while preserving the low version byte. The external attributes are replaced. The DOS directory bit uses the input Dir bit even when the Unix type combination falls back to regular. DOS readonly depends on the owner-write bit. The native overlay also checks the unchanged populated header fields.

The pinned native and Lean replay matched 139,264 fixed rows and 60,000 seeded rows at seeds 17, 2718 and 65537. Four independent literal rows also matched both implementations. The control rejected clearing the low version byte in 4 rows. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/archive/zip/file-mode/review` after building `godiff-zip-file-mode`.

Header methods cover nonnil receivers and the metadata fields described by the claim. The comparisons exercise metadata values without archive I/O.

Evidence SHA-256: `693e938325614ee5a8dee6004b80d28ae474a3b1458f1e1721b559a7842e13fa`
