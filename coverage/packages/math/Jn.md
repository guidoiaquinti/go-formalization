# Independent review of `math.Jn`

**Verdict: FAITHFUL**

Evidence SHA-256: `71eea3f794a3f41f05158aa228dd8837c08895453fcf26a3143cfa32dd38a55b`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `dea659fe6814d42dfafc37c4c77f7dfbe2a7e44602183173ee2f2e25f3b97d83`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b715220a7196c315723fbedb1d1c7618777d4ce05bb04aa26dde25afc03c53e7`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `890790ca4b88f58b2a68fe466fd0f6026c425f9bc30058c7f028143de0f6eaf4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `add0d035ae5568aac42e7f86295243c79779d51ce6e8b167d07098562f93898f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `8efb6ed8192e9e85b4f4d8d52775a1e93ae1361d82d4c6a1d66417313250190f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `22ef94935d5d4149003f3198a0db821b0af0d70dbd84c44f5bc09c3029952a87`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `4ff24b73b97e19159faf51e3701b64a79d23881e62f54e74ab3f7801c166b59a`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cab70349b0d12b69043c9792b96dd7792c61092034505579514e3e7dd6390e32`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `56b48f2fb52ddf83c124242211ed4201d6026ad36c38a5456fdc1eabc68c4269`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `cded26795e73c5970168c1bfc170e6d312295899292ed7e0e3f644d9ac1760ff`.

I checked the pinned Go 1.27.0 branches for bounded signed order,
including delegation, tiny and huge inputs, forward recurrence and
scaled backward recurrence. The scope is orders -64 through 64 and
all floating input classes. Larger native orders, native amd64 execution
and proof that the continued-fraction guard cannot exhaust remain open.

The independent control passed 55 exact darwin/arm64 inputs, 165
Lean observations across three profiles, and 120 observations against
an exact-rational defining series. The control confirms the pinned
negative infinity on arm64 for order 3 at argument word
0x40033d152e971b40. The separate series gives a finite value near
0.19899990535769083 there.

The author passed 1,467 boundaries and 60,000 seeded observations
with seeds 17, 2718 and 65537. No checked row exhausted the guard. The
author recorded 252 ill-conditioned rows near roots of J0 that
exceeded the independent numerical tolerance and still matched
the pinned Go words exactly.
