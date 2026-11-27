# Independent review of `net/netip.Addr.MarshalText`

**Verdict: FAITHFUL**

Evidence SHA-256: `122fa70617208da6194df605783af35c00092bdd7caa92690db464ad2af26bd9`

I reviewed evidence commit `f91be9c` against pinned Go 1.27.0 `net/netip/netip.go`. The source returns compact address text and a nil error. The pinned native control distinguishes the invalid address's nonnil empty result from its String text and a nil append destination.

The direct comparison on darwin/arm64 checked 14,336 fixed and 195,000 sampled observations across the four address profiles using seeds 17, 2718 and 65537. The independent control checks four literal rows against pinned Go and compiled Lean and rejects six deliberately incorrect observations. Scope covers reachable addresses and feasible append buffers. Exact growth capacity and allocation failures are outside scope. [Independent control](../../../../difftest/net/netip/addr-text/review/main.go)
