# Independent review of `net/netip.Addr.AppendText`

**Verdict: FAITHFUL**

Evidence SHA-256: `409179b2bd3deda57e31e66cda30342aa72bc37b4a0e39d6e066c79f319e580b`

I reviewed evidence commit `f91be9c` against pinned Go 1.27.0 `net/netip/netip.go`. The source returns the AppendTo result with a nil error. The controls check exact append bytes, preserved invalid-input destination and alias effects.

The direct comparison on darwin/arm64 checked 14,336 fixed and 195,000 sampled observations across the four address profiles using seeds 17, 2718 and 65537. The independent control checks four literal rows against pinned Go and compiled Lean and rejects six deliberately incorrect observations. Scope covers reachable addresses and feasible append buffers. Exact growth capacity and allocation failures are outside scope. [Independent control](../../../../difftest/net/netip/addr-text/review/main.go)
