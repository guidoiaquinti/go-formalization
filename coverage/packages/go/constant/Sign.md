# Independent review of `go/constant.Sign`

**Verdict: FAITHFUL**

The pinned numeric paths return a signed sign class. Complex values recursively combine both signs by signed OR, giving -1 whenever either component is negative. Unknown gives one, and nonnumeric values panic using their formatted representation. Native controls verify negative floating zero, first failing component order, and exact nonnumeric diagnostics. Lean abstracts the rendering as an input string.

The scope covers initialized numeric values and finite acyclic complex structures. Big-number magnitude and numeric construction details are abstracted. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 147 fixed and 60,000 sampled cases at seeds 17, 2718 and 65537 on darwin/arm64. Six independent literal controls and a native signed-zero and panic-order fixture passed. Replacing signed OR with a maximum gives +1 for a negative and a positive component and fails the mixed-sign controls.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/go/constant/numeric-sign/review
```

Evidence SHA-256: `038c75956bee0728dcc919d1176de9fb6434aa6a23946d69f2adcf866939760b`
