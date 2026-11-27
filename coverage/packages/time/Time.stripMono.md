# Independent review of `time.Time.stripMono`

**Verdict: FAITHFUL**

The source decodes wall seconds before clearing the monotonic marker and seconds field. The low 30 nanosecond bits survive and ext receives the wall seconds. The theorem establishes exact retention for unmarked words, including an unmarked synthetic value with otherwise unused bits.

The author scope covers raw packed clock words, opaque location identity, nonnil private receivers and default UTC and Local globals. Zone lookup and concurrent mutation are outside scope. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 1,000 fixed cases and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Six independent literal controls and a separate native pointer-identity trace passed. The trace distinguishes canonicalization by name, retaining monotonic ext during stripping, and accepting nil in Time.In.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/location/review
```

Evidence SHA-256: `961aa7fa65bbc4d1afcfb852a44336ed118fb4c0c5bc916945a40e21a5499eef`
