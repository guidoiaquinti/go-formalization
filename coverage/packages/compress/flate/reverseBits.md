# Independent review of `compress/flate.reverseBits`

**Verdict: FAITHFUL**

Evidence SHA-256: `85a04dc863761c26494e3b3e16dbf081d960f1a12f074cf4675168e6c7dd2702`

I checked evidence commit `9848145` against Go 1.27.0 `src/compress/flate/huffman_code.go` and the reused Lean word reversal. The byte subtraction wraps before masking the shift with 15. The word shifts in 16 bits and then reverses. The Lean theorem identifies the source bit or shifted-out zero for each output position. The pinned Go checks independently compare the full word against reversal of the selected low bits, with 16 bits selected when the count is a multiple of 16.

I ran `go run ./difftest/compress/flate/reverse-bits` on darwin/arm64. All 9,728 boundary cases covered the 256 count values, and 20,000 sampled cases each passed with seeds 17, 2718 and 65537. Six independent native and Lean fixtures passed with `go run ./difftest/compress/flate/reverse-bits/review`. Omitting the shift mask, ignoring the requested width and treating zero residue as a zero-bit selection failed four, four and two fixtures respectively. Scope covers every typed input to the helper. Huffman table construction is outside scope. [Independent control](../../../../difftest/compress/flate/reverse-bits/review/main.go)
