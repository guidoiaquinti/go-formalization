# Independent review of `net/url.URL.AppendBinary`

**Verdict: FAITHFUL**

Reviewed author commit `e3f3dbc` against Go 1.27.0 `src/net/url/url.go`. Appends the whole serialized URL after the prefix. The control covers reused and detached backing. Native verifies untouched spare bytes, unchanged backing on growth and before nil receiver failure, and correct output. The control rejects forced detachment on three rows.

Pinned Go and Lean matched 603 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537. Five independent literal cases passed in both implementations. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/net/url/serialize/review` after building `godiff-net-url-serialize`.

Scope covers direct URL values and valid slice geometry. Lean abstracts allocator-selected capacity and zero-capacity backing identity. Native checks nil receiver outcomes.

Evidence SHA-256: `b0ca80e163caf60b9501d8dd265fd15d41a0d6d5ee30c93047cadbd7cb4f7ba1`
