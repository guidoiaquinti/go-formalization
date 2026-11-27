# Independent review of `time.Time.UTC`

**Verdict: FAITHFUL**

The source passes the internal UTC pointer to setLoc, giving raw nil storage and a public UTC location. The native control checks the raw storage and stripped words as well as the original receiver remaining intact.

The author scope covers raw packed clock words, opaque location identity, nonnil private receivers and default UTC and Local globals. Zone lookup and concurrent mutation are outside scope. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 1,000 fixed cases and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Six independent literal controls and a separate native pointer-identity trace passed. The trace distinguishes canonicalization by name, retaining monotonic ext during stripping, and accepting nil in Time.In.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/location/review
```

Evidence SHA-256: `681c8778824edaad81ac279cc08f2546b47044f145b655323dc41c8bf9ff70ee`
