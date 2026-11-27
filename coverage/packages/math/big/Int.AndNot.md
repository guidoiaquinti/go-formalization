# Independent review of `math/big.Int.AndNot`

**Verdict: FAITHFUL**

For each sign combination, the source implements conjunction with the complement of the second operand. The theorem proves self-and-not gives zero. Mixed-sign controls distinguish this operation from ordinary And and from reversing the operands.

The scope covers normalized nonnil signed magnitudes and feasible shifts and bit indices. The Lean model abstracts raw limbs and allocation behavior. Native verification checks receiver identity and separate operand preservation for successful writes, including aliased receiver profiles. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 3,584 fixed cases and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Ten independent literal controls passed across three native receiver profiles, alongside direct native checks of the applicable algebraic laws.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/math/big/int-bitwise/review
```

Evidence SHA-256: `b02a7798927952421954d134203312ffb7d364cf45c63abae84add9168c06b01`
