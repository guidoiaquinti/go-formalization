# Independent review of `net/url.URL.Clone`

**Verdict: FAITHFUL**

Evidence SHA-256: `23f82b675c55cc6b1a686ac84ef2397d3a4f40660648d0746f63ef64aabe9b3c`

I checked evidence commit `7ae5464` against pinned Go 1.27.0 `src/net/url/url.go`. A nil receiver returns nil. Otherwise Go copies the URL struct into a new allocation and separately copies nonnil userinfo. The Lean model copies all value fields and uses distinct abstract identities for the allocations. Its freshness theorem assumes the supplied abstract identities differ from the original ones, and the direct native check exercises both pointer inequalities.

The pinned Go and Lean comparison at `go run ./difftest/net/url/clone` passed 49 boundaries and 20,000 samples each for seeds 17, 2718 and 65537 on darwin/arm64. Four independent native and Lean fixtures from `go run ./difftest/net/url/clone/review` distinguish nil from an allocated zero URL, nil from empty userinfo, and an empty password present from one absent. The native fixture also changes several fields through the returned pointer while checking that the original retains its values. Six altered observations fail, including a shared userinfo pointer and an omitted raw path. Allocation layout and unsafe string sharing are outside scope. [Independent control](../../../../difftest/net/url/clone/review/main.go)
