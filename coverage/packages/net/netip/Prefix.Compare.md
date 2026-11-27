# Independent review of `net/netip.Prefix.Compare`

**Verdict: FAITHFUL**

Evidence SHA-256: `9994b5c16114cf1a868206fe01c69fb8b25e1fe09b5201546e5070626630d732`

I traced the pinned Go 1.27.0 method, Lean model and private overlay. The comparison orders masked address, then prefix length, then retained unmasked address. Distinct host bits can break a tie between two equivalent networks.

I ran `go run ./difftest/net/netip/prefix-ops/review` on darwin/arm64. Five independent literal rows matched pinned native Go and compiled Lean. Four mutation witnesses distinguish omitting the unmasked comparison tiebreak, conflating mapped IPv6 with IPv4, equating invalid with the zero Prefix, and accepting zoned candidates. The author replay matched 6,480 boundaries and 60,000 samples at seeds 17, 2718 and 65537, with a per-bit native oracle. [Reviewer control](../../../../difftest/net/netip/prefix-ops/review/main.go)

The modeled states are reachable address and Prefix values and sequential calls.
