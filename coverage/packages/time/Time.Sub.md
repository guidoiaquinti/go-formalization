# Independent review of `time.Time.Sub`

**Verdict: FAITHFUL**

The source uses monotonic subtraction only when both operands are marked. Otherwise it forms a wrapped wall candidate and checks it by adding it back to the second operand. The native transcript distinguishes exact MinInt64 and MaxInt64 from the adjacent clipped durations by round-trip equality. A synthetic unnormalized wall nanosecond field produces a nonroundtripping result even though a naive linear subtraction is representable.

The author scope follows raw signed and packed Go fields, including synthetic words, and represents locations by opaque tags. Pinned native Go checks location identity where applicable. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 5,184 fixed and 60,000 sampled rows at seeds 17, 2718 and 65537 on darwin/arm64. Twelve independent literal fixtures and a separate native boundary transcript passed. The controls distinguish unconditional wall subtraction, unchecked arithmetic on the monotonic word, and skipping the wall round trip.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/clock-arithmetic/review
```

Evidence SHA-256: `14b24156c3b257ac6a011952cbe6a17a865fd5bc46920b16a401aebc2aed1464`
