# Independent review of `net/netip.Prefix.isZero`

**Verdict: FAITHFUL**

Evidence SHA-256: `309019cf2b8803aa54cd5ccd045fecc5d5d688a8f44b5ea2af4c9e81e04af7ee`

I traced the pinned Go 1.27.0 method, Lean model and private overlay. The entirety of the Prefix value participates in zero testing. An invalid length can coexist with a retained address and therefore yield false.

I ran `go run ./difftest/net/netip/prefix-ops/review` on darwin/arm64. Five independent literal rows matched pinned native Go and compiled Lean. Four mutation witnesses distinguish omitting the unmasked comparison tiebreak, conflating mapped IPv6 with IPv4, equating invalid with the zero Prefix, and accepting zoned candidates. The author replay matched 6,480 boundaries and 60,000 samples at seeds 17, 2718 and 65537, with a per-bit native oracle. [Reviewer control](../../../../difftest/net/netip/prefix-ops/review/main.go)

The modeled states are reachable address and Prefix values and sequential calls.
