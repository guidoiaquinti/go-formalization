# Independent review of `net/netip.AddrPort.AppendTo`

**Verdict: FAITHFUL**

Reviewed author commit `33799c2` against pinned Go 1.27.0 `src/net/netip/netip.go`. Addresses keep raw zones and bracket the IPv6 form. Ports use unpadded unsigned decimal. Native checks record destination alias, nilness, and the original bytes written before a later stage grows the allocation.

Seven literal fixtures passed against native Go and compiled Lean. Treating the append as a single atomic chunk changed 4 fixtures. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/net/netip/port-prefix-text/review` from the repository root. The program checks the exact rejection count for each declaration.

Author evidence passed 6,296 boundaries and 60,000 samples at seeds 17, 2718, and 65537. Scope covers reachable public values and feasible append buffers. Exact capacity after growth remains outside scope.
Evidence SHA-256: `ba60fc91bad32a537b4d28aae71f33ce721bfa90b7dce1ed00ec9609321dba25`
