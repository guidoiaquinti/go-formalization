# Mul64 source correspondence

The Go 1.27.0 source computes two 64-bit words from four 32-bit limb products. The source guard in `make verify` parses the pinned toolchain's `math/bits` source and checks the formatted declaration against the transcription. It covers the signature, constant, expression grouping and statement order. A source change requires checking and updating the transcription. [Source guard](../../../../scripts/coverage/mul64_source_test.go)

In `mul64Source`, Go's `uint64` values are `BitVec 64` values. The mask constant is evaluated exactly before conversion. The bitwise AND and logical right shifts use the same widths as Go. Addition and multiplication wrap at 64 bits, and the Lean shadowing of `w1` represents the Go compound assignment. Review of these operation mappings connects the checked source transcription to the Lean definition. [Source model](../../../../lean/GoLean/Math/Bits/Arithmetic.lean)

Lean checks `mul64Source_eq_mul64` for every input pair. The proof reduces the low-word mask to remainder modulo 2³², the right shift to division by 2³², and the wrapping operations to arithmetic modulo 2⁶⁴. The resulting expression equals the existing arithmetic model. The theorem `mul64Source_spec` then proves that recombining the output words gives the exact unsigned product.

## Compiled Go evidence

| Configuration | Default lowering observed for the direct call | Execution status |
| --- | --- | --- |
| Linux/amd64 | Go SSA multiply with a high and low result, emitted as `MULQ` | Cross-built and inspected. Hosted execution pending. |
| Linux/arm64 | Go SSA multiply with a high and low result, with `UMULH` for the high word | Cross-built and inspected. Hosted execution pending. |
| macOS/arm64 | Go SSA multiply with a high and low result, with `UMULH` for the high word | Both lowering modes passed locally. |

The registered differential check compares the compiled Go call and compiled Lean source model with a shift-and-add oracle. Each execution mode covers 144 boundary pairs and 20,000 pairs for each seed `17`, `2718` and `65537`. With intrinsics disabled in the driver and `math/bits`, the check exercises source fallback. The default build gives sampled evidence for the compiler lowering on the executing host. Instruction semantics, compiler correctness and input-independent timing are outside the Lean theorem. [Go driver](../../../../difftest/math/bits/mul64/check.go) · [CI matrix](../../../../.github/workflows/verify.yml)

Run the source guard, proofs, differential checks and controls from the repository root with `make verify`. To inspect the default lowering for a host or cross-build, use the pinned toolchain and run the following. Set `GOOS` and `GOARCH` for a cross-build.

```sh
GOTOOLCHAIN="go$(cat .go-version)"
export GOTOOLCHAIN
scratch=$(mktemp -d)
go build -gcflags='command-line-arguments=-d=ssa/intrinsics/debug=1' \
  -o "$scratch/mul64" ./difftest/math/bits/mul64/check.go
go tool objdump -s '^main.compare$' "$scratch/mul64"
```

The [independent review record](Mul64.md) tracks the verdict, evidence fingerprint and earlier runs.
