# Independent review of `encoding/binary.encoder.int8`

**Verdict: FAITHFUL**

Evidence SHA-256: `f4b35e35d40043d401bf2b90385aec6337e2af0fd9deb508b0d4795c0e839851`

Reviewed author commit `1556608` against Go 1.27.0 `src/encoding/binary/binary.go:791–794 and 831`. The signed input is written as its exact eight-bit representation. The fixtures use −128 and −1 alongside zero. Success advances the offset by one. A reached out-of-bounds access or nil receiver panics. The native observer checks the original offset and full backing after an index panic, including a case where the offset is inside capacity and outside length.

The independent reviewer program passed eight literal fixtures against the exact private method and the compiled Lean model. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/encoding/binary/coder-byte/review` from the repository root. Zeroing negative signed input bytes failed two fixtures. Reporting an advanced offset after an index panic failed four fixtures across the group.

Author evidence passed 7,171 boundaries and 60,000 sampled rows with seeds 17, 2718, and 65537. Scope is sequential access with valid slice geometry and native signed offsets. The author and reviewer executed the native 64-bit host profile.
