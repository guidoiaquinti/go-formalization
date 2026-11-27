# Independent review of `encoding/binary.bigEndian.PutUint32`

**Verdict: FAITHFUL**

Evidence SHA-256: `e5ec69201b177da2b033ef11e79595c3b195e5a2da7c09a882d6d59f6145b924`

Reviewed a759a3f against Go 1.27.0. The pinned method bounds-checks the last destination index before its assignments. The model writes the unsigned word in big-endian byte order and preserves the full backing outside the destination prefix. A short destination panics before any bytes change, even when its capacity is sufficient.

The shared reviewer control checked 24 fixed native and Lean geometry cases per method. It tests unaligned views, empty and one-short destinations, and sufficient destinations with poisoned prefix and suffix storage. The expected byte layouts are literal asymmetric words. Reversing the order, writing through capacity, or doing a partial write before panicking would differ on these fixtures.

The full independent replay passed 65,715 boundaries and 60,000 samples with seeds 17, 2718 and 65537. The review covers returned or bounds-panic outcomes and exact post-call backing bytes.

Run the shared reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/encoding/binary/endian-put/review/main.go
```
