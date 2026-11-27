# Independent review of `net/netip.Addr.v4`

**Verdict: FAITHFUL**

Evidence SHA-256: `76dca2ed828720879ab2dc2029a0fa0259994aac184bec64b4aa1047c6889e23`

I compared the Lean model and native overlay with the pinned Go 1.27.0 declaration. The pinned uint8 shift arithmetic wraps before the shift. Index 32 therefore selects the first byte of the low IPv4 word, which distinguishes ordinary unbounded out-of-range masking.

I ran `go run ./difftest/net/netip/addr-fields/review` on darwin/arm64. Five independent literal rows matched pinned native Go and compiled Lean. The checks distinguish premature index rejection, unmapping a sixteen-byte address, retaining an IPv6 zone, and treating a valid all-zero IPv6 address as invalid. The author comparison passed 2,308 fixed cases and 60,000 sampled cases at seeds 17, 2718 and 65537. [Reviewer control](../../../../difftest/net/netip/addr-fields/review/main.go)

The domain contains reachable address states and sequential calls. The tests use raw byte zones and check the value returned by each operation.
