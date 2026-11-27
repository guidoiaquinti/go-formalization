# Independent review of `math/big.Int.And`

**Verdict: FAITHFUL**

The Go sign branches implement infinite two’s-complement conjunction by operating on magnitudes or the complemented magnitude minus one. The theorem establishes idempotence. Literal controls check mixed signs across the 64-bit boundary and two negative operands.

The scope covers normalized nonnil signed magnitudes and feasible shifts and bit indices. The Lean model abstracts raw limbs and allocation behavior. Native verification checks receiver identity and separate operand preservation for successful writes, including aliased receiver profiles. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 3,584 fixed cases and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Ten independent literal controls passed across three native receiver profiles, alongside direct native checks of the applicable algebraic laws.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/math/big/int-bitwise/review
```

Evidence SHA-256: `477789323e3ac6691a5077949992caad6d1b47a254b2aa4aff4b20b9999aa649`
