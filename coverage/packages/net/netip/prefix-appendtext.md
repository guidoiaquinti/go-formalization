# Independent review of `net/netip.Prefix.AppendText`

**Verdict: FAITHFUL**

Reviewed author commit `33799c2` against pinned Go 1.27.0 `src/net/netip/netip.go`. Prefixes strip IPv6 zones, preserve host bits, and write the decimal prefix length. Native checks record destination alias, nilness, and the original bytes written before a later stage grows the allocation.

Seven literal fixtures passed against native Go and compiled Lean. Treating the append as a single atomic chunk changed 2 fixtures. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/net/netip/port-prefix-text/review` from the repository root. The program checks the exact rejection count for each declaration.

Author evidence passed 6,296 boundaries and 60,000 samples at seeds 17, 2718, and 65537. Scope covers reachable public values and feasible append buffers. Exact capacity after growth remains outside scope.
Evidence SHA-256: `0ddb3f4a7e5cd9e95d10736624b651c3d59c550d4b97b2b034660025c8388d1c`
