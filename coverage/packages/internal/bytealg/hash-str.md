# Independent review of `internal/bytealg.hashStr`

**Verdict: FAITHFUL**

Evidence SHA-256: `65723cf94d203ab9f3e8628dc0ef36be7f4fbcf9402780d67b481304075799d8`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `13cb47081a25699531f8a79f87d56440e498a51aee12b2e59946cc47bfbf2139`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b0a71903fd0dc1ffeebc8f6facc6f1dc4ba0e2d1705ef8543d70998965c9843b`.

Reviewed 1816c15 against Go 1.27.0. The forward hash multiplies by the pinned prime then adds each raw byte, wrapping at 32 bits. The square-and-multiply loop gives the same wrapped power.

The shared independent control checks ten pinned native and Lean cases, computes hash coefficients at arbitrary precision, and reduces the results modulo 2^32. It exercises differing first and last hits, overlapping matches, an empty separator, longer-separator bounds panic, malformed raw bytes, and both hash-collision directions. Rolling past a collision to a genuine match distinguishes byte confirmation from accepting hash equality alone.

The complete independent replay passed 597 boundaries and 60,000 samples using seeds 17, 2718 and 65537, through both native generic instantiations. The claim covers exact arithmetic and raw byte outcomes, including the bounds panic class. Timing and exact panic wording are outside scope.

Run the shared reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/internal/bytealg/rabin-karp/review/main.go
```
