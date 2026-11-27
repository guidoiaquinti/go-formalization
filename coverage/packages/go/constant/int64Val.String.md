# Independent review of `go/constant.int64Val.String`

**Verdict: FAITHFUL**

The source delegates to signed base-ten conversion. The model and native controls cover the minimum and maximum signed words and preserve the minus sign. The theorem checks the full minimum decimal string.

The author scope covers signed int64 words and canonical arbitrary-precision integers. Nil internal pointers are exercised for rendering only. Storage and allocation internals are outside scope. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 336 fixed cases and 60,000 seeded cases at seeds 17, 2718 and 65537 on darwin/arm64. Eleven independent controls passed, including signed width transitions, multiword rendering and the final accepted and first rejected magnitude bit lengths.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/go/constant/integer-render/review
```

Evidence SHA-256: `f6e0b0b6ee88d0602905bd48f5abb727e8813bf9119b97e2602be1772ae3f75c`
