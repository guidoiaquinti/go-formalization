# Independent review of `go/constant.intVal.Kind`

**Verdict: FAITHFUL**

The method returns Int independent of the magnitude. Native controls construct the internal representation directly, including a shifted negative minimum and zero.

The scope covers initialized numeric values and finite acyclic complex structures. Big-number magnitude and numeric construction details are abstracted. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 147 fixed and 60,000 sampled cases at seeds 17, 2718 and 65537 on darwin/arm64. Six independent literal controls and a native signed-zero and panic-order fixture passed. Replacing signed OR with a maximum gives +1 for a negative and a positive component and fails the mixed-sign controls.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/go/constant/numeric-sign/review
```

Evidence SHA-256: `b598334ed9ff3750f62161e7b59f668c89b6ec6e4f4575a811961f92f64b7a00`
