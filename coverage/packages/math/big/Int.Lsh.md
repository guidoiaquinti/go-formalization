# Independent review of `math/big.Int.Lsh`

**Verdict: FAITHFUL**

The source shifts the magnitude and restores its sign for nonzero results. The theorem checks zero shift. Controls include a 64-bit left shift of a multiword positive magnitude and negative input shifts.

The scope covers normalized nonnil signed magnitudes and feasible shifts and bit indices. The Lean model abstracts raw limbs and allocation behavior. Native verification checks receiver identity and separate operand preservation for successful writes, including aliased receiver profiles. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 3,584 fixed cases and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Ten independent literal controls passed across three native receiver profiles, alongside direct native checks of the applicable algebraic laws.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/math/big/int-bitwise/review
```

Evidence SHA-256: `f0f2f8bfb3c65302768035e1f9eb1b347f341e653e165997dc335127215b5c9d`
