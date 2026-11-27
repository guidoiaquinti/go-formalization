# Independent review of `bytes.Lines`

**Verdict: FAITHFUL**

Evidence SHA-256: `da80f463f86da1435c9f24e88734f861f89f01217850fcc0cddf977bcd8f24d9`

Reviewed author commit `01f389a` against pinned Go 1.27.0 `src/bytes/iter.go:18–33`. Each successful newline match includes the LF and advances the captured suffix before calling the callback. The final unterminated line also advances the captured suffix before calling the callback. The iterator clips capacity to length, and a terminal LF leaves no phantom line. The model tracks the suffix and the driver compares successive invocations.

Six literal fixtures passed against native Go and compiled Lean. They include consecutive LFs, a terminal LF, an unterminated final line, malformed bytes, nil input, and capacity clipping with spare backing. Delaying consumption until the callback returns true changed four observations. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/bytes/configured-iter/review` from the repository root.

Author coverage includes 1,952 boundaries and 60,000 sampled cases at seeds 17, 2718, and 65537. The claim covers sequential invocations with fixed bytes and deterministic callback budgets. Callback mutation, callback panics, and reentrancy are outside scope.
