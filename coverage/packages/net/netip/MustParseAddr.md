# Independent review of `net/netip.MustParseAddr`

**Verdict: FAITHFUL**

Reviewed author commit `0dfd0c8` against pinned Go 1.27.0. Go panics with the same structured error returned by the parser. The overlay checks the concrete panic type and fields. The empty-input control distinguishes the panic from an ordinary error result.

Pinned Go and Lean agreed on 374 fixed cases and 60,000 sampled cases at seeds 17, 2718 and 65537. This includes 48,000 constructed addresses with independent expected bytes, of which 36,000 are IPv6. The scoped Lean kernel audit found no axioms for the four registered theorems. The independent control ran four rows and rejected five deliberately wrong outcomes across the batch. Run it with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/net/netip/parse-address/review` from the repository root. The unmarshal receiver is nonnil. Zone interning and formatted diagnostic strings are outside the claims.

[Reviewer control](../../../../difftest/net/netip/parse-address/review/main.go)

Evidence SHA-256: `ede419014661dcb55fa5c591412486819c71a9298c8c5df322e018e1e0c8055d`
