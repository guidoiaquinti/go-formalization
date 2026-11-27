# Independent review of `archive/zip.Writer.SetComment`

**Verdict: FAITHFUL**

The source checks string length in bytes against 65,535 before reading the receiver. A 65,535-byte comment is stored unchanged, with malformed bytes allowed. At 65,536 bytes the native error is returned even for a nil receiver. The theorem proves the oversize branch. A literal valid UTF-8 case contains 65,535 runes but is rejected because it has 65,536 bytes.

The author scope covers metadata configuration for nil or synthetic writer states, raw comment bytes, signed counters and failures. Archive output, I/O and concurrency are outside scope. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 360 fixed and 60,000 sampled cases at seeds 17, 2718 and 65537 on darwin/arm64. Seven independent literal controls and a native state trace also passed. The trace checks fresh oversize errors, unaffected prior state and the guard on a repeated offset update.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/archive/zip/writer-config/review
```

Evidence SHA-256: `f21db7ab28b65eba8abdca34f245e6a2ec4a3e7b3426f7b9a1d7c03cc6c19b0d`
