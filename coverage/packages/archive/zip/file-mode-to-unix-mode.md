# Independent review of `archive/zip.fileModeToUnixMode`

**Verdict: FAITHFUL**

Reviewed author commit `1578735` against Go 1.27.0 `src/archive/zip/struct.go:360`. Switches on the exact Go type mask, including the character-device combination. Mixed unsupported masks fall back to regular while retaining permissions and special bits. The property checks the nine copied permission bits.

The pinned native and Lean replay matched 139,264 fixed rows and 60,000 seeded rows at seeds 17, 2718 and 65537. Four independent literal rows also matched both implementations. The control rejected preferring directory on a conflicting Go type mask in 1 rows. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/archive/zip/file-mode/review` after building `godiff-zip-file-mode`.

Header methods cover nonnil receivers and the metadata fields described by the claim. The comparisons exercise metadata values without archive I/O.

Evidence SHA-256: `b580735e83f47ca866a34d1ee489fd7a5f7065d97b4525c428a464243a041aa3`
