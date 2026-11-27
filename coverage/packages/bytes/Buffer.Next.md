# Independent review of `bytes.Buffer.Next`

**Verdict: FAITHFUL**

Evidence SHA-256: `b0448392907a6a9faf44cb4217e6f8b6de3e981bf421db9b459efd9c7669efd6`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c6ff7648dc5fbc34f83add8eacd1ae45bb5de9261eadb44cfbb3f3ccd8aea022`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d505bf89936951284f4c44416a8d6cbde08ec213ed4e29687c78805127889774`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e43bc1bf930ecf5b2e86245d411e578117ec0894f0457d153e01b0fe806d0b9b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e4045c33161ac7b1a7cf2f2d7820f42cf57238b13266918d10d7273b95a5d855`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d48f5ee8ebe221205c4543052fd3e85ead63ccc312da8cd99814075e24131550`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ed30af98ec9d54382b739fdef93e9e850dc919735f30c98e7c3f45ed97b3144d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d48f0accfa2ca6a86ee15f1c2b737176e292f98b3c6d730a7759add422e016d5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f51bd99af76bdb87ec16d09b0199d1b2d04730713d7dbafec777fcf83b0a7f32`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9ea1e946e976e3226b3ef5c2cb4ce030472347969668b590053fdc86dad02201`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1d6c220192147628a89950f428734df6323a1aa90c1ae9b6cb09170320ad9327`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2fcd5847dce4942c377f566aaaef7c8988520997a6fc16623094d9b71f66220a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `84c412919bf91fe0573b0f3641377e000753c59fec1e7655268e9063cbfd7a9f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e05b7ce58afedab31da955011b753704cd5df1e36f4d7cae871c0ebd5a87c419`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `df3b2d870ec76ef86b3f4fb3f6c30e5a75415fa2f27d5caebd81c52f8b9e59a6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9365f1447068fa4d432b96a811e817ce694ef6b7c41ab176d1d2853c91c8d6c6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `635b2b67185dfef04b3d559c756658d38d70704e5b0dd59a7140aa2a5f5ef799`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5aab8bbafc233bc9e25091624da3d59be0740422600f74cfa82e3e0cb1cf2b86`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9ad6a4aa2817973bdcdc28290c68bbed8fae23d095ee29f0af38377d13f13a19`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `98d6cab71c92a604a0f96d860d6899b2c62c78cd52ec5442b4157ade06ce977d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `702347f70f118c0206eb3a6318ab0e6544819b4bbc055f3eca56d0473533b94e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9782215abb63f3674d8b4ac8f0c031567a768d6e0257ea193d578dbfda798eee`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7a576cc640bd54325a71f9534a2ce960a7733908f618a0d84a805924bcb0f251`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `87bb8be2473705a5a23daa697cf20a97926a4f9525a72474cfa63fdcf472982a`.

The pinned method invalidates the prior read before obtaining the unread length and slicing. Nonnegative counts are clamped to the unread length, the resulting view begins at the original read offset, and only a positive returned length permits UnreadByte. A negative request triggers a bounds panic with the offset and backing unchanged and the previous read invalidated. Nil receiver panics are classified separately.

The direct darwin/arm64 comparison passed 11,056 boundary cases and 60,000 sampled cases from seeds 17, 2718 and 65537. The independent control passed 18 native and Lean fixtures, including the signed minimum and maximum, exact and adjacent clamp boundaries, empty views with retained capacity, and prior ordinary and rune reads. It mutates the view before any subsequent buffer operation and checks that the original backing and reread byte reflect the change. Copying the view, failing to clamp, leaving the byte-unread flag active and leaving the rune-unread flag active failed five, five, six and three fixtures respectively.

```sh
GOTOOLCHAIN=go1.27.0 GOWORK=off CGO_ENABLED=0 go run ./difftest/bytes/buffer-next/review
```

[Model](../../../lean/GoLean/Bytes/BufferNext.lean) · [Comparison](../../../difftest/bytes/buffer-next/check.go) · [Independent control](../../../difftest/bytes/buffer-next/review/main.go)
