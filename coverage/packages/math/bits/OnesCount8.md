# Independent review of `math/bits.OnesCount8`

**Verdict: FAITHFUL**

Evidence SHA-256: `c78b37e6735f658401c8d1fd3771f3b9a25c115346e5196466e45cdfb9ff2bff`

## Current independent review, 2026-11-25

A separate Codex reviewer that did not author the evidence checked Go 1.27.0 and Lean 4.35.0. The reviewer ran `make proof` and `make diff TARGET=math/bits`, including the registered controls and both default lowering and the build with intrinsics disabled in the driver and `math/bits`. Hosted CI execution remains pending.

I regenerated `bits_tables.go` with the pinned Go toolchain and obtained a byte-for-byte match. The Go function reads `pop8tab[x]`. The generator clears the lowest set bit on each iteration, and the Lean model uses the same operation with enough fuel for all eight positions. The kernel-checked `onesCount8_spec` counts exactly the set positions from 0 through 7.

On macOS/arm64, the driver enumerated all 256 bytes in both modes and summed bits independently of the model's loop. Go and the compiled Lean model returned that sum for every input. The independent control computed the sum by division and remainder, with 0 Go failures and one failure with a seven-iteration limit. At `0xff` the mutant returned 7 instead of 8.

## Earlier review evidence

The earlier verdict was FAITHFUL for the value returned on every `uint8` input.

Go 1.27.0 returns `int(pop8tab[x])`. I ran the pinned table generator in a temporary directory and verified that its output matched the shipped `bits_tables.go` byte for byte. For each nonzero byte the generator clears one set bit with `x &= x - 1` and increments the count. The Lean model performs the same operation on `BitVec 8`. Eight iterations suffice even for `0xff`, and zero needs none. Every possible result, 0 through 8, fits exactly in both Go's `int` and Lean's `Nat`. [Pinned Go function](https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L124-L127) · [Table generator](https://github.com/golang/go/blob/go1.27.0/src/math/bits/make_tables.go#L69-L75)

`GoLean.Math.Bits.onesCount8_spec` equates the result to the number of positions from 0 through 7 whose bit is set. This uniquely specifies every return value. The proof uses kernel-checked `decide`. Lean 4.35.0 reported only `propext` and `Quot.sound` when I inspected the theorem's axioms.

The migrated Go check sums the eight individual bits of each input. On darwin/arm64, pinned Go and the compiled Lean model returned this sum for all 256 bytes. The independent control evaluated the theorem by taking remainders modulo two and dividing by two at each position. It reproduced the failure counts below.

| Implementation under test | Inputs checked | Theorem failures |
| --- | ---: | ---: |
| go1.27.0 `bits.OnesCount8` | 256 | 0 |
| Clear at most seven set bits | 256 | 1 |

The seven-iteration mutant returned 7 at `0xff`, where all eight positions are set. That was its only failure.

The 2026-11-25 review passed the Lean build and proof audit, then ran the migrated comparison and property controls with the commands below. [Go check](../../../../difftest/math/bits/onescount8/check.go) · [Reviewer program](controls/OnesCount8/main.go)

```sh
make proof
export GOTOOLCHAIN="go$(cat .go-version)" GOWORK=off
go run difftest/math/bits/onescount8/check.go
go run coverage/packages/math/bits/controls/OnesCount8/main.go
```

With the renamed Go module, I checked that restoring the prior module declaration and driver import reproduced the original evidence digest. Go and Lean agreed on all 256 bytes. The seven-iteration control failed only at 0xff.

Earlier evidence digest SHA-256 was `b6edbfe67e265c7af80852cd787fabfca6db29d4471a4fea6a2d18d8002f1b4c`.
