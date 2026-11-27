# Independent review of `time.Duration.format`

**Verdict: FAITHFUL**

The source writes the units, fraction, digits and sign backward into a fixed 32-byte array. The resulting prefix outside the written suffix is untouched, and the offset counts bytes. Nil raises a runtime nil-pointer failure on the first write. The theorem states the exact bytes and retained prefix for zero.

The author scope covers signed int64 durations and fixed-size or nil format buffers. The Lean model summarizes backward intermediate writes by their final suffix and abstracts the nil failure classification. Allocation and panic diagnostics are outside scope. The scoped Lean axiom audit passed.

Pinned Go 1.27.0 and Lean passed 254 fixed cases and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Fourteen independent literal controls passed, including both signed endpoints, mixed-unit boundaries and nonuniform buffer prefixes. Substituting Greek mu for the micro sign disagrees on 4 rows, retaining full precision after a decimal point disagrees on 3, and clearing the prefix disagrees on all 12 nonnil buffer rows.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/duration-text/review
```

Evidence SHA-256: `8b2d6f3dd87e110dd0c0330e7d8a0e18ccea6e902584a8f417c04e79a04e7cbe`
