# Independent review of `archive/zip.Writer.SetOffset`

**Verdict: FAITHFUL**

The source permits any signed int64 offset when the count is zero, including negative values and the minimum signed word. The exact zero-count update is proved. Any prior nonzero count raises the pinned panic before writing, including a previously configured negative offset on a writer that has never emitted any bytes. The native trace checks count writer identity, count and unchanged comment after panic.

The author scope covers metadata configuration for nil or synthetic writer states, raw comment bytes, signed counters and failures. Archive output, I/O and concurrency are outside scope. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 360 fixed and 60,000 sampled cases at seeds 17, 2718 and 65537 on darwin/arm64. Seven independent literal controls and a native state trace also passed. The trace checks fresh oversize errors, unaffected prior state and the guard on a repeated offset update.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/archive/zip/writer-config/review
```

Evidence SHA-256: `80af5c29162e5863e2184a68cd46591ad914fb518cbf6f4c48643cf839a1f98f`
