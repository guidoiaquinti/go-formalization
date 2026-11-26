# Go formalisation in Lean

Lean 4 models and proofs for selected Go standard-library and compiler functions, checked against compiled Go.

## Run the checks

Install Make, Go 1.21+ and Elan with `lake` on your path. The first run may download Go 1.27.0 and Lean 4.35.0. Run from the repository root. [Go version](.go-version) · [Lean toolchain](lean/lean-toolchain) · [Elan setup](https://github.com/leanprover/elan#installation)

```sh
make verify
```

Run this command on your development machine with native Go and Lean toolchains. CI runs the same command on Linux/amd64, Linux/arm64 and macOS/arm64 on pushes and pull requests. It checks the generated files, builds the models and executables, audits the proofs, and runs every registered differential check and mutation control, regardless of `TARGET`. The Linux/amd64 job also runs the 32-bit Linux/386 `bits.Rem` wrapper. [Makefile](Makefile) · [CI workflow](.github/workflows/verify.yml)

## Explore coverage

The inventory uses Linux/amd64 v1 with cgo disabled as a reference profile for the Go standard-library and the Go compiler. It produces the same declaration list on every host with the pinned Go toolchain. The dashboard reports review progress across those selected declarations and lists each claim, exclusions and execution configuration. [Target definitions](coverage/targets.json) · [Coverage dashboard](coverage/README.md)

```sh
make coverage                     # print the dashboard
make inventory TARGET=math/bits   # list symbols, review states and pinned source lines
make diff TARGET=math/bits        # run this target's registered Go/Lean comparisons
```

## Formalisation workflow

For each declaration, translate the pinned implementation into Lean, compare the executable model with Go, prove a property that can distinguish plausible wrong implementations, and run that property against Go. An independent reviewer checks the source correspondence and tests the property against a deliberately broken implementation. Only a current `FAITHFUL` verdict contributes to reviewed progress. [Contributor guide](CONTRIBUTING.md) · [Formalizing CPython](https://broyojo.com/python-formalization/)

## What a proof establishes

Lean checks the stated property for every input admitted by the model and its hypotheses. Exhaustive differential checks establish agreement with the executable oracle across the stated domain for the tested build, assuming the driver, oracle and toolchains behave as intended. Larger checks provide sampled agreement. Source correspondence, compiler lowering, runtime dependencies and build settings determine how far that evidence transfers to Go. A Go-wide assurance claim would require a defined observation model and a verified connection from the language specification through each implementation and configuration it covers.

Go's `math/bits.ReverseBytes16` swaps the two bytes of a `uint16`, taking `0x1234` to `0x3412`. The model mirrors Go’s shifts and bitwise OR. The theorem `reverseBytes16_spec` extracts and concatenates the bytes. A second theorem proves the round trip. [Pinned Go source](https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L273-L278) · [Lean model and proofs](lean/GoLean/Math/Bits/ReverseBytes16.lean)

The differential check compares compiled Go and Lean with a byte-order conversion on all 65,536 inputs. In the historical review, the identity implementation passed the round-trip property on every word and failed the byte specification on all 65,280 words with unequal bytes. The review includes the Go control and its command, and limits its scope to return values. [Comparison](difftest/math/bits/reversebytes16/check.go) · [Historical review](coverage/packages/math/bits/ReverseBytes16.md)

Two other historical reviews cover every byte input to `math/bits.OnesCount8` and `math/bits.TrailingZeros8`, including what each function returns at zero. [Population count review](coverage/packages/math/bits/OnesCount8.md) · [Trailing zeros review](coverage/packages/math/bits/TrailingZeros8.md)

Use the coverage dashboard to check whether a historical verdict still matches the current evidence. [Coverage dashboard](coverage/README.md)

The `Mul64` worked example connects a guarded transcription of the pinned Go source to a Lean proof for every input pair. Separate differential runs exercise the source fallback and the compiler intrinsic on each execution host. [Mul64 source correspondence](coverage/packages/math/bits/Mul64-source.md)

## Contributing

The coverage prompt accepts a target and an optional function name. [Contributor guide](CONTRIBUTING.md) · [Agent instructions](AGENTS.md) · [Coverage prompt](prompts/extend-coverage.md)

## Credits

Maintained by [Guido Iaquinti](https://guido.iaquinti.com/).

The modelling and review workflow follows the CPython formalisation project, whose write-up shows how passing tests and proofs can miss a plausible implementation error. [Formalizing CPython](https://broyojo.com/python-formalization/) · [Code and review process](https://github.com/Broyojo/python-formalization)

Other inspirations include exhaustive comparisons with a language implementation, checking bit-level equivalence, and reproducing counterexamples in the original language. [Julia in Lean](https://github.com/MilesCranmer/julia-in-lean) · [lean-verify](https://github.com/hirosassa/lean-verify) · [Conao3’s Clojure study](https://s.conao3.com/blog/2026/e4f733) · [Koukyosyumei’s Rust study](https://zenn.dev/koukyosyumei/articles/31c8fcf4da45b7)

## Citation

The machine-readable citation is available in [CITATION.cff](CITATION.cff).

```bibtex
@misc{iaquinti2026formalizinggo,
  author = {Iaquinti, Guido},
  title  = {Formalizing {Go}: A Lean Model of {Go} 1.27.0},
  howpublished = {\url{https://github.com/guidoiaquinti/go-formalization}},
  year   = {2026}
}
```

## License

The project uses Apache License 2.0. The Go source mirrored in the model carries the Go Authors’ BSD license. [Project license](LICENSE) · [Go license notice](NOTICE)
