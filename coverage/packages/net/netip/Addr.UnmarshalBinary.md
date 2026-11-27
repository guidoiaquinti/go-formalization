# Independent review of `net/netip.Addr.UnmarshalBinary`

**Verdict: FAITHFUL**

Evidence SHA-256: `36649b3b476ff3b770f8ad962ea8208fe3a396fbd20abb080f30898491784fe4`

I checked evidence commit `b9ade71` against Go 1.27.0 `net/netip/netip.go`. A structural length error leaves a populated receiver intact. A successful long IPv6 payload uses the first sixteen bytes for the address and the tail for the raw zone. Native clears the original input before reading the zone back, which detects borrowed zone storage.

I ran `go run ./difftest/net/netip/binary-codec` on darwin/arm64 and reproduced 2,502 boundaries and 20,000 cases each for seeds 17, 2718, and 65537 against pinned native Go, the native oracle, and compiled Lean. Three independent literal rows also passed each declaration with `go run ./difftest/net/netip/binary-codec/review`. They distinguish staged append writes, size and exact wire bytes, little-endian tails, zone ownership, invalid prefix widths, and atomic decode failures. The scope covers reachable address states, nonnil decode receivers, valid append geometry, and disjoint source and destination. Exact growth capacity and resource-dependent allocation failures are outside scope. [Independent control](../../../../difftest/net/netip/binary-codec/review/main.go)
