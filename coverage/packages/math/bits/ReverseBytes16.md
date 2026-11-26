# Independent review of `math/bits.ReverseBytes16`

**Verdict: FAITHFUL**

Evidence SHA-256: `f91311b78234bc2614a361ddad4c786cd1518ed5db49af4c99ab10923c3cf3fe`

## Current independent review, 2026-11-25

A separate Codex reviewer that did not author the evidence checked Go 1.27.0 and Lean 4.35.0. The reviewer ran `make proof` and `make diff TARGET=math/bits`, including the registered controls and both default lowering and the build with intrinsics disabled in the driver and `math/bits`. Hosted CI execution remains pending.

I compared the pinned Go expression with the Lean shifts and bitwise OR at the same 16-bit width. The kernel checked that the output consists of the original low byte followed by the original high byte, and that applying the model twice recovers the input.

On macOS/arm64, all 65,536 inputs in each mode matched the differential oracle, which writes the input little-endian and reads it big-endian. The control independently checked the byte equation and executed both Go calls for the round trip, with 0 failures for either property. Identity failed the byte equation on 65,280 inputs but passed all round trips. The right-shift-only mutant failed the byte equation on 65,280 inputs and the round trip on 65,535 inputs. Timing is outside the value claim.

## Earlier review evidence

The earlier verdict was FAITHFUL for the value returned on every `uint16` input. The source also documents input-independent execution time, which is outside the scope of the value model.

I checked the implementation in the local go1.27.0 toolchain. The Lean model translates each operation in `x>>8 | x<<8` directly. Its `BitVec 16` type gives the right shift zero fill and discards overflow from the left shift. [Pinned Go source, lines 273–278](https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L273-L278)

`reverseBytes16_spec` uniquely determines the output by specifying both bytes. Its proof checks each bit position. The involution proof uses that byte specification twice. Lean 4.35.0 reported only `propext` and `Quot.sound` as dependencies for either theorem. In temporary copies the audit rejected a fresh axiom, a proof using `sorry`, and a proof using `native_decide`.

The migrated Go check writes each input in little-endian order and reads those bytes in big-endian order. On darwin/arm64, pinned Go and the compiled Lean model returned this value for every one of the 65,536 inputs. The independent control used unsigned arithmetic to compute the expected byte swap and exercised the round trip separately. It reproduced the failure counts below.

| Implementation under test | Byte specification failures | Involution failures |
| --- | ---: | ---: |
| go1.27.0 `bits.ReverseBytes16` | 0 | 0 |
| Identity mutant | 65,280 | 0 |
| Right-shift-only mutant (`x >> 8`) | 65,280 | 65,535 |

The identity control shows why involution alone would be insufficient. The byte specification detects every unequal-byte input on that mutant.

The 2026-11-25 review passed the Lean build and proof audit, then ran the migrated comparison and property controls with the commands below. [Go check](../../../../difftest/math/bits/reversebytes16/check.go) · [Reviewer program](controls/ReverseBytes16/main.go)

```sh
make proof
export GOTOOLCHAIN="go$(cat .go-version)" GOWORK=off
go run difftest/math/bits/reversebytes16/check.go
go run coverage/packages/math/bits/controls/ReverseBytes16/main.go
```

With the renamed Go module, I checked that restoring the prior module declaration and driver import reproduced the original evidence digest. Go and Lean agreed on all 65,536 words. The identity control again failed the byte-order specification 65,280 times and passed every round trip.

Earlier evidence digest SHA-256 was `7bfd1db79ade3bcdf0d7db7b6790c9368959d38c26d9cde0cf3ae7f2f435d71f`.
