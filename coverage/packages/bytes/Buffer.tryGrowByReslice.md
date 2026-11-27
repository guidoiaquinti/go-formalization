# Independent review of `bytes.Buffer.tryGrowByReslice`

**Verdict: FAITHFUL**

Evidence SHA-256: `402f081b2f50d15986b653e2287fe33974113174e890f6cab78539a99ff89e4e`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cdf13251bf71bcfa156dd96478539dfb68542142f24f81389b83d2b8653b3625`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `852cdd318b165bcc033cc4f4710cdec38701dafdbf4d558018900a910f6d2a5d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9178e75b07818695dbacdb64f1958b1f86e1aa54a553cc3b6c9a83f3e8bd0fb9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `fb6d4a84417d8e4ea9a36685ddf46c6084dd4bfaeb884db417512a213b65b9fd`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `764ae5591c7370f08f8d1597bb4ec996bccc559a8414b9d9cbf697f11ccf937c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8e989646982f13662d1e4aae0eab6d9d8c999c6bece2513fee0be1b3f4795b98`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c58319c7062c1fa8a4d247f57093ff6152156788762c76dcabb4fa8f66c7d611`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ea7dd59484671cd274f56d0c09c5248f6bc800a949733f028b1efe71c8fda6a4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0e2cf1e46d96131d6910adc1dacc81983d5f3ba7874b6fe33ca74b20348ed256`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4aa2c09e17bbf08a594f1aab9d9745d5b01ba23d8329b3766ab2082d4d2aa3d3`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1208a147ad79433a9ff0f0adeb0785d9dc1963d208cdd1498b7c1d6fc5aa1727`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f9c3ee9b1cf3d74c0690c5639aa806aed827541783aec46284c859f1f924047f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `f73b85864216bca81e8a6d0b42a72d39536c6f85eea8ab22db17c04fef8976fa`.

Go 1.27.0 applies the capacity test before forming the new length. A nonnegative accepted count cannot overflow because length plus count stays within capacity. A negative count plus a nonnegative length also stays within signed int bounds. The model preserves storage, capacity, nilness, cursor and read flag on both return branches and on a bounds panic. A successful shrink can retain a cursor beyond the new stored length. The earlier length is returned for every successful reslice. [Model](../../../lean/GoLean/Bytes/BufferTryGrowByReslice.lean)

Pinned Go and Lean matched all 40,400 boundary inputs and 60,000 samples from seeds 17, 2718 and 65537. Twelve independent fixtures matched both paths, covering nil receiver and storage, zero and negative counts, minimum and maximum machine integers, exact capacity, and shrinking below the cursor. Strict capacity, immediate rejection of negative counts and cursor clamping each disagree with a fixture. The checker also confirms the original backing array is retained. Results assume legal starting geometry and exclude concurrency. [Comparison](../../../difftest/bytes/buffer-try-grow-by-reslice/check.go) · [Reviewer control](../../../difftest/bytes/buffer-try-grow-by-reslice/review/main.go)
