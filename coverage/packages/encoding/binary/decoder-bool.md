# Independent review of `encoding/binary.decoder.bool`

**Verdict: FAITHFUL**

Evidence SHA-256: `eb43206db597e83d9cd4b16dba59905bb075570cda5b060d814e803cd3d8acb2`

Reviewed author commit `1556608` against Go 1.27.0 `src/encoding/binary/binary.go:770–774`. Every nonzero byte decodes to true. The fixtures decode 128 and 255 as true and zero as false. Success advances the offset by one. A reached out-of-bounds access or nil receiver panics. The native observer checks the original offset and full backing after an index panic, including a case where the offset is inside capacity and outside length.

The independent reviewer program passed eight literal fixtures against the exact private method and the compiled Lean model. Replay using `GOTOOLCHAIN=go1.27.0 go run ./difftest/encoding/binary/coder-byte/review` from the repository root. Treating only the byte 1 as true failed two fixtures. Reporting an advanced offset after an index panic failed four fixtures across the group.

Author evidence passed 7,171 boundaries and 60,000 sampled rows with seeds 17, 2718, and 65537. Scope is sequential access with valid slice geometry and native signed offsets. The author and reviewer executed the native 64-bit host profile.
