# Independent review of `net/netip.AddrPort.MarshalBinary`

**Verdict: FAITHFUL**

Evidence SHA-256: `0ae1e1ce0f3435c4a24a01add0203b14c0447691808f01abaeeacb98ea45320f`

I checked evidence commit `b9ade71` against Go 1.27.0 `net/netip/netip.go`. The port follows the encoded address as a little-endian uint16. The 0x1234 fixture must end in bytes 52,18 and rejects writing the reversed order.

I ran `go run ./difftest/net/netip/binary-codec` on darwin/arm64 and reproduced 2,502 boundaries and 20,000 cases each for seeds 17, 2718, and 65537 against pinned native Go, the native oracle, and compiled Lean. Three independent literal rows also passed each declaration with `go run ./difftest/net/netip/binary-codec/review`. They distinguish staged append writes, size and exact wire bytes, little-endian tails, zone ownership, invalid prefix widths, and atomic decode failures. The scope covers reachable address states, nonnil decode receivers, valid append geometry, and disjoint source and destination. Exact growth capacity and resource-dependent allocation failures are outside scope. [Independent control](../../../../difftest/net/netip/binary-codec/review/main.go)
