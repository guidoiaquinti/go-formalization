# Independent review of `go/constant.smallInt`

**Verdict: FAITHFUL**

The source uses magnitude bit length strictly below 4,096. Native controls distinguish both signs of 2^4095-1 from both signs of 2^4095. The theorem checks that zero counts as small.

The author scope covers signed int64 words and canonical arbitrary-precision integers. Nil internal pointers are exercised for rendering only. Storage and allocation internals are outside scope. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 336 fixed cases and 60,000 seeded cases at seeds 17, 2718 and 65537 on darwin/arm64. Eleven independent controls passed, including signed width transitions, multiword rendering and the final accepted and first rejected magnitude bit lengths.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/go/constant/integer-render/review
```

Evidence SHA-256: `64750e35a47a837652fcdc21e2dd95d4790bb9a2ea65f6e09ae971cd8124f9d7`
