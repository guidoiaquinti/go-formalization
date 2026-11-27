# Independent review of `net/netip.Prefix.Contains`

**Verdict: FAITHFUL**

Evidence SHA-256: `ed07024bcdb90f2bbc76ced1af89668db52c39b03e68ba2fed7051036320dbcc`

I traced the pinned Go 1.27.0 method, Lean model and private overlay. Membership compares only the prefix bits and requires the same family and an unzoned candidate. IPv4-mapped IPv6 stays in the IPv6 family.

I ran `go run ./difftest/net/netip/prefix-ops/review` on darwin/arm64. Five independent literal rows matched pinned native Go and compiled Lean. Four mutation witnesses distinguish omitting the unmasked comparison tiebreak, conflating mapped IPv6 with IPv4, equating invalid with the zero Prefix, and accepting zoned candidates. The author replay matched 6,480 boundaries and 60,000 samples at seeds 17, 2718 and 65537, with a per-bit native oracle. [Reviewer control](../../../../difftest/net/netip/prefix-ops/review/main.go)

The modeled states are reachable address and Prefix values and sequential calls.
