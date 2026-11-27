# Independent review of `time.Time.Location`

**Verdict: FAITHFUL**

The source returns the stored location as-is or the default UTC singleton when the stored pointer is nil. The theorem covers the nil branch. Native controls check return identity for two distinct locations both named UTC.

The author scope covers raw packed clock words, opaque location identity, nonnil private receivers and default UTC and Local globals. Zone lookup and concurrent mutation are outside scope. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 1,000 fixed cases and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Six independent literal controls and a separate native pointer-identity trace passed. The trace distinguishes canonicalization by name, retaining monotonic ext during stripping, and accepting nil in Time.In.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/location/review
```

Evidence SHA-256: `b83b89d1f4bfb04b532e877703aba2e39c971cba792357b0665e88f936ab4954`
