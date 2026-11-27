# Independent review of `net/url.URL.String`

**Verdict: FAITHFUL**

Reviewed author commit `e3f3dbc` against Go 1.27.0 `src/net/url/url.go`. Opaque values bypass authority and path rendering. Hierarchical rendering applies authority, relative-colon and omitted-host slash rules before optional query and fragment. The control rejects preserving an unescaped leading double slash when the empty host is omitted.

Pinned Go and Lean matched 603 fixed rows and 60,000 sampled rows at seeds 17, 2718 and 65537. Five independent literal cases passed in both implementations. Replay with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/net/url/serialize/review` after building `godiff-net-url-serialize`.

Scope covers direct URL values and valid slice geometry. Lean abstracts allocator-selected capacity and zero-capacity backing identity. Native checks nil receiver outcomes.

Evidence SHA-256: `1a89d40bf07104caa68ebf2a91a0a67e53adcfef709fcb556c4e5a6f4642b8d2`
