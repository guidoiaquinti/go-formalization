# Independent review of `archive/zip.msDosTimeToTime`

**Verdict: FAITHFUL**

The source extracts every packed field and delegates overflow and underflow to UTC time.Date. The model applies Gregorian normalization with month carry and proportional day, minute and second offsets. Literal controls check zero month/day, an overflow month and oversized clock fields, and distinguish February 29 in 2000 from the corresponding date in 2100.

The author scope covers nonnil header receivers, Gregorian UTC decoding, and packing from projected civil fields. Location object internals and monotonic clock implementation details are outside the Lean model. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 3,150 fixed and 60,000 sampled rows at seeds 17, 2718 and 65537 on darwin/arm64. Six independent literal controls passed. Packing the local fields in SetModTime disagrees on all 6, and rounding odd seconds upward disagrees on 4.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/archive/zip/dos-time/review
```

Evidence SHA-256: `304c45fa4f8092b2eb7fc3859e974bb4311a89d8d7bffa136a4777c092097e01`
