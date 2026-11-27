# Independent review of `archive/zip.FileHeader.Mode`

**Verdict: FAITHFUL**

Reviewed author commit `1578735` against Go 1.27.0 `src/archive/zip/struct.go:317`. Unix and macOS creators decode the upper attribute word. NTFS, VFAT and FAT apply the DOS rules. Unknown creators start at zero. A trailing raw slash sets Dir even on a different Unix type.

The pinned native and Lean replay matched 139,264 fixed rows and 60,000 seeded rows at seeds 17, 2718 and 65537. Four independent literal rows also matched both implementations. The control rejected omitting the slash override in 2 rows. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/archive/zip/file-mode/review` after building `godiff-zip-file-mode`.

Header methods cover nonnil receivers and the metadata fields described by the claim. The comparisons exercise metadata values without archive I/O.

Evidence SHA-256: `2b7a43e73204bd2ab1c4c6752f412590c7e5081ef83a2581bf66c3af2fe71607`
