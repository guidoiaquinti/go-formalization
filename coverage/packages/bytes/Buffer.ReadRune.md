# Independent review of `bytes.Buffer.ReadRune`

**Verdict: FAITHFUL**

Evidence SHA-256: `36713bf9a17fa61a7bfeb9c8839b428254fe36645c33ae97b3ad5111b585c84a`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `01d513a82c97380931c9fe001c9dff556105a1af5125fcccc3d51c63a74c85b3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5b1948bc83e9e324e1d247ec634772c68a06834cdbeb6469d89669aa4f0a1023`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a57d64dbd694142657cb5ccc2af996f463b1ef5bddd7c77ab25430fc0289a635`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2001b608acd9c4c86ae199d36789392a20a5d8e315af82bd7a72d9815cf17f97`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2c1ae47ae7a62e5a6626ffaf461c4743690c42870a3d1d2e066cb21947cc7035`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `809ebb6b9a7e89a288b8f4e34bf030ed3e62970a6f7c5720a208149820f2a7b9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2aa3dd195192e5ac1763aca28eaf693e02b557c6cb9dd49e8615e809a80c9808`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `93cf6e51504db00b04789d86751fa9088f985b7d39b7ca836383fe1b9030a717`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e2eac792c97c7ba6551c310bb595704379ed9a7c8ae08bd9c53564b5ba7992c5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7797f2887ad30bedfcb40398072e1b8cccbd06e37097bf53be73eddc8c96d5a5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `940e5004ae3ce0e6cfd17e5cb6432bfd476cd2a98b9fc488c7a3a845c904de39`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8d197937f28bdb25d527bca6ea1b2a8827dd3c40f58d08014223ca21558cbf1a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8393823538a4c5aee7c6eeff0c910dfaf7572960b05bb8fb783f36416306a443`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `466595c554641dad65193daa0b9df907cc02aed83bec46ebc273bfae29131402`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9b12a57365baf461878da22603dd517cfeae35ad3cd7e82c735759d620babed9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `56a20e5440141953348c913c9e4bc8b838fd5da9ae2ae8605f12ff2066012ec0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `457cf97610c23bdfef69ec726efc2e7f87c602bf2461a17670ed4dc2063cb40b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ecbab638eb50530e185159961289a8bd834c34d8f52359655c15915fb307fbe8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e049114449f708146f96dd3e2c333734cde78e02f9d099edd1db181708dbe5bb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6f3d01a7bfca258ab664b2c3cf275724b0e9dccdb9f42e426a15cfe2531ca27f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6feb171b77301e63df859ff9d75ab7fce6ad8b961155abcd267906237b33a66a`.

The pinned method advances by the decoded width and records that width as its read operation. A malformed byte returns RuneError at width one, while an encoded RuneError has width three. On exhaustion, the method resets the buffer and returns rune zero, width zero and the exact `io.EOF` sentinel. Nil receivers trigger the runtime pointer panic.

The darwin/arm64 Go and Lean comparison passed 68,668 boundaries and 60,000 seeded cases from seeds 17, 2718 and 65537. The independent control passed 15 native and Lean fixtures, spanning ASCII, two- through four-byte runes, surrogate and out-of-range sequences, truncation, overlong forms, and a prior byte read into a rune. UnreadByte replays the last raw byte, while UnreadRune rewinds the entire previous width. Omitting the EOF reset failed two fixtures, assigning width three to all RuneError outcomes failed five, and conflating the byte and rune rewinds failed four.

```sh
GOTOOLCHAIN=go1.27.0 GOWORK=off CGO_ENABLED=0 go run ./difftest/bytes/buffer-read-rune/review
```

[Model](../../../lean/GoLean/Bytes/BufferReadRune.lean) · [Comparison](../../../difftest/bytes/buffer-read-rune/check.go) · [Independent control](../../../difftest/bytes/buffer-read-rune/review/main.go)
