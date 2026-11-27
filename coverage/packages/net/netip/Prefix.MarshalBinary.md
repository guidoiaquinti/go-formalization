# Independent review of `net/netip.Prefix.MarshalBinary`

**Verdict: FAITHFUL**

Evidence SHA-256: `6ad004e76db91337d35f32cd96e8ea5c072448ce4ce185d0243e5c717a96bbe7`

I checked evidence commit `b9ade71` against Go 1.27.0 `net/netip/netip.go`. Prefix encoding strips the source zone and places the reported prefix width in one trailing byte. Invalid width is 255. The IPv6 fixture rejects appending the original zone before that width.

I ran `go run ./difftest/net/netip/binary-codec` on darwin/arm64 and reproduced 2,502 boundaries and 20,000 cases each for seeds 17, 2718, and 65537 against pinned native Go, the native oracle, and compiled Lean. Three independent literal rows also passed each declaration with `go run ./difftest/net/netip/binary-codec/review`. They distinguish staged append writes, size and exact wire bytes, little-endian tails, zone ownership, invalid prefix widths, and atomic decode failures. The scope covers reachable address states, nonnil decode receivers, valid append geometry, and disjoint source and destination. Exact growth capacity and resource-dependent allocation failures are outside scope. [Independent control](../../../../difftest/net/netip/binary-codec/review/main.go)
