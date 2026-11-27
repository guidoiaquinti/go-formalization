# Independent review of `net/url.InvalidHostError.Error`

**Verdict: FAITHFUL**

Evidence SHA-256: `7d174b35708bb0fc4f6e3af00637dc862be2ef5f7e12c685a78022209c09065c`

I reviewed evidence commit `de25499` against pinned Go 1.27.0 `src/net/url/url.go`. Pinned Go adds the exact prefix, Go-quoted raw input and the host-name suffix. The independent literal distinguishes a malformed byte escape from a literal replacement rune, and detects dropping the suffix.

The pinned comparison at `go run ./difftest/net/url/error-methods` passed 3,325 fixed rows and 20,000 sampled rows per seed for 17, 2718 and 65537 on darwin/arm64. Five independent cases passed with `go run ./difftest/net/url/error-methods/review` and rejected nine altered outcomes across the group. I regenerated the Unicode quote-print intervals using the pinned generator and obtained identical bytes. The underlying pinned strconv.IsPrint and unicode.IsPrint agreed on every code point from 0 through 0x10ffff. The scoped kernel audit passed. Nested error rendering and callback implementation remain supplied model inputs. [Independent control](../../../../difftest/net/url/error-methods/review/main.go)
