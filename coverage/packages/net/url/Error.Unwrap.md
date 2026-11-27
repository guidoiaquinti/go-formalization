# Independent review of `net/url.Error.Unwrap`

**Verdict: FAITHFUL**

Evidence SHA-256: `3816ec8b967a32570ff37319931306dbad93e3e6e48180526b021052578f91af`

I reviewed evidence commit `de25499` against pinned Go 1.27.0 `src/net/url/url.go`. Pinned Go returns the immediate stored error interface. A wrapper around another error remains the immediate result even when the nested error is reachable with errors.Is. Native tests verify typed nil and identity and the Lean model carries an abstract identity token.

The pinned comparison at `go run ./difftest/net/url/error-methods` passed 3,325 fixed rows and 20,000 sampled rows per seed for 17, 2718 and 65537 on darwin/arm64. Five independent cases passed with `go run ./difftest/net/url/error-methods/review` and rejected nine altered outcomes across the group. I regenerated the Unicode quote-print intervals using the pinned generator and obtained identical bytes. The underlying pinned strconv.IsPrint and unicode.IsPrint agreed on every code point from 0 through 0x10ffff. The scoped kernel audit passed. Nested error rendering and callback implementation remain supplied model inputs. [Independent control](../../../../difftest/net/url/error-methods/review/main.go)
