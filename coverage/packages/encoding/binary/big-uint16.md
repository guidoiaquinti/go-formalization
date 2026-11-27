# Independent review of `encoding/binary.bigEndian.Uint16`

**Verdict: FAITHFUL**

Evidence SHA-256: `086a751104e6d027e5f14c7ff176b3478a7d2f79a74204113e0adcb32145941c`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bee52a1584220a572a7a7905c81db31509c6f8b1e8c3eeda914afaf9bc6b417f`.

Reviewed 2fb84a1 and proof registration fix 989358f against Go 1.27.0. The pinned bigEndian.Uint16 checks input length before shifting and combining the first 2 raw bytes. The model matches the endian direction, ignores later bytes, and returns a bounds outcome for short inputs.

Nine independent fixed cases passed for both Go and Lean. The fixture covers nil, allocated empty, one-short inputs, unaligned views, asymmetric byte patterns and poisoned following bytes. It preserves caller backing unchanged. Byte-order reversal and using backing capacity to accept short input each disagree with a fixed witness.

The full independent replay passed 68,100 boundaries, including every 16-bit prefix, and 60,000 samples with seeds 17, 2718, and 65537. The scope concerns the exact scalar result and the bounds-panic class.

Run the shared reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/encoding/binary/endian-read/review/main.go
```
