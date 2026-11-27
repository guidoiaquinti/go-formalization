# Independent review of `bytes.Buffer.Grow`

**Verdict: FAITHFUL**

Evidence SHA-256: `86bf16c7896e2b0c9bf4e80b2c36c0f48dd8f4286ed4b016848ed120a2946a20`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5eaf9bfddb0712522c1872ef404ec0d65bc8f780bc5a4a6cb8cbe8e1fdc3b4a9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ac8cc199f15e9a74264478ac00cd46e9f28f9a91cb83b76def5f1cf976d0ac9f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `393403c7a72bfa9aac2f7d1ec6f287117b2dba5be26977021675b1d3ea311287`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `49c430b0bde866462fffff1048cc9317de7a2e5f82c8f2fe8ae48e81edf61034`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `592b06fa9ef4660043c4113c069c64ca0eaf5ae15c8c99b3aaada96f39c2e2d5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `92f72ad42dba6b244d3ae0f78d697b66fa0417206a02032efc0790fee3b1cea2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9bb23f1220774f833cf3fd1980e92346a9aa577e3596eec0a9ca48d4cfa84d75`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `05bd6066859a136ede9ef5797fede348b68c518b9dab74a9bc6735c431b820c9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `73047a5ad01086638d06edc6f7b76fafd70eb55e5ad65121dfd259979ba34e10`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2e5b8ec7e7e05b93f805afe9685ebbc420a7a780c01813b9fbba8fd0e22f8962`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `e6b7f832f08bc248aa402e64c271d089dd8bf1b095508237480fb6dee39bca5b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0b75f8e29ef6d63c1c931f3697bce0d43654f15f0df08ca525d8dddd2636436c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `03c276fd55cb94dcd69041196d325a6b02df01b83fb66b59958d7b335e8ca838`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9885c2eeb3a6459b138e33fa703a0a98b1031bf0c8ccaf7fd4b2bd7964b78f62`.

I traced Go 1.27.0 `src/bytes/buffer.go` through `Grow`, `grow`, `tryGrowByReslice`, and `growSlice`. The negative-count check runs before receiver access. A fully consumed receiver resets before the first reslice attempt, including for count zero and for requests that later panic. The reslice and the final shortening preserve the original bytes. A move occurs only when unread length plus requested space is at most half capacity. Otherwise the explicit `ErrTooLarge` guard is equivalent to `2 * capacity + count > maxInt` over legal machine words. The successful allocation path reserves the larger of old stored length plus the count and twice the old capacity after the offset, then uses the pinned small-object size classes. Grow preserves the read flag unless the exhausted receiver resets. The native bridge tests the exact overflow sentinel and separate nil-receiver and negative-count panics, along with the state left after panic.

The author comparison passed 21,715 boundary and 60,000 seeded rows using 17, 2718 and 65537. Independent controls passed 11 public-native darwin/arm64 and 12 Lean fixtures. The fixtures exercised the two count/receiver panic priorities, spare fit, half-capacity threshold, length preservation, exhausted zero-count reset, and both post-reset and unchanged overflow state. A synthetic small maximum made the overflow equality and its successor observable without dangerous large allocations. Scratch mutations checked the exhausted reset, inclusive compaction threshold, retained read flag, and strict overflow guard.

The claim covers legal buffer geometry, small successful allocation requests through 512 bytes, and the explicit overflow guard. Resource-dependent allocation failure, larger successful allocation geometry, and native execution on other platforms remain outside this review.
