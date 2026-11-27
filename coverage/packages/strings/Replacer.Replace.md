# Independent review of `strings.Replacer.Replace`

**Verdict: FAITHFUL**

Evidence SHA-256: `72e833de4e0fa449a425412143feb041eded9766fa5862924a2416dc98677ec1`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `862e8d0293f5de5b046ee5022b2df4241eae13e0993484c52f28b1070c3be4c5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ce475e77db0b9bf82eabe8f0288e2762ee4469468c5cd29a6d48e58988bc1604`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1ef11553d46baa3db1366b09ac1cf336381ce68728a10eabfd01b2a631f10922`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4adcba460a55da6a8e51fa05d3e2709e6ce4c94c6b80fc55455b5f27f0d51af9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `24da8a1ba2fae4d1de288a2f450457fc421645a40e98837ec931ce3badc922b2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f33703f3414f71f9fc9a31935f33c988d3970bd66f490ff4838787418d6dab1d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `135d4bbb846cdeb2dbe85221953c5e516d0458258fdae1fc5caa75592c2f0e73`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bd2a840dbb20a9343e7681ebd1196f3cdc7e1c3e3290ac656c2febabf54fe681`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `56b9d7c9467a4672fa55349b1c390aece912f75004b9182e7d7408eb62a450d4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f4cc290bbb3966af9a3623f6b744e943dd41c670ea117028bfbbb3195b0788e8`.

Reviewed 470a086 against Go 1.27.0. The wrapper calls sync.Once before dispatching to the built implementation. The sequential model represents this as pending and ready states, with unchanged implementation selection on a later call. Nil receivers panic and the zero value uses an identity byte table. Each selected implementation is delegated to its existing output model.

The independent control checks ten pinned native and Lean sequences covering nil and zero receivers, absent pairs, all four implementation families, first duplicate priority, nonoverlap, empty substitutions at raw byte positions, and malformed UTF-8. Matching empty keys only at rune boundaries changes the second result for the valid multi-byte input `é`. Rescanning replacement bytes changes the result for `a→b, b→c`. Two consecutive calls use the same receiver in every nonnil fixture. Dispatch caching follows from the pinned sync.Once wrapper and the model transition. The public result comparisons exercise first and later calls without observing private allocation or counting builds.

The differential checks 1,172 boundaries and 60,000 paired sequences with seeds 17, 2718 and 65537. The claim covers sequential calls to nil, zero, and constructor-created receivers with successful allocation. Concurrent initialization, mutated private state and backing identity are outside the scope.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/strings/replacer-replace/review/main.go
```
