# Contributing to Go formalisation in Lean

Use the byte-swap contribution as a worked example. [Example record](coverage/packages/math/bits/records.json) · [Example review](coverage/packages/math/bits/ReverseBytes16.md)

## Choose a declaration

Set up the pinned toolchains and run this from the repository root with a target from the dashboard. [Setup](README.md#run-the-checks) · [Coverage dashboard](coverage/README.md)

```sh
make inventory TARGET=math/bits
```

The inventory uses a fixed source-selection profile to keep the generated dashboard reproducible. Build and differential checks use the contributor’s native toolchains. CI repeats them on Linux/amd64, Linux/arm64 and macOS/arm64.

Use the inventory’s exact symbol, including the receiver for a method, and follow its pinned source link. A target can own several packages, so locate the record by **package import path**. Check existing Lean definitions and evidence before editing, including for `UNSTARTED` symbols that may have unregistered models.

## Check an independent contract first

For bug discovery, choose a public behavior with an expected result you can derive from the docs, an exact calculation, or a stated mathematical property. Write a small Go-only check against the pinned public entry point before translating the implementation. Choose inputs on both sides of representation boundaries and include signed counterparts where relevant. Record the expected and observed results, the platform, and any assumptions about the public contract.

Keep the original expectation when Go disagrees with it. A pinned replay can assert that the same witnesses still occur and flag new disagreements. Record which witnesses the direct check found and which extra guarantees the Lean proof supplies. Keep the evidence near the package records and record confirmed or potential Go bugs in the findings log. [Findings](FINDINGS.md)

`make contracts` replays the `Erfcinv` and `Ldexp` fixtures with normal compilation and with intrinsics disabled in the driver and math package. A changed result, including a fix, fails the replay so the evidence can be reviewed. The full verification suite and `make verify TARGET=math` include this replay. [Erfcinv evidence](coverage/packages/math/Erfcinv.md#go-only-contract-check) · [Ldexp evidence](coverage/packages/math/Ldexp.md#go-only-contract-check)

To check the original contract assertions directly, run the probe without `-replay`. It reports every disagreement and exits unsuccessfully on the pinned version.

```sh
GOWORK=off GOTOOLCHAIN="go$(cat .go-version)" go run ./difftest/math/contracts/check.go
```

## Translate the pinned implementation

Put the model and its theorems under `lean/GoLean/` in the package namespace. Cite the pinned Go file and lines, including relevant helpers or table generators. Preserve observable behavior on the stated domain, including word widths, overflow, branch order, errors and panics where relevant. Document preconditions, scope limits and assumptions about architecture or generic types. Where relevant, include callbacks that panic or mutate state, changes made before an error, and behavior on the first call in a fresh process. For example, the `TrailingZeros8` model follows the table generator and handles zero explicitly. [Model and proofs](lean/GoLean/Math/Bits/TrailingZeros8.lean)

Keep the implementation's intermediate operations in the model. For example, Go's `Mul64` splits each input into 32-bit limbs, multiplies the limbs and propagates carries. The Lean definition follows those steps. Its specification independently reconstructs the full product from the two output words. [Go implementation](https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L470-L484) · [Lean model and specification](lean/GoLean/Math/Bits/Arithmetic.lean)

## Compare the model with Go

Run the compiled model against pinned Go before proving properties. You can build the executable with `make build` and run the new driver from the repository root before adding the registry record.

1. Add a Lean executable under `lean/GoDiff/` that runs the model. Add `lean-exe.json` beside the Go driver with the executable `name` and Lean module `root`. Run `make lake-update` before building an unregistered driver with `cd lean && lake build <name>`. [Example executable](lean/GoDiff/Math/Bits/ReverseBytes16.lean) · [Lake configuration](lean/lakefile.toml)
2. Add a single `check.go` under `difftest/<import-path>/<function>/`. Call the actual Go function directly and compare both its result and the compiled Lean model's result with the specification on each input. The shared helper checks the Go version and the number of Lean output rows. [Example check](difftest/math/bits/reversebytes16/check.go) · [Shared helper](difftest/internal/check/check.go)
3. Exhaust small domains. For larger ones, exercise branch boundaries and run at least 50,000 reproducible cases across several recorded seeds. Preserve failing inputs, outcomes and seeds when diagnosing a mismatch. Record the platform and entry point exercised. Enable the source fallback run where compiler intrinsics can replace the selected implementation. The runner clears ambient Go build flags and experiments, disables cgo, and records the normalized target and architecture settings. Use a source-shaped model or check alongside an abstract model when delegated implementation details are in scope. Describe which equivalences are proved and which are tested.

For Go declarations inaccessible to an external driver, put a Go test under the driver's `testdata/` directory and call the declaration from that test in its original package. The driver still compares the compiled Lean result with the same specification. Use `check.TestInPackage` to add the test through a temporary Go build overlay. The helper runs the named test against the pinned package and requires that it passed. The overlay adds test files only. When an in-package test would introduce an import cycle, add a test-only file that exposes the private declaration and use a second fixture in the external test package to call it. [Two-fixture example](difftest/math/normalize/check.go) Use `make inventory EXTERNAL_BLOCKERS=1` to list unstarted declarations affected by Go export and internal package import rules. Add `TARGET=<target>` to filter the list. [In-package example](difftest/internal/dag/newgraph/check.go)

## Prove a property and check the exact claim

State a property a plausible wrong implementation would fail, and exercise the exact claim against the pinned Go function. Quantify over the supported inputs and justify each hypothesis from the Go contract or a type invariant. If a claim fails, preserve the input, expected and observed outcomes, and reproduction command. Preserve the original claim with its counterexample even if you subsequently prove a narrower, justified claim. Check the model and the interpretation of Go's contract before attributing a discrepancy to Go. Keep model and harness discrepancies with the relevant evidence. Record confirmed or potential Go bugs in [FINDINGS.md](FINDINGS.md). Label each entry by status and include the Go version, platform, reproducer, and expected and observed behavior. For potential bugs, state what still needs to be checked to attribute the discrepancy to Go.

Use kernel-checked proofs. Avoid placeholders, `native_decide` and new axioms. The audit checks that each registered name is a theorem whose statement mentions the record's model declaration and allows only `propext`, `Classical.choice` and `Quot.sound` as axiom dependencies of loaded `GoLean` declarations. Leave an unfinished claim out of the theorem list and describe the gap in the contribution. [Proof audit](lean/Audit.lean)

## Register and verify the evidence

Use `limits` to identify source operations summarized by the model. If those operations are part of the claim, include a source-shaped model and state how it is connected to the abstract one.

Edit `coverage/packages/<import-path>/records.json`. Each key is one inventory symbol. Paths are relative to the repository root. Use the following fields. [Example records](coverage/packages/math/bits/records.json) · [Registry validation](scripts/coverage/)

| Field | Value |
| --- | --- |
| `go_version` | Exact version from `.go-version`, without the `go` prefix. |
| `model` | Path to the Lean module containing the model. |
| `model_symbol` | Fully qualified Lean model declaration named in each registered theorem statement. |
| `claim` | Observable property covered by the differential check. |
| `limits` | Excluded behavior, abstractions and whether the check is exhaustive or sampled. |
| `configurations` | Execution profiles and entry points covered. |
| `theorems` | List of fully qualified theorem names beginning with `GoLean.`. |
| `check` | Path to the Go differential check that runs with `go run` from the repository root. |
| `source_fallback` | Set to `true` to run the check again with intrinsics disabled in the driver and the selected Go package. |
| `control` | Path to the runnable reviewer program, which exits unsuccessfully if the real implementation or mutants produce unexpected counts. Required for `FAITHFUL`. |
| `review_sources` | List the additional files the reviewer checked, such as additional fixtures. Repository Lean models and executables, differential Go code, registered controls and build configuration are included automatically. |
| `review` | Path to the independent review, once available. |

Inventory states follow the record. A missing record gives `UNSTARTED`, a model alone gives `MODELLED`, and a registered check gives `AWAITING_REVIEW` until an independent review supplies a current fingerprint and a `FAITHFUL` or `HOLLOW` verdict. Run verification to confirm the registered evidence passes.

```sh
make coverage-update
make verify
```

Include updates to the generated `coverage/README.md` and `lean/GoLean.lean` with the record and evidence. For local iteration, `make verify TARGET=math/bits` runs the shared coverage check, Lean build and proof audit, and Go tooling tests, then runs differential drivers, source fallbacks and reviewer controls for the selected target. To build only the selected target’s registered drivers and run their differential checks, use `make diff TARGET=math/bits`. CI uses plain `make verify` to run all registered evidence on each platform, with up to four differential commands running concurrently per platform.

## Review a contribution

A reviewer who did not write the model checks the pinned source correspondence, tests the theorem claims on Go, and tries a plausible broken implementation. Check whether the model omits any source branches, intermediate operations, or externally visible effects relevant to the claim. Verify that the Go check exercises the theorem’s conclusion and covers its hypotheses. Record the scope, commands, cases or seeds, and a concrete input or failure count showing what the properties detect. [Worked review](coverage/packages/math/bits/ReverseBytes16.md)

Keep the runnable reviewer program beside the review, register it in `control`, and include its invocation. Assert the recorded outcomes in the program so CI can replay them. For architecture-dependent functions, identify the actual Go entry point and word width exercised. On Linux/amd64, CI also builds and runs a 386 oracle to check the 32-bit `bits.Rem` wrapper.

Save the review beside the package record. The first line and the bold verdict must match the parser's format, as in this example.

```markdown
# Independent review of `math/bits.ReverseBytes16`

**Verdict: FAITHFUL**
```

After inspecting the evidence, run this command and paste its output into the review. The fingerprint binds the verdict to the symbol, Go version and module files, Lean toolchain pin and Lake build inputs, theorem and model names, claim scope, reviewer control, workflow and repository source used by the evidence. The source fingerprint is conservative, so changes anywhere in the model, harness or verification tooling can stale other verdicts. A changed fingerprint requires a fresh independent review.

```sh
GOTOOLCHAIN="go$(cat .go-version)" GOWORK=off go run ./scripts/coverage fingerprint --package math/bits --symbol ReverseBytes16
```

Use `HOLLOW` when the evidence falls short and explain the missing behavior. Only `FAITHFUL` earns reviewed progress and it requires a registered theorem and control. The audit verifies that the theorem mentions the model; the reviewer checks that its statement actually entails the advertised property and that the driver checks that property. If an independent reviewer is unavailable, omit the `review` field. Rerun `make coverage-update` and `make verify` after recording a verdict. In the submission, identify the symbol, claims, test domain and results, review state, and any remaining gaps.

### Protect review approval

On the default branch, enable a GitHub ruleset requiring pull requests, at least one approval from a reviewer who did not author the model, dismissal of stale approvals, and the aggregate `verify` job as a required status check. Limit bypass permissions. Keep the independent approval on the pull request that records the verdict and fingerprint. The fingerprint detects evidence changes; the protected pull request supplies reviewer provenance. These settings require repository administration.

## Parallel coverage work

Assign each agent a distinct package and symbol and give it a separate Git worktree. From a commit containing the intended starting files, create both worktrees before either agent edits. For the paired assignments in [the coverage prompt](prompts/extend-coverage.md), run these commands from the same checkout.

```sh
git worktree add -b coverage/alphabetical-first ../go-lean-agent-a HEAD
git worktree add -b coverage/alphabetical-last ../go-lean-agent-b HEAD
```

Open one agent in each directory. Each worktree has its own index, Lake build directory and generated reports. Worktrees start from committed files, so check that the baseline includes any in-progress models, drivers and registry edits the agents need.

Integrate the finished branches one at a time. Preserve both drivers’ `lean-exe.json` files. After combining the source and package records, run `make coverage-update` to regenerate the Lake configuration, dashboard and Lean imports. Resolve any conflicts in these generated files using the regenerated versions. Run `make verify` on the combined tree and inspect both target inventories. Review fingerprints cover shared inputs, so check review states again after integration.
