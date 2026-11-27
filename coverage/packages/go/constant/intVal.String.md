# Independent review of `go/constant.intVal.String`

**Verdict: FAITHFUL**

The source delegates to big.Int.String, which renders canonical decimal digits and returns the literal <nil> for a nil pointer. Native controls include positive and negative multiword magnitudes, zero, and a nil representation.

The author scope covers signed int64 words and canonical arbitrary-precision integers. Nil internal pointers are exercised for rendering only. Storage and allocation internals are outside scope. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 336 fixed cases and 60,000 seeded cases at seeds 17, 2718 and 65537 on darwin/arm64. Eleven independent controls passed, including signed width transitions, multiword rendering and the final accepted and first rejected magnitude bit lengths.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/go/constant/integer-render/review
```

Evidence SHA-256: `7d4a6cbd9d91b17ce7629f208be177ed9393a492a349d53240b47e08a6b0c95f`
