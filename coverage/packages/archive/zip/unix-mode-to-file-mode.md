# Independent review of `archive/zip.unixModeToFileMode`

**Verdict: FAITHFUL**

Reviewed author commit `1578735` against Go 1.27.0 `src/archive/zip/struct.go:390`. Maps Unix file types to Go mode bits and preserves permissions and special bits. Unrecognized types contribute no type bits. The property checks preservation of all nine permission bits.

The pinned native and Lean replay matched 139,264 fixed rows and 60,000 seeded rows at seeds 17, 2718 and 65537. Four independent literal rows also matched both implementations. The control rejected leaking raw Unix type bits into the output in 1 rows. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/archive/zip/file-mode/review` after building `godiff-zip-file-mode`.

Header methods cover nonnil receivers and the metadata fields described by the claim. The comparisons exercise metadata values without archive I/O.

Evidence SHA-256: `f0d26f1813e3a30db5670b204af99c6250721f902fac0b2b4c4617dc3bb2f75c`
