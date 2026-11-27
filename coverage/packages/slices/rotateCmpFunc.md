# Independent review of `slices.rotateCmpFunc`

**Verdict: FAITHFUL**

Evidence SHA-256: `a0f7093dc078aae0eb7e847d9780bacc6d57edeabc0fc9690f7cd373e1d82e38`

I reviewed evidence commit `e596c81` against the private Go 1.27.0 loops. The model starts with lengths `m-a` and `b-m`, swaps the same front or back block at each branch, and uses the same final swap at equal widths. With both blocks positive, each unequal step strictly decreases their sum while keeping both widths positive. No active block extends beyond the original range. The Lean fuel suffices on that domain. Both empty blocks cause a zero-count swap. The comparator variant passes its comparator through to the pinned swap helper, whose body never calls it.

I ran `go run ./difftest/slices/rotate-range` and reproduced 1,124 boundaries and 20,000 per seed 17, 2718, and 65537 on pinned native darwin/arm64 and compiled Lean. My separate `go run ./difftest/slices/rotate-range/review` passed seven literal fixtures per helper and rejected two block-swap mutations. Omitting the final swap leaves `[90,1,2,3,4,91,99]` in the equal-width witness where the required backing is `[90,3,4,1,2,91,99]`. Swapping against the start of a longer right block leaves `[90,1,2,3,4,5,91,99]` where the required result is `[90,2,3,4,5,1,91,99]`. The fixtures also cover longer left blocks, both-empty interior ranges, nil and nonnil empty slices, spare backing, and a panic-on-call comparator. [Independent control](../../../difftest/slices/rotate-range/review/main.go)

The evidence covers adjacent nonempty blocks and the both-empty case, with native comparisons instantiated at `int`. Exactly one empty block and invalid signed ranges are outside the claim. Lean proves a general both-empty result, comparator independence and equivalence of the wrappers, and a concrete unequal witness. The general block-order correspondence is checked by source inspection and sampled differential tests.
