# Independent review of `bytes.Buffer.String`

**Verdict: FAITHFUL**

Evidence SHA-256: `287fee937144e7a707897701b33fe8a13075e4efbfb1f7a3246c6b66a7c440e7`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `638e30904962bdd5c723f5f523077caa57e6a405aad8f74837df15575497297e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a6a3f0053a39a916f984cc1f183f631b3bd1e358c305f007d0f0bccb25213bf0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `75933bf0ca6ab36f23fa0f4799d983e6527f3641c11f67c2c1b88db6d0dfd0a4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4db75957ff453cf72b0ae37f1ad31a5348b7ff217936b1c6fe5d0f0ab75cdbab`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `16962685a1e7a1213ff6adab40b321bfb5974527cf46562c223be1180faa511a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bcb22187d50344ea6021d6c2d4096c79ae9820c289477a36a320ad31e2744fd9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cd4b6443f8da93d998b0add3376594fd1107a7aa6defa96966d7ba436eeefe91`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b7fc2c0b38878dd8b7c70c6a2e53daaf499b144f48b898bbccf5a28e054c2631`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c55dfdc92502afd00ab135ca99d9ceb023a9a94340f251ae6286de54f33480a3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4faed99f0ad98ca4ac9c824c5c0428170e562becf6e79c8575e25a3f8ba48063`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ab04fd30a9aaa0327f06ef0ef0fcf0c44b9c26181c0b85dadb2ada6791151554`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f57ac62b843809e6caef82859a30134487fd8f2b0a449d0d2e864b6c58ec4d13`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cb8c251523da598bedb50acc940a3d978676e30bac6fca70a01f72ef89a24e58`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f7b2a6b66e96a0f6e56c2643c573b592f1f0083183a94c9eb15dbeb8306c44b1`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4843a02c673cc9bdc1b41303dddf3cb152c4bf6bd1b0715acffcfab38675d13b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `348f455de2d285c6fe51958c98f9b4ee1e89a48aec9ee05712daa936eb30f9e9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7d6226f93b5c23f9fb334e690e56c72427cf2d31fc0050a90abcbe16405b4a17`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f591f725cbbc302b40a8d72e2908f69971c6fa46e1df76de2dff74fdc34122bf`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2426d5b8779987bdf282da669d589de3f60318665f9fb6e4afc5c84c6fc05636`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `92ad4f9857678664ffdfa2d56d9d1367b1ca3644ef512503a63d2ff0b431997f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `94a763c179b221a4b4b7bfbdfc25521be8c468814875e8d456e0c34fa1e7e0e2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `18fcf8ada67d82781ca79d197ddd904ca1c44193da4eead20b30a4e5ffcea54c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `815422e3967f57164d9c2ad2c6e84be46862cd1e96097c3dec0cd4b014bc1e49`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `90354219af247d81bf85b43fdf1b9f306d16b8e26a8f5af83e244a78bcc7028f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c626dcbf2df02bb140cf875f9e709a111b4db3c5e9295c58202ca6d1e0d881a9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4bc229fb5af914fb957149804283f56f91c92c727d69e1498931f41618a181fa`.

Pinned Go 1.27.0 on darwin/arm64 returns the literal nil marker when the receiver is nil. For a nonnil receiver it converts the backing slice from the read offset to a string. The model specifies the raw unread bytes and proves the before and after snapshots around an overwrite. The comparison uses actual strings retained across writes through the unread view, reset and overwrite.

The author comparison passed 2,604 boundaries and 60,000 samples from seeds 17, 2718 and 65537. The independent suite passed six fixed Go and Lean fixtures with the control in difftest/bytes/buffer-string/review. It checked nil and zero receivers, partial and complete consumption, spare bytes, embedded NUL, and malformed bytes. A delayed read from aliased backing storage failed two fixtures and returning an empty string for the nil receiver failed one.
