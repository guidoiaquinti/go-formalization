# Independent review of `path/filepath.matchChunk`

**Verdict: FAITHFUL**

Evidence SHA-256: `3132fb34a45ad20ca7d235a69633f3fd2c48585ff4d95002e96750258719c64a`

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `d014227c9d0ad9b11452f9beef549f495f97cac9f20d083b23d988709290d47d`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `49e925921cbe8f83fa9912e486ecc8e7182931aac3124aba63921423163b20eb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `7bfc84a492d40fb8a263a7dcedacf35885a96b21db480b38367bacaa42b07b75`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b7003f29185900cdca79a55958773e657dbe5b38edabdd7ddea98950e44787eb`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ce5a9960d0d3d679201553547d67bc2de721fd2fa1548c5ec6d6d53f8e75336b`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `292221f3731fc1b9657a8e8d5c6902eded5a2f060eb3f3a5de171833e3a47d57`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `9c2ee8b75ec8a5afd3a75d5102fa2d06dea3d7b2756ea51695d7000dad419879`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `a0828e42448323ee0a1954a3cec92871d45922bcf0ae0b100f05022f93214a56`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `ec5691d7b6e19014cea704ad6ecb5ab218f79cf689699bd016327029bb58599f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `06d9855b15f181478caa31a67120bf9657f865060d78b03a8ecdbf26896939ee`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `88451b23a45d1902ddc3e5a71143ea806982aadb27019fecbb8656bbcb0193fd`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `27c93a2b8f7f847f04a95bcc563a96f36b79189674c7c7857f41dc01173211a5`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c09a99cf406799594340efdd31b62e94c34391fedf8372d78d48d11d85543e8e`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c8b9529ac51b0865c7c9c790a6a1553a6f73a1550f1d703d2dd223a1f868e69f`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `649b738eec9f3c8283bb2cbf54a8a59c0f054d89bad4e6092a0d24a668aff3b6`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `b129281ecc78c0d9542ce83da219c035b9249a51f7bd9fd059a9cd6875d49b3c`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `985272eb2d049f20c9233cb83c66ede9a54d57e2182e12632d513728bed9c2c0`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `c68b72f62a14e8a36419cc512ff6ed4f6452e133d6edf3a118728bb4a1f28fe9`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `0149836e3011e7b3393e1f02cdba9ea883a0468c604e83eeca08205dc52b0dd8`.

The review scope matched the preceding evidence, and the full comparison and control suite passed with Go 1.27.0 on darwin/arm64. Previous digest `6f8f04d2959bdb8bc4e124753509abbb3369941d7bd5db2d929253a6bd251123`.

The private Go helper matches raw literal bytes and UTF-8 runes for question marks and ranges. It still validates the remaining pattern after a mismatch or input exhaustion. The Lean parse-then-match construction preserves that priority and returns the exact suffix on success. In the Windows model, question mark rejects the native backslash separator and can consume a forward slash. Non-Windows escapes preserve the next raw byte.

The pinned darwin/arm64 same-package comparison passed 524 boundaries and 60,000 samples from seeds 17, 2718 and 65537. The independent control passed 21 native fixtures and 42 platform model observations. It distinguishes malformed ranges after a mismatch, the exact ErrBadPattern identity, malformed input matching RuneError inside a class, and distinct malformed bytes in literal comparisons. Returning immediately after mismatch failed four fixtures. Treating forward slash as a Windows question-mark separator failed one. Comparing malformed literal bytes by their decoded rune failed one. Native Windows execution remains unverified.

```sh
GOTOOLCHAIN=go1.27.0 GOWORK=off CGO_ENABLED=0 go run ./difftest/path/filepath/match-chunk/review
```

[Model](../../../../lean/GoLean/Path/Filepath/MatchChunk.lean) · [Comparison](../../../../difftest/path/filepath/match-chunk/check.go) · [Independent control](../../../../difftest/path/filepath/match-chunk/review/main.go)
