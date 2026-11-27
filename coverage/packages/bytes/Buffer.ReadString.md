# Independent review of `bytes.Buffer.ReadString`

**Verdict: FAITHFUL**

Evidence SHA-256: `1134911e92ec178ea4e4a1ef0f5ff97b37e496c09f30f3d5baaed7f87e551886`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `93c176dcf6608278e0dba681e0de35189dc2092ac79c0244895004b319b4de30`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2fb7a435f80b67f271788d26d3defcaa355cdaa1577b477459da65eb8fefca53`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8bbba85d6a3e0138ac4c51a3ab53d2f918bbe2f12f80098b46b393921bfc407f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `43faecc5cf4b2de73f33492d349a80a3feed371051b263f54aeac0ee40615b8d`.

The pinned Go 1.27.0 method converts the delimiter helper's byte slice to a string and propagates its error. A hit includes the first delimiter after the cursor. A miss returns the full remaining string and io.EOF. The helper installs the ordinary read flag even for an empty EOF, without resetting the buffer. Conversion preserves invalid UTF-8 and embedded NUL as raw bytes and snapshots the source. [Model](../../../lean/GoLean/Bytes/BufferReadString.lean)

Pinned Go and Lean matched 11,525 boundaries and 60,000 samples at seeds 17, 2718 and 65537. Seven independent fixtures also matched, including a delimiter before the cursor, one only in spare storage, and an exact last-byte hit. The public control preserved a returned malformed string after mutating and reusing the original backing, and confirmed an empty EOF reinstates the ordinary read flag. Excluding the delimiter, preserving the previous flag on EOF, and suppressing EOF on a nonempty miss each disagree with a fixture. [Comparison](../../../difftest/bytes/buffer-read-string/check.go) · [Reviewer control](../../../difftest/bytes/buffer-read-string/review/main.go)
