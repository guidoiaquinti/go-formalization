# Independent review of `sort.rotate`

**Verdict: FAITHFUL**

Evidence SHA-256: `8049cc25c6adaf4644284dd2fdc0fc0af11f241bd897dee6bd36b4758da9cade`

I checked the Go 1.27.0 private implementation against the Lean trace model and the independent permutation oracle. The block-swap sequence reduces the longer side, then swaps the two equal residual blocks. Both unequal orientations, equal blocks and the all-empty case are checked. Triple reversal yields the expected rotated data but fails the callback trace. The native overlay calls the exact declaration, records every Swap, and panics if Len or Less is called.

The author replay passed 1,092 exhaustive geometries through width 12 and 60,000 seeded cases (17, 2718 and 65537) for this declaration. I ran `go run ./difftest/sort/range-moves/review` on darwin/arm64. Sixteen independent literal traces matched native Go and Lean. All six variant-specific mutants were rejected. [Reviewer control](../../../difftest/sort/range-moves/review/main.go)

Scope is valid index geometry with nonpanicking callbacks. Rotation requires two nonempty blocks or the all-empty case. Exactly one empty block is outside the pinned caller domain.
