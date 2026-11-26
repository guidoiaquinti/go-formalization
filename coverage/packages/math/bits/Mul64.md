# Independent review of `math/bits.Mul64`

**Verdict: FAITHFUL**

Evidence SHA-256: `52b27a5110e2fd8e801ea59f3aff7e3cc1d5aabeaf6501a73f8a13f3b3775f08`

## Current independent review, 2026-11-25

A separate Codex reviewer that did not author the evidence checked Go 1.27.0 and Lean 4.35.0. The reviewer ran `make proof` and `make diff TARGET=math/bits`, including the registered controls and both default lowering and the build with intrinsics disabled in the driver and `math/bits`. Hosted CI execution remains pending.

I compared the pinned Go declaration with `mul64Source` statement by statement. Go types the mask as `uint64` for the AND expressions, and the shifts bind before the surrounding additions. The Lean definition uses the same 64-bit mask, logical shifts and wrapping products and sums. Its second binding for `w1` consumes the first binding, matching Go's compound assignment. The source guard passed against the pinned declaration. Its checked transcription and the manual correspondence connect the Go source to Lean. [Source correspondence](Mul64-source.md)

The kernel checked `mul64Source_eq_mul64` and `mul64Source_spec` for every pair of words. The first proof reduces masks, shifts and wrapping operations to the corresponding natural-number arithmetic. The second requires that recombining the two bounded output words gives the full product, which determines their ordering and values uniquely.

On macOS/arm64, Go and the compiled Lean source model agreed with the shift-and-add oracle on 60,144 pairs in each mode. That includes 144 boundary pairs and 20,000 cases for each seed `17`, `2718` and `65537`. The independent reconstruction control found 0 failures on Go and 97 of 256 when omitting the last carry. At `x = 0x80000000` and `y = 0x8000000000000000`, the mutant returned zero for a product of `2^94`.

In the compiled driver I observed Go's `Mul64uhilo` substitution at the direct call, with `UMULH` on macOS/arm64 and in the Linux/arm64 cross-build and `MULQ` in the Linux/amd64 cross-build. The macOS fallback binary calls `math/bits.Mul64` directly. Host execution samples the compiler output. The universal theorems cover the Lean models and return values. Instruction semantics, compiler correctness and input-independent timing are outside their scope.

## Earlier review evidence

The earlier verdict was FAITHFUL for the high and low words returned for every pair of `uint64` inputs. Go's input-independent execution time is outside the scope of the value model.

I checked the model against the Go 1.27.0 source. For each input, the Lean remainder and quotient by `2^32` implement Go's low-word mask and right shift. The model follows the two cross products and their carries into the high word. The low word uses wrapping multiplication. The proof bounds both cross-product sums below `2^64`, so their unsigned reductions preserve the intermediate values used in the proof. [Pinned implementation](https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L465-L484) · [Model and proof](../../../../lean/GoLean/Math/Bits/Arithmetic.lean)

The theorem requires `2^64 * hi + lo = x * y` in natural-number arithmetic. Each output word is a `BitVec 64`, so the equation fixes both words uniquely. Lean 4.35.0 checked the proof and reported the axiom dependencies `propext`, `Classical.choice` and `Quot.sound`.

I ran the compiled Lean check on darwin/arm64. Pinned Go, Lean and a shift-and-add reference agreed on all 144 boundary pairs and 20,000 random pairs for each seed 17, 2718 and 65537. [Registered check](../../../../difftest/math/bits/mul64/check.go) · [Lean executable](../../../../lean/GoDiff/Math/Bits/Mul64.lean)

I separately exercised the theorem's reconstruction equation on 256 pairs. The reviewer program forms the reference product from sixteen 16-bit partial products. Each fits in 32 bits, and the program adds the shifted terms with arbitrary precision. The input grid includes zero, one, the top bit and the maximum value in each 32-bit limb. [Reviewer control](controls/Mul64/main.go)

| Implementation | Reconstruction failures |
| --- | ---: |
| Go 1.27.0 `bits.Mul64` | 0 of 256 |
| Omit the final carry into the high word | 97 of 256 |

At `x = 0x80000000` and `y = 0x8000000000000000`, the mutant returns two zero words. The true product is `2^94`, so the high word must be `0x40000000`.

Run from the repository root after `make build`.

```sh
GOTOOLCHAIN="go$(cat .go-version)" GOWORK=off go run ./difftest/math/bits/mul64/check.go
GOTOOLCHAIN="go$(cat .go-version)" GOWORK=off go run coverage/packages/math/bits/controls/Mul64/main.go
```

With the renamed Go module, I checked that restoring the prior module declaration and driver import reproduced the original evidence digest. Go and Lean agreed on all 60,144 input pairs. The independent control reproduced the 97 failures for the missing carry.

Earlier evidence digest SHA-256 was `6b2cef708cd16e644699c64fef7a517b015242e1945830bc4bd7ed67bc5a5701`.
