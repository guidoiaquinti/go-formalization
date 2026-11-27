# Independent review of `net/netip.Prefix.MarshalText`

**Verdict: FAITHFUL**

Reviewed author commit `33799c2` against pinned Go 1.27.0 `src/net/netip/netip.go`. Prefixes strip IPv6 zones, preserve host bits, and write the decimal prefix length. MarshalText returns a nonnil empty byte slice for the zero form and preserves the pinned invalid-prefix marker for the nonzero invalid form.

Seven literal fixtures passed against native Go and compiled Lean. Emitting the invalid marker for a zero-input marshal changed 2 fixtures. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/net/netip/port-prefix-text/review` from the repository root. The program checks the exact rejection count for each declaration.

Author evidence passed 6,296 boundaries and 60,000 samples at seeds 17, 2718, and 65537. Scope covers reachable public values and feasible append buffers. Exact capacity after growth remains outside scope.
Evidence SHA-256: `94dda1eac5a4a747893f2fe19c2966d7c16ba47364e79bee4de624a63639ea39`
