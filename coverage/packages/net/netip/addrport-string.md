# Independent review of `net/netip.AddrPort.String`

**Verdict: FAITHFUL**

Reviewed author commit `33799c2` against pinned Go 1.27.0 `src/net/netip/netip.go`. Addresses keep raw zones and bracket the IPv6 form. Ports use unpadded unsigned decimal. The string distinguishes an invalid value with an explicit marker.

Seven literal fixtures passed against native Go and compiled Lean. Omitting IPv6 brackets changed 3 fixtures. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/net/netip/port-prefix-text/review` from the repository root. The program checks the exact rejection count for each declaration.

Author evidence passed 6,296 boundaries and 60,000 samples at seeds 17, 2718, and 65537. Scope covers reachable public values and feasible append buffers. Exact capacity after growth remains outside scope.
Evidence SHA-256: `29cfb182b12bef83bb42b5b5df14a5c2f987fbe352b7c4924c5f710d18719469`
