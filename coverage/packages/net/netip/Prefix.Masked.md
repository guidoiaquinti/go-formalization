# Independent review of `net/netip.Prefix.Masked`

**Verdict: FAITHFUL**

Evidence SHA-256: `689a0ddc1cd7461b64b12ca8a1028d98d429981ce1765b2d23edcfd565180106`

I traced the pinned Go 1.27.0 method, Lean model and private overlay. The method returns the canonical masked address for a valid prefix and the zero Prefix for an invalid length.

I ran `go run ./difftest/net/netip/prefix-ops/review` on darwin/arm64. Five independent literal rows matched pinned native Go and compiled Lean. Four mutation witnesses distinguish omitting the unmasked comparison tiebreak, conflating mapped IPv6 with IPv4, equating invalid with the zero Prefix, and accepting zoned candidates. The author replay matched 6,480 boundaries and 60,000 samples at seeds 17, 2718 and 65537, with a per-bit native oracle. [Reviewer control](../../../../difftest/net/netip/prefix-ops/review/main.go)

The modeled states are reachable address and Prefix values and sequential calls.
