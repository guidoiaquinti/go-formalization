# Independent review of `math.stirling`

**Verdict: FAITHFUL**

Evidence SHA-256: `477c0cf71435ad23cab1559b0305b06332c1f849c2a578a0f1450a8f5cdc3227`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `79a3895e474bc0a1c583a4a103113cddf489e0fc552af267d2fde786382448d5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `69f94cc70f6f99a72c410a8ced32a4817aa65cf3fbc39fe421f22668b88f97d2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `108e2887f631f6bb64be545eeb4b93ee6139bb9606057f12c8f70ab267077453`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `045a698faf2fb30c3ca247149c3277bbb2f48bf67e1ecd4a114bcc0f65ec0973`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `291875e9991d332bd8b2bc25ae63fdc3ed103966267abb443378cddd173bdb38`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ef80879e92a7db0d1edf1a13acba7fb88653c7219c64d6937f988a5e318ca73e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `262f678b4e83fe82f7214238a8faa883870249b148a66ed7056fa816b4841063`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `15c570989022d795f7d26557810c589f9c6bde1846abef694416acf296bf5ad8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ad8af8dcf1a25156859f6e0064e6aecbe606212a5f643913d86098ce1a985c32`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a5ee7a8357ca74502e8843a1231fa55ce2958920edef5c77c0e0cd0b6ad33dfb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `71e95c5627368a875d0d4f6b8869393c1660c77095f2fd2cc44f12bf6e23ccb8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b96fa27df7ce161d99f7ed92a1a81b9a5e8e26da6018d24af1d977ca96ffceb6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `90da48bd06f97bc71ea56bd9e1c2ff0efe403221d4bb3a53a5fe3211d06e697c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a204ddc0c676d1660300fc3dd05e1e3625a16470637adef30f50203b135abd93`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a55fb40a9c923b853be0a32fc27f1f9904e879a6494cacb6d41da456ffacbe1b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c1dfff8da009df85220a49714163785cc71691382630dc9806f4333e131f4062`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2369062df4c4bb59b3d3415040f3e1ea24e595698b65c00082ca6cd154c8338f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c7414805b0057fb0bfd4378197eeb93420bbbfe54e47846388ed4db2da816166`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d3491e2f6786dee2c8ddf3f252a068f67238d0c4df6470ae1bad3ade28286fe4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `78ccdf8d032b4a8aae48b4e48f18ae6ecca1290d756c3ad83edbe026a46e7943`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7ba1150b1aea6320c21666760254ca1bfa8626008b1060c474304da9e7a2c952`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `52df66d1fa025ffa88e95544f49b9553c1dc4f1d3d473d10e40159d79bd4948a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `73811457a4bca93d0c7cdffc69dfdacbbc2885505bafa3eea2999013f354b23e`.

Pinned Go 1.27.0 on darwin/arm64 returns two factors for the modeled Gamma domain at and above 33. The implementation changes the allocation of the exponential contribution between the two factors strictly above 143.01608 to keep both finite. Above 200 it returns positive infinity and one. The model preserves those stages and the correction polynomial, including the separately represented compiler and runtime fusion choices. Direct inputs below 33 are outside the claim.

The author comparison matched both factors on 24 boundary cases and 60,000 seeded cases (17, 2718, 65537). The independent control in difftest/math/stirling/review passed 13 native fixtures and 65 observations across five Lean profiles. It used logarithmic evaluation of the first factor and factorial products for positive integer fixtures up to 170. It checked the immediate neighbors of both branch cutoffs. Collapsing the factors into a product and one failed 50 profile observations, and moving the infinity fallback down to 170 failed five. Native amd64 remains untested locally.
