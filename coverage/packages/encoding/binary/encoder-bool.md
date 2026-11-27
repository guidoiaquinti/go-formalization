# Independent review of `encoding/binary.encoder.bool`

**Verdict: FAITHFUL**

Evidence SHA-256: `0ba9d98e87d37c0ba4c6787a8c046c181a7060e6b0a3844b7208c896d78954d7`

Reviewed author commit `1556608` against Go 1.27.0 `src/encoding/binary/binary.go:776–783`. The Boolean write stores canonical 1 or 0 and preserves all other cells. Fixture arguments 255 and 128 select true. Success advances the offset by one. A reached out-of-bounds access or nil receiver panics. The native observer checks the original offset and full backing after an index panic, including a case where the offset is inside capacity and outside length.

The independent reviewer program passed eight literal fixtures against the exact private method and the compiled Lean model. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/encoding/binary/coder-byte/review` from the repository root. Writing 255 for true failed two fixtures. Reporting an advanced offset after an index panic failed four fixtures across the group.

Author evidence passed 7,171 boundaries and 60,000 sampled rows with seeds 17, 2718, and 65537. Scope is sequential access with valid slice geometry and native signed offsets. The author and reviewer executed the native 64-bit host profile.
