# Independent review of `net/netip.Addr.appendTo6`

**Verdict: FAITHFUL**

Evidence SHA-256: `aa7f31693b50531bf4753fa07de3478df278e6e8b9e6adf12a9f6b80d709a2d1`

I reviewed evidence commit `f91be9c` against pinned Go 1.27.0 `net/netip/netip.go`. The source compresses the earliest longest run of at least two zero words. The pinned private control distinguishes equal-length runs and checks that a two-byte double colon append preserves the old backing store when only one byte remains.

The direct comparison on darwin/arm64 checked 3,584 fixed and 60,000 sampled observations on the declared address family using seeds 17, 2718 and 65537. The independent control checks four literal rows against pinned Go and compiled Lean and rejects six deliberately incorrect observations. Scope covers reachable addresses and feasible append buffers. Exact growth capacity and allocation failures are outside scope. [Independent control](../../../../difftest/net/netip/addr-text/review/main.go)
