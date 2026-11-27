# Independent review of `math/big.Int.Bit`

**Verdict: FAITHFUL**

A nonnegative index observes the requested bit under infinite two’s-complement sign extension. Negative indices panic. The theorem checks the negative-index class, and a literal control checks bit 130 of -1.

The scope covers normalized nonnil signed magnitudes and feasible shifts and bit indices. The Lean model abstracts raw limbs and allocation behavior. Native verification checks receiver identity and separate operand preservation for successful writes, including aliased receiver profiles. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 3,584 fixed cases and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Ten independent literal controls passed across three native receiver profiles, alongside direct native checks of the applicable algebraic laws.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/math/big/int-bitwise/review
```

Evidence SHA-256: `0d498c8884a37b64e21c2e9f2bdef7ffd67a6e08f2ae869987a327b1e6645d81`
