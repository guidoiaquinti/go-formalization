# Independent review of `math/big.Int.SetBit`

**Verdict: FAITHFUL**

The source updates the requested bit in the infinite two’s-complement representation. Negative indices reject before invalid bit values. The theorems distinguish both error classes, and controls exercise high-bit clearing on negative numbers and exact diagnostics. Backing mutation after an invalid-value panic is outside scope.

The scope covers normalized nonnil signed magnitudes and feasible shifts and bit indices. The Lean model abstracts raw limbs and allocation behavior. Native verification checks receiver identity and separate operand preservation for successful writes, including aliased receiver profiles. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 3,584 fixed cases and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Ten independent literal controls passed across three native receiver profiles, alongside direct native checks of the applicable algebraic laws.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/math/big/int-bitwise/review
```

Evidence SHA-256: `ebe01f7b7f32703948d20679fe8d7fea898932a70991802ede863b27b8fe6736`
