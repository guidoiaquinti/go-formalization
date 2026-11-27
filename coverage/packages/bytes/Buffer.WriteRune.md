# Independent review of `bytes.Buffer.WriteRune`

**Verdict: FAITHFUL**

Evidence SHA-256: `1f6ba0a5616b2fa8f2394bf8016567cebb0fe11552e1ebba463c23b50dc5f688`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a13b01bf9e78aae9b1431f961a4257e50f732f09829ef4fda1edf98fa5562d69`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e2bb966055afd3c21fa23d2eb179570a4dca4dafee3e254cfb61bb8c0b706278`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ef66a5f8e9a8abf3efe7eb7522ee50ec3d0145121f753f558f26d96eadd94446`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `af2803d52ffa1dc2adc97f34a7b420dfea92e61c2c99371e3f412e22b82f4ef2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `54dff9d3523ab1edac485f9686113aa7707b617f25717972522a12425fca73f0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bd0aed507e3b2e783438ed7bdbc098a0db1813901497b78a240cf25dd8fde8fd`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f5489c017f6b351ef5d5f7a05a6a51665b87aa7c2634463a2adb183f4e8b074c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bbca8e4ec8d54295fa3a74bc9c41182c6d15b568e1e30500a48b5fd37895c85e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `df7d7bbcd3b86a325987c55c52d424376d90e395af5713ac90b99a32da2c1814`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fc776f10f77a5cd132265014fd864e278d10202a770b41d54ef3ebe927905c20`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ec337d0893e2848e9b54997fa6b7cd735a51c63bcdfd291f2029dc27aa62bf92`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `997d68b2ba17d06eeac7ed36b11de14f2bb7440b06385d077e979f5924bc86a8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `20d348967344f27a2c3333aeaab40999ea07d0e801adbaf89b12c89d6a5fdb2f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ede956e3684227b91a341df56dc568e2d67d21c18bd43d381f28bc9223adabff`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `17619a8a9d22a71aead17ca8513c2f00c85c41950a2c7a295fbe8f412b8cfc88`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ba4e1adcf946dd2a5b735ed0c6aa3179a7b294435a8b158019767a4d3841ef53`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `90ef444d2b7ea382d9e801a3ef2bcb3258cfe354cf5009663bbff7d712573fd6`.

The pinned Go 1.27.0 method uses one-byte WriteByte growth for ASCII and reserves four bytes before encoding all other signed int32 runes. The model matches invalid-scalar substitution, reported encoded width, flag reset, exact capacity and spare bytes, and alias reuse. Independent fixtures cover both nil-receiver branches, nil and nonnil empty storage, invalid rune extremes, two spare bytes versus four, exhausted reset with and without allocation, and both sides of the compaction threshold. The same-package overlay and model matched all 17,179 boundary rows and 60,000 samples from seeds 17, 2718 and 65537. [Model](../../../lean/GoLean/Bytes/BufferWriteRune.lean) · [Comparison](../../../difftest/bytes/buffer-write-rune/check.go)

Twelve independent fixtures matched Go and Lean. Two live public alias checks distinguished compaction from reallocation and verified invalidation of the previous read. Incorrect variants that reserve the encoded width, report the reserved width or encode a negative rune as ASCII each failed a fixture. The old WriteByte helper changed only in visibility, and its ten fixed fixtures and two alias checks also replayed successfully. Larger growth requests, allocation failure and invalid state geometry remain outside scope. [Reviewer control](../../../difftest/bytes/buffer-write-rune/review/main.go)
