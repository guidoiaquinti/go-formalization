# Independent review of `bytes.splitSeq`

**Verdict: FAITHFUL**

Evidence SHA-256: `c07d4398fd149af1d9752b238f94d7fffe9730ba922032681461ced72be9fd04`

Reviewed author commit `01f389a` against pinned Go 1.27.0 `src/bytes/iter.go:36–60`. For a nonempty separator, the two-index fragment slice may include spare backing before its capacity is clipped. Rejected yields leave the current suffix intact. At a bounds panic, previously yielded fragments remain visible and a resumed call starts at the failing suffix. The final yield keeps the captured suffix for subsequent calls. Empty separators use decoded rune width, including width one for malformed UTF-8.

Six literal fixtures passed against the exact private Go declaration and the compiled Lean model. They include spare capacity, prior yield then panic, repeated nil and empty final fragments, and restart after rejecting a malformed-byte fragment. Limiting the endpoint to length changed one observation. Consuming after a false callback changed four, and clearing the final suffix changed three. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/bytes/configured-iter/review` from the repository root.

Author coverage includes 1,952 boundaries and 60,000 sampled cases at seeds 17, 2718, and 65537. The claim covers sequential invocations with fixed backing and deterministic callback budgets. Callback mutation, callback panics, and reentrancy are outside scope.
