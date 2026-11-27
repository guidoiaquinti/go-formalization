# Independent review of `sort.Strings`

**Verdict: FAITHFUL**

Evidence SHA-256: `a155fb108351b02bd14ba66d02063dfa587a733bbc87aec02810785e5d4e9817`

I checked evidence commit `a6a654e` against Go 1.27.0. Raw Go string comparison orders bytes lexicographically, using length when one value is a prefix. Ranking strings by the number of lesser inputs preserves order and ties. The independent fixture detects a UTF-8 rune-key mutation using invalid high bytes, the encoded replacement rune, NUL bytes and proper prefixes.

I ran `go run ./difftest/sort/typed-sorts` and replayed 102 boundaries and 20,000 cases each at seeds 17, 2718, and 65537 against pinned native Go and compiled Lean on darwin/arm64. Nine additional literal cases passed both implementations with `go run ./difftest/sort/typed-sorts/review` across the six declarations. Exact float output observations preserve payload bits and reflect the pinned permutation. They make no stability guarantee for other inputs or versions. [Independent control](../../../difftest/sort/typed-sorts/review/main.go)
