# Independent review of `internal/bytealg.CountString`

**Verdict: FAITHFUL**

Evidence SHA-256: `80541a517f11c6aa5a26fdaac8195bc4da9cbcde19005182bc03c89098f106e2`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `660bc334fa1deb480a962e34290f1bf8ee073bdd063229a1971eac949fdcc02b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9e0642a4fd5dd3e81de0aa11b8c4f82fa8d59ca2a7aaf39003781c5398d933b5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5b8d9318126471a11ce796a0df111c0b4affa5fbc2abc3d7073da1eb9e20a873`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c6aa9283b120e78e1eae4c90832d574a26553cc7ae54579226d45b5a22b4d747`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9b1d12be7b434b3521f331b6dfd822c911648327afc840df93b3183da588a6a7`.

Reviewed 1abf5ea against Go 1.27.0. The arm64 string entry loads the byte length and shares the counting body with the slice entry. The Lean fold counts the same raw bytes, and its proof identifies the result with the length of the matching-byte filter.

The independent control checks ten fixed inputs through the four actual pinned entry points and the Lean models. It covers nil and allocated empty slices, byte offsets crossing the scalar and vector thresholds, a malformed byte adjacent to valid UTF-8, and target-filled backing outside each view. A rune scan fails the C3 witness. Counting capacity rather than length fails the empty offset view and the vector-tail cases.

The full independent replay passed 3,587 boundary cases and 60,000 samples from seeds 17, 2718 and 65537. The verdict concerns exact count on the directly exercised Darwin arm64 profile. Assembly layout, memory access traces, timing, and unexecuted architectures remain outside scope.

Run the shared reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/internal/bytealg/count/review/main.go
```
