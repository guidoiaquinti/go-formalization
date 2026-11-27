# Independent review of `archive/zip.FileHeader.SetModTime`

**Verdict: FAITHFUL**

The source converts the input instant to UTC, preserves seconds and nanoseconds in Modified, and packs UTC civil fields into the legacy words. The model accepts the projected UTC civil fields and its theorem preserves the stored instant. Native controls also check UTC location and unchanged populated header fields.

The author scope covers nonnil header receivers, Gregorian UTC decoding, and packing from projected civil fields. Location object internals and monotonic clock implementation details are outside the Lean model. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 3,150 fixed and 60,000 sampled rows at seeds 17, 2718 and 65537 on darwin/arm64. Six independent literal controls passed. Packing the local fields in SetModTime disagrees on all 6, and rounding odd seconds upward disagrees on 4.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/archive/zip/dos-time/review
```

Evidence SHA-256: `e7c7fe2ee697841140df4bfbac978e0a8f371e3d00d9fa020f82131dc6b38d54`
