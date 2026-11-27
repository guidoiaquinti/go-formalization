# Independent review of `bytes.growSlice`

**Verdict: FAITHFUL**

Evidence SHA-256: `10a2c5c6228e4e559d89a3251eba7adb077689cd3210d42bc0df93dccf91d078`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `372c24e04323e0cdb5e7326007b8b0cb74618a35f42581e8b6f5b2ef39e5d336`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `918c9b7674b8530c2d8b6dbbe4513b3a65760b653afce9967f331e12b79591f9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `42d89f694ee3490b1a6013e75aad2c8e2432b3f17f17246770ee224c16ef2038`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `563b926ea428ce2c80b86383c56c511ebc3727b4104822beb471153767688dfc`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0abf6fdbd30e8a917015ea4c8fa43c73f5001cb9a38472758772db04d76f1bf2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c19183e1274031b4696c8d7dfc44605137c5566ee5219e79d609006ae11d5f47`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `3066cb125766516d27a2ed955def14bf1be32fc5f9f3232fb0d7e43ad0f7ffe0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `04682c627c63f6cca35b67a250c9311f7d4390e1da8bdb8a54ab60d5205f4125`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e1a46365d67cde7983dcbdc6aa72bd8cc209c038e154bc58ea091de362e9009b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a186239a4767db56becafa576cacdb74ab4135e03fcc19e1d2b9c852d2ffe072`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0eb651a2dea21c84531b757cc9130bfd3806e6c59121758d944ee0bc8e201e83`.

Go 1.27.0 computes the signed length-plus-count sum first and replaces it with twice the input capacity when that is greater. Within the bounded capacity domain, doubling cannot wrap and a positive sum that wraps becomes negative, selecting the doubled capacity. A zero allocation request returns nil even for a nonnil empty input. Otherwise Go rounds the allocation to a byte size class, copies only the original payload, and retains zero-filled spare storage in the new backing array. [Model](../../../lean/GoLean/Bytes/GrowSlice.lean)

The direct pinned overlay and Lean matched all 24,906 boundaries and 60,000 samples from seeds 17, 2718 and 65537. Eleven independent fixed observations matched both paths, including native int extremes, an empty input with nonzero capacity, the 8-to-16 size-class transition, the 512-byte bound and original spare bytes set to a nonzero poison value. Variants that collapse empty input, skip rounding, reject signed wrap or copy spare bytes are distinguished. The review assumes legal slice geometry and successful allocation. Allocator panic translation and larger requests remain outside scope. [Comparison](../../../difftest/bytes/grow-slice/check.go) · [Reviewer control](../../../difftest/bytes/grow-slice/review/main.go)
