# Independent review of `net/netip.Addr.isZero`

**Verdict: FAITHFUL**

Evidence SHA-256: `df82a9c748f967fb4d26db12cc92297689efeb92ee5a95ecf6be8c5e2e4d8f59`

I compared the Lean model and native overlay with the pinned Go 1.27.0 declaration. The invalid detail marker determines zero-address status. An initialized IPv6 address with both limbs zero remains valid.

I ran `go run ./difftest/net/netip/addr-fields/review` on darwin/arm64. Five independent literal rows matched pinned native Go and compiled Lean. The checks distinguish premature index rejection, unmapping a sixteen-byte address, retaining an IPv6 zone, and treating a valid all-zero IPv6 address as invalid. The author comparison passed 2,308 fixed cases and 60,000 sampled cases at seeds 17, 2718 and 65537. [Reviewer control](../../../../difftest/net/netip/addr-fields/review/main.go)

The domain contains reachable address states and sequential calls. The tests use raw byte zones and check the value returned by each operation.
