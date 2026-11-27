# Independent review of `math/big.Int.TrailingZeroBits`

**Verdict: FAITHFUL**

The source counts low zero bits of the unsigned magnitude and yields zero for zero. The theorems check zero and negation invariance. Controls include exact counts at 64 and 129 bits.

The scope covers normalized nonnil signed magnitudes and feasible shifts and bit indices. The Lean model abstracts raw limbs and allocation behavior. Native verification checks receiver identity and separate operand preservation for successful writes, including aliased receiver profiles. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 3,584 fixed cases and 60,000 samples at seeds 17, 2718 and 65537 on darwin/arm64. Ten independent literal controls passed across three native receiver profiles, alongside direct native checks of the applicable algebraic laws.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/math/big/int-bitwise/review
```

Evidence SHA-256: `d75db13c1e941833beb375f4c0c9334ce0b2a025b866f6e7ecd4d8ba929aec7c`
