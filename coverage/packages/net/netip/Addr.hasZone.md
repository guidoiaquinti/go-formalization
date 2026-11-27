# Independent review of `net/netip.Addr.hasZone`

**Verdict: FAITHFUL**

Evidence SHA-256: `7ee830125e664a9f2f46061426b0ba10995af183463861cc0c76cdd26f4d91d4`

I compared the Lean model and native overlay with the pinned Go 1.27.0 declaration. A zone is present only for an IPv6 address with the zoned detail. The test includes arbitrary raw zone bytes.

I ran `go run ./difftest/net/netip/addr-fields/review` on darwin/arm64. Five independent literal rows matched pinned native Go and compiled Lean. The checks distinguish premature index rejection, unmapping a sixteen-byte address, retaining an IPv6 zone, and treating a valid all-zero IPv6 address as invalid. The author comparison passed 2,308 fixed cases and 60,000 sampled cases at seeds 17, 2718 and 65537. [Reviewer control](../../../../difftest/net/netip/addr-fields/review/main.go)

The domain contains reachable address states and sequential calls. The tests use raw byte zones and check the value returned by each operation.
