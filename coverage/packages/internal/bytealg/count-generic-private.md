# Independent review of `internal/bytealg.countGeneric`

**Verdict: FAITHFUL**

Evidence SHA-256: `37f2e8f2a940bfac608368aee4979b6af0f3bbaf9d2120c24d88f8b60037390e`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a898c82a7251a166858648a4c177b85ccb15e430372a9878e41afcacc7bc49b6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9c053412f665e447121171a753dcc78bea566e3e8a0b19af013f5c69e0cc266a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8f94c66ed8397dd876ce09e8b1d0babb99c01d5d3a029b25e0bb4d2a9257b129`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9a9c46adea535ab2bd2e8658012bdd90f9d32d07892989ae8ef69be5cb6bfd00`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e3b4a52a785f7313099aec2ada392c07d5993d4a70c21321304afcab6c2ff393`.

Reviewed 1abf5ea against Go 1.27.0. The pinned Go fallback ranges over the supplied byte slice and increments its counter on exact equality. The Lean fold counts the same raw bytes, and its proof identifies the result with the length of the matching-byte filter.

The independent control checks ten fixed inputs through the four actual pinned entry points and the Lean models. It covers nil and allocated empty slices, byte offsets crossing the scalar and vector thresholds, a malformed byte adjacent to valid UTF-8, and target-filled backing outside each view. A rune scan fails the C3 witness. Counting capacity rather than length fails the empty offset view and the vector-tail cases.

The full independent replay passed 3,587 boundary cases and 60,000 samples from seeds 17, 2718 and 65537. The verdict concerns exact count on the directly exercised Darwin arm64 profile. Assembly layout, memory access traces, timing, and unexecuted architectures remain outside scope.

Run the shared reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/internal/bytealg/count/review/main.go
```
