# Independent review of `internal/bytealg.Count`

**Verdict: FAITHFUL**

Evidence SHA-256: `5af3a4df4f20445699d4342d1466b6ce7005cedc3720ea48c64fec4fd1fd61dd`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bce3feb908d79ece23befa66f4a28767e80ad611d0f05a2b977785b0640a0754`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b67ade4471d962f1ed705b22c2869640b6d61c313f8837692d60dfb86d699118`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b14c013dfa4be881cbb5176402243604d70ad9371b47710bb8528950ddc3b86d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `99018e26e48503c20913b276fb5b22fe501eca1425b527d7f1b116408dc7f3bc`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2e8a8534aec2bdbbb43ccbd42298ba0f3f3e79d4d2272359150f65d58499d789`.

Reviewed 1abf5ea against Go 1.27.0. The arm64 slice entry loads the actual slice length and shares its counting body with the string entry. The Lean fold counts the same raw bytes, and its proof identifies the result with the length of the matching-byte filter.

The independent control checks ten fixed inputs through the four actual pinned entry points and the Lean models. It covers nil and allocated empty slices, byte offsets crossing the scalar and vector thresholds, a malformed byte adjacent to valid UTF-8, and target-filled backing outside each view. A rune scan fails the C3 witness. Counting capacity rather than length fails the empty offset view and the vector-tail cases.

The full independent replay passed 3,587 boundary cases and 60,000 samples from seeds 17, 2718 and 65537. The verdict concerns exact count on the directly exercised Darwin arm64 profile. Assembly layout, memory access traces, timing, and unexecuted architectures remain outside scope.

Run the shared reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/internal/bytealg/count/review/main.go
```
