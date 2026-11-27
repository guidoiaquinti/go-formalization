# Independent review of `net/netip.Addr.appendTo4`

**Verdict: FAITHFUL**

Evidence SHA-256: `e53a6c46633febb2b7d71d8155c6bd60150aad4964035a517a851666ee3df00d`

I reviewed evidence commit `f91be9c` against pinned Go 1.27.0 `net/netip/netip.go`. The source appends decimal digits individually. With one spare byte before detachment, the pinned private call writes the first digit of 255 into the old backing store.

The direct comparison on darwin/arm64 checked 3,584 fixed and 60,000 sampled observations on the declared address family using seeds 17, 2718 and 65537. The independent control checks four literal rows against pinned Go and compiled Lean and rejects six deliberately incorrect observations. Scope covers reachable addresses and feasible append buffers. Exact growth capacity and allocation failures are outside scope. [Independent control](../../../../difftest/net/netip/addr-text/review/main.go)
