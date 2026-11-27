# Independent review of `net/netip.Addr.MarshalBinary`

**Verdict: FAITHFUL**

Evidence SHA-256: `072748324c4dfb11cc29eea5ff9c34a5ccc6aa7abb82a0b831c8b0756a7589b3`

I checked evidence commit `b9ade71` against Go 1.27.0 `net/netip/netip.go`. The source allocates a nonnil exact-capacity buffer and appends network-order address words followed by raw zone bytes. Native checks distinguish empty marshal output from a zero-address nil append.

I ran `go run ./difftest/net/netip/binary-codec` on darwin/arm64 and reproduced 2,502 boundaries and 20,000 cases each for seeds 17, 2718, and 65537 against pinned native Go, the native oracle, and compiled Lean. Three independent literal rows also passed each declaration with `go run ./difftest/net/netip/binary-codec/review`. They distinguish staged append writes, size and exact wire bytes, little-endian tails, zone ownership, invalid prefix widths, and atomic decode failures. The scope covers reachable address states, nonnil decode receivers, valid append geometry, and disjoint source and destination. Exact growth capacity and resource-dependent allocation failures are outside scope. [Independent control](../../../../difftest/net/netip/binary-codec/review/main.go)
