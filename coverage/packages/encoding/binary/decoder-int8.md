# Independent review of `encoding/binary.decoder.int8`

**Verdict: FAITHFUL**

Evidence SHA-256: `6efc9d285675ff47af00d133bef31c588e2c8dfbd65658791bfb5feb48e1a5e3`

Reviewed author commit `1556608` against Go 1.27.0 `src/encoding/binary/binary.go:785–789 and 829`. The signed conversion sign-extends the bit pattern. The fixtures decode 128 as −128 and 255 as −1. Success advances the offset by one. A reached out-of-bounds access or nil receiver panics. The native observer checks the original offset and full backing after an index panic, including a case where the offset is inside capacity and outside length.

The independent reviewer program passed eight literal fixtures against the exact private method and the compiled Lean model. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/encoding/binary/coder-byte/review` from the repository root. Treating the signed result as unsigned failed two fixtures. Reporting an advanced offset after an index panic failed four fixtures across the group.

Author evidence passed 7,171 boundaries and 60,000 sampled rows with seeds 17, 2718, and 65537. Scope is sequential access with valid slice geometry and native signed offsets. The author and reviewer executed the native 64-bit host profile.
