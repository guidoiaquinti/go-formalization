# Independent review of `net/url.URL.MarshalBinary`

**Verdict: FAITHFUL**

Reviewed author commit `e3f3dbc` against Go 1.27.0 `src/net/url/url.go`. Delegates to append with a nil slice and returns the same serialized URL as String. Native checks that empty success retains a nil slice. The control rejects using RequestURI bytes on four rows.

Pinned Go and Lean matched 603 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537. Five independent literal cases passed in both implementations. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/net/url/serialize/review` after building `godiff-net-url-serialize`.

Scope covers direct URL values and valid slice geometry. Lean abstracts allocator-selected capacity and zero-capacity backing identity. Native checks nil receiver outcomes.

Evidence SHA-256: `152a3e26d38636ec37300003ec84bc2e19273bcb52da667ab28e8ec72c6fbd74`
