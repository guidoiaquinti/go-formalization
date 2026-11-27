# Independent review of `bytes.Reader.WriteTo`

**Verdict: FAITHFUL**

Evidence SHA-256: `326f2f30662db9e9ea4042b6ff1f66ed204321154c3b25cec384c4da77126c95`

I read the pinned Go 1.27.0 declaration and checked the model and native comparison. The pinned method invalidates the previous rune before dispatch. EOF bypasses the writer. A live reader passes one slice with the original alias and capacity. On return, an overcount panics before updating the cursor, any other count updates it, and a writer error takes precedence over the short-write sentinel. The negative count fixture preserves the signed count and cursor change.

I ran `go run ./difftest/bytes/reader-position/review` on darwin/arm64. Eight independent literal paired traces matched native Go and compiled Lean. The checks rejected retaining the previous rune after an invalid whence, saturating on overflow, substituting short-write for a supplied error, advancing the cursor before rejecting an overcount, and rejecting negative counts. The author comparison matched 8,403 boundaries and 60,000 rows at seeds 17, 2718 and 65537 per declaration. [Reviewer control](../../../difftest/bytes/reader-position/review/main.go)

The input scope is sequential access with a nonnegative initial cursor and legal slice geometry. Reentrant or mutating writer callbacks are outside scope.
