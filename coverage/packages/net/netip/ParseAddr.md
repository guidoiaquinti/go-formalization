# Independent review of `net/netip.ParseAddr`

**Verdict: FAITHFUL**

Reviewed author commit `0dfd0c8` against pinned Go 1.27.0. Go dispatches at the first dot, colon or percent byte. The control rejects interpreting a leading percent as an IPv6 zone. The mapped address case verifies that embedded IPv4 bytes remain in an IPv6 address.

Pinned Go and Lean agreed on 374 fixed cases and 60,000 sampled cases at seeds 17, 2718 and 65537. This includes 48,000 constructed addresses with independent expected bytes, of which 36,000 are IPv6. The scoped Lean kernel audit found no axioms for the four registered theorems. The independent control ran four rows and rejected five deliberately wrong outcomes across the batch. Run it with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/net/netip/parse-address/review` from the repository root. The unmarshal receiver is nonnil. Zone interning and formatted diagnostic strings are outside the claims.

[Reviewer control](../../../../difftest/net/netip/parse-address/review/main.go)

Evidence SHA-256: `dc970528225c1fa6cdddd434126a531896f535d8ccfec0caa3802d8b98ebb118`
