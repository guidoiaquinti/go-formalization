# Independent review of `net/netip.parseIPv6`

**Verdict: FAITHFUL**

Reviewed author commit `0dfd0c8` against pinned Go 1.27.0. The eight-group bound follows from consuming two bytes per hex group or four for the terminal IPv4 fields. Four hexadecimal digits cannot exceed 65535, so the Go accumulator check cannot independently fail. The control verifies that a bare double colon produces sixteen zero bytes.

Pinned Go and Lean agreed on 374 fixed cases and 60,000 sampled cases at seeds 17, 2718 and 65537. This includes 48,000 constructed addresses with independent expected bytes, of which 36,000 are IPv6. The scoped Lean kernel audit found no axioms for the four registered theorems. The independent control ran four rows and rejected five deliberately wrong outcomes across the batch. Run it with `GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/net/netip/parse-address/review` from the repository root. The unmarshal receiver is nonnil. Zone interning and formatted diagnostic strings are outside the claims.

[Reviewer control](../../../../difftest/net/netip/parse-address/review/main.go)

Evidence SHA-256: `a8ef7ffaec25cd0111716d96a8f2f6f030e9309cf1e2ba8fe76a35f4ffb8b85d`
