# Independent review of `internal/bytealg.Cutover`

**Verdict: FAITHFUL**

Evidence SHA-256: `d037462c5d4599336d761109ed7d381b635ad7792533f462761e2690340e675e`

I compared the pinned amd64 body with the signed 64-bit model. The
addition wraps before division, and division rounds toward zero.
Fourteen independent Lean fixtures checked the signed extrema,
the first wraparound values and both sides of multiples of eight.

The author's unchanged pinned source-body check passed 2,112 boundary
values and 60,000 seeded cases with 17, 2718 and 65537. The extracted
body ran on darwin/arm64 with the same 64-bit integer width. Native
amd64 execution remains open.
