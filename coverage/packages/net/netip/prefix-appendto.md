# Independent review of `net/netip.Prefix.AppendTo`

**Verdict: FAITHFUL**

Reviewed author commit `33799c2` against pinned Go 1.27.0 `src/net/netip/netip.go`. Prefixes strip IPv6 zones, preserve host bits, and write the decimal prefix length. Native checks record destination alias, nilness, and the original bytes written before a later stage grows the allocation.

Seven literal fixtures passed against native Go and compiled Lean. Treating the append as a single atomic chunk changed 2 fixtures. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/net/netip/port-prefix-text/review` from the repository root. The program checks the exact rejection count for each declaration.

Author evidence passed 6,296 boundaries and 60,000 samples at seeds 17, 2718, and 65537. Scope covers reachable public values and feasible append buffers. Exact capacity after growth remains outside scope.
Evidence SHA-256: `4fee5b936aed0b5ad181a712c54862e4a5590fafa74106449d6dff2bb9167bec`
