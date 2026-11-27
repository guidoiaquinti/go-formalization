# Independent review of `net/netip.Addr.marshalBinarySize`

**Verdict: FAITHFUL**

Evidence SHA-256: `f50c1d210ccdf4d46245a7fc568b4a68e0db99b35d697230a150eb5b9a63a857`

I checked evidence commit `b9ade71` against Go 1.27.0 `net/netip/netip.go`. The size is zero for an invalid Addr, four for IPv4, and sixteen plus the raw IPv6 zone length. The fixed two-byte zone rejects returning a constant sixteen for IPv6.

I ran `go run ./difftest/net/netip/binary-codec` on darwin/arm64 and reproduced 2,502 boundaries and 20,000 cases each for seeds 17, 2718, and 65537 against pinned native Go, the native oracle, and compiled Lean. Three independent literal rows also passed each declaration with `go run ./difftest/net/netip/binary-codec/review`. They distinguish staged append writes, size and exact wire bytes, little-endian tails, zone ownership, invalid prefix widths, and atomic decode failures. The scope covers reachable address states, nonnil decode receivers, valid append geometry, and disjoint source and destination. Exact growth capacity and resource-dependent allocation failures are outside scope. [Independent control](../../../../difftest/net/netip/binary-codec/review/main.go)
