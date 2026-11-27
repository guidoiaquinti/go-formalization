# Independent review of `bytes.Reader.Seek`

**Verdict: FAITHFUL**

Evidence SHA-256: `54aec40ca7dd481cded6cf7aed486c4364c1031cbcb71b046761837ca80f4c08`

I read the pinned Go 1.27.0 declaration and checked the model and native comparison. The pinned method invalidates the previous rune before whence validation. It selects the start, current position or end as the base. An invalid whence or negative resulting position preserves the cursor. The MaxInt64 plus one fixture exercises signed wrap into the negative position error.

I ran `go run ./difftest/bytes/reader-position/review` on darwin/arm64. Eight independent literal paired traces matched native Go and compiled Lean. The checks rejected retaining the previous rune after an invalid whence, saturating on overflow, substituting short-write for a supplied error, advancing the cursor before rejecting an overcount, and rejecting negative counts. The author comparison matched 8,403 boundaries and 60,000 rows at seeds 17, 2718 and 65537 per declaration. [Reviewer control](../../../difftest/bytes/reader-position/review/main.go)

The input scope is sequential access with a nonnegative initial cursor and legal slice geometry. Reentrant or mutating writer callbacks are outside scope.
