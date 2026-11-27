# Independent review of `bytes.Runes`

**Verdict: FAITHFUL**

Evidence SHA-256: `f5f835844802d5a3af7fd712d4dc2f5a0bf447b614cd368c925194e608684d31`

Reviewed author commit `e01cacb` against pinned Go 1.27.0 `src/bytes/bytes.go:1144–1154`. The count determines exact result capacity and the following decoder emits runes in input order. Nil and nonnil empty inputs both yield a nonnil empty slice. Malformed input contributes one RuneError for each invalid leading byte, and an encoded U+FFFD contributes one rune over three input bytes.

Nineteen literal pairs passed against native Go and compiled Lean, with rune equality, exact capacities and nilness checked for both inputs. The native controls also mutate every returned rune and check the original source and spare byte. Advancing by two bytes on RuneError changed three pairs. Allocating by byte count changed eighteen pairs. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/bytes/fold-runes/review` from the repository root.

Author evidence passed 72,644 boundary pairs and 60,000 seeded pairs at 17, 2718, and 65537. Allocation failure is outside scope.
