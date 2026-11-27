# Independent review of `net/netip.AddrFromSlice`

**Verdict: FAITHFUL**

Evidence SHA-256: `cd4492d953bd1cdeb8d5a52499e8a0fee69cb5eb068d24bdc0a47a1a9a21630a`

I compared the Lean model and native overlay with the pinned Go 1.27.0 declaration. Only four-byte and sixteen-byte inputs succeed. The sixteen-byte IPv4-mapped pattern remains an IPv6 address. The native comparison also checks that changing the input slice after construction does not affect the value.

I ran `go run ./difftest/net/netip/addr-fields/review` on darwin/arm64. Five independent literal rows matched pinned native Go and compiled Lean. The checks distinguish premature index rejection, unmapping a sixteen-byte address, retaining an IPv6 zone, and treating a valid all-zero IPv6 address as invalid. The author comparison passed 2,308 fixed cases and 60,000 sampled cases at seeds 17, 2718 and 65537. [Reviewer control](../../../../difftest/net/netip/addr-fields/review/main.go)

The domain contains reachable address states and sequential calls. The tests use raw byte zones and check the value returned by each operation.
