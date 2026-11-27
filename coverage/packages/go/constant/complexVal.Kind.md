# Independent review of `go/constant.complexVal.Kind`

**Verdict: FAITHFUL**

The method returns Complex independent of component signs or magnitudes. Native fixtures construct the internal complex representation directly, including both components zero and mixed signs.

The scope covers initialized numeric values and finite acyclic complex structures. Big-number magnitude and numeric construction details are abstracted. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 147 fixed and 60,000 sampled cases at seeds 17, 2718 and 65537 on darwin/arm64. Six independent literal controls and a native signed-zero and panic-order fixture passed. Replacing signed OR with a maximum gives +1 for a negative and a positive component and fails the mixed-sign controls.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/go/constant/numeric-sign/review
```

Evidence SHA-256: `8087168611c3b89069371147d872e164c95f79cc3d94a52bf23606199c6b002b`
