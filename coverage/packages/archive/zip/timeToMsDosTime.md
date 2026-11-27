# Independent review of `archive/zip.timeToMsDosTime`

**Verdict: FAITHFUL**

The source reads local civil fields, including after a zone offset crosses a date boundary. It packs by truncating odd seconds and converting to uint16, so years beyond the 1980–2107 range wrap in the date word. The theorem checks the first odd second. Independent controls include both year wrap directions.

The author scope covers nonnil header receivers, Gregorian UTC decoding, and packing from projected civil fields. Location object internals and monotonic clock implementation details are outside the Lean model. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 3,150 fixed and 60,000 sampled rows at seeds 17, 2718 and 65537 on darwin/arm64. Six independent literal controls passed. Packing the local fields in SetModTime disagrees on all 6, and rounding odd seconds upward disagrees on 4.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/archive/zip/dos-time/review
```

Evidence SHA-256: `0d66a8c81e413d8f3c7defc56e5b9a0b96f54be902c9246569ffd9fe65973e50`
