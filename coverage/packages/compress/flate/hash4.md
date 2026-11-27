# Independent review of `compress/flate.hash4`

**Verdict: FAITHFUL**

Evidence SHA-256: `e9f5c728d4bc18949a0892814cfcbddd9bd57f9de935cdec3950926786d6c3ec`

I checked evidence commit `4fc8602` against Go 1.27.0 `src/compress/flate/deflate.go` and `src/compress/flate/load_store.go`. The two-index slice in loadLE32 uses available backing capacity. A length-zero view with four capacity bytes hashes those bytes, and a three-byte backing panics. The word is little-endian and the hash keeps the high 17 bits of the wrapped product. The proof identifies the first four bytes independently of the suffix.

Pinned native Go and Lean agreed on 6,006 fixed cases and 20,000 cases for each of seeds 17, 2718 and 65537 on darwin/arm64 using `go run ./difftest/compress/flate/hash-window`. Six independent literal rows passed with `go run ./difftest/compress/flate/hash-window/review` and rejected seven wrong observations. Scope covers disjoint source and destination in legal slice geometry, the full destination backing and the bounds failure class. The wording of panic messages and concurrent mutation are outside scope. [Independent control](../../../../difftest/compress/flate/hash-window/review/main.go)
