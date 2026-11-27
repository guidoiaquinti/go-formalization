# Independent review of `go/constant.intVal.ExactString`

**Verdict: FAITHFUL**

The source delegates to intVal.String. The control checks exact equality for positive and negative multiword values and the pinned nil spelling.

The author scope covers signed int64 words and canonical arbitrary-precision integers. Nil internal pointers are exercised for rendering only. Storage and allocation internals are outside scope. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 336 fixed cases and 60,000 seeded cases at seeds 17, 2718 and 65537 on darwin/arm64. Eleven independent controls passed, including signed width transitions, multiword rendering and the final accepted and first rejected magnitude bit lengths.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/go/constant/integer-render/review
```

Evidence SHA-256: `dd11c4af459644df1831363d68d058feb7af6417222fd1b73c8365f6c8fea241`
