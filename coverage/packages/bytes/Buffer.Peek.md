# Independent review of `bytes.Buffer.Peek`

**Verdict: FAITHFUL**

Evidence SHA-256: `b8f6005e60085d34c7694a854f86daa7f7be640bd60cf694e560f861ee488352`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8f05cda795c2ef294dc1a9b9971daf410d61a9fdf941f8366b44bd5a95e56b1e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3b6e13eb8b86b9f5505e61936aee6dbc3d4cb7f48bcd9a1079af23441eb06791`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8de728c71be1cc27426a2f580f23dcd6fa28c91508f87098af465957665c64c5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5816ff425c49935983c2bd839540e54e7ed6750b600e8a45128d26c15e116b31`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ea5ed952cb320c99b94b22700473767fb9aa447af77f2540071ae5ebb4767156`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `47406b8cfb7168ff86af79eae067ced3596ec770617e089549acf5f374d9dfd2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3dc8869a551b3add84f178015a053fc89c736ebab1b9ebe424938146e6fa6859`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0fa268657c9d0778424bf41395b858dad63bcdddd37070de3e744a844bb60684`.

Go 1.27.0 returns a two-index view from the unread offset, so capacity extends through the original end of backing storage. Counts greater than unread length yield the full unread view with the exact io.EOF sentinel. Equal counts succeed, and negative counts panic before changing any state. Nil backing, nil receiver, spare bytes, cursor and previous read flag are preserved as applicable. All 19,643 boundaries and 60,000 samples at seeds 17, 2718 and 65537 matched. [Model](../../../lean/GoLean/Bytes/BufferPeek.lean) · [Comparison](../../../difftest/bytes/buffer-peek/check.go)

Nine independent fixed fixtures passed, including signed count extremes, nil backing, exact length and a zero-length view with positive remaining capacity. The public control mutated a retained spare byte through the short view, successfully unread the preceding byte, and checked that the overlong call returned the original backing and io.EOF. The fixtures distinguish an inclusive EOF condition, clipping returned capacity and treating negative counts as zero. [Reviewer control](../../../difftest/bytes/buffer-peek/review/main.go)
