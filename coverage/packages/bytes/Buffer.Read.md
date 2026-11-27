# Independent review of `bytes.Buffer.Read`

**Verdict: FAITHFUL**

Evidence SHA-256: `a2a96570c46b3436793ad8346905998e19e32b70c105c3dcaa1a6c77a4ab53b3`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6a5e1a38e486d694d0ecb65eaee2a336ef123d62e91d1b3054a60c2681f56c3a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4e53a598c1244c3153475613e7cf3774ef83d53ff5b84dc805df81b86ab9e0de`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `5746286e0759114c5ea88718683b35f6533db6615a0a3996f658c8ff89fc352b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `82364bad6930277244937784054ea76ed4628d5c1562f794dee3a06ac00a04fc`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `168521edac55b8c4100684c8b4af418cab47693c23b26bd8b410af57ef46f24c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `55911c3f148d883005366e3bef4c8980030c83bb04d859e5a32bc38f8aa8a26d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `069bdf753beb939fee1b403e7295237d78b86022654f20c349d127885f8c1988`.

The Go 1.27.0 method invalidates the previous read before inspecting storage. An already exhausted buffer resets even for a zero-length destination and returns io.EOF only when that destination is nonempty. A nonexhausted buffer copies the unread prefix, advances by exactly the count, and reports no EOF on a short final copy. Only a positive copy installs the ordinary read flag. The model matches these transitions and preserves the destination suffix for a disjoint output slice. [Model](../../../lean/GoLean/Bytes/BufferRead.lean)

Pinned-Go overlay and Lean matched all 16,834 boundaries and 60,000 samples using seeds 17, 2718 and 65537. Nine independent fixtures and the public sequence of a short final read, a zero-length read, and a subsequent EOF agreed. The fixtures distinguish skipping the drained reset, preserving a flag across an empty read and returning EOF with a short copy. Overlapping destination storage remains outside this evidence. [Comparison](../../../difftest/bytes/buffer-read/check.go) · [Reviewer control](../../../difftest/bytes/buffer-read/review/main.go)
