# Independent review of `net/netip.AddrPort.MarshalText`

**Verdict: FAITHFUL**

Reviewed author commit `33799c2` against pinned Go 1.27.0 `src/net/netip/netip.go`. Addresses keep raw zones and bracket the IPv6 form. Ports use unpadded unsigned decimal. MarshalText returns a nonnil empty byte slice for a zero address even with a nonzero port.

Seven literal fixtures passed against native Go and compiled Lean. Emitting the invalid marker for a zero-input marshal changed 2 fixtures. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/net/netip/port-prefix-text/review` from the repository root. The program checks the exact rejection count for each declaration.

Author evidence passed 6,296 boundaries and 60,000 samples at seeds 17, 2718, and 65537. Scope covers reachable public values and feasible append buffers. Exact capacity after growth remains outside scope.
Evidence SHA-256: `e40543d07569d86aa2398004feb2dec72b9b0ffef49bf85162e86156fcee1253`
