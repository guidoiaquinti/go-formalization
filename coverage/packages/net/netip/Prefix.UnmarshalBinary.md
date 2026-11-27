# Independent review of `net/netip.Prefix.UnmarshalBinary`

**Verdict: FAITHFUL**

Evidence SHA-256: `dc5b97cea51a162fa8567702ea56ecff2f9b3ea0cdd93aed2ccb30fab21f7abd`

I checked evidence commit `b9ade71` against Go 1.27.0 `net/netip/netip.go`. After structural decoding, PrefixFrom strips the zone and accepts an initialized address whose width is invalid, with the exposed width minus one. The fixture distinguishes this case from a structural length error, which must preserve the old receiver.

I ran `go run ./difftest/net/netip/binary-codec` on darwin/arm64 and reproduced 2,502 boundaries and 20,000 cases each for seeds 17, 2718, and 65537 against pinned native Go, the native oracle, and compiled Lean. Three independent literal rows also passed each declaration with `go run ./difftest/net/netip/binary-codec/review`. They distinguish staged append writes, size and exact wire bytes, little-endian tails, zone ownership, invalid prefix widths, and atomic decode failures. The scope covers reachable address states, nonnil decode receivers, valid append geometry, and disjoint source and destination. Exact growth capacity and resource-dependent allocation failures are outside scope. [Independent control](../../../../difftest/net/netip/binary-codec/review/main.go)
