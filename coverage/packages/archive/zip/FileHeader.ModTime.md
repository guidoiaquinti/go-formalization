# Independent review of `archive/zip.FileHeader.ModTime`

**Verdict: FAITHFUL**

The source uses only the two legacy fields and returns an integral UTC instant. The theorem checks the DOS epoch. The Go overlay initializes Modified to a different instant, and fixed controls include invalid legacy fields and normalization.

The author scope covers nonnil header receivers, Gregorian UTC decoding, and packing from projected civil fields. Location object internals and monotonic clock implementation details are outside the Lean model. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 3,150 fixed and 60,000 sampled rows at seeds 17, 2718 and 65537 on darwin/arm64. Six independent literal controls passed. Packing the local fields in SetModTime disagrees on all 6, and rounding odd seconds upward disagrees on 4.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/archive/zip/dos-time/review
```

Evidence SHA-256: `fd540ebd49bc7ff9aba89689b56afa3f4a323709456b6bb9ab6df48685158a6b`
