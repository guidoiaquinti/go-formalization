# Independent review of `math/big.Int.Xor`

**Verdict: FAITHFUL**

Both-negative operands cancel the infinite sign bits and give a nonnegative XOR of magnitudes minus one. Mixed signs retain a negative result. The theorem covers self-XOR and literal controls distinguish both sign paths.

The scope covers normalized nonnil signed magnitudes and feasible shifts and bit indices. The Lean model abstracts raw limbs and allocation behavior. Native verification checks receiver identity and separate operand preservation for successful writes, including aliased receiver profiles. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 3,584 fixed cases and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Ten independent literal controls passed across three native receiver profiles, alongside direct native checks of the applicable algebraic laws.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/math/big/int-bitwise/review
```

Evidence SHA-256: `0dd88d29add8893fd6e835a60780dd38d1a00811447f5ae07febde05adf7bde4`
