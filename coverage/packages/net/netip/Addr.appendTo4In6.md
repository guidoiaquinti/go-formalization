# Independent review of `net/netip.Addr.appendTo4In6`

**Verdict: FAITHFUL**

Evidence SHA-256: `5804629a2555d5a012039675e69b945ad85344e1c6dbc8899b3f698d0e859fea`

I reviewed evidence commit `f91be9c` against pinned Go 1.27.0 `net/netip/netip.go`. The source appends the seven-byte mapped marker as a single operation. The pinned private call with six spare bytes keeps the old backing store untouched and preserves a malformed byte in the zone.

The direct comparison on darwin/arm64 checked 3,584 fixed and 60,000 sampled observations on the declared address family using seeds 17, 2718 and 65537. The independent control checks four literal rows against pinned Go and compiled Lean and rejects six deliberately incorrect observations. Scope covers reachable addresses and feasible append buffers. Exact growth capacity and allocation failures are outside scope. [Independent control](../../../../difftest/net/netip/addr-text/review/main.go)
