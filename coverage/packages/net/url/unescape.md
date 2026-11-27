# Independent review of `net/url.unescape`

**Verdict: FAITHFUL**

Evidence SHA-256: `efd8861b2d7d866982a5bd503490582e50370a7dad46e14e220f9783bc557681`

I checked the pinned Go 1.27.0 scanner and byte reconstruction against compiled Lean and the exact private or public Go call. The pinned helper returns the first failing escape or host byte as a typed raw payload. In host mode only percent itself can be introduced using an ASCII escape. Zone mode permits an encoded space but preserves its separate host-byte checks.

I ran `go run ./difftest/net/url/unescape/review` on darwin/arm64. Nine independent literal rows matched all three native outputs and compiled Lean, including typed errors and raw payloads. Four mutation witnesses distinguish converting a path plus, keeping too much error input, accepting forbidden host escapes, and rejecting an escaped zone space. The author comparison passed 71,168 boundaries and 60,000 samples at seeds 17, 2718 and 65537. [Reviewer control](../../../../difftest/net/url/unescape/review/main.go)

The domain covers raw byte strings, observable decoded bytes and the first error type and payload. Error message formatting and allocation behavior are outside scope.
