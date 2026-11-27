# Independent review of `bytes.Buffer.Reset`

**Verdict: FAITHFUL**

Evidence SHA-256: `a30dab8b452f0d849e9081385a30e0e1344699acd047748e6830bab2ad2649e7`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ab50b4dd95cbab4a816a270fb5e2d3deba522dd11acece0797e47e2c1fc5a803`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `58fd88e5532fc1c462af7f6424e35c87369a09a56bc9440a690a054e54bdd589`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `222265a4bc4e518578469f6bb36630c3ba0f2ad0a2b0ee922d1af5e5a52deb81`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `18f59a7bae1e810341650921445d8281944f67a008177e46251cd8906cdd0c36`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `78dad5a417636d9b2fa6f26297e3a340447c976abbb7cc86b27c1e20ae3607be`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5ed0d4eae7d89f0c4ec8ecba0041662442181d28b233f8caeae97ca555375e49`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `149c494881ca720e659c655ae0c253808df2a99242a752e2ffb8567e2b92868a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d6b6322babc9267d5a1666d3a66a74dd71793117f96e7f2c4027b5d747c2ad9c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `43606b5d8394b35634d7913d4ac189fa681cf350f7aa2f6b06a9adfdda77f134`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ec023b2517b9604d38a2b85d54a477ebf8f6b798b92cdfc1c17dec991b4f7aa9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `67e9f0f5b6a5efd0a59c7758dc35e4f88194a2b05599cbe619acaeb9c5537821`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2e162fb0dab7e00e0afd18879fa51e71b9da7ff95191392dbd3b45ec00ee4b60`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b8739b4cc021e873baa7f77e147b3e29f8e99a888cf8857f9aeb75f1e4381817`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dba5eb30496cf55fe50f5e6facd051df5afd28b1f0e14f118292732c6fce800b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `398673d86e7fce26537b0972dbc8d21dac455769cfc9a663cee42d5b366625ca`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `54eca51395cb1979cb59e84a9db23b493a9c1d8c94b3636629659db369622878`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `525f18c6ab37b68db6d551f237c6a0d344117661305ef96d77fc2c55b51bd9a0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `577e61507b3eb6a571a3f292720d3c7df9fae79c6edbd810a4c893270e4c0a76`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dc6c9a332afa283145be5c1f51f9c1127d71203c428d789374581570f6877f6b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `52117f728a2b002b9ef92c396678a81c7361b21f0af3a57e70a851b7734610d9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4e85573c455a411d7ec95b8402ebb9858c18d79414d7bd68bc8d92114fd933d9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d23da41003a2741407cdcc49b2d9025c07bf2534ad87d788ba1be848d979951e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b818927daa0ea2940b9aad4b14cceb111c2e9b935c43c7a3e270393f89ec675b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d59800ff807fbec989a13d9f80b37a4cd4ccc752a209b94395af54991917afba`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3bea6b9ce1b28784da1613177191811898d1c7b5b3e18905f9960c89392e5bc2`.

The pinned method reslices the stored bytes to length zero, clears the read offset and invalidates the previous read. Lean retains the original stored and spare bytes, preserves slice nilness and represents the nil-pointer panic. The native check obtains the new empty view after Reset and shows that a later fitting write reuses the original backing at offset zero.

The direct comparison passed 2,609 boundary rows and 60,000 sampled rows from seeds 17, 2718 and 65537. The independent control passed 12 fixed Go and Lean fixtures. It separately confirmed which unread calls succeed before Reset after a regular read or a rune read, then confirmed both failures after Reset. The rune histories include widths one through four and a malformed byte. Preserving the offset failed seven fixtures, leaving the byte-unread flag active failed seven, leaving the rune-unread flag active failed five, and discarding storage failed nine.

```sh
GOTOOLCHAIN=go1.27.0 GOWORK=off CGO_ENABLED=0 go run ./difftest/bytes/buffer-reset/review
```

[Model](../../../lean/GoLean/Bytes/BufferReset.lean) · [Comparison](../../../difftest/bytes/buffer-reset/check.go) · [Independent control](../../../difftest/bytes/buffer-reset/review/main.go)
