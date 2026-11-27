# Independent review of `time.Time.Local`

**Verdict: FAITHFUL**

The source uses the default Local pointer and strips monotonic storage. The theorem proves retention of its location tag. Native control checks raw pointer equality with time.Local under the pinned default globals.

The author scope covers raw packed clock words, opaque location identity, nonnil private receivers and default UTC and Local globals. Zone lookup and concurrent mutation are outside scope. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 1,000 fixed cases and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Six independent literal controls and a separate native pointer-identity trace passed. The trace distinguishes canonicalization by name, retaining monotonic ext during stripping, and accepting nil in Time.In.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/location/review
```

Evidence SHA-256: `6478c153ad3e79edc677d8c8294b2ce6c3f9b0dcc157288004d548040d174d54`
