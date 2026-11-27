# Independent review of `net/netip.Addr.withoutZone`

**Verdict: FAITHFUL**

Evidence SHA-256: `52fb9ac2b4e9016f6d72df5a6ade649098cf0fb3fce8989585881146d4fef9ff`

I compared the Lean model and native overlay with the pinned Go 1.27.0 declaration. The helper leaves family and both limbs intact and clears the zone only for IPv6, returning an equivalent unzoned IPv6 address.

I ran `go run ./difftest/net/netip/addr-fields/review` on darwin/arm64. Five independent literal rows matched pinned native Go and compiled Lean. The checks distinguish premature index rejection, unmapping a sixteen-byte address, retaining an IPv6 zone, and treating a valid all-zero IPv6 address as invalid. The author comparison passed 2,308 fixed cases and 60,000 sampled cases at seeds 17, 2718 and 65537. [Reviewer control](../../../../difftest/net/netip/addr-fields/review/main.go)

The domain contains reachable address states and sequential calls. The tests use raw byte zones and check the value returned by each operation.
