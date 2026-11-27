# Independent review of `path/filepath.splitList`

**Verdict: FAITHFUL**

Evidence SHA-256: `8ce243df5753985ba9c17288761d83f0c009cd90734580c424510264e6dd5917`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `30eb22403fa81c7973c1613b9b9e574cf378ed6a8950e0e0f8e8337e3ad9dd7f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a1a85a278f5ab0f25362cdd074e50c9fd5e7343833f41773ced59b018c2f3498`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0e7121adec46b12f6d4208ee0fc143a63991c5301f749c3070bd08369b3dddc3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `badc81da2e48e6cbf7f08d039ded47be9b5193863ddaa8e5b64f05f6ab188126`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `45ca3c5290f7c6c17501de95a0f10953f0217e81e1452d5f3ecaeabfa2d59432`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `99491448d2cf10e9f47762d19f673643ab896043cccdada65b3b8f03bb79da96`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `55b44ea4873b6bbbc771325a8c78aaf93fe15565f8f6c508da314d10e3fb5740`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ed2164a8d8caa4fb842de3fc1c9ef874d9b6f46eb1e5f68f6790e9d959e188bd`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6f392c78a123118d190e1563fc1732f4fab40ad62416230a7ca9823e8834e381`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4bb7de25ce3200c6535eafc9982d2afe21c0ef8749e5d3ec23bfdb6bcab6be72`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c3519403e1d979052e48de647a6afb5e800d262f44e244646c51e37dee8972a5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `500f343f38c6d1dd76ec6c2852de290df1328672e5d6e6f10e86f3d7359f7bfa`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9f17ef239133bc2e1878d9e56dc6a3f091af7965020d6846461545d7ee297fe6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `13ada6acb045bc0466c56ced2b743498133a235c4debb8335831e45379bb461a`.

I compared the model with Go 1.27.0's private `splitList` bodies. Unix splits on colon and Plan 9 on NUL. Windows toggles quote state at every double quote, splits on semicolons outside quoted spans, and removes the quotes. An unmatched quote keeps the final span together, and a backslash is retained as data. The Lean fold emits the same fields directly and preserves malformed UTF-8 bytes. Empty input takes the dedicated zero-field branch, and the native overlay checks that the result is a nonnil empty slice.

The author comparison passed 524 boundary cases and 60,000 samples using seeds 17, 2718, and 65537. Seven additional literal fixtures passed the private native overlay on darwin/arm64 and all three Lean profiles. The fixtures distinguish adjacent quote toggles, unmatched final quotes, trailing separators, part boundaries at NUL and colon, and malformed UTF-8 from a literal replacement character.

Native Windows and Plan 9 execution and backing allocation remain outside the review.
