# Independent review of `net/netip.Addr.AppendTo`

**Verdict: FAITHFUL**

Evidence SHA-256: `8dbf34ea7f9c38b9f7d73112c8304a0e19d2b5d2cedb0a35db83ae43e0b0c859`

I reviewed evidence commit `f91be9c` against pinned Go 1.27.0 `net/netip/netip.go`. The source dispatches to the formatter for the correct address family and returns the original destination for invalid addresses. The controls distinguish old backing writes for individual decimal digits and all-or-nothing fixed chunks.

The direct comparison on darwin/arm64 checked 14,336 fixed and 195,000 sampled observations across the four address profiles using seeds 17, 2718 and 65537. The independent control checks four literal rows against pinned Go and compiled Lean and rejects six deliberately incorrect observations. Scope covers reachable addresses and feasible append buffers. Exact growth capacity and allocation failures are outside scope. [Independent control](../../../../difftest/net/netip/addr-text/review/main.go)
