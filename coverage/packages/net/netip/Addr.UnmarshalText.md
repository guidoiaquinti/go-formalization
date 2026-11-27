# Independent review of `net/netip.Addr.UnmarshalText`

**Verdict: FAITHFUL**

Reviewed author commit `0dfd0c8` against pinned Go 1.27.0. Go writes the zero address on empty or invalid input and replaces the receiver on success. The overlay starts with a nonzero zoned receiver. The empty-input control distinguishes clearing successfully from returning a parse error.

Pinned Go and Lean agreed on 374 fixed cases and 60,000 sampled cases at seeds 17, 2718 and 65537. This includes 48,000 constructed addresses with independent expected bytes, of which 36,000 are IPv6. The scoped Lean kernel audit found no axioms for the four registered theorems. The independent control ran four rows and rejected five deliberately wrong outcomes across the batch. Run it with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/net/netip/parse-address/review` from the repository root. The unmarshal receiver is nonnil. Zone interning and formatted diagnostic strings are outside the claims.

[Reviewer control](../../../../difftest/net/netip/parse-address/review/main.go)

Evidence SHA-256: `5b999034a6bed53719d703f0a33c5d3a4043856700183d93577076755e8b0556`
