# Independent review of `bytes.Buffer.ReadByte`

**Verdict: FAITHFUL**

Evidence SHA-256: `1ace6b723f07a537b6f61297c9b71502cfb9b467f1fe1e45e844af3bf4c1b4f4`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c9e957add3786866e1c0063c36d22d590f7616733784b5bcfec2dd314c9d4b55`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ab0e9f20b7dc50857aa2c878ce6773c4fe3cd24a2eba24e53f95f76075086997`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f77dfba82817e7344f26947d66625f7866626a05fda78231ca31e57a04884bc3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c1d2b220ccfdd6d55a536560279eb6f989e018f888a655735f1fc384d8816ab1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b8ccc5f9a966473fca9adef05a8a5fa6595118b2e60eb0983c05ed4c1625ab18`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `95536a9ec6da8fd149bfd6874644c96353ca5dd03568a12840e66e816d9cb5e0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `88049254555d298b84d31c50e58387e9ebd1ed933a65c2f6b3ad753deacd67f4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `da5f91be9b07e8a2312783f3d22916a0502a8eb7694de9bcdce8cab9e8373045`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1baf9e0833c4e3a99331f6470e3d5f3bb0e0a1ec101cb078183712d3762123c8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `52c6e433726d9fb9e7b10ff38525e880568d5ceb7f9d449d73ecf5aa7ed72213`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6a8da6824b1be1493d84305b89d9df63cd7998c949620afcdc965d4abd036712`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5f55a1651e0056da9d0751fef1de286d13e0839e2d57ef3253224bfe3fcd9174`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `35f21e782706fd379184571aa34b44dcd2bb21220dd461a5adb432505742a0f7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6256770bb21fe1a3e73752218fe486eb824dc317efad6e9bf5bdf6673f5dea46`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `77be506c29e0314c63bbdb8aef00447d23824887c3845f47f382b2b1c2e07d63`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2aeeba1fac71c73a03369c53948cee60300c5d9e18f564ba4bca26a963c07a2b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d266b2f7d6e842c959f8c3c3dfd9dbddf872788e0599a3da0806a87eb72480b7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b8b5f85b892e17f6ee626744cb7bcb15c26cdc7b83759c5a5df0a41a1ab5115d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a411bfa63a3961a1b04092bb23dc3fc7d392c68c82c6596ea6a7dde6df96684d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a2d76837a54a98d4efaf6800bc32a004f889fcdb1a864a779aaccdeac98d033c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `272e93ad03de97254870466a2b8965eaceb081d6deb04a7be81a7c568ee98b64`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `db68c67b2246a9db77ac77fb0c56b9caf60462f36b56b9e6f538fd737d430e69`.

The pinned method returns the byte at the read offset, advances by one and records an ordinary read. At exhaustion it resets the buffer before returning byte zero and the exact `io.EOF` sentinel. The reset frees the consumed prefix for reuse and preserves the backing data and storage nilness. A nil receiver triggers the runtime pointer panic.

The direct darwin/arm64 Go and Lean comparison passed 3,121 boundaries and 60,000 inputs using seeds 17, 2718 and 65537. The independent Go and Lean control passed 11 fixtures. It includes a mutation through an unread aliased view before the byte read, reread after UnreadByte, and ordinary and rune histories on both success and EOF. Omitting the EOF reset failed two fixtures, retaining the old byte-unread flag at EOF failed two, reading at offset zero failed three, and setting a rune-unread flag on success failed five.

```sh
GOTOOLCHAIN=go1.27.0 GOWORK=off CGO_ENABLED=0 go run ./difftest/bytes/buffer-read-byte/review
```

[Model](../../../lean/GoLean/Bytes/BufferReadByte.lean) · [Comparison](../../../difftest/bytes/buffer-read-byte/check.go) · [Independent control](../../../difftest/bytes/buffer-read-byte/review/main.go)
