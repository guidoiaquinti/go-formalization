# Independent review of `strings.byteStringReplacer.WriteString`

**Verdict: FAITHFUL**

Evidence SHA-256: `0448bded30874aea03b26bc11d167ff3f5960d6cd622c114f14fd98bfd122f79`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8524e65f811a11e983e0129eb2a2fb3a9918cb33a89e96ce7aba1553620ea65c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b9fecfcdd74fe02839e05f157c452a9b0cb071c56f550345f1610053f4200480`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0fad934667c3e7b739389e695e8b958dd38a7c93327742847a0b09aef26f4596`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dbc6d9a5c99151f4a40fa4d7fd0b161a3320907cc60fbbf4c6d59b0a984ea34e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4ba0490e33eb763c351625d88390c3e8c03ff5f5706bda983b535c14d201e256`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `72e38cf7ceabab523aeae198b7760a79b774c6d2a175e34f068674048dcf395d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1ba1ebba0c9d332908a24506cb7cf3f8a7c8cb5feaa6f9e3334e13634b414a7f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d46436635aa29bf11fe7b2f6baedc028b67f67d8b9ae222a5a92466224e4b514`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3740bc0d1eb7d1883fb884f5f9d3b33301eac2d2cdb1ea8e6e700e2314b6780d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f6641d772d1b50b98ef1e623afb926a7ae50360a6ef5a060936cc391ff4c1fc6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b99cc2ef91345225a227ac204070279524e2d96b83960767418cfda5b6e45a3f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `371198d043a1d9b27c4f0478ee4f45e3abbf30a70694d40c4514598e3c943a02`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a72f15534e2584e4ff19e9562ea3d34974045da9d2feff8b1501ffc58dcbf201`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c1c92d3bb830167e1d04a7c0c70989c33d4bc1a70ec15761ab9b0a08ad960b0b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ea887c6af678960ab8dd6c70c265a9da1f6e08cab4b741f90c4758200330074f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7370761bf6e988e51c6812e7e44d5b81132516ad4e4157313f14709cf104a1ee`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `13dafc0f5af4526c7ac0bd541758225a7c7f7029ea8933aa433e40d89c83cd98`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1ce210c07b784494853d5a7d33170de5780c074c19ca9e83cc4652d4897b7d82`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5fa2746f714bce3854d60dcac22fc41b2b74235b5aa8f5870abf97c3c37e9ad3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7a00f07dcd397fa55a5636558b28b7fcb222dd0198fe24b7e4f808ff74dd9d17`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5c644ac29cff0b851422c3d7dc18fa62994598f3cde6da2cfbfc9b0afb04ff8b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5534e18c9350eea194e7f5e0a4fe999220991169af8a5486f57d52d07e8ec816`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2fdf1b8f4f60f117bf654d141528072691795a4d9f5a9ba5039390a18ef96e16`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `56fb6801845de53987b2008a96104d9b2e52569e50203c526a70add7d0bc078c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2dc2869038885c31e8d402e94700d7fc4ea4f065406cbe460f8d24083b0fa372`.

The pinned method uses nil table entries to identify unchanged bytes. Present replacement entries go directly through the underlying byte writer, including empty slices and replacement bytes identical to the source byte. Nonempty unchanged spans use the selected string writer. The model adds each returned count, continues after short nil-error replies, and stops at the first error. The scope uses deterministic bounded replies and includes the nil-receiver distinction between empty and nonempty input. The native checker separately verifies nonnil empty replacement slices. [Pinned source](https://github.com/golang/go/blob/go1.27.0/src/strings/replace.go#L550-L578) · [Lean model](../../../lean/GoLean/Strings/ByteStringReplacerWriteString.lean)

The pinned Go overlay and compiled Lean model agreed with the independent span and reply oracle on 2,416 boundary histories and 60,000 samples using seeds 17, 2718, and 65537. The comparison checks exact method and payload order, accumulated count, and returned error identity. [Differential comparison](../../../difftest/strings/byte-string-replacer-write-string/check.go)

All eight reviewer fixtures passed. Treating a present empty entry as absent failed three, including an empty write that receives an error. Treating an identity entry as absent failed two because the method trace changed. The fixtures also check short successful writes and suppression of later writes after an error. [Reviewer control](../../../difftest/strings/byte-string-replacer-write-string/review/main.go)

```sh
make verify TARGET=strings
GOTOOLCHAIN=go1.27.0 GOWORK=off CGO_ENABLED=0 go run difftest/strings/byte-string-replacer-write-string/review/main.go
```
