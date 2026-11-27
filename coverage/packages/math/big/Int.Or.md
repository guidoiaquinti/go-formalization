# Independent review of `math/big.Int.Or`

**Verdict: FAITHFUL**

Both-negative Go operands complement the conjunction of magnitudes minus one, and the mixed sign branch complements the appropriate difference mask. The theorem proves idempotence. Fixed controls include negative sign extension and shared inputs.

The scope covers normalized nonnil signed magnitudes and feasible shifts and bit indices. The Lean model abstracts raw limbs and allocation behavior. Native verification checks receiver identity and separate operand preservation for successful writes, including aliased receiver profiles. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 3,584 fixed cases and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Ten independent literal controls passed across three native receiver profiles, alongside direct native checks of the applicable algebraic laws.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/math/big/int-bitwise/review
```

Evidence SHA-256: `03f1f90feccaa6d4e2b2457896a5dab0413c045cab943b138188ecce4d2d6ef1`
