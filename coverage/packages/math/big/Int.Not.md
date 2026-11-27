# Independent review of `math/big.Int.Not`

**Verdict: FAITHFUL**

For nonnegative Go inputs the result is negative with magnitude incremented by one. For negative inputs it is the magnitude minus one. The involution theorem and native round-trip controls check the two branches.

The scope covers normalized nonnil signed magnitudes and feasible shifts and bit indices. The Lean model abstracts raw limbs and allocation behavior. Native verification checks receiver identity and separate operand preservation for successful writes, including aliased receiver profiles. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 3,584 fixed cases and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Ten independent literal controls passed across three native receiver profiles, alongside direct native checks of the applicable algebraic laws.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/math/big/int-bitwise/review
```

Evidence SHA-256: `1a6af8d2ea5503f8c424f09c3bc69af2305158b51efb86a67501e31b67e2b689`
