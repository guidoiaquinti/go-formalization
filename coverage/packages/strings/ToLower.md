# Independent review of `strings.ToLower`

**Verdict: FAITHFUL**

Evidence SHA-256: `58e83f8fd5f1963253f8e018198b0410e2b183b2f5b1508948eeede8db8327b6`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `38a0fcfd18b418357cb62eba4789be612c4e6d53b614f212225587910e9e0dc8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `02f8dadd0032a5b9fd2442c7d1689b4111a794e9d457a0d8110f82e3983fa060`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `05a8794968ac15f69404200cd1d77f65ad463c9d086bfa30a24c209b07a8b8e8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7ba52ae806d695d387dc4dec78f9c32a639c1c5a1acbaa134bc4c2ec09a2ce1e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b106d7b6666e92c8e6c5942f9d0bcf400af074f9d88a1716023a9d583cade50d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `402b8d9a7e17ec63fa9153adf1d877835c752b1d08b7a3ef4c56cef8f03ca570`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5ae8a9bef86a34a56ff624946f8c4a7c22f809a85703a28fa2367eb6f262b001`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d379b92d33e2cf88547d707fce38d501c0f78708f954b674ba7a05bb57bc8338`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `34f98609b036441c893e8e53514d56e870cd01e563143bcfd0f8eb812184fe11`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `60def9bc5191f07dd3b571bfd2ac490525f1c2394a8e2391207bf27c326aaadc`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `300470f1e17cd85ea1aebd0ce07cc887d29a4d7d2153bcad0acd405b14e09055`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `67fa182b9b6afda1b822e3dd7d2185d9c3a848e2d6a63fedfc51aad0aabd1fba`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `440e926194ad8faed7f6b15fa0a1c1d74ee3b8517db5f8d29545f53d22df9359`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `35a35ae626e73eb331c2888b1ef34a2704933e4be5bc8140d7d0b83db1a99a26`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `eaa9b7201bd208b42b26e6dc943dba1ec2c4da0726f182e0a84341a9c2a18b8b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7c06a28c722e32ac0d70ca4fe76b671ad6995e79d2c198c8710f1285fd176d33`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c182e8afd3961a58da8c58d315cd7dd1db861d1342b6533d1c48b6fea75a9277`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1e57a392841ca9bdfec0601bc2326d2d2db003503f4b1b1eee92a024b8f5c585`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ac6364de818746e4cb99c6fcaf8ea2bbc021f55400d9b7b58571e28c83d6177a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `765c2d65d3494e1009ead5c2e19b068d3ddbd3c4efab4c6c5110a355979319c0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b7185a325c58d76ba400156514636a3fbaa1f81bf60d1cc26e1324f95c4f7124`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bed928176417d0565fb4aeb3edb6e996ffca5adffd6236840fe4b89f20556fd5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8473c79b58292f736d8fb8400396f4c6a7f08281479e409e7c18be93d5fe2cd9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `313119afc74724434fd9ef6e321b54d673ec79e64f0b7e2d272a73fe1c9d750c`.

Pinned Go 1.27.0 on darwin/arm64 has an ASCII path and a general path through Map and unicode.ToLower. The Lean ASCII path changes only A through Z. Its general path decodes one rune at a time, emits replacement runes for malformed bytes, and applies the pinned simple lowercase mapping without contextual expansion. String allocation and backing identity are outside scope.

The original comparison passed 3,504 boundary cases and 60,000 samples with seeds 17, 2718 and 65537, along with all 1,114,112 Unicode table lookups. Regenerating the checked-in Unicode 17.0.0 table with pinned Go produced identical bytes. The independent control in difftest/strings/to-lower/review passed ten Go and Lean fixtures, including width changes, an astral mapping, Greek sigma, literal replacement runes and truncated UTF-8. An ASCII-only mapper failed five fixtures and preserving malformed bytes failed two.
