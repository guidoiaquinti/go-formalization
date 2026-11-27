# Independent review of `net/netip.Prefix.AppendBinary`

**Verdict: FAITHFUL**

Evidence SHA-256: `65ea93e9ec26810dd442b5f2028d72ee0f3dfca5782da79cb40ffa46704c4fbe`

I checked evidence commit `b9ade71` against Go 1.27.0 `net/netip/netip.go`. Each source append is a separate write stage. A prefix plus eight spare bytes receives the IPv6 high word before the low-word append detaches. A four-byte IPv4 append can fill old capacity before the port or prefix tail detaches. The literal old-backing observations reject preallocating all output before the first source write.

I ran `go run ./difftest/net/netip/binary-codec` on darwin/arm64 and reproduced 2,502 boundaries and 20,000 cases each for seeds 17, 2718, and 65537 against pinned native Go, the native oracle, and compiled Lean. Three independent literal rows also passed each declaration with `go run ./difftest/net/netip/binary-codec/review`. They distinguish staged append writes, size and exact wire bytes, little-endian tails, zone ownership, invalid prefix widths, and atomic decode failures. The scope covers reachable address states, nonnil decode receivers, valid append geometry, and disjoint source and destination. Exact growth capacity and resource-dependent allocation failures are outside scope. [Independent control](../../../../difftest/net/netip/binary-codec/review/main.go)
