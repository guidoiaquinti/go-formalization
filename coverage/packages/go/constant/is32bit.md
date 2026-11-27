# Independent review of `go/constant.is32bit`

**Verdict: FAITHFUL**

The inclusive signed interval is from -2^31 through 2^31-1. Controls cover both endpoints and their immediate outside neighbors. An unsigned-fit or symmetric-magnitude rule disagrees at these boundaries.

The author scope covers signed int64 words and canonical arbitrary-precision integers. Nil internal pointers are exercised for rendering only. Storage and allocation internals are outside scope. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 336 fixed cases and 60,000 seeded cases at seeds 17, 2718 and 65537 on darwin/arm64. Eleven independent controls passed, including signed width transitions, multiword rendering and the final accepted and first rejected magnitude bit lengths.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/go/constant/integer-render/review
```

Evidence SHA-256: `7fcd3d322160af95223dc6346ab315c51602cfbc0e2c2b692df7d1d687b707fc`
