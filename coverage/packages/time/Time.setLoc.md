# Independent review of `time.Time.setLoc`

**Verdict: FAITHFUL**

The source canonicalizes the exact internal UTC pointer to nil, strips monotonic storage, and writes the resulting location. A nil argument is allowed here. The theorem checks the raw stored UTC tag. A native control uses a separately allocated location named UTC to distinguish pointer equality from name equality.

The author scope covers raw packed clock words, opaque location identity, nonnil private receivers and default UTC and Local globals. Zone lookup and concurrent mutation are outside scope. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 1,000 fixed cases and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Six independent literal controls and a separate native pointer-identity trace passed. The trace distinguishes canonicalization by name, retaining monotonic ext during stripping, and accepting nil in Time.In.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/location/review
```

Evidence SHA-256: `9c308c1015fdb0fd99dd9701ebc03d254ddee88b61bfe506cebaf93a0dffe34c`
