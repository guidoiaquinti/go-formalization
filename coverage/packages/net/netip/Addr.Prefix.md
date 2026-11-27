# Independent review of `net/netip.Addr.Prefix`

**Verdict: FAITHFUL**

Evidence SHA-256: `15859a1f0ae7234608e6a859368f712443ac02cde4689e5f46b493cc1d18c452`

I traced the pinned Go 1.27.0 method, Lean model and private overlay. The pinned method rejects negative lengths first, checks the family limit, and clears lower address bits. The prefix result loses an IPv6 zone.

I ran `go run ./difftest/net/netip/prefix-ops/review` on darwin/arm64. Five independent literal rows matched pinned native Go and compiled Lean. Four mutation witnesses distinguish omitting the unmasked comparison tiebreak, conflating mapped IPv6 with IPv4, equating invalid with the zero Prefix, and accepting zoned candidates. The author replay matched 6,480 boundaries and 60,000 samples at seeds 17, 2718 and 65537, with a per-bit native oracle. [Reviewer control](../../../../difftest/net/netip/prefix-ops/review/main.go)

The modeled states are reachable address and Prefix values and sequential calls.
