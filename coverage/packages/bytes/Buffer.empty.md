# Independent review of `bytes.Buffer.empty`

**Verdict: FAITHFUL**

Evidence SHA-256: `c9f4977133b688142901cf54c7edad2917aa98b898352e4f59ea0ac9019f7d9f`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1d15722989a24b6ead33ccff936a1b4f56c45bf1eb28ef8a4a7a82a46ad96f6e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `161cbdeb92d23bdaca6d7a8527a6a3f2a57a34b4bccc513f70f8cbc0b3f84f01`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f08a94e7f522b47bfd00ae80177722d45be0a758577d42f1e0568acfa4da4ca2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2979ede653f208b0a853b58d5d58956b425e181e2abffa66516c6218c8e96640`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e0571538f9ea813f007f35b98229d38b8f48cab51a1bf6610876bea565ba60b8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1163e8bcff6af887adb8c5f9f9a2a14f6d8a183aa13ab5aeb73eb9df5b5d5dc2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `69ff7fb851802ad780e77ca6d424f88e2b1f521f044f03ec9eb9aaea9f6773c7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4352271059cb22e8a5f956fb0ae54cfc408968072ac09d5215cbc43d7e554812`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f80975c78bc9ba7b352fcf5cfefd53027fb17be965f83eab406541b5e0e3e014`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f5126f239a830fee67267172aecc8c5ae520a1c7cffa2893e28cabdf05df7c48`.

Go 1.27.0 asks whether the stored slice length is at most the signed cursor. Neither spare capacity nor nilness of the slice changes the result. A nil receiver panics, a negative cursor returns false, and a cursor past the slice end returns true. Pinned Go and Lean matched 4,101 boundaries and 60,000 samples at seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Bytes/BufferEmpty.lean) · [Comparison](../../../difftest/bytes/buffer-empty/check.go)

Ten independent fixtures matched Go and Lean. Equality-only, capacity-based and unsigned comparisons failed two, three and three cases respectively. The fixtures cover nonnil empty storage with spare capacity, offsets around the stored end, and signed machine extremes. [Reviewer control](../../../difftest/bytes/buffer-empty/review/main.go)
