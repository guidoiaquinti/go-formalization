# Independent review of `net/netip.Addr.StringExpanded`

**Verdict: FAITHFUL**

Evidence SHA-256: `fe37901a47415cf1801e3e7792cfd0bc07179270b2077a2f098e708e4bce7eef`

I reviewed evidence commit `f91be9c` against pinned Go 1.27.0 `net/netip/netip.go`. The source emits eight zero-padded hexadecimal words for IPv6 including mapped IPv6. The controls reject substituting compact notation and verify the raw mapped zone suffix.

The direct comparison on darwin/arm64 checked 14,336 fixed and 195,000 sampled observations across the four address profiles using seeds 17, 2718 and 65537. The independent control checks four literal rows against pinned Go and compiled Lean and rejects six deliberately incorrect observations. Scope covers reachable addresses and feasible append buffers. Exact growth capacity and allocation failures are outside scope. [Independent control](../../../../difftest/net/netip/addr-text/review/main.go)
