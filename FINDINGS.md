## POTENTIAL — net/netip: Addr.Next changes the stored bits of the zero Addr

Go 1.27.0 on darwin/arm64. The documentation says that Next returns the
zero Addr when there is no following address. If that includes a zero
receiver, the expected result is Addr{} with all-zero As16 bytes. The
observed result has BitLen 0 and As16 bytes ending in 1, so it differs
from Addr{}. The pinned implementation increments the address words
before checking overflow.

Reproducer from the repository root after building the Lean executables:

```sh
GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/net/netip/addr-next
```

The boundary case with the zero receiver checks the observed behavior.
Confirm whether the documented no-successor rule is intended to apply
to the zero receiver before treating this as a Go implementation bug.

## POTENTIAL — internal/dag: TransitiveReduction can retain a shortcut in a closed DAG

Go 1.27.0 on darwin/arm64. The source comment says TransitiveReduction
removes edges that are transitively reachable and requires a transitively
closed graph. The overlay constructs four nodes in insertion order
`0, 1, 2, 3` and adds the full closure of the chain
`1 → 0 → 2 → 3`. The original claim was that the method returns exactly
the cover relation for every transitively closed DAG. The expected edges
are `0 → 2`, `1 → 0` and `2 → 3`. Go and the source-shaped Lean model also
retain `1 → 3`. Earlier deletions remove witnesses that later iterations
could use to remove that shortcut.

Reproducer from the repository root after registering the evidence:

```sh
make diff TARGET=internal/dag
GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/internal/dag/graph-transitive-reduction/check.go -cover-counterexample
```

The last command intentionally fails against the preserved exact-cover
claim and prints the input, Go result, Lean result and expected matrix.
The regular comparison checks the sequential scan and its deletion
invariants. Confirm whether a construction-order restriction is intended
in addition to transitive closure before attributing this to the
implementation. The repro establishes the behavior on a graph constructed
using newGraph, addNode and AddEdge. Whether the runtime's
lock-rank generator supplies such a graph remains unverified.

## POTENTIAL — math: Ldexp reverses overflow and underflow at int extrema

Go 1.27.0 on darwin/arm64. The documented scaling law is
`Ldexp(frac, exp) = frac × 2**exp`. It predicts positive infinity
for `Ldexp(2, math.MaxInt)` and positive zero for
`Ldexp(0.5, math.MinInt)` after rounding to binary64. Go returns
positive zero and positive infinity respectively. The pinned
implementation adds the normalized fraction's leading exponent to
the supplied `int` before checking floating-point bounds. The
integer addition wraps at both inputs. [Pinned Go implementation](https://github.com/golang/go/blob/go1.27.0/src/math/ldexp.go#L22-L42)

The reproducer intentionally exits unsuccessfully against the original
scaling claim and prints both expected and observed bit patterns.

```sh
GOWORK=off GOTOOLCHAIN=go1.27.0 go run ./difftest/math/ldexp/check.go -contract-counterexample
```

The regular comparison models the observed wraparound. Check whether
the public contract intends to exclude exponents at these native integer
extrema before treating the discrepancy as an implementation bug.

## CONFIRMED — slices: Grow nilness differs from its documentation

Go 1.27.0 on darwin/arm64. The documentation says Grow preserves the
nilness of its input. With `slices.Grow([]int(nil), 1)` the documented
nilness would be true. The observed result is nonnil with length zero
and capacity of at least one. The pinned source takes its allocation
branch and reslices the newly allocated result to length zero.
[Pinned Go source](https://github.com/golang/go/blob/go1.27.0/src/slices/slices.go#L417-L430)

```sh
GOTOOLCHAIN=go1.27.0 GOWORK=off CGO_ENABLED=0 go run coverage/packages/slices/controls/Grow/main.go
```

The control prints the observed nilness. The ordinary comparison checks
the pinned implementation behavior.

## CONFIRMED — math/big: Float.Sqrt overwrites the receiver rounding mode

Go 1.27.0 on darwin/arm64. The documentation specifies rounding with the
receiver's precision and mode. With input 2, a distinct 53-bit receiver
set to ToZero should keep ToZero and produce the binary64 word
`0x3ff6a09e667f3bcc`. It instead reports ToNearestEven and produces
`0x3ff6a09e667f3bcd`, whose exact square exceeds 2. In the finite-input
branch, `x.MantExp(z)` copies the input's rounding mode to the receiver.
The implementation restores precision only. [Sqrt implementation](https://github.com/golang/go/blob/go1.27.0/src/math/big/sqrt.go#L23-L60) · [MantExp implementation](https://github.com/golang/go/blob/go1.27.0/src/math/big/float.go#L268-L282)

Save this program as `sqrt-rounding.go` and run
`GOTOOLCHAIN=go1.27.0 GOWORK=off go run sqrt-rounding.go`.

```go
package main

import (
	"fmt"
	"math"
	"math/big"
)

func main() {
	x := new(big.Float).SetFloat64(2)
	z := new(big.Float).SetPrec(53).SetMode(big.ToZero)
	z.Sqrt(x)
	f, _ := z.Float64()
	r := new(big.Rat).SetFloat64(f)
	square := new(big.Rat).Mul(r, r)
	fmt.Printf("mode=%v bits=%#016x square-vs-two=%d\n",
		z.Mode(), math.Float64bits(f), square.Cmp(big.NewRat(2, 1)))
}
```

The registered `math.Sqrt` oracle uses ToNearestEven on both input and
receiver, so the mode copy has no effect on that comparison.

## POTENTIAL — math.Hypot returns infinity below the overflow rounding midpoint

On Go 1.27.0 darwin/arm64, `math.Hypot(x, x)` returns positive infinity
for `x := math.Float64frombits(0x7fe6a09e667f3bcc)`. All four operand
sign combinations give the same result. The exact mathematical norm
rounds to the largest finite binary64 value.

Writing the magnitude as `k * 2^971` with `k = 0x16a09e667f3bcc`,
the squared norm in units of `2^1940` is `8*k*k`. It lies strictly
between the squared midpoints `(2^54-3)^2` and `(2^54-1)^2`,
with respective gaps 50548499069883511 and 21509094968044417.
The first is the midpoint between the two largest finite values.
The second is the overflow rounding midpoint.

Reproducer from the repository root after building the Lean executable:

```sh
GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/math/hypot
```

The staged Lean model agrees with the pinned native result. Whether the
numerical contract requires the finite result remains open for triage.

## POTENTIAL — math.Exp2 returns zero above the binary64 underflow midpoint

On Go 1.27.0 darwin/arm64, `math.Exp2(-1074.5)` returns positive zero.
The exact value is the smallest positive subnormal divided by the
square root of two. Its nearest binary64 result is the smallest
positive subnormal. The pinned implementation cuts off inputs below
`-1074`, including `-1075 < x < -1074`, where the exact
result exceeds the midpoint between zero and the first subnormal.

Reproducer from the repository root after building the Lean executable:

```sh
GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/math/exp2
```

The source-shaped Lean model preserves the cutoff. The comparison
records agreement with Go and checks the high-precision reference in
the cutoff interval. Whether Go's accuracy contract requires the
subnormal result remains open for triage.

## POTENTIAL — math.Exp returns zero for an input that rounds to a subnormal

Go 1.27.0 on darwin/arm64 returns positive zero for `math.Exp(-745)`.
The 256-bit reference rounds the exact exponential to the smallest
positive binary64 subnormal. This is one encoded step from the pinned
Go result. The staged Lean model reproduces the native result. The
private generic helper returns the smallest subnormal on the same
platform, as checked by a direct overlay.

Reproducer from the repository root after building the Lean executable:

```sh
GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/math/exp
GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/math/exp/review
```

The native amd64 result for this input and the required rounding
accuracy remain open for triage.

## POTENTIAL — math.Atan2 selects positive π after lower-left quotient underflow

Go 1.27.0 on darwin/arm64 returns positive π for
`math.Atan2(-math.SmallestNonzeroFloat64, -2)`. Its quotient rounds
to positive zero. The comparison `q <= 0` then selects `q + Pi`,
giving `0x400921fb54442d18`. The mathematical angle rounds to
negative π, `0xc00921fb54442d18`. With `x` one representable step
toward zero from `-2`, the quotient rounds to the smallest positive
subnormal and Go returns negative π. At `-2` and its neighbor toward
negative infinity, Go returns positive π.
[Pinned Go implementation](https://github.com/golang/go/blob/go1.27.0/src/math/atan2.go#L37-L76)

Reproducer from the repository root after building the Lean executable:

```sh
GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/math/atan2
```

The driver confirms all three boundary words, retains the observed Go
behavior in the model and reports the independent mathematical reference.
Native amd64 execution and whether Go intends a wrapped angle when the
rounded quotient is zero remain open for triage.

## POTENTIAL — math.Sinh and math.Cosh overflow before halving their exponentials

Go 1.27.0 on darwin/arm64 returns positive infinity for
`math.Sinh(710)` and negative infinity for `math.Sinh(-710)`. Both signs give positive
infinity for `math.Cosh` at the same magnitude.
The 256-bit reference computes `(exp(|x|) - 1/exp(|x|))/2` before
rounding to binary64. It gives finite values with words
`0x7fe3e21a464507f9` and `0xffe3e21a464507f9`, respectively. The Cosh reference adds
the reciprocal before halving, giving `0x7fe3e21a464507f9`
for both signs.
Both pinned large-magnitude branches round Exp to binary64 before
multiplication by 0.5, so the intermediate overflows.
[Pinned Go implementation](https://github.com/golang/go/blob/go1.27.0/src/math/sinh.go#L33-L93)

Reproducer from the repository root after building the Lean executables:

```sh
GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/math/sinh
GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/math/cosh
```

Both comparisons preserve the pinned implementation and report the
independent mathematical outcomes for both signs. Confirm whether Go
intends finite results throughout the representable mathematical range
and reproduce the discrepancy on other architecture profiles before
attributing it to the public contract.

## CONFIRMED — math.Erfcinv returns positive infinity for some negative inputs

Go 1.27.0 on darwin/arm64 documents NaN for negative inputs.
`math.Erfcinv(-math.SmallestNonzeroFloat64)` and
`math.Erfcinv(-0x1p-53)` both return positive infinity,
word `0x7ff0000000000000`. The subtraction `1-x` rounds to
one before Erfinv is called. One representable step below
`-0x1p-53` produces staged word `0x3ff0000000000001`
and NaN. [Pinned source and documented domain](https://github.com/golang/go/blob/go1.27.0/src/math/erfinv.go#L119-L129)

The differential reproduces these words and preserves the pinned
behavior. Run it after building the Lean executable:

```sh
GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/math/erfcinv
```

## POTENTIAL — math.Erfcinv loses tiny positive inputs in its subtraction

The inverse of the complementary error function is finite for all
inputs strictly between zero and two. Go 1.27.0 on darwin/arm64
returns positive infinity for positive inputs through `0x1p-54`
because `1-x` rounds to one. The next binary64 input above that
cutoff returns a finite result. The same reproducer checks both
sides of the boundary, including the smallest positive subnormal.
Confirm the intended accuracy contract for tiny positive inputs
before treating the finite-precision loss as a separate Go bug.

## POTENTIAL — path/filepath: Rel can loop on a Windows UNC share

Go 1.27.0 Windows source, reproduced by executing the extracted pinned
algorithm on darwin/arm64. The call below is expected to return a lexical
result or an error. In the extracted implementation it did not return
within two seconds.

```go
filepath.Rel(`\\host\share`, `\\host\share\`)
```

The cleaned full paths differ by the trailing separator. After volume
removal the base UNC share is normalized to a lone separator, and the
target already has that body. Both cursors consume the separator.
Each subsequent iteration compares two empty tails and advances neither
cursor. The Lean comparison exposes the stationary state as `diverges`.
Native Windows reproduction remains pending.

## POTENTIAL — math.J1 returns a positive result for tiny negative inputs

Go 1.27.0 on darwin/arm64 returns `0x1p-29` for both `J1(0x1p-28)`
and `J1(-0x1p-28)`. The order-one Bessel function is odd. The pinned
implementation records the sign and makes the input positive, then
returns half that magnitude for inputs strictly below `0x1p-27`.
At the cutoff and the next representable magnitude above it, the
negative input again yields a negative result. The smallest negative
subnormal produces positive zero and the next produces the smallest
positive subnormal.

The differential checks the exact result words around this boundary and
compares the other finite inputs with an independent series and Hankel
reference. Run it after building the Lean executable:

```sh
GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/math/j1
```

## POTENTIAL — math.Jn can return infinity near a zero of J0

Go 1.27.0 on darwin/arm64 returns negative infinity for
`math.Jn(3, math.Float64frombits(0x40033d152e971b40))`. An independent
256-bit series gives approximately `0.19899990535769083`. At the two
adjacent input words the pinned results are approximately `0.1765760121`
and `0.1751340834`. The series remains close to `0.1989999054` on both
sides.

In the backward recurrence, the final normalization divides
`t * J0(x)` by the recurrence estimate `b`. For the central word, the
pinned sequence gives numerator word `0xbc7d8a1e96d116a3` and denominator
word `0x0000000000000000`. The differential preserves the pinned words,
checks the independent series at the witness, and records additional
rows that exceed its numerical tolerance close to J0 zeros.

```sh
GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/math/jn
```

## CONFIRMED — sort.SliceIsSorted documents a panic for accepted inputs

Go 1.27.0 on darwin/arm64. The function documentation says a nonslice
input panics. The pinned implementation instead calls reflected Len,
which accepts arrays, strings, maps and channels. For example,
`sort.SliceIsSorted([2]int{1, 2}, func(i, j int) bool { return i < j })`
returns true. A nonslice type without a reflected length, such as int,
panics before any callback. This is a documentation discrepancy.

Reproducer from the repository root after building the Lean executable:

```sh
GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/sort/comparison-entrypoints/review
```

The test calls the pinned public declaration with each accepted reflected
kind and checks the result and callback trace.

## POTENTIAL — encoding/binary.Decode panics on trailing bytes for two fast-path slice types

Go 1.27.0 on darwin/arm64. Decoding three bytes `[0,255,3]` into a
two-element bool slice writes `[false,true]` and then panics with an index
bounds error. The same source writes `[0,-1]` into a two-element int8 slice
and then panics. Byte slices and slices of a defined byte type consume two
bytes successfully. Decode checks that the source contains at least the
encoded size before passing the full source to decodeFast. The bool and
int8 slice loops iterate over the source length.

The documentation describes the number of consumed bytes as the result.
Confirm whether trailing bytes are intended to be allowed consistently
before treating this as a contract bug.

Reproducer from the repository root after building the Lean executable:

```sh
GOTOOLCHAIN=go1.27.0 GOWORK=off go run ./difftest/encoding/binary/public-buffer/review
```
