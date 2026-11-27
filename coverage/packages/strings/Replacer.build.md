# Independent review of `strings.Replacer.build`

**Verdict: FAITHFUL**

Evidence SHA-256: `31202e4b456aedbeda2d0b02987741529ef912bf58d8a88d80b75ada92cc27b5`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `608a20fc5ebeb96d2e937ad891c183f1836f8339c1f68169911d13e80a78216b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bb75128f4c673849e2bb7b63153430b8a6cf652084dad0eeac7bf9fd749bdad0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c45bca56151245e7b37fa4b46ddc8ffb0438dbdd3c5e1fb6db5591dd8d70f988`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `210985eb87fcbbf8286ce8ebeb34be5e090cedd7875ea629e53df448d9295ebb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5cebe6d78c30d0ae467eb29a94d65eb7b89e776aa13cf584bdd1b943e777ce6e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cf2af3a918621f777474bb5fde41353bb1518dabe2448c71ef1442de326712c8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `447a6699bd591d98120e34be0ff7fcecaeac99b9b16901eaba0a6a24b1038a47`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6fe83a1438af86666cf3e6ef86b61c6745033d6f18e8f6ea0ad561f25d138fb0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `175c78006e886e1b189a4efcf351bdd54d322a9dfd34e913303117607560fd55`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `bf0c6c0acb85d430c6d320caea4f513d800e2c7105febf677360f5f3f81a9ab4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `640e4b0c52738c3a55914c72359ec69268c5491522354018cc04e927fd62d7ec`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `47a544e24a89be5314747b53f0f8190d7c3e726b408989980009f74d9f532a39`.

Reviewed 49c15b9 against Go 1.27.0. A single key longer than one byte selects the single-string constructor before the general length checks. Every other configuration containing a key whose length differs from one selects the generic constructor. Otherwise, all one-byte values use the byte table and any other value length selects the byte-string implementation. An empty input therefore builds an identity byte table. Reverse insertion orders byte-string keys by decreasing last occurrence and leaves the value from the earliest pair. Present empty values remain distinguishable from absent keys.

The independent control checks eight complete dispatch snapshots, including single UTF-8 and malformed keys, a generic empty key, duplicate bytes, and the relative order of byte-string keys after duplicate insertion. The pinned Go overlay also compares the complete delegated structures against direct constructor calls. Always choosing generic for a long key, ordering byte-string keys by first appearance, and taking the final duplicate value each disagree with fixed witnesses.

The differential passes 274 boundaries and 60,000 samples under seeds 17, 2718 and 65537. The claim covers successful build from even pending argument lists. It excludes malformed private state, nil receivers, later lazy initialization and backing identity.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/strings/replacer-build/review/main.go
```
