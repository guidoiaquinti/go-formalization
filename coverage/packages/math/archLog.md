# Independent review of `math.archLog`

**Verdict: FAITHFUL**

Evidence SHA-256: `df0d00dd411e7a47b47e42d30328f5b7e7a10e87fa169943b2ca4f1f6f284958`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2d0129bb8221ae9d3caf0dd669393d8b9b69eaa2e5e4695775581dcb99fda7c0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `18753e1787167325975e5bfdd03174847075910922330b3ba94a71cae1ef817b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `186ff77646a97200f8df235afecb918c57d496210788a9120ba51cdbae82c998`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ba84410a85121ccc74fd0400b58690afe725f14376f3a05ea7e395bd9a64de39`.

I checked the pinned Go 1.27.0 non-amd64, non-s390x fallback and
the modeled amd64 implementation. The fallback reports false
availability, then panics with the exact string `not implemented`.
The amd64 model preserves the pinned special cases and handles positive
subnormal inputs according to the assembly's effective exponent.

The independent control passed 23 exact native darwin/arm64 stub calls
and 46 Lean observations across both profiles. Its separate logarithm
decomposition checked positive normals and the effective subnormal
input, with fixed tests for zero, infinity and the NaN class.

The author passed 2,164 boundaries and 60,000 seeded rows with seeds
17, 2718 and 65537, including 12,350 effective subnormal operands. The
amd64 model differed from the author's high-precision oracle by at most
one binary64 step. Native amd64 execution and s390x remain open.
