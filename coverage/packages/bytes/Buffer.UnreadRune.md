# Independent review of `bytes.Buffer.UnreadRune`

**Verdict: FAITHFUL**

Evidence SHA-256: `59b71c0190b91a21f4dd3b4ac8d4fe940f5581fc03d6e22050cd8b98c324e118`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `aadc37f722c65456af78251c5a7e788f4fe542730c4ea5d89186328c99c47579`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9207cedbad1b939803102514e05bba3a43057f8e77e7dfa7cd20d3fa1d216b36`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `41164aba113845ef339867c0ebea4adc9448b1304b13112ef00f5e7b07d31909`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c86436c204709618bf41cf0828bd4aa27bb16382963e0d3e6ce5b0a68c4084e3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `da9ca63a81e445f557ce37271b8c16091d2720ad66a6421762631da52f3f3c5b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b8ac71a570acfbf9efe17e95e0cdf5fde6d4dc91a3a4d97561ebafb93d06fcdc`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5a5fce6f3069b79650a0dd44e23e26aa894e921f2d1a813a326c32b2ba51ded2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `605a719b448459e8513853a08209c3460bf9e4d386631cf16d119e35f5c2b6bf`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `eb1a0e0b2385dd7bb73973fd0d9452afa75efe62ab1d7955fdfdd9ad013588bb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `46267a003b2ea12675f14ff8c76a3e7588f375ec74475b4d4f40f20b19a01dbb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d4e66eed6901a9271072fd04fc09ffad2d7e7f1a6981b660d049e1f041fce7e6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f15a6834beae989d4e4ff94d4884a0e99097d73c103da5e1b5898c475905431f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `58f28bfc53818c426e08db3e50e5b576231267cd5b17ea7785d3c1070ac8a590`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3e453cd2d3c61692a5a0de36883dffca4fd280815e7db2092b2144b26930c0b3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6af3ae559c3b9458e0da448ac52d1732742b26747028065695d68aeef696e194`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `12677549fc19fb231df83db2349a10cb60cf20f0d47039fd5f6aaa14666a471e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4de8cc7340fe2d2f043516f43459a74ae5671c8a0f3ae65e4ce608a4a986647b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ef5a5719de83314b13fca88e9bd145e8a16c7be5115fe7c7e5776590025a1561`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `731afd23f1fe3af77304565e1141ad82ad65286015a08d04468e043378f6a73f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `33736c2d5322e678f91daed2e6b1ca1c08ce9dc9d743999f5a9f119c5c2688ca`.

The pinned method rejects a nonpositive read flag with a newly allocated error and preserves the state. A positive flag is cleared on success, and the cursor rewinds only when its offset covers the width. The Lean model makes the offset guard explicit, including same-package states with positive widths that exceed the offset, and represents a nil pointer panic separately.

The direct darwin/arm64 comparison checks 33,281 boundaries and 60,000 sampled states from seeds 17, 2718 and 65537, including all representable signed read flags. The independent control passed 11 fixed native and Lean states and held three distinct errors live after repeated calls on one buffer. They had the expected text, preserved the cursor, and left the prior ordinary-read flag available for UnreadByte. Clearing a negative flag on failure failed two fixtures. Rewinding without the offset guard failed two.

```sh
GOTOOLCHAIN=go1.27.0 GOWORK=off CGO_ENABLED=0 go run ./difftest/bytes/buffer-unread-rune/review
```

[Model](../../../lean/GoLean/Bytes/BufferUnreadRune.lean) · [Comparison](../../../difftest/bytes/buffer-unread-rune/check.go) · [Independent control](../../../difftest/bytes/buffer-unread-rune/review/main.go)
