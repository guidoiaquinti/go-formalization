# Independent review of `net/url.URL.RequestURI`

**Verdict: FAITHFUL**

Reviewed author commit `e3f3dbc` against Go 1.27.0 `src/net/url/url.go`. Returns only the path or opaque target and raw query. Opaque values starting with two slashes carry their scheme, and an empty path defaults to slash. The control rejects adding a fragment in two rows.

Pinned Go and Lean matched 603 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537. Five independent literal cases passed in both implementations. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/net/url/serialize/review` after building `godiff-net-url-serialize`.

Scope covers direct URL values and valid slice geometry. Lean abstracts allocator-selected capacity and zero-capacity backing identity. Native checks nil receiver outcomes.

Evidence SHA-256: `34110f95f5bfade115abcd27b9b4848493156659dfd1a8c5f969294a618fbe17`
