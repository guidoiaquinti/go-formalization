# Independent review of `time.subMono`

**Verdict: FAITHFUL**

The source subtracts signed int64 readings, using the sign of the wrapped difference and the original ordering to detect overflow. Exact endpoints and saturated endpoints share the returned values. The theorem fixes both saturated directions, and native controls exercise exact MinInt64 and positive and negative overflow.

The author scope follows raw signed and packed Go fields, including synthetic words, and represents locations by opaque tags. Pinned native Go checks location identity where applicable. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 5,184 fixed and 60,000 sampled rows at seeds 17, 2718 and 65537 on darwin/arm64. Twelve independent literal fixtures and a separate native boundary transcript passed. The controls distinguish unconditional wall subtraction, unchecked arithmetic on the monotonic word, and skipping the wall round trip.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/clock-arithmetic/review
```

Evidence SHA-256: `371e78ea3aa774964605454aa3887c9f706fefa9bfac0045417c4c29f512e200`
