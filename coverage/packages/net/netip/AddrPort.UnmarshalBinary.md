# Independent review of `net/netip.AddrPort.UnmarshalBinary`

**Verdict: FAITHFUL**

Evidence SHA-256: `e68ec8efe9a49db72d1002a4c2f8798293cf6a7daf4602d4a5baea82311b1449`

I checked evidence commit `b9ade71` against Go 1.27.0 `net/netip/netip.go`. The final two bytes form the little-endian port. The preceding address must decode successfully before the receiver is replaced. The malformed-length witness preserves the old address and port.

I ran `go run ./difftest/net/netip/binary-codec` on darwin/arm64 and reproduced 2,502 boundaries and 20,000 cases each for seeds 17, 2718, and 65537 against pinned native Go, the native oracle, and compiled Lean. Three independent literal rows also passed each declaration with `go run ./difftest/net/netip/binary-codec/review`. They distinguish staged append writes, size and exact wire bytes, little-endian tails, zone ownership, invalid prefix widths, and atomic decode failures. The scope covers reachable address states, nonnil decode receivers, valid append geometry, and disjoint source and destination. Exact growth capacity and resource-dependent allocation failures are outside scope. [Independent control](../../../../difftest/net/netip/binary-codec/review/main.go)
