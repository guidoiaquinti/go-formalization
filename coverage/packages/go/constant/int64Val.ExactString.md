# Independent review of `go/constant.int64Val.ExactString`

**Verdict: FAITHFUL**

The source directly delegates to the signed value String method. The theorem fixes the equality and controls exercise both signed extrema and both width boundaries.

The author scope covers signed int64 words and canonical arbitrary-precision integers. Nil internal pointers are exercised for rendering only. Storage and allocation internals are outside scope. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 336 fixed cases and 60,000 seeded cases at seeds 17, 2718 and 65537 on darwin/arm64. Eleven independent controls passed, including signed width transitions, multiword rendering and the final accepted and first rejected magnitude bit lengths.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/go/constant/integer-render/review
```

Evidence SHA-256: `a2ba46838a5164d963b458242f764450b4ba540a0cfc59e4fd018aa2dac7cfef`
