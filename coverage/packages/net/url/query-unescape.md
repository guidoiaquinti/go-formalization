# Independent review of `net/url.QueryUnescape`

**Verdict: FAITHFUL**

Evidence SHA-256: `b1d562a2f35148f318e450569c0525036720b20acb9fd603da56c574ca99e1fe`

I checked the pinned Go 1.27.0 scanner and byte reconstruction against compiled Lean and the exact private or public Go call. The query wrapper decodes each valid percent pair and maps literal plus to space. An invalid escape returns the first malformed sequence, limited to three raw bytes.

I ran `go run ./difftest/net/url/unescape/review` on darwin/arm64. Nine independent literal rows matched all three native outputs and compiled Lean, including typed errors and raw payloads. Four mutation witnesses distinguish converting a path plus, keeping too much error input, accepting forbidden host escapes, and rejecting an escaped zone space. The author comparison passed 71,168 boundaries and 60,000 samples at seeds 17, 2718 and 65537. [Reviewer control](../../../../difftest/net/url/unescape/review/main.go)

The domain covers raw byte strings, observable decoded bytes and the first error type and payload. Error message formatting and allocation behavior are outside scope.
