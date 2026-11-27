# Independent review of `bytes.Buffer.WriteTo`

**Verdict: FAITHFUL**

Evidence SHA-256: `03473c44d94334bd17d19c6cdab54ee210997820440b5f5182ace4bb93f458c1`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `07499debb1ae61637cedd0b51c127503ef04bbe6d6c95b573ab82792dc2f9275`.

Go 1.27.0 clears the last-read flag before calling the writer. An empty view bypasses the writer and resets. A nonempty view makes one call with the unread slice and its original remaining capacity. An excessive count panics before cursor advancement, even when the writer also supplied an error. Otherwise the routine advances by the reported count, returns a supplied error by identity before checking for a short write, and resets only on success. [Model](../../../lean/GoLean/Bytes/BufferWriteTo.lean)

Pinned Go and Lean matched all 28,053 boundaries and 60,000 samples at seeds 17, 2718 and 65537. Nine independent Go and Lean fixtures and three public callback sequences passed. The public control checked the offered capacity, invalidated read flag before the write, writer error identity on a full write, and the exact io.ErrShortWrite sentinel and remaining bytes after partial acceptance. Negative reports, writer panics and mutation of the offered bytes remain outside the declared scope. [Comparison](../../../difftest/bytes/buffer-write-to/check.go) · [Reviewer control](../../../difftest/bytes/buffer-write-to/review/main.go)
