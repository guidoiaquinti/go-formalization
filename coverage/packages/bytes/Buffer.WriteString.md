# Independent review of `bytes.Buffer.WriteString`

**Verdict: FAITHFUL**

Evidence SHA-256: `fe1ed87e3286582d45f2d08381f773cca39e7b9ea9f5886f24d0103e2b592ef9`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6be11d7776e3ee5f0eba6f3f8bbe02326015d7e5ccaede992d40004da0bcfa55`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e002d6407f3bda5a580480f62a3040c30358961f3b95de8e6d0f15b8b04ac67b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e8d55fdb202333bc08dacac8e9d955d62c05124aaee2507ce6ef176377488cc2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `015b64e6088ce78ec42403e9f44f3217aaa77749e33fdcbbda6d7c6bb614312c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3ff9220bea5412ab0ffb761d7c25fd3182650a8aebcd4e5a366217bd77d4bd7a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bc1d0491a68957f438a31016b2542b2fd14ee8f10ce5263e685982c164ffa130`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e46bc026538d60bf6459a742ed227e6a7226f585b8982f18d97b87ea0ec6bfc1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c0b52bc2d83f02594de3e6f7447c26e84d206b1d77ab37387a49166e2cc9a27c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `00d010d98defe1c01ec22abaace7caf3470f68933d0429832a88424e17743cd0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f552a62045176b4053cec330c6e9ed41ea69e4a21d2cc3021cce5b0fb6081719`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a6b0e8c83b9893c8617a92ff9691e30639a30871957bc6bd3fb2e8f45ec84085`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c635d976a41f9a9e148a72428ad552a35f4162092b3ec9bf5b3117f16364f84e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5a2065acaa62dd04c83b50cd676b94f12aa4999dac640431c46e1b40da894464`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8c78fe1eae756d3cfcfd973916e5d7ee22549aba6299ab0f996bd1ed9b571527`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `010e67fefc6c8b91fe0b8f972aa792e97a0b7f76777494472220b0f1570da543`.

Pinned Go 1.27.0 clears the previous read flag before attempting to reslice. Empty writes always take that path, preserving cursor and nilness. Fitting nonempty writes preserve the cursor and unused bytes, while a failed fit can reset an exhausted buffer, compact at the inclusive half-capacity threshold, or allocate. The model reuses the shared growth logic with detached string bytes, preserving the source across compaction and returning its byte length. [Model](../../../lean/GoLean/Bytes/BufferWriteString.lean) · [Shared growth](../../../lean/GoLean/Bytes/BufferWrite.lean)

The direct Go overlay and Lean matched all 21,714 boundaries and 60,000 samples from seeds 17, 2718 and 65537. Eleven independent fixtures and two public state checks passed, including nil receiver, empty exhausted buffer, malformed string bytes, inclusive compaction, and the first allocation above 64 bytes. Incorrect variants that reset on an empty write, allocate at the compaction boundary, or conflate nil with nonnil empty storage each failed a fixture. Unsafe mutation of string backing, allocation failures and growth requests above 512 bytes are outside scope. [Comparison](../../../difftest/bytes/buffer-write-string/check.go) · [Reviewer control](../../../difftest/bytes/buffer-write-string/review/main.go)
