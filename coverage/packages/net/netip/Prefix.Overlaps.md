# Independent review of `net/netip.Prefix.Overlaps`

**Verdict: FAITHFUL**

Evidence SHA-256: `633eb4db768a17d9ec19fd28ac3337ca8db86eaa0c3ce41ce57460192a739546`

I traced the pinned Go 1.27.0 method, Lean model and private overlay. Two valid prefixes of the same family overlap when the bits of the shorter prefix agree. Either input may retain host bits.

I ran `go run ./difftest/net/netip/prefix-ops/review` on darwin/arm64. Five independent literal rows matched pinned native Go and compiled Lean. Four mutation witnesses distinguish omitting the unmasked comparison tiebreak, conflating mapped IPv6 with IPv4, equating invalid with the zero Prefix, and accepting zoned candidates. The author replay matched 6,480 boundaries and 60,000 samples at seeds 17, 2718 and 65537, with a per-bit native oracle. [Reviewer control](../../../../difftest/net/netip/prefix-ops/review/main.go)

The modeled states are reachable address and Prefix values and sequential calls.
