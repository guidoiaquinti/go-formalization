# Independent review of `encoding/binary.decoder.uint8`

**Verdict: FAITHFUL**

Evidence SHA-256: `efe35e805201396f9c4be6b9c3fea0dc7db64b4ce1cea7aa0d042f964196ae56`

Reviewed author commit `1556608` against Go 1.27.0 `src/encoding/binary/binary.go:785–789`. The decoded byte keeps its unsigned value. The fixtures decode the high-bit bytes as 128 and 255. Success advances the offset by one. A reached out-of-bounds access or nil receiver panics. The native observer checks the original offset and full backing after an index panic, including a case where the offset is inside capacity and outside length.

The independent reviewer program passed eight literal fixtures against the exact private method and the compiled Lean model. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/encoding/binary/coder-byte/review` from the repository root. Using backing capacity as the index bound failed two fixtures. Reporting an advanced offset after an index panic failed four fixtures across the group.

Author evidence passed 7,171 boundaries and 60,000 sampled rows with seeds 17, 2718, and 65537. Scope is sequential access with valid slice geometry and native signed offsets. The author and reviewer executed the native 64-bit host profile.
