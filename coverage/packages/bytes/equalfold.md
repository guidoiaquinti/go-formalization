# Independent review of `bytes.EqualFold`

**Verdict: FAITHFUL**

Evidence SHA-256: `d8ea9a28a7cfe8b616e8c5d46ed10439449e0a168dc5c648ab0962d7a489d7a4`

Reviewed author commit `e01cacb` against pinned Go 1.27.0 `src/bytes/bytes.go:1213–1285`. The shared string model preserves ASCII matching, Unicode simple fold and one-byte advancement for malformed input. A malformed byte may match a complete encoded RuneError, while two malformed bytes require two decoded counterparts. Full case-fold expansions and locale-specific equivalences are absent from the pinned routine.

Nineteen literal pairs and four exact successor transitions passed against native Go and compiled Lean. The cases include both orientations of all twelve distinct pairs of U+0345, U+0399, U+03B9 and U+1FBE. Requiring equal encoded byte lengths changed nine pairs, and truncating the simple fold walk to two steps changed two pairs. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/bytes/fold-runes/review` from the repository root.

Author evidence passed 72,644 boundary pairs and 60,000 seeded pairs at 17, 2718, and 65537. The author compared all 1,114,112 fold successors with pinned Go and verified a maximum orbit of four.
