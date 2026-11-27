# Independent review of `path/filepath.HasPrefix`

**Verdict: FAITHFUL**

Evidence SHA-256: `eb0034b1edab8d982c91931f1dc5568cdcc9036810bed07a42a1b37c22eb13d7`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7d99a32bab42a0b3c018fc2bd2bec77c0cb7daa11999be80e26f690f56183ce2`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `204ae26dcea54378f4f609d5b219dd77ca77a98a2634fdfd6b578ec6ca2cf5a0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d4f36e2a1818869d60de8e80ebf38f247f86c38cfa0f46979613332cfae4fe98`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `680dd103c06dd2b734697147741a0032e776e18796788f3a55e6849f7dfb2848`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `de8945f60cd0f509ce3b15cabe813e4d11cadcdbb9d76cb8eba9f0337629cb71`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ee2b450b04649e87152fb1fff0eb46bbade28ff9f5d8ef241d104a5b01f90da4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `61ca976dc34d2491e7cf3ed0cd36821f48e0b58ae8f5760e1c16c12d27fb0938`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `22850836c31dec1d3484dac84bbed840fc7a2f7f9cf69002b0ecb268d68e5e31`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a2f3ab4bee98dd95d26801959204b78793c408e7e33740c09d6a3aeaf90b22e4`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1fda95bac5654bcce767ba49a8c262a49be02a201e85615b4a9a979a1fe2430b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9bc322feca115e6fb96277e915b2224f54f1313df6609a854107c14d5eff1609`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `1cab2b10e2fbb63bfec28f3def00b8d1fa1092161e4603cdae8963075044129d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `919b12a048a4b42b52b60eae983e85864ad00059695337e4550d4a7a984242b9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `2990774ffe90dea4e4c5229354c578e18282d1b58d93845952a633e8a6383656`.

Reviewed 4055983 against Go 1.27.0. The Unix and Plan 9 sources use raw prefix comparison. Windows first performs that same raw check, then compares the prefix after applying the pinned ToLower mapping to both strings. The raw branch can accept a byte prefix that cuts through a valid UTF-8 sequence. In the fallback, distinct malformed bytes can both become the replacement rune, and lowercase conversion can change byte length before the prefix test. No path component boundary or separator equivalence is required.

The independent control checks 14 fixed native and portable Windows expectations and 28 Lean profile results. It distinguishes lowercase-only behavior using the first byte of `É`, simple folding using long s, and raw versus lowercased matching using malformed input and Kelvin sign. Additional cases cover empty prefixes, a prefix ending inside a path component, NUL, and distinct slash/backslash spellings.

The comparison checks 528 boundaries and 60,000 samples per profile using seeds 17, 2718 and 65537. The Windows Go baseline reuses the pinned portable source body with its ordinary strings dependencies. Native Windows and Plan 9 execution remain outside the local test matrix.

Run the reviewer control from the repository root after building the Lean executable.

```sh
GOTOOLCHAIN=go1.27.0 go run difftest/path/filepath/has-prefix/review/main.go
```
