# Independent review of `math/big.Int.Rsh`

**Verdict: FAITHFUL**

The negative branch subtracts one from the magnitude before shifting and complements the result, so division rounds down for negative inputs. The theorem checks zero shift. Truncating signed division disagrees on three literal controls.

The scope covers normalized nonnil signed magnitudes and feasible shifts and bit indices. The Lean model abstracts raw limbs and allocation behavior. Native verification checks receiver identity and separate operand preservation for successful writes, including aliased receiver profiles. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 3,584 fixed cases and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Ten independent literal controls passed across three native receiver profiles, alongside direct native checks of the applicable algebraic laws.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/math/big/int-bitwise/review
```

Evidence SHA-256: `d136f30253b02ef0a98da642beda199625dfda0c999a7e7bcbd89bd710fe8579`
