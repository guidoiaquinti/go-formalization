# Independent review of `net/netip.Addr.v6u16`

**Verdict: FAITHFUL**

Evidence SHA-256: `19689542287db59de3a0c5595ce21b2d117a113a1e57f362307d435a3935dfe8`

I compared the Lean model and native overlay with the pinned Go 1.27.0 declaration. The word selector alternates the 64-bit limb every four indices and selects the 16-bit network-order word.

I ran `go run ./difftest/net/netip/addr-fields/review` on darwin/arm64. Five independent literal rows matched pinned native Go and compiled Lean. The checks distinguish premature index rejection, unmapping a sixteen-byte address, retaining an IPv6 zone, and treating a valid all-zero IPv6 address as invalid. The author comparison passed 2,308 fixed cases and 60,000 sampled cases at seeds 17, 2718 and 65537. [Reviewer control](../../../../difftest/net/netip/addr-fields/review/main.go)

The domain contains reachable address states and sequential calls. The tests use raw byte zones and check the value returned by each operation.
