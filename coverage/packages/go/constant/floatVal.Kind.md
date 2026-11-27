# Independent review of `go/constant.floatVal.Kind`

**Verdict: FAITHFUL**

The method returns Float for initialized big.Float values. Native review also constructs a signed negative zero, which retains Float kind and has sign zero.

The scope covers initialized numeric values and finite acyclic complex structures. Big-number magnitude and numeric construction details are abstracted. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 147 fixed and 60,000 sampled cases at seeds 17, 2718 and 65537 on darwin/arm64. Six independent literal controls and a native signed-zero and panic-order fixture passed. Replacing signed OR with a maximum gives +1 for a negative and a positive component and fails the mixed-sign controls.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/go/constant/numeric-sign/review
```

Evidence SHA-256: `c44a1835b86cf3b53bea0b9a23d637208dd0dfe04d7b4dfa38ab39c95e564bba`
