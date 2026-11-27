# Independent review of `time.Duration.String`

**Verdict: FAITHFUL**

The source uses unsigned magnitude so the minimum int64 formats without overflow. It selects nanoseconds, microseconds, milliseconds or seconds by magnitude. Fractional digits keep leading zeros and trim trailing zeros, and the source uses U+00B5 for microseconds. Above a minute the source retains the minute and second fields, including zeros.

The author scope covers signed int64 durations and fixed-size or nil format buffers. The Lean model summarizes backward intermediate writes by their final suffix and abstracts the nil failure classification. Allocation and panic diagnostics are outside scope. The scoped Lean axiom audit passed.

Pinned Go 1.27.0 and Lean passed 254 fixed cases and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Fourteen independent literal controls passed, including both signed endpoints, mixed-unit boundaries and nonuniform buffer prefixes. Substituting Greek mu for the micro sign disagrees on 4 rows, retaining full precision after a decimal point disagrees on 3, and clearing the prefix disagrees on all 12 nonnil buffer rows.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/duration-text/review
```

Evidence SHA-256: `aea8bc35496fcdd7dad2642e0a3d1f37d8809d521c866b90d0c66ee1af171f3c`
