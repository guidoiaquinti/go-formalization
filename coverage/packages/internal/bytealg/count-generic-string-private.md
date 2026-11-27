# Independent review of `internal/bytealg.countGenericString`

**Verdict: FAITHFUL**

Evidence SHA-256: `39b765d376fcd28974388ec386b5af9bc07827640c619314987f128569010fd5`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `82aa92b2ed377dd751745c6368a281318c5db5cff8fb63d633da5a31a4c7be29`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a2e650b549b9f4aade2acb12894740e03a46b12b675b3456025eec7d2867263a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6c59328a665ab5d68131bd80efc4235b92a3275cea40dd0e2041b90b613f87c2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2698f48d438b805550629cbc9f7def4d8aaacaf06bb28feb6dc6a6b15c22013a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4a9f3578db9a4cad7b60eee78adb112d7b03d5bebad095fa634cd1d5da167520`.

Reviewed 1abf5ea against Go 1.27.0. The pinned Go fallback indexes the string by byte and increments its counter on exact equality. The Lean fold counts the same raw bytes, and its proof identifies the result with the length of the matching-byte filter.

The independent control checks ten fixed inputs through the four actual pinned entry points and the Lean models. It covers nil and allocated empty slices, byte offsets crossing the scalar and vector thresholds, a malformed byte adjacent to valid UTF-8, and target-filled backing outside each view. A rune scan fails the C3 witness. Counting capacity rather than length fails the empty offset view and the vector-tail cases.

The full independent replay passed 3,587 boundary cases and 60,000 samples from seeds 17, 2718 and 65537. The verdict concerns exact count on the directly exercised Darwin arm64 profile. Assembly layout, memory access traces, timing, and unexecuted architectures remain outside scope.

Run the shared reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/internal/bytealg/count/review/main.go
```
