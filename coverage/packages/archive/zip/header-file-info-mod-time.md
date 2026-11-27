# Independent review of `archive/zip.headerFileInfo.ModTime`

**Verdict: FAITHFUL**

Reviewed author commit `efb8dad` against Go 1.27.0 `src/archive/zip/struct.go:185`. The legacy DOS words supply the return value only when the modern timestamp is the zero instant. A non-UTC location on that instant still takes the fallback. One nanosecond after zero uses the modern timestamp. The native result has the canonical UTC location, and the modern path removes the monotonic clock reading.

Pinned Go and Lean matched 600 boundary rows and 60,000 rows using seeds 17, 2718 and 65537. Four independent literal cases passed. The control rejects unconditional legacy fallback on two rows, seconds-only zero detection on one and retaining the original location on two. The independent native time.Now run also detected a retained monotonic reading. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/archive/zip/header-modtime/review` after building `godiff-zip-header-modtime`.

The scope covers nonnil headers. Lean models the instant as seconds and nanoseconds. Native checks complete time equality, including UTC location identity and stripped monotonic state.
Evidence SHA-256: `d28091a1100d19c1fe6d920a77d750caa55fd2acea492d84cc86c6a00bcc3dc3`
