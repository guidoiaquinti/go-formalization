# Independent review of `time.Time.Add`

**Verdict: FAITHFUL**

The source splits the signed duration with truncation toward zero, performs one nanosecond carry or borrow, and advances wall seconds through addSec. A retained monotonic word advances by the original duration. Overflow of that word strips the marker and stores the advanced wall seconds. The literal controls check both packed boundaries and the native trace checks raw words and a preserved location pointer after monotonic overflow.

The author scope follows raw signed and packed Go fields, including synthetic words, and represents locations by opaque tags. Pinned native Go checks location identity where applicable. The scoped Lean axiom audit found only permitted dependencies.

Pinned Go 1.27.0 and Lean passed 5,184 fixed and 60,000 sampled rows at seeds 17, 2718 and 65537 on darwin/arm64. Twelve independent literal fixtures and a separate native boundary transcript passed. The controls distinguish unconditional wall subtraction, unchecked arithmetic on the monotonic word, and skipping the wall round trip.

```sh
GOTOOLCHAIN=go1.27.0 go run ./difftest/time/clock-arithmetic/review
```

Evidence SHA-256: `ca3af9b61a1439d1e6b9128a3147140a5cb1759e60a0584aec0c74a8577652fd`
