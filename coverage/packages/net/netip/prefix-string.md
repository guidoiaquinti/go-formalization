# Independent review of `net/netip.Prefix.String`

**Verdict: FAITHFUL**

Reviewed author commit `33799c2` against pinned Go 1.27.0 `src/net/netip/netip.go`. Prefixes strip IPv6 zones, preserve host bits, and write the decimal prefix length. The string distinguishes an invalid value with an explicit marker.

Seven literal fixtures passed against native Go and compiled Lean. Retaining zones on prefixes changed 2 fixtures. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/net/netip/port-prefix-text/review` from the repository root. The program checks the exact rejection count for each declaration.

Author evidence passed 6,296 boundaries and 60,000 samples at seeds 17, 2718, and 65537. Scope covers reachable public values and feasible append buffers. Exact capacity after growth remains outside scope.
Evidence SHA-256: `8a49608452fe6cde33e64392a2c1618879cbb3745f0685d54f5cce67f8924339`
