# Independent review of `archive/zip.msdosModeToFileMode`

**Verdict: FAITHFUL**

Reviewed author commit `1578735` against Go 1.27.0 `src/archive/zip/struct.go:348`. Returns directory with 0777 or regular permissions 0666. Readonly clears all three write bits. The registered theorem checks that only those two attribute bits affect the result.

The pinned native and Lean replay matched 139,264 fixed rows and 60,000 seeded rows at seeds 17, 2718 and 65537. Four independent literal rows also matched both implementations. The control rejected retaining write permissions for readonly entries in 3 rows. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/archive/zip/file-mode/review` after building `godiff-zip-file-mode`.

Header methods cover nonnil receivers and the metadata fields described by the claim. The comparisons exercise metadata values without archive I/O.

Evidence SHA-256: `aeabf4f62cf4557d9d91eb785b189a760faadf243c4103d26071161e7a135d93`
