# Independent review of `math/bits.TrailingZeros8`

**Verdict: FAITHFUL**

Evidence SHA-256: `de6e42bd79b28609b36bada3f269f7257ec3714ed01069d9faa491483bd7bfdb`

## Current independent review, 2026-11-25

A separate Codex reviewer that did not author the evidence checked Go 1.27.0 and Lean 4.35.0. The reviewer ran `make proof` and `make diff TARGET=math/bits`, including the registered controls and both default lowering and the build with intrinsics disabled in the driver and `math/bits`. Hosted CI execution remains pending.

I regenerated `bits_tables.go` with the pinned Go toolchain and obtained a byte-for-byte match. The Go function reads `ntz8tab[x]`. The Lean model tests the low bit and shifts an unsigned byte until the first set bit or eight iterations, matching the generator. The kernel proved the zero result separately and checked the nonzero remainder equation, which requires the reported position to be set and every lower position to be clear.

On macOS/arm64, Go and the compiled Lean model agreed with the bit-position oracle on all 256 inputs in both modes. The independent control checked zero and evaluated the theorem's remainder equation on the 255 nonzero inputs. Go passed both properties. The capped scan failed at zero, and the leading-zero mutant failed the remainder equation on 170 nonzero inputs, starting at `0x01`.

## Earlier review evidence

The earlier verdict was FAITHFUL for the value returned on every `uint8` input.

I inspected the local go1.27.0 implementation, the 256 entries of `ntz8tab`, and its generator. Go returns the byte at index `x` as an `int`. The Lean model computes that entry by testing the low bit and shifting an unsigned `BitVec 8` right, as the generator does. The eight-step bound gives 8 at zero. Every result fits exactly in both Go's `int` and Lean's `Nat`. [Pinned Go function](https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L66-L69) · [Table generator](https://github.com/golang/go/blob/go1.27.0/src/math/bits/make_tables.go#L61-L67)

`trailingZeros8_zero` fixes the result at zero. For each nonzero byte, `trailingZeros8_spec` states that if the result is `n`, then `x % 2^(n+1) = 2^n` in natural-number arithmetic. This fixes the result uniquely because it says bit `n` is set and every lower bit is clear. Both proofs use kernel-checked `decide`. Lean 4.35.0 reported only `propext` and `Quot.sound` for each theorem.

The migrated Go check counts consecutive clear bits starting at bit zero, up to the byte width. On darwin/arm64, pinned Go and the compiled Lean model returned the expected count for all 256 bytes. The independent control exercised the zero theorem at 0 and the remainder equation at each byte from 1 through 255. It reproduced the failure counts below.

| Implementation under test | Zero failures (1 input) | Remainder equation failures (255 inputs) |
| --- | ---: | ---: |
| go1.27.0 `bits.TrailingZeros8` | 0 | 0 |
| Scan capped at seven shifts | 1 | 0 |
| Count leading zeros instead | 0 | 170 |

The seven-shift mutant returns 7 for zero, where the zero theorem requires 8. The leading-zero mutant returns 7 for `0x01`, giving `1 % 256 = 1` where the nonzero theorem requires 128.

The 2026-11-25 review passed the Lean build and proof audit, then ran the migrated comparison and property controls with the commands below. [Go check](../../../../difftest/math/bits/trailingzeros8/check.go) · [Reviewer program](controls/TrailingZeros8/main.go)

```sh
make proof
export GOTOOLCHAIN="go$(cat .go-version)" GOWORK=off
go run difftest/math/bits/trailingzeros8/check.go
go run coverage/packages/math/bits/controls/TrailingZeros8/main.go
```

With the renamed Go module, I checked that restoring the prior module declaration and driver import reproduced the original evidence digest. Go and Lean agreed on all 256 bytes. The capped scan failed at zero and the reversed-direction control failed on 170 nonzero bytes.

Earlier evidence digest SHA-256 was `c5597b7f0da94ae35c4e187c75f5e4f9d3b5af4932c0c0a24136c769e6c42d16`.
