# Independent review of `archive/zip.timeZone`

**Verdict: FAITHFUL**

Pinned Go rounds the signed duration to a quarter hour before testing the inclusive -12-hour and +14-hour limits. Rounding uses ties away from zero and saturates at the signed duration endpoints. An out-of-range rounded result becomes zero before constructing an unnamed fixed zone. The pinned `time.FixedZone` implementation caches each unnamed whole-hour location in the resulting range.

The model observes offset seconds and cache reuse for all signed duration words. The native overlay checks the unnamed fixed offset, absence of daylight saving and open-ended bounds at four dates. The two theorems check the upper round-before-limit transition and the negative half-quarter tie. The scoped Lean axiom audit passed.

Pinned Go 1.27.0 and Lean passed 1,129 fixed cases and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Fourteen independent controls passed. A clamp-before-rounding mutation disagrees just outside both limits, on 2 rows. Additional native identity checks verify that distinct input durations rounding to the same whole hour share the standard-library cache, out-of-range durations share the unnamed zero cache distinct from UTC, and identical fractional-hour results use distinct pointers.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/archive/zip/time-zone/review
```

Evidence SHA-256: `0c147659a4c76bacf231c9b0ab1fc4ae15ef3bf50663079dde8c5353bbb9aaa9`
