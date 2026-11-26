# Independent review of `math/bits.Rem`

**Verdict: FAITHFUL**

Evidence SHA-256: `12e08378edb3a46bbab062204d165ddd650360524d74654cfd064b179022cb64`

## Current independent review, 2026-11-25

A separate Codex reviewer that did not author the evidence checked Go 1.27.0 and Lean 4.35.0. The reviewer ran `make proof` and `make diff TARGET=math/bits`, including the registered controls and both default lowering and the build with intrinsics disabled in the driver and `math/bits`. Hosted CI execution remains pending.

I checked both width branches against the pinned Go source. The 32-bit delegate widens the dividend before taking the remainder. The 64-bit delegate reduces the high word modulo the divisor before entering `Div64`, so a positive divisor always satisfies its quotient-overflow precondition. For a zero divisor, both delegate paths panic while evaluating the integer remainder. The auxiliary source model follows the normalization and both quotient-digit corrections in `Div64`, with a two-correction bound and an explicit exhaustion result.

The kernel-checked `remUint_spec` uniquely determines the result for either width. The panic case requires a zero divisor, and a returned remainder has to be below the divisor and reconstruct the dividend with some natural quotient. The theorem proves the arithmetic model. Agreement with the auxiliary source model, including the quotient-digit correction bound, is checked on the sampled inputs.

On macOS/arm64, the native `bits.Rem` wrapper on 64-bit words and `bits.Rem32` on 32-bit words agreed with both Lean models and the full-dividend oracle on 61,728 triples per width and mode. Each set includes 1,728 boundary triples and 20,000 cases for each seed `17`, `2718` and `65537`. The reviewer control found 0 Go failures per width and reproduced 792 failures for direct division, 756 for the low-word-only mutant, and 144 for the zero-divisor mutant. The direct-division mutant panicked on `(1, 0, 1)`, where Go returned zero. The Linux/386 wrapper's hosted execution remains pending. Earlier supplemental 32-bit work is recorded below.

## Earlier review evidence

The earlier verdict was FAITHFUL for the returned remainder and divide-by-zero panic on both Go `uint` widths.

I checked the local Go 1.27.0 source. `Rem` delegates by word width. `Rem32` widens both words to form the dividend, which the model represents exactly as `hi * 2^32 + lo`. `Rem64` first reduces `hi` modulo `y`, then calls `Div64`. For positive `y`, that reduction puts the high word strictly below `y`, even when the original quotient would overflow. Both paths evaluate a remainder with divisor zero when `y == 0`. The model captures the observable outcome and uses natural-number remainder for the result of the delegated division. This review credits `Rem` alone. [Dispatch and remainder implementations](https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L570-L599) · [Div64 preconditions and implementation](https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L514-L568)

`remUint_spec` covers every triple of words. A panic requires a zero divisor. A returned `r` must be below `y`, which forces a panic at zero. For positive `y`, the theorem requires a natural `q` satisfying `hi * 2^width + lo = y * q + r`. Together these clauses uniquely determine the outcome. The proof unfolds the model and shows that reducing the high word preserves the original dividend modulo `y`.

I independently exercised those clauses on darwin/arm64 with pinned Go. I used all 1,728 triples drawn from zero, one, two, three, and the boundaries around the half-word radix, the sign bit and the largest word. For each returned value I checked the bound and used `math/big` to check that subtracting it from the full dividend gives a nonnegative multiple of `y`. At zero I recovered the panic and checked its `runtime.Error` type and divide-by-zero message. I called native `bits.Rem` for 64 bits and the delegated `bits.Rem32` for 32 bits. A native 32-bit binary was outside this run.

| Implementation under test | 32-bit failures, positive / zero divisor | 64-bit failures, positive / zero divisor |
| --- | ---: | ---: |
| Pinned Go (1,584 positive and 144 zero cases per width) | 0 / 0 | 0 / 0 |
| Call word-width `Div` with the original high word | 792 / 0 | 792 / 0 |
| Compute only `lo % y` | 756 / 0 | 756 / 0 |
| Return zero for a zero divisor | 0 / 144 | 0 / 144 |

At `(hi, lo, y) = (1, 0, 1)` the direct-division mutant panicked with integer overflow. Go returned zero. At `(1, 0, 3)` the low-word mutant returned zero and Go returned one at both widths. The zero-return mutant failed even at `(0, 0, 0)`.

Reproduce the property and mutant counts from the repository root. [Reviewer program](controls/Rem/main.go)

```sh
GOTOOLCHAIN="go$(cat .go-version)" GOWORK=off go run coverage/packages/math/bits/controls/Rem/main.go
```

## Review of the Go differential check, 2026-11-25

I checked the width dispatch against the pinned Go source and type-checked `remUint_spec` with Lean 4.35.0. The Go driver sends each input triple to the compiled Lean model and calls the indicated `math/bits` function. For positive divisors, it constructs the two-word dividend with `math/big` and compares its remainder with both results. For a zero divisor, it checks the panic type and message. The shared driver helper enforces the Go version pin and the exact Lean row count. [Go driver](../../../../difftest/math/bits/rem/check.go) · [Driver helper](../../../../difftest/internal/check/check.go)

On darwin/arm64, the following command passed 61,728 cases per width using the 1,728 boundary triples and 20,000 cases for each seed (17, 2718 and 65537). It exercised `bits.Rem` on 64-bit inputs and `bits.Rem32` on 32-bit inputs. The reviewer program independently reproduced every failure count above.

```sh
GOTOOLCHAIN="go$(cat .go-version)" GOWORK=off go run ./difftest/math/bits/rem/check.go
```

On this host, `--native-width 32` exited with an error identifying the 64-bit binary. The Linux CI step sets `GOARCH=386` and uses that option to require the 32-bit `bits.Rem` wrapper. [CI workflow](../../../../.github/workflows/verify.yml)

### Supplemental 32-bit verification

A local Linux/386 build of the reviewer control exercised the actual `bits.Rem` wrapper on the 1,728 boundary triples and reproduced the table above. A second run exercised the wrapper on all 61,728 registered inputs. The local emulator failed while spawning the Lean subprocess, so this run used a temporary Go build overlay for the shared I/O helper to forward each batch to the compiled Lean executable on macOS. The registered case generator, Go calls, and outcome comparisons ran in the 32-bit process. All four batches agreed. CI runs the same cases on Linux/amd64 using the repository helper directly.

With the renamed Go module, I checked that restoring the prior module declaration and driver import reproduced the original evidence digest. On darwin/arm64, Go and Lean agreed on all 61,728 cases per width. The run exercised the 64-bit wrapper and 32-bit delegate, and the reviewer program reproduced every mutant count above.

Earlier evidence digest SHA-256 was `3847bdcf59364c0e0428599582a422049cb5f75c459cc80fd218b03da422249b`.
