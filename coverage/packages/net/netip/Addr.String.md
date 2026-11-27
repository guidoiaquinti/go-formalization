# Independent review of `net/netip.Addr.String`

**Verdict: FAITHFUL**

Evidence SHA-256: `e192d593ae9c69d883401f98b8879da793007359697656b412997c80406d2ff2`

I reviewed evidence commit `f91be9c` against pinned Go 1.27.0 `net/netip/netip.go`. The source splits invalid, IPv4, mapped IPv6 and ordinary IPv6. The controls distinguish invalid IP text, raw mapped zone bytes and earliest longest-run compaction.

The direct comparison on darwin/arm64 checked 14,336 fixed and 195,000 sampled observations across the four address profiles using seeds 17, 2718 and 65537. The independent control checks four literal rows against pinned Go and compiled Lean and rejects six deliberately incorrect observations. Scope covers reachable addresses and feasible append buffers. Exact growth capacity and allocation failures are outside scope. [Independent control](../../../../difftest/net/netip/addr-text/review/main.go)
