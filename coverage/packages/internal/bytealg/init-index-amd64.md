# Independent review of `internal/bytealg.init@index_amd64.go:11`

**Verdict: FAITHFUL**

Evidence SHA-256: `aa4a3c925f0d75de21f665f3b284e619677577f5721c8f2a52d5f92efa12a12c`

I checked the pinned build-specific initializer. It assigns MaxLen
31 with AVX2 absent and 63 with AVX2 present. The independent control
checked both model branches against those exact values. The author
exhausted both boolean inputs through the original source body,
with only the actual CPU flag supplied by the harness.

Native amd64 package initialization and CPU feature detection remain
outside the local observations.
